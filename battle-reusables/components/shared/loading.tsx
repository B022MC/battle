import { ActivityIndicator, View, ViewProps } from 'react-native';
import { Text } from '@/components/ui/text';
import { omitResponderProps } from '@/lib/utils';

type LoadingProps = ViewProps & {
  text?: string;
  size?: 'small' | 'large';
};

export const Loading = ({ text, size = 'large', className, ...props }: LoadingProps) => {
  return (
    <View className={className} {...omitResponderProps(props)}>
      <ActivityIndicator size={size} color="hsl(var(--primary-foreground))" />
      {text && <Text className="mt-4 text-foreground">{text}</Text>}
    </View>
  );
};
