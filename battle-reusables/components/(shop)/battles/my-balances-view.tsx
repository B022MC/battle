import React from 'react';
import { ScrollView, View } from 'react-native';
import { useRequest } from '@/hooks/use-request';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { alert } from '@/utils/alert';
import { getMyBalances } from '@/services/battles/query';
import type { MemberBalance } from '@/services/battles/query-typing';
import { useHouseSelector } from '@/hooks/use-house-selector';

export const MyBalancesView = () => {
  // 自动获取当前店铺（店铺管理员固定，超管可选择）
  const { houseGid, isReady } = useHouseSelector();

  // API 请求
  const { data: balancesData, loading: loadingBalances, run: runGetBalances } = useRequest(getMyBalances, { manual: true });

  // 加载余额（只需要 house_gid，game_id 由后端根据绑定账号解析）
  const handleLoadBalances = async () => {
    if (!isReady || !houseGid) {
      alert.show({ title: '未获取到店铺号', variant: 'error' });
      return;
    }
    const gid = Number(houseGid);
    if (!gid) {
      alert.show({ title: '未获取到店铺号', variant: 'error' });
      return;
    }
    if (isNaN(gid) || gid <= 0) {
      alert.show({ title: '店铺号格式错误', variant: 'error' });
      return;
    }

    const params = { house_gid: gid };
    try {
      await runGetBalances(params);
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
    const list = balancesData?.balances ?? [];
    return list.reduce((sum, b) => sum + b.balance_yuan, 0);
  };

  return (
    <ScrollView className="flex-1 bg-background p-4">
      {/* 查询表单（自动使用当前店铺，无需选择） */}
      <Card className="p-4 mb-4">
        <Text className="text-lg font-bold mb-4">查询我的余额</Text>
        <Text className="text-sm text-muted-foreground mb-3">
          自动使用当前店铺和当前圈子，无需选择
        </Text>

        <Button
          onPress={handleLoadBalances}
          disabled={loadingBalances}
        >
          <Text>{loadingBalances ? '查询中...' : '查询余额'}</Text>
        </Button>
      </Card>

      {/* 总余额 */}
      {balancesData?.balances && balancesData.balances.length > 0 && (
        <Card className="p-4 mb-4 bg-primary">
          <Text className="text-sm text-primary-foreground mb-1">总余额</Text>
          <Text className="text-3xl font-bold text-primary-foreground">
            ¥{getTotalBalance().toFixed(2)}
          </Text>
          <Text className="text-sm text-primary-foreground mt-1">
            {balancesData.balances.length} 个圈子
          </Text>
        </Card>
      )}

      {/* 余额列表 */}
      {balancesData?.balances && (
        <Card className="p-4">
          <Text className="text-lg font-bold mb-3">余额明细</Text>

          {balancesData.balances.length === 0 ? (
            <Text className="text-center text-muted-foreground py-8">
              暂无余额记录
            </Text>
          ) : (
            <View>
              {balancesData.balances.map((balance: MemberBalance, index: number) => (
                <View
                  key={`${balance.member_id}-${balance.group_id || 'default'}-${index}`}
                  className="border-b border-border py-4 last:border-b-0"
                >
                  <View className="flex-row justify-between items-start mb-2">
                    <View className="flex-1">
                      <Text className="font-medium mb-1">
                        {balance.group_name || '默认圈子'}
                      </Text>
                      <Text className="text-sm text-muted-foreground">
                        游戏ID: {balance.game_id}
                      </Text>
                      {balance.game_name && (
                        <Text className="text-sm text-muted-foreground">
                          昵称: {balance.game_name}
                        </Text>
                      )}
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
                    {balance.group_id && (
                      <Text className="text-sm text-muted-foreground">
                        圈子ID: {balance.group_id}
                      </Text>
                    )}
                    <Text className="text-sm text-muted-foreground">
                      更新: {formatDate(balance.updated_at)}
                    </Text>
                  </View>
                </View>
              ))}
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
          • 不填写圈子ID将显示您在该店铺所有圈子的余额
        </Text>
        <Text className="text-sm text-muted-foreground mt-1">
          • 填写圈子ID将只显示该圈子的余额
        </Text>
        <Text className="text-sm text-muted-foreground mt-1">
          • 余额以元为单位显示，实际存储为分
        </Text>
      </Card>
    </ScrollView>
  );
};

