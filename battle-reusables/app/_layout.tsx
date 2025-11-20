import { AlertProvider } from '@/components/shared/alert-provider';
import { ToastProvider } from '@/components/shared/toast-provider';
import { BubbleToastContainer } from '@/components/ui/bubble-toast-container';
import '@/global.css';

import { NAV_THEME } from '@/lib/theme';
import { ThemeProvider } from '@react-navigation/native';
import { PortalHost } from '@rn-primitives/portal';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { useColorScheme } from 'nativewind';
import { SafeAreaProvider } from 'react-native-safe-area-context';

export {
  // Catch any errors thrown by the Layout component.
  ErrorBoundary,
} from 'expo-router';

export default function RootLayout() {
  const { colorScheme } = useColorScheme();

  return (
    <SafeAreaProvider>
      <ThemeProvider value={NAV_THEME[colorScheme ?? 'light']}>
        <StatusBar style={colorScheme === 'dark' ? 'light' : 'dark'} />
        {/* 暂时使用 AlertProvider，Toast系统待修复 */}
        <AlertProvider>
          <Stack>
            <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
            <Stack.Screen name="(shop)" options={{ headerShown: false }} />
            <Stack.Screen name="auth/index" options={{ headerShown: false }} />
          </Stack>
          <PortalHost />
          <BubbleToastContainer />
        </AlertProvider>
      </ThemeProvider>
    </SafeAreaProvider>
  );
}
