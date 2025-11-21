import React, { useState } from 'react';
import { View, Modal, TextInput, TouchableOpacity, ActivityIndicator } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { showToast, toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';

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
    const amountNum = parseFloat(amount);
    if (!amountNum || amountNum <= 0) {
      showToast('请输入正确的金额', 'error');
      return;
    }

    // 二次确认
    toast.confirm({
      title: type === 'deposit' ? '确认上分' : '确认下分',
      description: memberName 
        ? `确定要为 ${memberName} ${type === 'deposit' ? '上分' : '下分'} ${amountNum} 元吗？`
        : `确定要${type === 'deposit' ? '上分' : '下分'} ${amountNum} 元吗？`,
      type: 'warning',
      confirmText: '确定',
      cancelText: '取消',
      onConfirm: async () => {
        setLoading(true);
        try {
          // 金额转换为分（cents）
          const amountInCents = Math.round(amountNum * 100);
          const bizNo = generateBizNo();

          const endpoint = type === 'deposit' 
            ? '/members/credit/deposit'
            : '/members/credit/withdraw';

          const response = await fetch(endpoint, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              house_gid: houseGid,
              member_id: memberId,
              amount: amountInCents,
              biz_no: bizNo,
              reason: reason || (type === 'deposit' ? '上分' : '下分'),
            }),
          });

          const res = await response.json();
          if (res.code === 0) {
            const actionText = type === 'deposit' ? '上分' : '下分';
            showSuccessBubble(`${actionText}成功`, `已为${memberName || '成员'}${actionText} ${amountNum} 元`);
            setAmount('');
            setReason('');
            onSuccess?.();
            onClose();
          } else {
            showToast(res.msg || '操作失败', 'error');
          }
        } catch (error) {
          showToast('操作失败', 'error');
        } finally {
          setLoading(false);
        }
      },
    });
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
            <Text className="text-xl font-semibold mb-4">
              {type === 'deposit' ? '上分' : '下分'}
            </Text>

            {memberName && (
              <View className="mb-4">
                <Text className="text-sm text-muted-foreground mb-1">用户</Text>
                <Text className="text-base">{memberName}</Text>
              </View>
            )}

            <View className="mb-4">
              <Text className="text-sm text-muted-foreground mb-2">
                金额（元）<Text className="text-red-500">*</Text>
              </Text>
              <TextInput
                className="border border-gray-300 rounded px-3 py-2"
                placeholder="请输入金额"
                keyboardType="numeric"
                value={amount}
                onChangeText={setAmount}
                editable={!loading}
              />
              <Text className="text-xs text-muted-foreground mt-1">
                1 元 = 100 分
              </Text>
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
