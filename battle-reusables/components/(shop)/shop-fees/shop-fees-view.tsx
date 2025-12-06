import React, { useState, useMemo, useEffect, useCallback } from 'react';
import { View, ScrollView, Switch } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Icon } from '@/components/ui/icon';
import { X, Plus, Pencil } from 'lucide-react-native';
import { showToast } from '@/utils/toast';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useHouseSelector } from '@/hooks/use-house-selector';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { MobileSelect } from '@/components/ui/mobile-select';
import { getShopFees, setShopFees, setShareFee } from '@/services/game/shop-fees';

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
  const [shareEnabled, setShareEnabled] = useState(false);
  
  // 费用规则表单
  const [threshold, setThreshold] = useState('');
  const [fee, setFee] = useState('');
  const [gameKind, setGameKind] = useState<string | undefined>(undefined);
  const [baseScore, setBaseScore] = useState('');
  const [editingIndex, setEditingIndex] = useState<number | null>(null); // 编辑中的规则索引

  // 获取游戏类型常量
  const { data: plazaData, maps } = usePlazaConsts();
  const gameKindSelectOptions = useMemo(() => {
    const options: { label: string; value: string }[] = [];
    if (plazaData?.game_kinds) {
      plazaData.game_kinds.forEach(item => {
        options.push({ label: item.label, value: String(item.value) });
      });
    }
    return options;
  }, [plazaData]);

  // 获取游戏类型名称
  const getGameKindName = (kindCode: string | number | undefined) => {
    // 未设置或为 0，表示全部游戏
    if (kindCode === undefined || kindCode === null || kindCode === '') {
      return '全部游戏';
    }

    const code = typeof kindCode === 'string' ? Number(kindCode) : kindCode;
    if (code === 0 || isNaN(code)) return '全部游戏';

    const name = maps.game_kinds.get(code);
    if (name) return name;

    // 常见游戏类型硬编码备用
    const fallback: Record<number, string> = {
      60: '血战到底',
      61: '血战换三张', 
      70: '跑得快',
      80: '斗地主',
    };
    return fallback[code] || `游戏${code}`;
  };

  // 加载店铺费用配置
  const loadFees = useCallback(async () => {
    if (!houseGid) return;
    
    try {
      setLoading(true);
      const res = await getShopFees({ house_gid: Number(houseGid) });

      if (res.data) {
        setShareEnabled(!!res.data.share_fee);
      }

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

  // 切换分运开关
  const handleToggleShareFee = async (value: boolean) => {
    if (!houseGid) {
      showToast('请选择店铺', 'error');
      return;
    }

    try {
      setLoading(true);
      await setShareFee({
        house_gid: Number(houseGid),
        share: value,
      });
      setShareEnabled(value);
      showToast(value ? '已开启分运' : '已关闭分运', 'success');
    } catch (error) {
      showToast('更新分运设置失败', 'error');
      console.error('更新分运设置失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 开始编辑规则
  const handleEditRule = (index: number) => {
    const rule = feesConfig.rules[index];
    setEditingIndex(index);
    setThreshold(String(rule.threshold));
    setFee(String(rule.fee));
    setGameKind(rule.kind || undefined);
    setBaseScore(String(rule.base || ''));
  };

  // 取消编辑
  const handleCancelEdit = () => {
    setEditingIndex(null);
    setThreshold('');
    setFee('');
    setGameKind(undefined);
    setBaseScore('');
  };

  // 添加/保存规则
  const handleSaveRule = () => {
    // 分数阈值与费用必须是有效的正数
    const thresholdNum = Number(threshold);
    const feeNum = Number(fee);
    const baseNum = baseScore ? Number(baseScore) : 0;

    if (!threshold || Number.isNaN(thresholdNum) || thresholdNum <= 0) {
      showToast('请输入正确的分数阈值', 'error');
      return;
    }
    if (!fee || Number.isNaN(feeNum) || feeNum <= 0) {
      showToast('请输入正确的费用', 'error');
      return;
    }

    const newRule: API.FeeRule = {
      threshold: thresholdNum,
      fee: feeNum,
    };

    // 如果选择了游戏类型，则必须设置大于 0 的底分，按具体【游戏类型 + 底分】生效
    if (gameKind) {
      if (Number.isNaN(baseNum) || baseNum <= 0) {
        showToast('选择了游戏类型时必须设置大于 0 的底分', 'error');
        return;
      }
      newRule.kind = gameKind;
      newRule.base = baseNum;
    } else {
      // 未选择游戏类型：可选设置底分>0，表示“所有游戏但指定底分”，否则为完全全局规则
      if (!Number.isNaN(baseNum) && baseNum > 0) {
        newRule.base = baseNum;
      }
    }

    setFeesConfig(prev => {
      const rules = [...prev.rules];
      if (editingIndex !== null && editingIndex >= 0 && editingIndex < rules.length) {
        rules[editingIndex] = newRule;
      } else {
        rules.push(newRule);
      }
      return { ...prev, rules };
    });

    showToast(editingIndex !== null ? '规则已修改，请点击保存' : '规则已添加，请点击保存', 'success');

    // 清空表单
    setEditingIndex(null);
    setThreshold('');
    setFee('');
    setGameKind(undefined);
    setBaseScore('');
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

        {/* 分运设置 */}
        <PermissionGate anyOf={['shop:sharefee:write', 'shop:fees:update']}>
          <Card className="mb-4 p-4 flex-row items-center justify-between">
            <View>
              <Text className="text-base font-semibold mb-1">分运设置</Text>
              <Text className="text-xs text-muted-foreground">
                开启后，运费在各圈之间按规则分摊
              </Text>
            </View>
            <Switch
              value={shareEnabled}
              onValueChange={handleToggleShareFee}
              disabled={!houseGid || loading}
            />
          </Card>
        </PermissionGate>

        {/* 添加/编辑费用规则 */}
        <PermissionGate anyOf={['shop:fees:update']}>
          <Card className="mb-4 p-4">
            <View className="flex-row items-center justify-between mb-3">
              <Text className="text-base font-semibold">
                {editingIndex !== null ? `编辑规则 #${editingIndex + 1}` : '添加费用规则'}
              </Text>
              {editingIndex !== null && (
                <Button variant="ghost" size="sm" onPress={handleCancelEdit}>
                  <Icon as={X} size={16} className="text-muted-foreground" />
                  <Text className="text-xs text-muted-foreground ml-1">取消</Text>
                </Button>
              )}
            </View>

            <View className="mb-3">
              <Text className="text-sm text-muted-foreground mb-1">
                游戏类型（可选，留空表示全部）
              </Text>
              <MobileSelect
                value={gameKind}
                placeholder="请选择游戏类型（留空表示全部）"
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
              onPress={handleSaveRule}
              disabled={!threshold || !fee || loading}
              variant="outline"
            >
              {editingIndex === null && <Icon as={Plus} size={16} className="mr-2" />}
              <Text>{editingIndex !== null ? '保存修改' : '添加规则'}</Text>
            </Button>

            <View className="mt-3 bg-muted/50 p-3 rounded">
              <Text className="text-xs text-muted-foreground">
                规则说明：
                {'\n'}• 游戏类型可选：留空表示全部游戏
                {'\n'}• 底分可选：0 或留空表示不限底分
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
                          {getGameKindName(rule.kind)}
                          {rule.base ? ` | 底分 ${rule.base}` : ''}
                        </Text>
                        <Text className="text-xs text-muted-foreground mt-1">
                          分数达到 {rule.threshold} 时收费
                        </Text>
                      </View>
                      <PermissionGate anyOf={['shop:fees:update']}>
                        <View className="flex-row">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-8 px-2"
                            onPress={() => handleEditRule(index)}
                            disabled={loading}
                          >
                            <Icon as={Pencil} size={14} className="text-primary" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-8 px-2"
                            onPress={() => handleDeleteRule(index)}
                            disabled={loading}
                          >
                            <Icon as={X} size={16} className="text-destructive" />
                          </Button>
                        </View>
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
