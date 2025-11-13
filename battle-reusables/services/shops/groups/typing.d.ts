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

  type AddMembersParams = {
    group_id: number;
    user_ids: number[];
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
}
