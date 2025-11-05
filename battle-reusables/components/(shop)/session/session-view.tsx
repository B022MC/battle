import React, { useMemo, useState } from 'react';
import { Platform, Pressable, ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { shopsCtrlAccountsListAll } from '@/services/shops/ctrlAccounts';
import { gameSessionStart, gameSessionStop } from '@/services/game/session';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardRow, InfoCardFooter, InfoCardContent } from '@/components/shared/info-card';
import { shopsHousesOptions } from '@/services/shops/houses';
import { Icon } from '@/components/ui/icon';
import { ChevronDown } from 'lucide-react-native';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';

export const SessionView = () => {
  const { getLoginModeLabel } = usePlazaConsts();
  const [houseGid, setHouseGid] = useState('');
  const { data: accounts, loading, run: getAccounts } = useRequest(shopsCtrlAccountsListAll, { manual: true });
  const { run: startSession, loading: startLoading } = useRequest(gameSessionStart, { manual: true });
  const { run: stopSession, loading: stopLoading } = useRequest(gameSessionStop, { manual: true });
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const [open, setOpen] = useState(false);
  const filtered = useMemo(() => {
    const list = (houseOptions ?? []).map((v) => String(v));
    const q = houseGid.trim();
    if (!q) return list;
    return list.filter((v) => v.includes(q));
  }, [houseOptions, houseGid]);

  const handleSearch = () => {
    if (!houseGid) return;
    setOpen(false);
    getAccounts({ house_gid: Number(houseGid) });
  };

  const handleStart = async (ctrlId: number) => {
    if (!houseGid) return;
    await startSession({ id: ctrlId, house_gid: Number(houseGid) });
  };

  const handleStop = async (ctrlId: number) => {
    if (!houseGid) return;
    await stopSession({ id: ctrlId, house_gid: Number(houseGid) });
  };

  return (
    <ScrollView className="flex-1 bg-secondary">
      <View className="border-b border-b-border p-4">
        <View className="flex flex-row gap-2">
          <View className="relative flex-1">
            <Input
              keyboardType="numeric"
              className="pr-8"
              placeholder="店铺号（可输入或下拉选择）"
              value={houseGid}
              onChangeText={(t) => { setHouseGid(t); setOpen(true); }}
              onSubmitEditing={handleSearch}
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
                {filtered.length > 0 ? (
                  <ScrollView>
                    {filtered.map((gid) => (
                      <Pressable
                        key={gid}
                        onPress={() => { setHouseGid(gid); setOpen(false); }}
                        className="px-3 py-2"
                        accessibilityRole="button"
                      >
                        <Text className="text-sm">店铺 {gid}</Text>
                      </Pressable>
                    ))}
                  </ScrollView>
                ) : (
                  <View className="px-3 py-2">
                    <Text className="text-sm text-muted-foreground">无匹配结果</Text>
                  </View>
                )}
              </View>
            )}
          </View>
          <Button disabled={!houseGid || loading} onPress={handleSearch}>
            <Text>查询</Text>
          </Button>
        </View>
      </View>
      
      {accounts && accounts.length > 0 ? (
        accounts.map((account) => (
          <InfoCard key={account.id}>
            <InfoCardHeader>
              <InfoCardTitle>中控账号 #{account.id}</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <InfoCardRow label="账号" value={account.identifier} />
              <InfoCardRow label="登录方式" value={getLoginModeLabel(account.login_mode as any)} />
              <InfoCardRow label="状态" value={account.status === 1 ? '启用' : '禁用'} />
            </InfoCardContent>
            <InfoCardFooter>
              <PermissionGate anyOf={["game:ctrl:create"]}>
                <Button disabled={startLoading} onPress={() => handleStart(account.id!)}>
                  启动会话
                </Button>
                <Button disabled={stopLoading} onPress={() => handleStop(account.id!)}>
                  停止会话
                </Button>
              </PermissionGate>
            </InfoCardFooter>
          </InfoCard>
        ))
      ) : (
        <View className="min-h-16 flex-row items-center justify-center p-4">
          <Text className="text-muted-foreground">请先查询中控账号</Text>
        </View>
      )}
    </ScrollView>
  );
};

