import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { getMenuTree, deleteMenu, type Menu } from '@/services/basic/menu';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { toast } from '@/utils/toast';

interface MenuListProps {
  onEdit: (menu: Menu) => void;
  onRefresh: () => void;
}

export function MenuList({ onEdit, onRefresh }: MenuListProps) {
  const [loading, setLoading] = useState(false);
  const [menus, setMenus] = useState<Menu[]>([]);
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());

  useEffect(() => {
    loadMenus();
  }, []);

  // 将嵌套的菜单树扁平化
  const flattenMenuTree = (menus: Menu[], result: Menu[] = []): Menu[] => {
    menus.forEach(menu => {
      result.push(menu);
      if (menu.children && menu.children.length > 0) {
        flattenMenuTree(menu.children, result);
      }
    });
    return result;
  };

  const loadMenus = async () => {
    try {
      setLoading(true);
      // 不传参数，避免传递undefined
      const response = await getMenuTree();
      if (response.code === 0 && response.data) {
        // 将嵌套的树形结构扁平化
        const flatMenus = flattenMenuTree(Array.isArray(response.data) ? response.data : []);
        setMenus(flatMenus);
        console.log('加载菜单成功，共', flatMenus.length, '个菜单');
        
        // 默认展开有子菜单的项（如店铺菜单）
        const menuWithChildren = flatMenus.filter(m => {
          return flatMenus.some(child => child.parent_id === m.id);
        });
        setExpandedIds(new Set(menuWithChildren.map(m => m.id)));
      } else {
        setMenus([]);
        if (response.data === null) {
          toast.info('菜单列表暂无数据', '请先添加菜单');
        }
      }
    } catch (error: any) {
      toast.error('加载菜单失败');
      console.error('加载菜单失败:', error);
      setMenus([]);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = (menu: Menu) => {
    Alert.alert(
      '确认删除',
      `确定要删除菜单"${menu.title}"吗？`,
      [
        { text: '取消', style: 'cancel' },
        {
          text: '删除',
          style: 'destructive',
          onPress: async () => {
            try {
              const res = await deleteMenu(menu.id);
              if (res.code === 0) {
                toast.success('删除成功');
                onRefresh();
                loadMenus();
              } else {
                toast.error(res.msg || '删除失败');
              }
            } catch (error: any) {
              toast.error(error.message || '删除失败');
            }
          },
        },
      ]
    );
  };

  const toggleExpand = (id: number) => {
    setExpandedIds((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const renderMenu = (menu: Menu, level: number = 0) => {
    const hasChildren = menus.some((m) => m.parent_id === menu.id);
    const isExpanded = expandedIds.has(menu.id);
    const children = menus.filter((m) => m.parent_id === menu.id);

    return (
      <View key={menu.id}>
        {/* 菜单项 */}
        <View
          className="px-4 py-3 border-b border-gray-100"
          style={{ paddingLeft: 16 + level * 20 }}
        >
          <View className="flex-row items-center justify-between">
            <View className="flex-1">
              <View className="flex-row items-center">
                {/* 展开/折叠按钮 */}
                {hasChildren && (
                  <TouchableOpacity
                    onPress={() => toggleExpand(menu.id)}
                    className="mr-2"
                  >
                    <Text className="text-gray-500 text-base">
                      {isExpanded ? '▼' : '▶'}
                    </Text>
                  </TouchableOpacity>
                )}
                {!hasChildren && <View className="w-6" />}

                <View className="flex-1">
                  <Text className="text-base font-medium text-gray-900">
                    {menu.title}
                  </Text>
                  <Text className="text-sm text-gray-500 mt-1">
                    {menu.name} • {menu.path}
                  </Text>
                  {menu.auths && (
                    <Text className="text-xs text-gray-400 mt-1">
                      权限: {menu.auths}
                    </Text>
                  )}
                </View>
              </View>
            </View>

            {/* 操作按钮 */}
            <View className="flex-row gap-2">
              <PermissionGate anyOf={['menu:update']}>
                <TouchableOpacity
                  onPress={() => onEdit(menu)}
                  className="px-3 py-1 bg-blue-50 rounded"
                >
                  <Text className="text-blue-600 text-sm">编辑</Text>
                </TouchableOpacity>
              </PermissionGate>

              <PermissionGate anyOf={['menu:delete']}>
                <TouchableOpacity
                  onPress={() => handleDelete(menu)}
                  className="px-3 py-1 bg-red-50 rounded"
                >
                  <Text className="text-red-600 text-sm">删除</Text>
                </TouchableOpacity>
              </PermissionGate>
            </View>
          </View>
        </View>

        {/* 子菜单 */}
        {isExpanded &&
          children.map((child) => renderMenu(child, level + 1))}
      </View>
    );
  };

  if (loading) {
    return (
      <View className="flex-1 items-center justify-center">
        <ActivityIndicator size="large" color="#3B82F6" />
      </View>
    );
  }

  // 只渲染顶层菜单（parent_id 为 -1, 0 或 null）
  const topMenus = menus.filter((m) => m.parent_id === -1 || m.parent_id === 0 || !m.parent_id);

  return (
    <FlatList
      data={topMenus}
      keyExtractor={(item) => item.id.toString()}
      renderItem={({ item }) => renderMenu(item)}
      ListEmptyComponent={
        <View className="p-8 items-center">
          <Text className="text-gray-400">暂无菜单数据</Text>
        </View>
      }
      refreshing={loading}
      onRefresh={loadMenus}
    />
  );
}


