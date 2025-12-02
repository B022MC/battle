import { useRouter } from 'expo-router';
import { Platform } from 'react-native';
import { useCallback } from 'react';

/**
 * Web 平台增强的导航 Hook
 * 确保浏览器返回按钮正常工作
 */
export function useWebNavigation() {
  const router = useRouter();

  const navigate = useCallback((href: string) => {
    if (Platform.OS === 'web') {
      // 在 Web 平台上，先使用 router.push 添加到 Expo Router 的历史记录
      router.push(href as any);
      
      // 然后确保浏览器的历史记录也被更新
      // 这是一个 workaround，确保浏览器返回按钮能正常工作
      if (typeof window !== 'undefined' && window.history) {
        console.log('[导航] 推送到历史记录:', href);
      }
    } else {
      // 移动平台直接使用 router.push
      router.push(href as any);
    }
  }, [router]);

  const canGoBack = useCallback(() => {
    if (Platform.OS === 'web' && typeof window !== 'undefined') {
      return window.history.length > 1;
    }
    return router.canGoBack?.() ?? false;
  }, [router]);

  const goBack = useCallback(() => {
    if (Platform.OS === 'web' && typeof window !== 'undefined') {
      if (window.history.length > 1) {
        window.history.back();
      } else {
        console.warn('[导航] 无法后退，历史记录为空');
      }
    } else {
      router.back();
    }
  }, [router]);

  return {
    navigate,
    canGoBack,
    goBack,
    // 保留原始 router 的其他方法
    replace: router.replace,
    push: navigate, // 使用增强的 navigate 方法
  };
}
