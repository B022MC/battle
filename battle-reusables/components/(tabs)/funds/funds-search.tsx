import React from 'react';
import { View, Pressable } from 'react-native';
import { Input } from '@/components/ui/input';
import { Button, ButtonProps } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import z from 'zod';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { omitResponderProps } from '@/lib/utils';
import { useRecentHouseIds } from '@/hooks/use-recent-house-ids';

const searchFormSchema = z.lazy(() =>
  z.object({
    houseGid: z.string().min(1, { message: '请输入店铺号' }),
    minBalance: z.string().optional(),
    maxBalance: z.string().optional(),
  })
);

type FundsSearchParams = z.infer<typeof searchFormSchema>;

type FundsSearchProps = {
  submitButtonProps?: ButtonProps & { loading?: boolean };
  onSubmit?: (data: API.MembersWalletListParams) => void;
};

export const FundsSearch = ({ submitButtonProps, onSubmit }: FundsSearchProps) => {
  const { add, suggestions } = useRecentHouseIds();
  const {
    control,
    handleSubmit,
    formState: { isValid },
  } = useForm<FundsSearchParams>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: { houseGid: '', minBalance: '', maxBalance: '' },
    mode: 'all',
  });

  const handleParams = (params: FundsSearchParams) => {
    const { houseGid, minBalance, maxBalance } = params;
    if (houseGid && houseGid.trim()) add(Number(houseGid.trim()));
    onSubmit?.({
      house_gid: Number(houseGid.trim()),
      min_balance: minBalance ? Number(minBalance) : undefined,
      max_balance: maxBalance ? Number(maxBalance) : undefined,
    });
  };

  const { loading, ...restSubmitButtonProps } = submitButtonProps ?? {};

  return (
    <View className="flex flex-col gap-2 border-b border-b-border p-4">
      <View className="flex flex-row gap-2">
        <Controller
          control={control}
          name="houseGid"
          render={({ field: { onChange, value } }) => (
            <View className="flex-1">
              <Input
                keyboardType="numeric"
                className="w-full"
                placeholder="店铺号"
                value={value}
                onChangeText={onChange}
              />
              {suggestions(value).length > 0 && (
                <View className="mt-1 gap-1">
                  {suggestions(value).map((s) => (
                    <Pressable key={s} onPress={() => onChange(s)}>
                      <Text className="text-muted-foreground">历史：{s}</Text>
                    </Pressable>
                  ))}
                </View>
              )}
            </View>
          )}
        />
      </View>
      <View className="flex flex-row gap-2">
        <Controller
          control={control}
          name="minBalance"
          render={({ field: { onChange, value } }) => (
            <Input
              keyboardType="numeric"
              className="flex-1"
              placeholder="最小分数"
              value={value}
              onChangeText={onChange}
            />
          )}
        />
        <Controller
          control={control}
          name="maxBalance"
          render={({ field: { onChange, value } }) => (
            <Input
              keyboardType="numeric"
              className="flex-1"
              placeholder="最大分数"
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
    </View>
  );
};

