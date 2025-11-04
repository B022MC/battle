import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { MembersItem } from './members-item';
import { Loading } from '@/components/shared/loading';

type MembersListProps = {
  houseId?: number;
  loading?: boolean;
  data?: API.ShopsMemberItem[];
};

export const MembersList = ({ houseId, loading, data }: MembersListProps) => {
  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无成员数据</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data}
      renderItem={({ item }) => <MembersItem houseId={houseId} data={item} />}
      keyExtractor={(item) => `${item.member_id}-${item.user_id}`}
    />
  );
};
