import { get, post } from '@/utils/request';

export function gameAccountVerify(data: API.GameAccountVerifyParams) {
  return post<API.GameAccountVerifyResult>('/game/accounts/verify', data);
}

export function gameAccountBind(data: API.GameAccountBindParams) {
  return post<API.GameAccountItem>('/game/accounts', data);
}

export function gameAccountMe() {
  return get<API.GameAccountItem | null>('/game/accounts/me');
}

export async function gameAccountDelete() {
  const response = await fetch('/game/accounts/me', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  return response.json();
}

