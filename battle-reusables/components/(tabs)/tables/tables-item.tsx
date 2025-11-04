import React, { useMemo, useState } from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { useRequest } from '@/hooks/use-request';
import { shopsTablesCheck, shopsTablesDetail, shopsTablesDismiss } from '@/services/shops/tables';
import { Icon } from '@/components/ui/icon';
import { Info, Search, Trash2 } from 'lucide-react-native';
import { PermissionGate } from '@/components/auth/PermissionGate';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardFooter,
  InfoCardContent,
} from '@/components/shared/info-card';

type TablesItemProps = {
  houseId?: number;
  data?: API.ShopsTableItem;
  onChanged?: () => void;
};

export const TablesItem = ({ houseId, data, onChanged }: TablesItemProps) => {
  const { table_id, group_id, mapped_num, kind_id, base_score } = data ?? {};

  const { run: dismiss, loading: dismissLoading } = useRequest(shopsTablesDismiss, {
    manual: true,
    onSuccess: () => {
      onChanged?.();
    },
  });
  const { run: fetchDetail, loading: detailLoading, data: detailData } = useRequest(shopsTablesDetail, { manual: true });
  const { run: runCheck, loading: checkLoading, data: checkData } = useRequest(shopsTablesCheck, { manual: true });

  const [expanded, setExpanded] = useState(false);
  const existsText = useMemo(() => {
    if (!checkData) return undefined;
    return checkData.exists_in_cache ? '缓存中存在' : '缓存中不存在';
  }, [checkData]);

  if (typeof houseId !== 'number' || typeof mapped_num !== 'number') return <Text>参数错误</Text>;

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>桌台 #{table_id}</InfoCardTitle>
        <InfoCardTitle>映射 {mapped_num}</InfoCardTitle>
        <InfoCardTitle>圈 {group_id}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="玩法" value={kind_id} />
        <InfoCardRow label="底分" value={base_score} />
        {expanded && (
          <View className="mt-2 rounded-md border border-border p-3">
            <Text className="mb-2 text-muted-foreground">详情</Text>
            {(() => {
              const triggered = (checkData as any)?.triggered ?? (detailData as any)?.triggered;
              const snap = (checkData as any)?.table ?? (detailData as any)?.table;
              return (
                <>
                  <InfoCardRow label="触发下发" value={triggered ? '是' : '否'} />
                  <InfoCardRow label="桌台ID" value={snap?.table_id ?? '-'} />
                  <InfoCardRow label="映射号" value={snap?.mapped_num ?? '-'} />
                  <InfoCardRow label="玩法" value={snap?.kind_id ?? '-'} />
                  <InfoCardRow label="底分" value={snap?.base_score ?? '-'} />
                </>
              );
            })()}
            {existsText && <InfoCardRow label="缓存检查" value={existsText} />}
          </View>
        )}
      </InfoCardContent>
      <InfoCardFooter>
        <View className="flex-row gap-2">
          <PermissionGate anyOf={["shop:table:detail","shop:table:view"]}>
            <Button
              variant="secondary"
              disabled={checkLoading}
              onPress={async () => {
                await runCheck({ house_gid: houseId, mapped_num });
                setExpanded(true);
              }}
            >
              <Icon as={Info} />
              <Text>详情</Text>
            </Button>
          </PermissionGate>
          <PermissionGate anyOf={["shop:table:check","shop:table:view"]}>
            <Button
              variant="outline"
              disabled={checkLoading}
              onPress={async () => {
                await runCheck({ house_gid: houseId, mapped_num });
                setExpanded(true);
              }}
            >
              <Icon as={Search} />
              <Text>检查</Text>
            </Button>
          </PermissionGate>
          <PermissionGate anyOf={["shop:table:dismiss"]}>
            <Button
              variant="destructive"
              disabled={dismissLoading}
              onPress={() => dismiss({ house_gid: houseId, mapped_num, kind_id })}
            >
              <Icon as={Trash2} />
              <Text>解散</Text>
            </Button>
          </PermissionGate>
        </View>
      </InfoCardFooter>
    </InfoCard>
  );
};
