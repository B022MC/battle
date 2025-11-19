import { View, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { useAuthStore } from '@/hooks/use-auth-store';

export default function Screen() {
  const { isAuthenticated, user } = useAuthStore();

  console.log('[Index Screen] isAuthenticated:', isAuthenticated);

  return (
    <ScrollView className="flex-1 bg-secondary">
      <View className="gap-4 p-4">
        <View className="rounded-lg border border-border bg-card p-4">
          <Text className="text-lg font-semibold mb-2">欢迎回来！</Text>
          <Text className="text-muted-foreground">用户: {user?.username || user?.nick_name || '未知'}</Text>
          <Text className="text-muted-foreground mt-2">认证状态: {isAuthenticated ? '已登录' : '未登录'}</Text>
        </View>
        
        <View className="rounded-lg border border-border bg-card p-4">
          <Text className="text-sm text-muted-foreground">
            首页功能正在加载中...
          </Text>
        </View>
      </View>
    </ScrollView>
  );
}
