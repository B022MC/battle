import React, { useState, useMemo, useEffect } from 'react';
import { View, ScrollView, ActivityIndicator, Platform, Pressable } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Icon } from '@/components/ui/icon';
import { ChevronDown, X } from 'lucide-react-native';
import { showToast } from '@/utils/toast';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useHouseSelector } from '@/hooks/use-house-selector';
import { useRequest } from '@/hooks/use-request';
import { listGroupsByHouse } from '@/services/shops/groups';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

export function RoomCreditsView() {
  const insets = useSafeAreaInsets();
  const {
    houseGid,
    setHouseGid,
    isSuperAdmin,
    isStoreAdmin,
    houseOptions,
    loading: loadingHouse,
    isReady,
    canSelectHouse,
  } = useHouseSelector();
  
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  
  // 额度设置参数
  const [groupName, setGroupName] = useState<string | undefined>(undefined);
  const [gameKind, setGameKind] = useState<string | undefined>(undefined);
  const [baseScore, setBaseScore] = useState('');
  const [creditLimit, setCreditLimit] = useState('');

  // 获取游戏类型常量
  const { maps, loading: loadingConsts } = usePlazaConsts();
  const gameKinds = useMemo(() => {
    const kinds: { label: string; value: string }[] = [];
    maps.game_kinds.forEach((label, value) => {
      kinds.push({ label, value: String(value) });
    });
    return kinds;
  }, [maps]);

  // 获取圈子列表
  const { data: groupsData, run: loadGroups } = useRequest(listGroupsByHouse, { manual: true });

  // 当店铺ID变化时加载圈子
  useEffect(() => {
    if (houseGid) {
      loadGroups({ house_gid: Number(houseGid) });
      // 重置选择
      setGroupName(undefined);
    }
  }, [houseGid, loadGroups]);

  const groups = useMemo(() => groupsData || [], [groupsData]);

  // 查询房间额度列表
  const handleQuery = async () => {
    if (!houseGid) {
      showToast('请输入店铺号', 'error');
      return;
    }

    try {
      setLoading(true);
      // TODO: 调用 API 获取房间额度列表
      // const res = await roomCreditList({ house_gid: Number(houseGid) });
      showToast('功能开发中...', 'info');
    } catch (error) {
      showToast('查询失败', 'error');
      console.error('查询房间额度失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 设置房间额度
  const handleSetCredit = async () => {
    if (!houseGid) {
      showToast('请输入店铺号', 'error');
      return;
    }

    if (!gameKind) {
      showToast('请选择游戏类型', 'error');
      return;
    }

    if (!creditLimit) {
      showToast('请输入额度限制', 'error');
      return;
    }

    try {
      setLoading(true);
      // TODO: 调用 API 设置房间额度
      // const res = await roomCreditSet({
      //   house_gid: Number(houseGid),
      //   group_name: groupName || '',
      //   game_kind: gameKind ? Number(gameKind) : 0,
      //   base_score: baseScore ? Number(baseScore) : 0,
      //   credit_limit: Number(creditLimit),
      // });
      showToast('功能开发中...', 'info');
    } catch (error) {
      showToast('设置失败', 'error');
      console.error('设置房间额度失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 下拉框选项
  const filtered = useMemo(() => {
    const list = (houseOptions ?? []).map((v) => String(v));
    const q = houseGid.trim();
    if (!q) return list;
    return list.filter((v) => v.includes(q));
  }, [houseOptions, houseGid]);

  return (
    <View className="flex-1 bg-gray-50">
      {/* 头部查询区 */}
      <View className="bg-white p-4 border-b border-gray-200">
        <Text className="text-lg font-semibold mb-3">房间额度管理</Text>
        
        {/* 店铺管理员：显示当前店铺 */}
        {isStoreAdmin && (
          <View className="mb-3 p-3 bg-blue-50 rounded border border-blue-200">
            <Text className="text-sm text-blue-700">
              当前店铺：{houseGid || '加载中...'}
            </Text>
          </View>
        )}
        
        {/* 超级管理员：下拉选择店铺 */}
        {isSuperAdmin && (
          <View className="flex-row gap-2 mb-3">
            <View className="relative flex-1">
              <Input
                keyboardType="numeric"
                className="pr-8"
                placeholder="店铺号（可输入或下拉选择）"
                value={houseGid}
                onChangeText={(t) => { setHouseGid(t); setOpen(true); }}
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
            <Button
              onPress={handleQuery}
              disabled={!houseGid || loading}
            >
              {loading ? (
                <ActivityIndicator size="small" color="white" />
              ) : (
                <Text className="text-white">查询</Text>
              )}
            </Button>
          </View>
        )}
      </View>

      <ScrollView className="flex-1 p-4">
        {/* 设置房间额度 */}
        <PermissionGate anyOf={['game:room_credit:set']}>
          <Card className="mb-4 p-4">
            <Text className="text-base font-semibold mb-3">设置房间额度</Text>
            
            <View className="mb-3 z-20">
              <Text className="text-sm text-gray-600 mb-1">
                圈子名称（可选）
              </Text>
              <View className="flex-row gap-2">
                <View className="flex-1">
                  <Select
                    value={groupName ? { value: groupName, label: groupName } : undefined}
                    onValueChange={(option) => setGroupName(option?.value)}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="选择圈子（留空表示全局）" />
                    </SelectTrigger>
                    <SelectContent portalHost="shop-layout-portal">
                      <SelectGroup>
                        <SelectLabel>所有圈子</SelectLabel>
                        {groups.map((g) => (
                          <SelectItem key={g.id} label={g.group_name} value={g.group_name} />
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </View>
                {groupName && (
                  <Button
                    variant="outline"
                    size="icon"
                    className="w-10"
                    onPress={() => setGroupName(undefined)}
                  >
                    <Icon as={X} size={16} />
                  </Button>
                )}
              </View>
            </View>

            <View className="mb-3 z-10">
              <Text className="text-sm text-gray-600 mb-1">
                游戏类型 *
              </Text>
              <Select
                value={gameKind ? { value: gameKind, label: maps.game_kinds.get(Number(gameKind)) || gameKind } : undefined}
                onValueChange={(option) => setGameKind(option?.value)}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="请选择游戏类型" />
                </SelectTrigger>
                <SelectContent portalHost="shop-layout-portal">
                  <SelectGroup>
                    <SelectLabel>游戏列表</SelectLabel>
                    {gameKinds.map((k) => (
                      <SelectItem key={k.value} label={k.label} value={k.value} />
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </View>

            <View className="mb-3">
              <Text className="text-sm text-gray-600 mb-1">
                底分（可选，0表示默认）
              </Text>
              <Input
                placeholder="0"
                keyboardType="numeric"
                value={baseScore}
                onChangeText={setBaseScore}
              />
            </View>

            <View className="mb-3">
              <Text className="text-sm text-gray-600 mb-1">
                额度限制（单位：分）*
              </Text>
              <Input
                placeholder="请输入额度限制，如 10000 表示 100 元"
                keyboardType="numeric"
                value={creditLimit}
                onChangeText={setCreditLimit}
              />
            </View>

            <Button
              onPress={handleSetCredit}
              disabled={!houseGid || !creditLimit || !gameKind || loading}
            >
              <Text className="text-white">设置额度</Text>
            </Button>

            <View className="mt-3 bg-blue-50 p-3 rounded">
              <Text className="text-xs text-blue-700">
                💡 提示：
                {'\n'}• 留空所有可选项表示设置全局默认额度
                {'\n'}• 指定圈子但不指定游戏类型/底分表示该圈子的默认额度
                {'\n'}• 完整指定表示特定房间的额度要求
                {'\n'}• 查找优先级：精确匹配 &gt; 圈子默认 &gt; 全局默认
              </Text>
            </View>
          </Card>
        </PermissionGate>

        {/* 额度列表区域 - TODO */}
        <PermissionGate anyOf={['game:room_credit:view']}>
          <Card className="mb-4 p-4">
            <Text className="text-base font-semibold mb-3">房间额度列表</Text>
            
            {!houseGid ? (
              <View className="items-center justify-center py-8">
                <Text className="text-gray-400">请输入店铺号查询额度列表</Text>
              </View>
            ) : (
              <View className="items-center justify-center py-8">
                <Text className="text-gray-400">功能开发中...</Text>
              </View>
            )}
          </Card>
        </PermissionGate>

        {/* API说明 */}
        <Card className="mb-4 p-4 bg-gray-50">
          <Text className="text-base font-semibold mb-2">API端点</Text>
          <View className="bg-white p-3 rounded border border-gray-200">
            <Text className="text-xs font-mono text-gray-700 mb-1">
              POST /room-credit/set - 设置房间额度
            </Text>
            <Text className="text-xs font-mono text-gray-700 mb-1">
              POST /room-credit/list - 查询额度列表
            </Text>
            <Text className="text-xs font-mono text-gray-700 mb-1">
              POST /room-credit/delete - 删除额度配置
            </Text>
            <Text className="text-xs font-mono text-gray-700">
              POST /room-credit/check - 检查玩家额度
            </Text>
          </View>
          <Text className="text-xs text-gray-500 mt-2">
            提示：API接口已就绪，前端功能待完善
          </Text>
        </Card>
      </ScrollView>
    </View>
  );
}
