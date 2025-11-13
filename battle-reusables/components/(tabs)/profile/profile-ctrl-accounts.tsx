// Super Admin Control Account Management Component
import React, { useEffect, useState, useRef } from 'react';
import { View, ScrollView, RefreshControl } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  InfoCard,
  InfoCardContent,
  InfoCardFooter,
  InfoCardHeader,
  InfoCardRow,
  InfoCardTitle,
} from '@/components/shared/info-card';
import { useRequest } from '@/hooks/use-request';
import {
  createCtrlAccount,
  listAllCtrlAccounts,
  bindCtrlAccount,
  unbindCtrlAccount,
  startSession,
  stopSession,
  getHouseOptions,
} from '@/services/game/ctrl-account';
import { gameAccountVerify } from '@/services/game/account';
import { alert } from '@/utils/alert';
import { md5Upper } from '@/utils/md5';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { Icon } from '@/components/ui/icon';
import { Plus, Trash2, Play, Square, Link, Unlink } from 'lucide-react-native';

export const ProfileCtrlAccounts = () => {
  const { getLoginModeLabel } = usePlazaConsts();
  const [refreshing, setRefreshing] = useState(false);

  // Form state for creating new account
  const [mode, setMode] = useState<'account' | 'mobile'>('account');
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);

  // Binding state
  const [bindingAccountId, setBindingAccountId] = useState<number | null>(null);
  const [bindHouseGid, setBindHouseGid] = useState('');

  // Select refs
  const modeRef = useRef<TriggerRef>(null);
  const modeOptions = [
    { label: '游戏账号', value: 'account' },
    { label: '手机号', value: 'mobile' },
  ] as const;
  const modeOption = modeOptions.find((o) => o.value === mode);

  // API requests
  const { data: accounts, run: runListAccounts, loading: loadingAccounts } = useRequest(
    listAllCtrlAccounts,
    { manual: true }
  );

  const { data: houseOptions, run: runGetHouseOptions } = useRequest(getHouseOptions, {
    manual: true,
  });

  const { run: runVerify, loading: verifying } = useRequest(gameAccountVerify, { manual: true });

  const { run: runCreate, loading: creating } = useRequest(createCtrlAccount, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '中控账号已创建' });
      setShowAddForm(false);
      setIdentifier('');
      setPassword('');
      runListAccounts();
    },
  });

  const { run: runBind, loading: binding } = useRequest(bindCtrlAccount, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '已绑定店铺' });
      setBindingAccountId(null);
      setBindHouseGid('');
      runListAccounts();
    },
  });

  const { run: runUnbind, loading: unbinding } = useRequest(unbindCtrlAccount, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '已解绑店铺' });
      runListAccounts();
    },
  });

  const { run: runStartSession, loading: startingSession } = useRequest(startSession, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '会话已启动' });
      runListAccounts();
    },
  });

  const { run: runStopSession, loading: stoppingSession } = useRequest(stopSession, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '会话已停止' });
      runListAccounts();
    },
  });

  useEffect(() => {
    runListAccounts();
    runGetHouseOptions();
  }, []);

  const onRefresh = async () => {
    setRefreshing(true);
    await runListAccounts();
    await runGetHouseOptions();
    setRefreshing(false);
  };

  const onCreateAccount = async () => {
    if (!identifier || !password) {
      return alert.show({ title: '参数错误', description: '请输入账号与密码' });
    }

    const pwdMD5 = md5Upper(password);

    // Verify first
    try {
      const vr = await runVerify({ mode, account: identifier, pwd_md5: pwdMD5 });
      if (!vr?.ok) {
        return alert.show({ title: '验证失败', description: '账号或密码不正确' });
      }
      await runCreate({ login_mode: mode, identifier, pwd_md5: pwdMD5, status: 1 });
    } catch (e) {
      // Error already handled by useRequest
    }
  };

  const onBindHouse = async (ctrlId: number) => {
    if (!bindHouseGid) {
      return alert.show({ title: '参数错误', description: '请输入店铺号' });
    }
    const houseGid = parseInt(bindHouseGid, 10);
    if (isNaN(houseGid)) {
      return alert.show({ title: '参数错误', description: '店铺号必须是数字' });
    }
    await runBind({ ctrl_id: ctrlId, house_gid: houseGid, status: 1 });
  };

  const onUnbindHouse = async (ctrlId: number, houseGid: number) => {
    await runUnbind({ ctrl_id: ctrlId, house_gid: houseGid });
  };

  const onStartSession = async (ctrlId: number, houseGid: number) => {
    await runStartSession({ id: ctrlId, house_gid: houseGid });
  };

  const onStopSession = async (ctrlId: number, houseGid: number) => {
    await runStopSession({ id: ctrlId, house_gid: houseGid });
  };

  function onTouchStart() {
    isWeb && modeRef.current?.open();
  }

  return (
    <ScrollView
      className="flex-1 bg-secondary"
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
    >
      <View className="gap-4 p-4">
        {/* Add New Account Button */}
        <Button onPress={() => setShowAddForm(!showAddForm)} variant="default">
          <Icon as={Plus} size={16} className="text-primary-foreground mr-2" />
          <Text className="text-primary-foreground">添加中控账号</Text>
        </Button>

        {/* Add Account Form */}
        {showAddForm && (
          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>创建中控账号</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-3">
                <View className="gap-1">
                  <Text variant="muted">登录方式</Text>
                  <Select
                    value={modeOption as any}
                    onValueChange={(option) => setMode((option?.value as 'account' | 'mobile') ?? 'account')}
                  >
                    <SelectTrigger ref={modeRef} className="w-full" onTouchStart={onTouchStart}>
                      <SelectValue placeholder="请选择登录方式" />
                    </SelectTrigger>
                    <SelectContent className="w-full">
                      {modeOptions.map((o) => (
                        <SelectItem key={o.value} label={o.label} value={o.value}>
                          {o.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </View>
                <View className="gap-1">
                  <Text variant="muted">账号或手机号</Text>
                  <Input value={identifier} onChangeText={setIdentifier} placeholder="账号或手机号" />
                </View>
                <View className="gap-1">
                  <Text variant="muted">密码</Text>
                  <Input value={password} onChangeText={setPassword} placeholder="请输入密码" secureTextEntry />
                </View>
              </View>
            </InfoCardContent>
            <InfoCardFooter>
              <Button variant="outline" onPress={() => setShowAddForm(false)} className="mr-2">
                <Text>取消</Text>
              </Button>
              <Button disabled={creating || verifying} onPress={onCreateAccount}>
                <Text>{verifying ? '验证中...' : creating ? '创建中...' : '创建'}</Text>
              </Button>
            </InfoCardFooter>
          </InfoCard>
        )}

        {/* Account List */}
        {loadingAccounts && !accounts ? (
          <Text className="text-center text-muted-foreground">加载中...</Text>
        ) : accounts && accounts.length > 0 ? (
          accounts.map((account) => (
            <InfoCard key={account.id}>
              <InfoCardHeader>
                <InfoCardTitle>中控账号 #{account.id}</InfoCardTitle>
              </InfoCardHeader>
              <InfoCardContent>
                <View className="gap-2">
                  <InfoCardRow label="账号" value={account.identifier} />
                  <InfoCardRow label="登录方式" value={getLoginModeLabel(account.login_mode as any)} />
                  <InfoCardRow label="状态" value={account.status === 1 ? '启用' : '禁用'} />
                  <InfoCardRow
                    label="绑定店铺"
                    value={account.houses.length > 0 ? account.houses.join(', ') : '未绑定'}
                  />

                  {/* Bind House Form */}
                  {bindingAccountId === account.id ? (
                    <View className="gap-2 mt-2 p-3 bg-muted rounded-lg">
                      <Text variant="muted">绑定店铺号</Text>
                      <Input
                        value={bindHouseGid}
                        onChangeText={setBindHouseGid}
                        placeholder="输入店铺号"
                        keyboardType="numeric"
                      />
                      <View className="flex-row gap-2">
                        <Button
                          variant="outline"
                          onPress={() => {
                            setBindingAccountId(null);
                            setBindHouseGid('');
                          }}
                          className="flex-1"
                        >
                          <Text>取消</Text>
                        </Button>
                        <Button disabled={binding} onPress={() => onBindHouse(account.id)} className="flex-1">
                          <Text>{binding ? '绑定中...' : '确认绑定'}</Text>
                        </Button>
                      </View>
                    </View>
                  ) : null}
                </View>
              </InfoCardContent>
              <InfoCardFooter>
                <View className="gap-2">
                  {/* 绑定店铺按钮 - 只有未绑定时显示 */}
                  {account.houses.length === 0 && (
                    <Button
                      variant="outline"
                      size="sm"
                      onPress={() => {
                        setBindingAccountId(account.id);
                        setBindHouseGid('');
                      }}
                    >
                      <Icon as={Link} size={14} className="mr-1" />
                      <Text>绑定店铺</Text>
                    </Button>
                  )}

                  {/* 每个店铺的操作按钮 - 垂直排列避免超出 */}
                  {account.houses.map((houseGid) => (
                    <View key={houseGid} className="gap-2">
                      <Text variant="muted" className="text-xs">店铺 {houseGid}</Text>
                      <View className="flex-row flex-wrap gap-2">
                        <Button
                          variant="default"
                          size="sm"
                          onPress={() => onStartSession(account.id, houseGid)}
                          disabled={startingSession}
                          className="flex-1 min-w-[80px]"
                        >
                          <Icon as={Play} size={14} className="mr-1 text-primary-foreground" />
                          <Text className="text-primary-foreground">启动</Text>
                        </Button>
                        <Button
                          variant="secondary"
                          size="sm"
                          onPress={() => onStopSession(account.id, houseGid)}
                          disabled={stoppingSession}
                          className="flex-1 min-w-[80px]"
                        >
                          <Icon as={Square} size={14} className="mr-1" />
                          <Text>停止</Text>
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          onPress={() => onUnbindHouse(account.id, houseGid)}
                          disabled={unbinding}
                          className="flex-1 min-w-[80px]"
                        >
                          <Icon as={Unlink} size={14} className="mr-1 text-destructive-foreground" />
                          <Text className="text-destructive-foreground">解绑</Text>
                        </Button>
                      </View>
                    </View>
                  ))}
                </View>
              </InfoCardFooter>
            </InfoCard>
          ))
        ) : (
          <Text className="text-center text-muted-foreground">暂无中控账号</Text>
        )}
      </View>
    </ScrollView>
  );
};

