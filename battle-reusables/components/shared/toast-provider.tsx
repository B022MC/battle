import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { View, Pressable } from 'react-native';
import { Text } from '@/components/ui/text';
import { Portal } from '@rn-primitives/portal';
import { NativeOnlyAnimatedView } from '@/components/ui/native-only-animated-view';
import { FadeIn, FadeOut, SlideInDown, SlideOutUp } from 'react-native-reanimated';
import { bindToast, unbindToast, type ToastOptions, type ConfirmOptions, type ToastType } from '@/utils/toast';
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';
import { cn } from '@/lib/utils';
import type { LucideIcon } from 'lucide-react-native';

/**
 * Toast 上下文值
 */
type ToastContextValue = {
  show: (options: ToastOptions) => void;
  success: (title: string, description?: string) => void;
  error: (title: string, description?: string) => void;
  warning: (title: string, description?: string) => void;
  info: (title: string, description?: string) => void;
  confirm: (options: ConfirmOptions) => void;
  close: () => void;
};

const ToastContext = createContext<ToastContextValue | null>(null);

/**
 * Toast 内部数据结构
 */
type ToastData = {
  type: 'toast' | 'confirm';
  toastOptions?: ToastOptions;
  confirmOptions?: ConfirmOptions;
};

/**
 * 根据类型获取默认图标
 */
function getDefaultIcon(type: ToastType): LucideIcon {
  switch (type) {
    case 'success':
      return CheckCircle;
    case 'error':
      return XCircle;
    case 'warning':
      return AlertTriangle;
    case 'info':
    default:
      return Info;
  }
}

/**
 * 根据类型获取样式类名
 */
function getTypeClassName(type: ToastType): string {
  switch (type) {
    case 'success':
      return 'border-green-500 bg-green-50 dark:bg-green-950';
    case 'error':
      return 'border-red-500 bg-red-50 dark:bg-red-950';
    case 'warning':
      return 'border-yellow-500 bg-yellow-50 dark:bg-yellow-950';
    case 'info':
    default:
      return 'border-blue-500 bg-blue-50 dark:bg-blue-950';
  }
}

/**
 * 根据类型获取图标颜色类名
 */
function getIconClassName(type: ToastType): string {
  switch (type) {
    case 'success':
      return 'text-green-600 dark:text-green-400';
    case 'error':
      return 'text-red-600 dark:text-red-400';
    case 'warning':
      return 'text-yellow-600 dark:text-yellow-400';
    case 'info':
    default:
      return 'text-blue-600 dark:text-blue-400';
  }
}

/**
 * Toast 提供者组件
 */
export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [current, setCurrent] = useState<ToastData | null>(null);
  const timerRef = React.useRef<NodeJS.Timeout>();

  const close = useCallback(() => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = undefined;
    }
    setCurrent(null);
  }, []);

  const show = useCallback((options: ToastOptions) => {
    // 清除之前的定时器
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = undefined;
    }
    
    setCurrent({ type: 'toast', toastOptions: options });
    
    // 如果 duration 为 0，不自动关闭
    const duration = options.duration ?? 3000;
    if (duration > 0) {
      timerRef.current = setTimeout(() => {
        setCurrent(null);
        timerRef.current = undefined;
      }, duration);
    }
  }, []);

  const success = useCallback((title: string, description?: string) => {
    show({ title, description, type: 'success', duration: 3000 });
  }, [show]);

  const error = useCallback((title: string, description?: string) => {
    show({ title, description, type: 'error', duration: 4000 });
  }, [show]);

  const warning = useCallback((title: string, description?: string) => {
    show({ title, description, type: 'warning', duration: 3500 });
  }, [show]);

  const info = useCallback((title: string, description?: string) => {
    show({ title, description, type: 'info', duration: 3000 });
  }, [show]);

  const confirm = useCallback((options: ConfirmOptions) => {
    // 清除之前的定时器
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = undefined;
    }
    setCurrent({ type: 'confirm', confirmOptions: options });
  }, []);

  // 绑定全局 toast API
  useEffect(() => {
    bindToast({ show, success, error, warning, info, confirm, close });
    return () => unbindToast();
  }, [show, success, error, warning, info, confirm, close]);

  // 清理定时器
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  return (
    <ToastContext.Provider value={{ show, success, error, warning, info, confirm, close }}>
      {children}
      {current && (
        <Portal name="toast-portal">
          {current.type === 'toast' && current.toastOptions && (
            <ToastContent options={current.toastOptions} onClose={close} />
          )}
          {current.type === 'confirm' && current.confirmOptions && (
            <ConfirmDialog options={current.confirmOptions} onClose={close} />
          )}
        </Portal>
      )}
    </ToastContext.Provider>
  );
}

