import React, { useState, useEffect } from 'react';
import { ScrollView, View } from 'react-native';
import { TablesSearch } from './tables-search';
import { TablesList } from './tables-list';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { shopsTablesList } from '@/services/shops/tables';
import { shopsAdminsMe } from '@/services/shops/admins';

export const TablesView = () => {
  const { isSuperAdmin, isStoreAdmin } = usePermission();
  const [searchParams, setSearchParams] = useState<API.ShopsTablesListParams | undefined>();
  const { data, loading, run } = useRequest(shopsTablesList, { manual: true });
  const { data: myAdminInfo } = useRequest(shopsAdminsMe, { manual: !isStoreAdmin });

  // 店铺管理员自动加载自己店铺的桌台
  useEffect(() => {
    if (isStoreAdmin && myAdminInfo?.house_gid) {
      const params = { house_gid: myAdminInfo.house_gid };
      setSearchParams(params);
      run(params);
    }
  }, [isStoreAdmin, myAdminInfo?.house_gid]);

  const handleSubmit = (params: API.ShopsTablesListParams) => {
    setSearchParams(params);
    run(params);
  };

  return (
    <View className="flex-1">
      {/* 超级管理员显示搜索功能，店铺管理员只显示刷新按钮 */}
      <TablesSearch 
        onSubmit={handleSubmit} 
        submitButtonProps={{ loading }}
        hideSearch={isStoreAdmin}
        defaultHouseGid={isStoreAdmin ? myAdminInfo?.house_gid : undefined}
      />
      <ScrollView className="flex-1 bg-secondary">
        <View className="p-3 gap-3">
          {/* 桌台列表 */}
          <View>
            <TablesList
              houseId={searchParams?.house_gid}
              data={data?.items}
              loading={loading}
              onChanged={() => {
                if (searchParams) run(searchParams);
              }}
            />
          </View>
        </View>
      </ScrollView>
    </View>
  );
};
