import React, { useRef, useState } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { InfoCard, InfoCardContent, InfoCardFooter, InfoCardHeader, InfoCardTitle } from '@/components/shared/info-card';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { useRequest } from '@/hooks/use-request';
import { gameAccountVerify } from '@/services/game/account';
import { shopsCtrlAccountsCreate } from '@/services/shops/ctrlAccounts';
import { md5Upper } from '@/utils/md5';
import { alert } from '@/utils/alert';

export const ProfileCtrlAccount = () => {
  const [loginMode, setLoginMode] = useState<'account' | 'mobile'>('account');
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [houseGid, setHouseGid] = useState('');

  const options = [
    { label: '游戏账号', value: 'account' as const },
    { label: '手机号', value: 'mobile' as const },
  ];
  const selected = options.find(o => o.value === loginMode);
  const triggerRef = useRef<TriggerRef>(null);
  function onTouchStart() { isWeb && triggerRef.current?.open(); }

  const { run: runVerify, loading: verifying } = useRequest(gameAccountVerify, { manual: true });
  const { run: runCreate, loading: saving } = useRequest(shopsCtrlAccountsCreate, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '已保存', description: '中控账号已创建/更新' });
      setPassword('');
    },
  });

  const onSubmit = async () => {
    if (!identifier || !password) {
      return alert.show({ title: '参数错误', description: '请输入账号与密码' });
    }
    const pwd_md5 = md5Upper(password);
    try {
      const v = await runVerify({ mode: loginMode, account: identifier, pwd_md5 });
      if (!v?.ok) return alert.show({ title: '校验失败', description: '账号或密码不正确' });
      await runCreate({ login_mode: loginMode, identifier, pwd_md5, status: 1, house_gid: houseGid ? Number(houseGid) : undefined });
    } catch {}
  };

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>中控账号（仅超级管理员）</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <View className="gap-3">
          <View className="gap-1">
            <Text variant="muted">登录方式</Text>
            <Select value={selected as any} onValueChange={(opt) => setLoginMode((opt?.value as 'account' | 'mobile') ?? 'account')}>
              <SelectTrigger ref={triggerRef} className="w-full" onTouchStart={onTouchStart}>
                <SelectValue placeholder="请选择登录方式" />
              </SelectTrigger>
              <SelectContent className="w-full">
                {options.map(o => (
                  <SelectItem key={o.value} label={o.label} value={o.value}>{o.label}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </View>
          <View className="gap-1">
            <Text variant="muted">账号或手机号</Text>
            <Input value={identifier} onChangeText={setIdentifier} placeholder="请输入账号或手机号" />
          </View>
          <View className="gap-1">
            <Text variant="muted">密码（本地MD5后提交）</Text>
            <Input value={password} onChangeText={setPassword} placeholder="请输入密码" secureTextEntry />
          </View>
          <View className="gap-1">
            <Text variant="muted">店铺号（可选，填写则创建后立即绑定）</Text>
            <Input value={houseGid} onChangeText={setHouseGid} placeholder="如：10001" />
          </View>
        </View>
      </InfoCardContent>
      <InfoCardFooter>
        <Button disabled={verifying || saving} onPress={onSubmit}><Text>验证并保存</Text></Button>
      </InfoCardFooter>
    </InfoCard>
  );
};


