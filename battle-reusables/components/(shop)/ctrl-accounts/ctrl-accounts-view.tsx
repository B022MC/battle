import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { CtrlAccountsSearch } from './ctrl-accounts-search';
import { CtrlAccountsList } from './ctrl-accounts-list';
import { useRequest } from '@/hooks/use-request';
import { shopsCtrlAccountsList } from '@/services/shops/ctrlAccounts';

export const CtrlAccountsView = () => {
  const [searchParams, setSearchParams] = useState<API.ShopsCtrlAccountsListParams>();
  const { data, loading, run } = useRequest(shopsCtrlAccountsList, { manual: true });

  const handleSubmit = (params: API.ShopsCtrlAccountsListParams) => {
    setSearchParams(params);
    run(params);
  };

  return (
    <View className="flex-1">
      <CtrlAccountsSearch onSubmit={handleSubmit} submitButtonProps={{ loading }} />
      <ScrollView className="flex-1 bg-secondary">
        <CtrlAccountsList houseId={searchParams?.house_gid} data={data} />
      </ScrollView>
    </View>
  );
};

