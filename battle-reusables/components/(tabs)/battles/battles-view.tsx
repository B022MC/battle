import React, { useState } from 'react';
import { ScrollView, View, RefreshControl } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { InfoCard, InfoCardContent, InfoCardHeader, InfoCardTitle } from '@/components/shared/info-card';
import { useRequest } from '@/hooks/use-request';
import { listMyBattles, getMyStats } from '@/services/battles';
import type { BattleRecord } from '@/services/battles/typing';

export const BattlesView = () => {
  const [page, setPage] = useState(1);
  const [refreshing, setRefreshing] = useState(false);

  // 查询战绩列表
  const { data: battlesData, run: runListBattles, loading: loadingBattles } = useRequest(
    () => listMyBattles({ page, size: 20 }),
    { manual: true }
  );

  // 查询战绩统计
  const { data: stats, run: runGetStats, loading: loadingStats } = useRequest(
    () => getMyStats({}),
    { manual: true }
  );

  // 初始加载
  React.useEffect(() => {
    runListBattles();
    runGetStats();
  }, []);

  // 刷新
  const onRefresh = async () => {
    setRefreshing(true);
    setPage(1);
    await Promise.all([
      runListBattles(),
      runGetStats(),
    ]);
    setRefreshing(false);
  };

  // 加载更多
  const onLoadMore = () => {
    if (loadingBattles) return;
    const nextPage = page + 1;
    setPage(nextPage);
    runListBattles();
  };

  const battles = battlesData?.list || [];
  const total = battlesData?.total || 0;

  return (
    <View className="flex-1">
      <ScrollView
        className="flex-1 bg-secondary"
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        <View className="gap-4 p-4">
          {/* 统计卡片 */}
          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>战绩统计</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-3">
                <View className="flex-row justify-between">
                  <Text variant="muted">总局数</Text>
                  <Text className="font-semibold">{stats?.total_games || 0}</Text>
                </View>
                <View className="flex-row justify-between">
                  <Text variant="muted">总分数</Text>
                  <Text className={`font-semibold ${(stats?.total_score || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {stats?.total_score || 0}
                  </Text>
                </View>
                <View className="flex-row justify-between">
                  <Text variant="muted">总手续费</Text>
                  <Text className="font-semibold">{stats?.total_fee || 0}</Text>
                </View>
              </View>
            </InfoCardContent>
          </InfoCard>

          {/* 战绩列表 */}
          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>战绩记录 ({total})</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              {battles.length === 0 ? (
                <View className="py-8">
                  <Text className="text-center text-muted-foreground">暂无战绩记录</Text>
                </View>
              ) : (
                <View className="gap-3">
                  {battles.map((battle: BattleRecord, index: number) => {
                    // 解析对战信息
                    let playersInfo: { players: Array<{ UserGameID: number; Score: number }>; opponents: Array<{ UserGameID: number; Score: number }> } | null = null;
                    try {
                      const players = JSON.parse(battle.players_json || '[]') as Array<{ UserGameID: number; Score: number }>;
                      if (players.length > 0) {
                        const opponents = players.filter(p => p.UserGameID !== battle.player_game_id);
                        playersInfo = { players, opponents };
                      }
                    } catch {}
                    
                    // 胜负状态
                    const result = battle.score > 0 
                      ? { text: '胜', color: 'text-green-600', bg: 'bg-green-100' }
                      : battle.score < 0 
                        ? { text: '负', color: 'text-red-600', bg: 'bg-red-100' }
                        : { text: '平', color: 'text-gray-600', bg: 'bg-gray-100' };
                    
                    const formatScore = (s: number) => s > 0 ? `+${s}` : String(s);
                    
                    return (
                      <View
                        key={battle.id || index}
                        className="border-b border-border pb-3 last:border-b-0"
                      >
                        {/* 第一行：时间 + 胜负 + 得分 */}
                        <View className="flex-row justify-between items-center mb-2">
                          <View className="flex-row items-center gap-2">
                            <Text className="text-xs text-muted-foreground">
                              {new Date(battle.battle_at).toLocaleString('zh-CN')}
                            </Text>
                            <View className={`px-1.5 py-0.5 rounded ${result.bg}`}>
                              <Text className={`text-xs font-bold ${result.color}`}>{result.text}</Text>
                            </View>
                          </View>
                          <Text className={`text-lg font-bold ${result.color}`}>
                            {formatScore(battle.score)}
                          </Text>
                        </View>
                        
                        {/* 第二行：对战详情 */}
                        {playersInfo && playersInfo.opponents.length > 0 && (
                          <View className="bg-secondary/50 rounded p-2 mb-2">
                            <View className="flex-row flex-wrap gap-2">
                              {playersInfo.players.map((p) => (
                                <View 
                                  key={p.UserGameID}
                                  className={`flex-row items-center gap-1 px-2 py-0.5 rounded ${
                                    p.UserGameID === battle.player_game_id 
                                      ? 'bg-primary/10' 
                                      : 'bg-background'
                                  }`}
                                >
                                  <Text className="text-xs">
                                    {p.UserGameID === battle.player_game_id ? '我' : '对手'}
                                  </Text>
                                  <Text className={`text-xs font-medium ${
                                    p.Score > 0 ? 'text-green-600' : p.Score < 0 ? 'text-red-600' : ''
                                  }`}>
                                    {formatScore(p.Score)}
                                  </Text>
                                </View>
                              ))}
                            </View>
                          </View>
                        )}
                        
                        {/* 第三行：房间、底分、手续费 */}
                        <View className="flex-row justify-between items-center">
                          <View className="flex-row gap-3">
                            <Text className="text-xs text-muted-foreground">房间 {battle.room_uid}</Text>
                            <Text className="text-xs text-muted-foreground">底分 {battle.base_score}</Text>
                          </View>
                          {battle.fee > 0 && (
                            <Text className="text-xs text-orange-600">费用 {battle.fee}</Text>
                          )}
                        </View>
                      </View>
                    );
                  })}
                </View>
              )}
            </InfoCardContent>
          </InfoCard>

          {/* 加载更多按钮 */}
          {battles.length > 0 && battles.length < total && (
            <Button
              variant="outline"
              onPress={onLoadMore}
              disabled={loadingBattles}
            >
              <Text>{loadingBattles ? '加载中...' : '加载更多'}</Text>
            </Button>
          )}
        </View>
      </ScrollView>
    </View>
  );
};

