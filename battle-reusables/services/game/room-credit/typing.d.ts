declare namespace API {
  /**
   * 设置房间额度限制 - 请求参数
   */
  type SetRoomCreditLimitParams = {
    house_gid: number;
    group_name?: string; // 圈子名称，空表示全局
    game_kind?: number; // 游戏类型，0表示默认
    base_score?: number; // 底分，0表示默认
    credit_limit: number; // 额度限制（分）
  };

  /**
   * 查询房间额度限制 - 请求参数
   */
  type GetRoomCreditLimitParams = {
    house_gid: number;
    group_name?: string; // 圈子名称
    game_kind?: number; // 游戏类型
    base_score?: number; // 底分
  };

  /**
   * 列出房间额度限制 - 请求参数
   */
  type ListRoomCreditLimitParams = {
    house_gid: number;
    group_name?: string; // 圈子名称，空表示查询所有
  };

  /**
   * 删除房间额度限制 - 请求参数
   */
  type DeleteRoomCreditLimitParams = {
    house_gid: number;
    group_name?: string; // 圈子名称
    game_kind?: number; // 游戏类型
    base_score?: number; // 底分
  };

  /**
   * 检查玩家是否满足房间额度要求 - 请求参数
   */
  type CheckPlayerCreditParams = {
    house_gid: number;
    game_id: number; // 玩家游戏ID
    group_name?: string; // 圈子名称
    game_kind: number; // 游戏类型
    base_score: number; // 底分
  };

  /**
   * 房间额度限制项
   */
  type RoomCreditLimitItem = {
    id: number;
    house_gid: number;
    group_name: string;
    game_kind: number;
    game_kind_name?: string; // 游戏类型名称（如"红中"）
    base_score: number;
    credit_limit: number; // 单位：分
    credit_yuan: number; // 单位：元
    created_at: string;
    updated_at: string;
    updated_by: number;
  };

  /**
   * 房间额度限制列表 - 响应数据
   */
  type RoomCreditLimitListResult = {
    total: number;
    items: RoomCreditLimitItem[];
  };

  /**
   * 检查玩家额度 - 响应数据
   */
  type CheckPlayerCreditResult = {
    can_enter: boolean; // 是否可以进入
    player_balance: number; // 玩家余额（分）
    required_credit: number; // 需要的额度（分）
    player_credit: number; // 玩家个人额度调整（分）
    effective_credit: number; // 有效额度要求（分）
    balance_yuan: number; // 余额（元）
    required_yuan: number; // 需要的额度（元）
  };
}

