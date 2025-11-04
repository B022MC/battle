import { post } from '@/utils/request';

export function membersWalletGet(data: API.MembersWalletGetParams) {
  return post<API.MembersWalletItem>('/members/wallet/get', data);
}

export function membersWalletList(data: API.MembersWalletListParams) {
  return post<API.MembersWalletListResult>('/members/wallet/list', data);
}

