import React, { useState, useEffect } from 'react';
import { View, FlatList, RefreshControl } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Loading } from '@/components/shared/loading';
import { 
  listGameApplications, 
  approveGameApplication, 
  rejectGameApplication 
} from '@/services/shops/game-applications';
import type { GameApplication } from '@/services/shops/game-applications/typing';
import { showToast } from '@/utils/toast';

type GameApplicationListProps = {
  houseGid: number;
};

export const GameApplicationList = ({ houseGid }: GameApplicationListProps) => {
  const [applications, setApplications] = useState<GameApplication[]>([]);
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [processingId, setProcessingId] = useState<number | null>(null);

  // 加载申请列表
  const loadApplications = async (isRefresh = false) => {
    if (isRefresh) {
      setRefreshing(true);
    } else {
      setLoading(true);
    }

    try {
      const res = await listGameApplications({ house_gid: houseGid });
      if (res.code === 0) {
        setApplications(res.data || []);
      } else {
        showToast(res.msg || '加载失败', 'error');
      }
    } catch (error) {
      showToast('加载失败', 'error');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  // 通过申请
  const handleApprove = async (messageId: number) => {
    setProcessingId(messageId);
    try {
      const res = await approveGameApplication({
        house_gid: houseGid,
        message_id: messageId,
      });
      if (res.code === 0) {
        showToast('已通过申请', 'success');
        // 立即刷新列表
        loadApplications(true);
      } else {
        showToast(res.msg || '操作失败', 'error');
      }
    } catch (error) {
      showToast('操作失败', 'error');
    } finally {
      setProcessingId(null);
    }
  };

  // 拒绝申请
  const handleReject = async (messageId: number) => {
    setProcessingId(messageId);
    try {
      const res = await rejectGameApplication({
        house_gid: houseGid,
        message_id: messageId,
      });
      if (res.code === 0) {
        showToast('已拒绝申请', 'success');
        // 立即刷新列表
        loadApplications(true);
      } else {
        showToast(res.msg || '操作失败', 'error');
      }
    } catch (error) {
      showToast('操作失败', 'error');
    } finally {
      setProcessingId(null);
    }
  };

  // 格式化时间戳
  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp * 1000);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  };

  // 初始化和定时刷新
  useEffect(() => {
    loadApplications();

    // 定时刷新（10秒一次）
    const timer = setInterval(() => {
      loadApplications(true);
    }, 10000);

    return () => clearInterval(timer);
  }, [houseGid]);

  if (loading) {
    return <Loading text="加载中..." />;
  }

  if (applications.length === 0) {
    return (
      <View className="min-h-32 flex-row items-center justify-center">
        <Text className="text-muted-foreground">暂无待处理申请</Text>
      </View>
    );
  }

  return (
    <FlatList
      data={applications}
      keyExtractor={(item) => item.message_id.toString()}
      renderItem={({ item }) => (
        <Card className="mb-3 p-4">
          <View className="flex-row items-center justify-between mb-3">
            <View className="flex-1">
              <Text className="text-lg font-semibold mb-1">
                {item.applier_gname}
              </Text>
              <Text className="text-sm text-muted-foreground">
                游戏ID: {item.applier_gid}
              </Text>
              <Text className="text-xs text-muted-foreground mt-1">
                {formatTime(item.applied_at)}
              </Text>
            </View>
            <View className="bg-yellow-100 px-2 py-1 rounded">
              <Text className="text-yellow-800 text-xs">待处理</Text>
            </View>
          </View>

          <View className="flex-row gap-2">
            <Button
              variant="outline"
              size="sm"
              onPress={() => handleReject(item.message_id)}
              disabled={processingId === item.message_id}
              className="flex-1"
            >
              {processingId === item.message_id ? '处理中...' : '拒绝'}
            </Button>
            <Button
              size="sm"
              onPress={() => handleApprove(item.message_id)}
              disabled={processingId === item.message_id}
              className="flex-1"
            >
              {processingId === item.message_id ? '处理中...' : '通过'}
            </Button>
          </View>
        </Card>
      )}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={() => loadApplications(true)} />
      }
    />
  );
};
