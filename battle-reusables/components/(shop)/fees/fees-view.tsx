import React, { useState, useEffect } from 'react';
import { View, ScrollView, ActivityIndicator, RefreshControl } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Switch } from 'react-native';
import { shopsFeesGet, shopsShareFeeSet } from '@/services/shops/fees';
import { showToast } from '@/utils/toast';
import { PermissionGate } from '@/components/auth/PermissionGate';

export function FeesView() {
  const [houseGid, setHouseGid] = useState('');
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  
  // 设置数据
  const [shareFee, setShareFee] = useState(false);
  const [pushCredit, setPushCredit] = useState('');
  const [feesJson, setFeesJson] = useState('');
  
  // 操作状态
  const [updating, setUpdating] = useState(false);

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
        setFeesJson(res.data.fees_json || '');
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


  return (
    <View className="flex-1 bg-gray-50">
      {/* 头部查询区 */}
      <View className="bg-white p-4 border-b border-gray-200">
        <Text className="text-lg font-semibold mb-3">店铺费用设置</Text>
        <View className="flex-row gap-2">
          <Input
            className="flex-1"
            placeholder="请输入店铺号"
            keyboardType="numeric"
            value={houseGid}
            onChangeText={setHouseGid}
          />
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
      </View>

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
                  ⚠️ 此功能已废弃：新版本不再使用机器人推送，此配额仅供查看历史数据
                </Text>
              </View>
            </Card>
          </PermissionGate>
        )}

        {/* 费用规则信息（只读显示） */}
        {feesJson && (
          <PermissionGate anyOf={['shop:fees:view']}>
            <Card className="mb-4 p-4">
              <Text className="text-base font-semibold mb-3">费用规则配置</Text>
              <View className="bg-gray-100 p-3 rounded">
                <Text className="text-xs font-mono text-gray-700">
                  {feesJson}
                </Text>
              </View>
              <Text className="text-xs text-gray-400 mt-2">
                提示：费用规则配置较复杂，请联系技术人员进行修改
              </Text>
            </Card>
          </PermissionGate>
        )}

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
