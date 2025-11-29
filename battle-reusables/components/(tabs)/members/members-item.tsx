import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useRequest } from '@/hooks/use-request';
import { shopsMembersKick, shopsMembersLogout, shopsMembersAddToPlatform, shopsMembersRemovePlatform, shopsMembersPin, shopsMembersUnpin, shopsMembersUpdateRemark } from '@/services/shops/members';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { useAuthStore } from '@/hooks/use-auth-store';
import { toast } from '@/utils/toast';
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
  const { user_id, member_id, game_id, nick_name, member_type, user_status, group_name, game_player_id, is_pinned, pin_order, remark } = data ?? {};
  const { getLabel } = usePlazaConsts();
  const currentUser = useAuthStore((state) => state.user);
  const [editingRemark, setEditingRemark] = React.useState(false);
  const [remarkValue, setRemarkValue] = React.useState(remark || '');

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
  const { run: pinRun, loading: pinLoading } = useRequest(shopsMembersPin, {
    manual: true,
  });
  const { run: unpinRun, loading: unpinLoading } = useRequest(shopsMembersUnpin, {
    manual: true,
  });
  const { run: updateRemarkRun, loading: updateRemarkLoading } = useRequest(shopsMembersUpdateRemark, {
    manual: true,
  });

  if (typeof houseId !== 'number' || typeof member_id !== 'number') return <Text>参数错误</Text>;

  const handleKick = () => {
    toast.confirm({
      title: '确认踢出',
      description: `确定要踢出成员 "${nick_name || member_id}" 吗？`,
      type: 'error',
      confirmText: '踢出',
      cancelText: '取消',
      confirmVariant: 'destructive',
      onConfirm: async () => {
        kickRun({ house_gid: houseId, member_id });
      },
    });
  };

  const handleLogout = () => {
    toast.confirm({
      title: '确认下线',
      description: `确定要让成员 "${nick_name || member_id}" 下线吗？`,
      type: 'warning',
      confirmText: '下线',
      cancelText: '取消',
      onConfirm: async () => {
        logoutRun({ house_gid: houseId, member_id });
      },
    });
  };

  const handleAddToGroup = () => {
    if (typeof user_id !== 'number') return;
    toast.confirm({
      title: '确认拉入圈子',
      description: `确定要将成员 "${nick_name || member_id}" 拉入圈子吗？`,
      type: 'warning',
      confirmText: '确定',
      cancelText: '取消',
      onConfirm: async () => {
        addToGroupRun({ house_gid: houseId, member_user_id: user_id });
      },
    });
  };

  const handleRemoveFromGroup = () => {
    if (typeof user_id !== 'number') return;
    toast.confirm({
      title: '确认踢出圈子',
      description: `确定要将成员 "${nick_name || member_id}" 踢出圈子吗？`,
      type: 'error',
      confirmText: '踢出',
      cancelText: '取消',
      confirmVariant: 'destructive',
      onConfirm: async () => {
        removeFromGroupRun({ house_gid: houseId, member_user_id: user_id });
      },
    });
  };

  const handlePin = () => {
    if (!game_player_id) return;
    toast.confirm({
      title: '确认置顶',
      description: `确定要置顶成员 "${nick_name || member_id}" 吗？`,
      type: 'info',
      confirmText: '置顶',
      cancelText: '取消',
      onConfirm: async () => {
        pinRun({ house_gid: houseId, game_player_id, pin_order: 0 });
      },
    });
  };

  const handleUnpin = () => {
    if (!game_player_id) return;
    toast.confirm({
      title: '确认取消置顶',
      description: `确定要取消置顶成员 "${nick_name || member_id}" 吗？`,
      type: 'info',
      confirmText: '确定',
      cancelText: '取消',
      onConfirm: async () => {
        unpinRun({ house_gid: houseId, game_player_id });
      },
    });
  };

  const handleUpdateRemark = async () => {
    if (!game_player_id) return;
    await updateRemarkRun({ house_gid: houseId, game_player_id, remark: remarkValue });
    setEditingRemark(false);
  };

  // 同步remark值的变化
  React.useEffect(() => {
    setRemarkValue(remark || '');
  }, [remark]);

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>成员 #{member_id} {is_pinned && '📌'}</InfoCardTitle>
        <InfoCardTitle>用户 {user_id}</InfoCardTitle>
        <InfoCardTitle>昵称 {nick_name}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="游戏ID" value={game_id} />
        <InfoCardRow label="成员类型" value={typeof member_type === 'number' ? getLabel('member_types', member_type) : member_type} />
        <InfoCardRow label="用户状态" value={typeof user_status === 'number' ? getLabel('user_status', user_status) : user_status} />
        {group_name && <InfoCardRow label="圈子" value={group_name} />}
        {/* 备注显示和编辑 */}
        <View className="flex-row items-center gap-2 mt-2">
          <Text className="text-muted-foreground min-w-[80px]">备注:</Text>
          {editingRemark ? (
            <Input 
              value={remarkValue} 
              onChangeText={setRemarkValue}
              placeholder="输入备注"
              className="flex-1"
            />
          ) : (
            <Text className="flex-1">{remark || '(无备注)'}</Text>
          )}
        </View>
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
        {/* 置顶/取消置顶按钮 */}
        {game_player_id && (
          <PermissionGate anyOf={["shop:member:update"]}>
            {is_pinned ? (
              <Button disabled={unpinLoading} onPress={handleUnpin} variant="outline">
                取消置顶
              </Button>
            ) : (
              <Button disabled={pinLoading} onPress={handlePin} variant="outline">
                置顶
              </Button>
            )}
          </PermissionGate>
        )}
        {/* 备注编辑按钮 */}
        {game_player_id && (
          <PermissionGate anyOf={["shop:member:update"]}>
            {editingRemark ? (
              <>
                <Button disabled={updateRemarkLoading} onPress={handleUpdateRemark} size="sm">
                  保存备注
                </Button>
                <Button disabled={updateRemarkLoading} onPress={() => {
                  setEditingRemark(false);
                  setRemarkValue(remark || '');
                }} variant="outline" size="sm">
                  取消
                </Button>
              </>
            ) : (
              <Button onPress={() => setEditingRemark(true)} variant="outline" size="sm">
                {remark ? '编辑备注' : '添加备注'}
              </Button>
            )}
          </PermissionGate>
        )}
      </InfoCardFooter>
    </InfoCard>
  );
};
