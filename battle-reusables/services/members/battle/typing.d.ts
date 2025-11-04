declare namespace API {
  type MembersBattleDetailsParams = {
    house_gid: number;
    group_id: number;
    period: 'today' | 'yesterday' | 'thisweek';
    game_id?: number;
  };

  type MembersBattlePlayer = {
    game_id?: number;
    score?: number;
  };

  type MembersBattleRecord = {
    room_id?: number;
    kind_id?: number;
    base_score?: number;
    time?: number;
    players?: MembersBattlePlayer[];
  };
}

