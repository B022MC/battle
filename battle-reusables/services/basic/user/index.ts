import { get, post } from '@/utils/request';

export function basicUserAddOne(data: API.BasicUserAddParams) {
  return post<API.BasicUserItem>('/basic/user/addOne', data);
}

export function basicUserDelOne(data: API.BasicUserDelParams) {
  return get<null>('/basic/user/delOne', data);
}

export function basicUserDelMany(data: API.BasicUserDelManyParams) {
  return post<null>('/basic/user/delMany', data);
}

export function basicUserGetOne(data: API.BasicUserGetParams) {
  return get<API.BasicUserItem>('/basic/user/getOne', data);
}

export function basicUserGetList(data: API.BasicUserListParams) {
  return get<API.BasicUserList>('/basic/user/getList', data);
}

export function basicUserGetOption(data: API.BasicUserListParams) {
  return get<API.BasicUserList>('/basic/user/getOption', data);
}

export function basicUserUpdateOne(data: API.BasicUserUpdateParams) {
  return post<API.BasicUserItem>('/basic/user/updateOne', data);
}

export function basicUserMeRoles() {
  return get<API.BasicUserRoles>('/basic/user/me/roles');
}

export function basicUserMePerms() {
  return get<API.BasicUserPerms>('/basic/user/me/perms');
}

export function basicUserMe() {
  return get<API.BasicUserItem>('/basic/user/me');
}

export function basicUserChangePassword(data: { old_password: string; new_password: string }) {
  return post<null>('/basic/user/changePassword', data);
}

