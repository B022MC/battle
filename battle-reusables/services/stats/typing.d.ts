declare namespace API {
  type StatsParams = {
    house_gid: number;
    group_id?: number;
  };

  type StatsMemberParams = {
    house_gid: number;
    member_id: number;
  };

  type StatsLedger = {
    /**
     * 调整（可正可负）
     */
    adjust?: number;
    /**
     * 上分总额（正数）
     */
    income?: number;
    /**
     * 参与流水成员数
     */
    members_involved?: number;
    /**
     * 净变动
     */
    net?: number;
    /**
     * 下分总额（正数）
     */
    payout?: number;
    /**
     * 流水条数
     */
    records?: number;
  };

  type StatsSession = {
    /**
     * 当前在线会话数
     */
    active?: number;
  };

  type StatsWallet = {
    balance_total?: number;
    low_balance_members?: number;
    members?: number;
  };

  type StatsResult = {
    house_gid?: number;
    ledger?: StatsLedger;
    range_end?: string;
    range_start?: string;
    session?: StatsSession;
    wallet?: StatsWallet;
  };

  type ActiveByHouseItem = {
    house_gid?: number;
    active?: number;
  };
}
