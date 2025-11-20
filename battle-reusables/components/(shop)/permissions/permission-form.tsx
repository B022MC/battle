import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  Modal,
} from 'react-native';
import {
  createPermission,
  updatePermission,
  type Permission,
  type CreatePermissionRequest,
  type UpdatePermissionRequest,
} from '@/services/basic/permission';
import { showToast, toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';

interface PermissionFormProps {
  visible: boolean;
  onClose: () => void;
  onSuccess: () => void;
  permission?: Permission | null;
}

export function PermissionForm({
  visible,
  onClose,
  onSuccess,
  permission,
}: PermissionFormProps) {
  const [code, setCode] = useState('');
  const [name, setName] = useState('');
  const [category, setCategory] = useState('system');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);

  const categories = [
    { key: 'stats', label: '统计' },
    { key: 'fund', label: '资金' },
    { key: 'shop', label: '店铺' },
    { key: 'game', label: '游戏' },
    { key: 'system', label: '系统' },
  ];

  useEffect(() => {
    if (permission) {
      setCode(permission.code);
      setName(permission.name);
      setCategory(permission.category);
      setDescription(permission.description || '');
    } else {
      setCode('');
      setName('');
      setCategory('system');
      setDescription('');
    }
  }, [permission, visible]);

  const handleSubmit = async () => {
    if (!code.trim() || !name.trim()) {
      showToast('请填写完整信息', 'error');
      return;
    }

    // 二次确认
    toast.confirm({
      title: permission ? '确认更新' : '确认创建',
      description: permission 
        ? `确定要更新权限"${permission.name}"吗？`
        : `确定要创建权限"${name}"吗？`,
      type: 'warning',
      confirmText: permission ? '更新' : '创建',
      cancelText: '取消',
      onConfirm: async () => {
        setLoading(true);
        try {
          if (permission) {
            // 更新
            const data: UpdatePermissionRequest = {
              id: permission.id,
              name,
              category,
              description,
            };
            const res = await updatePermission(data);
            if (res.code === 0) {
              showSuccessBubble('更新成功', `权限"${name}"已更新`);
              onSuccess();
              onClose();
            } else {
              showToast(res.msg || '更新失败', 'error');
            }
          } else {
            // 创建
            const data: CreatePermissionRequest = {
              code,
              name,
              category,
              description,
            };
            const res = await createPermission(data);
            if (res.code === 0) {
              showSuccessBubble('创建成功', `权限"${name}"已创建`);
              onSuccess();
              onClose();
            } else {
              showToast(res.msg || '创建失败', 'error');
            }
          }
        } catch (error) {
          showToast(permission ? '更新失败' : '创建失败', 'error');
          console.error('Submit permission error:', error);
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
        <View className="bg-white rounded-t-2xl" style={{ maxHeight: '90%' }}>
          {/* 头部 */}
          <View className="flex-row items-center justify-between px-4 py-4 border-b border-gray-200">
            <Text className="text-lg font-semibold">
              {permission ? '编辑权限' : '创建权限'}
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Text className="text-gray-500 text-base">✕</Text>
            </TouchableOpacity>
          </View>

          {/* 表单内容 */}
          <ScrollView className="px-4 py-4">
            {/* 权限编码 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                权限编码 *
              </Text>
              <TextInput
                value={code}
                onChangeText={setCode}
                placeholder="例如: shop:admin:view"
                className="px-4 py-3 bg-gray-50 rounded-lg text-base"
                editable={!permission} // 编辑时不允许修改编码
              />
              <Text className="text-xs text-gray-500 mt-1">
                格式: 模块:功能:操作，例如 shop:admin:view
              </Text>
            </View>

            {/* 权限名称 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                权限名称 *
              </Text>
              <TextInput
                value={name}
                onChangeText={setName}
                placeholder="例如: 查看管理员"
                className="px-4 py-3 bg-gray-50 rounded-lg text-base"
              />
            </View>

            {/* 权限分类 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                权限分类 *
              </Text>
              <View className="flex-row flex-wrap gap-2">
                {categories.map((cat) => (
                  <TouchableOpacity
                    key={cat.key}
                    onPress={() => setCategory(cat.key)}
                    className={`px-4 py-2 rounded-lg ${
                      category === cat.key ? 'bg-blue-500' : 'bg-gray-100'
                    }`}
                  >
                    <Text
                      className={`text-sm ${
                        category === cat.key ? 'text-white' : 'text-gray-700'
                      }`}
                    >
                      {cat.label}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </View>

            {/* 权限描述 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                权限描述
              </Text>
              <TextInput
                value={description}
                onChangeText={setDescription}
                placeholder="描述这个权限的用途"
                multiline
                numberOfLines={3}
                className="px-4 py-3 bg-gray-50 rounded-lg text-base"
                style={{ minHeight: 80 }}
                textAlignVertical="top"
              />
            </View>
          </ScrollView>

          {/* 底部按钮 */}
          <View className="px-4 py-4 border-t border-gray-200">
            <View className="flex-row gap-3">
              <TouchableOpacity
                onPress={onClose}
                className="flex-1 py-3 bg-gray-100 rounded-lg"
                disabled={loading}
              >
                <Text className="text-center text-gray-700 font-medium">
                  取消
                </Text>
              </TouchableOpacity>
              <TouchableOpacity
                onPress={handleSubmit}
                className={`flex-1 py-3 rounded-lg ${
                  loading ? 'bg-blue-300' : 'bg-blue-500'
                }`}
                disabled={loading}
              >
                <Text className="text-center text-white font-medium">
                  {loading ? '提交中...' : permission ? '更新' : '创建'}
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
}

