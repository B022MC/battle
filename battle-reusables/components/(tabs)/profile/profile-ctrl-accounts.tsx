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
  updateCtrlAccountStatus,
  deleteCtrlAccount,
} from '@/services/game/ctrl-account';
import { gameAccountVerify } from '@/services/game/account';
import { alert } from '@/utils/alert';
import { md5Upper } from '@/utils/md5';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';
import { Icon } from '@/components/ui/icon';
import { Plus, Trash2 } from 'lucide-react-native';

export const ProfileCtrlAccounts = () => {
  const { getLoginModeLabel } = usePlazaConsts();
  const [refreshing, setRefreshing] = useState(false);

  // Form state for creating new account
  const [mode, setMode] = useState<'account' | 'mobile'>('account');
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);



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



  const { run: runUpdateStatus, loading: updatingStatus } = useRequest(updateCtrlAccountStatus, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '状态已更新' });
      runListAccounts();
    },
  });

  const { run: runDelete, loading: deleting } = useRequest(deleteCtrlAccount, {
    manual: true,
    onSuccess: () => {
      alert.show({ title: '成功', description: '中控账号已删除' });
      runListAccounts();
    },
  });

  useEffect(() => {
    runListAccounts();
  }, []);

  const onRefresh = async () => {
    setRefreshing(true);
    await runListAccounts();
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



  const onEnableAccount = async (ctrlId: number) => {
    // 启用账号,后端的 session monitor 会自动检测并启动所有绑定店铺的会话
    await runUpdateStatus({ ctrl_id: ctrlId, status: 1 });
  };

  const onDisableAccount = async (ctrlId: number) => {
    // 停用账号,后端会自动停止所有会话
    await runUpdateStatus({ ctrl_id: ctrlId, status: 0 });
  };

  const onDeleteAccount = async (ctrlId: number, identifier: string) => {
    const confirmed = await new Promise<boolean>((resolve) => {
      alert.show({
        title: '确认删除',
        description: `确定要删除中控账号 "${identifier}" 吗?此操作不可恢复。`,
        actions: [
          { text: '取消', onPress: () => resolve(false) },
          { text: '删除', onPress: () => resolve(true), style: 'destructive' },
        ],
      });
    });

    if (confirmed) {
      await runDelete(ctrlId);
    }
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
                </View>
              </InfoCardContent>
              <InfoCardFooter>
                <View className="flex-row gap-2">
                  {/* 启用按钮 - 只在停用状态显示 */}
                  {account.status === 0 && (
                    <Button
                      variant="default"
                      size="sm"
                      onPress={() => onEnableAccount(account.id)}
                      disabled={updatingStatus}
                      className="flex-1"
                    >
                      <Text className="text-primary-foreground">
                        {updatingStatus ? '启用中...' : '启用'}
                      </Text>
                    </Button>
                  )}

                  {/* 停用按钮 - 只在启用状态显示 */}
                  {account.status === 1 && (
                    <Button
                      variant="secondary"
                      size="sm"
                      onPress={() => onDisableAccount(account.id)}
                      disabled={updatingStatus}
                      className="flex-1"
                    >
                      <Text>{updatingStatus ? '停用中...' : '停用'}</Text>
                    </Button>
                  )}

                  {/* 删除按钮 */}
                  <Button
                    variant="destructive"
                    size="sm"
                    onPress={() => onDeleteAccount(account.id, account.identifier)}
                    disabled={deleting}
                    className="flex-1"
                  >
                    <Text className="text-destructive-foreground">
                      {deleting ? '删除中...' : '删除'}
                    </Text>
                  </Button>
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

