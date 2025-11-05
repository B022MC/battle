import { post, request } from '@/utils/request';

export function shopsAdminsAssign(data: API.ShopsAdminsAssignParams) {
  return post<null>('/shops/admins', data);
}

export function shopsAdminsRevoke(data: API.ShopsAdminsRevokeParams) {
  return request<null>('/shops/admins', { method: 'DELETE', data });
}

export function shopsAdminsList(data: API.ShopsAdminsListParams) {
  return post<API.ShopsAdminsListResult>('/shops/admins/list', data);
}

