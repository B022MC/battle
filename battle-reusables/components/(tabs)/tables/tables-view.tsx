import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { TablesSearch } from './tables-search';
import { TablesList } from './tables-list';
import { useRequest } from '@/hooks/use-request';
import { shopsTablesList } from '@/services/shops/tables';

export const TablesView = () => {
  const [searchParams, setSearchParams] = useState<API.ShopsTablesListParams | undefined>();
  const { data, loading, run } = useRequest(shopsTablesList, { manual: true });

  const handleSubmit = (params: API.ShopsTablesListParams) => {
    setSearchParams(params);
    run(params);
  };

  return (
    <View className="flex-1">
      <TablesSearch onSubmit={handleSubmit} submitButtonProps={{ loading }} />
      <ScrollView className="flex-1 bg-secondary">
        <View className="p-3">
          <TablesList
            houseId={searchParams?.house_gid}
            data={data?.items}
            loading={loading}
            onChanged={() => {
              if (searchParams) run(searchParams);
            }}
          />
        </View>
      </ScrollView>
    </View>
  );
};
