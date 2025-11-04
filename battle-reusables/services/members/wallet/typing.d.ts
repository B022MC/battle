declare namespace API {
  type MembersWalletGetParams = {
    house_gid: number;
    member_id: number;
  };

  type MembersWalletListParams = {
    house_gid: number;
    min_balance?: number;
    max_balance?: number;
    has_custom_limit?: boolean;
    page?: number;
    page_size?: number;
  };

  type MembersWalletItem = {
    member_id?: number;
    house_gid?: number;
    balance?: number;
    forbid?: boolean;
    limit_min?: number;
  };

  type MembersWalletListResult = {
    list?: MembersWalletItem[];
    total?: number;
    page?: number;
    page_size?: number;
  };
}

