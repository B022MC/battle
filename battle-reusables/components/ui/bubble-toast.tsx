import React, { useEffect, useState } from 'react';
import { View, Animated, Pressable } from 'react-native';
import { Text } from '@/components/ui/text';
import { Icon } from '@/components/ui/icon';
import { CheckCircle2, XCircle, AlertCircle, Info, X } from 'lucide-react-native';
import { cn } from '@/lib/utils';

export type BubbleToastType = 'success' | 'error' | 'warning' | 'info';

export interface BubbleToastProps {
  type: BubbleToastType;
  title: string;
  description?: string;
  duration?: number;
  onClose?: () => void;
}

const iconMap = {
  success: CheckCircle2,
  error: XCircle,
  warning: AlertCircle,
  info: Info,
};

const colorMap = {
  success: {
    bg: 'bg-green-50 dark:bg-green-950/30',
    border: 'border-green-200 dark:border-green-800',
    icon: 'text-green-600 dark:text-green-400',
    title: 'text-green-900 dark:text-green-100',
    description: 'text-green-700 dark:text-green-300',
  },
  error: {
    bg: 'bg-red-50 dark:bg-red-950/30',
    border: 'border-red-200 dark:border-red-800',
    icon: 'text-red-600 dark:text-red-400',
    title: 'text-red-900 dark:text-red-100',
    description: 'text-red-700 dark:text-red-300',
  },
  warning: {
    bg: 'bg-yellow-50 dark:bg-yellow-950/30',
    border: 'border-yellow-200 dark:border-yellow-800',
    icon: 'text-yellow-600 dark:text-yellow-400',
    title: 'text-yellow-900 dark:text-yellow-100',
    description: 'text-yellow-700 dark:text-yellow-300',
  },
  info: {
    bg: 'bg-blue-50 dark:bg-blue-950/30',
    border: 'border-blue-200 dark:border-blue-800',
    icon: 'text-blue-600 dark:text-blue-400',
    title: 'text-blue-900 dark:text-blue-100',
    description: 'text-blue-700 dark:text-blue-300',
  },
};

export function BubbleToast({ type, title, description, duration = 3000, onClose }: BubbleToastProps) {
  const [opacity] = useState(new Animated.Value(0));
  const [translateY] = useState(new Animated.Value(-20));
  const [scale] = useState(new Animated.Value(0.9));

  useEffect(() => {
    // 入场动画
    Animated.parallel([
      Animated.timing(opacity, {
        toValue: 1,
        duration: 300,
        useNativeDriver: true,
      }),
      Animated.spring(translateY, {
        toValue: 0,
        tension: 100,
        friction: 8,
        useNativeDriver: true,
      }),
      Animated.spring(scale, {
        toValue: 1,
        tension: 100,
        friction: 8,
        useNativeDriver: true,
      }),
    ]).start();

    // 自动关闭
    if (duration > 0) {
      const timer = setTimeout(() => {
        handleClose();
      }, duration);
      return () => clearTimeout(timer);
    }
  }, []);

  const handleClose = () => {
    // 出场动画
    Animated.parallel([
      Animated.timing(opacity, {
        toValue: 0,
        duration: 200,
        useNativeDriver: true,
      }),
      Animated.timing(translateY, {
        toValue: -20,
        duration: 200,
        useNativeDriver: true,
      }),
      Animated.timing(scale, {
        toValue: 0.9,
        duration: 200,
        useNativeDriver: true,
      }),
    ]).start(() => {
      onClose?.();
    });
  };

  const IconComponent = iconMap[type];
  const colors = colorMap[type];

  return (
    <Animated.View
      style={{
        opacity,
        transform: [{ translateY }, { scale }],
      }}
      className="px-4 py-2"
    >
      <View
        className={cn(
          'flex-row items-start gap-3 rounded-2xl border-2 p-4 shadow-lg',
          colors.bg,
          colors.border
        )}
        style={{
          shadowColor: '#000',
          shadowOffset: { width: 0, height: 4 },
          shadowOpacity: 0.15,
          shadowRadius: 12,
          elevation: 8,
        }}
      >
        {/* 图标 */}
        <View className="pt-0.5">
          <Icon as={IconComponent} size={24} className={colors.icon} />
        </View>

        {/* 内容 */}
        <View className="flex-1">
          <Text className={cn('text-base font-semibold mb-0.5', colors.title)}>
            {title}
          </Text>
          {description && (
            <Text className={cn('text-sm', colors.description)}>
              {description}
            </Text>
          )}
        </View>

        {/* 关闭按钮 */}
        <Pressable
          onPress={handleClose}
          className="pt-0.5 active:opacity-50"
          hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
        >
          <Icon as={X} size={20} className={cn('opacity-50', colors.icon)} />
        </Pressable>
      </View>
    </Animated.View>
  );
}

