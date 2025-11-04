import { get } from '@/utils/request';

export function shopsHousesOptions() {
  return get<number[]>('/shops/houses/options');
}
