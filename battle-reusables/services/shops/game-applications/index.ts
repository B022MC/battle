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
 * 通过游戏内申请（自动处理游戏账号入圈）
 * 
 * 后端逻辑：
 * 1. 发送通过命令到游戏服务器
 * 2. 查找或创建游戏账号（user_id 为 NULL）
 * 3. 确保管理员有圈子（如果没有则自动创建）
 * 4. 将游戏账号加入管理员的圈子
 * 5. 记录操作日志
 * 
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
