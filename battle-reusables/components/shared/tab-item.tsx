import { Icon } from '@/components/ui/icon';
import { Text } from '@/components/ui/text';
import { LucideIcon } from 'lucide-react-native';

type TabItemProps = {
  focused: boolean;
  icon: LucideIcon;
  label: string;
};

export const TabIcon = ({
  focused,
  icon: IconComponent,
}: Pick<TabItemProps, 'focused' | 'icon'>) => {
  return (
    <Icon
      as={IconComponent}
      className={focused ? 'text-primary' : 'text-muted-foreground'}
      size={20}
    />
  );
};

export const TabLabel = ({ focused, label }: Pick<TabItemProps, 'focused' | 'label'>) => {
  return (
    <Text
      className={`text-center text-xs leading-4 ${focused ? 'text-primary' : 'text-muted-foreground'}`}>
      {label}
    </Text>
  );
};
