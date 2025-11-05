import { post, patch } from '@/utils/request';

export function membersCreditDeposit(data: API.MembersCreditDepositParams) {
  return post<API.MembersFundsBalanceResult>('/members/credit/deposit', data);
}

export function membersCreditWithdraw(data: API.MembersCreditWithdrawParams) {
  return post<API.MembersFundsBalanceResult>('/members/credit/withdraw', data);
}

export function membersCreditForceWithdraw(data: API.MembersCreditForceWithdrawParams) {
  return post<API.MembersFundsBalanceResult>('/members/credit/force_withdraw', data);
}

export function membersLimitUpdate(data: API.MembersLimitUpdateParams) {
  return patch<unknown>('/members/limit', data);
}

