import React from 'react';
import { View, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';
import { useAuthStore } from '@/hooks/use-auth-store';
import { useRequest } from '@/hooks/use-request';
import { shopsMyGroups } from '@/services/shops/groups';
import { shopsApplicationsList, shopsApplicationsHistory } from '@/services/shops/applications';
import { ApplicationsList } from './applications-list';
import { alert } from '@/utils/alert';

export const ApplicationsListPage = () => {
  const { isSuperAdmin } = usePermission();
  const userId = useAuthStore((s) => s.user?.id) ?? 0;

  const [houseInput, setHouseInput] = React.useState('');
  const [items, setItems] = React.useState<API.ShopsApplicationsItem[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [historyItems, setHistoryItems] = React.useState<API.ShopsApplicationsItem[]>([]);
  const [historyLoading, setHistoryLoading] = React.useState(false);
  const [lastSuperHouse, setLastSuperHouse] = React.useState<number | undefined>(undefined);

  const { run: fetchMyGroups } = useRequest(shopsMyGroups, { manual: true });
  const { run: fetchList } = useRequest(shopsApplicationsList, { manual: true });
  const { run: fetchHistory } = useRequest(shopsApplicationsHistory, { manual: true });

  const loadForSuperAdmin = async (houseGid: number) => {
    setLoading(true);
    try {
      const res = await fetchList({ house_gid: houseGid, type: 1 });
      setItems(res?.items ?? []);
      alert.show({ title: '已加载', duration: 800 });
      setLastSuperHouse(houseGid);
    } finally {
      setLoading(false);
    }
  };

  const loadForAdmin = async () => {
    if (!userId || userId <= 0) return;
    setLoading(true);
    try {
      const groups = await fetchMyGroups();
      const houses = Array.from(
        new Set(
          (groups ?? [])
            .map((g: any) => Number((g?.house_gid ?? g?.HouseGID) || 0))
            .filter((v) => !!v)
        )
      );
      const all: API.ShopsApplicationsItem[] = [];
      for (const hg of houses) {
        const res = await fetchList({ house_gid: hg, type: 2, admin_user_id: userId });
        if (Array.isArray(res?.items)) all.push(...res.items);
      }
      setItems(all);
      alert.show({ title: '已加载', duration: 800 });
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    if (isSuperAdmin) return; // 超管走手动或默认逻辑
    if (userId && userId > 0) loadForAdmin();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isSuperAdmin, userId]);

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
        <ApplicationsList
          loading={loading}
          data={items}
          onChanged={() => {
            if (isSuperAdmin) {
              if (lastSuperHouse) loadForSuperAdmin(lastSuperHouse);
            } else {
              loadForAdmin();
            }
          }}
        />
        <View className="p-4 gap-2">
          <Text className="text-lg font-medium">申请记录</Text>
          <View className="flex-row items-center gap-2">
            <Input className="w-40" keyboardType="numeric" placeholder="店铺号" value={houseInput} onChangeText={setHouseInput} />
            <Button disabled={!houseInput || historyLoading} onPress={async () => {
              setHistoryLoading(true);
              try {
                const res = await fetchHistory({ house_gid: Number(String(houseInput).trim()) });
                setHistoryItems(res?.items ?? []);
                alert.show({ title: '已加载', duration: 800 });
              } finally {
                setHistoryLoading(false);
              }
            }}>
              <Text>查询记录</Text>
            </Button>
          </View>
          <ApplicationsList loading={historyLoading} data={historyItems} />
        </View>
      </ScrollView>
    </View>
  );
};


