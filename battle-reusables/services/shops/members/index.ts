import { post } from '@/utils/request';

export function shopsMembersList(data: API.ShopsMembersListParams) {
  return post<API.ShopsMembersList>('/shops/members/list', data);
}

export function shopsMembersPull(data: API.ShopsMembersListParams) {
  return post<null>('/shops/members/pull', data);
}

export function shopsMembersKick(data: API.ShopsMembersItemParams) {
  return post<null>('/shops/members/kick', data);
}

export function shopsMembersLogout(data: API.ShopsMembersItemParams) {
  return post<null>('/shops/members/logout', data);
}

export function shopsMembersDiamondQuery(data: API.ShopsMembersListParams) {
  return post<API.ShopsMembersDiamond>('/shops/diamond/query', data);
}

export function shopsMembersRulesVip(data: API.ShopsMembersRulesVipParams) {
  return post<null>('/shops/members/rules/vip', data);
}

export function shopsMembersRulesMulti(data: API.ShopsMembersRulesMultiParams) {
  return post<null>('/shops/members/rules/multi', data);
}

export function shopsMembersRulesTempRelease(data: API.ShopsMembersRulesTempReleaseParams) {
  return post<null>('/shops/members/rules/temp_release', data);
}