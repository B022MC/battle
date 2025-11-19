import React, { useState } from 'react';
import { View, TouchableOpacity, Text } from 'react-native';
import { PermissionList } from './permission-list';
import { PermissionForm } from './permission-form';
import { PermissionGate } from '@/components/auth/PermissionGate';
import type { Permission } from '@/services/basic/permission';

export function PermissionsView() {
  const [showForm, setShowForm] = useState(false);
  const [editingPermission, setEditingPermission] = useState<Permission | null>(
    null
  );
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setEditingPermission(null);
    setShowForm(true);
  };

  const handleEdit = (permission: Permission) => {
    setEditingPermission(permission);
    setShowForm(true);
  };

  const handleFormSuccess = () => {
    setRefreshKey((prev) => prev + 1);
  };

  const handleFormClose = () => {
    setShowForm(false);
    setEditingPermission(null);
  };

  return (
    <View className="flex-1 bg-white">
      {/* 头部 */}
      <View className="px-4 py-3 bg-white border-b border-gray-200 flex-row items-center justify-between">
        <Text className="text-lg font-semibold">权限管理</Text>
        <PermissionGate anyOf={['permission:create']}>
          <TouchableOpacity
            onPress={handleCreate}
            className="px-4 py-2 bg-blue-500 rounded-lg"
          >
            <Text className="text-white font-medium">创建权限</Text>
          </TouchableOpacity>
        </PermissionGate>
      </View>

      {/* 权限列表 */}
      <PermissionList
        key={refreshKey}
        onEdit={handleEdit}
        onRefresh={handleFormSuccess}
      />

      {/* 创建/编辑表单 */}
      <PermissionForm
        visible={showForm}
        onClose={handleFormClose}
        onSuccess={handleFormSuccess}
        permission={editingPermission}
      />
    </View>
  );
}

