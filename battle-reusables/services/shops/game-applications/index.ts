import { post } from '@/utils/request';
import type { GameApplication, GameApplicationListParams, GameApplicationRespondParams } from './typing';

/**
 * 查询游戏内申请列表（从 Plaza Session 内存读取）
 * @param data - house_gid
 * @returns 申请列表
 */
export function listGameApplications(data: GameApplicationListParams) {
  return post<GameApplication[]>('/shops/game-applications/list', data);
}

/**
 * 通过游戏内申请
 * @param data - house_gid, message_id
 * @returns 成功响应
 */
export function approveGameApplication(data: GameApplicationRespondParams) {
  return post('/shops/game-applications/approve', data);
}

/**
 * 拒绝游戏内申请
 * @param data - house_gid, message_id
 * @returns 成功响应
 */
export function rejectGameApplication(data: GameApplicationRespondParams) {
  return post('/shops/game-applications/reject', data);
}
