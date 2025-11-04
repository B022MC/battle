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
    share: number;
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

  type ShopsFeesItem = {
    house_gid?: number;
    fees_json?: string;
    share_fee?: number;
    push_credit?: number;
  };

  type ShopsFeesSettleInsertResult = {
    ok?: boolean;
  };

  type ShopsFeesSettleSumResult = {
    sum?: number;
  };
}

