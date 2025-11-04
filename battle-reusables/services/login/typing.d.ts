declare namespace API {
  type LoginUsernameParams = {
    username: string;
    password: string;
  };

  type LoginRegisterParams = LoginUsernameParams & {
    nick_name?: string;
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
    refresh_token?: string;
    user?: UserInfo;
  };
}
