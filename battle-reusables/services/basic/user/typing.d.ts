declare namespace API {
  type BasicUserAddParams = {
    username: string;
    password?: string;
    nick_name?: string;
    avatar?: string;
    role_ids?: number[];
  };

  type BasicUserUpdateParams = BasicUserAddParams & {
    id: number;
  };

  type BasicUserDelParams = {
    id: number;
  };

  type BasicUserDelManyParams = {
    id: number[];
  };

  type BasicUserGetParams = {
    id: number;
  };

  type BasicUserListParams = {
    page_no?: number;
    page_size?: number;
    not_page?: boolean;
    keyword?: string;
  };

  type BasicUserItem = {
    id?: number;
    username?: string;
    nick_name?: string;
    avatar?: string;
    role_ids?: number[];
    created_at?: string;
    updated_at?: string;
  };

  type BasicUserList = {
    list?: BasicUserItem[];
    total?: number;
    page_no?: number;
    page_size?: number;
    not_page?: boolean;
  };

  type BasicUserRoles = {
    role_ids?: number[];
  };

  type BasicUserPerms = {
    perms?: string[];
  };
}

