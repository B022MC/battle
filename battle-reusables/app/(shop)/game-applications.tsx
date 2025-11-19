import React, { useState } from 'react';
import { View } from 'react-native';
import { RouteGuard } from '@/components/auth/RouteGuard';
import { GameApplicationList } from '@/components/(shop)/game-applications/game-application-list';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

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
  const [houseGid, setHouseGid] = useState('');
  const [confirmedHouseGid, setConfirmedHouseGid] = useState<number | null>(null);

  const handleConfirm = () => {
    const gid = Number(houseGid);
    if (!isNaN(gid) && gid > 0) {
      setConfirmedHouseGid(gid);
    }
  };

  return (
    <View className="flex-1 p-4">
      {!confirmedHouseGid ? (
        <View className="items-center justify-center flex-1">
          <View className="w-full max-w-sm p-6 bg-card rounded-lg">
            <Text className="text-xl font-bold mb-4">选择店铺</Text>
            <Text className="text-muted-foreground mb-4">
              请输入要管理的店铺号
            </Text>
            <Input
              placeholder="输入店铺号"
              value={houseGid}
              onChangeText={setHouseGid}
              keyboardType="numeric"
              className="mb-4"
            />
            <Button 
              onPress={handleConfirm}
              disabled={!houseGid}
              className="w-full"
            >
              <Text>确认</Text>
            </Button>
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
