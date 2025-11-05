declare namespace API {
  type ShopsFeesGetParams = {
    house_gid: number;
  };

  type ShopsFeesSetParams = {
    house_gid: number;
    fees_json: string;
  };

  type ShopsShareFeeSetParams = {
    house_gid: number;
    share: boolean;
  };

  type ShopsPushCreditSetParams = {
    house_gid: number;
    credit: number;
  };

  type ShopsFeesSettleInsertParams = {
    house_gid: number;
    amount: number;
    date: string;
    note?: string;
  };

  type ShopsFeesSettleSumParams = {
    house_gid: number;
    start_date?: string;
    end_date?: string;
  };

  type ShopsFeesSettlePayoffsParams = {
    house_gid: number;
    start_at: string; // RFC3339
    end_at: string;   // RFC3339
  };

  type ShopsFeesItem = {
    house_gid?: number;
    fees_json?: string;
    share_fee?: boolean;
    push_credit?: number;
  };

  type ShopsFeesSettleInsertResult = {
    ok?: boolean;
  };

  type ShopsFeesSettleSumResult = {
    sum?: number;
  };

  type ShopsFeesSettlePayoffsResult = {
    group_sums: { play_group: string; sum: number }[];
    payoffs: { from_group: string; to_group: string; value: number }[];
  };
}

