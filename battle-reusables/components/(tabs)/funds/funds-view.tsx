import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { FundsSearch } from './funds-search';
import { FundsList } from './funds-list';
import { useRequest } from '@/hooks/use-request';
import { membersWalletList } from '@/services/members/wallet';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { Text } from '@/components/ui/text';

export const FundsView = () => {
  const [searchParams, setSearchParams] = useState<API.MembersWalletListParams>();
  const { data, loading, run } = useRequest(membersWalletList, { manual: true });

  const handleSubmit = (params: API.MembersWalletListParams) => {
    setSearchParams(params);
    run(params);
  };

  return (
    <View className="flex-1">
      <PermissionGate anyOf={["fund:wallet:view"]} fallback={null}>
        <FundsSearch onSubmit={handleSubmit} submitButtonProps={{ loading }} />
      </PermissionGate>
      <ScrollView className="flex-1 bg-secondary">
        <PermissionGate anyOf={["fund:wallet:view"]} fallback={<View className="p-4"><Text className="text-muted-foreground">无查看分数权限</Text></View>}>
          <FundsList houseId={searchParams?.house_gid} data={data?.list} />
        </PermissionGate>
      </ScrollView>
    </View>
  );
};
