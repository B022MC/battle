import React from 'react';
import { View } from 'react-native';
import { AccountView } from '@/components/(shop)/account/account-view';
import { usePermission } from '@/hooks/use-permission';
import { Text } from '@/components/ui/text';

export default function AccountScreen() {
  const { isSuperAdmin, hasAny } = usePermission();
  const isAdmin = hasAny([
    'shop:admin:assign', 'shop:admin:view',
    'shop:member:view', 'shop:table:view', 'shop:apply:view',
    'game:ctrl:view', 'game:ctrl:update', 'game:ctrl:create',
  ]);

  if (isSuperAdmin || isAdmin) {
    return (
      <View className="flex-1 items-center justify-center p-6">
        <Text className="text-muted-foreground">当前身份无需配置个人游戏账号</Text>
      </View>
    );
  }
  return <AccountView />;
}
