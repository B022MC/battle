import React, { useState, useEffect, useMemo, useRef } from 'react';
import { View, ScrollView, ActivityIndicator, RefreshControl, Platform, Pressable, Modal } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Icon } from '@/components/ui/icon';
import { Switch } from 'react-native';
import { Trash2, Plus, ChevronDown, Pencil, X } from 'lucide-react-native';
import { shopsFeesGet, shopsShareFeeSet, shopsFeesSet } from '@/services/shops/fees';
import { showToast } from '@/utils/toast';
import { PermissionGate } from '@/components/auth/PermissionGate';
import { useHouseSelector } from '@/hooks/use-house-selector';
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

type GameFee = {
  kind: string;        // 游戏类型ID ("0"表示所有)
  base: number;        // 底分 (0表示所有)
  threshold: number;   // 门槛金额（分）
  fee: number;         // 运费金额（分）
};

export function FeesView() {
  const {
    houseGid,
    setHouseGid,
    isSuperAdmin,
    isStoreAdmin,
    houseOptions,
  } = useHouseSelector();
  
  const { maps } = usePlazaConsts();
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [open, setOpen] = useState(false);
  const inputRef = useRef<View>(null);
  const [dropdownLayout, setDropdownLayout] = useState({ top: 0, left: 0, width: 0 });
  
  // 设置数据
  const [shareFee, setShareFee] = useState(false);
  const [pushCredit, setPushCredit] = useState('');
  const [gameFees, setGameFees] = useState<GameFee[]>([]);
  
  // 新增运费表单
  const [newThreshold, setNewThreshold] = useState('');
  const [newFee, setNewFee] = useState('');
  const [newGameKind, setNewGameKind] = useState<string | undefined>(undefined);
  const [newBaseScore, setNewBaseScore] = useState('');
  
  // 操作状态
  const [updating, setUpdating] = useState(false);
  const [editingIndex, setEditingIndex] = useState<number | null>(null); // 正在编辑的规则索引
  
  // 游戏类型选项
  const gameKinds = useMemo(() => {
    const kinds: { label: string; value: string }[] = [];
    maps.game_kinds.forEach((label, value) => {
      kinds.push({ label, value: String(value) });
    });
    return kinds;
  }, [maps]);

  // 店铺管理员自动加载数据
  useEffect(() => {
    if (isStoreAdmin && houseGid) {
      loadSettings(true);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isStoreAdmin, houseGid]);

  // 加载店铺设置
  const loadSettings = async (showLoading = true) => {
    if (!houseGid) {
      showToast('请输入店铺号', 'error');
      return;
    }

    try {
      if (showLoading) setLoading(true);
      const res = await shopsFeesGet({ house_gid: Number(houseGid) });
      
      if (res.code === 0 && res.data) {
        setShareFee(res.data.share_fee || false);
        setPushCredit(String(res.data.push_credit || 0));
        
        // 解析运费规则
        try {
          const parsed = res.data.fees_json ? JSON.parse(res.data.fees_json) : { rules: [] };
          const rules = parsed.rules || [];
          setGameFees(Array.isArray(rules) ? rules : []);
        } catch (e) {
          console.error('解析运费规则失败:', e);
          setGameFees([]);
        }
      } else {
        showToast(res.msg || '加载失败', 'error');
      }
    } catch (error) {
      showToast('加载失败', 'error');
      console.error('加载店铺设置失败:', error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  // 下拉刷新
  const handleRefresh = () => {
    setRefreshing(true);
    loadSettings(false);
  };

  // 切换分运费
  const handleToggleShareFee = async (value: boolean) => {
    if (!houseGid) return;

    try {
      setUpdating(true);
      const res = await shopsShareFeeSet({
        house_gid: Number(houseGid),
        share: value,
      });

      if (res.code === 0) {
        setShareFee(value);
        showToast(value ? '已开启分运费' : '已关闭分运费', 'success');
      } else {
        showToast(res.msg || '设置失败', 'error');
      }
    } catch (error) {
      showToast('设置失败', 'error');
    } finally {
      setUpdating(false);
    }
  };

  // 开始编辑规则
  const handleEditFee = (index: number) => {
    const fee = gameFees[index];
    setEditingIndex(index);
    setNewThreshold(String(fee.threshold));
    setNewFee(String(fee.fee));
    setNewGameKind(fee.kind || '0');
    setNewBaseScore(String(fee.base || 0));
  };

  // 取消编辑
  const handleCancelEdit = () => {
    setEditingIndex(null);
    setNewThreshold('');
    setNewFee('');
    setNewGameKind(undefined);
    setNewBaseScore('');
  };

  // 添加/保存运费规则
  const handleSaveFee = async () => {
    if (!houseGid || !newThreshold || !newFee) {
      showToast('请填写门槛和运费', 'error');
      return;
    }

    const threshold = Number(newThreshold);
    const fee = Number(newFee);
    const gameKind = newGameKind ? Number(newGameKind) : 0;
    const baseScore = newBaseScore ? Number(newBaseScore) : 0;

    if (threshold <= 0 || fee <= 0) {
      showToast('门槛和运费必须大于0', 'error');
      return;
    }

    const newRule: GameFee = {
      kind: String(gameKind),
      base: baseScore,
      threshold: threshold,
      fee: fee,
    };

    let updatedFees: GameFee[];
    if (editingIndex !== null) {
      // 编辑模式：替换对应索引的规则
      updatedFees = [...gameFees];
      updatedFees[editingIndex] = newRule;
    } else {
      // 新增模式
      updatedFees = [...gameFees, newRule];
    }

    try {
      setUpdating(true);
      const res = await shopsFeesSet({
        house_gid: Number(houseGid),
        fees_json: JSON.stringify({ rules: updatedFees }),
      });

      if (res.code === 0) {
        setGameFees(updatedFees);
        setNewThreshold('');
        setNewFee('');
        setNewGameKind(undefined);
        setNewBaseScore('');
        setEditingIndex(null);
        showToast(editingIndex !== null ? '修改成功' : '添加成功', 'success');
      } else {
        showToast(res.msg || '保存失败', 'error');
      }
    } catch (error) {
      showToast('保存失败', 'error');
      console.error('保存运费规则失败:', error);
    } finally {
      setUpdating(false);
    }
  };

  // 删除运费规则
  const handleDeleteFee = async (index: number) => {
    if (!houseGid) return;

    const updatedFees = gameFees.filter((_, i) => i !== index);

    try {
      setUpdating(true);
      const res = await shopsFeesSet({
        house_gid: Number(houseGid),
        fees_json: JSON.stringify({ rules: updatedFees }),
      });

      if (res.code === 0) {
        setGameFees(updatedFees);
        showToast('删除成功', 'success');
      } else {
        showToast(res.msg || '删除失败', 'error');
      }
    } catch (error) {
      showToast('删除失败', 'error');
      console.error('删除运费规则失败:', error);
    } finally {
      setUpdating(false);
    }
  };

  // 获取游戏类型名称
  const getGameKindName = (kindCode: string | number) => {
    const code = typeof kindCode === 'string' ? Number(kindCode) : kindCode;
    if (code === 0 || isNaN(code)) return '所有游戏';
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

  // 店铺下拉框选项
  const filtered = useMemo(() => {
    // houseOptions 是 MobileSelectOption[] 类型，需要提取 value 字段
    const list = (houseOptions ?? []).map((opt) => opt.value);
    const q = houseGid.trim();
    if (!q) return list;
    return list.filter((v) => v.includes(q));
  }, [houseOptions, houseGid]);

  return (
    <View className="flex-1 bg-gray-50">
      {/* 头部查询区 */}
      <View className="bg-white p-4 border-b border-gray-200">
        <Text className="text-lg font-semibold mb-3">店铺费用设置</Text>
        
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
          <View className="flex-row gap-2">
            <View className="relative flex-1">
              <View
                ref={inputRef}
                onLayout={(e) => {
                  inputRef.current?.measureInWindow((x, y, width, height) => {
                    setDropdownLayout({ top: y + height, left: x, width });
                  });
                }}
              >
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
              </View>
            </View>
            <Button
              onPress={() => loadSettings(true)}
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

      {/* 店铺下拉列表 Modal */}
      {open && Platform.OS === 'web' && (
        <Pressable
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            zIndex: 9998,
          }}
          onPress={() => setOpen(false)}
        >
          <View
            style={{
              position: 'fixed',
              top: dropdownLayout.top,
              left: dropdownLayout.left,
              width: dropdownLayout.width,
              maxHeight: 224,
              zIndex: 9999,
              backgroundColor: 'white',
              borderRadius: 6,
              borderWidth: 1,
              borderColor: '#e5e7eb',
              shadowColor: '#000',
              shadowOffset: { width: 0, height: 2 },
              shadowOpacity: 0.1,
              shadowRadius: 4,
              overflow: 'hidden',
            }}
          >
            <ScrollView style={{ maxHeight: 224 }}>
              {(filtered.length > 0 ? filtered : ['无匹配结果']).map((gid) => (
                <Pressable
                  key={gid}
                  onPress={() => {
                    if (gid !== '无匹配结果') {
                      setHouseGid(gid);
                      setOpen(false);
                    }
                  }}
                  className="px-3 py-2"
                  accessibilityRole="button"
                >
                  <Text className="text-sm">{gid === '无匹配结果' ? gid : `店铺 ${gid}`}</Text>
                </Pressable>
              ))}
            </ScrollView>
          </View>
        </Pressable>
      )}

      {open && Platform.OS !== 'web' && (
        <Modal transparent visible={open} onRequestClose={() => setOpen(false)}>
          <Pressable
            style={{ flex: 1, backgroundColor: 'rgba(0,0,0,0.3)' }}
            onPress={() => setOpen(false)}
          >
            <View
              style={{
                position: 'absolute',
                top: dropdownLayout.top,
                left: dropdownLayout.left,
                width: dropdownLayout.width,
                maxHeight: 224,
                backgroundColor: 'white',
                borderRadius: 6,
                borderWidth: 1,
                borderColor: '#e5e7eb',
                overflow: 'hidden',
              }}
            >
              <ScrollView style={{ maxHeight: 224 }}>
                {(filtered.length > 0 ? filtered : ['无匹配结果']).map((gid) => (
                  <Pressable
                    key={gid}
                    onPress={() => {
                      if (gid !== '无匹配结果') {
                        setHouseGid(gid);
                        setOpen(false);
                      }
                    }}
                    className="px-3 py-2"
                    accessibilityRole="button"
                  >
                    <Text className="text-sm">{gid === '无匹配结果' ? gid : `店铺 ${gid}`}</Text>
                  </Pressable>
                ))}
              </ScrollView>
            </View>
          </Pressable>
        </Modal>
      )}

      <ScrollView
        className="flex-1 p-4"
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
        }
      >
        {/* 分运费设置 */}
        <PermissionGate anyOf={['shop:fees:view']}>
          <Card className="mb-4 p-4">
            <View className="flex-row items-center justify-between mb-2">
              <View className="flex-1">
                <Text className="text-base font-semibold mb-1">分运费</Text>
                <Text className="text-sm text-gray-500">
                  {shareFee ? '已开启：店铺将收取分运费用' : '已关闭：店铺不收取分运费用'}
                </Text>
              </View>
              <PermissionGate anyOf={['shop:fees:update']}>
                <Switch
                  value={shareFee}
                  onValueChange={handleToggleShareFee}
                  disabled={!houseGid || updating}
                />
              </PermissionGate>
            </View>
          </Card>
        </PermissionGate>

        {/* 推送配额设置 - 已废弃，仅保留查看 */}
        {pushCredit && pushCredit !== '0' && (
          <PermissionGate anyOf={['shop:fees:view']}>
            <Card className="mb-4 p-4 bg-gray-50">
              <Text className="text-base font-semibold mb-2 text-gray-600">
                推送配额（历史数据）
              </Text>
              <Text className="text-sm text-gray-500 mb-2">
                当前配额：{pushCredit} 条
              </Text>
              <View className="bg-yellow-50 p-3 rounded border border-yellow-200">
                <Text className="text-xs text-yellow-700">
                  此功能已废弃：新版本不再使用机器人推送，此配额仅供查看历史数据
                </Text>
              </View>
            </Card>
          </PermissionGate>
        )}

        {/* 运费规则列表 */}
        <PermissionGate anyOf={['shop:fees:view']}>
          <Card className="mb-4 p-4">
            <Text className="text-base font-semibold mb-3">运费规则</Text>
            
            {gameFees.length > 0 ? (
              <View className="space-y-2">
                {gameFees.map((fee, index) => (
                  <View key={index} className="flex-row items-center justify-between p-3 bg-gray-50 rounded">
                    <View className="flex-1">
                      <Text className="text-sm font-medium">
                        {(fee.kind === '0' || !fee.kind) && (fee.base === 0 || !fee.base)
                          ? `通用规则: ${fee.threshold}分/${fee.fee}分`
                          : `${getGameKindName(Number(fee.kind))} ${fee.base}底分: ${fee.threshold}分/${fee.fee}分`}
                      </Text>
                      <Text className="text-xs text-gray-500 mt-1">
                        门槛: {fee.threshold}分 → 运费: {fee.fee}分
                      </Text>
                    </View>
                    <PermissionGate anyOf={['shop:fees:update']}>
                      <View className="flex-row">
                        <Button
                          variant="ghost"
                          size="icon"
                          onPress={() => handleEditFee(index)}
                          disabled={updating}
                        >
                          <Icon as={Pencil} size={16} className="text-blue-500" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onPress={() => handleDeleteFee(index)}
                          disabled={updating}
                        >
                          <Icon as={Trash2} size={16} className="text-red-500" />
                        </Button>
                      </View>
                    </PermissionGate>
                  </View>
                ))}
              </View>
            ) : (
              <View className="items-center py-6">
                <Text className="text-gray-400 text-sm">暂无运费规则</Text>
              </View>
            )}
          </Card>
        </PermissionGate>

        {/* 添加/编辑运费规则 */}
        <PermissionGate anyOf={['shop:fees:update']}>
          <Card className="mb-4 p-4">
            <View className="flex-row items-center justify-between mb-3">
              <Text className="text-base font-semibold">
                {editingIndex !== null ? `编辑规则 #${editingIndex + 1}` : '添加运费规则'}
              </Text>
              {editingIndex !== null && (
                <Button variant="ghost" size="sm" onPress={handleCancelEdit}>
                  <Icon as={X} size={16} className="text-gray-500" />
                  <Text className="text-xs text-gray-500 ml-1">取消</Text>
                </Button>
              )}
            </View>
            
            <View className="mb-3">
              <Text className="text-sm text-gray-600 mb-1">门槛金额（分）*</Text>
              <Input
                placeholder="如：50"
                keyboardType="numeric"
                value={newThreshold}
                onChangeText={setNewThreshold}
              />
            </View>

            <View className="mb-3">
              <Text className="text-sm text-gray-600 mb-1">运费金额（分）*</Text>
              <Input
                placeholder="如：800"
                keyboardType="numeric"
                value={newFee}
                onChangeText={setNewFee}
              />
            </View>

            <View className="mb-3" style={{ zIndex: 50 }}>
              <Text className="text-sm text-gray-600 mb-1">游戏类型 *</Text>
              <Select
                value={newGameKind && newGameKind !== '0' ? { 
                  value: newGameKind, 
                  label: maps.game_kinds.get(Number(newGameKind)) || getGameKindName(newGameKind)
                } : undefined}
                onValueChange={(option) => setNewGameKind(option?.value)}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="请选择游戏类型" />
                </SelectTrigger>
                <SelectContent portalHost="shop-layout-portal" style={{ zIndex: 9999 }}>
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
              <Text className="text-sm text-gray-600 mb-1">底分（可选）</Text>
              <Input
                placeholder="留空表示所有底分"
                keyboardType="numeric"
                value={newBaseScore}
                onChangeText={setNewBaseScore}
              />
            </View>

            <Button
              onPress={handleSaveFee}
              disabled={!houseGid || !newThreshold || !newFee || !newGameKind || newGameKind === '0' || updating}
            >
              <View className="flex-row items-center gap-2">
                {editingIndex === null && <Icon as={Plus} size={16} className="text-white" />}
                <Text className="text-white">{editingIndex !== null ? '保存修改' : '添加规则'}</Text>
              </View>
            </Button>

            <View className="mt-3 bg-blue-50 p-3 rounded">
              <Text className="text-xs text-blue-700">
                说明：
              </Text>
              <Text className="text-xs text-blue-700 mt-1">
                • 门槛、运费和游戏类型必填，底分可选
              </Text>
              <Text className="text-xs text-blue-700">
                • 底分留空表示适用该游戏所有底分
              </Text>
              <Text className="text-xs text-blue-700">
                • 规则示例：门槛50分，运费800分，表示达到50分收取800分运费
              </Text>
            </View>
          </Card>
        </PermissionGate>

        {/* 空状态提示 */}
        {!houseGid && (
          <View className="items-center justify-center py-12">
            <Text className="text-gray-400">请输入店铺号查询设置</Text>
          </View>
        )}
      </ScrollView>
    </View>
  );
}
