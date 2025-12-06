import React, { useState, useMemo, useEffect, useCallback } from 'react';
import { View, ScrollView } from 'react-native';
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
import { listRoomCreditLimits, setRoomCreditLimit, deleteRoomCreditLimit } from '@/services/game/room-credit';
import { getGroupOptions } from '@/services/shops/groups';

export function RoomCreditView() {
  const {
    houseGid,
    setHouseGid,
    isSuperAdmin,
    isStoreAdmin,
    houseOptions,
  } = useHouseSelector();

  const [creditLimits, setCreditLimits] = useState<API.RoomCreditLimitItem[]>([]);
  const [creditLoading, setCreditLoading] = useState(false);
  const [groupOptions, setGroupOptions] = useState<{ label: string; value: string }[]>([]);
  
  // 额度表单
  const [creditAmount, setCreditAmount] = useState(''); // 额度（元）
  const [creditGameKind, setCreditGameKind] = useState<string | undefined>(undefined);
  const [creditBaseScore, setCreditBaseScore] = useState('');
  const [creditGroupName, setCreditGroupName] = useState<string | undefined>(undefined); // 圈子名称
  const [editingCreditIndex, setEditingCreditIndex] = useState<number | null>(null);

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

  // 加载额度限制列表
  const loadCreditLimits = useCallback(async () => {
    if (!houseGid) {
      setCreditLimits([]);
      return;
    }
    
    try {
      setCreditLoading(true);
      const res = await listRoomCreditLimits({ house_gid: Number(houseGid) });
      if (res.data?.items) {
        setCreditLimits(res.data.items);
      } else {
        setCreditLimits([]);
      }
    } catch (error: any) {
      console.error('加载额度限制失败:', error);
      // 不显示错误提示，因为可能是后端服务未启动或路由未注册
      setCreditLimits([]);
    } finally {
      setCreditLoading(false);
    }
  }, [houseGid]);

  // 加载圈子选项
  const loadGroupOptions = useCallback(async () => {
    if (!houseGid) {
      setGroupOptions([{ label: '全局', value: '' }]);
      return;
    }
    
    try {
      const res = await getGroupOptions({ house_gid: Number(houseGid) });
      if (res.data) {
        // 添加"全局"选项（空字符串表示全局）
        const options: { label: string; value: string }[] = [
          { label: '全局', value: '' }
        ];
        res.data.forEach(group => {
          options.push({ label: group.name, value: group.name });
        });
        setGroupOptions(options);
      } else {
        setGroupOptions([{ label: '全局', value: '' }]);
      }
    } catch (error) {
      console.error('加载圈子选项失败:', error);
      // 即使失败也显示"全局"选项
      setGroupOptions([{ label: '全局', value: '' }]);
    }
  }, [houseGid]);

  // 店铺改变时加载配置
  useEffect(() => {
    loadCreditLimits();
    loadGroupOptions();
  }, [loadCreditLimits, loadGroupOptions]);

  // 开始编辑额度
  const handleEditCredit = (index: number) => {
    const credit = creditLimits[index];
    setEditingCreditIndex(index);
    setCreditAmount(String(credit.credit_yuan));
    setCreditGameKind(credit.game_kind > 0 ? String(credit.game_kind) : undefined);
    setCreditBaseScore(credit.base_score > 0 ? String(credit.base_score) : '');
    // 空字符串或 undefined 都显示为"全局"（空字符串）
    setCreditGroupName(credit.group_name || '');
  };

  // 取消编辑额度
  const handleCancelEditCredit = () => {
    setEditingCreditIndex(null);
    setCreditAmount('');
    setCreditGameKind(undefined);
    setCreditBaseScore('');
    setCreditGroupName(undefined); // undefined 会在 MobileSelect 中显示为空字符串（全局）
  };

  // 保存额度设置
  const handleSaveCredit = async () => {
    if (!houseGid) {
      showToast('请选择店铺', 'error');
      return;
    }

    const amountNum = Number(creditAmount);
    if (!creditAmount || Number.isNaN(amountNum) || amountNum <= 0) {
      showToast('请输入正确的额度（元）', 'error');
      return;
    }

    const baseNum = creditBaseScore ? Number(creditBaseScore) : 0;
    const kindNum = creditGameKind ? Number(creditGameKind) : 0;

    // 如果选择了游戏类型，则必须设置大于 0 的底分
    if (creditGameKind && (Number.isNaN(baseNum) || baseNum <= 0)) {
      showToast('选择了游戏类型时必须设置大于 0 的底分', 'error');
      return;
    }

    try {
      setCreditLoading(true);
      const creditLimitInCents = Math.round(amountNum * 100); // 转换为分

      await setRoomCreditLimit({
        house_gid: Number(houseGid),
        group_name: creditGroupName || undefined,
        game_kind: kindNum || undefined,
        base_score: baseNum || undefined,
        credit_limit: creditLimitInCents,
      });

      showToast('额度设置成功', 'success');
      loadCreditLimits(); // 重新加载
      handleCancelEditCredit(); // 清空表单
    } catch (error) {
      showToast('设置额度失败', 'error');
      console.error('设置额度失败:', error);
    } finally {
      setCreditLoading(false);
    }
  };

  // 删除额度
  const handleDeleteCredit = async (index: number) => {
    const credit = creditLimits[index];
    if (!houseGid) return;

    try {
      setCreditLoading(true);
      await deleteRoomCreditLimit({
        house_gid: Number(houseGid),
        group_name: credit.group_name || undefined,
        game_kind: credit.game_kind || undefined,
        base_score: credit.base_score || undefined,
      });

      showToast('删除成功', 'success');
      loadCreditLimits(); // 重新加载
    } catch (error) {
      showToast('删除失败', 'error');
      console.error('删除额度失败:', error);
    } finally {
      setCreditLoading(false);
    }
  };

  // 格式化额度显示
  const formatCreditDisplay = (credit: API.RoomCreditLimitItem) => {
    const parts: string[] = [];
    if (credit.group_name) {
      parts.push(`圈：${credit.group_name}`);
    }
    if (credit.game_kind > 0) {
      parts.push(getGameKindName(credit.game_kind));
    }
    if (credit.base_score > 0) {
      parts.push(`底分 ${credit.base_score}`);
    }
    if (parts.length === 0) {
      parts.push('全局默认');
    }
    return `🈲 ${credit.credit_yuan}元 (${parts.join(' / ')})`;
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

        {/* 额度设置 */}
        <PermissionGate anyOf={['room:credit:set', 'room:credit:view']}>
          <Card className="mb-4 p-4">
            <View className="flex-row items-center justify-between mb-3">
              <Text className="text-base font-semibold">额度设置</Text>
              <Text className="text-xs text-muted-foreground">
                设置玩家进入房间所需的最低余额
              </Text>
            </View>

            {/* 添加/编辑额度表单 */}
            <PermissionGate anyOf={['room:credit:set']}>
              <View className="mb-4 p-3 bg-muted/30 rounded">
                <View className="flex-row items-center justify-between mb-3">
                  <Text className="text-sm font-medium">
                    {editingCreditIndex !== null ? `编辑额度 #${editingCreditIndex + 1}` : '添加额度设置'}
                  </Text>
                  {editingCreditIndex !== null && (
                    <Button variant="ghost" size="sm" onPress={handleCancelEditCredit}>
                      <Icon as={X} size={16} className="text-muted-foreground" />
                      <Text className="text-xs text-muted-foreground ml-1">取消</Text>
                    </Button>
                  )}
                </View>

                <View className="mb-3">
                  <Text className="text-sm text-muted-foreground mb-1">
                    圈子名称（可选，选择"全局"表示全局设置）
                  </Text>
                  <MobileSelect
                    value={creditGroupName || ''}
                    placeholder='请选择圈子（选择"全局"表示全局设置）'
                    options={groupOptions}
                    onValueChange={(value) => setCreditGroupName(value === '' ? undefined : value)}
                    className="w-full"
                  />
                </View>

                <View className="mb-3">
                  <Text className="text-sm text-muted-foreground mb-1">
                    游戏类型（可选，留空表示全部）
                  </Text>
                  <MobileSelect
                    value={creditGameKind}
                    placeholder="请选择游戏类型（留空表示全部）"
                    options={gameKindSelectOptions}
                    onValueChange={(value) => setCreditGameKind(value)}
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
                    value={creditBaseScore}
                    onChangeText={setCreditBaseScore}
                  />
                </View>

                <View className="mb-3">
                  <Text className="text-sm text-muted-foreground mb-1">
                    额度（单位：元）*
                  </Text>
                  <Input
                    placeholder="如：399（表示399元）"
                    keyboardType="numeric"
                    value={creditAmount}
                    onChangeText={setCreditAmount}
                  />
                </View>

                <Button
                  onPress={handleSaveCredit}
                  disabled={!creditAmount || creditLoading}
                  variant="outline"
                >
                  {editingCreditIndex === null && <Icon as={Plus} size={16} className="mr-2" />}
                  <Text>{editingCreditIndex !== null ? '保存修改' : '添加额度'}</Text>
                </Button>

                <View className="mt-3 bg-muted/50 p-3 rounded">
                  <Text className="text-xs text-muted-foreground">
                    额度说明：
                    {'\n'}• 玩家余额必须达到设置的额度才能进入房间
                    {'\n'}• 圈子名称：留空表示全局设置，填写表示该圈子的设置
                    {'\n'}• 游戏类型：留空表示全部游戏，选择后必须设置底分
                    {'\n'}• 底分：0 或留空表示不限底分
                    {'\n'}• 优先级：圈子+游戏类型+底分 {'>'} 圈子默认 {'>'} 全局+游戏类型+底分 {'>'} 全局默认
                  </Text>
                </View>
              </View>
            </PermissionGate>

            {/* 当前额度列表 */}
            <PermissionGate anyOf={['room:credit:view']}>
              <View>
                <Text className="text-sm font-semibold mb-3">当前额度设置</Text>
                {!houseGid ? (
                  <View className="items-center justify-center py-8">
                    <Text className="text-muted-foreground">请选择店铺</Text>
                  </View>
                ) : creditLoading ? (
                  <View className="items-center justify-center py-8">
                    <Text className="text-muted-foreground">加载中...</Text>
                  </View>
                ) : creditLimits.length === 0 ? (
                  <View className="items-center justify-center py-8">
                    <Text className="text-muted-foreground">暂无额度设置</Text>
                  </View>
                ) : (
                  <View className="gap-2">
                    {creditLimits.map((credit, index) => (
                      <View
                        key={credit.id}
                        className="border border-border rounded p-3 bg-card"
                      >
                        <View className="flex-row justify-between items-start mb-2">
                          <View className="flex-1">
                            <Text className="text-sm font-medium">
                              {formatCreditDisplay(credit)}
                            </Text>
                            <Text className="text-xs text-muted-foreground mt-1">
                              更新时间：{new Date(credit.updated_at).toLocaleString('zh-CN')}
                            </Text>
                          </View>
                          <PermissionGate anyOf={['room:credit:set']}>
                            <View className="flex-row">
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 px-2"
                                onPress={() => handleEditCredit(index)}
                                disabled={creditLoading}
                              >
                                <Icon as={Pencil} size={14} className="text-primary" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 px-2"
                                onPress={() => handleDeleteCredit(index)}
                                disabled={creditLoading}
                              >
                                <Icon as={X} size={16} className="text-destructive" />
                              </Button>
                            </View>
                          </PermissionGate>
                        </View>
                      </View>
                    ))}
                  </View>
                )}
              </View>
            </PermissionGate>
          </Card>
        </PermissionGate>
      </ScrollView>
    </View>
  );
}

