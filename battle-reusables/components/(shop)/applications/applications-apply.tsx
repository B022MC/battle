import React from 'react';
import { View } from 'react-native';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import z from 'zod';
import { useRequest } from '@/hooks/use-request';
import { shopsApplicationsApplyAdmin, shopsApplicationsApplyJoin } from '@/services/shops/applications';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectGroup, SelectItem, SelectLabel } from '@/components/ui/select';
import { shopsAdminsList } from '@/services/shops/admins';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';

const adminSchema = z.object({
  houseGid: z.string().min(1, { message: '请输入店铺号' }),
  note: z.string().optional(),
});

type AdminForm = z.infer<typeof adminSchema>;
const joinSchema = z.object({
  houseGid: z.string().min(1, { message: '请输入店铺号' }),
  adminUserId: z.string().min(1, { message: '请选择圈主' }),
  note: z.string().optional(),
});
type JoinForm = z.infer<typeof joinSchema>;

export const ApplicationsApply = () => {
  const { control: c1, handleSubmit: h1, reset: r1, formState: f1 } = useForm<AdminForm>({
    resolver: zodResolver(adminSchema),
    defaultValues: { houseGid: '', note: '' },
    mode: 'all',
  });
  const { run: applyAdmin, loading: la } = useRequest(shopsApplicationsApplyAdmin, { manual: true });
  const { control: c2, handleSubmit: h2, reset: r2, formState: f2, watch } = useForm<JoinForm>({
    resolver: zodResolver(joinSchema),
    defaultValues: { houseGid: '', adminUserId: '', note: '' },
    mode: 'all',
  });
  const { run: applyJoin, loading: lj } = useRequest(shopsApplicationsApplyJoin, { manual: true });
  const adminTriggerRef = React.useRef<TriggerRef>(null);
  const onTouchStartAdmin = () => { isWeb && adminTriggerRef.current?.open(); };
  const { data: adminList, run: fetchAdmins } = useRequest<API.ShopsAdminsListResult, any>(shopsAdminsList, { manual: true });
  React.useEffect(() => {
    const sub = watch((v) => {
      const gid = String(v.houseGid ?? '').trim();
      if (gid && /^\d+$/.test(gid)) fetchAdmins({ house_gid: Number(gid) });
    });
    return () => sub.unsubscribe();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onSubmitAdmin = (v: AdminForm) => {
    applyAdmin({ house_gid: Number(v.houseGid.trim()), note: v.note }).then(() => r1());
  };
  const onSubmitJoin = (v: JoinForm) => {
    applyJoin({ house_gid: Number(v.houseGid.trim()), admin_user_id: Number(v.adminUserId), note: v.note }).then(() => r2());
  };

  return (
    <View className="gap-4 border-b border-b-border p-4 bg-background">
      <Text className="font-medium">申请成为店铺管理员</Text>
      <View className="flex-row items-center gap-2 flex-wrap">
        <Controller control={c1} name="houseGid" render={({ field: { onChange, value } }) => (
          <Input keyboardType="numeric" className="w-36" placeholder="店铺号" value={value} onChangeText={onChange} />
        )} />
        <Controller control={c1} name="note" render={({ field: { onChange, value } }) => (
          <Input className="w-64" placeholder="备注（可选）" value={value} onChangeText={onChange} />
        )} />
        <Button disabled={!f1.isValid || la} onPress={h1(onSubmitAdmin)}>
          提交管理员申请
        </Button>
      </View>
      <Text className="font-medium mt-2">申请加入圈子</Text>
      <View className="flex-row items-center gap-2 flex-wrap">
        <Controller control={c2} name="houseGid" render={({ field: { onChange, value } }) => (
          <Input keyboardType="numeric" className="w-36" placeholder="店铺号" value={value} onChangeText={onChange} />
        )} />
        <Controller control={c2} name="adminUserId" render={({ field: { onChange, value } }) => (
          <Select value={value ? ({ label: value, value } as any) : undefined} onValueChange={(opt) => onChange(String(opt?.value ?? ''))}>
            <SelectTrigger ref={adminTriggerRef} onTouchStart={onTouchStartAdmin} className="w-44">
              <SelectValue placeholder="选择圈主" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>圈主</SelectLabel>
                {(adminList ?? [])
                  .filter((it: API.ShopsAdminsItem) => typeof it.user_id === 'number' && it.user_id !== 1)
                  .map((it: API.ShopsAdminsItem) => {
                    const uid = String(it.user_id);
                    const name = (it as any).nick_name as string | undefined;
                    const label = name ? `${name}（ID ${uid}）` : `用户 ${uid}`;
                    return (
                      <SelectItem key={uid} label={label} value={uid}>
                        {label}
                      </SelectItem>
                    );
                  })}
              </SelectGroup>
            </SelectContent>
          </Select>
        )} />
        <Controller control={c2} name="note" render={({ field: { onChange, value } }) => (
          <Input className="w-64" placeholder="备注（可选）" value={value} onChangeText={onChange} />
        )} />
        <Button disabled={!f2.isValid || lj} onPress={h2(onSubmitJoin)}>
          提交入圈申请
        </Button>
      </View>
    </View>
  );
}
