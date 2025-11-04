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
