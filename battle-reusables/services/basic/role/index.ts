import { get } from '@/utils/request';

export function basicRoleList(data?: API.BasicRoleListParams) {
  return get<API.BasicRoleListResult>('/basic/role/list', data);
}

export function basicRoleGetOne(data: API.BasicRoleGetParams) {
  return get<API.BasicRoleItem>('/basic/role/getOne', data);
}

export function basicRoleAll() {
  return get<API.BasicRoleAllResult>('/basic/role/all');
}

