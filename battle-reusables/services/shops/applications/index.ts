import { post } from '@/utils/request';

export function shopsApplicationsList(data: API.ShopsApplicationsListParams) {
  return post<API.ShopsApplicationsListResult>('/shops/applications/list', data);
}

export function shopsApplicationsApplyAdmin(data: API.ShopsApplicationsApplyAdminParams) {
  return post<API.ShopsApplicationsApplyResult>('/shops/applications/applyAdmin', data);
}

export function shopsApplicationsApplyJoin(data: API.ShopsApplicationsApplyJoinParams) {
  return post<API.ShopsApplicationsApplyResult>('/shops/applications/applyJoin', data);
}

