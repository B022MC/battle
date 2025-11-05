import { post, request } from '@/utils/request';

export function shopsCtrlAccountsCreate(data: API.ShopsCtrlAccountsCreateParams) {
  return post<API.ShopsCtrlAccountsItem>('/shops/ctrlAccounts', data);
}

export function shopsCtrlAccountsBind(data: API.ShopsCtrlAccountsBindParams) {
  return post<null>('/shops/ctrlAccounts/bind', data);
}

export function shopsCtrlAccountsUnbind(data: API.ShopsCtrlAccountsUnbindParams) {
  return request<null>('/shops/ctrlAccounts/bind', { method: 'DELETE', data });
}

export function shopsCtrlAccountsList(data: API.ShopsCtrlAccountsListParams) {
  return post<API.ShopsCtrlAccountsListResult>('/shops/ctrlAccounts/list', data);
}

export function shopsCtrlAccountsListAll(data: API.ShopsCtrlAccountsListAllParams) {
  return post<API.ShopsCtrlAccountsAllResult>('/shops/ctrlAccounts/listAll', data);
}

