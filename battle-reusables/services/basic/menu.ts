import { request } from '@/utils/request';

export interface Menu {
  id: number;
  parent_id: number;
  menu_type: number;
  title: string;
  name: string;
  path: string;
  component: string;
  rank?: string;
  redirect: string;
  icon: string;
  extra_icon: string;
  enter_transition: string;
  leave_transition: string;
  active_path: string;
  auths: string;
  frame_src: string;
  frame_loading: boolean;
  keep_alive: boolean;
  hidden_tag: boolean;
  fixed_tag: boolean;
  show_link: boolean;
  show_parent: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  is_del: number;
  children?: Menu[]; // 子菜单
}

export interface CreateMenuRequest {
  parent_id: number;
  menu_type: number;
  title: string;
  name: string;
  path: string;
  component: string;
  rank?: string;
  icon?: string;
  auths?: string;
  show_link?: boolean;
  show_parent?: boolean;
}

export interface UpdateMenuRequest {
  id: number;
  parent_id?: number;
  menu_type?: number;
  title?: string;
  name?: string;
  path?: string;
  component?: string;
  rank?: string;
  icon?: string;
  auths?: string;
  show_link?: boolean;
  show_parent?: boolean;
}

// 查询所有菜单
export async function getAllMenus() {
  return request<Menu[]>({
    url: '/basic/baseMenu/getOption',
    method: 'GET',
    params: {
      page: 1,
      page_size: 1000, // 设置较大值以获取所有菜单
    },
  });
}

// 查询菜单树
export async function getMenuTree(keyword?: string) {
  const params: any = {};
  // 只有当keyword有值时才添加到参数中
  if (keyword) {
    params.keyword = keyword;
  }
  return request<Menu[]>({
    url: '/basic/baseMenu/getTree',
    method: 'GET',
    params,
  });
}

// 查询用户菜单树
export async function getUserMenuTree() {
  return request<Menu[]>({
    url: '/basic/menu/me/tree',
    method: 'GET',
  });
}

// 查询单个菜单
export async function getMenu(id: number) {
  return request<Menu>({
    url: '/basic/baseMenu/getOne',
    method: 'GET',
    params: { id },
  });
}

// 创建菜单
export async function createMenu(data: CreateMenuRequest) {
  return request({
    url: '/basic/baseMenu/addOne',
    method: 'POST',
    data,
  });
}

// 更新菜单
export async function updateMenu(data: UpdateMenuRequest) {
  return request({
    url: '/basic/baseMenu/updateOne',
    method: 'POST',
    data,
  });
}

// 删除菜单
export async function deleteMenu(id: number) {
  return request({
    url: '/basic/baseMenu/delOne',
    method: 'GET',
    params: { id },
  });
}

