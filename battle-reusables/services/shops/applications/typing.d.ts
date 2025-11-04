declare namespace API {
  type ShopsApplicationsListParams = {
    house_gid: number;
  };

  type ShopsApplicationsItem = {
    id?: number;
    status?: number;
    applier_id?: number;
    applier_gid?: number;
    applier_name?: string;
    house_gid?: number;
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
}

