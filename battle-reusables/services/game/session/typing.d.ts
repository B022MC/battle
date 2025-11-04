declare namespace API {
  type GameSessionStartParams = {
    id: number;
    house_gid: number;
  };

  type GameSessionStopParams = {
    id: number;
    house_gid: number;
  };

  type GameSessionStateResult = {
    state?: 'online' | 'offline';
  };
}

