import React, { useState, useRef } from 'react';
import { ScrollView, View, ActivityIndicator, RefreshControl } from 'react-native';
import { useRequest } from '@/hooks/use-request';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { alert } from '@/utils/alert';
import { listMyBattles, getMyStats } from '@/services/battles/query';
import type { BattleRecord, BattleStats } from '@/services/battles/query-typing';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { shopsHousesOptions } from '@/services/shops/houses';
import { getGroupOptions } from '@/services/shops/groups';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';

export const MyBattlesView = () => {
  // 状态管理
  const [houseGid, setHouseGid] = useState<string>('');
  const [groupId, setGroupId] = useState<string>('');
  const [page, setPage] = useState(1);

  // 计算最近7天的时间范围
  const getTimeRange = () => {
    const now = new Date();
    const endTime = Math.floor(now.getTime() / 1000);
    const startTime = endTime - 7 * 24 * 60 * 60; // 7天前
    return { startTime, endTime };
  };

  // API 请求
  const { data: battlesData, loading: loadingBattles, run: runListBattles } = useRequest(listMyBattles, { manual: true });
  const { data: statsData, loading: loadingStats, run: runGetStats } = useRequest(getMyStats, { manual: true });
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const { data: groupOptions, run: runGetGroupOptions } = useRequest(getGroupOptions, { manual: true });
  const houseRef = useRef<TriggerRef>(null);
  const groupRef = useRef<TriggerRef>(null);

  function onHouseTouchStart() {
    isWeb && houseRef.current?.open();
  }

  function onGroupTouchStart() {
    isWeb && groupRef.current?.open();
  }

  // 当店铺改变时，加载该店铺的圈子列表
  const handleHouseChange = async (newHouseGid: string) => {
    setHouseGid(newHouseGid);
    setGroupId(''); // 重置圈子选择
    if (newHouseGid) {
      try {
        await runGetGroupOptions({ house_gid: Number(newHouseGid) });
      } catch (error) {
        console.error('加载圈子列表失败:', error);
      }
    }
  };

  // 加载战绩
  const handleLoadBattles = async () => {
    if (!houseGid) {
      alert.show({ title: '请选择店铺号', variant: 'error' });
      return;
    }
    const gid = Number(houseGid);
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误', variant: 'error' });
      return;
    }

    const { startTime, endTime } = getTimeRange();
    const params: any = {
      house_gid: gid,
      page,
      size: 20,
      start_time: startTime,
      end_time: endTime,
    };

    if (groupId) {
      const gidNum = Number(groupId);
      if (!isNaN(gidNum) && gidNum > 0) {
        params.group_id = gidNum;
      }
    }

    try {
      await runListBattles(params);
    } catch (error: any) {
      console.error('加载战绩失败:', error);
    }
  };

  // 加载统计
  const handleLoadStats = async () => {
    if (!houseGid) {
      alert.show({ title: '请选择店铺号', variant: 'error' });
      return;
    }
    const gid = Number(houseGid);
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误', variant: 'error' });
      return;
    }

    const { startTime, endTime } = getTimeRange();
    const params: any = {
      house_gid: gid,
      start_time: startTime,
      end_time: endTime,
    };

    if (groupId) {
      const gidNum = Number(groupId);
      if (!isNaN(gidNum) && gidNum > 0) {
        params.group_id = gidNum;
      }
    }

    try {
      await runGetStats(params);
    } catch (error: any) {
      console.error('加载统计失败:', error);
    }
  };

  // 格式化时间
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // 格式化分数
  const formatScore = (score: number) => {
    return score > 0 ? `+${score}` : score.toString();
  };

  return (
    <ScrollView className="flex-1 bg-background p-4">
      {/* 查询表单 */}
      <Card className="p-4 mb-4">
        <Text className="text-lg font-bold mb-4">查询我的战绩</Text>
        
        <View className="mb-3">
          <Text className="mb-2">店铺号 *</Text>
          <Select
            value={houseGid ? ({ label: `店铺 ${houseGid}`, value: houseGid } as any) : undefined}
            onValueChange={(opt) => handleHouseChange(String(opt?.value ?? ''))}
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
        </View>

        <View className="mb-3">
          <Text className="mb-2">圈子 (可选)</Text>
          <Select
            value={groupId ? ({ label: groupOptions?.find(g => String(g.id) === groupId)?.name || `圈子 ${groupId}`, value: groupId } as any) : undefined}
            onValueChange={(opt) => setGroupId(String(opt?.value ?? ''))}
          >
            <SelectTrigger ref={groupRef} onTouchStart={onGroupTouchStart} className="min-w-[160px]">
              <SelectValue placeholder={groupId ? `圈子 ${groupId}` : '选择圈子 (可选)'} />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>圈子</SelectLabel>
                {(groupOptions ?? []).map((group) => (
                  <SelectItem key={String(group.id)} label={group.name} value={String(group.id)}>
                    {group.name}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </View>

        <View className="mb-3">
          <Text className="text-sm text-muted-foreground">
            查询最近7天的战绩数据
          </Text>
        </View>

        <View className="flex-row gap-2">
          <Button
            className="flex-1"
            onPress={handleLoadBattles}
            disabled={loadingBattles}
          >
            <Text>{loadingBattles ? '查询中...' : '查询战绩'}</Text>
          </Button>
          <Button
            className="flex-1"
            variant="secondary"
            onPress={handleLoadStats}
            disabled={loadingStats}
          >
            <Text>{loadingStats ? '统计中...' : '查询统计'}</Text>
          </Button>
        </View>
      </Card>

      {/* 统计信息 */}
      {statsData && (
        <Card className="p-4 mb-4">
          <Text className="text-lg font-bold mb-3">统计信息</Text>
          {statsData.group_name && (
            <Text className="mb-2">圈子: {statsData.group_name}</Text>
          )}
          <Text className="mb-2">总局数: {statsData.total_games}</Text>
          <Text className="mb-2">总分数: {formatScore(statsData.total_score)}</Text>
          <Text className="mb-2">总费用: {statsData.total_fee}</Text>
          <Text className="mb-2">平均分: {statsData.avg_score.toFixed(2)}</Text>
        </Card>
      )}

      {/* 战绩列表 */}
      {battlesData && (
        <Card className="p-4">
          <View className="flex-row justify-between items-center mb-3">
            <Text className="text-lg font-bold">战绩列表</Text>
            <Text className="text-sm text-muted-foreground">
              共 {battlesData.total} 条
            </Text>
          </View>

          {battlesData.list.length === 0 ? (
            <Text className="text-center text-muted-foreground py-8">
              暂无战绩记录
            </Text>
          ) : (
            <View>
              {battlesData.list.map((battle: BattleRecord) => (
                <View
                  key={battle.id}
                  className="border-b border-border py-3 last:border-b-0"
                >
                  <View className="flex-row justify-between items-start mb-2">
                    <View className="flex-1">
                      <Text className="font-medium mb-1">
                        {battle.group_name || '未知圈子'}
                      </Text>
                      <Text className="text-sm text-muted-foreground">
                        {formatDate(battle.battle_at)}
                      </Text>
                    </View>
                    <View className="items-end">
                      <Text
                        className={`text-lg font-bold ${
                          battle.score > 0
                            ? 'text-green-600'
                            : battle.score < 0
                            ? 'text-red-600'
                            : 'text-foreground'
                        }`}
                      >
                        {formatScore(battle.score)}
                      </Text>
                      <Text className="text-sm text-muted-foreground">
                        费用: {battle.fee}
                      </Text>
                    </View>
                  </View>
                  <View className="flex-row justify-between">
                    <Text className="text-sm text-muted-foreground">
                      房间: {battle.room_uid}
                    </Text>
                    <Text className="text-sm text-muted-foreground">
                      余额: {battle.player_balance}
                    </Text>
                  </View>
                </View>
              ))}

              {/* 分页 */}
              {battlesData.total > 20 && (
                <View className="flex-row justify-center gap-2 mt-4">
                  <Button
                    variant="outline"
                    onPress={() => {
                      if (page > 1) {
                        setPage(page - 1);
                        handleLoadBattles();
                      }
                    }}
                    disabled={page === 1 || loadingBattles}
                  >
                    <Text>上一页</Text>
                  </Button>
                  <Text className="py-2 px-4">
                    第 {page} 页
                  </Text>
                  <Button
                    variant="outline"
                    onPress={() => {
                      if (page * 20 < battlesData.total) {
                        setPage(page + 1);
                        handleLoadBattles();
                      }
                    }}
                    disabled={page * 20 >= battlesData.total || loadingBattles}
                  >
                    <Text>下一页</Text>
                  </Button>
                </View>
              )}
            </View>
          )}
        </Card>
      )}
    </ScrollView>
  );
};

