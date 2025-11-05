import { Stack, useFocusEffect, useRouter } from 'expo-router';
import { useAuthStore } from '@/hooks/use-auth-store';

export default function ShopLayout() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();

  useFocusEffect(() => {
    if (!isAuthenticated) router.replace('/auth');
  });

  return (
    <Stack
      screenOptions={{
        headerShown: true,
        headerTitleAlign: 'center',
        headerShadowVisible: false,
      }}>
      <Stack.Screen name="account" options={{ title: '游戏账号' }} />
      <Stack.Screen name="session" options={{ title: '会话管理' }} />
      <Stack.Screen name="applications" options={{ title: '申请管理' }} />
      <Stack.Screen name="admins" options={{ title: '管理员' }} />
      <Stack.Screen name="rooms" options={{ title: '中控账号' }} />
      <Stack.Screen name="fees" options={{ title: '费用设置' }} />
      <Stack.Screen name="balances" options={{ title: '余额筛查' }} />
      <Stack.Screen name="members" options={{ title: '成员管理' }} />
    </Stack>
  );
}
