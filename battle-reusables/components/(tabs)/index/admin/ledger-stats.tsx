import React from 'react';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardContent,
} from '@/components/shared/info-card';

type LedgerStatsProps = {
  data?: API.StatsLedger;
};

export const LedgerStats = ({ data }: LedgerStatsProps) => {
  const { income, payout, adjust, net, records, members_involved } = data ?? {};

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>流水汇总</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="上分总额" value={income} />
        <InfoCardRow label="下分总额" value={payout} />
        <InfoCardRow label="调整" value={adjust} />
        <InfoCardRow label="净变动" value={net} />
        <InfoCardRow label="流水条数" value={records} />
        <InfoCardRow label="参与成员数" value={members_involved} />
      </InfoCardContent>
    </InfoCard>
  );
};
