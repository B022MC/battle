import { useEffect } from 'react';
import { Platform } from 'react-native';

/**
 * 调试浏览器历史记录的 Hook
 * 仅在 Web 平台上生效
 */
export function useNavigationDebug() {
  useEffect(() => {
    if (Platform.OS !== 'web') return;

    const logHistory = () => {
      console.log('=== 浏览器历史记录调试 ===');
      console.log('当前 URL:', window.location.href);
      console.log('历史记录长度:', window.history.length);
      console.log('可以后退:', window.history.length > 1);
    };

    // 初始记录
    logHistory();

    // 监听 popstate 事件（浏览器后退/前进）
    const handlePopState = (event: PopStateEvent) => {
      console.log('=== 检测到历史记录变化 ===');
      console.log('状态:', event.state);
      logHistory();
    };

    window.addEventListener('popstate', handlePopState);

    return () => {
      window.removeEventListener('popstate', handlePopState);
    };
  }, []);
}
