import React, { useRef } from 'react';
import { View, Pressable } from 'react-native';
import { Input } from '@/components/ui/input';
import { Button, ButtonProps } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import z from 'zod';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { omitResponderProps } from '@/lib/utils';
import { useRecentHouseIds } from '@/hooks/use-recent-house-ids';
import { useRequest } from '@/hooks/use-request';
import { shopsHousesOptions } from '@/services/shops/houses';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';

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
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const ref = useRef<TriggerRef>(null);

  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  return (
    <View className="flex flex-col gap-2 border-b border-b-border p-4">
      <View className="flex flex-row gap-2">
        <Controller
          control={control}
          name="houseGid"
          render={({ field: { onChange, value } }) => (
            <Select
              value={value ? ({ label: `店铺 ${value}`, value } as any) : undefined}
              onValueChange={(opt) => onChange(String(opt?.value ?? ''))}
            >
              <SelectTrigger ref={ref} onTouchStart={onTouchStart} className="min-w-[160px]">
                <SelectValue placeholder={value ? `店铺 ${value}` : '选择店铺号'} />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>店铺号</SelectLabel>
                  {(houseOptions ?? []).map((gid) => (
                    <SelectItem key={String(gid)} label={`店铺 ${gid}`} value={String(gid)}>
                      店铺 {gid}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
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

