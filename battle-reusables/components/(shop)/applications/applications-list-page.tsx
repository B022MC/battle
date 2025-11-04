import React from 'react';
import { View, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';
import { useAuthStore } from '@/hooks/use-auth-store';
import { useRequest } from '@/hooks/use-request';
import { shopsMyGroups } from '@/services/shops/groups';
import { shopsApplicationsList } from '@/services/shops/applications';
import { ApplicationsList } from './applications-list';

export const ApplicationsListPage = () => {
  const { isSuperAdmin } = usePermission();
  const userId = useAuthStore((s) => s.user?.id) ?? 0;

  const [houseInput, setHouseInput] = React.useState('');
  const [items, setItems] = React.useState<API.ShopsApplicationsItem[]>([]);
  const [loading, setLoading] = React.useState(false);

  const { run: fetchMyGroups } = useRequest(shopsMyGroups, { manual: true });
  const { run: fetchList } = useRequest(shopsApplicationsList, { manual: true });

  const loadForSuperAdmin = async (houseGid: number) => {
    setLoading(true);
    try {
      const res = await fetchList({ house_gid: houseGid, type: 1 });
      setItems(res?.items ?? []);
    } finally {
      setLoading(false);
    }
  };

  const loadForAdmin = async () => {
    setLoading(true);
    try {
      const groups = await fetchMyGroups();
      const houses = Array.from(new Set((groups ?? []).map((g) => g.house_gid).filter(Boolean)));
      const all: API.ShopsApplicationsItem[] = [];
      for (const hg of houses) {
        const res = await fetchList({ house_gid: hg, type: 2, admin_user_id: userId });
        if (Array.isArray(res?.items)) all.push(...res.items);
      }
      setItems(all);
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    if (isSuperAdmin) return; // 超管走手动或默认逻辑
    loadForAdmin();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isSuperAdmin]);

  return (
    <View className="flex-1">
      <View className="p-4 border-b border-b-border gap-2">
        <Text className="text-lg font-medium">申请列表</Text>
        {isSuperAdmin ? (
          <View className="flex-row items-center gap-2">
            <Input className="w-40" keyboardType="numeric" placeholder="店铺号" value={houseInput} onChangeText={setHouseInput} />
            <Button disabled={!houseInput} onPress={() => loadForSuperAdmin(Number(String(houseInput).trim()))}>
              <Text>查询</Text>
            </Button>
          </View>
        ) : (
          <Text className="text-muted-foreground">管理员：自动加载你所在圈子的入圈申请</Text>
        )}
      </View>
      <ScrollView className="flex-1 bg-secondary">
        <ApplicationsList loading={loading} data={items} />
      </ScrollView>
    </View>
  );
};


