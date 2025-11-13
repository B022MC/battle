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
    state?: 'online' | 'offline' | 'error';
  };

  // Session detail types (for future backend implementation)
  type GameSessionDetail = {
    id: number;
    game_ctrl_account_id: number;
    user_id: number;
    house_gid: number;
    state: 'active' | 'inactive' | 'error';
    device_ip?: string;
    error_msg?: string;
    auto_sync_enabled: boolean;
    last_sync_at?: string;
    sync_status: 'idle' | 'syncing' | 'error';
    game_account_id?: number;
    created_at?: string;
    updated_at?: string;
  };

  // Sync log types (for future backend implementation)
  type GameSyncLog = {
    id: number;
    session_id: number;
    sync_type: 'battle_record' | 'member_list' | 'wallet_update' | 'room_list' | 'group_member';
    status: 'success' | 'failed' | 'partial';
    records_synced: number;
    error_message?: string;
    started_at: string;
    completed_at?: string;
  };

  type GameSyncLogListParams = {
    session_id: number;
    sync_type?: string;
    status?: string;
    page?: number;
    size?: number;
  };

  type GameSessionQueryParams = {
    house_gid?: number;
    ctrl_account_id?: number;
    state?: string;
  };
}

