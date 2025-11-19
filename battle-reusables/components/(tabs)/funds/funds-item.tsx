import React, { useState } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { membersCreditDeposit, membersCreditWithdraw, membersLimitUpdate } from '@/services/game/funds';
import { PermissionGate } from '@/components/auth/PermissionGate';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardFooter,
  InfoCardContent,
} from '@/components/shared/info-card';

type FundsItemProps = {
  houseId?: number;
  data?: API.MembersWalletItem;
};

export const FundsItem = ({ houseId, data }: FundsItemProps) => {
  const { member_id, balance, forbid, limit_min } = data ?? {};

  const [amount, setAmount] = useState('');
  const [limitMin, setLimitMin] = useState(String(limit_min ?? ''));

  const { run: depositRun, loading: depositLoading } = useRequest(membersCreditDeposit, { manual: true });
  const { run: withdrawRun, loading: withdrawLoading } = useRequest(membersCreditWithdraw, { manual: true });
  const { run: limitRun, loading: limitLoading } = useRequest(membersLimitUpdate, { manual: true });

  if (typeof houseId !== 'number' || typeof member_id !== 'number') return <Text>参数错误</Text>;

  const handleDeposit = () => {
    if (!amount) return;
    depositRun({
      house_gid: houseId,
      member_id,
      amount: Number(amount),
      biz_no: `deposit-${Date.now()}`,
    });
  };

  const handleWithdraw = () => {
    if (!amount) return;
    withdrawRun({
      house_gid: houseId,
      member_id,
      amount: Number(amount),
      biz_no: `withdraw-${Date.now()}`,
    });
  };

  const handleUpdateLimit = () => {
    limitRun({
      house_gid: houseId,
      member_id,
      limit_min: limitMin ? Number(limitMin) : undefined,
    });
  };

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>成员 #{member_id}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="分数" value={balance} />
        <InfoCardRow label="禁分状态" value={forbid ? '已禁' : '正常'} />
        <InfoCardRow label="禁分阈值" value={limit_min} />
        <View className="mt-2 flex flex-row gap-2">
          <Input
            keyboardType="numeric"
            className="flex-1"
            placeholder="金额"
            value={amount}
            onChangeText={setAmount}
          />
        </View>
        <View className="mt-2 flex flex-row gap-2">
          <Input
            keyboardType="numeric"
            className="flex-1"
            placeholder="禁分阈值"
            value={limitMin}
            onChangeText={setLimitMin}
          />
        </View>
      </InfoCardContent>
      <InfoCardFooter>
        <PermissionGate anyOf={['fund:deposit']}>
          <Button disabled={depositLoading || !amount} onPress={handleDeposit}>
            上分
          </Button>
        </PermissionGate>
        <PermissionGate anyOf={['fund:withdraw']}>
          <Button disabled={withdrawLoading || !amount} onPress={handleWithdraw}>
            下分
          </Button>
        </PermissionGate>
        <PermissionGate anyOf={['fund:limit:update']}>
          <Button disabled={limitLoading} onPress={handleUpdateLimit}>
            设置阈值
          </Button>
        </PermissionGate>
      </InfoCardFooter>
    </InfoCard>
  );
};

