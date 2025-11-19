// Type definitions for Game Control Account Management (Super Admin)

export interface CtrlAccountVO {
  id: number;
  house_gid?: number;
  login_mode: 'account' | 'mobile';
  identifier: string;
  status: number; // 0=disabled, 1=enabled
}

export interface CtrlAccountAllVO {
  id: number;
  login_mode: 'account' | 'mobile';
  identifier: string;
  status: number;
  last_verify_at?: string;
  houses: number[]; // List of bound house_gid
}

export interface CreateCtrlAccountRequest {
  login_mode: 'account' | 'mobile';
  identifier: string;
  pwd_md5: string; // 32-char uppercase MD5
  status?: number; // 0 or 1, default 1
  house_gid?: number; // Optional: if provided, will bind immediately after creation
}

export interface BindCtrlAccountRequest {
  ctrl_id: number;
  house_gid: number;
  status?: number; // 0 or 1, default 1
}

export interface UnbindCtrlAccountRequest {
  ctrl_id: number;
  house_gid: number;
}

export interface ListCtrlAccountsRequest {
  house_gid: number;
}

export interface ListAllCtrlAccountsRequest {
  login_mode?: 'account' | 'mobile';
  status?: number;
  keyword?: string;
  page?: number;
  size?: number;
}

export interface StartSessionRequest {
  id: number; // game_ctrl_account primary key
  house_gid: number;
}

export interface StopSessionRequest {
  id: number; // game_ctrl_account primary key
  house_gid: number;
}

export interface SessionStateResponse {
  state: 'online' | 'offline' | 'error';
}

export interface UpdateStatusRequest {
  ctrl_id: number;
  status: 0 | 1; // 0=disabled, 1=enabled
}

