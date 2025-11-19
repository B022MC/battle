/**
 * 游戏内申请信息
 */
export interface GameApplication {
  message_id: number;      // 游戏消息ID
  house_gid: number;       // 店铺游戏ID
  applier_gid: number;     // 申请人游戏ID
  applier_gname: string;   // 申请人游戏昵称
  applied_at: number;      // 申请时间戳（秒）
}

/**
 * 查询游戏内申请列表参数
 */
export interface GameApplicationListParams {
  house_gid: number;       // 店铺游戏ID
}

/**
 * 处理游戏内申请参数（通过/拒绝）
 */
export interface GameApplicationRespondParams {
  house_gid: number;       // 店铺游戏ID
  message_id: number;      // 游戏消息ID
}
