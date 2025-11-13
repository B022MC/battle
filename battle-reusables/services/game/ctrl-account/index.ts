// Game Control Account Management API Services (Super Admin)
import { request } from '@/utils/request';
import type {
  CtrlAccountVO,
  CtrlAccountAllVO,
  CreateCtrlAccountRequest,
  BindCtrlAccountRequest,
  UnbindCtrlAccountRequest,
  ListCtrlAccountsRequest,
  ListAllCtrlAccountsRequest,
  StartSessionRequest,
  StopSessionRequest,
  SessionStateResponse,
} from './typing';

/**
 * Create or update a control account (without binding to store)
 * If house_gid is provided, will bind immediately after creation
 * Requires permission: game:ctrl:create
 */
export const createCtrlAccount = (data: CreateCtrlAccountRequest) => {
  return request<CtrlAccountVO>({
    url: '/shops/ctrlAccounts',
    method: 'POST',
    data,
  });
};

/**
 * Bind control account to store
 * Requires permission: game:ctrl:update
 */
export const bindCtrlAccount = (data: BindCtrlAccountRequest) => {
  return request<{ ok: boolean }>({
    url: '/shops/ctrlAccounts/bind',
    method: 'POST',
    data,
  });
};

/**
 * Unbind control account from store
 * Requires permission: game:ctrl:update
 */
export const unbindCtrlAccount = (data: UnbindCtrlAccountRequest) => {
  return request<{ ok: boolean }>({
    url: '/shops/ctrlAccounts/bind',
    method: 'DELETE',
    data,
  });
};

/**
 * List control accounts by store (house_gid)
 * Requires permission: game:ctrl:view
 */
export const listCtrlAccountsByHouse = (data: ListCtrlAccountsRequest) => {
  return request<CtrlAccountVO[]>({
    url: '/shops/ctrlAccounts/list',
    method: 'POST',
    data,
  });
};

/**
 * List all control accounts with filters and pagination
 * Requires permission: game:ctrl:view
 */
export const listAllCtrlAccounts = (data: ListAllCtrlAccountsRequest = {}) => {
  return request<CtrlAccountAllVO[]>({
    url: '/shops/ctrlAccounts/listAll',
    method: 'POST',
    data,
  });
};

/**
 * Get available house options (distinct house_gid from bindings)
 * Requires authentication
 */
export const getHouseOptions = () => {
  return request<number[]>({
    url: '/shops/houses/options',
    method: 'GET',
  });
};

/**
 * Start session for control account
 * Requires permission: game:ctrl:create
 */
export const startSession = (data: StartSessionRequest) => {
  return request<SessionStateResponse>({
    url: '/game/accounts/sessionStart',
    method: 'POST',
    data,
  });
};

/**
 * Stop session for control account
 * Requires permission: game:ctrl:update
 */
export const stopSession = (data: StopSessionRequest) => {
  return request<SessionStateResponse>({
    url: '/game/accounts/sessionStop',
    method: 'POST',
    data,
  });
};

