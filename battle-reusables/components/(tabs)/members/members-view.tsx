import React, { useState, useRef } from 'react';
import { ScrollView, View, ActivityIndicator, RefreshControl } from 'react-native';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { alert } from '@/utils/alert';
import { CreditDialog } from '@/components/(shop)/members/credit-dialog';
import { listAllUsers } from '@/services/members';
import {
  getMyGroup,
  listGroupsByHouse,
  removeMemberFromGroup,
  listGroupMembers,
  pullMembersToGroup,
  removeFromGroup,
} from '@/services/shops/groups';
import { shopsMembersList } from '@/services/shops/tables';
import {
  shopsAdminsAssign,
  shopsAdminsRevoke,
  shopsAdminsList,
  shopsAdminsMe
} from '@/services/shops/admins';
import { shopsHousesOptions } from '@/services/shops/houses';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { SetShopAdminModal } from './set-shop-admin-modal';
import { MembersList } from '@/components/(tabs)/tables/members-list';

export const MembersView = () => {
  const { isSuperAdmin, isStoreAdmin } = usePermission();

  // 状态管理
  const [houseGid, setHouseGid] = useState<string>('');
  const [keyword, setKeyword] = useState<string>('');
  const [page, setPage] = useState(1);
  const [selectedUsers, setSelectedUsers] = useState<number[]>([]);
  const [activeTab, setActiveTab] = useState<'all' | 'group' | 'game' | 'blocked'>('all');
  const [showSetAdminModal, setShowSetAdminModal] = useState(false);
  const [selectedUserForAdmin, setSelectedUserForAdmin] = useState<{ id: number; name: string } | null>(null);
  const [selectedGroupFilter, setSelectedGroupFilter] = useState<number | null>(null); // 超级管理员的圈子筛选
  const [creditDialog, setCreditDialog] = useState<{ visible: boolean; type: 'deposit' | 'withdraw'; memberId: number; memberName: string; gameId: number } | null>(null);

  // API 请求
  const { data: allUsers, loading: loadingUsers, run: runListUsers } = useRequest(listAllUsers, { manual: true });
  const { data: myGroup, loading: loadingMyGroup, run: runGetMyGroup } = useRequest(getMyGroup, { manual: true });
  const { data: allGroups, loading: loadingGroups, run: runListGroups } = useRequest(listGroupsByHouse, { manual: true });
  const { data: groupMembers, loading: loadingGroupMembers, run: runListGroupMembers } = useRequest(listGroupMembers, { manual: true });
  const { run: runRemoveMember, loading: removingMember } = useRequest(removeMemberFromGroup, { manual: true });
  const { run: runAssignAdmin, loading: assigningAdmin } = useRequest(shopsAdminsAssign, { manual: true });
  const { run: runRevokeAdmin, loading: revokingAdmin } = useRequest(shopsAdminsRevoke, { manual: true });
  const { data: shopAdmins, run: runListAdmins } = useRequest(shopsAdminsList, { manual: true });
  const { data: myAdminInfo, run: runGetMyAdminInfo } = useRequest(shopsAdminsMe, { manual: true }); // 改为手动加载
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const { data: gameMembersData, loading: loadingGameMembers, run: runListGameMembers } = useRequest(shopsMembersList, { manual: true });
  const { run: runPullToGroup } = useRequest(pullMembersToGroup, { manual: true });
  const { run: runRemoveFromGroup } = useRequest(removeFromGroup, { manual: true });
  const houseRef = useRef<TriggerRef>(null);
  
  // 无感刷新：独立管理成员数据
  const [silentMembersData, setSilentMembersData] = useState<API.ShopsMembersList | undefined>(undefined);
  const isInitialLoadRef = useRef(true); // 跟踪是否是首次加载

  function onHouseTouchStart() {
    isWeb && houseRef.current?.open();
  }

  // 加载所有用户
  const handleLoadUsers = async () => {
    await runListUsers({ page, size: 20, keyword: keyword || undefined });
  };

  // 拉圈处理（店铺管理员拉入自己的圈子）
  const handlePullToGroup = async (gamePlayerID: string, memberName: string, currentGroupName?: string) => {
    if (!myGroup) {
      alert.show({ title: '请先加载圈子信息' });
      return;
    }

    const confirmMsg = currentGroupName
      ? `确定要将 ${memberName} 从「${currentGroupName}」转移到「${myGroup.group_name}」吗？`
      : `确定要将 ${memberName} 拉入「${myGroup.group_name}」吗？`;

    alert.show({
      title: '确认拉圈',
      description: confirmMsg,
      confirmText: '确定',
      cancelText: '取消',
      onCancel: () => {
        // 用户取消操作
      },
      onConfirm: async () => {
        try {
          if (!myAdminInfo?.house_gid) {
            alert.show({ title: '无法获取店铺信息' });
            return;
          }

          await runPullToGroup({
            house_gid: myAdminInfo.house_gid,
            group_id: myGroup.id,
            game_player_ids: [gamePlayerID]
          });
          alert.show({ title: '拉圈成功' });
          // 静默刷新成员列表
          const response = await shopsMembersList({ house_gid: myAdminInfo.house_gid });
          if (response?.data) setSilentMembersData(response.data);
        } catch (err: any) {
          console.error('拉圈失败:', err);
          alert.show({ title: '拉圈失败', description: err.message || '未知错误' });
        }
      }
    });
  };

  // 踢出圈子处理
  const handleRemoveFromGroup = async (gamePlayerID: string, memberName: string, currentGroupName: string) => {
    alert.show({
      title: '确认踢出圈子',
      description: `确定要将 ${memberName} 从「${currentGroupName}」中移除吗？`,
      confirmText: '确定',
      cancelText: '取消',
      onCancel: () => {
        // 用户取消操作
      },
      onConfirm: async () => {
        try {
          if (!myAdminInfo?.house_gid) {
            alert.show({ title: '无法获取店铺信息' });
            return;
          }

          await runRemoveFromGroup({
            house_gid: myAdminInfo.house_gid,
            game_player_ids: [gamePlayerID]
          });
          alert.show({ title: '踢出圈子成功' });
          // 静默刷新成员列表
          const response = await shopsMembersList({ house_gid: myAdminInfo.house_gid });
          if (response?.data) setSilentMembersData(response.data);
        } catch (err: any) {
          console.error('踢出圈子失败:', err);
          alert.show({ title: '踢出圈子失败', description: err.message || '未知错误' });
        }
      }
    });
  };

  // 店铺管理员自动加载管理员信息和圈子
  React.useEffect(() => {
    // 只有当用户是店铺管理员时才加载
    if (isStoreAdmin) {
      runGetMyAdminInfo()
        .then((adminInfo) => {
          if (adminInfo && adminInfo.house_gid) {
            // 加载店铺管理员的圈子
            return runGetMyGroup({ house_gid: adminInfo.house_gid });
          }
        })
        .then((group) => {
          if (group) {
            // 加载圈子成员
            runListGroupMembers({ group_id: group.id, page: 1, size: 100 });
          }
        })
        .catch((error) => {
          console.error('自动加载圈子失败:', error);
        });
    }
  }, [isStoreAdmin]); // 只依赖 isStoreAdmin，避免重复调用

  // "所有用户"标签页自动加载
  React.useEffect(() => {
    if (activeTab === 'all') {
      // 进入"所有用户"标签时自动加载
      handleLoadUsers();
    }
  }, [activeTab]);

  // 游戏成员列表自动刷新（每10秒）- 无感刷新
  React.useEffect(() => {
    if (activeTab !== 'game' && activeTab !== 'blocked') return;

    const effectiveHouseGid = isStoreAdmin && myAdminInfo?.house_gid 
      ? myAdminInfo.house_gid 
      : (houseGid ? Number(houseGid) : null);

    if (!effectiveHouseGid) return;

    // 首次加载：显示loading状态
    if (isInitialLoadRef.current) {
      isInitialLoadRef.current = false;
      runListGameMembers({ house_gid: effectiveHouseGid }).then((data) => {
        if (data) setSilentMembersData(data);
      });
    }

    // 设置定时静默刷新（每10秒）
    const intervalId = setInterval(async () => {
      try {
        // 静默刷新：直接调用API，不触发loading
        const response = await shopsMembersList({ house_gid: effectiveHouseGid });
        if (response?.data) {
          setSilentMembersData(response.data);
        }
      } catch (error) {
        // 静默失败，不显示错误提示
        console.error('Silent refresh failed:', error);
      }
    }, 10000);

    // 清理定时器
    return () => {
      clearInterval(intervalId);
      isInitialLoadRef.current = true; // 重置标志
    };
  }, [activeTab, houseGid, myAdminInfo?.house_gid, isStoreAdmin]);

  // 加载我的圈子（手动触发，用于超级管理员）
  const handleLoadMyGroup = async () => {
    if (!houseGid) {
      alert.show({ title: '请输入店铺号' });
      return;
    }
    const gid = Number(houseGid);
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误' });
      return;
    }

    try {
      const group = await runGetMyGroup({ house_gid: gid });
      if (group) {
        // 加载圈子成员
        await runListGroupMembers({ group_id: group.id, page: 1, size: 100 });
      }
    } catch (error: any) {
      // 错误已由 request 函数自动显示
      console.error('加载圈子失败:', error);
    }
  };

  // 加载所有圈子（超级管理员）
  const handleLoadAllGroups = async () => {
    if (!houseGid) {
      alert.show({ title: '请输入店铺号' });
      return;
    }
    const gid = Number(houseGid);
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误' });
      return;
    }

    await runListGroups({ house_gid: gid });
  };

  // 从圈子移除成员
  const handleRemoveMember = async (userId: number) => {
    if (!myGroup) return;

    try {
      await runRemoveMember({ group_id: myGroup.id, user_id: userId });
      alert.show({ title: '移除成功' });
      // 重新加载圈子成员
      await runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
    } catch (error: any) {
      // 错误已由 request 函数自动显示
      console.error('移除成员失败:', error);
    }
  };

  // 打开设置管理员弹框
  const handleOpenSetAdminModal = (userId: number, userName: string) => {
    setSelectedUserForAdmin({ id: userId, name: userName });
    setShowSetAdminModal(true);
  };

  // 设置店铺管理员（超级管理员功能）
  const handleSetShopAdmin = async (houseGid: number) => {
    if (!selectedUserForAdmin) return;

    try {
      await runAssignAdmin({ house_gid: houseGid, user_id: selectedUserForAdmin.id, role: 'admin' });
      alert.show({ title: '设置成功' });
      setShowSetAdminModal(false);
      setSelectedUserForAdmin(null);
      // 重新加载用户列表
      await handleLoadUsers();
    } catch (error: any) {
      // 错误已由 request 函数自动显示，这里不需要再次显示
      console.error('设置管理员失败:', error);
    }
  };

  // 移除店铺管理员（超级管理员功能）
  const handleRemoveShopAdmin = async (userId: number) => {
    try {
      // 不传 house_gid，后端会自动查找
      await runRevokeAdmin({ house_gid: 0, user_id: userId });
      alert.show({ title: '移除成功' });
      // 重新加载用户列表
      await handleLoadUsers();
    } catch (error: any) {
      // 错误已由 request 函数自动显示，这里不需要再次显示
      console.error('移除管理员失败:', error);
    }
  };

  // 切换用户选择
  const toggleUserSelection = (userId: number) => {
    setSelectedUsers(prev =>
      prev.includes(userId)
        ? prev.filter(id => id !== userId)
        : [...prev, userId]
    );
  };

  // 判断用户是否在圈子中
  const isUserInGroup = (userId: number) => {
    return groupMembers?.items?.some(member => member.id === userId) || false;
  };

  // 踢出圈子（店铺管理员专用）
  const handleKickUserFromGroup = async (userId: number) => {
    if (!myGroup) return;

    try {
      await runRemoveMember({ group_id: myGroup.id, user_id: userId });
      alert.show({ title: '踢出成功' });
      // 重新加载圈子成员
      await runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
    } catch (error: any) {
      // 错误已由 request 函数自动显示
      console.error('踢出圈子失败:', error);
    }
  };

  return (
    <View className="flex-1 bg-background">
      <View className="flex-row items-center gap-2 px-4 py-3 border-b border-border bg-card">
        <Button
          variant={activeTab === 'all' ? 'default' : 'outline'}
          onPress={() => setActiveTab('all')}
          className="flex-1"
        >
          <Text>所有用户</Text>
        </Button>
        <Button
          variant={activeTab === 'group' ? 'default' : 'outline'}
          onPress={() => setActiveTab('group')}
          className="flex-1"
        >
          <Text>我的圈子</Text>
        </Button>
        <Button
          variant={activeTab === 'game' ? 'default' : 'outline'}
          onPress={() => setActiveTab('game')}
          className="flex-1"
        >
          <Text>游戏成员</Text>
        </Button>
        <Button
          variant={activeTab === 'blocked' ? 'default' : 'outline'}
          onPress={() => setActiveTab('blocked')}
          className="flex-1"
        >
          <Text>禁用名单</Text>
        </Button>
      </View>
      <View className="px-4 py-3 gap-3 border-b border-border bg-card">
        {activeTab === 'game' || activeTab === 'blocked' ? (
          <>
            {!isStoreAdmin ? (
              <View className="flex-row gap-2">
                <Select
                  value={houseGid ? ({ label: `店铺 ${houseGid}`, value: houseGid } as any) : undefined}
                  onValueChange={(opt) => setHouseGid(String(opt?.value ?? ''))}
                >
                  <SelectTrigger ref={houseRef} onTouchStart={onHouseTouchStart} className="min-w-[160px]">
                    <SelectValue placeholder={houseGid ? `店铺 ${houseGid}` : '选择店铺号'} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>店铺号</SelectLabel>
                      {(houseOptions ?? []).map((gid) => (
                        <SelectItem key={String(gid)} label={`店铺 ${gid}`} value={String(gid)}>
                          {`店铺 ${gid}`}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Button
                  onPress={() => {
                    const gid = Number(houseGid);
                    if (gid > 0) {
                      runListGameMembers({ house_gid: gid }).then(data => {
                        if (data) setSilentMembersData(data);
                      });
                    } else {
                      alert.show({ title: '请选择店铺号' });
                    }
                  }}
                  disabled={loadingGameMembers}
                >
                  <Text>{loadingGameMembers ? '加载中...' : '加载成员'}</Text>
                </Button>
              </View>
            ) : null}
            {isStoreAdmin && myAdminInfo && myAdminInfo.house_gid ? (
              <Button
                onPress={() => {
                  runListGameMembers({ house_gid: myAdminInfo.house_gid! }).then(data => {
                    if (data) setSilentMembersData(data);
                  });
                }}
                disabled={loadingGameMembers}
                variant="default"
              >
                <Text>{loadingGameMembers ? '加载中...' : '刷新游戏成员'}</Text>
              </Button>
            ) : null}
          </>
        ) : activeTab === 'all' ? (
          <>
            <View className="flex-row gap-2">
              <Input
                className="flex-1"
                placeholder="搜索用户名或昵称"
                value={keyword}
                onChangeText={setKeyword}
              />
              <Button onPress={handleLoadUsers} disabled={loadingUsers}>
                <Text>{loadingUsers ? '加载中...' : '搜索'}</Text>
              </Button>
            </View>
          </>
        ) : (
          <>
            {!isStoreAdmin ? (
              <View className="flex-row gap-2">
                <Select
                  value={houseGid ? ({ label: `店铺 ${houseGid}`, value: houseGid } as any) : undefined}
                  onValueChange={(opt) => setHouseGid(String(opt?.value ?? ''))}
                >
                  <SelectTrigger ref={houseRef} onTouchStart={onHouseTouchStart} className="min-w-[160px]">
                    <SelectValue placeholder={houseGid ? `店铺 ${houseGid}` : '选择店铺号'} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>店铺号</SelectLabel>
                      {(houseOptions ?? []).map((gid) => (
                        <SelectItem key={String(gid)} label={`店铺 ${gid}`} value={String(gid)}>
                          {`店铺 ${gid}`}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Button onPress={handleLoadMyGroup} disabled={loadingMyGroup}>
                  <Text>{loadingMyGroup ? '加载中...' : '加载圈子'}</Text>
                </Button>
              </View>
            ) : null}
            {isSuperAdmin ? (
              <Button onPress={handleLoadAllGroups} disabled={loadingGroups} variant="outline">
                <Text>查看所有圈子</Text>
              </Button>
            ) : null}
          </>
        )}
      </View>
      {myGroup && activeTab === 'group' ? (
        <View className="px-4 py-3 bg-secondary border-b border-border">
          <Text className="font-semibold text-lg">{myGroup.group_name}</Text>
          <Text className="text-muted-foreground text-sm mt-1">{myGroup.description || '暂无描述'}</Text>
          <Text className="text-muted-foreground text-sm mt-1">{`成员数量: ${myGroup.member_count || groupMembers?.total || 0}`}</Text>
        </View>
      ) : null}

      <ScrollView
        className="flex-1"
        refreshControl={
          <RefreshControl
            refreshing={loadingUsers || loadingGroupMembers || loadingGameMembers}
            onRefresh={async () => {
              if (activeTab === 'all') {
                handleLoadUsers();
              } else if (activeTab === 'group' && myGroup) {
                runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
              } else if (activeTab === 'game' || activeTab === 'blocked') {
                const effectiveHouseGid = isStoreAdmin && myAdminInfo?.house_gid 
                  ? myAdminInfo.house_gid 
                  : (houseGid ? Number(houseGid) : null);
                if (effectiveHouseGid) {
                  runListGameMembers({ house_gid: effectiveHouseGid }).then(data => {
                    if (data) setSilentMembersData(data);
                  });
                }
              }
            }}
          />
        }
      >
        {activeTab === 'game' || activeTab === 'blocked' ? (
          <View className="p-4">
            {myGroup && activeTab === 'game' ? (
              <View className="mb-4 p-3 bg-secondary/50 rounded-lg">
                <Text className="text-sm font-medium">{`拉入圈子: ${myGroup.group_name}`}</Text>
                <Text className="text-xs text-muted-foreground mt-1">点击成员卡片中的「拉入圈子」按钮即可将成员拉入您的圈子</Text>
              </View>
            ) : null}
            <MembersList 
              loading={loadingGameMembers && !silentMembersData} 
              data={(silentMembersData || gameMembersData)?.items?.filter(item => {
                // 过滤逻辑：
                // game 标签页: 显示未禁用的成员
                // blocked 标签页: 显示已禁用的成员
                if (activeTab === 'blocked') {
                  return item.forbid === true;
                }
                // game 标签页只显示未禁用的
                return !item.forbid;
              })}
              houseGid={
                isStoreAdmin && myAdminInfo?.house_gid 
                  ? myAdminInfo.house_gid 
                  : (houseGid ? Number(houseGid) : undefined)
              }
              myGroupId={myGroup?.id}
              onPullToGroup={myGroup ? handlePullToGroup : undefined}
              onRemoveFromGroup={handleRemoveFromGroup}
              onCreditChange={async () => {
                const effectiveHouseGid = houseGid || myAdminInfo?.house_gid;
                if (effectiveHouseGid) {
                  const response = await shopsMembersList({ house_gid: Number(effectiveHouseGid) });
                  if (response?.data) setSilentMembersData(response.data);
                }
              }}
              isBlockedList={activeTab === 'blocked'}
            />
          </View>
        ) : activeTab === 'all' ? (
          <View className="p-4 gap-2">
            {loadingUsers && !allUsers ? (
              <View className="py-8 items-center">
                <ActivityIndicator size="large" />
                <Text className="text-muted-foreground mt-2">加载中...</Text>
              </View>
            ) : null}
            {(allUsers?.items || []).map((user) => (
              <View
                key={user.id}
                className="bg-card rounded-lg border border-border p-4 gap-2"
              >
                <View className="flex-row justify-between items-start gap-2">
                  <View className="flex-1">
                    <View className="flex-row items-center gap-2">
                      <Text className="font-semibold">{user.nick_name || user.username}</Text>
                      {user.role === 'store_admin' ? (
                        <View className="bg-primary px-2 py-0.5 rounded">
                          <Text className="text-primary-foreground text-xs">店铺管理员</Text>
                        </View>
                      ) : null}
                      {user.role === 'super_admin' ? (
                        <View className="bg-destructive px-2 py-0.5 rounded">
                          <Text className="text-destructive-foreground text-xs">超级管理员</Text>
                        </View>
                      ) : null}
                    </View>
                    <Text className="text-muted-foreground text-sm">{`ID: ${user.id}`}</Text>
                    {user.phone ? (
                      <Text className="text-muted-foreground text-sm">{`手机: ${user.phone}`}</Text>
                    ) : null}
                  </View>
                  <View className="gap-2">
                    {isSuperAdmin && user.role !== 'super_admin' ? (
                      user.role === 'store_admin' ? (
                        <Button
                          variant="destructive"
                          onPress={() => handleRemoveShopAdmin(user.id)}
                          disabled={revokingAdmin}
                          size="sm"
                        >
                          <Text>移除管理员</Text>
                        </Button>
                      ) : (
                        <Button
                          variant="secondary"
                          onPress={() => handleOpenSetAdminModal(user.id, user.nick_name || user.username)}
                          disabled={assigningAdmin}
                          size="sm"
                        >
                          <Text>设为管理员</Text>
                        </Button>
                      )
                    ) : null}
                    {isStoreAdmin && myGroup && user.role !== 'super_admin' && user.role !== 'store_admin' && isUserInGroup(user.id) ? (
                      <Button
                        variant="destructive"
                        onPress={() => handleKickUserFromGroup(user.id)}
                        disabled={removingMember}
                        size="sm"
                      >
                        <Text>踢出圈子</Text>
                      </Button>
                    ) : null}
                  </View>
                </View>
              </View>
            ))}
            {allUsers && (allUsers.items?.length || 0) === 0 ? (
              <View className="py-8 items-center">
                <Text className="text-muted-foreground">暂无用户</Text>
              </View>
            ) : null}
            {allUsers && allUsers.total > (allUsers.items?.length || 0) ? (
              <Button
                variant="outline"
                onPress={() => {
                  setPage(p => p + 1);
                  handleLoadUsers();
                }}
                disabled={loadingUsers}
              >
                <Text>加载更多</Text>
              </Button>
            ) : null}
          </View>
        ) : (
          <View className="p-4 gap-2">
            {loadingGroupMembers && !groupMembers ? (
              <View className="py-8 items-center">
                <ActivityIndicator size="large" />
                <Text className="text-muted-foreground mt-2">加载中...</Text>
              </View>
            ) : null}
            {(groupMembers?.items || []).map((user) => {
              const gameIdMatch = (user as any).introduction?.match(/game_id:(\d+)/);
              const gameId = gameIdMatch ? parseInt(gameIdMatch[1]) : null;
              const isUnboundUser = user.id === 0 && gameId;
              return (
                <View
                  key={user.id || `game-${gameId}`}
                  className="bg-card rounded-lg border border-border p-4 gap-2"
                >
                  <View className="gap-2">
                    <View className="flex-row justify-between items-start">
                      <View className="flex-1">
                        <Text className="font-semibold">{user.nick_name || user.username}</Text>
                        {isUnboundUser ? (
                          <View className="mt-1">
                            <View className="flex-row items-center gap-2">
                              <View className="rounded-full bg-orange-500/20 px-2 py-0.5">
                                <Text className="text-xs text-orange-700 dark:text-orange-400">未绑定平台</Text>
                              </View>
                              <Text className="text-xs text-muted-foreground">{`GameID: ${gameId}`}</Text>
                            </View>
                          </View>
                        ) : (
                          <View className="mt-1">
                            <Text className="text-muted-foreground text-sm">{`ID: ${user.id}`}</Text>
                            {user.phone ? (
                              <Text className="text-muted-foreground text-sm">{`手机: ${user.phone}`}</Text>
                            ) : null}
                          </View>
                        )}
                      </View>
                    </View>
                    <View className="flex-row gap-2 mt-2 border-t border-border pt-2">
                      <Button
                        variant="destructive"
                        onPress={() => handleRemoveMember(user.id)}
                        disabled={removingMember}
                        size="sm"
                        className="flex-1"
                      >
                        <Text className="text-xs">移除</Text>
                      </Button>
                      {gameId && myGroup ? (
                        <>
                          <Button
                            variant="default"
                            size="sm"
                            className="flex-1"
                            onPress={() => setCreditDialog({
                              visible: true,
                              type: 'deposit',
                              memberId: gameId,
                              memberName: user.nick_name || user.username,
                              gameId: gameId
                            })}
                          >
                            <Text className="text-xs">上分</Text>
                          </Button>
                          <Button
                            variant="secondary"
                            size="sm"
                            className="flex-1"
                            onPress={() => setCreditDialog({
                              visible: true,
                              type: 'withdraw',
                              memberId: gameId,
                              memberName: user.nick_name || user.username,
                              gameId: gameId
                            })}
                          >
                            <Text className="text-xs">下分</Text>
                          </Button>
                        </>
                      ) : null}
                    </View>
                  </View>
                </View>
              );
            })}
            {groupMembers && (groupMembers.items?.length || 0) === 0 ? (
              <View className="py-8 items-center">
                <Text className="text-muted-foreground">圈子暂无成员</Text>
              </View>
            ) : null}
            {!myGroup && !loadingMyGroup ? (
              <View className="py-8 items-center">
                <Text className="text-muted-foreground">请先加载圈子</Text>
              </View>
            ) : null}
          </View>
        )}
        {isSuperAdmin && activeTab === 'group' && allGroups && allGroups.length > 0 ? (
          <View className="p-4 gap-2 border-t border-border">
            <Text className="font-semibold text-lg mb-2">所有圈子</Text>
            {(allGroups || []).map((group) => (
              <View
                key={group.id}
                className="bg-card rounded-lg border border-border p-4 gap-2"
              >
                <Text className="font-semibold">{group.group_name}</Text>
                <Text className="text-muted-foreground text-sm">{group.description || '暂无描述'}</Text>
                <Text className="text-muted-foreground text-sm">{`管理员ID: ${group.admin_user_id} | 成员: ${group.member_count || 0}`}</Text>
              </View>
            ))}
          </View>
        ) : null}
      </ScrollView>
      <SetShopAdminModal
        visible={showSetAdminModal}
        onClose={() => {
          setShowSetAdminModal(false);
          setSelectedUserForAdmin(null);
        }}
        onConfirm={handleSetShopAdmin}
        userName={selectedUserForAdmin?.name || ''}
        loading={assigningAdmin}
      />
      {creditDialog && myGroup ? (
        <CreditDialog
          visible={creditDialog.visible}
          type={creditDialog.type}
          houseGid={myGroup.house_gid}
          memberId={creditDialog.memberId}
          memberName={creditDialog.memberName}
          onClose={() => setCreditDialog(null)}
          onSuccess={() => {
            setCreditDialog(null);
            if (myGroup) {
              runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
            }
          }}
        />
      ) : null}
    </View>
  );
};
