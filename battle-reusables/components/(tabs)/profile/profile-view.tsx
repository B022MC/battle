import React, { useMemo, useState } from 'react';
import { ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { InfoCard, InfoCardContent, InfoCardFooter, InfoCardHeader, InfoCardRow, InfoCardTitle } from '@/components/shared/info-card';
import { useAuthStore } from '@/hooks/use-auth-store';
import { useRequest } from '@/hooks/use-request';
import { basicUserMe, basicUserMePerms, basicUserMeRoles, basicUserUpdateOne, basicUserChangePassword } from '@/services/basic/user';
import { alert } from '@/utils/alert';
import { ProfileGameAccount } from './profile-game-account';
import { ProfileCtrlAccounts } from './profile-ctrl-accounts';
import { router } from 'expo-router';
import { usePermission } from '@/hooks/use-permission';
import { useColorScheme } from 'nativewind';

export const ProfileView = () => {
  const { user, roles, perms, platform, clearAuth, updateAuth, isAuthenticated } = useAuthStore();
  const { isSuperAdmin, hasAny } = usePermission();
  const { colorScheme, setColorScheme } = useColorScheme();
  const isAdmin = hasAny([
    'shop:admin:assign', 'shop:admin:view',
    'shop:member:view', 'shop:table:view', 'shop:apply:view',
    'game:ctrl:view', 'game:ctrl:update', 'game:ctrl:create',
  ]);

  // 调试信息：打印权限和角色
  React.useEffect(() => {
    console.log('=== Profile Debug Info ===');
    console.log('isSuperAdmin:', isSuperAdmin);
    console.log('isAdmin:', isAdmin);
    console.log('roles:', roles);
    console.log('perms:', perms);
    console.log('Should show ProfileGameAccount:', !(isSuperAdmin || isAdmin));
  }, [isSuperAdmin, isAdmin, roles, perms]);

  const { data: me, run: runMe } = useRequest(basicUserMe, { manual: true, onSuccess: (res) => {
    setNick(res?.nick_name ?? '');
    setAvatar(res?.avatar ?? '');
  }});

  const [nick, setNick] = useState<string>(user?.nick_name ?? '');
  const [avatar, setAvatar] = useState<string>(user?.avatar ?? '');

  const roleText = useMemo(() => (roles && roles.length ? roles.join(', ') : '-'), [roles]);
  const permsText = useMemo(() => (perms && perms.length ? `${perms.length} 个权限` : '-'), [perms]);

  const { run: runUpdate, loading: saving } = useRequest(basicUserUpdateOne, {
    manual: true,
    onSuccess: (res) => {
      updateAuth({ user: { id: res?.id, username: res?.username, nick_name: res?.nick_name, avatar: res?.avatar } });
      alert.show({ title: '已保存', description: '个人资料已更新' });
    },
  });

  const { run: runMeRoles, loading: loadingRoles } = useRequest(basicUserMeRoles, {
    manual: true,
    onSuccess: (res) => {
      updateAuth({ roles: res?.role_ids });
      alert.show({ title: '已刷新', description: '角色已刷新' });
    },
  });
  const { run: runMePerms, loading: loadingPerms } = useRequest(basicUserMePerms, {
    manual: true,
    onSuccess: (res) => {
      updateAuth({ perms: res?.perms });
      alert.show({ title: '已刷新', description: '权限已刷新' });
    },
  });

  const [oldPwd, setOldPwd] = useState('');
  const [newPwd, setNewPwd] = useState('');
  const { run: runChangePwd, loading: changingPwd } = useRequest(basicUserChangePassword, {
    manual: true,
    onSuccess: () => {
      setOldPwd('');
      setNewPwd('');
      alert.show({
        title: '修改成功',
        description: '请重新登录以生效',
        confirmText: '确定并退出',
        cancelText: '取消',
        onConfirm: () => { clearAuth(); router.push('/auth'); },
      });
    },
  });

  const onSave = () => {
    const uid = me?.id ?? user?.id;
    const uname = me?.username ?? user?.username;
    if (!uid) return alert.show({ title: '无法保存', description: '缺少用户ID' });
    runUpdate({ id: uid, username: uname ?? '', nick_name: nick, avatar });
  };

  React.useEffect(() => {
    if (!isAuthenticated) return;
    runMe();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  return (
    <View className="flex-1">
      <ScrollView className="flex-1 bg-secondary">
        <View className="gap-4 p-4">
          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>账号信息</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-3">
                <InfoCardRow label="用户名" value={me?.username ?? user?.username ?? '-'} />
                <InfoCardRow label="平台" value={platform ?? '-'} />
                <View className="gap-1">
                  <Text variant="muted">昵称</Text>
                  <Input value={nick} onChangeText={setNick} placeholder="请输入昵称" />
                </View>
                <View className="gap-1">
                  <Text variant="muted">头像URL</Text>
                  <Input value={avatar} onChangeText={setAvatar} placeholder="https://..." />
                </View>
              </View>
            </InfoCardContent>
            <InfoCardFooter>
              <View className="flex-row gap-2">
                <Button disabled={saving} onPress={onSave}><Text>保存</Text></Button>
                <Button variant="outline" onPress={() => { setNick(user?.nick_name ?? ''); setAvatar(user?.avatar ?? ''); }}><Text>重置</Text></Button>
              </View>
            </InfoCardFooter>
          </InfoCard>

          {/* 普通用户可以绑定游戏账号，管理员无需绑定 */}
          {!isSuperAdmin && !isAdmin && <ProfileGameAccount />}

          {/* 中控账号区域仅对超级管理员可见 - 使用新的综合管理组件 */}
          {isSuperAdmin && <ProfileCtrlAccounts />}

          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>角色与权限</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-3">
                <InfoCardRow label="角色ID" value={roleText} />
                <InfoCardRow label="权限" value={permsText} />
              </View>
            </InfoCardContent>
            <InfoCardFooter>
              <View className="flex-row gap-2">
                <Button variant="outline" disabled={loadingRoles} onPress={() => runMeRoles()}><Text>刷新角色</Text></Button>
                <Button variant="outline" disabled={loadingPerms} onPress={() => runMePerms()}><Text>刷新权限</Text></Button>
              </View>
            </InfoCardFooter>
          </InfoCard>

          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>外观设置 (调试)</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-3">
                <InfoCardRow label="当前主题" value={colorScheme === 'dark' ? '🌙 夜间模式' : '☀️ 浅色模式'} />
                <InfoCardRow label="主题值" value={colorScheme || 'undefined'} />
                <View className="p-3 bg-muted rounded-md">
                  <Text className="text-xs">
                    点击下方按钮切换主题。如果按钮颜色没变化，说明主题系统有问题。
                  </Text>
                </View>
              </View>
            </InfoCardContent>
            <InfoCardFooter>
              <View className="gap-2">
                <View className="flex-row gap-2">
                  <Button 
                    variant={colorScheme === 'light' ? 'default' : 'outline'}
                    onPress={() => {
                      console.log('[主题切换] 切换到浅色模式');
                      setColorScheme('light');
                    }}
                    className="flex-1"
                  >
                    <Text>☀️ 浅色</Text>
                  </Button>
                  <Button 
                    variant={colorScheme === 'dark' ? 'default' : 'outline'}
                    onPress={() => {
                      console.log('[主题切换] 切换到夜间模式');
                      setColorScheme('dark');
                    }}
                    className="flex-1"
                  >
                    <Text>🌙 夜间</Text>
                  </Button>
                </View>
                <Text className="text-xs text-center text-muted-foreground">
                  当前按钮背景色：{colorScheme === 'dark' ? '应该是蓝色' : '应该是黄色'}
                </Text>
              </View>
            </InfoCardFooter>
          </InfoCard>

          <InfoCard>
            <InfoCardHeader>
              <InfoCardTitle>安全</InfoCardTitle>
            </InfoCardHeader>
            <InfoCardContent>
              <View className="gap-2">
                <Text className="text-muted-foreground">修改密码后需要重新登录。</Text>
                <View className="gap-1">
                  <Text variant="muted">旧密码</Text>
                  <Input value={oldPwd} onChangeText={setOldPwd} secureTextEntry placeholder="请输入旧密码" />
                </View>
                <View className="gap-1">
                  <Text variant="muted">新密码</Text>
                  <Input value={newPwd} onChangeText={setNewPwd} secureTextEntry placeholder="请输入新密码" />
                </View>
              </View>
            </InfoCardContent>
            <InfoCardFooter>
              <View className="flex-row gap-2">
                <Button variant="outline" disabled={changingPwd || !oldPwd || !newPwd} onPress={() => runChangePwd({ old_password: oldPwd, new_password: newPwd })}>
                  <Text>修改密码</Text>
                </Button>
                <Button variant="destructive" onPress={() => { clearAuth(); router.push('/auth'); }}>
                  <Text>退出登录</Text>
                </Button>
              </View>
            </InfoCardFooter>
          </InfoCard>
        </View>
      </ScrollView>
    </View>
  );
};


