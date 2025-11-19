/**
 * RouteGuard 路由守卫组件
 * 用于控制页面级别的访问权限
 */
import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { usePermission } from '@/hooks/use-permission';

type Props = {
  /** 所需的权限码（满足任一即可） */
  anyOf?: string[];
  /** 所需的权限码（必须全部满足） */
  allOf?: string[];
  /** 无权限时显示的内容 */
  fallback?: React.ReactNode;
  /** 子组件 */
  children: React.ReactNode;
};

export function RouteGuard({ anyOf, allOf, fallback, children }: Props) {
  const { isSuperAdmin, hasAny, hasAll } = usePermission();
  
  // 超级管理员直接放行
  if (isSuperAdmin) {
    return <>{children}</>;
  }
  
  // 检查权限
  const hasAllPermissions = allOf && allOf.length > 0 ? hasAll(allOf) : true;
  const hasAnyPermission = anyOf && anyOf.length > 0 ? hasAny(anyOf) : true;
  
  const hasPermission = hasAllPermissions && hasAnyPermission;
  
  if (!hasPermission) {
    return fallback ? (
      <>{fallback}</>
    ) : (
      <View className="flex-1 items-center justify-center p-6">
        <Text className="text-xl font-bold text-destructive mb-2">
          访问受限
        </Text>
        <Text className="text-center text-muted-foreground">
          您没有权限访问此页面
        </Text>
      </View>
    );
  }
  
  return <>{children}</>;
}

