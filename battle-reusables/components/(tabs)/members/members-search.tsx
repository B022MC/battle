import React, { useRef } from 'react';
import { View } from 'react-native';
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
import { useRequest } from '@/hooks/use-request';
import { shopsHousesOptions } from '@/services/shops/houses';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';

const searchFormSchema = z.lazy(() =>
  z.object({
    houseGid: z.string().min(1, { message: '请输入店铺号' }),
  })
);

type MembersSearchParams = z.infer<typeof searchFormSchema>;

type MembersSearchProps = {
  submitButtonProps?: ButtonProps & { loading?: boolean };
  onSubmit?: (data: API.ShopsMembersListParams) => void;
};

export const MembersSearch = ({ submitButtonProps, onSubmit }: MembersSearchProps) => {
  const {
    control,
    handleSubmit,
    formState: { isValid },
  } = useForm<MembersSearchParams>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: { houseGid: '' },
    mode: 'all',
  });

  const handleParms = (params: MembersSearchParams) => {
    const { houseGid } = params;
    onSubmit?.({ house_gid: Number(houseGid.trim()) });
  };

  const { loading, ...restSubmitButtonProps } = submitButtonProps ?? {};
  const { data: houseOptions } = useRequest(shopsHousesOptions);

  const ref = useRef<TriggerRef>(null);
  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  return (
    <View className="flex flex-row items-center gap-2 border-b border-b-border p-4">
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
      {/* TODO pull 刷新按钮 */}
      <Button
        onPress={handleSubmit(handleParms)}
        disabled={!isValid || loading}
        className="px-4"
        {...restSubmitButtonProps}>
        <Text className="font-medium">搜索</Text>
      </Button>
    </View>
  );
};
