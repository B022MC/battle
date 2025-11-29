import React, { useState, useEffect } from 'react';
import { View, ScrollView, RefreshControl, ActivityIndicator } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { useAuthStore } from '@/hooks/use-auth-store';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { AnnouncementsCard } from '@/components/(tabs)/home/announcements-card';
import { HouseStatsCard } from '@/components/(tabs)/home/house-stats-card';
import { MembersOnlineCard } from '@/components/(tabs)/home/members-online-card';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { getMyGroupInfo, getGroupMembers } from '@/services/home';
import { shopsAdminsMe } from '@/services/shops/admins';
import { shopsHousesOptions } from '@/services/shops/houses';
import { getHouseStats } from '@/services/battles/query';
import type { HouseStats } from '@/services/battles/query-typing';

/**
 * 获取时间范围
 */
function getTimeRange(type: 'today' | 'yesterday' | 'thisWeek' | 'lastWeek'): { start_time: number; end_time: number } {
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  
  if (type === 'today') {
    // 今日：今天0点到现在
    return {
      start_time: Math.floor(today.getTime() / 1000),
      end_time: Math.floor(now.getTime() / 1000)
    };
  } else if (type === 'yesterday') {
    // 昨日：昨天0点到今天0点
    const yesterdayStart = new Date(today);
    yesterdayStart.setDate(yesterdayStart.getDate() - 1);
    return {
      start_time: Math.floor(yesterdayStart.getTime() / 1000),
      end_time: Math.floor(today.getTime() / 1000)
    };
  } else if (type === 'thisWeek') {
    // 本周：本周一0点到现在
    const dayOfWeek = today.getDay() || 7; // 0是周日，转为7
    const thisWeekStart = new Date(today);
    thisWeekStart.setDate(thisWeekStart.getDate() - dayOfWeek + 1); // 周一
    return {
      start_time: Math.floor(thisWeekStart.getTime() / 1000),
      end_time: Math.floor(now.getTime() / 1000)
    };
  } else {
    // 上周：上周一0点到上周日24点
    const dayOfWeek = today.getDay() || 7;
    const lastWeekStart = new Date(today);
    lastWeekStart.setDate(lastWeekStart.getDate() - dayOfWeek - 6); // 上周一
    const lastWeekEnd = new Date(lastWeekStart);
    lastWeekEnd.setDate(lastWeekEnd.getDate() + 7); // 上周一+7天
    return {
      start_time: Math.floor(lastWeekStart.getTime() / 1000),
      end_time: Math.floor(lastWeekEnd.getTime() / 1000)
    };
  }
}

