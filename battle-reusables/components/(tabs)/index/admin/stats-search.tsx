import React, { useRef } from 'react';
import { Platform, View } from 'react-native';
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
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import z from 'zod';
import { useForm, Controller, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { StatsPath } from '@/services/stats';
import { useRequest } from '@/hooks/use-request';
import { shopsHousesOptions } from '@/services/shops/houses';
import { getGroupOptions } from '@/services/shops/groups';

// 时间选项配置
const timeOptions: { label: string; value: StatsPath }[] = [
  { label: '上周', value: '/stats/lastweek' },
  { label: '本周', value: '/stats/week' },
  { label: '昨日', value: '/stats/yesterday' },
  { label: '今日', value: '/stats/today' },
];

const searchFormSchema = z.lazy(() =>
  z.object({
    houseGid: z.string().min(1, { message: '请输入店铺号' }),
    groupId: z.string().optional(),
    timeOption: z.enum(['/stats/today', '/stats/yesterday', '/stats/week', '/stats/lastweek']),
  })
);

type StatsSearchParams = z.infer<typeof searchFormSchema>;

type SearchFormProps = {
  submitButtonProps?: ButtonProps & { loading?: boolean };
  onSubmit?: (path: StatsPath, data: API.StatsParams) => void;
};

export const StatsSearch = ({ submitButtonProps, onSubmit }: SearchFormProps) => {
  const ref = useRef<TriggerRef>(null);
  const insets = useSafeAreaInsets();
  const contentInsets = {
    top: insets.top,
    bottom: Platform.select({ ios: insets.bottom, android: insets.bottom + 24 }),
    left: 12,
    right: 12,
  };

  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  const {
    control,
    handleSubmit,
    formState: { isValid },
    setValue,
  } = useForm<StatsSearchParams>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: {
      houseGid: '',
      groupId: '',
      timeOption: timeOptions[3].value, // 默认选择今日
    },
    mode: 'all',
  });

  const handleParms = (params: StatsSearchParams) => {
    const { houseGid, groupId, timeOption } = params;
    const data: API.StatsParams = { house_gid: Number(houseGid.trim()) };
    if (groupId && groupId.trim()) data.group_id = Number(groupId.trim());
    onSubmit?.(timeOption, data);
  };

  const { loading, ...restSubmitButtonProps } = submitButtonProps ?? {};
  const { data: houseOptions } = useRequest(shopsHousesOptions);
  const houseRef = useRef<TriggerRef>(null);
  const groupRef = useRef<TriggerRef>(null);

  function onHouseTouchStart() {
    isWeb && houseRef.current?.open();
  }

  function onGroupTouchStart() {
    isWeb && groupRef.current?.open();
  }

  // 获取当前选中的店铺ID
  const currentHouseGid = useWatch({
    control,
    name: 'houseGid',
  });

  // 根据店铺ID获取圈子选项
  const { data: groupOptions, run: fetchGroupOptions } = useRequest(
    () => getGroupOptions({ house_gid: Number(currentHouseGid) }),
    {
      manual: true,
    }
  );

  // 当店铺改变时，重新获取圈子列表
  React.useEffect(() => {
    if (currentHouseGid) {
      // 清空圈子选择
      setValue('groupId', '');
      // 获取新的圈子列表
      fetchGroupOptions();
    }
  }, [currentHouseGid]);

  return (
    <View className="flex flex-row items-center gap-2 border-b border-b-border p-4">
      {/* 店铺号下拉框 */}
      <Controller
        control={control}
        name="houseGid"
        render={({ field: { onChange, value } }) => (
          <Select
            value={value ? ({ label: `店铺 ${value}`, value } as any) : undefined}
            onValueChange={(opt) => onChange(String(opt?.value ?? ''))}
          >
            <SelectTrigger ref={houseRef} onTouchStart={onHouseTouchStart} className="min-w-[160px]">
              <SelectValue placeholder={value ? `店铺 ${value}` : '选择店铺号'} />
            </SelectTrigger>
            <SelectContent insets={contentInsets}>
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

      {/* 圈子ID（可选） */}
      <Controller
        control={control}
        name="groupId"
        render={({ field: { onChange, value } }) => (
          <Select
            value={value ? ({ label: `圈子 ${value}`, value } as any) : undefined}
            onValueChange={(opt) => onChange(String(opt?.value ?? ''))}
          >
            <SelectTrigger ref={groupRef} onTouchStart={onGroupTouchStart} className="min-w-[160px]">
              <SelectValue placeholder={value ? `圈子 ${value}` : '选择圈子(选填)'} />
            </SelectTrigger>
            <SelectContent insets={contentInsets}>
              <SelectGroup>
                <SelectLabel>圈子</SelectLabel>
                {(groupOptions ?? []).map((group) => (
                  <SelectItem key={String(group.id)} label={group.name} value={String(group.id)}>
                    {group.name}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        )}
      />

      {/* 时间选择器 */}
      <Controller
        control={control}
        name="timeOption"
        render={({ field: { onChange } }) => (
          <Select defaultValue={timeOptions[3]} onValueChange={(option) => onChange(option?.value)}>
            <SelectTrigger ref={ref} onTouchStart={onTouchStart}>
              <SelectValue placeholder="请选择时间范围" />
            </SelectTrigger>
            <SelectContent insets={contentInsets}>
              <SelectGroup>
                <SelectLabel>时间范围</SelectLabel>
                {timeOptions.map(({ label, value }) => (
                  <SelectItem key={value} label={label} value={value}>
                    {label}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        )}
      />

      {/* 查询按钮 */}
      <Button
        onPress={handleSubmit(handleParms)}
        disabled={!isValid || loading}
        className="px-4"
        {...restSubmitButtonProps}>
        <Text className="font-medium">统计</Text>
      </Button>
    </View>
  );
};
