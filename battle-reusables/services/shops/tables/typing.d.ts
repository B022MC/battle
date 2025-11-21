declare namespace API {
  type ShopsTablesListParams = { house_gid: number };

  type ShopsTablesItemParams = ShopsTablesListParams & { mapped_num: number };

  type ShopsTablesDismissParams = ShopsTablesItemParams & { kind_id?: number };

  type ShopsTableItem = {
    table_id?: number;
    group_id?: number;
    mapped_num?: number;
    kind_id?: number;
    base_score?: number;
  };

  type ShopsTablesList = { items?: ShopsTableItem[] };

  type ShopsTablesDetail = { table?: ShopsTableItem; triggered?: boolean };

  type ShopsTablesCheck = ShopsTablesDetail & { exists_in_cache?: boolean };

  type UserInfo = {
    id: number;
    username: string;
    nick_name: string;
    avatar: string;
    role: string;
    introduction: string;
    created_at: string;
    updated_at: string;
  };

  type ShopsMemberItem = {
    user_id: number;
    user_status: number;
    game_id: number;
    member_id: number;
    member_type: number;
    nick_name: string;
    group_id: number;
    group_name?: string;
    // 平台用户关联信息
    platform_user?: UserInfo;
    game_account_id?: number;
    is_bind_platform: boolean;
    // 拉圈功能字段
    game_player_id?: string; // 游戏玩家ID（用于拉圈）
    current_group_id?: number; // 当前所在圈子ID
    current_group_name?: string; // 当前所在圈子名称
  };

  type ShopsMembersList = { items: ShopsMemberItem[] };
}
