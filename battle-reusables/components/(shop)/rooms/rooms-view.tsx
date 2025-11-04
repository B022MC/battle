import React, { useMemo, useState } from 'react';
import { Platform, Pressable, ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { shopsTablesList, shopsTablesPull, shopsTablesCheck, shopsTablesDetail, shopsTablesDismiss } from '@/services/shops/tables';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardRow, InfoCardFooter, InfoCardContent } from '@/components/shared/info-card';
import { shopsHousesOptions } from '@/services/shops/houses';
import { Icon } from '@/components/ui/icon';
import { ChevronDown } from 'lucide-react-native';

export const RoomsView = () => {
  const [houseGid, setHouseGid] = useState('');
  const [mappedNum, setMappedNum] = useState('');

  const { data: list, loading, run: getList } = useRequest(shopsTablesList, { manual: true });
  const { run: pull, loading: pulling } = useRequest(shopsTablesPull, { manual: true });
  const { run: check, loading: checking } = useRequest(shopsTablesCheck, { manual: true });
  const { data: detail, loading: detailing, run: getDetail } = useRequest(shopsTablesDetail, { manual: true });
  const { run: dismiss, loading: dismissing } = useRequest(shopsTablesDismiss, { manual: true });
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const [open, setOpen] = useState(false);
  const filtered = useMemo(() => {
    const list = (houseOptions ?? []).map((v) => String(v));
    const q = houseGid.trim();
    if (!q) return list;
    return list.filter((v) => v.includes(q));
  }, [houseOptions, houseGid]);

  const items = useMemo(() => list?.items ?? [], [list]);

  const handleQuery = async () => {
    if (!houseGid) return;
    await getList({ house_gid: Number(houseGid) });
  };

  const handlePull = async () => {
    if (!houseGid) return;
    await pull({ house_gid: Number(houseGid) });
    await getList({ house_gid: Number(houseGid) });
  };

  const handleCheck = async () => {
    if (!houseGid || !mappedNum) return;
    await check({ house_gid: Number(houseGid), mapped_num: Number(mappedNum) });
    await getDetail({ house_gid: Number(houseGid), mapped_num: Number(mappedNum) });
  };

  const handleDismiss = async (mn: number) => {
    if (!houseGid) return;
    await dismiss({ house_gid: Number(houseGid), mapped_num: mn });
    await getList({ house_gid: Number(houseGid) });
  };

  return (
    <ScrollView className="flex-1 bg-secondary p-4">
      <View className="mb-4">
        <View className="flex flex-row gap-2">
          <View className="relative flex-1">
            <Input
              keyboardType="numeric"
              className="pr-8"
              placeholder="店铺号（可输入或下拉选择）"
              value={houseGid}
              onChangeText={(t) => { setHouseGid(t); setOpen(true); }}
              onSubmitEditing={handleQuery}
            />
            <Pressable
              accessibilityRole="button"
              onPress={() => setOpen((v) => !v)}
              className="absolute right-2 top-1/2 -translate-y-1/2"
            >
              <Icon as={ChevronDown} className="text-muted-foreground" size={16} />
            </Pressable>
            {open && (
              <View
                className={
                  Platform.select({
                    web: 'bg-popover border-border absolute left-0 top-full z-50 mt-1 max-h-56 w-full overflow-y-auto rounded-md border shadow-sm shadow-black/5',
                    default: 'bg-popover border-border absolute left-0 top-full z-50 mt-1 max-h-56 w-full rounded-md border',
                  }) as string
                }
              >
                {(filtered.length > 0 ? filtered : ['无匹配结果']).map((gid) => (
                  <Pressable
                    key={gid}
                    onPress={() => { if (gid !== '无匹配结果') { setHouseGid(gid); setOpen(false); } }}
                    className="px-3 py-2"
                    accessibilityRole="button"
                  >
                    <Text className="text-sm">{gid === '无匹配结果' ? gid : `店铺 ${gid}`}</Text>
                  </Pressable>
                ))}
              </View>
            )}
          </View>
          <Button disabled={!houseGid || loading} onPress={handleQuery}>
            <Text>查询</Text>
          </Button>
          <Button disabled={!houseGid || pulling} onPress={handlePull} variant="outline">
            <Text>刷新</Text>
          </Button>
        </View>
      </View>

      <InfoCard className="mb-4">
        <InfoCardHeader>
          <InfoCardTitle>查桌/详情</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          <View className="flex-row gap-2">
            <Input
              keyboardType="numeric"
              className="flex-1"
              placeholder="桌映射号 MappedNum"
              value={mappedNum}
              onChangeText={setMappedNum}
            />
            <Button disabled={!houseGid || !mappedNum || checking || detailing} onPress={handleCheck}>
              <Text>查桌</Text>
            </Button>
          </View>
          {detail?.table && (
            <View className="mt-3 gap-1">
              <Text className="text-muted-foreground">桌详情</Text>
              <InfoCardRow label="映射号" value={String(detail.table.mapped_num)} />
              <InfoCardRow label="KindID" value={String(detail.table.kind_id)} />
              <InfoCardRow label="底分" value={String(detail.table.base_score)} />
              <InfoCardRow label="GroupID" value={String(detail.table.group_id)} />
            </View>
          )}
        </InfoCardContent>
      </InfoCard>

      <InfoCard>
        <InfoCardHeader>
          <InfoCardTitle>房间列表</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          {items.length === 0 ? (
            <View className="min-h-16 items-center justify-center">
              <Text className="text-muted-foreground">无房间</Text>
            </View>
          ) : (
            <View className="gap-2">
              {items.map((t, idx) => (
                <View key={String(t.mapped_num ?? idx)} className="flex-row items-center justify-between rounded-md border border-border p-3">
                  <View>
                    <Text>映射号: {t.mapped_num ?? '-'}</Text>
                    <Text className="text-muted-foreground text-xs">KindID: {t.kind_id} 底分: {t.base_score}</Text>
                  </View>
                  <View className="flex-row gap-2">
                    <Button variant="outline" onPress={() => { if (t.mapped_num != null) { setMappedNum(String(t.mapped_num)); handleCheck(); } }}>
                      <Text>详情</Text>
                    </Button>
                    <Button disabled={dismissing || t.mapped_num == null} onPress={() => handleDismiss(Number(t.mapped_num))}>
                      <Text>解散</Text>
                    </Button>
                  </View>
                </View>
              ))}
            </View>
          )}
        </InfoCardContent>
        <InfoCardFooter>
          <Text className="text-muted-foreground text-xs">共 {items.length} 个</Text>
        </InfoCardFooter>
      </InfoCard>
    </ScrollView>
  );
};


