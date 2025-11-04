import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { FundsItem } from './funds-item';
import { Loading } from '@/components/shared/loading';

type FundsListProps = {
  houseId?: number;
  loading?: boolean;
  data?: API.MembersWalletItem[];
};

export const FundsList = ({ houseId, loading, data }: FundsListProps) => {
  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无数据</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data}
      renderItem={({ item }) => <FundsItem houseId={houseId} data={item} />}
      keyExtractor={(item) => `${item.member_id}-${item.house_gid}`}
    />
  );
};

