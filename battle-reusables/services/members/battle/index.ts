import { post } from '@/utils/request';

export function membersBattleDetails(data: API.MembersBattleDetailsParams) {
  return post<API.MembersBattleRecord[]>('/members/battle/details', data);
}

export function membersBattleExportHtml(data: API.MembersBattleDetailsParams) {
  return post<Blob>('/members/battle/export/html', data);
}

