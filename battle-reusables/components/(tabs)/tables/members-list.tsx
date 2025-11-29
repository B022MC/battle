import React, { useState } from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loading } from '@/components/shared/loading';
import { CreditDialog } from '@/components/(shop)/members/credit-dialog';
import { shopsMembersUpdateRemark } from '@/services/shops/members';
import { useRequest } from '@/hooks/use-request';

type MembersListProps = {
  loading: boolean;
  data?: API.ShopsMemberItem[];
  houseGid?: number;
  myGroupId?: number; // 店铺管理员的圈子ID，用于判断是否可以上分下分
  onPullToGroup?: (gamePlayerID: string, memberName: string, currentGroupName?: string) => void;
  onRemoveFromGroup?: (gamePlayerID: string, memberName: string, currentGroupName: string) => void;
  onCreditChange?: () => void; // 上分/下分后的回调
};

export const MembersList = ({ loading, data, houseGid, myGroupId, onPullToGroup, onRemoveFromGroup, onCreditChange }: MembersListProps) => {
  const [creditDialog, setCreditDialog] = useState<{ visible: boolean; type: 'deposit' | 'withdraw'; memberId: number; memberName: string } | null>(null);
  const [editingRemarkId, setEditingRemarkId] = useState<string | null>(null);
  const [remarkValues, setRemarkValues] = useState<Record<string, string>>({});
  const { run: updateRemarkRun } = useRequest(shopsMembersUpdateRemark, { manual: true });

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
        <Text className="text-sm text-muted-foreground">共 {data.length} 人</Text>
      </View>
      <FlatList
        data={data}
        renderItem={({ item }) => (
          <Card className="mb-2 p-3">
            <View className="gap-2">
              {/* 游戏信息区 */}
              <View className="flex-row items-center justify-between">
                <View className="flex-1">
                  <Text className="font-medium">{item.nick_name || '未命名'}</Text>
                  <View className="mt-1 flex-row gap-2">
                    <Text className="text-xs text-muted-foreground">
                      GameID: {item.game_id}
                    </Text>
                    <Text className="text-xs text-muted-foreground">
                      MemberID: {item.member_id}
                    </Text>
                  </View>
                  {/* 当前圈子状态 */}
                  <View className="mt-1 flex-row items-center gap-2">
                    {item.current_group_name ? (
                      <View className="flex-row items-center gap-1">
                        <View className="h-2 w-2 rounded-full bg-blue-500" />
                        <Text className="text-xs text-blue-600 dark:text-blue-400">
                          {item.current_group_name}
                        </Text>
                      </View>
                    ) : (
                      <View className="flex-row items-center gap-1">
                        <View className="h-2 w-2 rounded-full bg-orange-500" />
                        <Text className="text-xs text-orange-600 dark:text-orange-400">
                          无圈子
                        </Text>
                      </View>
                    )}
                  </View>
                  {/* 备注显示 */}
                  {item.remark && (
                    <View className="mt-1">
                      <Text className="text-xs text-muted-foreground">
                        💬 {item.remark}
                      </Text>
                    </View>
                  )}
                </View>
                <View className="ml-2">
                  {item.member_type === 2 && (
                    <View className="rounded-md bg-primary px-2 py-1">
                      <Text className="text-xs text-primary-foreground">管理员</Text>
                    </View>
                  )}
                  {item.member_type === 0 && (
                    <View className="rounded-md bg-secondary px-2 py-1">
                      <Text className="text-xs text-secondary-foreground">普通成员</Text>
                    </View>
                  )}
                </View>
              </View>
              
              {/* 平台用户关联信息 */}
              {item.is_bind_platform && item.platform_user ? (
                <View className="mt-2 border-t border-border pt-2">
                  <View className="flex-row items-center justify-between">
                    <View className="flex-1">
                      <View className="flex-row items-center gap-2">
                        <View className="rounded-full bg-green-500/20 px-2 py-0.5">
                          <Text className="text-xs text-green-700 dark:text-green-400">
                            已绑定
                          </Text>
                        </View>
                        <Text className="font-medium text-sm">
                          {item.platform_user.nick_name || item.platform_user.username}
                        </Text>
                      </View>
                      <Text className="mt-1 text-xs text-muted-foreground">
                        用户名: {item.platform_user.username}
                      </Text>
                    </View>
                  </View>
                </View>
              ) : (
                <View className="mt-2 border-t border-border pt-2">
                  <View className="flex-row items-center gap-2">
                    <View className="rounded-full bg-orange-500/20 px-2 py-0.5">
                      <Text className="text-xs text-orange-700 dark:text-orange-400">
                        暂未绑定
                      </Text>
                    </View>
                    <Text className="text-xs text-muted-foreground">
                      该游戏账号尚未绑定平台用户
                    </Text>
                  </View>
                </View>
              )}

              {/* 备注编辑区 - 独立显示 */}
              {item.game_player_id && houseGid && (
                <View className="mt-2 border-t border-border pt-2">
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
                        <Text className="text-xs">
                          {item.remark ? '✏️ 编辑备注' : '➕ 添加备注'}
                        </Text>
                      </Button>
                    </View>
                  )}
                </View>
              )}

              {/* 操作按钮区 */}
              {item.game_player_id && (onPullToGroup || onRemoveFromGroup || houseGid) && (
                <View className="mt-2 border-t border-border pt-2">
                  {/* 圈子管理按钮 */}
                  {(onPullToGroup || onRemoveFromGroup) && (
                    <View className="flex-row gap-2 mb-2">
                      {onPullToGroup && (
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
                          <Text className="text-xs">
                            {item.current_group_name ? '转移圈子' : '拉入圈子'}
                          </Text>
                        </Button>
                      )}
                      {onRemoveFromGroup && item.current_group_name && (
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
                      )}
                    </View>
                  )}
                  {/* 资金操作按钮 - 只对自己圈子的成员显示 */}
                  {houseGid && item.member_id && myGroupId && item.current_group_id === myGroupId && (
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
                  )}
                </View>
              )}
            </View>
          </Card>
        )}
        keyExtractor={(item) => `${item.user_id}-${item.game_id}-${item.member_id}`}
        scrollEnabled={false}
      />
      
      {/* 上分/下分对话框 */}
      {creditDialog && houseGid && (
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
      )}
    </View>
  );
};
