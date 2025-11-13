import { Tabs, useFocusEffect, useRouter } from 'expo-router';
import {
  CircleDashedIcon,
  HomeIcon,
  UsersIcon,
  WalletIcon,
  UserCircleIcon,
  StoreIcon
} from 'lucide-react-native';
import { TabIcon, TabLabel } from '@/components/shared/tab-item';
import { useAuthStore } from '@/hooks/use-auth-store';
import { usePermission } from '@/hooks/use-permission';

export default function TabLayout() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const { hasAny } = usePermission();
  const canViewTables = hasAny(['shop:table:view']);
  const canViewMembers = hasAny(['shop:member:view']);

  useFocusEffect(() => {
    if (!isAuthenticated) router.replace('/auth');
  });

  return (
    <Tabs
      screenOptions={{
        headerShown: true,
        headerTitleAlign: 'center',
        headerShadowVisible: false,
        headerStyle: { borderBottomWidth: 0 },
      }}>
      {/* 基于权限动态控制 Tab 显示 */}
      <Tabs.Screen
        name="index"
        options={{
          title: '首页',
          tabBarLabel: ({ focused }) => <TabLabel focused={focused} label="首页" />,
          tabBarIcon: ({ focused }) => <TabIcon focused={focused} icon={HomeIcon} />,
        }}
      />
      <Tabs.Screen
        name="tables"
        options={{
          title: '桌台',
          tabBarLabel: ({ focused }) => <TabLabel focused={focused} label="桌台" />,
          tabBarIcon: ({ focused }) => <TabIcon focused={focused} icon={CircleDashedIcon} />,
          href: canViewTables ? undefined : null,
        }}
      />
      <Tabs.Screen
        name="members"
        options={{
          title: '成员',
          tabBarLabel: ({ focused }) => <TabLabel focused={focused} label="成员" />,
          tabBarIcon: ({ focused }) => <TabIcon focused={focused} icon={UsersIcon} />,
          href: canViewMembers ? undefined : null,
        }}
      />
      <Tabs.Screen
        name="funds"
        options={{
          title: '资金',
          tabBarLabel: ({ focused }) => <TabLabel focused={focused} label="资金" />,
          tabBarIcon: ({ focused }) => <TabIcon focused={focused} icon={WalletIcon} />,
        }}
      />
      <Tabs.Screen
        name="shop"
        options={{
          title: '店铺',
          tabBarLabel: ({ focused }) => <TabLabel focused={focused} label="店铺" />,
          tabBarIcon: ({ focused }) => <TabIcon focused={focused} icon={StoreIcon} />,
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: '我的',
          tabBarLabel: ({ focused }) => <TabLabel focused={focused} label="我的" />,
          tabBarIcon: ({ focused }) => <TabIcon focused={focused} icon={UserCircleIcon} />,
        }}
      />
    </Tabs>
  );
}
