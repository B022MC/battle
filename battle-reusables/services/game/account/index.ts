import { get, post, del } from '@/utils/request';

export function gameAccountVerify(data: API.GameAccountVerifyParams) {
  return post<API.GameAccountVerifyResult>('/game/accounts/verify', data);
}

export function gameAccountBind(data: API.GameAccountBindParams) {
  return post<API.GameAccountItem>('/game/accounts', data);
}

export function gameAccountMe() {
  return get<API.GameAccountItem | null>('/game/accounts/me');
}

export function gameAccountMeHouses() {
  return get<API.GameAccountHouseItem[]>('/game/accounts/me/houses');
}

export function gameAccountDelete() {
  return del<null>('/game/accounts/me');
}

