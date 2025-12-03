import React, { useState } from 'react';
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
  const [creditDialog, setCreditDialog] = useState<{ visible: boolean; type: 'deposit' | 'withdraw'; memberId: number; memberName: string } | null>(null);
  const [editingRemarkId, setEditingRemarkId] = useState<string | null>(null);
  const [remarkValues, setRemarkValues] = useState<Record<string, string>>({});
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
      <View className="mb-2 flex-row items-center justify-between">
        <Text className="text-lg font-semibold">游戏成员列表</Text>
        <Text className="text-sm text-muted-foreground">{`共 ${data.length} 人`}</Text>
      </View>
      {data.map((item) => (
        <Card key={`${item.user_id}-${item.game_id}-${item.member_id}`} className="mb-2 p-3">
          <View className="gap-2">
            <View className="flex-row items-center justify-between">
              <View className="flex-1">
                <View className="flex-row items-center">
                  <Text className="font-medium">{item.nick_name || '未命名'}</Text>
                  {item.forbid ? (
                    <View className="bg-destructive px-2 py-0.5 rounded ml-2">
                      <Text className="text-destructive-foreground text-xs">已禁用</Text>
                    </View>
                  ) : null}
                </View>
                <View className="mt-1 flex-row gap-2">
                  <Text className="text-xs text-muted-foreground">{`GameID: ${item.game_id}`}</Text>
                  <Text className="text-xs text-muted-foreground">{`MemberID: ${item.member_id}`}</Text>
                </View>
                <View className="mt-1 flex-row items-center gap-2">
                  {item.current_group_name ? (
                    <View className="flex-row items-center gap-1">
                      <View className="h-2 w-2 rounded-full bg-blue-500" />
                      <Text className="text-xs text-blue-600 dark:text-blue-400">{item.current_group_name}</Text>
                    </View>
                  ) : (
                    <View className="flex-row items-center gap-1">
                      <View className="h-2 w-2 rounded-full bg-orange-500" />
                      <Text className="text-xs text-orange-600 dark:text-orange-400">无圈子</Text>
                    </View>
                  )}
                </View>
                {item.remark ? (
                  <View className="mt-1">
                    <Text className="text-xs text-muted-foreground">{`💬 ${item.remark}`}</Text>
                  </View>
                ) : null}
              </View>
              <View className="ml-2">
                {item.member_type === 2 ? (
                  <View className="rounded-md bg-primary px-2 py-1">
                    <Text className="text-xs text-primary-foreground">管理员</Text>
                  </View>
                ) : null}
                {item.member_type === 0 ? (
                  <View className="rounded-md bg-secondary px-2 py-1">
                    <Text className="text-xs text-secondary-foreground">普通成员</Text>
                  </View>
                ) : null}
              </View>
            </View>
            {item.is_bind_platform && item.platform_user ? (
              <View className="mt-2 border-t border-border pt-2">
                <View className="flex-row items-center justify-between">
                  <View className="flex-1">
                    <View className="flex-row items-center gap-2">
                      <View className="rounded-full bg-green-500/20 px-2 py-0.5">
                        <Text className="text-xs text-green-700 dark:text-green-400">已绑定</Text>
                      </View>
                      <Text className="font-medium text-sm">{item.platform_user.nick_name || item.platform_user.username}</Text>
                    </View>
                    <Text className="mt-1 text-xs text-muted-foreground">{`用户名: ${item.platform_user.username}`}</Text>
                  </View>
                </View>
              </View>
            ) : (
              <View className="mt-2 border-t border-border pt-2">
                <View className="flex-row items-center gap-2">
                  <View className="rounded-full bg-orange-500/20 px-2 py-0.5">
                    <Text className="text-xs text-orange-700 dark:text-orange-400">暂未绑定</Text>
                  </View>
                  <Text className="text-xs text-muted-foreground">该游戏账号尚未绑定平台用户</Text>
                </View>
              </View>
            )}
            {item.game_player_id && houseGid ? (
              <View className="mt-2 border-t border-border pt-2">
                <View className="mb-2 flex-row gap-2">
                  {isBlockedList ? (
                    <Button
                      variant="default" 
                      size="sm"
                      className="flex-1 bg-green-600"
                      onPress={() => handleUnforbid(item.game_player_id!, item.nick_name || '未命名')}
                    >
                      <Text className="text-xs">✅ 解禁成员</Text>
                    </Button>
                  ) : (
                    !item.forbid ? (
                      <Button
                        variant="destructive"
                        size="sm"
                        className="flex-1"
                        onPress={() => handleForbid(item.game_player_id!, item.nick_name || '未命名')}
                      >
                        <Text className="text-xs">🚫 禁用成员</Text>
                      </Button>
                    ) : null
                  )}
                </View>
                {editingRemarkId === item.game_player_id ? (
                  <View className="mb-2 gap-2">
                    <Input
                      value={remarkValues[item.game_player_id] || ''}
                      onChangeText={(text) => setRemarkValues(prev => ({ ...prev, [item.game_player_id!]: text }))}
                      placeholder="输入备注"
                    />
                    <View className="flex-row gap-2">
                      <Button
                        variant="default"
                        size="sm"
                        className="flex-1"
                        onPress={() => handleSaveRemark(item.game_player_id!)}
                      >
                        <Text className="text-xs">保存备注</Text>
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onPress={handleCancelRemark}
                      >
                        <Text className="text-xs">取消</Text>
                      </Button>
                    </View>
                  </View>
                ) : (
                  <View className="mb-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onPress={() => handleEditRemark(item.game_player_id!, item.remark || '')}
                    >
                      <Text className="text-xs">{item.remark ? '✏️ 编辑备注' : '➕ 添加备注'}</Text>
                    </Button>
                  </View>
                )}
              </View>
            ) : null}
            {item.game_player_id && item.game_id && myGroupId && item.current_group_id === myGroupId ? (
              <View className="mt-2 border-t border-border pt-2">
                <View className="mb-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onPress={() => toggleBattleExpand(item.game_player_id!, item.game_id!)}
                      disabled={loadingBattles.has(item.game_player_id!)}
                    >
                      <Text className="text-xs">{loadingBattles.has(item.game_player_id!) ? '📊 加载中...' : expandedBattleIds.has(item.game_player_id!) ? '📊 收起战绩' : '📊 查看战绩'}</Text>
                    </Button>
                </View>
                {expandedBattleIds.has(item.game_player_id!) && battleRecords[item.game_player_id!] ? (
                  <View className="gap-2">
                    {battleRecords[item.game_player_id!].length === 0 ? (
                      <View className="py-4 items-center">
                        <Text className="text-xs text-muted-foreground">暂无战绩记录</Text>
                      </View>
                    ) : (
                      battleRecords[item.game_player_id!].map((record) => (
                        <View key={record.id} className="bg-secondary/30 rounded-md p-2">
                          <View className="flex-row items-center justify-between mb-1">
                            <Text className="text-xs font-medium">{formatTime(record.battle_at)}</Text>
                            <View className="flex-row items-center gap-1">
                              <Text className={`text-xs font-bold ${record.score >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>{`${record.score >= 0 ? '+' : ''}${formatScore(record.score)}`}</Text>
                            </View>
                          </View>
                          <View className="flex-row items-center justify-between">
                            <Text className="text-xs text-muted-foreground">{`房间 ${record.room_uid} · 底分 ${record.base_score}`}</Text>
                            <Text className="text-xs text-muted-foreground">{`余额 ${formatScore(record.player_balance)}`}</Text>
                          </View>
                          {record.fee > 0 ? (
                            <Text className="text-xs text-muted-foreground mt-1">{`手续费 -${formatScore(record.fee)}`}</Text>
                          ) : null}
                        </View>
                      ))
                    )}
                  </View>
                ) : null}
              </View>
            ) : null}
            {item.game_player_id && (onPullToGroup || onRemoveFromGroup || houseGid) ? (
              <View className="mt-2 border-t border-border pt-2">
                {(onPullToGroup || onRemoveFromGroup) ? (
                  <View className="flex-row gap-2 mb-2">
                    {onPullToGroup ? (
                      <Button
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onPress={() => onPullToGroup(
                          item.game_player_id!,
                          item.nick_name || '未命名',
                          item.current_group_name
                        )}
                      >
                        <Text className="text-xs">{item.current_group_name ? '转移圈子' : '拉入圈子'}</Text>
                      </Button>
                    ) : null}
                    {onRemoveFromGroup && item.current_group_name ? (
                      <Button
                        variant="destructive"
                        size="sm"
                        className="flex-1"
                        onPress={() => onRemoveFromGroup(
                          item.game_player_id!,
                          item.nick_name || '未命名',
                          item.current_group_name!
                        )}
                      >
                        <Text className="text-xs">踢出圈子</Text>
                      </Button>
                    ) : null}
                  </View>
                ) : null}
                {houseGid && item.member_id && myGroupId && item.current_group_id === myGroupId ? (
                  <View className="flex-row gap-2">
                    <Button
                      variant="default"
                      size="sm"
                      className="flex-1"
                      onPress={() => setCreditDialog({
                        visible: true,
                        type: 'deposit',
                        memberId: item.member_id!,
                        memberName: item.nick_name || '未命名'
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
                        memberId: item.member_id!,
                        memberName: item.nick_name || '未命名'
                      })}
                    >
                      <Text className="text-xs">下分</Text>
                    </Button>
                  </View>
                ) : null}
              </View>
            ) : null}
          </View>
        </Card>
      ))}
      {creditDialog && houseGid ? (
        <CreditDialog
          visible={creditDialog.visible}
          type={creditDialog.type}
          houseGid={houseGid}
          memberId={creditDialog.memberId}
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
