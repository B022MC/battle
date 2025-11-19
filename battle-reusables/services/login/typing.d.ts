declare namespace API {
  type LoginUsernameParams = {
    username: string;
    password: string;
  };

  type LoginRegisterParams = LoginUsernameParams & {
    nick_name?: string;
    game_account_mode?: string;
    game_account?: string;
    game_password?: string;
  };

  type UserInfo = {
    id?: number;
    username?: string;
    avatar?: string;
    nick_name?: string;
  };

  type LoginRegisterResult = {
    access_token?: string;
    expires_in?: number;
    platform?: string;
    role?: string; // 用户角色：super_admin, store_admin, user
    refresh_token?: string;
    user?: UserInfo;
  };
}
