import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useRequest } from '@/hooks/use-request';
import { shopsMembersKick, shopsMembersLogout } from '@/services/shops/members';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardFooter,
  InfoCardContent,
} from '@/components/shared/info-card';

type MembersItemProps = {
  houseId?: number;
  data?: API.ShopsMemberItem;
};

export const MembersItem = ({ houseId, data }: MembersItemProps) => {
  const { user_id, member_id, game_id, nick_name, member_type, user_status } = data ?? {};
  const { getLabel } = usePlazaConsts();

  const { run: kickRun, loading: kickLoading } = useRequest(shopsMembersKick, { manual: true });
  const { run: logoutRun, loading: logoutLoading } = useRequest(shopsMembersLogout, {
    manual: true,
  });

  if (typeof houseId !== 'number' || typeof member_id !== 'number') return <Text>参数错误</Text>;

  const handleKick = () => {
    kickRun({ house_gid: houseId, member_id });
  };

  const handleLogout = () => {
    logoutRun({ house_gid: houseId, member_id });
  };

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>成员 #{member_id}</InfoCardTitle>
        <InfoCardTitle>用户 {user_id}</InfoCardTitle>
        <InfoCardTitle>昵称 {nick_name}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="游戏ID" value={game_id} />
        <InfoCardRow label="成员类型" value={typeof member_type === 'number' ? getLabel('member_types', member_type) : member_type} />
        <InfoCardRow label="用户状态" value={typeof user_status === 'number' ? getLabel('user_status', user_status) : user_status} />
      </InfoCardContent>
      <InfoCardFooter>
        <PermissionGate anyOf={["shop:member:kick"]}>
          <Button disabled={kickLoading} onPress={handleKick}>
            踢出
          </Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:member:logout"]}>
          <Button disabled={logoutLoading} onPress={handleLogout}>
            下线
          </Button>
        </PermissionGate>
      </InfoCardFooter>
    </InfoCard>
  );
};
