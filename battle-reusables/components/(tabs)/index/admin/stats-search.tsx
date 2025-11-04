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
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { StatsPath } from '@/services/stats';

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

  return (
    <View className="flex flex-row items-center gap-2 border-b border-b-border p-4">
      {/* 店铺号输入框 */}
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

      {/* 圈子ID（可选） */}
      <Controller
        control={control}
        name="groupId"
        render={({ field: { onChange, value } }) => (
          <Input
            keyboardType="numeric"
            className="w-28"
            placeholder="圈子ID(选填)"
            value={value}
            onChangeText={onChange}
          />
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
