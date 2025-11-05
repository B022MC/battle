declare namespace API {
  type ShopsApplicationsListParams = {
    house_gid: number;
    type?: number; // 1=管理员,2=入圈
    admin_user_id?: number; // 圈主ID（仅对入圈申请有意义）
  };

  type ShopsApplicationsItem = {
    id?: number;
    status?: number;
    applier_id?: number;
    applier_gid?: number;
    applier_name?: string;
    house_gid?: number;
    type?: number;
    admin_user_id?: number;
    created_at?: number;
  };

  type ShopsApplicationsListResult = {
    items?: ShopsApplicationsItem[];
  };

  type ShopsApplicationsApplyAdminParams = {
    house_gid: number;
    note?: string;
  };

  type ShopsApplicationsApplyJoinParams = {
    house_gid: number;
    admin_user_id: number;
    note?: string;
  };

  type ShopsApplicationsApplyResult = {
    ok?: boolean;
  };

  type ShopsApplicationsHistoryParams = {
    house_gid: number;
    type?: number;   // 1=管理员,2=入圈
    status?: number; // 0待审,1通过,2拒绝
    start_at?: string; // RFC3339
    end_at?: string;   // RFC3339
  };
}

