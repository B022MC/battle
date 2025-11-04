import React from 'react';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardContent,
} from '@/components/shared/info-card';

type SessionStatsProps = {
  data?: API.StatsSession;
};

export const SessionStats = ({ data }: SessionStatsProps) => {
  const { active } = data ?? {};

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>会话</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="当前在线会话数" value={active} />
      </InfoCardContent>
    </InfoCard>
  );
};
