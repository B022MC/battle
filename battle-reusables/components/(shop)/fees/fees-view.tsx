import React, { useState } from 'react';
import { ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { shopsFeesGet, shopsFeesSet, shopsShareFeeSet, shopsPushCreditSet } from '@/services/shops/fees';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardRow, InfoCardFooter, InfoCardContent } from '@/components/shared/info-card';

export const FeesView = () => {
  const [houseGid, setHouseGid] = useState('');
  const [feesJson, setFeesJson] = useState('');
  const [shareFee, setShareFee] = useState('');
  const [pushCredit, setPushCredit] = useState('');

  const { data: fees, loading, run: getFees } = useRequest(shopsFeesGet, { manual: true });
  const { run: setFees, loading: setFeesLoading } = useRequest(shopsFeesSet, { manual: true });
  const { run: setShare, loading: setShareLoading } = useRequest(shopsShareFeeSet, { manual: true });
  const { run: setCredit, loading: setCreditLoading } = useRequest(shopsPushCreditSet, { manual: true });

  const handleQuery = () => {
    if (!houseGid) return;
    getFees({ house_gid: Number(houseGid) });
  };

  const handleSetFees = async () => {
    if (!houseGid || !feesJson) return;
    await setFees({ house_gid: Number(houseGid), fees_json: feesJson });
    getFees({ house_gid: Number(houseGid) });
  };

  const handleSetShare = async () => {
    if (!houseGid || !shareFee) return;
    await setShare({ house_gid: Number(houseGid), share: Number(shareFee) });
    getFees({ house_gid: Number(houseGid) });
  };

  const handleSetCredit = async () => {
    if (!houseGid || !pushCredit) return;
    await setCredit({ house_gid: Number(houseGid), credit: Number(pushCredit) });
    getFees({ house_gid: Number(houseGid) });
  };

  React.useEffect(() => {
    if (fees) {
      setFeesJson(fees.fees_json || '');
      setShareFee(String(fees.share_fee || ''));
      setPushCredit(String(fees.push_credit || ''));
    }
  }, [fees]);

  return (
    <ScrollView className="flex-1 bg-secondary p-4">
      <View className="mb-4">
        <View className="flex flex-row gap-2">
          <Input
            keyboardType="numeric"
            className="flex-1"
            placeholder="店铺号"
            value={houseGid}
            onChangeText={setHouseGid}
          />
          <Button disabled={!houseGid || loading} onPress={handleQuery}>
            <Text>查询</Text>
          </Button>
        </View>
      </View>

      <InfoCard className="mb-4">
        <InfoCardHeader>
          <InfoCardTitle>费用设置</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          <View className="gap-2">
            <Input placeholder="费用JSON" value={feesJson} onChangeText={setFeesJson} multiline />
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <Button disabled={setFeesLoading || !houseGid || !feesJson} onPress={handleSetFees}>
            设置费用
          </Button>
        </InfoCardFooter>
      </InfoCard>

      <InfoCard className="mb-4">
        <InfoCardHeader>
          <InfoCardTitle>分运设置</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          <View className="gap-2">
            <Input
              keyboardType="numeric"
              placeholder="分运比例"
              value={shareFee}
              onChangeText={setShareFee}
            />
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <Button disabled={setShareLoading || !houseGid || !shareFee} onPress={handleSetShare}>
            设置分运
          </Button>
        </InfoCardFooter>
      </InfoCard>

      <InfoCard>
        <InfoCardHeader>
          <InfoCardTitle>推送额度设置</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          <View className="gap-2">
            <Input
              keyboardType="numeric"
              placeholder="推送额度"
              value={pushCredit}
              onChangeText={setPushCredit}
            />
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <Button disabled={setCreditLoading || !houseGid || !pushCredit} onPress={handleSetCredit}>
            设置额度
          </Button>
        </InfoCardFooter>
      </InfoCard>
    </ScrollView>
  );
};

