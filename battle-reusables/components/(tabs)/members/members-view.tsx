import React, { useMemo, useState } from 'react';
import { ScrollView, View } from 'react-native';
import { MembersSearch } from './members-search';
import { MembersList } from './members-list';
import { useRequest } from '@/hooks/use-request';
import { shopsMembersList } from '@/services/shops/members';
import { usePermission } from '@/hooks/use-permission';
import { shopsGroupsList, shopsMyGroups } from '@/services/shops/groups';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { shopsMembersListPlatform, shopsMembersRemovePlatform } from '@/services/shops/members';
import { alert } from '@/utils/alert';

export const MembersView = () => {
  const [searchParams, setSearchParams] = useState<API.ShopsMembersListParams>();
  const { data, loading, run } = useRequest(shopsMembersList, { manual: true });
  const { isSuperAdmin } = usePermission();

  const [groupId, setGroupId] = useState<string>('');
  const { run: runGroups, data: groups } = useRequest(shopsGroupsList, { manual: true });
  const { run: runMyGroups, data: myGroups } = useRequest(shopsMyGroups, { manual: true });

  // 平台端数据（基于平台入圈记录）
  const [mode, setMode] = useState<'game' | 'platform'>('game');
  const { data: platData, loading: platLoading, run: runPlat } = useRequest(shopsMembersListPlatform, { manual: true });
  const { run: runRemovePlat, loading: removing } = useRequest(shopsMembersRemovePlatform, { manual: true });

  const handleSubmit = async (params: API.ShopsMembersListParams) => {
    setSearchParams(params);
    if (mode === 'game') {
      await run(params);
      alert.show({ title: '已加载', duration: 800 });
    } else {
      // 平台端：管理员需带上 group_id（=我的 user_id）；超管可不带
      let gidForHouse: number | undefined = undefined;
      if (!isSuperAdmin) {
        const groupsData = myGroups ?? await runMyGroups();
        const found = (groupsData ?? []).find((g: any) => Number(g.house_gid ?? g.HouseGID) === params.house_gid);
        if (found) gidForHouse = Number(found.group_id ?? found.GroupID);
      }
      await runPlat({ house_gid: params.house_gid, group_id: gidForHouse });
      alert.show({ title: '已加载', duration: 800 });
    }
    if (isSuperAdmin) runGroups({ house_gid: params.house_gid }); else runMyGroups();
  };

  const filtered = useMemo(() => {
    const list = data?.items ?? [];
    const gid = Number(groupId || '');
    if (isNaN(gid) || gid <= 0) return list;
    return list.filter(it => (it.group_id ?? 0) === gid);
  }, [data, groupId]);

  return (
    <View className="flex-1">
      <View className="flex-row items-center gap-2 px-4 py-2 border-b border-b-border">
        <Button variant={mode==='game'?'secondary':'outline'} onPress={() => setMode('game')}><Text>游戏端</Text></Button>
        <Button variant={mode==='platform'?'secondary':'outline'} onPress={() => setMode('platform')}><Text>平台端</Text></Button>
      </View>
      <MembersSearch onSubmit={handleSubmit} submitButtonProps={{ loading }} />
      {isSuperAdmin ? (
        <View className="flex-row items-center gap-2 px-4 py-2 border-b border-b-border">
          <Text>圈ID</Text>
          <Input className="w-28" keyboardType="numeric" value={groupId} onChangeText={setGroupId} placeholder="可选" />
          <Button variant="outline" onPress={() => setGroupId('')}><Text>清除</Text></Button>
          {Array.isArray(groups) && groups.length > 0 && (
            <Text className="text-muted-foreground">可见圈: {groups.join(', ')}</Text>
          )}
        </View>
      ) : (
        <View className="px-4 py-2 border-b border-b-border">
          <Text className="text-muted-foreground">我的圈: {Array.isArray(myGroups) && myGroups.length > 0 ? myGroups.map(g => `${g.house_gid}/${g.group_id}`).join(', ') : '-'}</Text>
        </View>
      )}
      <ScrollView className="flex-1 bg-secondary">
        {mode === 'game' ? (
          <MembersList houseId={searchParams?.house_gid} data={filtered} />
        ) : (
          <View className="p-4 gap-2">
            {((platData as any)?.items ?? []).map((it: any, idx: number) => {
              const gid = (() => {
                const found = (myGroups ?? []).find((g: any) => Number(g.house_gid ?? g.HouseGID) === (searchParams?.house_gid ?? 0));
                return Number(found?.group_id ?? found?.GroupID ?? 0) || undefined;
              })();
              return (
                <View key={idx} className="rounded-md border border-border p-3 gap-2">
                  <Text>成员 {it.nick_name || it.user_id || '-'}</Text>
                  <Button variant="outline" disabled={removing} onPress={async () => {
                    await runRemovePlat({ house_gid: searchParams?.house_gid as number, group_id: gid, member_user_id: Number(it.user_id || 0) });
                    alert.show({ title: '已移除', duration: 800 });
                    await runPlat({ house_gid: searchParams?.house_gid as number, group_id: gid });
                  }}>
                    <Text>踢出</Text>
                  </Button>
                </View>
              );
            })}
            {(!platData || !(platData as any).items || (platData as any).items.length===0) && <Text className="text-muted-foreground">暂无平台成员数据</Text>}
          </View>
        )}
      </ScrollView>
    </View>
  );
};
