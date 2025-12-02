import { View, Platform } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { useRouter } from 'expo-router';
import { useEffect, useState } from 'react';

/**
 * 导航调试组件 - 显示浏览器历史记录状态
 * 仅在 Web 平台显示
 */
export function NavigationDebugger() {
  const router = useRouter();
  const [historyLength, setHistoryLength] = useState(0);
  const [currentUrl, setCurrentUrl] = useState('');

  useEffect(() => {
    if (Platform.OS !== 'web') return;

    const updateInfo = () => {
      if (typeof window !== 'undefined') {
        setHistoryLength(window.history.length);
        setCurrentUrl(window.location.pathname);
      }
    };

    updateInfo();

    const handlePopState = () => {
      updateInfo();
    };

    window.addEventListener('popstate', handlePopState);
    
    // 监听路由变化
    const interval = setInterval(updateInfo, 1000);

    return () => {
      window.removeEventListener('popstate', handlePopState);
      clearInterval(interval);
    };
  }, []);

  // 仅在 Web 平台显示
  if (Platform.OS !== 'web') {
    return null;
  }

  return (
    <View className="fixed bottom-4 right-4 z-50 bg-card border border-border rounded-lg p-3 shadow-lg max-w-xs">
      <Text className="text-xs font-bold mb-2">🔍 导航调试</Text>
      <View className="gap-1">
        <Text className="text-xs">
          <Text className="font-semibold">当前路径:</Text> {currentUrl}
        </Text>
        <Text className="text-xs">
          <Text className="font-semibold">历史长度:</Text> {historyLength}
        </Text>
        <Text className="text-xs">
          <Text className="font-semibold">可以后退:</Text> {historyLength > 1 ? '✅ 是' : '❌ 否'}
        </Text>
      </View>
      <View className="flex-row gap-2 mt-2">
        <Button 
          size="sm" 
          variant="outline"
          onPress={() => {
            if (typeof window !== 'undefined' && window.history.length > 1) {
              window.history.back();
            } else {
              alert('无法后退：历史记录为空');
            }
          }}
        >
          <Text className="text-xs">← 后退</Text>
        </Button>
        <Button 
          size="sm" 
          variant="outline"
          onPress={() => {
            console.log('=== 历史记录详情 ===');
            console.log('URL:', window.location.href);
            console.log('历史长度:', window.history.length);
            console.log('State:', window.history.state);
          }}
        >
          <Text className="text-xs">📋 日志</Text>
        </Button>
      </View>
    </View>
  );
}
