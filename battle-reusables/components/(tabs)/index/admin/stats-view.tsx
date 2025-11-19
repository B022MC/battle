import React from 'react';
import { View, ScrollView } from 'react-native';
import { StatsSearch } from './stats-search';
import { SessionStats } from './session-stats';
import { WalletStats } from './wallet-stats';
import { LedgerStats } from './ledger-stats';
import { statsByPath, statsActiveByHouse } from '@/services/stats';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardContent, InfoCardRow } from '@/components/shared/info-card';
import { useRequest } from '@/hooks/use-request';
import { Text } from '@/components/ui/text';
import { useAuthStore } from '@/hooks/use-auth-store';

export const StatsView = () => {
  const { isAuthenticated } = useAuthStore();
  const { data, loading, run } = useRequest(statsByPath, { manual: true });
  const { data: activeSessions, loading: activeLoading } = useRequest(statsActiveByHouse, {
    manual: !isAuthenticated, // 未登录时不自动请求
  });

  // 未登录时显示提示信息
  if (!isAuthenticated) {
    return (
      <ScrollView className="flex-1 bg-secondary">
        <View className="gap-4 p-4">
          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>提示</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <Text className="text-muted-foreground">请先登录后查看统计数据</Text>
            </InfoCardContent>
          </InfoCard>
        </View>
      </ScrollView>
    );
  }

  return (
    <View className="flex-1">
      <StatsSearch onSubmit={run} submitButtonProps={{ loading }} />
      <ScrollView className="flex-1 bg-secondary">
        <View className="gap-4 p-4">
          {/* 所有店铺会话状态 */}
          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>店铺会话状态</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-2">
                {activeLoading ? (
                  <Text className="text-muted-foreground">加载中...</Text>
                ) : Array.isArray(activeSessions) && activeSessions.length > 0 ? (
                  activeSessions.map((it) => (
                    <InfoCardRow key={String(it!.house_gid)} label={`店铺 ${it!.house_gid}`} value={`在线会话：${it!.active ?? 0}`} />
                  ))
                ) : (
                  <Text className="text-muted-foreground">暂无在线会话</Text>
                )}
              </View>
            </InfoCardContent>
          </InfoCard>
          {/* 统计主体（按接口返回渲染） */}
          {loading ? (
            <InfoCard>
              <InfoCardHeader><InfoCardTitle>统计</InfoCardTitle></InfoCardHeader>
              <InfoCardContent>
                <Text className="text-muted-foreground">加载中...</Text>
              </InfoCardContent>
            </InfoCard>
          ) : data ? (
            <>
              <SessionStats data={data?.session} />
              <WalletStats data={data?.wallet} />
              <LedgerStats data={data?.ledger} />
            </>
          ) : (
            <InfoCard>
              <InfoCardHeader><InfoCardTitle>统计</InfoCardTitle></InfoCardHeader>
              <InfoCardContent>
                <Text className="text-muted-foreground">请选择店铺并点击统计</Text>
              </InfoCardContent>
            </InfoCard>
          )}
        </View>
      </ScrollView>
    </View>
  );
};
