import React, { useState, useRef } from 'react';
import { ScrollView, View } from 'react-native';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { alert } from '@/utils/alert';
import { listGroupMemberBalances } from '@/services/battles/query';
import type { MemberBalance } from '@/services/battles/query-typing';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { shopsHousesOptions } from '@/services/shops/houses';
import { getGroupOptions } from '@/services/shops/groups';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';

export const GroupBalancesView = () => {
  const { isStoreAdmin } = usePermission();

  // 状态管理
  const [houseGid, setHouseGid] = useState<string>('');
  const [groupId, setGroupId] = useState<string>('');
  const [minYuan, setMinYuan] = useState<string>('');
  const [maxYuan, setMaxYuan] = useState<string>('');
  const [page, setPage] = useState(1);

  // API 请求
  const { data: balancesData, loading: loadingBalances, run: runListBalances } = useRequest(listGroupMemberBalances, { manual: true });
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const { data: groupOptions, run: runGetGroupOptions } = useRequest(getGroupOptions, { manual: true });
  const houseRef = useRef<TriggerRef>(null);
  const groupRef = useRef<TriggerRef>(null);

  function onHouseTouchStart() {
    isWeb && houseRef.current?.open();
  }

  function onGroupTouchStart() {
    isWeb && groupRef.current?.open();
  }

  // 当店铺改变时，加载该店铺的圈子列表
  const handleHouseChange = async (newHouseGid: string) => {
    setHouseGid(newHouseGid);
    setGroupId(''); // 重置圈子选择
    if (newHouseGid) {
      try {
        await runGetGroupOptions({ house_gid: Number(newHouseGid) });
      } catch (error) {
        console.error('加载圈子列表失败:', error);
      }
    }
  };

  // 加载余额
  const handleLoadBalances = async () => {
    if (!houseGid) {
      alert.show({ title: '请输入店铺号', variant: 'error' });
      return;
    }
    if (!groupId) {
      alert.show({ title: '请输入圈子ID', variant: 'error' });
      return;
    }

    const gid = Number(houseGid);
    const gidNum = Number(groupId);
    if (isNaN(gid) || gid <= 0 || isNaN(gidNum) || gidNum <= 0) {
      alert.show({ title: '店铺号或圈子ID格式错误', variant: 'error' });
      return;
    }

    const params: any = {
      house_gid: gid,
      group_id: gidNum,
      page,
      size: 20,
    };

    if (minYuan) {
      const min = Number(minYuan);
      if (!isNaN(min) && min >= 0) {
        params.min_yuan = min;
      }
    }

    if (maxYuan) {
      const max = Number(maxYuan);
      if (!isNaN(max) && max >= 0) {
        params.max_yuan = max;
      }
    }

    try {
      await runListBalances(params);
    } catch (error: any) {
      console.error('加载余额失败:', error);
    }
  };

  // 格式化时间
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // 计算总余额
  const getTotalBalance = () => {
    if (!balancesData?.list) return 0;
    return balancesData.list.reduce((sum, b) => sum + b.balance_yuan, 0);
  };

  // 计算平均余额
  const getAvgBalance = () => {
    if (!balancesData?.list || balancesData.list.length === 0) return 0;
    return getTotalBalance() / balancesData.list.length;
  };

  if (!isStoreAdmin) {
    return (
      <View className="flex-1 bg-background p-4 justify-center items-center">
        <Text className="text-lg text-muted-foreground">
          仅店铺管理员可访问此页面
        </Text>
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-background p-4">
      {/* 查询表单 */}
      <Card className="p-4 mb-4">
        <Text className="text-lg font-bold mb-4">查询圈子成员余额</Text>
        
        <View className="mb-3">
          <Text className="mb-2">店铺号 *</Text>
          <Select
            value={houseGid ? ({ label: `店铺 ${houseGid}`, value: houseGid } as any) : undefined}
            onValueChange={(opt) => handleHouseChange(String(opt?.value ?? ''))}
          >
            <SelectTrigger ref={houseRef} onTouchStart={onHouseTouchStart} className="min-w-[160px]">
              <SelectValue placeholder={houseGid ? `店铺 ${houseGid}` : '选择店铺号'} />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>店铺号</SelectLabel>
                {(houseOptions ?? []).map((gid) => (
                  <SelectItem key={String(gid)} label={`店铺 ${gid}`} value={String(gid)}>
                    店铺 {gid}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </View>

        <View className="mb-3">
          <Text className="mb-2">圈子 *</Text>
          <Select
            value={groupId ? ({ label: groupOptions?.find(g => String(g.id) === groupId)?.name || `圈子 ${groupId}`, value: groupId } as any) : undefined}
            onValueChange={(opt) => setGroupId(String(opt?.value ?? ''))}
          >
            <SelectTrigger ref={groupRef} onTouchStart={onGroupTouchStart} className="min-w-[160px]">
              <SelectValue placeholder={groupId ? `圈子 ${groupId}` : '选择圈子'} />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>圈子</SelectLabel>
                {(groupOptions ?? []).map((group) => (
                  <SelectItem key={String(group.id)} label={group.name} value={String(group.id)}>
                    {group.name}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </View>

        <View className="mb-3">
          <Text className="mb-2">最小余额(元) (可选)</Text>
          <Input
            placeholder="例如: 100"
            value={minYuan}
            onChangeText={setMinYuan}
            keyboardType="numeric"
          />
        </View>

        <View className="mb-3">
          <Text className="mb-2">最大余额(元) (可选)</Text>
          <Input
            placeholder="例如: 10000"
            value={maxYuan}
            onChangeText={setMaxYuan}
            keyboardType="numeric"
          />
        </View>

        <Button
          onPress={handleLoadBalances}
          disabled={loadingBalances}
        >
          <Text>{loadingBalances ? '查询中...' : '查询余额'}</Text>
        </Button>
      </Card>

      {/* 统计信息 */}
      {balancesData && balancesData.list.length > 0 && (
        <Card className="p-4 mb-4">
          <Text className="text-lg font-bold mb-3">统计信息</Text>
          <View className="flex-row justify-between mb-2">
            <Text>成员数量:</Text>
            <Text className="font-bold">{balancesData.list.length}</Text>
          </View>
          <View className="flex-row justify-between mb-2">
            <Text>总余额:</Text>
            <Text className="font-bold text-primary">
              ¥{getTotalBalance().toFixed(2)}
            </Text>
          </View>
          <View className="flex-row justify-between">
            <Text>平均余额:</Text>
            <Text className="font-bold">
              ¥{getAvgBalance().toFixed(2)}
            </Text>
          </View>
        </Card>
      )}

      {/* 余额列表 */}
      {balancesData && (
        <Card className="p-4">
          <View className="flex-row justify-between items-center mb-3">
            <Text className="text-lg font-bold">成员余额</Text>
            <Text className="text-sm text-muted-foreground">
              共 {balancesData.total} 人
            </Text>
          </View>

          {balancesData.list.length === 0 ? (
            <Text className="text-center text-muted-foreground py-8">
              暂无成员余额记录
            </Text>
          ) : (
            <View>
              {balancesData.list.map((balance: MemberBalance, index: number) => (
                <View
                  key={`${balance.member_id}-${index}`}
                  className="border-b border-border py-4 last:border-b-0"
                >
                  <View className="flex-row justify-between items-start mb-2">
                    <View className="flex-1">
                      <Text className="font-medium mb-1">
                        {balance.game_name || `用户${balance.game_id}`}
                      </Text>
                      <Text className="text-sm text-muted-foreground">
                        游戏ID: {balance.game_id}
                      </Text>
                      <Text className="text-sm text-muted-foreground">
                        成员ID: {balance.member_id}
                      </Text>
                    </View>
                    <View className="items-end">
                      <Text className="text-2xl font-bold text-primary">
                        ¥{balance.balance_yuan.toFixed(2)}
                      </Text>
                      <Text className="text-xs text-muted-foreground">
                        {balance.balance} 分
                      </Text>
                    </View>
                  </View>
                  <View className="flex-row justify-between">
                    <Text className="text-sm text-muted-foreground">
                      {balance.group_name}
                    </Text>
                    <Text className="text-sm text-muted-foreground">
                      更新: {formatDate(balance.updated_at)}
                    </Text>
                  </View>
                </View>
              ))}

              {/* 分页 */}
              {balancesData.total > 20 && (
                <View className="flex-row justify-center gap-2 mt-4">
                  <Button
                    variant="outline"
                    onPress={() => {
                      if (page > 1) {
                        setPage(page - 1);
                        handleLoadBalances();
                      }
                    }}
                    disabled={page === 1 || loadingBalances}
                  >
                    <Text>上一页</Text>
                  </Button>
                  <Text className="py-2 px-4">
                    第 {page} 页
                  </Text>
                  <Button
                    variant="outline"
                    onPress={() => {
                      if (page * 20 < balancesData.total) {
                        setPage(page + 1);
                        handleLoadBalances();
                      }
                    }}
                    disabled={page * 20 >= balancesData.total || loadingBalances}
                  >
                    <Text>下一页</Text>
                  </Button>
                </View>
              )}
            </View>
          )}
        </Card>
      )}

      {/* 说明 */}
      <Card className="p-4 mt-4 bg-muted">
        <Text className="text-sm text-muted-foreground">
          提示：
        </Text>
        <Text className="text-sm text-muted-foreground mt-2">
          • 可以通过最小/最大余额筛选成员
        </Text>
        <Text className="text-sm text-muted-foreground mt-1">
          • 余额以元为单位显示，实际存储为分
        </Text>
        <Text className="text-sm text-muted-foreground mt-1">
          • 仅显示您管理的圈子成员
        </Text>
      </Card>
    </ScrollView>
  );
};

