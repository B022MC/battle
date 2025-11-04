declare namespace API {
  type GameAccountVerifyParams = {
    mode: 'account' | 'mobile';
    account: string;
    pwd_md5: string;
  };

  type GameAccountVerifyResult = {
    ok?: boolean;
  };

  type GameAccountBindParams = {
    mode: 'account' | 'mobile';
    account: string;
    pwd_md5: string;
    nickname?: string;
  };

  type GameAccountItem = {
    id?: number;
    account?: string;
    nickname?: string;
    is_default?: boolean;
    status?: number;
    login_mode?: string;
  };
}

