import React, { useState, useMemo } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loading } from '@/components/shared/loading';
import { CreditDialog } from '@/components/(shop)/members/credit-dialog';
import { shopsMembersUpdateRemark, shopsMembersForbid, shopsMembersUnforbid } from '@/services/shops/members';
import { useRequest } from '@/hooks/use-request';
import { listGroupBattles } from '@/services/battles/query';
import type { BattleRecord } from '@/services/battles/query-typing';
import { alert } from '@/utils/alert';

type MembersListProps = {
  loading: boolean;
  data?: API.ShopsMemberItem[];
  houseGid?: number;
  myGroupId?: number; // 店铺管理员的圈子ID，用于判断是否可以上分下分
  onPullToGroup?: (gamePlayerID: string, memberName: string, currentGroupName?: string) => void;
  onRemoveFromGroup?: (gamePlayerID: string, memberName: string, currentGroupName: string) => void;
  onCreditChange?: () => void; // 上分/下分后的回调，也可用于刷新列表
  isBlockedList?: boolean; // 是否为禁用名单模式
};

export const MembersList = ({ loading, data, houseGid, myGroupId, onPullToGroup, onRemoveFromGroup, onCreditChange, isBlockedList }: MembersListProps) => {
  const [creditDialog, setCreditDialog] = useState<{ visible: boolean; type: 'deposit' | 'withdraw'; gameId: number; memberName: string } | null>(null);
  const [editingRemarkId, setEditingRemarkId] = useState<string | null>(null);
  const [remarkValues, setRemarkValues] = useState<Record<string, string>>({});
  const [searchText, setSearchText] = useState('');
  const [sortBy, setSortBy] = useState<'default' | 'gameId' | 'name'>('default');

  // 搜索和排序后的数据
  const filteredData = useMemo(() => {
    if (!data) return [];
    let result = [...data];
    // 搜索过滤
    if (searchText.trim()) {
      const keyword = searchText.trim().toLowerCase();
      result = result.filter(item => 
        (item.nick_name?.toLowerCase().includes(keyword)) ||
        (item.game_id?.toString().includes(keyword)) ||
        (item.game_player_id?.includes(keyword)) ||
        (item.remark?.toLowerCase().includes(keyword)) ||
        (item.current_group_name?.toLowerCase().includes(keyword))
      );
    }
    // 排序
    if (sortBy === 'gameId') {
      result.sort((a, b) => (b.game_id || 0) - (a.game_id || 0)); // GameID 降序（新用户在前）
    } else if (sortBy === 'name') {
      result.sort((a, b) => (a.nick_name || '').localeCompare(b.nick_name || ''));
    }
    return result;
  }, [data, searchText, sortBy]);
  const [expandedBattleIds, setExpandedBattleIds] = useState<Set<string>>(new Set());
  const [battleRecords, setBattleRecords] = useState<Record<string, BattleRecord[]>>({});
  const [loadingBattles, setLoadingBattles] = useState<Set<string>>(new Set());
  const { run: updateRemarkRun } = useRequest(shopsMembersUpdateRemark, { manual: true });
  const { run: forbidRun } = useRequest(shopsMembersForbid, { manual: true });
  const { run: unforbidRun } = useRequest(shopsMembersUnforbid, { manual: true });

  // 处理禁用
  const handleForbid = async (gamePlayerId: string, memberName: string) => {
    if (!houseGid) return;
    
    alert.show({
      title: '确认禁用',
      description: `确定要禁用玩家 ${memberName} 吗？禁用后玩家将无法进入游戏。`,
      confirmText: '禁用',
      cancelText: '取消',
      onConfirm: async () => {
        try {
          await forbidRun({ house_gid: houseGid, game_player_id: gamePlayerId });
          alert.show({ title: '已禁用' });
          onCreditChange?.(); // 刷新列表
        } catch (error: any) {
          console.error('禁用失败:', error);
          alert.show({ title: '禁用失败', description: error.message || '未知错误' });
        }
      }
    });
  };

  // 处理解禁
  const handleUnforbid = async (gamePlayerId: string, memberName: string) => {
    if (!houseGid) return;
    
    alert.show({
      title: '确认解禁',
      description: `确定要解禁玩家 ${memberName} 吗？`,
      confirmText: '解禁',
      cancelText: '取消',
      onConfirm: async () => {
        try {
          await unforbidRun({ house_gid: houseGid, game_player_id: gamePlayerId });
          alert.show({ title: '已解禁' });
          onCreditChange?.(); // 刷新列表
        } catch (error: any) {
          console.error('解禁失败:', error);
          alert.show({ title: '解禁失败', description: error.message || '未知错误' });
        }
      }
    });
  };

  // 处理备注编辑
  const handleEditRemark = (gamePlayerId: string, currentRemark: string) => {
    setEditingRemarkId(gamePlayerId);
    setRemarkValues(prev => ({ ...prev, [gamePlayerId]: currentRemark || '' }));
  };

  // 保存备注
  const handleSaveRemark = async (gamePlayerId: string) => {
    if (!houseGid) return;
    const remark = remarkValues[gamePlayerId] || '';
    await updateRemarkRun({ house_gid: houseGid, game_player_id: gamePlayerId, remark });
    setEditingRemarkId(null);
    onCreditChange?.(); // 刷新列表
  };

  // 取消编辑
  const handleCancelRemark = () => {
    setEditingRemarkId(null);
  };

  // 加载战绩
  const loadBattleRecords = async (gamePlayerId: string, gameId: number) => {
    if (!houseGid || !myGroupId) return;
    
    setLoadingBattles(prev => new Set(prev).add(gamePlayerId));
    try {
      const response = await listGroupBattles({
        house_gid: houseGid,
        group_id: myGroupId,
        player_game_id: gameId,
        page: 1,
        size: 10, // 获取最近10条战绩
      });
      
      if (response?.data?.list) {
        setBattleRecords(prev => ({ ...prev, [gamePlayerId]: response.data?.list || [] }));
      }
    } catch (error) {
      console.error('加载战绩失败:', error);
    } finally {
      setLoadingBattles(prev => {
        const newSet = new Set(prev);
        newSet.delete(gamePlayerId);
        return newSet;
      });
    }
  };

  // 切换战绩展开/收起
  const toggleBattleExpand = async (gamePlayerId: string, gameId: number) => {
    const isExpanded = expandedBattleIds.has(gamePlayerId);
    
    if (isExpanded) {
      // 收起
      setExpandedBattleIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(gamePlayerId);
        return newSet;
      });
    } else {
      // 展开，如果还没有加载过战绩则先加载
      if (!battleRecords[gamePlayerId]) {
        await loadBattleRecords(gamePlayerId, gameId);
      }
      setExpandedBattleIds(prev => new Set(prev).add(gamePlayerId));
    }
  };

  // 格式化时间
  const formatTime = (timeStr: string) => {
    const date = new Date(timeStr);
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // 格式化分数（分转元）
  const formatScore = (score: number) => {
    return (score / 100).toFixed(2);
  };

  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无成员数据</Text>
      </View>
    );
  }

  return (
    <View>
      {/* 搜索和排序 */}
      <View className="mb-2 gap-2">
        <Input
          value={searchText}
          onChangeText={setSearchText}
          placeholder="搜索昵称/ID/备注/圈子..."
          className="h-8 text-sm"
        />
        <View className="flex-row items-center justify-between">
          <View className="flex-row items-center gap-1">
            <Text className="text-xs text-muted-foreground">排序:</Text>
            <Button variant={sortBy === 'default' ? 'default' : 'ghost'} size="sm" className="h-6 px-2" onPress={() => setSortBy('default')}><Text className="text-xs">默认</Text></Button>
            <Button variant={sortBy === 'gameId' ? 'default' : 'ghost'} size="sm" className="h-6 px-2" onPress={() => setSortBy('gameId')}><Text className="text-xs">最新</Text></Button>
            <Button variant={sortBy === 'name' ? 'default' : 'ghost'} size="sm" className="h-6 px-2" onPress={() => setSortBy('name')}><Text className="text-xs">昵称</Text></Button>
          </View>
          <Text className="text-xs text-muted-foreground">{searchText ? `${filteredData.length}/${data.length}` : `共 ${data.length} 人`} {myGroupId ? `(我的圈:${myGroupId})` : '(无圈子)'}</Text>
        </View>
      </View>
      {filteredData.length === 0 && searchText ? (
        <View className="min-h-16 flex-row items-center justify-center">
          <Text className="text-muted-foreground">未找到匹配的成员</Text>
        </View>
      ) : null}
      {filteredData.map((item) => (
        <Card key={`${item.user_id}-${item.game_id}-${item.member_id}`} className="mb-1.5 p-2">
          {/* 第一行：基本信息 */}
          <View className="flex-row items-center justify-between">
            <View className="flex-1 flex-row items-center gap-2">
              <Text className="font-medium text-sm">{item.nick_name || '未命名'}</Text>
              {item.forbid && <View className="bg-destructive px-1.5 py-0.5 rounded"><Text className="text-destructive-foreground text-xs">禁</Text></View>}
              {item.member_type === 2 && <View className="bg-primary px-1.5 py-0.5 rounded"><Text className="text-primary-foreground text-xs">管</Text></View>}
              {item.is_bind_platform && <View className="bg-green-500/20 px-1.5 py-0.5 rounded"><Text className="text-green-700 dark:text-green-400 text-xs">绑</Text></View>}
              {item.current_group_name && <Text className="text-xs text-blue-600 dark:text-blue-400">{item.current_group_name}</Text>}
            </View>
            <Text className="text-xs text-muted-foreground">{item.game_id} {item.current_group_id ? `(圈:${item.current_group_id})` : ''}</Text>
          </View>
          {/* 第二行：备注和操作按钮 */}
          {item.game_player_id && houseGid ? (
            <View className="mt-1.5 flex-row items-center justify-between">
              <View className="flex-1">
                {editingRemarkId === item.game_player_id ? (
                  <View className="flex-row items-center gap-1">
                    <Input
                      value={remarkValues[item.game_player_id] || ''}
                      onChangeText={(text) => setRemarkValues(prev => ({ ...prev, [item.game_player_id!]: text }))}
                      placeholder="备注"
                      className="flex-1 h-7 text-xs"
                    />
                    <Button variant="default" size="sm" className="h-7 px-2" onPress={() => handleSaveRemark(item.game_player_id!)}><Text className="text-xs">保存</Text></Button>
                    <Button variant="ghost" size="sm" className="h-7 px-2" onPress={handleCancelRemark}><Text className="text-xs">取消</Text></Button>
                  </View>
                ) : (
                  <Text className="text-xs text-muted-foreground" onPress={() => handleEditRemark(item.game_player_id!, item.remark || '')}>
                    {item.remark || '点击添加备注'}
                  </Text>
                )}
              </View>
              <View className="flex-row gap-1 ml-2">
                <Button variant="destructive" size="sm" className="h-7 px-2" onPress={() => handleForbid(item.game_player_id!, item.nick_name || '未命名')}><Text className="text-xs">禁用</Text></Button>
                <Button variant="default" size="sm" className="h-7 px-2 bg-green-600" onPress={() => handleUnforbid(item.game_player_id!, item.nick_name || '未命名')}><Text className="text-xs">解禁</Text></Button>
              </View>
            </View>
          ) : null}
          {/* 第三行：圈子操作和上下分 */}
          {item.game_player_id && (onPullToGroup || onRemoveFromGroup || (houseGid && item.member_id)) ? (
            <View className="mt-1.5 flex-row items-center gap-1 flex-wrap">
              {onPullToGroup && (
                <Button variant="outline" size="sm" className="h-7 px-2" onPress={() => onPullToGroup(item.game_player_id!, item.nick_name || '未命名', item.current_group_name)}>
                  <Text className="text-xs">{item.current_group_name ? '转移' : '拉入'}</Text>
                </Button>
              )}
              {onRemoveFromGroup && item.current_group_name && (
                <Button variant="destructive" size="sm" className="h-7 px-2" onPress={() => onRemoveFromGroup(item.game_player_id!, item.nick_name || '未命名', item.current_group_name!)}>
                  <Text className="text-xs">踢出</Text>
                </Button>
              )}
              {houseGid && item.game_id && myGroupId && item.current_group_id === myGroupId && (
                <>
                  <Button variant="default" size="sm" className="h-7 px-2" onPress={() => setCreditDialog({ visible: true, type: 'deposit', gameId: item.game_id!, memberName: item.nick_name || '未命名' })}>
                    <Text className="text-xs">上分</Text>
                  </Button>
                  <Button variant="secondary" size="sm" className="h-7 px-2" onPress={() => setCreditDialog({ visible: true, type: 'withdraw', gameId: item.game_id!, memberName: item.nick_name || '未命名' })}>
                    <Text className="text-xs">下分</Text>
                  </Button>
                  <Text className="text-xs text-primary ml-1" onPress={() => toggleBattleExpand(item.game_player_id!, item.game_id!)}>
                    {loadingBattles.has(item.game_player_id!) ? '...' : expandedBattleIds.has(item.game_player_id!) ? '收起' : '战绩'}
                  </Text>
                </>
              )}
            </View>
          ) : null}
          {/* 战绩展开区域 */}
          {item.game_player_id && expandedBattleIds.has(item.game_player_id) && battleRecords[item.game_player_id] ? (
            <View className="mt-1.5 gap-1">
              {battleRecords[item.game_player_id].length === 0 ? (
                <Text className="text-xs text-muted-foreground py-2 text-center">暂无战绩</Text>
              ) : (
                battleRecords[item.game_player_id].map((record) => (
                  <View key={record.id} className="bg-secondary/30 rounded p-1.5 flex-row items-center justify-between">
                    <Text className="text-xs text-muted-foreground">{formatTime(record.battle_at)}</Text>
                    <Text className={`text-xs font-medium ${record.score >= 0 ? 'text-green-600' : 'text-red-600'}`}>{`${record.score >= 0 ? '+' : ''}${formatScore(record.score)}`}</Text>
                  </View>
                ))
              )}
            </View>
          ) : null}
        </Card>
      ))}
      {creditDialog && houseGid ? (
        <CreditDialog
          visible={creditDialog.visible}
          type={creditDialog.type}
          houseGid={houseGid}
          memberId={creditDialog.gameId}
          memberName={creditDialog.memberName}
          onClose={() => setCreditDialog(null)}
          onSuccess={() => {
            setCreditDialog(null);
            onCreditChange?.();
          }}
        />
      ) : null}
    </View>
  );
};
