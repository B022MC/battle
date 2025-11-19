import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  Modal,
  ActivityIndicator,
} from 'react-native';
import { getRoleMenus, assignMenusToRole } from '@/services/basic/role';
import { getMenuTree } from '@/services/basic/menu';
import type { Role } from '@/services/basic/role';
import { showToast } from '@/utils/toast';

interface Menu {
  id: number;
  parent_id: number;
  menu_type: number;
  title: string;
  name: string;
  path: string;
}

interface AssignMenusModalProps {
  visible: boolean;
  onClose: () => void;
  onSuccess: () => void;
  role: Role | null;
}

export function AssignMenusModal({
  visible,
  onClose,
  onSuccess,
  role,
}: AssignMenusModalProps) {
  const [allMenus, setAllMenus] = useState<Menu[]>([]);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  // 分组菜单：一级和二级
  const primaryMenus = allMenus.filter((m) => m.menu_type === 1);
  const secondaryMenus = allMenus.filter((m) => m.menu_type === 2);

  useEffect(() => {
    if (visible && role) {
      loadData();
    }
  }, [visible, role]);

  const flattenMenus = (menus: any[]): Menu[] => {
    const result: Menu[] = [];
    const flatten = (items: any[]) => {
      items.forEach((item) => {
        result.push({
          id: item.id,
          parent_id: item.parent_id,
          menu_type: item.menu_type,
          title: item.title,
          name: item.name,
          path: item.path,
        });
        if (item.children && item.children.length > 0) {
          flatten(item.children);
        }
      });
    };
    flatten(menus);
    return result;
  };

  const loadData = async () => {
    if (!role) return;

    setLoading(true);
    try {
      // 加载所有菜单（树形结构）
      const allRes: any = await getMenuTree();
      if (allRes.code === 0 && allRes.data) {
        // 将树形结构扁平化
        const flatMenus = flattenMenus(allRes.data);
        setAllMenus(flatMenus);
      }

      // 加载角色已有菜单
      const roleRes = await getRoleMenus(role.id);
      if (roleRes.code === 0 && roleRes.data) {
        setSelectedIds(new Set(roleRes.data.menu_ids));
      }
    } catch (error) {
      showToast('加载数据失败', 'error');
      console.error('Load menus error:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleMenu = (menuId: number) => {
    const newSet = new Set(selectedIds);
    if (newSet.has(menuId)) {
      newSet.delete(menuId);
    } else {
      newSet.add(menuId);
    }
    setSelectedIds(newSet);
  };

  const toggleParentMenu = (parentId: number) => {
    const children = secondaryMenus.filter((m) => m.parent_id === parentId);
    const allChildrenSelected = children.every((c) => selectedIds.has(c.id));

    const newSet = new Set(selectedIds);
    
    if (allChildrenSelected) {
      // 取消父菜单和所有子菜单
      newSet.delete(parentId);
      children.forEach((c) => newSet.delete(c.id));
    } else {
      // 选中父菜单和所有子菜单
      newSet.add(parentId);
      children.forEach((c) => newSet.add(c.id));
    }
    setSelectedIds(newSet);
  };

  const handleSubmit = async () => {
    if (!role) return;

    setSaving(true);
    try {
      const res = await assignMenusToRole({
        role_id: role.id,
        menu_ids: Array.from(selectedIds),
      });
      if (res.code === 0) {
        showToast('分配菜单成功', 'success');
        onSuccess();
        onClose();
      } else {
        showToast(res.msg || '分配菜单失败', 'error');
      }
    } catch (error) {
      showToast('分配菜单失败', 'error');
      console.error('Assign menus error:', error);
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
              <Text className="text-lg font-semibold">分配菜单</Text>
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

          {/* 菜单列表 */}
          {loading ? (
            <View className="py-12 items-center">
              <ActivityIndicator size="large" color="#3b82f6" />
              <Text className="text-gray-500 mt-2">加载中...</Text>
            </View>
          ) : (
            <ScrollView className="flex-1 px-4 py-4">
              {primaryMenus.map((menu) => {
                const children = secondaryMenus.filter(
                  (m) => m.parent_id === menu.id
                );
                const hasChildren = children.length > 0;

                return (
                  <View key={menu.id} className="mb-4">
                    {/* 一级菜单 */}
                    <TouchableOpacity
                      onPress={() =>
                        hasChildren
                          ? toggleParentMenu(menu.id)
                          : toggleMenu(menu.id)
                      }
                      className="flex-row items-center p-3 bg-blue-50 rounded-lg mb-2"
                    >
                      <View
                        className={`w-5 h-5 rounded border-2 items-center justify-center mr-3 ${
                          selectedIds.has(menu.id)
                            ? 'bg-blue-500 border-blue-500'
                            : 'border-gray-300'
                        }`}
                      >
                        {selectedIds.has(menu.id) && (
                          <Text className="text-white text-xs">✓</Text>
                        )}
                      </View>
                      <View className="flex-1">
                        <Text className="text-sm font-semibold text-gray-900">
                          {menu.title}
                        </Text>
                        <Text className="text-xs text-gray-500 mt-0.5">
                          {menu.path}
                        </Text>
                      </View>
                      {hasChildren && (
                        <Text className="text-xs text-gray-500">
                          {children.filter((c) => selectedIds.has(c.id)).length}
                          /{children.length}
                        </Text>
                      )}
                    </TouchableOpacity>

                    {/* 二级菜单 */}
                    {hasChildren && (
                      <View className="ml-6 gap-2">
                        {children.map((child) => (
                          <TouchableOpacity
                            key={child.id}
                            onPress={() => toggleMenu(child.id)}
                            className="flex-row items-center p-3 bg-gray-50 rounded-lg"
                          >
                            <View
                              className={`w-5 h-5 rounded border-2 items-center justify-center mr-3 ${
                                selectedIds.has(child.id)
                                  ? 'bg-blue-500 border-blue-500'
                                  : 'border-gray-300'
                              }`}
                            >
                              {selectedIds.has(child.id) && (
                                <Text className="text-white text-xs">✓</Text>
                              )}
                            </View>
                            <View className="flex-1">
                              <Text className="text-sm font-medium text-gray-900">
                                {child.title}
                              </Text>
                              <Text className="text-xs text-gray-500 mt-0.5">
                                {child.path}
                              </Text>
                            </View>
                          </TouchableOpacity>
                        ))}
                      </View>
                    )}
                  </View>
                );
              })}
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

