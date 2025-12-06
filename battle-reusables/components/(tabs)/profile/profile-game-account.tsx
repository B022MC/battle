import React, { useEffect, useState } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { InfoCard, InfoCardContent, InfoCardFooter, InfoCardHeader, InfoCardRow, InfoCardTitle } from '@/components/shared/info-card';
import { useRequest } from '@/hooks/use-request';
import { gameAccountBind, gameAccountDelete, gameAccountMe, gameAccountMeHouses, gameAccountVerify } from '@/services/game/account';
import { alert } from '@/utils/alert';
import { toast } from '@/utils/toast';
import { showSuccessBubble } from '@/utils/bubble-toast';
import { md5Upper } from '@/utils/md5';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { useRef } from 'react';
import { isWeb } from '@/utils/platform';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { useAuthStore } from '@/hooks/use-auth-store';

export const ProfileGameAccount = () => {
  const { getLoginModeLabel } = usePlazaConsts();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const { data: me, run: runMe, loading: loadingMe } = useRequest(gameAccountMe, { manual: true });
  const { data: houses, run: runHouses, loading: loadingHouses } = useRequest(gameAccountMeHouses, { manual: true });
  const { run: runVerify, loading: verifying } = useRequest(gameAccountVerify, { manual: true });
  const { run: runBind, loading: binding } = useRequest(gameAccountBind, {
    manual: true,
  });
  const { run: runUnbind, loading: unbinding } = useRequest(gameAccountDelete, {
    manual: true,
    onSuccess: () => {
      showSuccessBubble('解绑成功', '游戏账号已成功解绑');
      runMe();
      runHouses();
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
    if (!isAuthenticated) return;
    runMe();
    runHouses();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  // 调试信息
  useEffect(() => {
    console.log('=== ProfileGameAccount Debug ===');
    console.log('me:', me);
    console.log('houses:', houses);
    console.log('isBound:', !!(me && me.account));
    console.log('loadingMe:', loadingMe);
    console.log('loadingHouses:', loadingHouses);
  }, [me, houses, loadingMe, loadingHouses]);

  const onBind = async () => {
    if (!account || !password) {
      return alert.show({ title: '参数错误', description: '请输入账号与密码' });
    }
    const digest = md5Upper(password);
    
    toast.confirm({
      title: '确认绑定',
      description: `确定要绑定游戏账号 "${account}" 吗？`,
      onConfirm: async () => {
        try {
          const vr = await runVerify({ mode, account, pwd_md5: digest });
          if (!vr?.ok) {
            return toast.error('账号或密码不正确');
          }
          await runBind({ mode, account, pwd_md5: digest });
          showSuccessBubble('绑定成功', '游戏账号已成功绑定');
          runMe();
          runHouses();
          setAccount('');
          setPassword('');
        } catch (e: any) {
          // 处理特殊错误
          if (e?.message?.includes('already bound')) {
            toast.error('该游戏账号已被其他用户绑定');
          } else {
            toast.error(e?.message || '绑定失败');
          }
        }
      },
      confirmText: '绑定',
    });
  };

  const isBound = !!(me && me.account);

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>我的游戏账号</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        {loadingMe ? (
          <Text className="text-muted-foreground">加载中...</Text>
        ) : isBound ? (
          <View className="gap-3">
            <InfoCardRow label="账号" value={me.account ?? '-'} />
            <InfoCardRow label="登录方式" value={getLoginModeLabel(me.login_mode as any)} />
            <InfoCardRow label="状态" value={String(me.status ?? '-') } />

            {/* 显示绑定的店铺ID */}
            {loadingHouses ? (
              <Text className="text-muted-foreground text-sm mt-2">加载店铺信息...</Text>
            ) : houses ? (
              <View className="gap-2 mt-2">
                <Text variant="muted" className="font-semibold">绑定的店铺：</Text>
                <View className="flex-row items-center gap-2 pl-2">
                  <Text className="text-sm">
                    {houses.is_default ? '[默认]' : ''} 店铺 {houses.house_gid}
                    {houses.status === 1 ? ' (启用)' : ' (禁用)'}
                  </Text>
                </View>
              </View>
            ) : null}
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
          <Button 
            variant="destructive" 
            disabled={unbinding || loadingMe} 
            onPress={() => {
              toast.confirm({
                title: '确认解绑',
                description: `确定要解绑游戏账号 "${me.account}" 吗？`,
                type: 'error',
                confirmText: '确定解绑',
                cancelText: '取消',
                confirmVariant: 'destructive',
                onConfirm: async () => {
                  runUnbind();
                },
              });
            }}
          >
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


