import React, { useEffect, useMemo, useState } from 'react';
import { Platform, Pressable, RefreshControl, ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { shopsCtrlAccountsListAll } from '@/services/shops/ctrlAccounts';
import { gameSessionStart, gameSessionStop } from '@/services/game/session';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardRow, InfoCardFooter, InfoCardContent } from '@/components/shared/info-card';
import { shopsHousesOptions } from '@/services/shops/houses';
import { Icon } from '@/components/ui/icon';
import { ChevronDown, Activity, Clock, CheckCircle, XCircle, AlertCircle } from 'lucide-react-native';
import { omitResponderProps } from '@/lib/utils';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { usePermission } from '@/hooks/use-permission';

export const SessionView = () => {
  const { getLoginModeLabel } = usePlazaConsts();
  const { isSuperAdmin, isStoreAdmin } = usePermission();
  const [houseGid, setHouseGid] = useState('');
  const { data: accounts, loading, run: getAccounts, refresh } = useRequest(shopsCtrlAccountsListAll, { manual: true });
  const { run: startSession, loading: startLoading } = useRequest(gameSessionStart, { manual: true });
  const { run: stopSession, loading: stopLoading } = useRequest(gameSessionStop, { manual: true });
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const [open, setOpen] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const filtered = useMemo(() => {
    const list = (houseOptions ?? []).map((v) => String(v));
    const q = houseGid.trim();
    if (!q) return list;
    return list.filter((v) => v.includes(q));
  }, [houseOptions, houseGid]);

  // Auto-load for store admins (they typically manage one store)
  useEffect(() => {
    if (isStoreAdmin && houseOptions && houseOptions.length > 0 && !houseGid) {
      const firstHouse = String(houseOptions[0]);
      setHouseGid(firstHouse);
      getAccounts({ house_gid: Number(firstHouse) });
    }
  }, [isStoreAdmin, houseOptions, houseGid, getAccounts]);

  const handleSearch = () => {
    if (!houseGid) return;
    setOpen(false);
    getAccounts({ house_gid: Number(houseGid) });
  };

  const handleRefresh = async () => {
    if (!houseGid) return;
    setRefreshing(true);
    try {
      await refresh();
    } finally {
      setRefreshing(false);
    }
  };

  const handleStart = async (ctrlId: number) => {
    if (!houseGid) return;
    await startSession({ id: ctrlId, house_gid: Number(houseGid) });
    // Refresh to get updated status
    await handleRefresh();
  };

  const handleStop = async (ctrlId: number) => {
    if (!houseGid) return;
    await stopSession({ id: ctrlId, house_gid: Number(houseGid) });
    // Refresh to get updated status
    await handleRefresh();
  };

  // Helper to render session status badge
  const renderStatusBadge = (status: number) => {
    if (status === 1) {
      return (
        <View className="flex-row items-center gap-1">
          <Icon as={CheckCircle} size={14} className="text-green-600" />
          <Text className="text-sm text-green-600">启用</Text>
        </View>
      );
    }
    return (
      <View className="flex-row items-center gap-1">
        <Icon as={XCircle} size={14} className="text-red-600" />
        <Text className="text-sm text-red-600">禁用</Text>
      </View>
    );
  };

  return (
    <ScrollView
      className="flex-1 bg-secondary"
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
      }
    >
      {/* Header Section */}
      <View className="bg-card border-b border-b-border p-4">
        <View className="mb-3">
          <Text className="text-lg font-semibold">店铺会话管理</Text>
          <Text className="text-sm text-muted-foreground mt-1">
            {isStoreAdmin ? '查看和管理您的店铺会话状态' : '查询和管理店铺中控账号会话'}
          </Text>
        </View>

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

      {/* Store Info Card (for Store Admins) */}
      {isStoreAdmin && houseGid && (
        <InfoCard>
          <InfoCardHeader>
            <View className="flex-row items-center gap-2">
              <Icon as={Activity} size={18} className="text-primary" />
              <InfoCardTitle>店铺信息</InfoCardTitle>
            </View>
          </InfoCardHeader>
          <InfoCardContent>
            <InfoCardRow label="店铺号" value={houseGid} />
            <InfoCardRow
              label="会话数量"
              value={accounts ? `${accounts.length} 个中控账号` : '加载中...'}
            />
            <InfoCardRow
              label="活跃会话"
              value={accounts ? `${accounts.filter(a => a.status === 1).length} 个` : '-'}
            />
          </InfoCardContent>
        </InfoCard>
      )}

      {/* Account Cards */}
      {accounts && accounts.length > 0 ? (
        accounts.map((account) => (
          <InfoCard key={account.id}>
            <InfoCardHeader>
              <View className="flex-row items-center justify-between">
                <InfoCardTitle>中控账号 #{account.id}</InfoCardTitle>
                {renderStatusBadge(account.status!)}
              </View>
            </InfoCardHeader>
            <InfoCardContent>
              <InfoCardRow label="账号标识" value={account.identifier} />
              <InfoCardRow label="登录方式" value={getLoginModeLabel(account.login_mode as any)} />

              {/* Session Status Info */}
              <View className="mt-3 pt-3 border-t border-border">
                <Text className="text-sm font-medium mb-2">会话状态</Text>
                <View className="flex-row items-center gap-2 mb-2">
                  <Icon as={Activity} size={14} className="text-muted-foreground" />
                  <Text className="text-sm text-muted-foreground">
                    会话状态: <Text className="text-foreground">待查询</Text>
                  </Text>
                </View>
                <View className="flex-row items-center gap-2">
                  <Icon as={Clock} size={14} className="text-muted-foreground" />
                  <Text className="text-sm text-muted-foreground">
                    最后同步: <Text className="text-foreground">-</Text>
                  </Text>
                </View>
              </View>

              {/* Sync Status Info */}
              {account.status === 1 && (
                <View className="mt-3 pt-3 border-t border-border">
                  <Text className="text-sm font-medium mb-2">同步信息</Text>
                  <View className="space-y-1">
                    <View className="flex-row items-center gap-2">
                      <Icon as={AlertCircle} size={12} className="text-blue-600" />
                      <Text className="text-xs text-muted-foreground">
                        对战记录: 每 5 秒同步
                      </Text>
                    </View>
                    <View className="flex-row items-center gap-2">
                      <Icon as={AlertCircle} size={12} className="text-blue-600" />
                      <Text className="text-xs text-muted-foreground">
                        成员列表: 每 30 秒同步
                      </Text>
                    </View>
                    <View className="flex-row items-center gap-2">
                      <Icon as={AlertCircle} size={12} className="text-blue-600" />
                      <Text className="text-xs text-muted-foreground">
                        钱包交易: 每 10 秒同步
                      </Text>
                    </View>
                  </View>
                </View>
              )}
            </InfoCardContent>

            <InfoCardFooter>
              <PermissionGate anyOf={["game:ctrl:create"]}>
                <Button
                  disabled={startLoading || account.status !== 1}
                  onPress={() => handleStart(account.id!)}
                  variant="default"
                >
                  <Text>启动会话</Text>
                </Button>
                <Button
                  disabled={stopLoading}
                  onPress={() => handleStop(account.id!)}
                  variant="outline"
                >
                  <Text>停止会话</Text>
                </Button>
              </PermissionGate>
            </InfoCardFooter>
          </InfoCard>
        ))
      ) : loading ? (
        <View className="min-h-32 flex-row items-center justify-center p-4">
          <Text className="text-muted-foreground">加载中...</Text>
        </View>
      ) : houseGid ? (
        <View className="min-h-32 flex-row items-center justify-center p-4">
          <View className="items-center">
            <Icon as={AlertCircle} size={48} className="text-muted-foreground mb-2" />
            <Text className="text-muted-foreground">该店铺暂无中控账号</Text>
            {isSuperAdmin && (
              <Text className="text-xs text-muted-foreground mt-1">
                请在个人中心添加并绑定中控账号
              </Text>
            )}
          </View>
        </View>
      ) : (
        <View className="min-h-32 flex-row items-center justify-center p-4">
          <View className="items-center">
            <Icon as={Activity} size={48} className="text-muted-foreground mb-2" />
            <Text className="text-muted-foreground">请输入店铺号查询会话信息</Text>
          </View>
        </View>
      )}
    </ScrollView>
  );
};

