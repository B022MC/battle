import React, { useState, useRef } from 'react';
import { ScrollView, View } from 'react-native';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { alert } from '@/utils/alert';
import { listGroupBattles, getGroupStats } from '@/services/battles/query';
import type { BattleRecord, GroupStats } from '@/services/battles/query-typing';
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

export const GroupBattlesView = () => {
  const { isStoreAdmin } = usePermission();

  // 状态管理
  const [houseGid, setHouseGid] = useState<string>('');
  const [groupId, setGroupId] = useState<string>('');
  const [playerGameId, setPlayerGameId] = useState<string>('');
  const [page, setPage] = useState(1);
  const [startTime, setStartTime] = useState<string>('');
  const [endTime, setEndTime] = useState<string>('');

  // API 请求
  const { data: battlesData, loading: loadingBattles, run: runListBattles } = useRequest(listGroupBattles, { manual: true });
  const { data: statsData, loading: loadingStats, run: runGetStats } = useRequest(getGroupStats, { manual: true });
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
      alert.show({ title: '请输入店铺号', variant: 'error' });
      return;
    }
    if (!groupId) {
      alert.show({ title: '请输入圈子ID', variant: 'error' });
      return;
    }

    const gid = Number(houseGid);
    const gidNum = Number(groupId);
    if (isNaN(gid) || gid <= 0 || isNaN(gidNum) || gidNum <= 0) {
      alert.show({ title: '店铺号或圈子ID格式错误', variant: 'error' });
      return;
    }

    const params: any = {
      house_gid: gid,
      group_id: gidNum,
      page,
      size: 20,
    };

    if (playerGameId) {
      const pid = Number(playerGameId);
      if (!isNaN(pid) && pid > 0) {
        params.player_game_id = pid;
      }
    }

    if (startTime) {
      const timestamp = new Date(startTime).getTime() / 1000;
      if (!isNaN(timestamp)) {
        params.start_time = Math.floor(timestamp);
      }
    }

    if (endTime) {
      const timestamp = new Date(endTime).getTime() / 1000;
      if (!isNaN(timestamp)) {
        params.end_time = Math.floor(timestamp);
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
      alert.show({ title: '请输入店铺号', variant: 'error' });
      return;
    }
    if (!groupId) {
      alert.show({ title: '请输入圈子ID', variant: 'error' });
      return;
    }

    const gid = Number(houseGid);
    const gidNum = Number(groupId);
    if (isNaN(gid) || gid <= 0 || isNaN(gidNum) || gidNum <= 0) {
      alert.show({ title: '店铺号或圈子ID格式错误', variant: 'error' });
      return;
    }

    const params: any = {
      house_gid: gid,
      group_id: gidNum,
    };

    if (startTime) {
      const timestamp = new Date(startTime).getTime() / 1000;
      if (!isNaN(timestamp)) {
        params.start_time = Math.floor(timestamp);
      }
    }

    if (endTime) {
      const timestamp = new Date(endTime).getTime() / 1000;
      if (!isNaN(timestamp)) {
        params.end_time = Math.floor(timestamp);
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

  // 解析 players_json 获取对战信息
  const parsePlayersInfo = (playersJson: string, currentPlayerId?: number) => {
    try {
      const players = JSON.parse(playersJson) as Array<{ UserGameID: number; Score: number }>;
      if (!players || players.length === 0) return null;
      
      const currentPlayer = players.find(p => p.UserGameID === currentPlayerId);
      const opponents = players.filter(p => p.UserGameID !== currentPlayerId);
      
      return {
        players,
        currentPlayer,
        opponents,
        playerCount: players.length,
      };
    } catch {
      return null;
    }
  };

  // 获取胜负状态
  const getBattleResult = (score: number) => {
    if (score > 0) return { text: '胜', color: 'text-green-600', bg: 'bg-green-100' };
    if (score < 0) return { text: '负', color: 'text-red-600', bg: 'bg-red-100' };
    return { text: '平', color: 'text-gray-600', bg: 'bg-gray-100' };
  };

  if (!isStoreAdmin) {
    return (
      <View className="flex-1 bg-background p-4 justify-center items-center">
        <Text className="text-lg text-muted-foreground">
          仅店铺管理员可访问此页面
        </Text>
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-background p-4">
      {/* 查询表单 */}
      <Card className="p-4 mb-4">
        <Text className="text-lg font-bold mb-4">查询圈子战绩</Text>
        
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
          <Text className="mb-2">圈子 *</Text>
          <Select
            value={groupId ? ({ label: groupOptions?.find(g => String(g.id) === groupId)?.name || `圈子 ${groupId}`, value: groupId } as any) : undefined}
            onValueChange={(opt) => setGroupId(String(opt?.value ?? ''))}
          >
            <SelectTrigger ref={groupRef} onTouchStart={onGroupTouchStart} className="min-w-[160px]">
              <SelectValue placeholder={groupId ? `圈子 ${groupId}` : '选择圈子'} />
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
          <Text className="mb-2">玩家游戏ID (可选)</Text>
          <Input
            placeholder="不填则查询所有成员"
            value={playerGameId}
            onChangeText={setPlayerGameId}
            keyboardType="numeric"
          />
        </View>

        <View className="mb-3">
          <Text className="mb-2">开始时间 (可选)</Text>
          <Input
            placeholder="YYYY-MM-DD HH:mm:ss"
            value={startTime}
            onChangeText={setStartTime}
          />
        </View>

        <View className="mb-3">
          <Text className="mb-2">结束时间 (可选)</Text>
          <Input
            placeholder="YYYY-MM-DD HH:mm:ss"
            value={endTime}
            onChangeText={setEndTime}
          />
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
          <Text className="text-lg font-bold mb-3">圈子统计</Text>
          <Text className="mb-2">圈子: {statsData.group_name}</Text>
          <Text className="mb-2">总局数: {statsData.total_games}</Text>
          <Text className="mb-2">总分数: {formatScore(statsData.total_score)}</Text>
          <Text className="mb-2">总费用: {statsData.total_fee}</Text>
          <Text className="mb-2">活跃成员: {statsData.active_members}</Text>
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
              {battlesData.list.map((battle: BattleRecord) => {
                const playersInfo = parsePlayersInfo(battle.players_json, battle.player_game_id);
                const result = getBattleResult(battle.score);
                
                return (
                  <View
                    key={battle.id}
                    className="border-b border-border py-3 last:border-b-0"
                  >
                    {/* 第一行：玩家信息 + 胜负标签 + 得分 */}
                    <View className="flex-row justify-between items-center mb-2">
                      <View className="flex-row items-center gap-2">
                        <Text className="font-medium">
                          {battle.player_game_id}
                        </Text>
                        <View className={`px-2 py-0.5 rounded ${result.bg}`}>
                          <Text className={`text-xs font-bold ${result.color}`}>
                            {result.text}
                          </Text>
                        </View>
                      </View>
                      <Text
                        className={`text-xl font-bold ${result.color}`}
                      >
                        {formatScore(battle.score)}
                      </Text>
                    </View>

                    {/* 第二行：对战详情 */}
                    {playersInfo && playersInfo.opponents.length > 0 && (
                      <View className="bg-secondary/50 rounded-lg p-2 mb-2">
                        <Text className="text-xs text-muted-foreground mb-1">对战详情</Text>
                        <View className="flex-row flex-wrap gap-2">
                          {playersInfo.players.map((p, idx) => (
                            <View 
                              key={p.UserGameID} 
                              className={`flex-row items-center gap-1 px-2 py-1 rounded ${
                                p.UserGameID === battle.player_game_id 
                                  ? 'bg-primary/10 border border-primary/30' 
                                  : 'bg-background'
                              }`}
                            >
                              <Text className={`text-xs ${p.UserGameID === battle.player_game_id ? 'font-bold' : ''}`}>
                                {p.UserGameID === battle.player_game_id ? '我' : `对手${playersInfo.opponents.length > 1 ? idx : ''}`}
                              </Text>
                              <Text 
                                className={`text-xs font-medium ${
                                  p.Score > 0 ? 'text-green-600' : p.Score < 0 ? 'text-red-600' : 'text-muted-foreground'
                                }`}
                              >
                                {formatScore(p.Score)}
                              </Text>
                            </View>
                          ))}
                        </View>
                      </View>
                    )}

                    {/* 第三行：时间、房间、费用等信息 */}
                    <View className="flex-row justify-between items-center">
                      <View className="flex-row gap-3">
                        <Text className="text-xs text-muted-foreground">
                          {formatDate(battle.battle_at)}
                        </Text>
                        <Text className="text-xs text-muted-foreground">
                          房间 {battle.room_uid}
                        </Text>
                      </View>
                      <View className="flex-row gap-3">
                        {battle.fee > 0 && (
                          <Text className="text-xs text-orange-600">
                            费用 {battle.fee}
                          </Text>
                        )}
                        <Text className="text-xs text-muted-foreground">
                          余额 {battle.player_balance}
                        </Text>
                      </View>
                    </View>
                  </View>
                );
              })}

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

