import { post } from '@/utils/request';

/**
 * 设置房间额度限制
 */
export function setRoomCreditLimit(data: API.SetRoomCreditLimitParams) {
  return post<API.RoomCreditLimitItem>('/room-credit/set', data);
}

/**
 * 查询房间额度限制
 */
export function getRoomCreditLimit(data: API.GetRoomCreditLimitParams) {
  return post<API.RoomCreditLimitItem>('/room-credit/get', data);
}

/**
 * 列出房间额度限制
 */
export function listRoomCreditLimits(data: API.ListRoomCreditLimitParams) {
  return post<API.RoomCreditLimitListResult>('/room-credit/list', data);
}

/**
 * 删除房间额度限制
 */
export function deleteRoomCreditLimit(data: API.DeleteRoomCreditLimitParams) {
  return post<null>('/room-credit/delete', data);
}

/**
 * 检查玩家是否满足房间额度要求
 */
export function checkPlayerCredit(data: API.CheckPlayerCreditParams) {
  return post<API.CheckPlayerCreditResult>('/room-credit/check', data);
}

