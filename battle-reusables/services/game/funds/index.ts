import { post } from '@/utils/request';

export function membersCreditDeposit(data: API.MembersCreditDepositParams) {
  return post<API.MembersFundsBalanceResult>('/members/credit/deposit', data);
}

export function membersCreditWithdraw(data: API.MembersCreditWithdrawParams) {
  return post<API.MembersFundsBalanceResult>('/members/credit/withdraw', data);
}

export function membersCreditForceWithdraw(data: API.MembersCreditForceWithdrawParams) {
  return post<API.MembersFundsBalanceResult>('/members/credit/force_withdraw', data);
}

export async function membersLimitUpdate(data: API.MembersLimitUpdateParams) {
  const response = await fetch('/members/limit', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return response.json();
}

