declare namespace API {
  type platform = {
    platform?: string;
    name?: string;
    db_name?: string;
    created_at?: string;
  };

  type LabeledInt = { key: string; value: number; label: string };
  type LabeledStr = { key: string; value: string; label: string };

  type PlazaConsts = {
    modes: LabeledInt[];
    scenes: LabeledInt[];
    user_status: LabeledInt[];
    member_types: LabeledInt[];
    game_genre: LabeledInt[];
    game_kinds: LabeledInt[];
    table_genre: LabeledInt[];
    system_message_types: LabeledInt[];
    versions: Record<string, number | string>;
    urls: LabeledStr[];
  };
}
