import { post } from '@/utils/request';

/**
 * 列出所有用户（支持分页和搜索）
 */
export function listAllUsers(data: API.ListAllUsersParams) {
  return post<API.AllUsersResponse>('/members/list', data);
}

/**
 * 获取用户信息
 */
export function getUser(data: API.GetUserParams) {
  return post<API.BasicUser>('/members/get', data);
}

/**
 * 列出店铺管理员
 */
export function listShopAdmins(data: API.ListShopAdminsParams) {
  return post<API.BasicUser[]>('/members/shop-admins', data);
}

