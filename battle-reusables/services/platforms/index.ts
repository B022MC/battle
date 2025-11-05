import { get } from '@/utils/request';

export async function platformsList() {
  return get<API.platform[]>('/platforms/list');
}

export async function platformsPlazaConsts() {
  return get<API.PlazaConsts>('/platforms/plaza/consts');
}
