import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  Modal,
  ActivityIndicator,
} from 'react-native';
import {
  getAllPermissions,
  getRolePermissions,
  assignPermissionsToRole,
  type Permission,
} from '@/services/basic/permission';
import type { Role } from '@/services/basic/role';
import { showToast } from '@/utils/toast';

interface AssignPermissionsModalProps {
  visible: boolean;
  onClose: () => void;
  onSuccess: () => void;
  role: Role | null;
}

export function AssignPermissionsModal({
  visible,
  onClose,
  onSuccess,
  role,
}: AssignPermissionsModalProps) {
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  // 按分类分组权限
  const groupedPermissions = allPermissions.reduce((acc, perm) => {
    if (!acc[perm.category]) {
      acc[perm.category] = [];
    }
    acc[perm.category].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  const categoryLabels: Record<string, string> = {
    stats: '统计',
    fund: '资金',
    shop: '店铺',
    game: '游戏',
    system: '系统',
  };

  useEffect(() => {
    if (visible && role) {
      loadData();
    }
  }, [visible, role]);

  const loadData = async () => {
    if (!role) return;

    setLoading(true);
    try {
      // 加载所有权限
      const allRes = await getAllPermissions();
      if (allRes.success && allRes.data) {
        setAllPermissions(allRes.data);
      }

      // 加载角色已有权限
      const roleRes = await getRolePermissions(role.id);
      if (roleRes.success && roleRes.data) {
        const ids = new Set(roleRes.data.map((p) => p.id));
        setSelectedIds(ids);
      }
    } catch (error) {
      showToast('加载数据失败', 'error');
      console.error('Load permissions error:', error);
    } finally {
      setLoading(false);
    }
  };

  const togglePermission = (permId: number) => {
    const newSet = new Set(selectedIds);
    if (newSet.has(permId)) {
      newSet.delete(permId);
    } else {
      newSet.add(permId);
    }
    setSelectedIds(newSet);
  };

  const toggleCategory = (category: string) => {
    const categoryPerms = groupedPermissions[category] || [];
    const categoryIds = categoryPerms.map((p) => p.id);
    const allSelected = categoryIds.every((id) => selectedIds.has(id));

    const newSet = new Set(selectedIds);
    if (allSelected) {
      // 取消全选
      categoryIds.forEach((id) => newSet.delete(id));
    } else {
      // 全选
      categoryIds.forEach((id) => newSet.add(id));
    }
    setSelectedIds(newSet);
  };

  const handleSubmit = async () => {
    if (!role) return;

    setSaving(true);
    try {
      const res = await assignPermissionsToRole({
        role_id: role.id,
        permission_ids: Array.from(selectedIds),
      });
      if (res.success) {
        showToast('分配权限成功', 'success');
        onSuccess();
        onClose();
      } else {
        showToast(res.message || '分配权限失败', 'error');
      }
    } catch (error) {
      showToast('分配权限失败', 'error');
      console.error('Assign permissions error:', error);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="slide"
      onRequestClose={onClose}
    >
      <View className="flex-1 bg-black/50 justify-end">
        <View className="bg-white rounded-t-2xl" style={{ maxHeight: '90%' }}>
          {/* 头部 */}
          <View className="flex-row items-center justify-between px-4 py-4 border-b border-gray-200">
            <View>
              <Text className="text-lg font-semibold">分配权限</Text>
              {role && (
                <Text className="text-sm text-gray-500 mt-1">
                  {role.name} (已选 {selectedIds.size} 个)
                </Text>
              )}
            </View>
            <TouchableOpacity onPress={onClose}>
              <Text className="text-gray-500 text-base">✕</Text>
            </TouchableOpacity>
          </View>

          {/* 权限列表 */}
          {loading ? (
            <View className="py-12 items-center">
              <ActivityIndicator size="large" color="#3b82f6" />
              <Text className="text-gray-500 mt-2">加载中...</Text>
            </View>
          ) : (
            <ScrollView className="flex-1 px-4 py-4">
              {Object.entries(groupedPermissions).map(([category, perms]) => (
                <View key={category} className="mb-4">
                  {/* 分类标题 */}
                  <TouchableOpacity
                    onPress={() => toggleCategory(category)}
                    className="flex-row items-center justify-between py-2 mb-2"
                  >
                    <Text className="text-base font-semibold text-gray-900">
                      {categoryLabels[category] || category}
                    </Text>
                    <View className="flex-row items-center">
                      <Text className="text-xs text-gray-500 mr-2">
                        {perms.filter((p) => selectedIds.has(p.id)).length}/
                        {perms.length}
                      </Text>
                      <View
                        className={`w-5 h-5 rounded border-2 items-center justify-center ${
                          perms.every((p) => selectedIds.has(p.id))
                            ? 'bg-blue-500 border-blue-500'
                            : 'border-gray-300'
                        }`}
                      >
                        {perms.every((p) => selectedIds.has(p.id)) && (
                          <Text className="text-white text-xs">✓</Text>
                        )}
                      </View>
                    </View>
                  </TouchableOpacity>

                  {/* 权限列表 */}
                  <View className="gap-2">
                    {perms.map((perm) => (
                      <TouchableOpacity
                        key={perm.id}
                        onPress={() => togglePermission(perm.id)}
                        className="flex-row items-center p-3 bg-gray-50 rounded-lg"
                      >
                        <View
                          className={`w-5 h-5 rounded border-2 items-center justify-center mr-3 ${
                            selectedIds.has(perm.id)
                              ? 'bg-blue-500 border-blue-500'
                              : 'border-gray-300'
                          }`}
                        >
                          {selectedIds.has(perm.id) && (
                            <Text className="text-white text-xs">✓</Text>
                          )}
                        </View>
                        <View className="flex-1">
                          <Text className="text-sm font-medium text-gray-900">
                            {perm.name}
                          </Text>
                          <Text className="text-xs text-gray-500 mt-0.5">
                            {perm.code}
                          </Text>
                        </View>
                      </TouchableOpacity>
                    ))}
                  </View>
                </View>
              ))}
            </ScrollView>
          )}

          {/* 底部按钮 */}
          <View className="px-4 py-4 border-t border-gray-200">
            <View className="flex-row gap-3">
              <TouchableOpacity
                onPress={onClose}
                className="flex-1 py-3 bg-gray-100 rounded-lg"
                disabled={saving}
              >
                <Text className="text-center text-gray-700 font-medium">
                  取消
                </Text>
              </TouchableOpacity>
              <TouchableOpacity
                onPress={handleSubmit}
                className={`flex-1 py-3 rounded-lg ${
                  saving ? 'bg-blue-300' : 'bg-blue-500'
                }`}
                disabled={saving || loading}
              >
                <Text className="text-center text-white font-medium">
                  {saving ? '保存中...' : '保存'}
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
}

