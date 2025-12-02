import { post } from '@/utils/request';

/**
 * 查询店铺费用配置
 */
export function getShopFees(data: API.GetShopFeesParams) {
  return post<API.ShopFeesResult>('/shops/fees/get', data);
}

/**
 * 设置店铺费用规则
 */
export function setShopFees(data: API.SetShopFeesParams) {
  return post<unknown>('/shops/fees/set', data);
}

/**
 * 设置分运开关
 */
export function setShareFee(data: API.SetShareFeeParams) {
  return post<unknown>('/shops/sharefee/set', data);
}

/**
 * 设置推送额度
 */
export function setPushCredit(data: API.SetPushCreditParams) {
  return post<unknown>('/shops/pushcredit/set', data);
}
