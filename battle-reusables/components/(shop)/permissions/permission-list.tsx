import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import {
  getAllPermissions,
  deletePermission,
  type Permission,
} from '@/services/basic/permission';
import { toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';
import { PermissionGate } from '@/components/auth/PermissionGate';

interface PermissionListProps {
  onEdit: (permission: Permission) => void;
  onRefresh?: () => void;
}

export function PermissionList({ onEdit, onRefresh }: PermissionListProps) {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  // 权限分类
  const categories = [
    { key: 'all', label: '全部' },
    { key: 'stats', label: '统计' },
    { key: 'fund', label: '资金' },
    { key: 'shop', label: '店铺' },
    { key: 'game', label: '游戏' },
    { key: 'system', label: '系统' },
  ];

  const loadPermissions = async () => {
    setLoading(true);
    try {
      const res = await getAllPermissions();
      if (res.code === 0 && res.data) {
        setPermissions(res.data);
      } else {
        setPermissions([]);
        toast.info('权限列表暂无数据');
      }
    } catch (error: any) {
      if (error.message?.includes('404')) {
        // API未实现，显示提示
        console.warn('权限API返回404，需要后端实现 /basic/permission/listAll 接口');
        setPermissions([]);
        toast.warning('权限管理功能暂未开放', '请联系后端实现 /basic/permission/listAll 接口');
      } else {
        toast.error('加载权限列表失败');
        console.error('Load permissions error:', error);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPermissions();
  }, []);

  const handleDelete = async (permission: Permission) => {
    toast.confirm({
      title: '确认删除',
      description: `确定要删除权限"${permission.name}"吗？删除后将无法恢复。`,
      type: 'error',
      confirmText: '删除',
      cancelText: '取消',
      confirmVariant: 'destructive',
      onConfirm: async () => {
        try {
          const res = await deletePermission(permission.id);
          if (res.code === 0) {
            showSuccessBubble('删除成功', `权限"${permission.name}"已删除`);
            loadPermissions();
            onRefresh?.();
          } else {
            toast.error(res.msg || '删除失败');
          }
        } catch (error) {
          toast.error('删除失败');
          console.error('Delete permission error:', error);
        }
      },
    });
  };

  // 过滤权限
  const filteredPermissions =
    selectedCategory === 'all'
      ? permissions
      : permissions.filter((p) => p.category === selectedCategory);

  return (
    <View className="flex-1 bg-gray-50">
      {/* 分类筛选 */}
      <View className="px-4 py-3 bg-white">
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          <View className="flex-row gap-2">
            {categories.map((cat) => (
              <TouchableOpacity
                key={cat.key}
                onPress={() => setSelectedCategory(cat.key)}
                className={`px-4 py-2 rounded-full ${
                  selectedCategory === cat.key
                    ? 'bg-blue-500'
                    : 'bg-gray-100'
                }`}
              >
                <Text
                  className={`text-sm font-medium ${
                    selectedCategory === cat.key
                      ? 'text-white'
                      : 'text-gray-700'
                  }`}
                >
                  {cat.label}
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </ScrollView>
      </View>

      {/* 权限列表 */}
      <ScrollView className="flex-1 px-4">
        {loading ? (
          <View className="py-8 items-center">
            <Text className="text-gray-500">加载中...</Text>
          </View>
        ) : filteredPermissions.length === 0 ? (
          <View className="py-8 items-center">
            <Text className="text-gray-500">暂无权限数据</Text>
          </View>
        ) : (
          <View className="py-4 gap-3">
            {filteredPermissions.map((permission) => (
              <View
                key={permission.id}
                className="bg-white rounded-lg p-4 shadow-sm"
              >
                <View className="flex-row items-start justify-between mb-2">
                  <View className="flex-1">
                    <Text className="text-base font-semibold text-gray-900">
                      {permission.name}
                    </Text>
                    <Text className="text-xs text-gray-500 mt-1">
                      {permission.code}
                    </Text>
                  </View>
                  <View className="flex-row gap-2">
                    <PermissionGate anyOf={['permission:update']}>
                      <TouchableOpacity
                        onPress={() => onEdit(permission)}
                        className="px-3 py-1.5 bg-blue-50 rounded"
                      >
                        <Text className="text-xs text-blue-600 font-medium">
                          编辑
                        </Text>
                      </TouchableOpacity>
                    </PermissionGate>
                    <PermissionGate anyOf={['permission:delete']}>
                      <TouchableOpacity
                        onPress={() => handleDelete(permission)}
                        className="px-3 py-1.5 bg-red-50 rounded"
                      >
                        <Text className="text-xs text-red-600 font-medium">
                          删除
                        </Text>
                      </TouchableOpacity>
                    </PermissionGate>
                  </View>
                </View>
                {permission.description && (
                  <Text className="text-sm text-gray-600 mt-1">
                    {permission.description}
                  </Text>
                )}
                <View className="flex-row items-center mt-2">
                  <View className="px-2 py-1 bg-gray-100 rounded">
                    <Text className="text-xs text-gray-600">
                      {
                        categories.find((c) => c.key === permission.category)
                          ?.label
                      }
                    </Text>
                  </View>
                </View>
              </View>
            ))}
          </View>
        )}
      </ScrollView>
    </View>
  );
}

