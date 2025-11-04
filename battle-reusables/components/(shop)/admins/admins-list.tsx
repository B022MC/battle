import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { AdminsItem } from './admins-item';
import { Loading } from '@/components/shared/loading';

type AdminsListProps = {
  houseId?: number;
  loading?: boolean;
  data?: API.ShopsAdminsListResult;
};

export const AdminsList = ({ houseId, loading, data }: AdminsListProps) => {
  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无管理员数据</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data}
      renderItem={({ item }) => <AdminsItem houseId={houseId} data={item} />}
      keyExtractor={(item) => `${item.id}`}
    />
  );
};

