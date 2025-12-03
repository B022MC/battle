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

export function shopsMembersListPlatform(data: { house_gid: number; group_id?: number; admin_user_id?: number }) {
  return post<API.ShopsMembersList>('/shops/members/list_platform', data);
}

export function shopsMembersRemovePlatform(data: { house_gid: number; group_id?: number; member_user_id: number }) {
  return post<null>('/shops/members/remove_platform', data);
}

export function shopsMembersAddToPlatform(data: { house_gid: number; group_id?: number; member_user_id: number }) {
  return post<null>('/shops/members/add_platform', data);
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

export function shopsMembersPin(data: API.ShopsMembersPinParams) {
  return post<null>('/shops/members/pin', data);
}

export function shopsMembersUnpin(data: API.ShopsMembersUnpinParams) {
  return post<null>('/shops/members/unpin', data);
}

export function shopsMembersUpdateRemark(data: API.ShopsMembersUpdateRemarkParams) {
  return post<null>('/shops/members/update-remark', data);
}

export function shopsMembersForbid(data: { house_gid: number; game_player_id: string }) {
  return post<null>('/shops/members/forbid', data);
}

export function shopsMembersUnforbid(data: { house_gid: number; game_player_id: string }) {
  return post<null>('/shops/members/unforbid', data);
}