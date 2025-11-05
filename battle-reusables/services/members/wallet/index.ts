import { post } from '@/utils/request';

export function membersWalletGet(data: API.MembersWalletGetParams) {
  return post<API.MembersWalletItem>('/members/wallet/get', data);
}

export function membersWalletList(data: API.MembersWalletListParams) {
  return post<API.MembersWalletListResult>('/members/wallet/list', data);
}

export function membersWalletListByGroup(data: API.MembersWalletListByGroupParams) {
  return post<API.MembersWalletListResult>('/members/wallet/list_by_group', data);
}

