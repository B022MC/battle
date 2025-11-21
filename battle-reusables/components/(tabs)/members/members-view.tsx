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
  const [activeTab, setActiveTab] = useState<'all' | 'group' | 'game'>('all');
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
          // 刷新成员列表
          runListGameMembers({ house_gid: myAdminInfo.house_gid });
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
          // 刷新成员列表
          runListGameMembers({ house_gid: myAdminInfo.house_gid });
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

  // 加载我的圈子（手动触发，用于超级管理员）
  const handleLoadMyGroup = async () => {
    if (!houseGid) {
      alert.show({ title: '请输入店铺号', variant: 'error' });
      return;
    }
    const gid = Number(houseGid);
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误', variant: 'error' });
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
      alert.show({ title: '请输入店铺号', variant: 'error' });
      return;
    }
    const gid = Number(houseGid);
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误', variant: 'error' });
      return;
    }

    await runListGroups({ house_gid: gid });
  };

  // 从圈子移除成员
  const handleRemoveMember = async (userId: number) => {
    if (!myGroup) return;

    try {
      await runRemoveMember({ group_id: myGroup.id, user_id: userId });
      alert.show({ title: '移除成功', variant: 'success' });
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
      alert.show({ title: '设置成功', variant: 'success' });
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
      alert.show({ title: '移除成功', variant: 'success' });
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
      alert.show({ title: '踢出成功', variant: 'success' });
      // 重新加载圈子成员
      await runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
    } catch (error: any) {
      // 错误已由 request 函数自动显示
      console.error('踢出圈子失败:', error);
    }
  };

  return (
    <View className="flex-1 bg-background">
      {/* 顶部标签切换 */}
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
      </View>

      {/* 搜索和操作栏 */}
      <View className="px-4 py-3 gap-3 border-b border-border bg-card">
        {activeTab === 'game' ? (
          <>
            {!isStoreAdmin && (
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
                          店铺 {gid}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Button
                  onPress={() => {
                    const gid = Number(houseGid);
                    if (gid > 0) {
                      runListGameMembers({ house_gid: gid });
                    } else {
                      alert.show({ title: '请选择店铺号', variant: 'error' });
                    }
                  }}
                  disabled={loadingGameMembers}
                >
                  <Text>{loadingGameMembers ? '加载中...' : '加载成员'}</Text>
                </Button>
              </View>
            )}
            {isStoreAdmin && myAdminInfo && myAdminInfo.house_gid && (
              <Button
                onPress={() => runListGameMembers({ house_gid: myAdminInfo.house_gid! })}
                disabled={loadingGameMembers}
                variant="default"
              >
                <Text>{loadingGameMembers ? '加载中...' : '刷新游戏成员'}</Text>
              </Button>
            )}
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
            {/* 店铺管理员不需要输入店铺号 */}
            {!isStoreAdmin && (
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
                          店铺 {gid}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Button onPress={handleLoadMyGroup} disabled={loadingMyGroup}>
                  <Text>{loadingMyGroup ? '加载中...' : '加载圈子'}</Text>
                </Button>
              </View>
            )}
            {isSuperAdmin && (
              <Button onPress={handleLoadAllGroups} disabled={loadingGroups} variant="outline">
                <Text>查看所有圈子</Text>
              </Button>
            )}
          </>
        )}
      </View>

      {/* 圈子信息 */}
      {myGroup && activeTab === 'group' && (
        <View className="px-4 py-3 bg-secondary border-b border-border">
          <Text className="font-semibold text-lg">{myGroup.group_name}</Text>
          <Text className="text-muted-foreground text-sm mt-1">{myGroup.description || '暂无描述'}</Text>
          <Text className="text-muted-foreground text-sm mt-1">
            成员数量: {myGroup.member_count || groupMembers?.total || 0}
          </Text>
        </View>
      )}

      {/* 内容区域 */}
      <ScrollView
        className="flex-1"
        refreshControl={
          <RefreshControl
            refreshing={loadingUsers || loadingGroupMembers || loadingGameMembers}
            onRefresh={() => {
              if (activeTab === 'all') {
                handleLoadUsers();
              } else if (activeTab === 'group' && myGroup) {
                runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
              } else if (activeTab === 'game') {
                if (isStoreAdmin && myAdminInfo && myAdminInfo.house_gid) {
                  runListGameMembers({ house_gid: myAdminInfo.house_gid });
                } else if (houseGid) {
                  runListGameMembers({ house_gid: Number(houseGid) });
                }
              }
            }}
          />
        }
      >
        {activeTab === 'game' ? (
          // 游戏成员列表
          <View className="p-4">
            {myGroup && (
              <View className="mb-4 p-3 bg-secondary/50 rounded-lg">
                <Text className="text-sm font-medium">拉入圈子: {myGroup.group_name}</Text>
                <Text className="text-xs text-muted-foreground mt-1">
                  点击成员卡片中的「拉入圈子」按钮即可将成员拉入您的圈子
                </Text>
              </View>
            )}
            <MembersList 
              loading={loadingGameMembers} 
              data={gameMembersData?.items}
              houseGid={houseGid ? Number(houseGid) : undefined}
              onPullToGroup={myGroup ? handlePullToGroup : undefined}
              onRemoveFromGroup={handleRemoveFromGroup}
              onCreditChange={() => {
                // 上分/下分后刷新成员列表
                const effectiveHouseGid = houseGid || myAdminInfo?.house_gid;
                if (effectiveHouseGid) {
                  runListGameMembers({ house_gid: Number(effectiveHouseGid) });
                }
              }}
            />
          </View>
        ) : activeTab === 'all' ? (
          // 所有用户列表
          <View className="p-4 gap-2">
            {loadingUsers && !allUsers && (
              <View className="py-8 items-center">
                <ActivityIndicator size="large" />
                <Text className="text-muted-foreground mt-2">加载中...</Text>
              </View>
            )}
            {(allUsers?.items || []).map((user) => (
              <View
                key={user.id}
                className="bg-card rounded-lg border border-border p-4 gap-2"
              >
                <View className="flex-row justify-between items-start gap-2">
                  <View className="flex-1">
                    <View className="flex-row items-center gap-2">
                      <Text className="font-semibold">{user.nick_name || user.username}</Text>
                      {user.role === 'store_admin' && (
                        <View className="bg-primary px-2 py-0.5 rounded">
                          <Text className="text-primary-foreground text-xs">店铺管理员</Text>
                        </View>
                      )}
                      {user.role === 'super_admin' && (
                        <View className="bg-destructive px-2 py-0.5 rounded">
                          <Text className="text-destructive-foreground text-xs">超级管理员</Text>
                        </View>
                      )}
                    </View>
                    <Text className="text-muted-foreground text-sm">ID: {user.id}</Text>
                    {user.phone && (
                      <Text className="text-muted-foreground text-sm">手机: {user.phone}</Text>
                    )}
                  </View>
                  <View className="gap-2">
                    {/* 超级管理员的管理员操作按钮 */}
                    {isSuperAdmin && user.role !== 'super_admin' && (
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
                    )}

                    {/* 店铺管理员的踢出按钮 */}
                    {isStoreAdmin && myGroup && user.role !== 'super_admin' && user.role !== 'store_admin' && isUserInGroup(user.id) && (
                      <Button
                        variant="destructive"
                        onPress={() => handleKickUserFromGroup(user.id)}
                        disabled={removingMember}
                        size="sm"
                      >
                        <Text>踢出圈子</Text>
                      </Button>
                    )}
                    {/* 已删除拉入圈子按钮，请使用游戏成员列表进行拉圈 */}
                  </View>
                </View>
              </View>
            ))}
            {allUsers && (allUsers.items?.length || 0) === 0 && (
              <View className="py-8 items-center">
                <Text className="text-muted-foreground">暂无用户</Text>
              </View>
            )}
            {allUsers && allUsers.total > (allUsers.items?.length || 0) && (
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
            )}
          </View>
        ) : (
          // 圈子成员列表
          <View className="p-4 gap-2">
            {loadingGroupMembers && !groupMembers && (
              <View className="py-8 items-center">
                <ActivityIndicator size="large" />
                <Text className="text-muted-foreground mt-2">加载中...</Text>
              </View>
            )}
            {(groupMembers?.items || []).map((user) => {
              // 从 introduction 字段提取 game_id (格式: "game_id:21309263")
              const gameIdMatch = (user as any).introduction?.match(/game_id:(\d+)/);
              const gameId = gameIdMatch ? parseInt(gameIdMatch[1]) : null;
              const isUnboundUser = user.id === 0 && gameId; // 未绑定用户标记

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
                              <Text className="text-xs text-muted-foreground">GameID: {gameId}</Text>
                            </View>
                          </View>
                        ) : (
                          <View className="mt-1">
                            <Text className="text-muted-foreground text-sm">ID: {user.id}</Text>
                            {user.phone && (
                              <Text className="text-muted-foreground text-sm">手机: {user.phone}</Text>
                            )}
                          </View>
                        )}
                      </View>
                    </View>
                    {/* 操作按钮区 */}
                    <View className="flex-row gap-2 mt-2 border-t border-border pt-2">
                      {/* 移除按钮 */}
                      <Button
                        variant="destructive"
                        onPress={() => handleRemoveMember(user.id)}
                        disabled={removingMember}
                        size="sm"
                        className="flex-1"
                      >
                        <Text className="text-xs">移除</Text>
                      </Button>
                      {/* 上分/下分按钮 - 需要 gameId */}
                      {gameId && myGroup && (
                        <>
                          <Button
                            variant="default"
                            size="sm"
                            className="flex-1"
                            onPress={() => setCreditDialog({
                              visible: true,
                              type: 'deposit',
                              memberId: gameId, // 使用 gameId 作为 memberId
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
                              memberId: gameId, // 使用 gameId 作为 memberId
                              memberName: user.nick_name || user.username,
                              gameId: gameId
                            })}
                          >
                            <Text className="text-xs">下分</Text>
                          </Button>
                        </>
                      )}
                    </View>
                  </View>
                </View>
              );
            })}
            {groupMembers && (groupMembers.items?.length || 0) === 0 && (
              <View className="py-8 items-center">
                <Text className="text-muted-foreground">圈子暂无成员</Text>
              </View>
            )}
            {!myGroup && !loadingMyGroup && (
              <View className="py-8 items-center">
                <Text className="text-muted-foreground">请先加载圈子</Text>
              </View>
            )}
          </View>
        )}

        {/* 超级管理员查看所有圈子 */}
        {isSuperAdmin && activeTab === 'group' && allGroups && allGroups.length > 0 && (
          <View className="p-4 gap-2 border-t border-border">
            <Text className="font-semibold text-lg mb-2">所有圈子</Text>
            {(allGroups || []).map((group) => (
              <View
                key={group.id}
                className="bg-card rounded-lg border border-border p-4 gap-2"
              >
                <Text className="font-semibold">{group.group_name}</Text>
                <Text className="text-muted-foreground text-sm">{group.description || '暂无描述'}</Text>
                <Text className="text-muted-foreground text-sm">
                  管理员ID: {group.admin_user_id} | 成员: {group.member_count || 0}
                </Text>
              </View>
            ))}
          </View>
        )}
      </ScrollView>

      {/* 设置店铺管理员弹框 */}
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

      {/* 上分/下分对话框 */}
      {creditDialog && myGroup && (
        <CreditDialog
          visible={creditDialog.visible}
          type={creditDialog.type}
          houseGid={myGroup.house_gid}
          memberId={creditDialog.memberId}
          memberName={creditDialog.memberName}
          onClose={() => setCreditDialog(null)}
          onSuccess={() => {
            setCreditDialog(null);
            // 刷新圈子成员列表
            if (myGroup) {
              runListGroupMembers({ group_id: myGroup.id, page: 1, size: 100 });
            }
          }}
        />
      )}
    </View>
  );
};
