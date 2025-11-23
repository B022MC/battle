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
        <Button onPress={() => router.push('/(shop)/account')}><Text>游戏账号</Text></Button>
        <PermissionGate anyOf={["shop:admin:view"]}>
          <Button onPress={() => router.push('/(shop)/admins')}><Text>管理员</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["game:ctrl:view"]}>
          <Button onPress={() => router.push('/(shop)/rooms')}><Text>中控账号</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:fees:view"]}>
          <Button onPress={() => router.push('/(shop)/fees')}><Text>费用设置</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["fund:wallet:view"]}>
          <Button onPress={() => router.push('/(shop)/balances')}><Text>余额筛查</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:member:view"]}>
          <Button onPress={() => router.push('/(shop)/members')}><Text>成员管理</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["game:room_credit:view"]}>
          <Button onPress={() => router.push('/(shop)/room-credits')}><Text>房间额度</Text></Button>
        </PermissionGate>
        <Button onPress={() => router.push('/(shop)/game-applications')}><Text>游戏申请</Text></Button>

        {/* 战绩查询功能 */}
        <Text className="text-lg font-semibold mt-4">战绩查询</Text>
        <Button onPress={() => router.push('/(shop)/my-battles')}><Text>我的战绩</Text></Button>
        <Button onPress={() => router.push('/(shop)/my-balances')}><Text>我的余额</Text></Button>

        <PermissionGate anyOf={["shop:member:view"]}>
          <Button onPress={() => router.push('/(shop)/group-battles')}><Text>圈子战绩(管理员)</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["shop:member:view"]}>
          <Button onPress={() => router.push('/(shop)/group-balances')}><Text>圈子余额(管理员)</Text></Button>
        </PermissionGate>

        {/* 系统管理功能 */}
        <Text className="text-lg font-semibold mt-4">系统管理</Text>
        <PermissionGate anyOf={["permission:view"]}>
          <Button onPress={() => router.push('/(shop)/permissions')}><Text>权限管理</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["role:view"]}>
          <Button onPress={() => router.push('/(shop)/roles')}><Text>角色管理</Text></Button>
        </PermissionGate>
        <PermissionGate anyOf={["menu:view"]}>
          <Button onPress={() => router.push('/(shop)/menus')}><Text>菜单管理</Text></Button>
        </PermissionGate>
      </View>

      </View>
    </ScrollView>
  );
}
