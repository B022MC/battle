import { Text } from '@/components/ui/text';

export const ErrorMessage = ({ name, errors }: { name: string; errors: any }) => {
  if (!errors[name]) return null;
  return <Text className="text-destructive text-xs">{errors[name]?.message}</Text>;
};
