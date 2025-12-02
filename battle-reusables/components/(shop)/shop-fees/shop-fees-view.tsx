import React, { useState, useMemo, useEffect, useCallback } from 'react';
import { View, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Icon } from '@/components/ui/icon';
import { X, Plus } from 'lucide-react-native';
import { showToast } from '@/utils/toast';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useHouseSelector } from '@/hooks/use-house-selector';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { MobileSelect } from '@/components/ui/mobile-select';
import { getShopFees, setShopFees } from '@/services/game/shop-fees';

export function ShopFeesView() {
  const {
    houseGid,
    setHouseGid,
    isSuperAdmin,
    isStoreAdmin,
    houseOptions,
  } = useHouseSelector();

  const [loading, setLoading] = useState(false);
  const [feesConfig, setFeesConfig] = useState<API.FeesConfig>({ rules: [] });
  
  // 费用规则表单
  const [threshold, setThreshold] = useState('');
  const [fee, setFee] = useState('');
  const [gameKind, setGameKind] = useState<string | undefined>(undefined);
  const [baseScore, setBaseScore] = useState('');

  // 获取游戏类型常量
  const { data: plazaData } = usePlazaConsts();
  const gameKindSelectOptions = useMemo(() => {
    const options: { label: string; value: string }[] = [
      { label: '全局默认（所有游戏）', value: '0' }
    ];
    if (plazaData?.game_kinds) {
      plazaData.game_kinds.forEach(item => {
        options.push({ label: item.label, value: String(item.value) });
      });
    }
    return options;
  }, [plazaData]);

  // 加载店铺费用配置
  const loadFees = useCallback(async () => {
    if (!houseGid) return;
    
    try {
      setLoading(true);
      const res = await getShopFees({ house_gid: Number(houseGid) });
      
      if (res.data?.fees_json) {
        try {
          const config = JSON.parse(res.data.fees_json) as API.FeesConfig;
          setFeesConfig(config);
        } catch {
          setFeesConfig({ rules: [] });
        }
      } else {
        setFeesConfig({ rules: [] });
      }
    } catch (error: any) {
      if (error?.message?.includes('record not found')) {
        // 店铺还没有配置，显示空配置
        setFeesConfig({ rules: [] });
      } else {
        showToast('加载失败', 'error');
        console.error('加载店铺费用配置失败:', error);
      }
    } finally {
      setLoading(false);
    }
  }, [houseGid]);

  // 店铺改变时加载配置
  useEffect(() => {
    loadFees();
  }, [loadFees]);

  // 添加规则
  const handleAddRule = () => {
    if (!threshold || !fee) {
      showToast('请输入分数阈值和费用', 'error');
      return;
    }

    const newRule: API.FeeRule = {
      threshold: Number(threshold),
      fee: Number(fee),
    };

    // 如果指定了游戏类型，必须同时指定底分
    if (gameKind && gameKind !== '0') {
      if (!baseScore || baseScore === '0') {
        showToast('指定游戏类型时必须设置底分', 'error');
        return;
      }
      newRule.kind = gameKind;  // 使用字符串
      newRule.base = Number(baseScore);
    }
    // 全局规则：不设置 kind 和 base 字段

    setFeesConfig(prev => ({
      rules: [...prev.rules, newRule]
    }));

    // 清空表单
    setThreshold('');
    setFee('');
    setGameKind(undefined);
    setBaseScore('');
    
    showToast('规则已添加，请点击保存', 'success');
  };

  // 删除规则
  const handleDeleteRule = (index: number) => {
    setFeesConfig(prev => ({
      rules: prev.rules.filter((_, i) => i !== index)
    }));
    showToast('规则已删除，请点击保存', 'success');
  };

  // 保存配置
  const handleSave = async () => {
    if (!houseGid) {
      showToast('请选择店铺', 'error');
      return;
    }

    try {
      setLoading(true);
      const feesJSON = JSON.stringify(feesConfig);
      
      await setShopFees({
        house_gid: Number(houseGid),
        fees_json: feesJSON,
      });

      showToast('保存成功', 'success');
      loadFees(); // 重新加载
    } catch (error) {
      showToast('保存失败', 'error');
      console.error('保存店铺费用配置失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <View className="flex-1 bg-background">
      <ScrollView className="flex-1 p-4">
        {/* 店铺选择 */}
        <Card className="mb-4 p-4">
          <Text className="text-base font-semibold mb-3">选择店铺</Text>
          <View className="mb-3">
            <Text className="text-sm text-muted-foreground mb-1">店铺号 *</Text>
            {isSuperAdmin ? (
              <MobileSelect
                value={houseGid}
                placeholder="请选择店铺"
                options={houseOptions}
                onValueChange={(value) => setHouseGid(value)}
                className="w-full"
              />
            ) : (
              <Input
                value={houseGid}
                editable={false}
                className="bg-muted"
              />
            )}
          </View>
        </Card>

        {/* 添加费用规则 */}
        <PermissionGate anyOf={['shop:fees:update']}>
          <Card className="mb-4 p-4">
            <Text className="text-base font-semibold mb-3">添加费用规则</Text>

            <View className="mb-3">
              <Text className="text-sm text-muted-foreground mb-1">
                游戏类型
              </Text>
              <MobileSelect
                value={gameKind}
                placeholder="全局默认（所有游戏）"
                options={gameKindSelectOptions}
                onValueChange={(value) => setGameKind(value)}
                className="w-full"
              />
            </View>

            <View className="mb-3">
              <Text className="text-sm text-muted-foreground mb-1">
                底分（0表示不限）
              </Text>
              <Input
                placeholder="0"
                keyboardType="numeric"
                value={baseScore}
                onChangeText={setBaseScore}
              />
            </View>

            <View className="mb-3">
              <Text className="text-sm text-muted-foreground mb-1">
                分数阈值 *
              </Text>
              <Input
                placeholder="如：100（表示达到100分时收费）"
                keyboardType="numeric"
                value={threshold}
                onChangeText={setThreshold}
              />
            </View>

            <View className="mb-3">
              <Text className="text-sm text-muted-foreground mb-1">
                费用（单位：分）*
              </Text>
              <Input
                placeholder="如：1000（表示10元）"
                keyboardType="numeric"
                value={fee}
                onChangeText={setFee}
              />
            </View>

            <Button
              onPress={handleAddRule}
              disabled={!threshold || !fee || loading}
              variant="outline"
            >
              <Icon as={Plus} size={16} className="mr-2" />
              <Text>添加规则</Text>
            </Button>

            <View className="mt-3 bg-muted/50 p-3 rounded">
              <Text className="text-xs text-muted-foreground">
                💡 规则说明：
                {'\n'}• 全局规则：不指定游戏类型和底分，适用所有房间
                {'\n'}• 特定规则：指定游戏类型或底分，精确匹配
                {'\n'}• 匹配优先级：全局规则优先，然后是特定规则
                {'\n'}• 分数阈值：最高分达到该值时收取费用
              </Text>
            </View>
          </Card>
        </PermissionGate>

        {/* 当前规则列表 */}
        <PermissionGate anyOf={['shop:fees:view']}>
          <Card className="mb-4 p-4">
            <View className="flex-row items-center justify-between mb-3">
              <Text className="text-base font-semibold">当前费用规则</Text>
              {feesConfig.rules.length > 0 && (
                <PermissionGate anyOf={['shop:fees:update']}>
                  <Button
                    variant="default"
                    size="sm"
                    onPress={handleSave}
                    disabled={loading}
                  >
                    <Text className="text-xs">保存配置</Text>
                  </Button>
                </PermissionGate>
              )}
            </View>

            {!houseGid ? (
              <View className="items-center justify-center py-8">
                <Text className="text-muted-foreground">请选择店铺</Text>
              </View>
            ) : loading ? (
              <View className="items-center justify-center py-8">
                <Text className="text-muted-foreground">加载中...</Text>
              </View>
            ) : feesConfig.rules.length === 0 ? (
              <View className="items-center justify-center py-8">
                <Text className="text-muted-foreground">暂无费用规则</Text>
              </View>
            ) : (
              <View className="gap-2">
                {feesConfig.rules.map((rule, index) => (
                  <View
                    key={index}
                    className="border border-border rounded p-3 bg-card"
                  >
                    <View className="flex-row justify-between items-start mb-2">
                      <View className="flex-1">
                        <Text className="text-sm font-medium">
                          {rule.kind ? `游戏类型 ${rule.kind}` : '全局默认'}
                          {rule.base ? ` | 底分 ${rule.base}` : ''}
                        </Text>
                        <Text className="text-xs text-muted-foreground mt-1">
                          分数达到 {rule.threshold} 时收费
                        </Text>
                      </View>
                      <PermissionGate anyOf={['shop:fees:update']}>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-8 px-2"
                          onPress={() => handleDeleteRule(index)}
                          disabled={loading}
                        >
                          <Icon as={X} size={16} className="text-destructive" />
                        </Button>
                      </PermissionGate>
                    </View>
                    <View className="flex-row items-center justify-between">
                      <Text className="text-lg font-semibold text-primary">
                        ¥{(rule.fee / 100).toFixed(2)}
                      </Text>
                      <Text className="text-xs text-muted-foreground">
                        ({rule.fee} 分)
                      </Text>
                    </View>
                  </View>
                ))}
              </View>
            )}
          </Card>
        </PermissionGate>
      </ScrollView>
    </View>
  );
}
