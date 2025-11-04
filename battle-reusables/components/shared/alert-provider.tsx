import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { View, Pressable } from 'react-native';
import { Text } from '@/components/ui/text';
import { Portal } from '@rn-primitives/portal';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { NativeOnlyAnimatedView } from '@/components/ui/native-only-animated-view';
import { FadeIn, FadeOut } from 'react-native-reanimated';
import { bindAlert, unbindAlert, type AlertOptions } from '@/utils/alert';
import { InfoIcon, X } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';

type AlertContextValue = {
  show: (options: AlertOptions) => void;
  close: () => void;
};

const AlertContext = createContext<AlertContextValue | null>(null);

export function AlertProvider({ children }: { children: React.ReactNode }) {
  const [current, setCurrent] = useState<AlertOptions | null>(null);

  const close = useCallback(() => setCurrent(null), []);

  const show = useCallback((options: AlertOptions) => {
    setCurrent(options);
    // 有确认操作时，不自动关闭；否则按 duration 自动关闭
    if (options.onConfirm || options.onCancel) {
      return () => {};
    }
    const duration = !options.duration || options.duration <= 0 ? 800 : options.duration;
    const timer = setTimeout(() => setCurrent(null), duration);
    return () => clearTimeout(timer);
  }, []);

  // 绑定全局 alert API（供非组件代码调用）
  useEffect(() => {
    bindAlert({ show });
    return () => unbindAlert();
  }, [show]);

  return (
    <AlertContext.Provider value={{ show, close }}>
      {children}
      {current && (
        <Portal name="alert-portal">
          <View className="absolute inset-0 z-50 items-center justify-center">
            {/* 弹层内容：使用 Reusables 的 Alert */}
            <NativeOnlyAnimatedView entering={FadeIn} exiting={FadeOut}>
              <Alert icon={current.icon ?? InfoIcon} className="max-w-[400px]">
                <AlertTitle>{current.title}</AlertTitle>
                {current.description ? (
                  <AlertDescription>{current.description}</AlertDescription>
                ) : null}
                {(current.onConfirm || current.onCancel) && (
                  <View className="mt-3 flex-row justify-end gap-2">
                    {current.onCancel && (
                      <Pressable
                        className="rounded-md border border-border px-3 py-2 active:bg-accent"
                        onPress={() => {
                          try { current.onCancel?.(); } finally { close(); }
                        }}>
                        <Text>{current.cancelText ?? '取消'}</Text>
                      </Pressable>
                    )}
                    {current.onConfirm && (
                      <Pressable
                        className="rounded-md bg-primary px-3 py-2 active:bg-primary/90"
                        onPress={() => {
                          try { current.onConfirm?.(); } finally { close(); }
                        }}>
                        <Text className="text-primary-foreground">{current.confirmText ?? '确定'}</Text>
                      </Pressable>
                    )}
                  </View>
                )}
                {current.showClose && (
                  <View className="mt-3 flex-row justify-end">
                    <Pressable
                      className="rounded-md border border-border px-3 py-2 active:bg-accent"
                      onPress={close}>
                      <Icon as={X} size={16} />
                    </Pressable>
                  </View>
                )}
              </Alert>
            </NativeOnlyAnimatedView>
          </View>
        </Portal>
      )}
    </AlertContext.Provider>
  );
}

export function useAlert() {
  const ctx = useContext(AlertContext);
  if (!ctx) throw new Error('useAlert must be used within AlertProvider');
  return ctx;
}
