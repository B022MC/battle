declare namespace API {
  type MembersCreditDepositParams = {
    house_gid: number;
    member_id: number;
    amount: number;
    biz_no: string;
    reason?: string;
  };

  type MembersCreditWithdrawParams = {
    house_gid: number;
    member_id: number;
    amount: number;
    biz_no: string;
    reason?: string;
  };

  type MembersCreditForceWithdrawParams = {
    house_gid: number;
    member_id: number;
    amount: number;
    biz_no: string;
    reason?: string;
  };

  type MembersLimitUpdateParams = {
    house_gid: number;
    member_id: number;
    limit_min?: number;
    forbid?: boolean;
    reason?: string;
  };

  type MembersFundsBalanceResult = {
    balance?: number;
  };

  type MembersFundsLimitResult = {
    balance?: number;
    forbid?: boolean;
    limit_min?: number;
  };
}

