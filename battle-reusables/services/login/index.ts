import { post } from '@/utils/request';

export async function loginRegister(data: API.LoginRegisterParams) {
  return post<API.LoginRegisterResult>('/login/register', data);
}

export async function loginUsername(data: API.LoginUsernameParams) {
  return post<API.LoginRegisterResult>('/login/username', data);
}
