import React from 'react';
import { View, TouchableOpacity } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Trophy, Calendar, ChevronRight } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';

/**
 * 单个时间段的统计数据
 */
interface PeriodStats {
  total_games: number;      // 总局数
  total_score: number;      // 总得分（分）
  total_fee: number;        // 总手续费（分）
  active_members: number;   // 活跃成员数
}

/**
 * 圈子战绩卡片属性
 */
interface GroupBattlesCardProps {
  groupName?: string;
  yesterdayStats?: PeriodStats;  // 昨日统计
  thisWeekStats?: PeriodStats;   // 本周统计
  lastWeekStats?: PeriodStats;   // 上周统计
  loading?: boolean;
  onViewDetails?: () => void;
}

export function GroupBattlesCard({ 
  groupName,
  yesterdayStats,
  thisWeekStats,
  lastWeekStats,
  loading = false,
  onViewDetails 
}: GroupBattlesCardProps) {
  // 格式化分数（分转元，保留2位小数）
  const formatScore = (score: number) => {
    return (score / 100).toFixed(2);
  };

  // 格式化费用（分转元，保留2位小数）
  const formatFee = (fee: number) => {
    return (fee / 100).toFixed(2);
  };

  if (loading) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Trophy} className="size-5 text-yellow-600" />
          <Text className="text-lg font-semibold">圈子战绩统计</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">加载中...</Text>
        </View>
      </Card>
    );
  }

  const hasData = yesterdayStats || thisWeekStats || lastWeekStats;
  if (!hasData) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Trophy} className="size-5 text-yellow-600" />
          <Text className="text-lg font-semibold">圈子战绩统计</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">暂无战绩数据</Text>
        </View>
      </Card>
    );
  }

  // 渲染单个统计卡片
  const renderStatCard = (title: string, stats: PeriodStats, colorScheme: 'blue' | 'green' | 'purple') => {
    const bgColor = colorScheme === 'blue' 
      ? 'bg-blue-50 dark:bg-blue-950' 
      : colorScheme === 'green'
      ? 'bg-green-50 dark:bg-green-950'
      : 'bg-purple-50 dark:bg-purple-950';
    
    const textColor = colorScheme === 'blue'
      ? 'text-blue-900 dark:text-blue-100'
      : colorScheme === 'green'
      ? 'text-green-900 dark:text-green-100'
      : 'text-purple-900 dark:text-purple-100';

    const accentColor = colorScheme === 'blue'
      ? 'text-blue-700 dark:text-blue-300'
      : colorScheme === 'green'
      ? 'text-green-700 dark:text-green-300'
      : 'text-purple-700 dark:text-purple-300';

    return (
      <View className={`p-3 ${bgColor} rounded-lg`}>
        <View className="flex-row items-center justify-between mb-2">
          <View className="flex-row items-center gap-1">
            <Icon as={Calendar} className={`size-4 ${accentColor}`} />
            <Text className={`text-sm font-medium ${textColor}`}>
              {title}
            </Text>
          </View>
          <Text className={`text-xs ${accentColor}`}>
            {stats.total_games} 场
          </Text>
        </View>
        <View className="flex-row gap-3">
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">得分</Text>
            <Text className={`text-base font-bold ${stats.total_score >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {stats.total_score >= 0 ? '+' : ''}{formatScore(stats.total_score)}
            </Text>
          </View>
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">手续费</Text>
            <Text className="text-base font-bold text-orange-600">
              {formatFee(stats.total_fee)}
            </Text>
          </View>
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">活跃人数</Text>
            <Text className="text-base font-bold text-blue-600">
              {stats.active_members}
            </Text>
          </View>
        </View>
      </View>
    );
  };

  return (
    <Card className="p-4">
      {/* 标题栏 */}
      <View className="flex-row items-center justify-between mb-4">
        <View className="flex-row items-center gap-2">
          <Icon as={Trophy} className="size-5 text-yellow-600" />
          <Text className="text-lg font-semibold">圈子战绩统计</Text>
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
        {/* 昨日统计 */}
        {yesterdayStats && renderStatCard('昨日统计', yesterdayStats, 'blue')}
        
        {/* 本周统计 */}
        {thisWeekStats && renderStatCard('本周统计', thisWeekStats, 'green')}
        
        {/* 上周统计 */}
        {lastWeekStats && renderStatCard('上周统计', lastWeekStats, 'purple')}
      </View>
    </Card>
  );
}
