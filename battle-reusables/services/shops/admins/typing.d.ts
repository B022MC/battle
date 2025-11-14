declare namespace API {
  type ShopsAdminsAssignParams = {
    house_gid: number;
    user_id: number;
    role?: string;
  };

  type ShopsAdminsRevokeParams = {
    house_gid: number;
    user_id: number;
  };

  type ShopsAdminsListParams = {
    house_gid: number;
  };

  type ShopsAdminsItem = {
    id?: number;
    house_gid?: number;
    user_id?: number;
    role?: string;
    nick_name?: string;
  };

  type ShopsAdminsListResult = ShopsAdminsItem[];

  type ShopAdminInfo = ShopsAdminsItem;
}

