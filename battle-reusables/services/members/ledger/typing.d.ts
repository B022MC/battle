declare namespace API {
  type MembersLedgerListParams = {
    house_gid: number;
    member_id?: number;
    type?: string;
    start_at?: string;
    end_at?: string;
    page?: number;
    page_size?: number;
  };

  type MembersLedgerItem = {
    id?: number;
    house_gid?: number;
    member_id?: number;
    type?: string;
    amount?: number;
    balance_before?: number;
    balance_after?: number;
    reason?: string;
    biz_no?: string;
    created_at?: string;
    operator_id?: number;
  };

  type MembersLedgerListResult = {
    list?: MembersLedgerItem[];
    total?: number;
    page?: number;
    page_size?: number;
  };
}

