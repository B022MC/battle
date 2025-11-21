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
}
