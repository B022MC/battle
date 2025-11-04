import { post } from '@/utils/request';

export function applicationsApprove(data: API.ApplicationsDecideParams) {
  return post<API.ApplicationsDecideResult>('/applications/approve', data);
}

export function applicationsReject(data: API.ApplicationsDecideParams) {
  return post<API.ApplicationsDecideResult>('/applications/reject', data);
}

