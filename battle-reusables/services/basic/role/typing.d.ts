declare namespace API {
  type BasicRoleListParams = {
    keyword?: string;
    enable?: boolean;
    page_no?: number;
    page_size?: number;
  };

  type BasicRoleGetParams = {
    id: number;
  };

  type BasicRoleItem = {
    id?: number;
    code?: string;
    name?: string;
    parent_id?: number;
    remark?: string;
    created_at?: string;
    created_user?: number;
    updated_at?: string;
    updated_user?: number;
    first_letter?: string;
    pinyin_code?: string;
    enable?: boolean;
    is_deleted?: boolean;
  };

  type BasicRoleListResult = {
    list?: BasicRoleItem[];
    page_no?: number;
    page_size?: number;
    total?: number;
  };

  type BasicRoleAllResult = {
    list?: BasicRoleItem[];
  };
}

