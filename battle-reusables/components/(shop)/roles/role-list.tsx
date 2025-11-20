import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import { getRoleList, deleteRole, type Role } from '@/services/basic/role';
import { showToast, toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';
import { PermissionGate } from '@/components/auth/PermissionGate';

interface RoleListProps {
  onEdit: (role: Role) => void;
  onAssignPermissions: (role: Role) => void;
  onAssignMenus: (role: Role) => void;
  onRefresh?: () => void;
}

export function RoleList({
  onEdit,
  onAssignPermissions,
  onAssignMenus,
  onRefresh,
}: RoleListProps) {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);

  const loadRoles = async () => {
    setLoading(true);
    try {
      const res = await getRoleList({ page_size: 100 });
      if (res.code === 0 && res.data) {
        setRoles(res.data.list);
      }
    } catch (error) {
      showToast('加载角色列表失败', 'error');
      console.error('Load roles error:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadRoles();
  }, []);

  const handleDelete = async (role: Role) => {
    // 不允许删除系统预定义角色
    if (role.id === 1 || role.id === 2 || role.id === 3) {
      showToast('不能删除系统预定义角色', 'error');
      return;
    }

    toast.confirm({
      title: '确认删除',
      description: `确定要删除角色"${role.name}"吗？删除后将无法恢复。`,
      type: 'error',
      confirmText: '删除',
      cancelText: '取消',
      confirmVariant: 'destructive',
      onConfirm: async () => {
        try {
          const res = await deleteRole(role.id);
          if (res.code === 0) {
            showSuccessBubble('删除成功', `角色"${role.name}"已删除`);
            loadRoles();
            onRefresh?.();
          } else {
            showToast(res.msg || '删除失败', 'error');
          }
        } catch (error) {
          showToast('删除失败', 'error');
          console.error('Delete role error:', error);
        }
      },
    });
  };

  const getRoleBadgeColor = (roleId: number) => {
    switch (roleId) {
      case 1:
        return 'bg-red-100 text-red-700';
      case 2:
        return 'bg-blue-100 text-blue-700';
      case 3:
        return 'bg-green-100 text-green-700';
      default:
        return 'bg-gray-100 text-gray-700';
    }
  };

  return (
    <View className="flex-1 bg-gray-50">
      <ScrollView className="flex-1 px-4">
        {loading ? (
          <View className="py-8 items-center">
            <Text className="text-gray-500">加载中...</Text>
          </View>
        ) : roles.length === 0 ? (
          <View className="py-8 items-center">
            <Text className="text-gray-500">暂无角色数据</Text>
          </View>
        ) : (
          <View className="py-4 gap-3">
            {roles.map((role) => (
              <View
                key={role.id}
                className="bg-white rounded-lg p-4 shadow-sm"
              >
                <View className="flex-row items-start justify-between mb-3">
                  <View className="flex-1">
                    <View className="flex-row items-center gap-2 mb-1">
                      <Text className="text-base font-semibold text-gray-900">
                        {role.name}
                      </Text>
                      <View
                        className={`px-2 py-0.5 rounded ${getRoleBadgeColor(
                          role.id
                        )}`}
                      >
                        <Text className="text-xs font-medium">
                          {role.code}
                        </Text>
                      </View>
                      {!role.enable && (
                        <View className="px-2 py-0.5 bg-gray-100 rounded">
                          <Text className="text-xs text-gray-600">已禁用</Text>
                        </View>
                      )}
                    </View>
                    {role.remark && (
                      <Text className="text-sm text-gray-600 mt-1">
                        {role.remark}
                      </Text>
                    )}
                  </View>
                </View>

                {/* 操作按钮 */}
                <View className="flex-row flex-wrap gap-2">
                  <PermissionGate anyOf={['role:update']}>
                    <TouchableOpacity
                      onPress={() => onEdit(role)}
                      className="px-3 py-1.5 bg-blue-50 rounded flex-row items-center"
                    >
                      <Text className="text-xs text-blue-600 font-medium">
                        编辑
                      </Text>
                    </TouchableOpacity>
                  </PermissionGate>

                  <PermissionGate anyOf={['permission:assign']}>
                    <TouchableOpacity
                      onPress={() => onAssignPermissions(role)}
                      className="px-3 py-1.5 bg-purple-50 rounded flex-row items-center"
                    >
                      <Text className="text-xs text-purple-600 font-medium">
                        分配权限
                      </Text>
                    </TouchableOpacity>
                  </PermissionGate>

                  <PermissionGate anyOf={['role:update']}>
                    <TouchableOpacity
                      onPress={() => onAssignMenus(role)}
                      className="px-3 py-1.5 bg-green-50 rounded flex-row items-center"
                    >
                      <Text className="text-xs text-green-600 font-medium">
                        分配菜单
                      </Text>
                    </TouchableOpacity>
                  </PermissionGate>

                  <PermissionGate anyOf={['role:delete']}>
                    {role.id > 3 && (
                      <TouchableOpacity
                        onPress={() => handleDelete(role)}
                        className="px-3 py-1.5 bg-red-50 rounded flex-row items-center"
                      >
                        <Text className="text-xs text-red-600 font-medium">
                          删除
                        </Text>
                      </TouchableOpacity>
                    )}
                  </PermissionGate>
                </View>
              </View>
            ))}
          </View>
        )}
      </ScrollView>
    </View>
  );
}

