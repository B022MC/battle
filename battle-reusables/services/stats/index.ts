import { post, get } from '@/utils/request';

export type StatsPath = '/stats/today' | '/stats/yesterday' | '/stats/week' | '/stats/lastweek';

export async function statsByPath(path: StatsPath, data: API.StatsParams) {
  return post<API.StatsResult>(path, data);
}

export type StatsMemberPath =
  | '/stats/member/today'
  | '/stats/member/yesterday'
  | '/stats/member/thisweek'
  | '/stats/member/lastweek';

export async function statsMemberByPath(path: StatsMemberPath, data: API.StatsMemberParams) {
  return post<API.StatsResult>(path, data);
}

export async function statsActiveByHouse() {
  return get<API.ActiveByHouseItem[]>('/stats/sessions/activeByHouse');
}