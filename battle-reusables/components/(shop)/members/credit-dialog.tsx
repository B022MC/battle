import React, { useState } from 'react';
import { View, Modal, TextInput, TouchableOpacity, ActivityIndicator } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { showToast } from '@/utils/toast';
import { alert } from '@/utils/alert';
import { showSuccessBubble } from '@/utils/bubble-toast';
import { membersCreditDeposit, membersCreditWithdraw } from '@/services/game/funds';

type CreditDialogProps = {
  visible: boolean;
  type: 'deposit' | 'withdraw';
  houseGid: number;
  memberId: number;
  memberName?: string;
  onClose: () => void;
  onSuccess?: () => void;
};

// 生成业务单号
function generateBizNo() {
  return `${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
}

export const CreditDialog = ({
  visible,
  type,
  houseGid,
  memberId,
  memberName,
  onClose,
  onSuccess,
}: CreditDialogProps) => {
  const [amount, setAmount] = useState('');
  const [reason, setReason] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    const amountNum = parseInt(amount, 10);
    if (!amountNum || amountNum <= 0) {
      showToast('请输入正确的分数', 'error');
      return;
    }

    // 先关闭弹框再显示确认
    onClose();
    
    alert.show({
      title: type === 'deposit' ? '确认上分' : '确认下分',
      description: memberName 
        ? `确定要为 ${memberName} ${type === 'deposit' ? '上分' : '下分'} ${amountNum} 分吗？`
        : `确定要${type === 'deposit' ? '上分' : '下分'} ${amountNum} 分吗？`,
      confirmText: '确定',
      cancelText: '取消',
      onConfirm: async () => {
        setLoading(true);
        try {
          // 直接使用输入的分数
          const bizNo = generateBizNo();

          const params = {
            house_gid: houseGid,
            member_id: memberId,
            amount: amountNum,
            biz_no: bizNo,
            reason: reason || (type === 'deposit' ? '上分' : '下分'),
          };

          const res = type === 'deposit' 
            ? await membersCreditDeposit(params)
            : await membersCreditWithdraw(params);

          const actionText = type === 'deposit' ? '上分' : '下分';
          showSuccessBubble(`${actionText}成功`, `已为${memberName || '成员'}${actionText} ${amountNum} 分`);
          onSuccess?.();
        } catch (error: any) {
          showToast(error?.message || '操作失败', 'error');
        } finally {
          setLoading(false);
        }
      },
    });
    
    // 重置表单
    setAmount('');
    setReason('');
  };

  const handleCancel = () => {
    setAmount('');
    setReason('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={handleCancel}
    >
      <TouchableOpacity 
        className="flex-1 bg-black/50 justify-center items-center p-4"
        activeOpacity={1}
        onPress={handleCancel}
      >
        <TouchableOpacity activeOpacity={1} onPress={(e) => e.stopPropagation()}>
          <Card className="w-80 max-w-full p-6">
            <Text className="text-xl font-semibold mb-4">{type === 'deposit' ? '上分' : '下分'}</Text>
            {memberName ? (
              <View className="mb-4">
                <Text className="text-sm text-muted-foreground mb-1">用户</Text>
                <Text className="text-base">{memberName}</Text>
              </View>
            ) : null}
            <View className="mb-4">
              <View className="flex-row mb-2">
                <Text className="text-sm text-muted-foreground">分数</Text>
                <Text className="text-sm text-red-500">*</Text>
              </View>
              <TextInput
                className="border border-gray-300 rounded px-3 py-2"
                placeholder="请输入分数"
                keyboardType="numeric"
                value={amount}
                onChangeText={setAmount}
                editable={!loading}
              />
            </View>
            <View className="mb-6">
              <Text className="text-sm text-muted-foreground mb-2">原因</Text>
              <TextInput
                className="border border-gray-300 rounded px-3 py-2"
                placeholder={type === 'deposit' ? '上分原因（可选）' : '下分原因（可选）'}
                multiline
                numberOfLines={3}
                value={reason}
                onChangeText={setReason}
                editable={!loading}
              />
            </View>
            <View className="flex-row gap-3">
              <Button
                variant="outline"
                onPress={handleCancel}
                disabled={loading}
                className="flex-1"
              >
                <Text>取消</Text>
              </Button>
              <Button
                onPress={handleSubmit}
                disabled={loading}
                className="flex-1"
              >
                {loading ? <ActivityIndicator size="small" color="white" /> : <Text>确定</Text>}
              </Button>
            </View>
          </Card>
        </TouchableOpacity>
      </TouchableOpacity>
    </Modal>
  );
};
