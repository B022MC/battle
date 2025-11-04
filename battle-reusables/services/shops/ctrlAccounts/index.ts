import { post } from '@/utils/request';

export function shopsCtrlAccountsCreate(data: API.ShopsCtrlAccountsCreateParams) {
  return post<API.ShopsCtrlAccountsItem>('/shops/ctrlAccounts', data);
}

export function shopsCtrlAccountsBind(data: API.ShopsCtrlAccountsBindParams) {
  return post<null>('/shops/ctrlAccounts/bind', data);
}

export async function shopsCtrlAccountsUnbind(data: API.ShopsCtrlAccountsUnbindParams) {
  const response = await fetch('/shops/ctrlAccounts/bind', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return response.json();
}

export function shopsCtrlAccountsList(data: API.ShopsCtrlAccountsListParams) {
  return post<API.ShopsCtrlAccountsListResult>('/shops/ctrlAccounts/list', data);
}

export function shopsCtrlAccountsListAll(data: API.ShopsCtrlAccountsListAllParams) {
  return post<API.ShopsCtrlAccountsAllResult>('/shops/ctrlAccounts/listAll', data);
}

