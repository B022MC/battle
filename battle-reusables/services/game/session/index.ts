import { post, get } from '@/utils/request';

/**
 * Start a session for a control account
 * Requires permission: game:ctrl:create
 */
export function gameSessionStart(data: API.GameSessionStartParams) {
  return post<API.GameSessionStateResult>('/game/accounts/sessionStart', data);
}

/**
 * Stop a session for a control account
 * Requires permission: game:ctrl:update
 */
export function gameSessionStop(data: API.GameSessionStopParams) {
  return post<API.GameSessionStateResult>('/game/accounts/sessionStop', data);
}

/**
 * Query session details by house_gid or ctrl_account_id
 * NOTE: Backend endpoint not yet implemented
 * TODO: Add backend endpoint /game/sessions/query
 */
export function gameSessionQuery(params: API.GameSessionQueryParams) {
  return post<API.GameSessionDetail[]>('/game/sessions/query', params);
}

/**
 * Get session detail by session ID
 * NOTE: Backend endpoint not yet implemented
 * TODO: Add backend endpoint /game/sessions/:id
 */
export function gameSessionDetail(sessionId: number) {
  return get<API.GameSessionDetail>(`/game/sessions/${sessionId}`);
}

/**
 * Query sync logs for a session
 * NOTE: Backend endpoint not yet implemented
 * TODO: Add backend endpoint /game/sessions/:id/logs
 */
export function gameSessionSyncLogs(params: API.GameSyncLogListParams) {
  return post<API.GameSyncLog[]>(`/game/sessions/${params.session_id}/logs`, params);
}

