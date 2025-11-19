import { request } from '@/utils/request';

export interface Role {
  id: number;
  code: string;
  name: string;
  parent_id: number;
  remark?: string;
  created_at: string;
  created_user?: number;
  updated_at?: string;
  updated_user?: number;
  first_letter: string;
  pinyin_code: string;
  enable: boolean;
  is_deleted: boolean;
}

export interface CreateRoleRequest {
  code: string;
  name: string;
  remark?: string;
}

export interface UpdateRoleRequest {
  id: number;
  name?: string;
  remark?: string;
  enable?: boolean;
}

export interface AssignMenusRequest {
  role_id: number;
  menu_ids: number[];
}

export interface RoleListResponse {
  list: Role[];
  page_no: number;
  page_size: number;
  total: number;
}

// 查询角色列表（分页）
export async function getRoleList(params?: {
  keyword?: string;
  enable?: boolean;
  page_no?: number;
  page_size?: number;
}) {
  return request<RoleListResponse>({
    url: '/basic/role/list',
    method: 'GET',
    params,
  });
}

// 查询单个角色
export async function getRole(id: number) {
  return request<Role>({
    url: '/basic/role/getOne',
    method: 'GET',
    params: { id },
  });
}

// 查询所有角色（不分页）
export async function getAllRoles() {
  return request<{ list: Role[] }>({
    url: '/basic/role/all',
    method: 'GET',
  });
}

// 创建角色
export async function createRole(data: CreateRoleRequest) {
  return request({
    url: '/basic/role/create',
    method: 'POST',
    data,
  });
}

// 更新角色
export async function updateRole(data: UpdateRoleRequest) {
  return request({
    url: '/basic/role/update',
    method: 'POST',
    data,
  });
}

// 删除角色
export async function deleteRole(id: number) {
  return request({
    url: '/basic/role/delete',
    method: 'POST',
    data: { id },
  });
}

// 查询角色菜单
export async function getRoleMenus(roleId: number) {
  return request<{ menu_ids: number[] }>({
    url: '/basic/role/menus',
    method: 'GET',
    params: { role_id: roleId },
  });
}

// 为角色分配菜单
export async function assignMenusToRole(data: AssignMenusRequest) {
  return request({
    url: '/basic/role/menus/assign',
    method: 'POST',
    data,
  });
}

