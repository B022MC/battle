import React, { useState } from 'react';
import { View, TouchableOpacity, Text } from 'react-native';
import { MenuList } from './menu-list';
import { MenuForm } from './menu-form';
import { PermissionGate } from '@/components/auth/PermissionGate';
import type { Menu } from '@/services/basic/menu';

export function MenusView() {
  const [showForm, setShowForm] = useState(false);
  const [editingMenu, setEditingMenu] = useState<Menu | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setEditingMenu(null);
    setShowForm(true);
  };

  const handleEdit = (menu: Menu) => {
    setEditingMenu(menu);
    setShowForm(true);
  };

  const handleFormSuccess = () => {
    setRefreshKey((prev) => prev + 1);
  };

  const handleFormClose = () => {
    setShowForm(false);
    setEditingMenu(null);
  };

  return (
    <View className="flex-1 bg-white">
      {/* 头部 */}
      <View className="px-4 py-3 bg-white border-b border-gray-200 flex-row items-center justify-between">
        <Text className="text-lg font-semibold">菜单管理</Text>
        <PermissionGate anyOf={['menu:create']}>
          <TouchableOpacity
            onPress={handleCreate}
            className="px-4 py-2 bg-blue-500 rounded-lg"
          >
            <Text className="text-white font-medium">创建菜单</Text>
          </TouchableOpacity>
        </PermissionGate>
      </View>

      {/* 菜单列表 */}
      <MenuList
        key={refreshKey}
        onEdit={handleEdit}
        onRefresh={handleFormSuccess}
      />

      {/* 创建/编辑表单 */}
      <MenuForm
        visible={showForm}
        onClose={handleFormClose}
        onSuccess={handleFormSuccess}
        menu={editingMenu}
      />
    </View>
  );
}



