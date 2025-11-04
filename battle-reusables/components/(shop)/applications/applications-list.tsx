import React from 'react';
import { View, FlatList } from 'react-native';
import { Text } from '@/components/ui/text';
import { ApplicationsItem } from './applications-item';
import { Loading } from '@/components/shared/loading';

type ApplicationsListProps = {
  loading?: boolean;
  data?: API.ShopsApplicationsItem[];
};

export const ApplicationsList = ({ loading, data }: ApplicationsListProps) => {
  if (loading) return <Loading text="加载中..." />;

  if (!data || data.length === 0) {
    return (
      <View className="min-h-16 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无申请数据</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data}
      renderItem={({ item }) => <ApplicationsItem data={item} />}
      keyExtractor={(item) => `${item.id}`}
    />
  );
};

