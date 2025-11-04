declare namespace API {
  type ShopsTablesListParams = { house_gid: number };

  type ShopsTablesItemParams = ShopsTablesListParams & { mapped_num: number };

  type ShopsTablesDismissParams = ShopsTablesItemParams & { kind_id?: number };

  type ShopsTableItem = {
    table_id?: number;
    group_id?: number;
    mapped_num?: number;
    kind_id?: number;
    base_score?: number;
  };

  type ShopsTablesList = { items?: ShopsTableItem[] };

  type ShopsTablesDetail = { table?: ShopsTableItem; triggered?: boolean };

  type ShopsTablesCheck = ShopsTablesDetail & { exists_in_cache?: boolean };
}
