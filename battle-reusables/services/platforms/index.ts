import { get } from '@/utils/request';

export async function platformsList() {
  return get<API.platform[]>('/platforms/list');
}
