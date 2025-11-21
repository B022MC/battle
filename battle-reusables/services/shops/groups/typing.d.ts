declare namespace API {
  type ShopsGroupsBaseParams = {
    house_gid: number;
  };

  type ShopsGroupsForbidParams = ShopsGroupsBaseParams & {
    key: string;
    member_ids?: number[];
  };

  type ShopsGroupsBindParams = ShopsGroupsBaseParams & {
    message_id: number;
  };

  // 新的圈子系统类型定义
  type ShopGroup = {
    id: number;
    house_gid: number;
    group_name: string;
    admin_user_id: number;
    description: string;
    member_count?: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };

  type CreateGroupParams = {
    house_gid: number;
    group_name: string;
    description?: string;
  };

  type GetMyGroupParams = {
    house_gid: number;
  };

  type ListGroupsByHouseParams = {
    house_gid: number;
  };

  type RemoveMemberParams = {
    group_id: number;
    user_id: number;
  };

  type ListGroupMembersParams = {
    group_id: number;
    page?: number;
    size?: number;
  };

  type GroupMembersResponse = {
    items: BasicUser[];
    total: number;
  };

  type BasicUser = {
    id: number;
    username: string;
    nick_name: string;
    role?: string; // 用户角色：super_admin, store_admin, user
    phone?: string;
    email?: string;
    created_at: string;
  };

  type ListAllUsersParams = {
    page?: number;
    size?: number;
    keyword?: string;
  };

  type AllUsersResponse = {
    items: BasicUser[];
    total: number;
  };

  type GetUserParams = {
    user_id: number;
  };

  type ListShopAdminsParams = {
    house_gid: number;
  };

  type GetGroupOptionsParams = {
    house_gid: number;
  };

  type GroupOption = {
    id: number;
    name: string;
  };

  // 拉圈请求参数
  type PullMembersToGroupParams = {
    house_gid: number;          // 店铺ID
    group_id: number;           // 圈子ID
    game_player_ids: string[];  // game_player_id 列表
  };

  // 拉圈响应
  type PullMembersToGroupResponse = {
    message: string;
    count: number;
  };

  // 踢出圈子请求参数
  type RemoveFromGroupParams = {
    house_gid: number;          // 店铺ID
    game_player_ids: string[];  // game_player_id 列表
  };

  // 踢出圈子响应
  type RemoveFromGroupResponse = {
    message: string;
    count: number;
  };
}

// 兼容旧命名空间
declare namespace Shops {
  namespace Groups {
    type PullMembersToGroupParams = API.PullMembersToGroupParams;
    type PullMembersToGroupResponse = API.PullMembersToGroupResponse;
    type RemoveFromGroupParams = API.RemoveFromGroupParams;
    type RemoveFromGroupResponse = API.RemoveFromGroupResponse;
  }
}