/**
 * Toast 内容组件
 */
function ToastContent({ options, onClose }: { options: ToastOptions; onClose: () => void }) {
  const type = options.type ?? 'info';
  const icon = options.icon ?? getDefaultIcon(type);
  const position = options.position ?? 'top';

  const positionClass = {
    top: 'top-16',
    center: 'top-1/2 -translate-y-1/2',
    bottom: 'bottom-16',
  }[position];

  return (
    <View className={cn('absolute inset-x-0 z-50 items-center px-4', positionClass)}>
      <NativeOnlyAnimatedView
        entering={position === 'top' ? SlideInDown : FadeIn}
        exiting={position === 'top' ? SlideOutUp : FadeOut}
      >
        <View
          className={cn(
            'relative w-full max-w-md rounded-lg border-2 p-4 shadow-lg',
            getTypeClassName(type)
          )}
        >
          <View className="flex-row items-start gap-3">
            <Icon as={icon} className={cn('size-5 shrink-0 mt-0.5', getIconClassName(type))} />
            <View className="flex-1">
              <Text className="font-semibold text-base mb-1">{options.title}</Text>
              {options.description && (
                <Text className="text-sm opacity-90">{options.description}</Text>
              )}
            </View>
            {options.showClose && (
              <Pressable
                onPress={onClose}
                className="ml-2 rounded p-1 active:bg-black/10 dark:active:bg-white/10"
              >
                <Icon as={X} className="size-4 opacity-70" />
              </Pressable>
            )}
          </View>
        </View>
      </NativeOnlyAnimatedView>
    </View>
  );
}

/**
 * 确认对话框组件
 */
function ConfirmDialog({ options, onClose }: { options: ConfirmOptions; onClose: () => void }) {
  const type = options.type ?? 'warning';
  const icon = options.icon ?? getDefaultIcon(type);
  const [loading, setLoading] = useState(false);

  const handleConfirm = async () => {
    if (options.onConfirm) {
      try {
        setLoading(true);
        await options.onConfirm();
      } finally {
        setLoading(false);
        onClose();
      }
    } else {
      onClose();
    }
  };

  const handleCancel = () => {
    options.onCancel?.();
    onClose();
  };

  return (
    <Pressable 
      className="absolute inset-0 z-50 items-center justify-center bg-black/50 px-4"
      onPress={handleCancel}
    >
      <NativeOnlyAnimatedView entering={FadeIn} exiting={FadeOut}>
        <Pressable onPress={(e) => e.stopPropagation()}>
          <View className="w-full max-w-xl rounded-xl border border-border bg-card p-8 shadow-xl">
            <View className="mb-8 flex-row items-start gap-4">
              <Icon as={icon} className={cn('size-7 shrink-0 mt-0.5', getIconClassName(type))} />
              <View className="flex-1">
                <Text className="font-semibold text-xl mb-3">{options.title}</Text>
                {options.description && (
                  <Text className="text-base text-muted-foreground leading-relaxed">
                    {options.description}
                  </Text>
                )}
              </View>
            </View>

            <View className="flex-row justify-end gap-4">
              {options.onCancel !== undefined && (
                <Pressable
                  onPress={handleCancel}
                  disabled={loading}
                  className={cn(
                    'min-w-[120px] h-11 px-8 py-2.5 flex-row items-center justify-center rounded-md border border-border bg-background shadow-sm active:bg-accent',
                    loading && 'opacity-50'
                  )}
                >
                  <Text className="text-base font-medium">{options.cancelText ?? '取消'}</Text>
                </Pressable>
              )}
              {options.onConfirm !== undefined && (
                <Pressable
                  onPress={handleConfirm}
                  disabled={loading}
                  className={cn(
                    'min-w-[120px] h-11 px-8 py-2.5 flex-row items-center justify-center rounded-md shadow-sm active:opacity-90',
                    options.confirmVariant === 'destructive'
                      ? 'bg-destructive dark:bg-destructive/60'
                      : 'bg-primary',
                    loading && 'opacity-50'
                  )}
                >
                  <Text className={cn(
                    'text-base font-medium',
                    options.confirmVariant === 'destructive' ? 'text-white' : 'text-primary-foreground'
                  )}>
                    {loading ? '处理中...' : options.confirmText ?? '确定'}
                  </Text>
                </Pressable>
              )}
            </View>
          </View>
        </Pressable>
      </NativeOnlyAnimatedView>
    </Pressable>
  );
}

/**
 * 使用 Toast Hook
 */
export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast must be used within ToastProvider');
  return ctx;
}

