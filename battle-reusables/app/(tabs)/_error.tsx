import React from 'react';
import { View, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { router, useGlobalSearchParams } from 'expo-router';
import { useAuthStore } from '@/hooks/use-auth-store';

export default function TabsErrorBoundary() {
  const { isAuthenticated } = useAuthStore();
  const params = useGlobalSearchParams();

  // 如果未登录，提示并跳转到登录页
  if (!isAuthenticated) {
    return (
      <ScrollView className="flex-1 bg-background p-4">
        <View className="flex-1 items-center justify-center gap-4">
          <Text className="text-xl font-bold text-primary">需要登录</Text>
          <Text className="text-center text-muted-foreground">
            请先登录后再访问此页面
          </Text>
          <View className="gap-2">
            <Button onPress={() => router.replace('/auth')}>
              <Text>去登录</Text>
            </Button>
          </View>
        </View>
      </ScrollView>
    );
  }

  return (
    <ScrollView className="flex-1 bg-background p-4">
      <View className="flex-1 items-center justify-center gap-4">
        <Text className="text-xl font-bold text-destructive">出错了</Text>
        <Text className="text-center text-muted-foreground">
          页面加载失败，请重试
        </Text>
        <Text className="text-xs text-muted-foreground/50">
          {JSON.stringify(params)}
        </Text>
        <View className="gap-2">
          <Button onPress={() => router.replace('/auth')}>
            <Text>返回登录</Text>
          </Button>
          <Button variant="outline" onPress={() => router.back()}>
            <Text>返回上一页</Text>
          </Button>
        </View>
      </View>
    </ScrollView>
  );
}
