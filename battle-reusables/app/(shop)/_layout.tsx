import { Stack, useFocusEffect, useRouter } from 'expo-router';
import { useAuthStore } from '@/hooks/use-auth-store';
import { Platform } from 'react-native';
import { HeaderBackButton } from '@/components/shared/header-back-button';

export default function ShopLayout() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();

  useFocusEffect(() => {
    if (!isAuthenticated) router.push('/auth');
  });

  return (
    <Stack
      screenOptions={{
        headerShown: true,
        headerTitleAlign: 'center',
        headerShadowVisible: false,
        // 使用自定义返回按钮（特别是为了 Web 平台）
        headerLeft: () => <HeaderBackButton />,
        // 确保 header 在最顶层，避免被遮罩层挡住
        headerStyle: {
          // @ts-ignore - Web only CSS
          zIndex: 1000,
        },
      }}>
      <Stack.Screen name="account" options={{ title: '游戏账号' }} />
      <Stack.Screen name="admins" options={{ title: '管理员' }} />
      <Stack.Screen name="rooms" options={{ title: '中控账号' }} />
      <Stack.Screen name="fees" options={{ title: '费用设置' }} />
      <Stack.Screen name="balances" options={{ title: '余额筛查' }} />
      <Stack.Screen name="members" options={{ title: '成员管理' }} />
      <Stack.Screen name="my-battles" options={{ title: '我的战绩' }} />
      <Stack.Screen name="my-balances" options={{ title: '我的余额' }} />
      <Stack.Screen name="group-battles" options={{ title: '圈子战绩' }} />
      <Stack.Screen name="group-balances" options={{ title: '圈子余额' }} />
      <Stack.Screen name="permissions" options={{ title: '权限管理' }} />
      <Stack.Screen name="roles" options={{ title: '角色管理' }} />
      <Stack.Screen name="menus" options={{ title: '菜单管理' }} />
      <Stack.Screen name="game-applications" options={{ title: '游戏申请' }} />
      <Stack.Screen name="shop-fees" options={{ title: '店铺费用' }} />
      <Stack.Screen name="room-credit" options={{ title: '额度设置' }} />
    </Stack>
  );
}
