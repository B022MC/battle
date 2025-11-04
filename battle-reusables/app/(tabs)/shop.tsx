import React from 'react';
import { View, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { router } from 'expo-router';
import { PermissionGate } from '@/components/auth/PermissionGate';
 

export default function ShopHubScreen() {

  return (
    <ScrollView className="flex-1 p-4">
      <View className="gap-4">
        <Text className="text-xl font-bold">店铺</Text>
        <View className="gap-2">
        <PermissionGate anyOf={["shop:table:view"]}>
          <Button onPress={() => router.push('/(shop)/rooms')}><Text>房间管理</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:member:view"]}>
          <Button onPress={() => router.push('/(shop)/members')}><Text>成员管理</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:apply:view"]}>
          <Button onPress={() => router.push('/(shop)/applications')}><Text>提交申请</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:apply:view"]}>
          <Button onPress={() => router.push('/(shop)/applications-list')}><Text>申请列表</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["game:ctrl:view","game:ctrl:update","game:ctrl:create"]}>
          <Button onPress={() => router.push('/(shop)/session')}><Text>会话控制</Text></Button>
        </PermissionGate>
      </View>
        
      </View>
    </ScrollView>
  );
}
