import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { useRequest } from '@/hooks/use-request';
import { shopsAdminsRevoke } from '@/services/shops/admins';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardFooter,
  InfoCardContent,
} from '@/components/shared/info-card';

type AdminsItemProps = {
  houseId?: number;
  data?: API.ShopsAdminsItem;
};

export const AdminsItem = ({ houseId, data }: AdminsItemProps) => {
  const { id, user_id, role } = data ?? {};

  const { run: revokeRun, loading: revokeLoading } = useRequest(shopsAdminsRevoke, { manual: true });

  if (typeof houseId !== 'number' || typeof user_id !== 'number') return <Text>参数错误</Text>;

  const handleRevoke = () => {
    revokeRun({ house_gid: houseId, user_id });
  };

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>管理员 #{id}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="用户ID" value={user_id} />
        <InfoCardRow label="角色" value={role || '管理员'} />
      </InfoCardContent>
      <InfoCardFooter>
        <Button disabled={revokeLoading} onPress={handleRevoke}>
          撤销
        </Button>
      </InfoCardFooter>
    </InfoCard>
  );
};

