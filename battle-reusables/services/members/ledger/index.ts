import { post } from '@/utils/request';

export function membersLedgerList(data: API.MembersLedgerListParams) {
  return post<API.MembersLedgerListResult>('/members/ledger/list', data);
}

