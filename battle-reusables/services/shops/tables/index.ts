import { post } from '@/utils/request';

export async function shopsTablesList(data: API.ShopsTablesListParams) {
  // 先触发拉取，确保会话服务端刷新了最新桌台快照
  await post<null>('/shops/tables/pull', data);
  return post<API.ShopsTablesList>('/shops/tables/list', data);
}

export function shopsTablesPull(data: API.ShopsTablesListParams) {
  return post<null>('/shops/tables/pull', data);
}

export function shopsTablesDetail(data: API.ShopsTablesItemParams) {
  return post<API.ShopsTablesDetail>('/shops/tables/detail', data);
}

export function shopsTablesCheck(data: API.ShopsTablesItemParams) {
  return post<API.ShopsTablesCheck>('/shops/tables/check', data);
}

export function shopsTablesDismiss(data: API.ShopsTablesDismissParams) {
  return post<null>('/shops/tables/dismiss', data);
}

export async function shopsMembersList(data: API.ShopsTablesListParams) {
  // 先触发拉取，确保会话服务端刷新了最新成员快照
  await post<null>('/shops/members/pull', data);
  return post<API.ShopsMembersList>('/shops/members/list', data);
}

export function shopsMembersPull(data: API.ShopsTablesListParams) {
  return post<null>('/shops/members/pull', data);
}
