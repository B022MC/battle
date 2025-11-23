import { post } from '@/utils/request';

// 获取我的圈子信息（复用现有API）
export function getMyGroupInfo(data: { house_gid: number }) {
  return post<any>('/groups/my', data);
}

// 获取圈子成员列表（复用现有API）
export function getGroupMembers(data: { group_id: number; page: number; size: number }) {
  return post<API.GroupMembersResponse>('/groups/members/list', data);
}
