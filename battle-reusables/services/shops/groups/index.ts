import { get, post } from '@/utils/request';

// 旧的圈子 API（保留兼容）
export function shopsMyGroups() {
  return get<{ house_gid: number; group_id: number }[]>('/shops/my/groups');
}

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

// ========== 新的圈子系统 API ==========

/**
 * 创建圈子
 */
export function createGroup(data: API.CreateGroupParams) {
  return post<API.ShopGroup>('/groups/create', data);
}

/**
 * 获取我的圈子
 */
export function getMyGroup(data: API.GetMyGroupParams) {
  return post<API.ShopGroup>('/groups/my', data);
}

/**
 * 列出店铺下的所有圈子
 */
export function listGroupsByHouse(data: API.ListGroupsByHouseParams) {
  return post<API.ShopGroup[]>('/groups/list', data);
}

/**
 * 从圈子移除成员
 */
export function removeMemberFromGroup(data: API.RemoveMemberParams) {
  return post<{ success: boolean }>('/groups/members/remove', data);
}

/**
 * 列出圈子成员
 */
export function listGroupMembers(data: API.ListGroupMembersParams) {
  return post<API.GroupMembersResponse>('/groups/members/list', data);
}

/**
 * 列出我加入的所有圈子（通过游戏账号反向查询）
 * 
 * 后端逻辑：
 * 1. 根据用户ID查询绑定的游戏账号
 * 2. 根据游戏账号ID查询 game_account_group 表
 * 3. 返回游戏账号加入的所有圈子
 * 
 * @returns 圈子列表
 */
export function listMyGroups() {
  return post<API.ShopGroup[]>('/groups/my/list', {});
}

/**
 * 获取圈子选项列表（用于下拉框）
 */
export function getGroupOptions(data: API.GetGroupOptionsParams) {
  return post<API.GroupOption[]>('/groups/options', data);
}

/**
 * 拉圈：将游戏成员拉入圈子（通过 game_player_id）
 */
export function pullMembersToGroup(data: API.PullMembersToGroupParams) {
  return post<API.PullMembersToGroupResponse>('/shops/members/pull-to-group', data);
}

/**
 * 踢出圈子
 */
export function removeFromGroup(data: API.RemoveFromGroupParams) {
  return post<API.RemoveFromGroupResponse>('/shops/members/remove-from-group', data);
}
