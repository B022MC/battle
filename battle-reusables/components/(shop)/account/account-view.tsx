import React from 'react';
import { ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { gameAccountMe, gameAccountBind, gameAccountDelete } from '@/services/game/account';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardRow, InfoCardFooter, InfoCardContent } from '@/components/shared/info-card';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';

export const AccountView = () => {
  const { getLoginModeLabel } = usePlazaConsts();
  const [account, setAccount] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [nickname, setNickname] = React.useState('');

  const { data: myAccount, loading, run: getAccount } = useRequest(gameAccountMe);
  const { run: bindAccount, loading: bindLoading } = useRequest(gameAccountBind, { manual: true });
  const { run: deleteAccount, loading: deleteLoading } = useRequest(gameAccountDelete, { manual: true });

  const handleBind = async () => {
    if (!account || !password) return;
    await bindAccount({
      mode: 'account',
      account,
      pwd_md5: password,
      nickname: nickname || undefined,
    });
    getAccount();
  };

  const handleDelete = async () => {
    await deleteAccount();
    getAccount();
  };

  return (
    <ScrollView className="flex-1 bg-secondary p-4">
      {myAccount ? (
        <InfoCard>
          <InfoCardHeader>
            <InfoCardTitle>我的游戏账号</InfoCardTitle>
          </InfoCardHeader>
          <InfoCardContent>
            <InfoCardRow label="账号" value={myAccount.account} />
            <InfoCardRow label="昵称" value={myAccount.nickname} />
            <InfoCardRow label="登录方式" value={getLoginModeLabel(myAccount.login_mode as any)} />
            <InfoCardRow label="状态" value={myAccount.status === 1 ? '正常' : '禁用'} />
          </InfoCardContent>
          <InfoCardFooter>
            <Button disabled={deleteLoading} onPress={handleDelete}>
              <Text>解绑</Text>
            </Button>
          </InfoCardFooter>
        </InfoCard>
      ) : (
        <InfoCard>
          <InfoCardHeader>
            <InfoCardTitle>绑定游戏账号</InfoCardTitle>
          </InfoCardHeader>
          <InfoCardContent>
            <View className="gap-2">
              <Input placeholder="账号" value={account} onChangeText={setAccount} />
              <Input placeholder="密码(MD5)" value={password} onChangeText={setPassword} secureTextEntry />
              <Input placeholder="昵称(可选)" value={nickname} onChangeText={setNickname} />
            </View>
          </InfoCardContent>
          <InfoCardFooter>
            <Button disabled={bindLoading || !account || !password} onPress={handleBind}>
              <Text>绑定</Text>
            </Button>
          </InfoCardFooter>
        </InfoCard>
      )}
    </ScrollView>
  );
};

