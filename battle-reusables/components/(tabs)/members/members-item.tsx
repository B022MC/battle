import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useRequest } from '@/hooks/use-request';
import { shopsMembersKick, shopsMembersLogout, shopsMembersAddToPlatform, shopsMembersRemovePlatform } from '@/services/shops/members';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { useAuthStore } from '@/hooks/use-auth-store';
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
  const { user_id, member_id, game_id, nick_name, member_type, user_status, group_name } = data ?? {};
  const { getLabel } = usePlazaConsts();
  const currentUser = useAuthStore((state) => state.user);

  // 判断是否可以拉入圈子
  // 条件: 1. 不是管理员 2. 没有圈子 3. 不是超级管理员 4. 不是店铺管理员
  const canAddToGroup = member_type !== 2 && !group_name;
  // 判断是否可以踢出圈子
  // 条件: 1. 不是管理员 2. 有圈子
  const canRemoveFromGroup = member_type !== 2 && !!group_name;

  const { run: kickRun, loading: kickLoading } = useRequest(shopsMembersKick, { manual: true });
  const { run: logoutRun, loading: logoutLoading } = useRequest(shopsMembersLogout, {
    manual: true,
  });
  const { run: addToGroupRun, loading: addToGroupLoading } = useRequest(shopsMembersAddToPlatform, {
    manual: true,
  });
  const { run: removeFromGroupRun, loading: removeFromGroupLoading } = useRequest(shopsMembersRemovePlatform, {
    manual: true,
  });

  if (typeof houseId !== 'number' || typeof member_id !== 'number') return <Text>参数错误</Text>;

  const handleKick = () => {
    kickRun({ house_gid: houseId, member_id });
  };

  const handleLogout = () => {
    logoutRun({ house_gid: houseId, member_id });
  };

  const handleAddToGroup = () => {
    if (typeof user_id !== 'number') return;
    addToGroupRun({ house_gid: houseId, member_user_id: user_id });
  };

  const handleRemoveFromGroup = () => {
    if (typeof user_id !== 'number') return;
    removeFromGroupRun({ house_gid: houseId, member_user_id: user_id });
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
        {group_name && <InfoCardRow label="圈子" value={group_name} />}
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
        {/* 拉入圈子按钮: 只对没有圈子的普通用户显示 */}
        {canAddToGroup && (
          <PermissionGate anyOf={["shop:member:update"]}>
            <Button disabled={addToGroupLoading} onPress={handleAddToGroup}>
              拉入圈子
            </Button>
          </PermissionGate>
        )}
        {/* 踢出圈子按钮: 只对有圈子的普通用户显示 */}
        {canRemoveFromGroup && (
          <PermissionGate anyOf={["shop:member:kick"]}>
            <Button disabled={removeFromGroupLoading} onPress={handleRemoveFromGroup}>
              踢出圈子
            </Button>
          </PermissionGate>
        )}
      </InfoCardFooter>
    </InfoCard>
  );
};
