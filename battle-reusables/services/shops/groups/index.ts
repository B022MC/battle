import { get } from '@/utils/request';

export function shopsMyGroups() {
  return get<{ house_gid: number; group_id: number }[]>('/shops/my/groups');
}

import { post } from '@/utils/request';

export function shopsGroupsForbid(data: API.ShopsGroupsForbidParams) {
  return post<null>('/shops/groups/forbid', data);
}

export function shopsGroupsUnforbid(data: API.ShopsGroupsForbidParams) {
  return post<null>('/shops/groups/unforbid', data);
}

export function shopsGroupsDelete(data: API.ShopsGroupsBaseParams) {
  return post<null>('/shops/groups/delete', data);
}

export function shopsGroupsBind(data: API.ShopsGroupsBindParams) {
  return post<null>('/shops/groups/bind', data);
}

export function shopsGroupsUnbind(data: API.ShopsGroupsBaseParams) {
  return post<null>('/shops/groups/unbind', data);
}

export function shopsGroupsList(data: API.ShopsGroupsBaseParams) {
  return post<number[]>('/shops/groups/list', data);
}

