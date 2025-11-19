import React from 'react';
import { View, TouchableOpacity } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Trophy, TrendingUp, ChevronRight } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';

interface BattleStats {
  today_battles: number;
  today_wins: number;
  today_losses: number;
  week_battles: number;
  week_wins: number;
  week_winrate: number;
  total_battles: number;
}

interface GroupBattlesCardProps {
  groupName?: string;
  stats?: BattleStats;
  loading?: boolean;
  onViewDetails?: () => void;
}

export function GroupBattlesCard({ 
  groupName,
  stats, 
  loading = false,
  onViewDetails 
}: GroupBattlesCardProps) {
  if (loading) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Trophy} className="size-5 text-yellow-600" />
          <Text className="text-lg font-semibold">我的圈子战绩</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">加载中...</Text>
        </View>
      </Card>
    );
  }

  if (!stats) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Trophy} className="size-5 text-yellow-600" />
          <Text className="text-lg font-semibold">我的圈子战绩</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">暂无战绩数据</Text>
        </View>
      </Card>
    );
  }

  return (
    <Card className="p-4">
      {/* 标题栏 */}
      <View className="flex-row items-center justify-between mb-4">
        <View className="flex-row items-center gap-2">
          <Icon as={Trophy} className="size-5 text-yellow-600" />
          <Text className="text-lg font-semibold">我的圈子战绩</Text>
        </View>
        {onViewDetails && (
          <TouchableOpacity
            onPress={onViewDetails}
            className="flex-row items-center gap-1"
          >
            <Text className="text-sm text-blue-600">查看详情</Text>
            <Icon as={ChevronRight} className="size-4 text-blue-600" />
          </TouchableOpacity>
        )}
      </View>

      {groupName && (
        <Text className="text-sm text-muted-foreground mb-3">
          圈子：{groupName}
        </Text>
      )}

      {/* 统计数据 */}
      <View className="gap-3">
        {/* 今日数据 */}
        <View className="p-3 bg-blue-50 dark:bg-blue-950 rounded-lg">
          <View className="flex-row items-center justify-between mb-2">
            <Text className="text-sm font-medium text-blue-900 dark:text-blue-100">
              今日战绩
            </Text>
            <Text className="text-xs text-blue-700 dark:text-blue-300">
              {stats.today_battles} 场
            </Text>
          </View>
          <View className="flex-row gap-4">
            <View className="flex-1">
              <Text className="text-xs text-muted-foreground">胜</Text>
              <Text className="text-lg font-bold text-green-600">
                {stats.today_wins}
              </Text>
            </View>
            <View className="flex-1">
              <Text className="text-xs text-muted-foreground">负</Text>
              <Text className="text-lg font-bold text-red-600">
                {stats.today_losses}
              </Text>
            </View>
            <View className="flex-1">
              <Text className="text-xs text-muted-foreground">胜率</Text>
              <Text className="text-lg font-bold text-blue-600">
                {stats.today_battles > 0 
                  ? ((stats.today_wins / stats.today_battles) * 100).toFixed(0)
                  : 0}%
              </Text>
            </View>
          </View>
        </View>

        {/* 本周数据 */}
        <View className="flex-row gap-3">
          <View className="flex-1 p-3 bg-secondary rounded-lg">
            <View className="flex-row items-center gap-1 mb-1">
              <Icon as={TrendingUp} className="size-4 text-green-600" />
              <Text className="text-xs text-muted-foreground">本周胜率</Text>
            </View>
            <Text className="text-2xl font-bold text-foreground">
              {stats.week_winrate.toFixed(0)}%
            </Text>
            <Text className="text-xs text-muted-foreground mt-1">
              {stats.week_battles} 场对局
            </Text>
          </View>

          <View className="flex-1 p-3 bg-secondary rounded-lg">
            <Text className="text-xs text-muted-foreground mb-1">总对局数</Text>
            <Text className="text-2xl font-bold text-foreground">
              {stats.total_battles}
            </Text>
            <Text className="text-xs text-muted-foreground mt-1">
              累计战绩
            </Text>
          </View>
        </View>
      </View>
    </Card>
  );
}
