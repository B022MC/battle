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
import { Select, SelectTrigger, SelectValue, SelectContent, SelectGroup, SelectItem, SelectLabel } from '@/components/ui/select';
import { shopsAdminsList } from '@/services/shops/admins';
import { useRequest } from '@/hooks/use-request';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import ReactNative from 'react-native';

const searchFormSchema = z.lazy(() =>
  z.object({
    houseGid: z.string().min(1, { message: '请输入店铺号' }),
    type: z.string().optional(), // 'all' | 'admin' | 'join'
    adminUserId: z.string().optional(),
  })
);

type ApplicationsSearchParams = z.infer<typeof searchFormSchema>;

type ApplicationsSearchProps = {
  submitButtonProps?: ButtonProps & { loading?: boolean };
  onSubmit?: (data: API.ShopsApplicationsListParams) => void;
};

export const ApplicationsSearch = ({ submitButtonProps, onSubmit }: ApplicationsSearchProps) => {
  const { add, suggestions } = useRecentHouseIds();
  const typeTriggerRef = React.useRef<TriggerRef>(null);
  const adminTriggerRef = React.useRef<TriggerRef>(null);
  const onTouchStartType = () => { isWeb && typeTriggerRef.current?.open(); };
  const onTouchStartAdmin = () => { isWeb && adminTriggerRef.current?.open(); };
  const {
    control,
    handleSubmit,
    watch,
    formState: { isValid },
  } = useForm<ApplicationsSearchParams>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: { houseGid: '', type: 'all', adminUserId: '' },
    mode: 'all',
  });

  const handleParams = (params: ApplicationsSearchParams) => {
    const { houseGid, type, adminUserId } = params;
    if (houseGid && houseGid.trim()) add(Number(houseGid.trim()));
    const payload: any = { house_gid: Number(houseGid.trim()) };
    // 可选筛选条件（后端支持则透传，不支持则忽略）
    if (type && type !== 'all') {
      const map: Record<string, number> = { admin: 1, join: 2 };
      if (map[type]) payload.type = map[type];
    }
    if (adminUserId && adminUserId.trim()) payload.admin_user_id = Number(adminUserId.trim());
    onSubmit?.(payload);
  };

  const { loading, ...restSubmitButtonProps } = submitButtonProps ?? {};
  // 管理员下拉：当 houseGid 合法时拉取
  const { data: adminList, run: fetchAdmins } = useRequest(shopsAdminsList, { manual: true });
  React.useEffect(() => {
    // 监听当前表单的 houseGid 值
    const sub = watch((v) => {
      const gid = String(v.houseGid ?? '').trim();
      if (gid && /^\d+$/.test(gid)) fetchAdmins({ house_gid: Number(gid) });
    });
    return () => sub.unsubscribe();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <View className="flex flex-row flex-wrap items-center gap-2 border-b border-b-border p-4">
      <Controller
        control={control}
        name="houseGid"
        render={({ field: { onChange, value } }) => (
          <View className="flex-1 min-w-[160px]">
            <Input
              keyboardType="numeric"
              className="w-full"
              placeholder="请输入店铺号"
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
      <Controller
        control={control}
        name="type"
        render={({ field: { onChange, value } }) => (
          <Select value={value ? ({ label: value === 'admin' ? '管理员' : value === 'join' ? '加圈' : '全部', value } as any) : undefined} onValueChange={(opt) => onChange(String(opt?.value ?? 'all'))}>
            <SelectTrigger ref={typeTriggerRef} onTouchStart={onTouchStartType} className="w-28">
              <SelectValue placeholder="类型" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>类型</SelectLabel>
                <SelectItem label="全部" value="all">全部</SelectItem>
                <SelectItem label="管理员" value="admin">管理员</SelectItem>
                <SelectItem label="加圈" value="join">加圈</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        )}
      />
      <Controller
        control={control}
        name="adminUserId"
        render={({ field: { onChange, value } }) => (
          <Select
            value={value ? ({ label: value, value } as any) : undefined}
            onValueChange={(opt) => onChange(String(opt?.value ?? ''))}>
            <SelectTrigger ref={adminTriggerRef} onTouchStart={onTouchStartAdmin} className="w-36">
              <SelectValue placeholder="圈主（可选）" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>圈主</SelectLabel>
                {(adminList?.items ?? []).map((it) => (
                  <SelectItem key={String(it.user_id)} label={`用户 ${it.user_id}`} value={String(it.user_id)}>
                    用户 {it.user_id}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
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

