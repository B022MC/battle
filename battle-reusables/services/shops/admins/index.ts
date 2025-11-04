import { post } from '@/utils/request';

export function shopsAdminsAssign(data: API.ShopsAdminsAssignParams) {
  return post<null>('/shops/admins', data);
}

export async function shopsAdminsRevoke(data: API.ShopsAdminsRevokeParams) {
  const response = await fetch('/shops/admins', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return response.json();
}

export function shopsAdminsList(data: API.ShopsAdminsListParams) {
  return post<API.ShopsAdminsListResult>('/shops/admins/list', data);
}

