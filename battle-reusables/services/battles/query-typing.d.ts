/**
 * ===== 普通用户接口参数 =====
 */

/**
 * 查询我的战绩参数
 */
export interface ListMyBattlesParams {
  house_gid: number;
  group_id?: number; // 可选,不传则查询所有圈子
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
  page?: number;
  size?: number;
}

/**
 * 查询我的余额参数
 */
export interface GetMyBalancesParams {
  house_gid: number;
  group_id?: number; // 可选,不传则查询所有圈子
}

/**
 * 查询我的统计参数
 */
export interface GetMyStatsParams {
  house_gid: number;
  group_id?: number; // 可选,不传则查询所有圈子
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
}

/**
 * ===== 管理员接口参数 =====
 */

/**
 * 查询圈子战绩参数
 */
export interface ListGroupBattlesParams {
  house_gid: number;
  group_id: number;
  player_game_id?: number; // 可选,不传则查询所有成员
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
  page?: number;
  size?: number;
}

/**
 * 查询圈子成员余额参数
 */
export interface ListGroupMemberBalancesParams {
  house_gid: number;
  group_id: number;
  min_yuan?: number; // 最小余额(元)
  max_yuan?: number; // 最大余额(元)
  page?: number;
  size?: number;
}

/**
 * 查询圈子统计参数
 */
export interface GetGroupStatsParams {
  house_gid: number;
  group_id: number;
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
}

/**
 * ===== 超级管理员接口参数 =====
 */

/**
 * 查询店铺统计参数
 */
export interface GetHouseStatsParams {
  house_gid: number;
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
}

/**
 * ===== 响应数据类型 =====
 */

/**
 * 战绩记录
 */
export interface BattleRecord {
  id: number;
  house_gid: number;
  group_id: number;
  group_name: string;
  room_uid: number;
  kind_id: number;
  base_score: number;
  battle_at: string;
  players_json: string;
  player_id?: number;
  player_game_id?: number;
  score: number;
  fee: number;
  factor: number;
  player_balance: number;
  created_at: string;
}

/**
 * 成员余额信息
 */
export interface MemberBalance {
  member_id: number;
  game_id: number;
  game_name: string;
  group_id?: number;
  group_name: string;
  balance: number; // 余额(分)
  balance_yuan: number; // 余额(元)
  updated_at: string;
}

/**
 * 战绩统计
 */
export interface BattleStats {
  total_games: number;
  total_score: number;
  total_fee: number;
  avg_score: number;
  group_id?: number;
  group_name?: string;
}

/**
 * 圈子统计
 */
export interface GroupStats {
  group_id: number;
  group_name: string;
  total_games: number;
  total_score: number;
  total_fee: number;
  active_members: number;
}

/**
 * 店铺统计
 */
export interface HouseStats {
  house_gid: number;
  total_games: number;
  total_score: number;
  total_fee: number;
}

