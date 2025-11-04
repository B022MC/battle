import React, { useEffect, useState } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { InfoCard, InfoCardContent, InfoCardFooter, InfoCardHeader, InfoCardRow, InfoCardTitle } from '@/components/shared/info-card';
import { useRequest } from '@/hooks/use-request';
import { gameAccountBind, gameAccountDelete, gameAccountMe, gameAccountVerify } from '@/services/game/account';
import { alert } from '@/utils/alert';
import { md5Upper } from '@/utils/md5';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { useRef } from 'react';
import { isWeb } from '@/utils/platform';

export const ProfileGameAccount = () => {
  const { data: me, run: runMe, loading: loadingMe } = useRequest(gameAccountMe, { manual: true });
  const { run: runVerify, loading: verifying } = useRequest(gameAccountVerify, { manual: true });
  const { run: runBind, loading: binding } = useRequest(gameAccountBind, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '已绑定', description: '游戏账号绑定成功' });
      runMe();
      setAccount('');
      setPassword('');
    },
  });
  const { run: runUnbind, loading: unbinding } = useRequest(gameAccountDelete, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '已解绑', description: '已解绑我的游戏账号' });
      runMe();
    },
  });

  const [mode, setMode] = useState<'account' | 'mobile'>('account');
  const [account, setAccount] = useState('');
  const [password, setPassword] = useState('');

  // Select helper (open on web)
  const triggerRef = useRef<TriggerRef>(null);
  function onTouchStart() { isWeb && triggerRef.current?.open(); }
  const modeOptions = [
    { label: '游戏账号', value: 'account' },
    { label: '手机号', value: 'mobile' },
  ] as const;
  const modeOption = modeOptions.find(o => o.value === mode);

  useEffect(() => {
    runMe();
  }, []);

  const onBind = async () => {
    if (!account || !password) {
      return alert.show({ title: '参数错误', description: '请输入账号与密码' });
    }
    const digest = md5Upper(password);
    try {
      const vr = await runVerify({ mode, account, pwd_md5: digest });
      if (!vr?.ok) {
        return alert.show({ title: '校验失败', description: '账号或密码不正确' });
      }
      await runBind({ mode, account, pwd_md5: digest });
    } catch (e) {
      // useRequest 已统一错误提示
    }
  };

  const isBound = !!(me && Object.keys(me as any).length > 0);

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>我的游戏账号</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        {isBound ? (
          <View className="gap-3">
            <InfoCardRow label="账号" value={me.account ?? '-'} />
            <InfoCardRow label="登录方式" value={me.login_mode ?? '-'} />
            <InfoCardRow label="状态" value={String(me.status ?? '-') } />
          </View>
        ) : (
          <View className="gap-3">
            <Text className="text-muted-foreground">当前未绑定游戏账号</Text>
            <View className="gap-1">
              <Text variant="muted">登录方式</Text>
              <Select
                value={modeOption as any}
                onValueChange={(option) => setMode((option?.value as 'account' | 'mobile') ?? 'account')}
              >
                <SelectTrigger ref={triggerRef} className="w-full" onTouchStart={onTouchStart}>
                  <SelectValue placeholder="请选择登录方式" />
                </SelectTrigger>
                <SelectContent className="w-full">
                  {modeOptions.map((o) => (
                    <SelectItem key={o.value} label={o.label} value={o.value}>{o.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </View>
            <View className="gap-1">
              <Text variant="muted">账号或手机号</Text>
              <Input value={account} onChangeText={setAccount} placeholder="账号或手机号" />
            </View>
            <View className="gap-1">
              <Text variant="muted">密码（本地加密为MD5后提交）</Text>
              <Input value={password} onChangeText={setPassword} placeholder="请输入密码" secureTextEntry />
            </View>
          </View>
        )}
      </InfoCardContent>
      <InfoCardFooter>
        {isBound ? (
          <Button variant="destructive" disabled={unbinding || loadingMe} onPress={() => runUnbind()}>
            <Text>解绑我的账号</Text>
          </Button>
        ) : (
          <Button disabled={binding || verifying} onPress={onBind}>
            <Text>绑定</Text>
          </Button>
        )}
      </InfoCardFooter>
    </InfoCard>
  );
};


