import { get, post, del } from '@/utils/request';

export function gameAccountVerify(data: API.GameAccountVerifyParams) {
  return post<API.GameAccountVerifyResult>('/game/accounts/verify', data);
}

/**
 * 绑定游戏账号
 * 
 * 后端逻辑：
 * 1. 如果游戏账号不存在 → 创建新账号并绑定
 * 2. 如果游戏账号存在但未绑定用户 → 直接绑定到当前用户
 * 3. 如果游戏账号已被其他用户绑定 → 返回错误
 * 
 * @param data - 绑定参数
 * @returns 游戏账号信息
 */
export function gameAccountBind(data: API.GameAccountBindParams) {
  return post<API.GameAccountItem>('/game/accounts', data);
}

export function gameAccountMe() {
  return get<API.GameAccountItem | null>('/game/accounts/me');
}

export function gameAccountMeHouses() {
  return get<API.GameAccountHouseItem | null>('/game/accounts/me/houses');
}

export function gameAccountDelete() {
  return del<null>('/game/accounts/me');
}

