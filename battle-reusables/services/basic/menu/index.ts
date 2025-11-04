import { get, post } from '@/utils/request';

export function basicMenuAddOne(data: API.BasicMenuAddParams) {
  return post<API.BasicMenuItem>('/basic/baseMenu/addOne', data);
}

export function basicMenuUpdateOne(data: API.BasicMenuUpdateParams) {
  return post<null>('/basic/baseMenu/updateOne', data);
}

export function basicMenuGetOne(data: API.BasicMenuGetParams) {
  return get<API.BasicMenuItem>('/basic/baseMenu/getOne', data);
}

export function basicMenuGetPage(data: API.BasicMenuPageParams) {
  return get<API.BasicMenuList>('/basic/baseMenu/getPage', data);
}

export function basicMenuGetOption(data: API.BasicMenuPageParams) {
  return get<API.BasicMenuList>('/basic/baseMenu/getOption', data);
}

export function basicMenuGetTree(data?: API.BasicMenuTreeParams) {
  return get<API.BasicMenuTree>('/basic/baseMenu/getTree', data);
}

export function basicMenuSaveTree(data: API.BasicMenuSaveTreeParams) {
  return post<API.BasicMenuTree>('/basic/baseMenu/saveTree', data);
}

export function basicMenuDelOne(data: API.BasicMenuGetParams) {
  return get<null>('/basic/baseMenu/delOne', data);
}

export function basicMenuDelMany(data: API.BasicMenuDelManyParams) {
  return post<null>('/basic/baseMenu/delMany', data);
}

