import React from 'react';
import { View } from 'react-native';
import { Input } from '@/components/ui/input';
import { Button, ButtonProps } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import z from 'zod';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

const searchFormSchema = z.lazy(() =>
  z.object({
    houseGid: z.string().min(1, { message: '请输入店铺号' }),
  })
);

type CtrlAccountsSearchParams = z.infer<typeof searchFormSchema>;

type CtrlAccountsSearchProps = {
  submitButtonProps?: ButtonProps & { loading?: boolean };
  onSubmit?: (data: API.ShopsCtrlAccountsListParams) => void;
};

export const CtrlAccountsSearch = ({ submitButtonProps, onSubmit }: CtrlAccountsSearchProps) => {
  const {
    control,
    handleSubmit,
    formState: { isValid },
  } = useForm<CtrlAccountsSearchParams>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: { houseGid: '' },
    mode: 'all',
  });

  const handleParams = (params: CtrlAccountsSearchParams) => {
    const { houseGid } = params;
    onSubmit?.({ house_gid: Number(houseGid.trim()) });
  };

  const { loading, ...restSubmitButtonProps } = submitButtonProps ?? {};

  return (
    <View className="flex flex-row items-center gap-2 border-b border-b-border p-4">
      <Controller
        control={control}
        name="houseGid"
        render={({ field: { onChange, value } }) => (
          <Input
            keyboardType="numeric"
            className="flex-1"
            placeholder="请输入店铺号"
            value={value}
            onChangeText={onChange}
          />
        )}
      />
      <Button
        onPress={handleSubmit(handleParams)}
        disabled={!isValid || loading}
        className="px-4"
        {...restSubmitButtonProps}>
        <Text className="font-medium">搜索</Text>
      </Button>
    </View>
  );
};

