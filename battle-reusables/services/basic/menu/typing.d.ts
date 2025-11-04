declare namespace API {
  type BasicMenuAddParams = {
    name: string;
    path?: string;
    component?: string;
    icon?: string;
    parent_id?: number;
    rank?: string;
    menu_type?: string;
    perm_code?: string;
  };

  type BasicMenuUpdateParams = BasicMenuAddParams & {
    id: number;
  };

  type BasicMenuGetParams = {
    id: number;
  };

  type BasicMenuPageParams = {
    page?: number;
    page_size?: number;
  };

  type BasicMenuTreeParams = {
    keyword?: string;
  };

  type BasicMenuDelManyParams = {
    ids: number[];
  };

  type BasicMenuItem = {
    id?: number;
    name?: string;
    path?: string;
    component?: string;
    icon?: string;
    parent_id?: number;
    rank?: string;
    menu_type?: string;
    perm_code?: string;
    created_at?: string;
    updated_at?: string;
  };

  type BasicMenuTree = {
    id?: number;
    name?: string;
    path?: string;
    icon?: string;
    children?: BasicMenuTree[];
  };

  type BasicMenuList = {
    list?: BasicMenuItem[];
    total?: number;
    page_no?: number;
    page_size?: number;
  };

  type BasicMenuSaveTreeParams = {
    menu_tree: BasicMenuTree[];
  };
}

