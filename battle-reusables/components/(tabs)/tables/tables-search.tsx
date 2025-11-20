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
import { Icon } from '@/components/ui/icon';
import { RefreshCw, Search } from 'lucide-react-native';
import { useRequest } from '@/hooks/use-request';
import { shopsTablesPull } from '@/services/shops/tables';
import { shopsHousesOptions } from '@/services/shops/houses';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { useAuthStore } from '@/hooks/use-auth-store';

const searchFormSchema = z.lazy(() =>
  z.object({
    houseGid: z.string().min(1, { message: '请输入店铺号' }),
  })
);

type TablesSearchParams = z.infer<typeof searchFormSchema>;

type TablesSearchProps = {
  submitButtonProps?: ButtonProps & { loading?: boolean };
  onSubmit?: (data: API.ShopsTablesListParams) => void;
  hideSearch?: boolean; // 隐藏搜索框（店铺管理员不需要搜索）
  defaultHouseGid?: number; // 默认店铺id（店铺管理员自动填充）
};

export const TablesSearch = ({ submitButtonProps, onSubmit, hideSearch = false, defaultHouseGid }: TablesSearchProps) => {
  const { isAuthenticated } = useAuthStore();
  const {
    control,
    handleSubmit,
    getValues,
    setValue,
    formState: { isValid },
  } = useForm<TablesSearchParams>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: { houseGid: defaultHouseGid ? String(defaultHouseGid) : '' },
    mode: 'all',
  });

  // 店铺管理员自动设置店铺id
  React.useEffect(() => {
    if (defaultHouseGid) {
      setValue('houseGid', String(defaultHouseGid));
    }
  }, [defaultHouseGid]);

  const handleParms = (params: TablesSearchParams) => {
    const { houseGid } = params;
    onSubmit?.({ house_gid: Number(houseGid.trim()) });
  };

  const { loading, ...restSubmitButtonProps } = submitButtonProps ?? {};
  const { run: pull, loading: pulling } = useRequest(shopsTablesPull, { manual: true });
  const { data: houseOptions } = useRequest(shopsHousesOptions, {
    manual: !isAuthenticated, // 未登录时不自动加载
  });

  const ref = useRef<TriggerRef>(null);
  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  const handlePull = async () => {
    const { houseGid } = getValues();
    if (!houseGid?.trim()) return;
    const house_gid = Number(houseGid.trim());
    await pull({ house_gid });
    onSubmit?.({ house_gid });
  };

  // 店铺管理员：只显示刷新按钮
  if (hideSearch) {
    return (
      <View className="flex flex-row items-center justify-end gap-2 border-b border-b-border p-4">
        <Button
          onPress={handlePull}
          disabled={!isValid || pulling}
          variant="outline"
          className="px-3"
        >
          <Icon as={RefreshCw} />
          <Text className="font-medium">刷新桌台</Text>
        </Button>
      </View>
    );
  }

  // 超级管理员：显示搜索功能
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
      <View className="flex-row gap-2">
        <Button
          onPress={handlePull}
          disabled={!isValid || pulling}
          variant="outline"
          className="px-3"
        >
          <Icon as={RefreshCw} />
          <Text className="font-medium">刷新</Text>
        </Button>
        <Button
          onPress={handleSubmit(handleParms)}
          disabled={!isValid || loading}
          className="px-4"
          {...restSubmitButtonProps}>
          <Icon as={Search} />
          <Text className="font-medium">搜索</Text>
        </Button>
      </View>
    </View>
  );
};
