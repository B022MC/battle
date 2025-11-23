import React, { useState, useMemo, useEffect } from 'react';
import { View, Platform, Pressable } from 'react-native';
import { RouteGuard } from '@/components/auth/RouteGuard';
import { GameApplicationList } from '@/components/(shop)/game-applications/game-application-list';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Icon } from '@/components/ui/icon';
import { ChevronDown } from 'lucide-react-native';
import { useHouseSelector } from '@/hooks/use-house-selector';

/**
 * 游戏内申请管理页面
 * 
 * 功能：
 * - 查看游戏内玩家提交的申请（加入店铺、换圈等）
 * - 批准或拒绝申请
 * - 自动刷新申请列表
 * 
 * 权限要求：shop:applications:view
 */
function GameApplicationsContent() {
  const {
    houseGid,
    setHouseGid,
    isSuperAdmin,
    isStoreAdmin,
    houseOptions,
    loading: loadingHouse,
    isReady,
  } = useHouseSelector();
  
  const [confirmedHouseGid, setConfirmedHouseGid] = useState<number | null>(null);
  const [open, setOpen] = useState(false);

  // 店铺管理员自动确认店铺
  useEffect(() => {
    if (isStoreAdmin && houseGid && isReady) {
      const gid = Number(houseGid);
      if (!isNaN(gid) && gid > 0) {
        setConfirmedHouseGid(gid);
      }
    }
  }, [isStoreAdmin, houseGid, isReady]);

  const handleConfirm = () => {
    const gid = Number(houseGid);
    if (!isNaN(gid) && gid > 0) {
      setConfirmedHouseGid(gid);
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
    <View className="flex-1 p-4">
      {!confirmedHouseGid ? (
        <View className="items-center justify-center flex-1">
          <View className="w-full max-w-sm p-6 bg-card rounded-lg">
            <Text className="text-xl font-bold mb-4">
              {isStoreAdmin ? '店铺信息' : '选择店铺'}
            </Text>
            
            {/* 店铺管理员：显示当前店铺 */}
            {isStoreAdmin && (
              <View className="mb-4 p-3 bg-blue-50 rounded border border-blue-200">
                <Text className="text-sm text-blue-700">
                  当前店铺：{houseGid || '加载中...'}
                </Text>
                <Text className="text-xs text-blue-600 mt-1">
                  正在加载申请列表...
                </Text>
              </View>
            )}
            
            {/* 超级管理员：下拉选择店铺 */}
            {isSuperAdmin && (
              <>
                <Text className="text-muted-foreground mb-4">
                  请选择要管理的店铺
                </Text>
                <View className="relative mb-4">
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
                  onPress={handleConfirm}
                  disabled={!houseGid}
                  className="w-full"
                >
                  <Text>确认</Text>
                </Button>
              </>
            )}
          </View>
        </View>
      ) : (
        <GameApplicationList houseGid={confirmedHouseGid} />
      )}
    </View>
  );
}

export default function GameApplicationsPage() {
  return (
    <RouteGuard anyOf={['shop:applications:view']}>
      <GameApplicationsContent />
    </RouteGuard>
  );
}
