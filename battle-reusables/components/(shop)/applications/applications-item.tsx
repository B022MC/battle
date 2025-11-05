import React from 'react';
import { View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { useRequest } from '@/hooks/use-request';
import { alert } from '@/utils/alert';
import { applicationsApprove, applicationsReject } from '@/services/applications';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardFooter,
  InfoCardContent,
} from '@/components/shared/info-card';

type ApplicationsItemProps = {
  data?: API.ShopsApplicationsItem;
  readOnly?: boolean;
  onChanged?: () => void;
};

export const ApplicationsItem = ({ data, readOnly, onChanged }: ApplicationsItemProps) => {
  const { id, status, applier_id, applier_gid, applier_name, house_gid, created_at, type, admin_user_id } = data ?? {};

  const { run: approveRun, loading: approveLoading } = useRequest(applicationsApprove, { manual: true });
  const { run: rejectRun, loading: rejectLoading } = useRequest(applicationsReject, { manual: true });

  if (typeof id !== 'number') return <Text>参数错误</Text>;

  const handleApprove = async () => {
    await approveRun({ id });
    alert.show({ title: '已通过', duration: 800 });
    onChanged?.();
  };

  const handleReject = async () => {
    await rejectRun({ id });
    alert.show({ title: '已拒绝', duration: 800 });
    onChanged?.();
  };

  const getStatusText = (status?: number) => {
    if (status === 0) return '待审批';
    if (status === 1) return '已通过';
    if (status === 2) return '已拒绝';
    return '未知';
  };
  const getTypeText = (t?: number) => {
    if (t === 1) return '管理员';
    if (t === 2) return '加圈';
    return '未知';
  };

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>申请 #{id}</InfoCardTitle>
        <InfoCardTitle>{applier_name}</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="申请人ID" value={applier_id} />
        <InfoCardRow label="申请人GID" value={applier_gid} />
        <InfoCardRow label="店铺号" value={house_gid} />
        <InfoCardRow label="类型" value={getTypeText(type)} />
        {admin_user_id != null && <InfoCardRow label="圈主ID" value={admin_user_id} />}
        <InfoCardRow label="状态" value={getStatusText(status)} />
        <InfoCardRow label="申请时间" value={created_at ? new Date(created_at * 1000).toLocaleString() : '-'} />
      </InfoCardContent>
      {!readOnly && (
        <InfoCardFooter>
          <Button disabled={approveLoading || status !== 0} onPress={handleApprove}>
            通过
          </Button>
          <Button disabled={rejectLoading || status !== 0} onPress={handleReject}>
            拒绝
          </Button>
        </InfoCardFooter>
      )}
    </InfoCard>
  );
};

