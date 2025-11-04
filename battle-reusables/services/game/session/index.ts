import { post } from '@/utils/request';

export function gameSessionStart(data: API.GameSessionStartParams) {
  return post<API.GameSessionStateResult>('/game/accounts/sessionStart', data);
}

export function gameSessionStop(data: API.GameSessionStopParams) {
  return post<API.GameSessionStateResult>('/game/accounts/sessionStop', data);
}

