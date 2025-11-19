import React from 'react';
import { View, TouchableOpacity } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Users, ChevronRight } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';

interface OnlineMember {
  id: number;
  nickname: string;
  is_online: boolean;
  last_active_at?: string;
}

interface MembersOnlineCardProps {
  totalMembers: number;
  onlineMembers: OnlineMember[];
  loading?: boolean;
  onViewAll?: () => void;
}

export function MembersOnlineCard({
  totalMembers,
  onlineMembers,
  loading = false,
  onViewAll,
}: MembersOnlineCardProps) {
  const onlineCount = onlineMembers.filter(m => m.is_online).length;

  if (loading) {
    return (
      <Card className="p-4">
        <View className="flex-row items-center gap-2 mb-4">
          <Icon as={Users} className="size-5 text-green-600" />
          <Text className="text-lg font-semibold">成员在线情况</Text>
        </View>
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">加载中...</Text>
        </View>
      </Card>
    );
  }

  return (
    <Card className="p-4">
      {/* 标题栏 */}
      <View className="flex-row items-center justify-between mb-4">
        <View className="flex-row items-center gap-2">
          <Icon as={Users} className="size-5 text-green-600" />
          <Text className="text-lg font-semibold">成员在线情况</Text>
        </View>
        {onViewAll && (
          <TouchableOpacity
            onPress={onViewAll}
            className="flex-row items-center gap-1"
          >
            <Text className="text-sm text-blue-600">查看全部</Text>
            <Icon as={ChevronRight} className="size-4 text-blue-600" />
          </TouchableOpacity>
        )}
      </View>

      {/* 在线统计 */}
      <View className="flex-row items-center gap-3 mb-4 p-3 bg-green-50 dark:bg-green-950 rounded-lg">
        <View className="flex-1">
          <Text className="text-xs text-muted-foreground">在线人数</Text>
          <Text className="text-2xl font-bold text-green-600">
            {onlineCount}
          </Text>
        </View>
        <View className="h-8 w-px bg-border" />
        <View className="flex-1">
          <Text className="text-xs text-muted-foreground">总人数</Text>
          <Text className="text-2xl font-bold text-foreground">
            {totalMembers}
          </Text>
        </View>
        <View className="h-8 w-px bg-border" />
        <View className="flex-1">
          <Text className="text-xs text-muted-foreground">在线率</Text>
          <Text className="text-2xl font-bold text-blue-600">
            {totalMembers > 0 ? ((onlineCount / totalMembers) * 100).toFixed(0) : 0}%
          </Text>
        </View>
      </View>

      {/* 在线成员列表 */}
      {onlineMembers.length === 0 ? (
        <View className="py-4 items-center">
          <Text className="text-sm text-muted-foreground">暂无在线成员</Text>
        </View>
      ) : (
        <View className="gap-2">
          <Text className="text-xs font-medium text-muted-foreground mb-1">
            最近在线成员
          </Text>
          <View className="flex-row flex-wrap gap-2">
            {onlineMembers.slice(0, 10).map((member) => (
              <View
                key={member.id}
                className="flex-row items-center gap-2 px-3 py-2 bg-secondary rounded-full"
              >
                <View
                  className={`size-2 rounded-full ${
                    member.is_online ? 'bg-green-500' : 'bg-gray-400'
                  }`}
                />
                <Text className="text-sm">{member.nickname}</Text>
              </View>
            ))}
          </View>
          {onlineMembers.length > 10 && (
            <Text className="text-xs text-muted-foreground text-center mt-2">
              还有 {onlineMembers.length - 10} 位成员...
            </Text>
          )}
        </View>
      )}
    </Card>
  );
}
