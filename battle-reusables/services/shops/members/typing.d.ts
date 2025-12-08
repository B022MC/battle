declare namespace API {
  type ShopsMembersListParams = { house_gid: number };

  type ShopsMembersItemParams = ShopsMembersListParams & { member_id: number };

  type ShopsMemberItem = {
    user_id?: number;
    member_id?: number;
    game_id?: number;
    nick_name?: string;
    member_type?: number;
    user_status?: number;
    group_id?: number;
    group_name?: string;
    // 平台用户关联信息
    platform_user?: UserInfo;
    game_account_id?: number;
    is_bind_platform?: boolean;
    // 拉圈功能字段
    game_player_id?: string; // 游戏玩家ID（用于拉圈）
    current_group_id?: number; // 当前所在圈子ID
    current_group_name?: string; // 当前所在圈子名称
    // 置顶功能字段
    is_pinned?: boolean; // 是否置顶
    pin_order?: number; // 置顶排序（数字越小越靠前）
    // 备注功能字段
    remark?: string; // 管理员备注
    // 禁用功能字段
    forbid?: boolean; // 是否禁用
    // 余额字段（来自 game_member 表）
    balance?: number; // 余额（单位：分）
    credit?: number; // 信用额度（单位：分）
  };

  type ShopsMembersList = { items?: ShopsMemberItem[] };

  type ShopsMembersDiamond = { triggered?: boolean };

  type ShopsMembersRulesVipParams = {
    house_gid: number;
    member_id: number;
    vip: boolean;
  };

  type ShopsMembersRulesMultiParams = {
    house_gid: number;
    member_id: number;
    allow: boolean;
  };

  type ShopsMembersRulesTempReleaseParams = {
    house_gid: number;
    member_id: number;
    limit: number;
  };

  type ShopsMembersPinParams = {
    house_gid: number;
    game_player_id: string;
    pin_order?: number; // 置顶顺序（可选，默认为0）
  };

  type ShopsMembersUnpinParams = {
    house_gid: number;
    game_player_id: string;
  };

  type ShopsMembersUpdateRemarkParams = {
    house_gid: number;
    game_player_id: string;
    remark: string; // 备注内容
  };
}
