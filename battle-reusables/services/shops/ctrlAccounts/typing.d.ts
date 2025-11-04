declare namespace API {
  type ShopsCtrlAccountsCreateParams = {
    login_mode: 'account' | 'mobile';
    identifier: string;
    pwd_md5: string;
    status?: number;
    house_gid?: number;
  };

  type ShopsCtrlAccountsBindParams = {
    ctrl_id: number;
    house_gid: number;
    status?: number;
  };

  type ShopsCtrlAccountsUnbindParams = {
    ctrl_id: number;
    house_gid: number;
  };

  type ShopsCtrlAccountsListParams = {
    house_gid: number;
  };

  type ShopsCtrlAccountsListAllParams = {
    login_mode?: string;
    status?: number;
    keyword?: string;
    page?: number;
    size?: number;
  };

  type ShopsCtrlAccountsItem = {
    id?: number;
    house_gid?: number;
    login_mode?: string;
    identifier?: string;
    status?: number;
  };

  type ShopsCtrlAccountsHouse = {
    house_gid?: number;
    status?: number;
  };

  type ShopsCtrlAccountsAllItem = {
    id?: number;
    login_mode?: string;
    identifier?: string;
    status?: number;
    last_verify_at?: string;
    houses?: ShopsCtrlAccountsHouse[];
  };

  type ShopsCtrlAccountsListResult = ShopsCtrlAccountsItem[];

  type ShopsCtrlAccountsAllResult = ShopsCtrlAccountsAllItem[];
}

