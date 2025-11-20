import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  Modal,
  TextInput,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Switch,
} from 'react-native';
import {
  createMenu,
  updateMenu,
  getMenuTree,
  type Menu,
  type CreateMenuRequest,
  type UpdateMenuRequest,
} from '@/services/basic/menu';
import { showToast, toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';

interface MenuFormProps {
  visible: boolean;
  onClose: () => void;
  onSuccess: () => void;
  menu: Menu | null;
}

export function MenuForm({ visible, onClose, onSuccess, menu }: MenuFormProps) {
  const [loading, setLoading] = useState(false);
  const [parentMenus, setParentMenus] = useState<Menu[]>([]);
  
  // 表单字段
  const [title, setTitle] = useState('');
  const [name, setName] = useState('');
  const [path, setPath] = useState('');
  const [component, setComponent] = useState('');
  const [parentId, setParentId] = useState(-1);
  const [menuType, setMenuType] = useState(1);
  const [rank, setRank] = useState('');
  const [icon, setIcon] = useState('');
  const [auths, setAuths] = useState('');
  const [showLink, setShowLink] = useState(true);
  const [showParent, setShowParent] = useState(true);

  useEffect(() => {
    if (visible) {
      loadParentMenus();
      if (menu) {
        // 编辑模式
        setTitle(menu.title);
        setName(menu.name);
        setPath(menu.path);
        setComponent(menu.component);
        setParentId(menu.parent_id);
        setMenuType(menu.menu_type);
        setRank(menu.rank || '');
        setIcon(menu.icon || '');
        setAuths(menu.auths || '');
        setShowLink(menu.show_link);
        setShowParent(menu.show_parent);
      } else {
        // 创建模式
        resetForm();
      }
    }
  }, [visible, menu]);

  const flattenMenus = (menus: any[]): Menu[] => {
    const result: Menu[] = [];
    const flatten = (items: any[]) => {
      items.forEach((item) => {
        result.push(item);
        if (item.children && item.children.length > 0) {
          flatten(item.children);
        }
      });
    };
    flatten(menus);
    return result;
  };

  const loadParentMenus = async () => {
    try {
      const response: any = await getMenuTree();
      if (response.code === 0 && response.data) {
        // 将树形结构扁平化
        const flatMenus = flattenMenus(response.data);
        setParentMenus(flatMenus);
      }
    } catch (error) {
      console.error('加载父级菜单失败:', error);
    }
  };

  const resetForm = () => {
    setTitle('');
    setName('');
    setPath('');
    setComponent('');
    setParentId(-1);
    setMenuType(1);
    setRank('');
    setIcon('');
    setAuths('');
    setShowLink(true);
    setShowParent(true);
  };

  const handleSubmit = async () => {
    if (!title.trim()) {
      showToast('请输入菜单标题', 'error');
      return;
    }
    if (!name.trim()) {
      showToast('请输入菜单名称', 'error');
      return;
    }
    if (!path.trim()) {
      showToast('请输入菜单路径', 'error');
      return;
    }
    if (!component.trim()) {
      showToast('请输入组件路径', 'error');
      return;
    }

    // 二次确认
    toast.confirm({
      title: menu ? '确认更新' : '确认创建',
      description: menu 
        ? `确定要更新菜单"${menu.title}"吗？`
        : `确定要创建菜单"${title}"吗？`,
      type: 'warning',
      confirmText: menu ? '更新' : '创建',
      cancelText: '取消',
      onConfirm: async () => {
        try {
          setLoading(true);
          
          if (menu) {
            // 更新
            const data: UpdateMenuRequest = {
              id: menu.id,
              title,
              name,
              path,
              component,
              parent_id: parentId,
              menu_type: menuType,
              rank: rank || undefined,
              icon: icon || undefined,
              auths: auths || undefined,
              show_link: showLink,
              show_parent: showParent,
            };
            const res = await updateMenu(data);
            if (res.code === 0) {
              showSuccessBubble('更新成功', `菜单"${title}"已更新`);
              onSuccess();
              onClose();
            } else {
              showToast(res.msg || '更新失败', 'error');
            }
          } else {
            // 创建
            const data: CreateMenuRequest = {
              title,
              name,
              path,
              component,
              parent_id: parentId,
              menu_type: menuType,
              rank: rank || undefined,
              icon: icon || undefined,
              auths: auths || undefined,
              show_link: showLink,
              show_parent: showParent,
            };
            const res = await createMenu(data);
            if (res.code === 0) {
              showSuccessBubble('创建成功', `菜单"${title}"已创建`);
              onSuccess();
              onClose();
            } else {
              showToast(res.msg || '创建失败', 'error');
            }
          }
        } catch (error: any) {
          showToast(menu ? '更新失败' : '创建失败', 'error');
          console.error('Menu form error:', error);
        } finally {
          setLoading(false);
        }
      },
    });
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="slide"
      onRequestClose={onClose}
    >
      <View className="flex-1 bg-black/50 justify-end">
        <View className="bg-white rounded-t-3xl max-h-[90%]">
          {/* 头部 */}
          <View className="px-4 py-4 border-b border-gray-200 flex-row items-center justify-between">
            <Text className="text-lg font-semibold">
              {menu ? '编辑菜单' : '创建菜单'}
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Text className="text-gray-500 text-lg">✕</Text>
            </TouchableOpacity>
          </View>

          {/* 表单内容 */}
          <ScrollView className="px-4 py-4">
            {/* 菜单标题 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">
                菜单标题 <Text className="text-red-500">*</Text>
              </Text>
              <TextInput
                value={title}
                onChangeText={setTitle}
                placeholder="请输入菜单标题"
                className="px-3 py-2 border border-gray-300 rounded-lg"
              />
            </View>

            {/* 菜单名称 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">
                菜单名称 <Text className="text-red-500">*</Text>
              </Text>
              <TextInput
                value={name}
                onChangeText={setName}
                placeholder="例如: shop.menus"
                className="px-3 py-2 border border-gray-300 rounded-lg"
              />
            </View>

            {/* 路由路径 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">
                路由路径 <Text className="text-red-500">*</Text>
              </Text>
              <TextInput
                value={path}
                onChangeText={setPath}
                placeholder="例如: /(shop)/menus"
                className="px-3 py-2 border border-gray-300 rounded-lg"
              />
            </View>

            {/* 组件路径 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">
                组件路径 <Text className="text-red-500">*</Text>
              </Text>
              <TextInput
                value={component}
                onChangeText={setComponent}
                placeholder="例如: shop/menus"
                className="px-3 py-2 border border-gray-300 rounded-lg"
              />
            </View>

            {/* 菜单类型 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">菜单类型</Text>
              <View className="flex-row gap-2">
                <TouchableOpacity
                  onPress={() => setMenuType(1)}
                  className={`flex-1 px-3 py-2 rounded-lg border ${
                    menuType === 1
                      ? 'bg-blue-50 border-blue-500'
                      : 'border-gray-300'
                  }`}
                >
                  <Text
                    className={`text-center ${
                      menuType === 1 ? 'text-blue-600' : 'text-gray-700'
                    }`}
                  >
                    一级菜单
                  </Text>
                </TouchableOpacity>
                <TouchableOpacity
                  onPress={() => setMenuType(2)}
                  className={`flex-1 px-3 py-2 rounded-lg border ${
                    menuType === 2
                      ? 'bg-blue-50 border-blue-500'
                      : 'border-gray-300'
                  }`}
                >
                  <Text
                    className={`text-center ${
                      menuType === 2 ? 'text-blue-600' : 'text-gray-700'
                    }`}
                  >
                    二级菜单
                  </Text>
                </TouchableOpacity>
              </View>
            </View>

            {/* 排序 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">排序</Text>
              <TextInput
                value={rank}
                onChangeText={setRank}
                placeholder="数字越小越靠前"
                keyboardType="numeric"
                className="px-3 py-2 border border-gray-300 rounded-lg"
              />
            </View>

            {/* 图标 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">图标</Text>
              <TextInput
                value={icon}
                onChangeText={setIcon}
                placeholder="图标名称"
                className="px-3 py-2 border border-gray-300 rounded-lg"
              />
            </View>

            {/* 权限标识 */}
            <View className="mb-4">
              <Text className="text-sm text-gray-700 mb-2">权限标识</Text>
              <TextInput
                value={auths}
                onChangeText={setAuths}
                placeholder="多个权限用逗号分隔"
                className="px-3 py-2 border border-gray-300 rounded-lg"
                multiline
              />
              <Text className="text-xs text-gray-500 mt-1">
                例如: menu:view,menu:create
              </Text>
            </View>

            {/* 显示链接 */}
            <View className="mb-4 flex-row items-center justify-between">
              <Text className="text-sm text-gray-700">显示链接</Text>
              <Switch value={showLink} onValueChange={setShowLink} />
            </View>

            {/* 显示父级 */}
            <View className="mb-4 flex-row items-center justify-between">
              <Text className="text-sm text-gray-700">显示父级</Text>
              <Switch value={showParent} onValueChange={setShowParent} />
            </View>
          </ScrollView>

          {/* 底部按钮 */}
          <View className="px-4 py-3 border-t border-gray-200 flex-row gap-3">
            <TouchableOpacity
              onPress={onClose}
              className="flex-1 px-4 py-3 bg-gray-100 rounded-lg"
            >
              <Text className="text-center text-gray-700 font-medium">
                取消
              </Text>
            </TouchableOpacity>
            <TouchableOpacity
              onPress={handleSubmit}
              disabled={loading}
              className="flex-1 px-4 py-3 bg-blue-500 rounded-lg"
            >
              {loading ? (
                <ActivityIndicator color="#fff" />
              ) : (
                <Text className="text-center text-white font-medium">
                  {menu ? '更新' : '创建'}
                </Text>
              )}
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
}



