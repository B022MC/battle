import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { CtrlAccountsItem } from './ctrl-accounts-item';
import { Loading } from '@/components/shared/loading';

type CtrlAccountsListProps = {
  houseId?: number;
  loading?: boolean;
  data?: API.ShopsCtrlAccountsListResult;
};

export const CtrlAccountsList = ({ houseId, loading, data }: CtrlAccountsListProps) => {
  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无中控账号数据</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data}
      renderItem={({ item }) => <CtrlAccountsItem houseId={houseId} data={item} />}
      keyExtractor={(item) => `${item.id}`}
    />
  );
};

