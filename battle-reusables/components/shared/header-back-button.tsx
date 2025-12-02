import { Pressable, Platform } from 'react-native';
import { useRouter } from 'expo-router';
import { ChevronLeft } from 'lucide-react-native';
import { useColorScheme } from 'nativewind';

/**
 * 自定义返回按钮组件
 * 确保在 Web 平台上也能正常工作
 */
export function HeaderBackButton() {
  const router = useRouter();
  const { colorScheme } = useColorScheme();
  
  // 根据主题选择颜色
  const iconColor = colorScheme === 'dark' ? '#ffffff' : '#000000';

  const handlePress = () => {
    console.log('[返回按钮] 点击返回');
    
    if (Platform.OS === 'web' && typeof window !== 'undefined') {
      // 在 Web 平台上，优先使用浏览器的历史记录
      if (window.history.length > 1) {
        console.log('[返回按钮] 使用浏览器后退');
        window.history.back();
      } else {
        console.log('[返回按钮] 使用 router.back()');
        router.back();
      }
    } else {
      // 移动平台使用 router.back()
      router.back();
    }
  };

  return (
    <Pressable
      onPress={handlePress}
      style={({ pressed }) => ({
        opacity: pressed ? 0.6 : 1,
        paddingHorizontal: 16,
        paddingVertical: 8,
        cursor: Platform.OS === 'web' ? 'pointer' : undefined,
      })}
      // 添加可访问性标签
      accessibilityLabel="返回"
      accessibilityRole="button"
    >
      <ChevronLeft size={24} color={iconColor} />
    </Pressable>
  );
}
