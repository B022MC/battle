import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { AdminsSearch } from './admins-search';
import { AdminsList } from './admins-list';
import { useRequest } from '@/hooks/use-request';
import { shopsAdminsList } from '@/services/shops/admins';
import { AdminsAssign } from './admins-assign';
import { PermissionGate } from '@/components/auth/PermissionGate';

export const AdminsView = () => {
  const [searchParams, setSearchParams] = useState<API.ShopsAdminsListParams>();
  const { data, loading, run } = useRequest(shopsAdminsList, { manual: true });

  const handleSubmit = (params: API.ShopsAdminsListParams) => {
    setSearchParams(params);
    run(params);
  };

  return (
    <View className="flex-1">
      <AdminsSearch onSubmit={handleSubmit} submitButtonProps={{ loading }} />
      <PermissionGate anyOf={["shop:admin:assign"]}>
        <AdminsAssign houseId={searchParams?.house_gid} />
      </PermissionGate>
      <ScrollView className="flex-1 bg-secondary">
        <AdminsList houseId={searchParams?.house_gid} data={data} />
      </ScrollView>
    </View>
  );
};

