import { post } from '@/utils/request';
import type {
  ListMyBattlesParams,
  GetMyBalancesParams,
  GetMyStatsParams,
  ListGroupBattlesParams,
  ListGroupMemberBalancesParams,
  GetGroupStatsParams,
  GetHouseStatsParams,
  BattleRecord,
  MemberBalance,
  BattleStats,
  GroupStats,
  HouseStats,
} from './query-typing';

/**
 * 普通用户 - 查询我的战绩
 */
export function listMyBattles(data: ListMyBattlesParams) {
  return post<{
    list: BattleRecord[];
    total: number;
  }>('/battle-query/my/battles', data);
}

/**
 * 普通用户 - 查询我的余额
 */
export function getMyBalances(data: GetMyBalancesParams) {
  return post<{
    balances: MemberBalance[];
  }>('/battle-query/my/balances', data);
}

/**
 * 普通用户 - 查询我的统计
 */
export function getMyStats(data: GetMyStatsParams) {
  return post<BattleStats>('/battle-query/my/stats', data);
}

/**
 * 管理员 - 查询圈子战绩
 */
export function listGroupBattles(data: ListGroupBattlesParams) {
  return post<{
    list: BattleRecord[];
    total: number;
  }>('/battle-query/group/battles', data);
}

/**
 * 管理员 - 查询圈子成员余额
 */
export function listGroupMemberBalances(data: ListGroupMemberBalancesParams) {
  return post<{
    list: MemberBalance[];
    total: number;
  }>('/battle-query/group/balances', data);
}

/**
 * 管理员 - 查询圈子统计
 */
export function getGroupStats(data: GetGroupStatsParams) {
  return post<GroupStats>('/battle-query/group/stats', data);
}

/**
 * 超级管理员 - 查询店铺统计
 */
export function getHouseStats(data: GetHouseStatsParams) {
  return post<HouseStats>('/battle-query/house/stats', data);
}

