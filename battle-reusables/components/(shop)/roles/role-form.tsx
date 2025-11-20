import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  Modal,
  Switch,
} from 'react-native';
import {
  createRole,
  updateRole,
  type Role,
  type CreateRoleRequest,
  type UpdateRoleRequest,
} from '@/services/basic/role';
import { showToast, toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';

interface RoleFormProps {
  visible: boolean;
  onClose: () => void;
  onSuccess: () => void;
  role?: Role | null;
}

export function RoleForm({ visible, onClose, onSuccess, role }: RoleFormProps) {
  const [code, setCode] = useState('');
  const [name, setName] = useState('');
  const [remark, setRemark] = useState('');
  const [enable, setEnable] = useState(true);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (role) {
      setCode(role.code);
      setName(role.name);
      setRemark(role.remark || '');
      setEnable(role.enable);
    } else {
      setCode('');
      setName('');
      setRemark('');
      setEnable(true);
    }
  }, [role, visible]);

  const handleSubmit = async () => {
    if (!code.trim() || !name.trim()) {
      showToast('请填写完整信息', 'error');
      return;
    }

    // 二次确认
    toast.confirm({
      title: role ? '确认更新' : '确认创建',
      description: role 
        ? `确定要更新角色"${role.name}"吗？`
        : `确定要创建角色"${name}"吗？`,
      type: 'warning',
      confirmText: role ? '更新' : '创建',
      cancelText: '取消',
      onConfirm: async () => {
        setLoading(true);
        try {
          if (role) {
            // 更新
            const data: UpdateRoleRequest = {
              id: role.id,
              name,
              remark: remark || undefined,
              enable,
            };
            const res = await updateRole(data);
            if (res.code === 0) {
              showSuccessBubble('更新成功', `角色"${name}"已更新`);
              onSuccess();
              onClose();
            } else {
              showToast(res.msg || '更新失败', 'error');
            }
          } else {
            // 创建
            const data: CreateRoleRequest = {
              code,
              name,
              remark: remark || undefined,
            };
            const res = await createRole(data);
            if (res.code === 0) {
              showSuccessBubble('创建成功', `角色"${name}"已创建`);
              onSuccess();
              onClose();
            } else {
              showToast(res.msg || '创建失败', 'error');
            }
          }
        } catch (error) {
          showToast(role ? '更新失败' : '创建失败', 'error');
          console.error('Submit role error:', error);
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
              {role ? '编辑角色' : '创建角色'}
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Text className="text-gray-500 text-base">✕</Text>
            </TouchableOpacity>
          </View>

          {/* 表单内容 */}
          <ScrollView className="px-4 py-4">
            {/* 角色编码 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                角色编码 *
              </Text>
              <TextInput
                value={code}
                onChangeText={setCode}
                placeholder="例如: custom_admin"
                className="px-4 py-3 bg-gray-50 rounded-lg text-base"
                editable={!role} // 编辑时不允许修改编码
              />
              <Text className="text-xs text-gray-500 mt-1">
                英文字母和下划线，不能与现有角色重复
              </Text>
            </View>

            {/* 角色名称 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                角色名称 *
              </Text>
              <TextInput
                value={name}
                onChangeText={setName}
                placeholder="例如: 自定义管理员"
                className="px-4 py-3 bg-gray-50 rounded-lg text-base"
              />
            </View>

            {/* 角色备注 */}
            <View className="mb-4">
              <Text className="text-sm font-medium text-gray-700 mb-2">
                角色备注
              </Text>
              <TextInput
                value={remark}
                onChangeText={setRemark}
                placeholder="描述这个角色的用途"
                multiline
                numberOfLines={3}
                className="px-4 py-3 bg-gray-50 rounded-lg text-base"
                style={{ minHeight: 80 }}
                textAlignVertical="top"
              />
            </View>

            {/* 是否启用 */}
            {role && (
              <View className="mb-4">
                <View className="flex-row items-center justify-between">
                  <View className="flex-1">
                    <Text className="text-sm font-medium text-gray-700">
                      启用状态
                    </Text>
                    <Text className="text-xs text-gray-500 mt-1">
                      禁用后，该角色将无法被分配给用户
                    </Text>
                  </View>
                  <Switch value={enable} onValueChange={setEnable} />
                </View>
              </View>
            )}
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
                  {loading ? '提交中...' : role ? '更新' : '创建'}
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
}