export default function Screen() {
  const { isAuthenticated, user, perms } = useAuthStore();
  const { isSuperAdmin, isStoreAdmin } = usePermission();
  const [refreshing, setRefreshing] = useState(false);

  // 获取店铺管理员信息
  const { data: myAdminInfo, run: runGetMyAdminInfo } = useRequest(shopsAdminsMe, { manual: true });
  
  // 获取我的圈子
  const { data: myGroup, loading: loadingGroup, run: runGetMyGroup } = useRequest(getMyGroupInfo, { manual: true });
  
  // 获取圈子成员
  const { data: groupMembers, loading: loadingMembers, run: runGetGroupMembers } = useRequest(getGroupMembers, { manual: true });

  // 店铺统计（店铺管理员）
  const [shopTodayStats, setShopTodayStats] = useState<HouseStats | null>(null);
  const [shopYesterdayStats, setShopYesterdayStats] = useState<HouseStats | null>(null);
  const [shopLastWeekStats, setShopLastWeekStats] = useState<HouseStats | null>(null);
  const [loadingShopStats, setLoadingShopStats] = useState(false);

  // 全平台统计（超级管理员）
  const [allHouses, setAllHouses] = useState<number[]>([]);
  const [houseTodayStats, setHouseTodayStats] = useState<HouseStats[]>([]);
  const [houseYesterdayStats, setHouseYesterdayStats] = useState<HouseStats[]>([]);
  const [houseLastWeekStats, setHouseLastWeekStats] = useState<HouseStats[]>([]);
  const [loadingHouseStats, setLoadingHouseStats] = useState(false);

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
    if (isAuthenticated) {
      if (isSuperAdmin) {
        loadSuperAdminData();
      } else if (isStoreAdmin) {
        loadStoreAdminData();
      }
    }
  }, [isAuthenticated, isSuperAdmin, isStoreAdmin]);

  // 店铺管理员加载数据
  const loadStoreAdminData = async () => {
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
        
        // 加载店铺统计（今日、昨日、上周）
        await loadShopStats(adminInfo.house_gid);
      }
    } catch (error) {
      console.error('[Index] Load store admin data failed:', error);
    }
  };

  // 超级管理员加载数据
  const loadSuperAdminData = async () => {
    try {
      // 获取所有店铺列表
      const housesResponse = await shopsHousesOptions();
      const houses = housesResponse.data || [];
      setAllHouses(houses);
      
      if (houses.length > 0) {
        // 加载全平台统计
        await loadAllHousesStats(houses);
      }
    } catch (error) {
      console.error('[Index] Load super admin data failed:', error);
    }
  };

  // 加载店铺统计数据
  const loadShopStats = async (houseGid: number) => {
    setLoadingShopStats(true);
    try {
      // 并发请求三个时间段的统计
      const [today, yesterday, lastWeek] = await Promise.all([
        getHouseStats({
          house_gid: houseGid,
          ...getTimeRange('today')
        }),
        getHouseStats({
          house_gid: houseGid,
          ...getTimeRange('yesterday')
        }),
        getHouseStats({
          house_gid: houseGid,
          ...getTimeRange('lastWeek')
        })
      ]);

      setShopTodayStats(today.data || null);
      setShopYesterdayStats(yesterday.data || null);
      setShopLastWeekStats(lastWeek.data || null);
    } catch (error) {
      console.error('[Index] Load shop stats failed:', error);
    } finally {
      setLoadingShopStats(false);
    }
  };

  // 加载全平台统计数据
  const loadAllHousesStats = async (houses: number[]) => {
    setLoadingHouseStats(true);
    try {
      // 为每个店铺并发请求三个时间段的统计
      const allRequests = houses.flatMap(houseGid => [
        getHouseStats({ house_gid: houseGid, ...getTimeRange('today') }),
        getHouseStats({ house_gid: houseGid, ...getTimeRange('yesterday') }),
        getHouseStats({ house_gid: houseGid, ...getTimeRange('lastWeek') })
      ]);

      const results = await Promise.all(allRequests);
      
      // 分类结果
      const todayResults: HouseStats[] = [];
      const yesterdayResults: HouseStats[] = [];
      const lastWeekResults: HouseStats[] = [];
      
      results.forEach((result, index) => {
        const periodIndex = index % 3;
        if (result.data) {
          if (periodIndex === 0) {
            todayResults.push(result.data);
          } else if (periodIndex === 1) {
            yesterdayResults.push(result.data);
          } else {
            lastWeekResults.push(result.data);
          }
        }
      });

      setHouseTodayStats(todayResults);
      setHouseYesterdayStats(yesterdayResults);
      setHouseLastWeekStats(lastWeekResults);
    } catch (error) {
      console.error('[Index] Load all houses stats failed:', error);
    } finally {
      setLoadingHouseStats(false);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    if (isSuperAdmin) {
      await loadSuperAdminData();
    } else if (isStoreAdmin) {
      await loadStoreAdminData();
    }
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

        {/* 超级管理员 - 全平台统计 */}
        {isSuperAdmin && (
          <HouseStatsCard
            title="全平台统计"
            subtitle="所有店铺数据汇总"
            todayStats={houseTodayStats}
            yesterdayStats={houseYesterdayStats}
            lastWeekStats={houseLastWeekStats}
            loading={loadingHouseStats}
            onViewDetails={() => console.log('View all houses details')}
          />
        )}

        {/* 店铺管理员 - 店铺统计 */}
        {isStoreAdmin && (
          <HouseStatsCard
            title="店铺统计"
            subtitle={myAdminInfo ? `店铺 ID: ${myAdminInfo.house_gid}` : undefined}
            todayStats={shopTodayStats ? [shopTodayStats] : undefined}
            yesterdayStats={shopYesterdayStats ? [shopYesterdayStats] : undefined}
            lastWeekStats={shopLastWeekStats ? [shopLastWeekStats] : undefined}
            loading={loadingShopStats}
            onViewDetails={() => console.log('View shop details')}
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
