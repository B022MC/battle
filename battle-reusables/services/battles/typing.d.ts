/**
 * 查询我的战绩参数
 */
export interface ListMyBattlesParams {
  house_gid?: number; // 可选，如果不传则查询所有店铺
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
  page?: number;
  size?: number;
}

/**
 * 查询我的统计参数
 */
export interface GetMyStatsParams {
  house_gid?: number; // 可选，如果不传则查询所有店铺
  start_time?: number; // Unix timestamp
  end_time?: number; // Unix timestamp
}

/**
 * 战绩记录
 */
export interface BattleRecord {
  id: number;
  house_gid: number;
  group_id: number;
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
 * 战绩统计
 */
export interface BattleStats {
  total_games: number;
  total_score: number;
  total_fee: number;
}

