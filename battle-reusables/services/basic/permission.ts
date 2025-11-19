import { request } from '@/utils/request';

export interface Permission {
  id: number;
  code: string;
  name: string;
  category: string;
  description: string;
  created_at: string;
  updated_at: string;
  is_deleted: boolean;
}

export interface CreatePermissionRequest {
  code: string;
  name: string;
  category: string;
  description?: string;
}

export interface UpdatePermissionRequest {
  id: number;
  name?: string;
  category?: string;
  description?: string;
}

export interface AssignPermissionsRequest {
  role_id: number;
  permission_ids: number[];
}

// 查询权限列表
export async function getPermissionList(category?: string) {
  return request<Permission[]>({
    url: '/basic/permission/list',
    method: 'GET',
    params: { category },
  });
}

// 查询所有权限
export async function getAllPermissions() {
  return request<Permission[]>({
    url: '/basic/permission/listAll',
    method: 'GET',
  });
}

// 创建权限
export async function createPermission(data: CreatePermissionRequest) {
  return request({
    url: '/basic/permission/create',
    method: 'POST',
    data,
  });
}

// 更新权限
export async function updatePermission(data: UpdatePermissionRequest) {
  return request({
    url: '/basic/permission/update',
    method: 'POST',
    data,
  });
}

// 删除权限
export async function deletePermission(id: number) {
  return request({
    url: '/basic/permission/delete',
    method: 'POST',
    data: { id },
  });
}

// 查询角色权限
export async function getRolePermissions(roleId: number) {
  return request<Permission[]>({
    url: '/basic/permission/role/permissions',
    method: 'GET',
    params: { role_id: roleId },
  });
}

// 为角色分配权限
export async function assignPermissionsToRole(data: AssignPermissionsRequest) {
  return request({
    url: '/basic/permission/role/assign',
    method: 'POST',
    data,
  });
}

// 从角色移除权限
export async function removePermissionsFromRole(data: AssignPermissionsRequest) {
  return request({
    url: '/basic/permission/role/remove',
    method: 'POST',
    data,
  });
}

