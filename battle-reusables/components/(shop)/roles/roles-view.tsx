import React, { useState } from 'react';
import { View, TouchableOpacity, Text } from 'react-native';
import { RoleList } from './role-list';
import { RoleForm } from './role-form';
import { AssignPermissionsModal } from './assign-permissions-modal';
import { AssignMenusModal } from './assign-menus-modal';
import { PermissionGate } from '@/components/auth/PermissionGate';
import type { Role } from '@/services/basic/role';

export function RolesView() {
  const [showForm, setShowForm] = useState(false);
  const [showPermissions, setShowPermissions] = useState(false);
  const [showMenus, setShowMenus] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setEditingRole(null);
    setShowForm(true);
  };

  const handleEdit = (role: Role) => {
    setEditingRole(role);
    setShowForm(true);
  };

  const handleAssignPermissions = (role: Role) => {
    setSelectedRole(role);
    setShowPermissions(true);
  };

  const handleAssignMenus = (role: Role) => {
    setSelectedRole(role);
    setShowMenus(true);
  };

  const handleSuccess = () => {
    setRefreshKey((prev) => prev + 1);
  };

  const handleFormClose = () => {
    setShowForm(false);
    setEditingRole(null);
  };

  const handlePermissionsClose = () => {
    setShowPermissions(false);
    setSelectedRole(null);
  };

  const handleMenusClose = () => {
    setShowMenus(false);
    setSelectedRole(null);
  };

  return (
    <View className="flex-1 bg-white">
      {/* 头部 */}
      <View className="px-4 py-3 bg-white border-b border-gray-200 flex-row items-center justify-between">
        <Text className="text-lg font-semibold">角色管理</Text>
        <PermissionGate anyOf={['role:create']}>
          <TouchableOpacity
            onPress={handleCreate}
            className="px-4 py-2 bg-blue-500 rounded-lg"
          >
            <Text className="text-white font-medium">创建角色</Text>
          </TouchableOpacity>
        </PermissionGate>
      </View>

      {/* 角色列表 */}
      <RoleList
        key={refreshKey}
        onEdit={handleEdit}
        onAssignPermissions={handleAssignPermissions}
        onAssignMenus={handleAssignMenus}
        onRefresh={handleSuccess}
      />

      {/* 创建/编辑表单 */}
      <RoleForm
        visible={showForm}
        onClose={handleFormClose}
        onSuccess={handleSuccess}
        role={editingRole}
      />

      {/* 分配权限弹窗 */}
      <AssignPermissionsModal
        visible={showPermissions}
        onClose={handlePermissionsClose}
        onSuccess={handleSuccess}
        role={selectedRole}
      />

      {/* 分配菜单弹窗 */}
      <AssignMenusModal
        visible={showMenus}
        onClose={handleMenusClose}
        onSuccess={handleSuccess}
        role={selectedRole}
      />
    </View>
  );
}

