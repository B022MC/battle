import React, { useRef, useState } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { useRequest } from '@/hooks/use-request';
import { shopsAdminsAssign } from '@/services/shops/admins';
import { alert } from '@/utils/alert';

type Props = { houseId?: number };

export const AdminsAssign = ({ houseId }: Props) => {
  const [houseGid, setHouseGid] = useState(houseId ? String(houseId) : '');
  const [userId, setUserId] = useState('');
  const [role, setRole] = useState<'admin' | 'operator'>('admin');

  const triggerRef = useRef<TriggerRef>(null);
  function onTouchStart() { isWeb && triggerRef.current?.open(); }

  const { run, loading } = useRequest(shopsAdminsAssign, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '已生效', description: '已分配管理员/操作员' });
      setUserId('');
    },
  });

  const onSubmit = () => {
    if (!houseGid || !userId) return;
    run({ house_gid: Number(houseGid), user_id: Number(userId), role });
  };

  return (
    <View className="gap-2 border-b border-b-border p-4">
      <Text className="font-medium">直接指定管理员/操作员</Text>
      <View className="flex-row gap-2">
        <Input
          keyboardType="numeric"
          className="flex-1"
          placeholder="店铺号"
          value={houseGid}
          onChangeText={setHouseGid}
        />
        <Input
          keyboardType="numeric"
          className="flex-1"
          placeholder="用户ID"
          value={userId}
          onChangeText={setUserId}
        />
      </View>
      <View className="flex-row items-center gap-2">
        <Select value={{ label: role === 'admin' ? '管理员' : '操作员', value: role } as any} onValueChange={(opt) => setRole((opt?.value as 'admin' | 'operator') ?? 'admin')}>
          <SelectTrigger ref={triggerRef} className="flex-1" onTouchStart={onTouchStart}>
            <SelectValue placeholder="请选择角色" />
          </SelectTrigger>
          <SelectContent className="w-full">
            <SelectItem label="管理员" value="admin">管理员</SelectItem>
            <SelectItem label="操作员" value="operator">操作员</SelectItem>
          </SelectContent>
        </Select>
        <Button disabled={loading || !houseGid || !userId} onPress={onSubmit}><Text>指定</Text></Button>
      </View>
    </View>
  );
};


