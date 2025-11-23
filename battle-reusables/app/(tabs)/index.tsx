import React, { useState, useEffect } from 'react';
import { View, ScrollView, RefreshControl, ActivityIndicator } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { useAuthStore } from '@/hooks/use-auth-store';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { AnnouncementsCard } from '@/components/(tabs)/home/announcements-card';
import { GroupBattlesCard } from '@/components/(tabs)/home/group-battles-card';
import { MembersOnlineCard } from '@/components/(tabs)/home/members-online-card';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { getMyGroupInfo, getGroupMembers } from '@/services/home';
import { shopsAdminsMe } from '@/services/shops/admins';

export default function Screen() {
  const { isAuthenticated, user, perms } = useAuthStore();
  const { isStoreAdmin } = usePermission();
  const [refreshing, setRefreshing] = useState(false);

  // 获取店铺管理员信息
  const { data: myAdminInfo, run: runGetMyAdminInfo } = useRequest(shopsAdminsMe, { manual: true });
  
  // 获取我的圈子
  const { data: myGroup, loading: loadingGroup, run: runGetMyGroup } = useRequest(getMyGroupInfo, { manual: true });
  
  // 获取圈子成员
  const { data: groupMembers, loading: loadingMembers, run: runGetGroupMembers } = useRequest(getGroupMembers, { manual: true });

  console.log('[Index Screen] isAuthenticated:', isAuthenticated);

  // 模拟数据 - TODO: 替换为真实API调用
  const mockAnnouncements = [
    {
      id: 1,
      title: '系统升级通知',
      content: '系统将于本周六凌晨2:00-4:00进行维护升级，期间服务可能短暂中断',
      priority: 'important' as const,
      created_at: new Date().toISOString(),
    },
    {
      id: 2,
      title: '新功能上线',
      content: '首页功能全新改版，支持查看圈子战绩和成员在线情况',
      priority: 'normal' as const,
      created_at: new Date(Date.now() - 86400000).toISOString(),
    },
  ];

  const mockBattleStats = {
    today_battles: 12,
    today_wins: 8,
    today_losses: 4,
    week_battles: 45,
    week_wins: 30,
    week_winrate: 66.7,
    total_battles: 156,
  };

  const mockOnlineMembers = [
    { id: 1, nickname: '张三', is_online: true },
    { id: 2, nickname: '李四', is_online: true },
    { id: 3, nickname: '王五', is_online: false },
    { id: 4, nickname: '赵六', is_online: true },
    { id: 5, nickname: '钱七', is_online: true },
  ];

  // 初始化加载数据
  useEffect(() => {
    if (isAuthenticated && isStoreAdmin) {
      loadData();
    }
  }, [isAuthenticated, isStoreAdmin]);

  const loadData = async () => {
    try {
      // 获取管理员信息
      const adminInfo = await runGetMyAdminInfo();
      
      if (adminInfo && adminInfo.house_gid) {
        // 获取我的圈子
        const group = await runGetMyGroup({ house_gid: adminInfo.house_gid });
        
        // 获取圈子成员
        if (group && group.id) {
          await runGetGroupMembers({ group_id: group.id, page: 1, size: 100 });
        }
      }
    } catch (error) {
      console.error('[Index] Load data failed:', error);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    await loadData();
    setRefreshing(false);
  };

  return (
    <ScrollView
      className="flex-1 bg-secondary"
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
      }
    >
      <View className="gap-4 p-4">
        {/* 欢迎卡片 */}
        <Card className="p-4">
          <Text className="text-lg font-semibold mb-2">欢迎回来！</Text>
          <Text className="text-muted-foreground">
            用户: {user?.username || user?.nick_name || '未知'}
          </Text>
          <Text className="text-muted-foreground mt-1">
            认证状态: {isAuthenticated ? '已登录' : '未登录'}
          </Text>
          {myAdminInfo && (
            <Text className="text-muted-foreground mt-1">
              店铺: {myAdminInfo.house_gid}
            </Text>
          )}
          {myGroup && (
            <Text className="text-muted-foreground mt-1">
              我的圈子: {myGroup.group_name}
            </Text>
          )}
        </Card>

        {/* 系统公告 - 所有人可见 */}
        <AnnouncementsCard
          announcements={mockAnnouncements}
          onViewAll={() => console.log('View all announcements')}
        />

        {/* 圈子战绩 - 普通用户可见 */}
        {isStoreAdmin && (
          <GroupBattlesCard
            groupName={myGroup?.group_name}
            stats={myGroup ? {
              today_battles: 0,  // TODO: 需要后端API支持
              today_wins: 0,
              today_losses: 0,
              week_battles: 0,
              week_wins: 0,
              week_winrate: 0,
              total_battles: 0
            } : undefined}
            loading={loadingGroup}
            onViewDetails={() => console.log('View battle details')}
          />
        )}

        {/* 成员在线情况 - 仅管理员可见 */}
        {isStoreAdmin && (
          <PermissionGate anyOf={['shop:member:view']}>
            <MembersOnlineCard
              totalMembers={groupMembers?.total || 0}
              onlineMembers={(groupMembers?.items || []).map(member => ({
                id: member.id,
                nickname: member.nick_name || member.username,
                is_online: false  // TODO: 需要后端API支持在线状态
              }))}
              loading={loadingMembers}
              onViewAll={() => console.log('View all members')}
            />
          </PermissionGate>
        )}
      </View>
    </ScrollView>
  );
}
