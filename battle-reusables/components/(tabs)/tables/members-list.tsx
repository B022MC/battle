import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Loading } from '@/components/shared/loading';

type MembersListProps = {
  loading?: boolean;
  data?: API.ShopsMemberItem[];
};

export const MembersList = ({ loading, data }: MembersListProps) => {
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
                  {item.group_name && (
                    <Text className="mt-1 text-xs text-muted-foreground">
                      圈子: {item.group_name}
                    </Text>
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
                        用户名: {item.platform_user.username} | 角色: {item.platform_user.role}
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
            </View>
          </Card>
        )}
        keyExtractor={(item) => `${item.user_id}-${item.game_id}-${item.member_id}`}
        scrollEnabled={false}
      />
    </View>
  );
};
