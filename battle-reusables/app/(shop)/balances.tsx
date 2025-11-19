import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { InfoCard, InfoCardContent, InfoCardFooter, InfoCardHeader, InfoCardRow, InfoCardTitle } from '@/components/shared/info-card';
import { RouteGuard } from '@/components/auth/RouteGuard';
import { useRequest } from '@/hooks/use-request';
import { membersWalletListByGroup } from '@/services/members/wallet';

function BalancesContent() {
  const [houseGid, setHouseGid] = useState('');
  const [groupId, setGroupId] = useState('');
  const [maxBalance, setMaxBalance] = useState('');

  const { data, loading, run } = useRequest(membersWalletListByGroup, { manual: true });

  const onQuery = () => {
    if (!houseGid || !maxBalance) return;
    run({ house_gid: Number(houseGid), group_id: groupId ? Number(groupId) : undefined, max_balance: Number(maxBalance), page: 1, page_size: 100 });
  };

  return (
    <ScrollView className="flex-1 bg-secondary p-4">
      <View className="flex-row gap-2 mb-4">
        <Input className="flex-1" keyboardType="numeric" placeholder="店铺号" value={houseGid} onChangeText={setHouseGid} />
        <Input className="w-28" keyboardType="numeric" placeholder="圈ID(可选)" value={groupId} onChangeText={setGroupId} />
        <Input className="w-28" keyboardType="numeric" placeholder="≤余额(分)" value={maxBalance} onChangeText={setMaxBalance} />
        <Button disabled={loading || !houseGid || !maxBalance} onPress={onQuery}><Text>查询</Text></Button>
      </View>

      <InfoCard>
        <InfoCardHeader><InfoCardTitle>筛选结果</InfoCardTitle></InfoCardHeader>
        <InfoCardContent>
          <View className="gap-2">
            {(data?.list ?? []).map((it, idx) => (
              <InfoCardRow key={idx} label={`成员 ${it.member_id}`} value={`${it.balance}`} />
            ))}
            {!data?.list?.length && <Text className="text-muted-foreground">无数据</Text>}
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <Text className="text-muted-foreground">共 {data?.total ?? 0} 人</Text>
        </InfoCardFooter>
      </InfoCard>
    </ScrollView>
  );
}

export default function BalancesScreen() {
  return (
    <RouteGuard anyOf={['fund:wallet:view']}>
      <BalancesContent />
    </RouteGuard>
  );
}


