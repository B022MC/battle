import React, { useMemo, useState } from 'react';
import { ScrollView, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRequest } from '@/hooks/use-request';
import { shopsFeesGet, shopsFeesSet, shopsShareFeeSet, shopsPushCreditSet, shopsFeesSettlePayoffs } from '@/services/shops/fees';
import { InfoCard, InfoCardHeader, InfoCardTitle, InfoCardRow, InfoCardFooter, InfoCardContent } from '@/components/shared/info-card';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { useRef } from 'react';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';
import { usePlazaConsts } from '@/hooks/use-plaza-consts';

export const FeesView = () => {
  const [houseGid, setHouseGid] = useState('');
  const [feesJson, setFeesJson] = useState('');
  const [shareEnabled, setShareEnabled] = useState<boolean | undefined>(undefined);
  const [pushCredit, setPushCredit] = useState('');
  const [jsonMode, setJsonMode] = useState(false);
  const [formError, setFormError] = useState<string | undefined>(undefined);

  const { data: fees, loading, run: getFees } = useRequest(shopsFeesGet, { manual: true });
  const { run: setFees, loading: setFeesLoading } = useRequest(shopsFeesSet, { manual: true });
  const { run: setShare, loading: setShareLoading } = useRequest(shopsShareFeeSet, { manual: true });
  const { run: setCredit, loading: setCreditLoading } = useRequest(shopsPushCreditSet, { manual: true });
  const { run: getPayoffs, loading: payoffsLoading, data: payoffsData } = useRequest(shopsFeesSettlePayoffs, { manual: true });

  // Kind 选择
  const { data: consts, getLabel } = usePlazaConsts();
  const kindOptions = useMemo(() => (consts?.game_kinds ?? []).map(k => ({ label: k.label, value: String(k.value) })), [consts]);
  const triggerRef = useRef<TriggerRef>(null);
  function onTouchStart() { isWeb && triggerRef.current?.open(); }

  type FeeRule = { threshold: number; fee: number; kind?: string; base?: number };
  const [rules, setRules] = useState<FeeRule[]>([]);

  const handleQuery = () => {
    if (!houseGid) return;
    getFees({ house_gid: Number(houseGid) });
  };

  const handleSetFees = async () => {
    if (!houseGid) return;
    setFormError(undefined);
    let payload = feesJson;
    if (!jsonMode) {
      // 从表单规则生成 JSON
      for (const r of rules) {
        if (!(r.threshold > 0 && r.fee > 0)) {
          setFormError('请填写有效的阈值与费用（均需>0）');
          return;
        }
        if (r.kind && !(r.base && r.base > 0)) {
          setFormError('按玩法配置时需填写有效底分');
          return;
        }
      }
      payload = JSON.stringify({ rules });
      setFeesJson(payload);
    } else {
      // JSON 模式下做基本校验
      try {
        const p = JSON.parse(feesJson || '{}');
        if (p.rules && !Array.isArray(p.rules)) throw new Error('rules 应为数组');
      } catch (e: any) {
        setFormError(`JSON不合法: ${e?.message ?? e}`);
        return;
      }
    }
    await setFees({ house_gid: Number(houseGid), fees_json: payload });
    getFees({ house_gid: Number(houseGid) });
  };

  const [startAt, setStartAt] = useState('');
  const [endAt, setEndAt] = useState('');
  const handleQueryPayoffs = async () => {
    if (!houseGid || !startAt || !endAt) return;
    await getPayoffs({ house_gid: Number(houseGid), start_at: startAt, end_at: endAt });
  };

  const handleSetShare = async (enable: boolean) => {
    if (!houseGid) return;
    await setShare({ house_gid: Number(houseGid), share: enable });
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
      setShareEnabled(Boolean(fees.share_fee));
      setPushCredit(String(fees.push_credit || ''));
      // 尝试解析 fees_json -> rules
      try {
        const p = JSON.parse(fees.fees_json || '{}');
        if (Array.isArray(p?.rules)) setRules(p.rules as FeeRule[]);
      } catch {}
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
            {!jsonMode && (
              <View className="gap-3">
                <Text className="text-muted-foreground">规则按从上到下匹配：满足首条规则即采用其费用。</Text>
                {rules.map((r, idx) => (
                  <View key={idx} className="rounded-md border border-border p-3 gap-2">
                    <View className="flex-row gap-2">
                      <View className="flex-1 gap-1">
                        <Text variant="muted">阈值（分数）</Text>
                        <Input keyboardType="numeric" value={String(r.threshold ?? '')} onChangeText={(t) => setRules(rs => rs.map((it, i) => i===idx ? { ...it, threshold: Number(t||0) } : it))} />
                      </View>
                      <View className="flex-1 gap-1">
                        <Text variant="muted">费用（单位：分）</Text>
                        <Input keyboardType="numeric" value={String(r.fee ?? '')} onChangeText={(t) => setRules(rs => rs.map((it, i) => i===idx ? { ...it, fee: Number(t||0) } : it))} />
                      </View>
                    </View>
                    <View className="flex-row gap-2 items-end">
                      <View className="flex-1 gap-1">
                        <Text variant="muted">玩法（可选）</Text>
                        <Select
                          value={r.kind ? { label: r.kind, value: r.kind } as any : undefined}
                          onValueChange={(opt) => setRules(rs => rs.map((it, i) => i===idx ? { ...it, kind: opt?.label } : it))}
                        >
                          <SelectTrigger ref={triggerRef} className="w-full" onTouchStart={onTouchStart}>
                            <SelectValue placeholder="选择玩法（不选表示通用规则）" />
                          </SelectTrigger>
                          <SelectContent className="w-full max-h-72">
                            {kindOptions.map(o => (
                              <SelectItem key={o.value} label={o.label} value={o.value}>{o.label}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </View>
                      <View className="flex-1 gap-1">
                        <Text variant="muted">底分（选了玩法必填）</Text>
                        <Input keyboardType="numeric" value={String(r.base ?? '')} onChangeText={(t) => setRules(rs => rs.map((it, i) => i===idx ? { ...it, base: t ? Number(t) : undefined } : it))} />
                      </View>
                    </View>
                    <View className="flex-row gap-2">
                      <Button variant="outline" onPress={() => setRules(rs => rs.filter((_, i) => i!==idx))}><Text>删除</Text></Button>
                      {idx>0 && (
                        <Button variant="secondary" onPress={() => setRules(rs => { const a=[...rs]; const t=a[idx-1]; a[idx-1]=a[idx]; a[idx]=t; return a; })}><Text>上移</Text></Button>
                      )}
                      {idx<rules.length-1 && (
                        <Button variant="secondary" onPress={() => setRules(rs => { const a=[...rs]; const t=a[idx+1]; a[idx+1]=a[idx]; a[idx]=t; return a; })}><Text>下移</Text></Button>
                      )}
                    </View>
                  </View>
                ))}
                <Button variant="secondary" onPress={() => setRules(rs => [...rs, { threshold: 1, fee: 1 }])}><Text>新增规则</Text></Button>
              </View>
            )}
            {jsonMode && (
              <View className="gap-2">
                <Text className="text-muted-foreground">高级：直接编辑 JSON（{`{"rules":[...]}` }）。</Text>
                <Input placeholder="费用JSON" value={feesJson} onChangeText={setFeesJson} multiline />
              </View>
            )}
            {formError && <Text className="text-destructive">{formError}</Text>}
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <View className="flex-row gap-2">
            <Button disabled={setFeesLoading || !houseGid} onPress={handleSetFees}>
              <Text>保存费用设置</Text>
            </Button>
            <Button variant="outline" onPress={() => setJsonMode(v => !v)}>
              <Text>{jsonMode ? '切换到表单模式' : '切换到JSON模式'}</Text>
            </Button>
          </View>
        </InfoCardFooter>
      </InfoCard>

      <InfoCard className="mb-4">
        <InfoCardHeader>
          <InfoCardTitle>分运设置</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          <View className="gap-2">
            <InfoCardRow label="当前状态" value={shareEnabled === undefined ? '-' : (shareEnabled ? '已开启' : '已关闭')} />
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <View className="flex-row gap-2">
            <Button variant="secondary" disabled={setShareLoading || !houseGid} onPress={() => handleSetShare(true)}>
              开启分运
            </Button>
            <Button variant="outline" disabled={setShareLoading || !houseGid} onPress={() => handleSetShare(false)}>
              关闭分运
            </Button>
          </View>
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
              placeholder="推送额度（单位：分，例如 1000=¥10.00）"
              value={pushCredit}
              onChangeText={setPushCredit}
            />
            <Text className="text-muted-foreground">说明：推送额度用于机器人推送可用额度上限，单位为分。</Text>
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <Button disabled={setCreditLoading || !houseGid || !pushCredit} onPress={handleSetCredit}>
            设置额度
          </Button>
        </InfoCardFooter>
      </InfoCard>

      <InfoCard className="mt-4">
        <InfoCardHeader>
          <InfoCardTitle>费用结算分析（圈间结转）</InfoCardTitle>
        </InfoCardHeader>
        <InfoCardContent>
          <View className="gap-2">
            <View className="flex-row gap-2">
              <View className="flex-1 gap-1">
                <Text variant="muted">开始时间（RFC3339）</Text>
                <Input placeholder="2025-11-01T00:00:00Z" value={startAt} onChangeText={setStartAt} />
              </View>
              <View className="flex-1 gap-1">
                <Text variant="muted">结束时间（RFC3339）</Text>
                <Input placeholder="2025-11-08T00:00:00Z" value={endAt} onChangeText={setEndAt} />
              </View>
            </View>
            {payoffsData && (
              <View className="gap-3">
                <View>
                  <Text className="mb-1 text-muted-foreground">各圈费用汇总</Text>
                  {(payoffsData.group_sums ?? []).map((g, idx) => (
                    <InfoCardRow key={idx} label={g.play_group || '主圈'} value={g.sum} />
                  ))}
                </View>
                <View>
                  <Text className="mb-1 text-muted-foreground">圈间结转（从 → 到）</Text>
                  {(payoffsData.payoffs ?? []).map((p, idx) => (
                    <InfoCardRow key={idx} label={`${p.from_group || '主圈'} → ${p.to_group || '主圈'}`} value={p.value} />
                  ))}
                </View>
              </View>
            )}
          </View>
        </InfoCardContent>
        <InfoCardFooter>
          <Button disabled={payoffsLoading || !houseGid || !startAt || !endAt} onPress={handleQueryPayoffs}>
            计算结转
          </Button>
        </InfoCardFooter>
      </InfoCard>
    </ScrollView>
  );
};

