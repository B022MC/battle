import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { useRequest } from '@/hooks/use-request';
import { shopsCtrlAccountsUnbind } from '@/services/shops/ctrlAccounts';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardFooter,
  InfoCardContent,
} from '@/components/shared/info-card';

type CtrlAccountsItemProps = {
  houseId?: number;
  data?: API.ShopsCtrlAccountsItem;
};

export const CtrlAccountsItem = ({ houseId, data }: CtrlAccountsItemProps) => {
  const { id, identifier, login_mode, status } = data ?? {};
  const { getLoginModeLabel } = usePlazaConsts();

  const { run: unbindRun, loading: unbindLoading } = useRequest(shopsCtrlAccountsUnbind, { manual: true });

  if (typeof houseId !== 'number' || typeof id !== 'number') return <Text>参数错误</Text>;

  const handleUnbind = () => {
    unbindRun({ ctrl_id: id, house_gid: houseId });
  };

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>中控账号 #{id}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="账号" value={identifier} />
        <InfoCardRow label="登录方式" value={getLoginModeLabel(login_mode as any)} />
        <InfoCardRow label="状态" value={status === 1 ? '启用' : '禁用'} />
      </InfoCardContent>
      <InfoCardFooter>
        <Button disabled={unbindLoading} onPress={handleUnbind}>
          <Text>解绑</Text>
        </Button>
      </InfoCardFooter>
    </InfoCard>
  );
};

