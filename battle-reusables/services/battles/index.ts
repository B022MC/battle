import { post } from '@/utils/request';
import type { BattleRecord, BattleStats, ListMyBattlesParams, GetMyStatsParams } from './typing';

/**
 * 查询我的战绩
 */
export function listMyBattles(data: ListMyBattlesParams) {
  return post<{
    list: BattleRecord[];
    total: number;
    page: number;
    size: number;
  }>('/battles/my/list', data);
}

/**
 * 查询我的战绩统计
 */
export function getMyStats(data: GetMyStatsParams) {
  return post<BattleStats>('/battles/my/stats', data);
}

