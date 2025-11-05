import React from 'react';
import { View } from 'react-native';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { useForm, Controller, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import z from 'zod';
import { useRequest } from '@/hooks/use-request';
import { shopsApplicationsApplyAdmin, shopsApplicationsApplyJoin, shopsApplicationsHistory } from '@/services/shops/applications';
import { ApplicationsList } from './applications-list';
import { useAuthStore } from '@/hooks/use-auth-store';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectGroup, SelectItem, SelectLabel } from '@/components/ui/select';
import { shopsAdminsList } from '@/services/shops/admins';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { alert } from '@/utils/alert';

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
  const userId = useAuthStore((s) => s.user?.id) ?? 0;
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
  const { run: fetchHistory } = useRequest(shopsApplicationsHistory, { manual: true });
  const [myItems, setMyItems] = React.useState<API.ShopsApplicationsItem[]>([]);
  const [myLoading, setMyLoading] = React.useState(false);
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

  // 自动搜索我的申请记录：监听两个表单的店铺号
  const w1 = useWatch({ control: c1, name: 'houseGid' });
  const w2 = useWatch({ control: c2, name: 'houseGid' });
  React.useEffect(() => {
    const gid = (w1 || w2 || '').trim();
    if (!gid || !/^\d+$/.test(gid)) return;
    setMyLoading(true);
    fetchHistory({ house_gid: Number(gid) })
      .then((res) => {
        const list = (res?.items ?? []).filter((it) => Number(it.applier_id) === Number(userId));
        setMyItems(list);
      })
      .finally(() => setMyLoading(false));
  }, [w1, w2, userId]);

  const onSubmitAdmin = (v: AdminForm) => {
    applyAdmin({ house_gid: Number(v.houseGid.trim()), note: v.note }).then(() => {
      alert.show({ title: '已提交管理员申请', duration: 1200 });
      r1();
      // 刷新我的申请记录
      const gid = String(v.houseGid ?? '').trim();
      if (gid && /^\d+$/.test(gid)) fetchHistory({ house_gid: Number(gid) }).then((res) => {
        const list = (res?.items ?? []).filter((it) => Number(it.applier_id) === Number(userId));
        setMyItems(list);
      });
    });
  };
  const onSubmitJoin = (v: JoinForm) => {
    applyJoin({ house_gid: Number(v.houseGid.trim()), admin_user_id: Number(v.adminUserId), note: v.note }).then(() => {
      alert.show({ title: '已提交入圈申请', duration: 1200 });
      r2();
      const gid = String(v.houseGid ?? '').trim();
      if (gid && /^\d+$/.test(gid)) fetchHistory({ house_gid: Number(gid) }).then((res) => {
        const list = (res?.items ?? []).filter((it) => Number(it.applier_id) === Number(userId));
        setMyItems(list);
      });
    });
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
          <Text>提交管理员申请</Text>
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
          <Text>提交入圈申请</Text>
        </Button>
      </View>
      <Text className="font-medium mt-4">我的申请记录</Text>
      <ApplicationsList loading={myLoading} data={myItems} readOnly />
    </View>
  );
}
