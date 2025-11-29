import React from 'react';
import { View, TouchableOpacity, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Building2, Calendar, ChevronRight, TrendingUp, TrendingDown, Wallet, Users } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';
import type { HouseStats, GroupPayoff } from '@/services/battles/query-typing';

// 使用后端的 HouseStats 类型
type HousePeriodStats = HouseStats;

/**
 * 店铺统计卡片属性
 */
interface HouseStatsCardProps {
  title?: string;                       // 卡片标题（默认“店铺统计”）
  subtitle?: string;                    // 卡片副标题
  todayStats?: HousePeriodStats[];      // 今日统计
  yesterdayStats?: HousePeriodStats[];  // 昨日统计
  lastWeekStats?: HousePeriodStats[];   // 上周统计
  loading?: boolean;
  onViewDetails?: () => void;
}

export function HouseStatsCard({ 
  title = '店铺统计',
  subtitle,
  todayStats,
  yesterdayStats,
  lastWeekStats,
  loading = false,
  onViewDetails 
}: HouseStatsCardProps) {
  // 格式化分数（分转元，保留2位小数）
  const formatScore = (score: number) => {
    return (score / 100).toFixed(2);
  };

  // 格式化费用（分转元，保留2位小数）
  const formatFee = (fee: number) => {
    return (fee / 100).toFixed(2);
  };

  // 计算汇总数据
  const calculateTotal = (stats?: HousePeriodStats[]) => {
    if (!stats || stats.length === 0) {
      return {
        total_games: 0,
        total_score: 0,
        total_fee: 0,
        recharge_shang: 0,
        recharge_xia: 0,
        balance_pay: 0,
        balance_take: 0,
        balance_payoffs: [] as GroupPayoff[],
        fee_payoffs: [] as GroupPayoff[],
      };
    }
    return stats.reduce(
      (acc, item) => ({
        total_games: acc.total_games + item.total_games,
        total_score: acc.total_score + item.total_score,
        total_fee: acc.total_fee + item.total_fee,
        recharge_shang: acc.recharge_shang + (item.recharge_shang || 0),
        recharge_xia: acc.recharge_xia + (item.recharge_xia || 0),
        balance_pay: acc.balance_pay + (item.balance_pay || 0),
        balance_take: acc.balance_take + (item.balance_take || 0),
        balance_payoffs: [...acc.balance_payoffs, ...(item.balance_payoffs || [])],
        fee_payoffs: [...acc.fee_payoffs, ...(item.fee_payoffs || [])],
      }),
      {
        total_games: 0,
        total_score: 0,
        total_fee: 0,
        recharge_shang: 0,
        recharge_xia: 0,
        balance_pay: 0,
        balance_take: 0,
        balance_payoffs: [] as GroupPayoff[],
        fee_payoffs: [] as GroupPayoff[],
      }
    );
  };

  if (loading) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Building2} className="size-5 text-purple-600" />
          <Text className="text-lg font-semibold">{title}</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">加载中...</Text>
        </View>
      </Card>
    );
  }

  const hasData = todayStats || yesterdayStats || lastWeekStats;
  if (!hasData) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Building2} className="size-5 text-purple-600" />
          <Text className="text-lg font-semibold">{title}</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">暂无统计数据</Text>
        </View>
      </Card>
    );
  }

  // 计算各时间段的汇总
  const todayTotal = calculateTotal(todayStats);
  const yesterdayTotal = calculateTotal(yesterdayStats);
  const lastWeekTotal = calculateTotal(lastWeekStats);

  // 渲染单个统计卡片
  const renderStatCard = (title: string, stats: ReturnType<typeof calculateTotal>, colorScheme: 'blue' | 'green' | 'purple') => {
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
            <Text className={`text-sm font-medium ${textColor}`}>{title}</Text>
          </View>
          <Text className={`text-xs ${accentColor}`}>{`${stats.total_games} 场`}</Text>
        </View>
        <View className="flex-row gap-2 mb-2">
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">成绩</Text>
            <Text className={`text-sm font-bold ${stats.total_score >= 0 ? 'text-green-600' : 'text-red-600'}`}>{`${stats.total_score >= 0 ? '+' : ''}${formatScore(stats.total_score)}`}</Text>
          </View>
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">金额</Text>
            <Text className="text-sm font-bold text-orange-600">{formatFee(stats.total_fee)}</Text>
          </View>
        </View>
        <View className="flex-row gap-2 mb-2">
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">上分</Text>
            <Text className="text-sm font-semibold text-green-600">{`+${formatScore(stats.recharge_shang)}`}</Text>
          </View>
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">下分</Text>
            <Text className="text-sm font-semibold text-red-600">{formatScore(stats.recharge_xia)}</Text>
          </View>
        </View>
        <View className="flex-row gap-2">
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">待提</Text>
            <Text className="text-sm font-semibold text-blue-600">{formatScore(stats.balance_pay)}</Text>
          </View>
          <View className="flex-1">
            <Text className="text-xs text-muted-foreground">欠费</Text>
            <Text className="text-sm font-semibold text-red-600">{formatScore(Math.abs(stats.balance_take))}</Text>
          </View>
        </View>
        {stats.balance_payoffs && stats.balance_payoffs.length > 0 ? (
          <View className="mt-2 pt-2 border-t border-border">
            <Text className="text-xs text-muted-foreground mb-1">圈子转账</Text>
            {stats.balance_payoffs.slice(0, 3).map((payoff, index) => (
              <View key={index} className="flex-row justify-between items-center mb-0.5">
                <Text className="text-xs">{payoff.group || '主圈'}</Text>
                <Text className={`text-xs font-semibold ${payoff.value > 0 ? 'text-green-600' : 'text-red-600'}`}>{`${payoff.value > 0 ? '+' : ''}${formatScore(payoff.value)}`}</Text>
              </View>
            ))}
          </View>
        ) : null}
      </View>
    );
  };

  return (
    <Card className="p-4">
      <View className="flex-row items-center justify-between mb-4">
        <View className="flex-row items-center gap-2">
          <Icon as={Building2} className="size-5 text-purple-600" />
          <Text className="text-lg font-semibold">{title}</Text>
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

      {subtitle && (
        <Text className="text-sm text-muted-foreground mb-3">
          {subtitle}
        </Text>
      )}

      <ScrollView className="max-h-[600px]">
        <View className="gap-3">
          {todayTotal.total_games > 0 && renderStatCard('今日统计', todayTotal, 'green')}
          {yesterdayTotal.total_games > 0 && renderStatCard('昨日统计', yesterdayTotal, 'blue')}
          {lastWeekTotal.total_games > 0 && renderStatCard('上周统计', lastWeekTotal, 'purple')}
        </View>
      </ScrollView>
    </Card>
  );
}
