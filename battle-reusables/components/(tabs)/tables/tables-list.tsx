import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { TablesItem } from './tables-item';
import { Loading } from '@/components/shared/loading';

type TablesListProps = {
  houseId?: number;
  loading?: boolean;
  data?: API.ShopsTableItem[];
  onChanged?: () => void;
};

export const TablesList = ({ houseId, loading, data, onChanged }: TablesListProps) => {
  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无桌台数据</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data}
      renderItem={({ item }) => <TablesItem houseId={houseId} data={item} onChanged={onChanged} />}
      keyExtractor={(item) => `${item.table_id}-${item.mapped_num}`}
      ItemSeparatorComponent={() => <View className="h-3" />}
      contentContainerStyle={{ paddingBottom: 12 }}
    />
  );
};
