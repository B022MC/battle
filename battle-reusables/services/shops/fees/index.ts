import { post } from '@/utils/request';

export function shopsFeesGet(data: API.ShopsFeesGetParams) {
  return post<API.ShopsFeesItem>('/shops/fees/get', data);
}

export function shopsFeesSet(data: API.ShopsFeesSetParams) {
  return post<null>('/shops/fees/set', data);
}

export function shopsShareFeeSet(data: API.ShopsShareFeeSetParams) {
  return post<null>('/shops/sharefee/set', data);
}

export function shopsPushCreditGet(data: API.ShopsFeesGetParams) {
  return post<API.ShopsFeesItem>('/shops/pushcredit/get', data);
}

export function shopsPushCreditSet(data: API.ShopsPushCreditSetParams) {
  return post<null>('/shops/pushcredit/set', data);
}

export function shopsFeesSettleInsert(data: API.ShopsFeesSettleInsertParams) {
  return post<API.ShopsFeesSettleInsertResult>('/shops/fees/settle/insert', data);
}

export function shopsFeesSettleSum(data: API.ShopsFeesSettleSumParams) {
  return post<API.ShopsFeesSettleSumResult>('/shops/fees/settle/sum', data);
}

