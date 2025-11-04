import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { MembersSearch } from './members-search';
import { MembersList } from './members-list';
import { useRequest } from '@/hooks/use-request';
import { shopsMembersList } from '@/services/shops/members';

export const MembersView = () => {
  const [searchParams, setSearchParams] = useState<API.ShopsMembersListParams>();
  const { data, loading, run } = useRequest(shopsMembersList, { manual: true });

  const handleSubmit = (params: API.ShopsMembersListParams) => {
    setSearchParams(params);
    run(params);
  };

  return (
    <View className="flex-1">
      <MembersSearch onSubmit={handleSubmit} submitButtonProps={{ loading }} />
      <ScrollView className="flex-1 bg-secondary">
        <MembersList houseId={searchParams?.house_gid} data={data?.items} />
      </ScrollView>
    </View>
  );
};
