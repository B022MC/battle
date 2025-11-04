import React from 'react';
import {
  InfoCard,
  InfoCardHeader,
  InfoCardTitle,
  InfoCardRow,
  InfoCardContent,
} from '@/components/shared/info-card';

type WalletStatsProps = {
  data?: API.StatsWallet;
};

export const WalletStats = ({ data }: WalletStatsProps) => {
  const { balance_total, members, low_balance_members } = data ?? {};

  return (
    <InfoCard>
      <InfoCardHeader>
        <InfoCardTitle>汇总</InfoCardTitle>
      </InfoCardHeader>
      <InfoCardContent>
        <InfoCardRow label="分数总和" value={balance_total} />
        <InfoCardRow label="成员总数" value={members} />
        <InfoCardRow label="低分数成员数" value={low_balance_members} />
      </InfoCardContent>
    </InfoCard>
  );
};
