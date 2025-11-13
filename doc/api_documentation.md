# API Documentation - User Management System

## Base URL
```
Production: https://api.battle-tiles.com
Staging: https://staging-api.battle-tiles.com
Development: http://localhost:8080
```

## Authentication
All authenticated endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

---

## Table of Contents
1. [Authentication APIs](#authentication-apis)
2. [User Management APIs](#user-management-apis)
3. [Game Account Management APIs](#game-account-management-apis)
4. [Store Management APIs](#store-management-apis)
5. [Game Record APIs](#game-record-apis)

---

## Authentication APIs

### 1. User Registration

**Endpoint:** `POST /api/login/register`

**Description:** Register a new user with game account binding

**Authentication:** None (Public)

**Request Body:**
```json
{
  "username": "string",          // Required, 3-50 chars, alphanumeric + underscore
  "password": "string",          // Required, 6-30 chars
  "game_login_mode": "string",   // Required, "account" or "mobile"
  "game_account": "string",      // Required, game account or mobile number
  "game_password_md5": "string"  // Required, MD5 hash of game password
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "nick_name": "Game Nickname",  // From game profile
      "avatar": "https://...",        // From game profile
      "roles": [1],
      "perms": ["user:read"]
    }
  }
}
```

**Error Responses:**
```json
// Username already exists
{
  "code": 400,
  "msg": "username already exists",
  "data": null
}

// Invalid game account
{
  "code": 400,
  "msg": "invalid game account credentials",
  "data": null
}

// Validation failed
{
  "code": 400,
  "msg": "validation failed: password must be at least 6 characters",
  "data": null
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/login/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456",
    "game_login_mode": "account",
    "game_account": "gameuser123",
    "game_password_md5": "5f4dcc3b5aa765d61d8327deb882cf99"
  }'
```

---

### 2. User Login

**Endpoint:** `POST /api/login/username`

**Description:** Login with username and password

**Authentication:** None (Public)

**Request Body:**
```json
{
  "username": "string",  // Required
  "password": "string"   // Required
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "nick_name": "Game Nickname",
      "avatar": "https://...",
      "roles": [1, 2],
      "perms": ["user:read", "user:write"]
    }
  }
}
```

**Error Responses:**
```json
// Invalid credentials
{
  "code": 401,
  "msg": "invalid username or password",
  "data": null
}
```

---

## User Management APIs

### 3. Get Current User Info

**Endpoint:** `GET /api/basic/user/me`

**Description:** Get current logged-in user information

**Authentication:** Required

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "nick_name": "Game Nickname",
    "avatar": "https://...",
    "introduction": "User bio",
    "last_login_at": "2025-11-11T10:00:00Z"
  }
}
```

---

### 4. Get Current User Roles

**Endpoint:** `GET /api/basic/user/me/roles`

**Description:** Get current user's role IDs

**Authentication:** Required

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "role_ids": [1, 2]
  }
}
```

---

### 5. Get Current User Permissions

**Endpoint:** `GET /api/basic/user/me/perms`

**Description:** Get current user's permission codes

**Authentication:** Required

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "perms": [
      "user:read",
      "user:write",
      "game:ctrl:create",
      "game:ctrl:update"
    ]
  }
}
```

---

### 6. Change Password

**Endpoint:** `POST /api/basic/user/changePassword`

**Description:** Change current user's password

**Authentication:** Required

**Request Body:**
```json
{
  "old_password": "string",  // Required
  "new_password": "string"   // Required, 6-30 chars
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

**Error Responses:**
```json
// Incorrect old password
{
  "code": 400,
  "msg": "old password incorrect",
  "data": null
}
```

---

## Game Account Management APIs

### 7. Verify Game Account

**Endpoint:** `POST /api/game/accounts/verify`

**Description:** Verify game account credentials without creating binding

**Authentication:** Required

**Request Body:**
```json
{
  "mode": "string",      // Required, "account" or "mobile"
  "account": "string",   // Required
  "pwd_md5": "string"    // Required, MD5 hash
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "ok": true
  }
}
```

**Error Responses:**
```json
// Invalid credentials
{
  "code": 400,
  "msg": "invalid game account credentials",
  "data": null
}
```

---

### 8. Bind Game Account (User)

**Endpoint:** `POST /api/game/accounts`

**Description:** Bind a game account to current user (limited to 1 for regular users)

**Authentication:** Required

**Request Body:**
```json
{
  "mode": "string",      // Required, "account" or "mobile"
  "account": "string",   // Required
  "pwd_md5": "string"    // Required, MD5 hash
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "account": "gameuser123",
    "nickname": "Game Nickname",
    "is_default": true,
    "status": 1
  }
}
```

**Error Responses:**
```json
// Already have game account
{
  "code": 409,
  "msg": "you have already bound a game account",
  "data": null
}
```

---

### 9. Get My Game Account

**Endpoint:** `GET /api/game/accounts/me`

**Description:** Get current user's game account

**Authentication:** Required

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "account": "gameuser123",
    "nickname": "Game Nickname",
    "is_default": true,
    "status": 1,
    "last_login_at": "2025-11-11T10:00:00Z"
  }
}
```

---

### 10. Unbind My Game Account

**Endpoint:** `DELETE /api/game/accounts/me`

**Description:** Unbind current user's game account

**Authentication:** Required

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

---

## Super Admin APIs

### 11. Bind Game Account to Store (Super Admin)

**Endpoint:** `POST /api/game/admin/ctrl-accounts/bind`

**Description:** Bind a game account to a store (super admin only)

**Authentication:** Required (Super Admin)

**Permissions:** `game:ctrl:create`

**Request Body:**
```json
{
  "login_mode": 1,        // Required, 1=account, 2=mobile
  "identifier": "string", // Required, game account or mobile
  "pwd_md5": "string",    // Required, MD5 hash
  "house_gid": 123        // Required, store ID
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "game_user_id": "12345",
    "game_id": "67890",
    "house_gid": 123,
    "session_state": "online",      // Auto-started
    "sync_status": "in_progress"    // Auto-sync triggered
  }
}
```

**Error Responses:**
```json
// Account already bound to another store
{
  "code": 409,
  "msg": "game account already bound to another store",
  "data": null
}

// Insufficient permissions
{
  "code": 403,
  "msg": "super admin role required",
  "data": null
}
```

---

### 12. List Game Accounts (Super Admin)

**Endpoint:** `GET /api/game/admin/ctrl-accounts`

**Description:** List all game accounts bound by super admin

**Authentication:** Required (Super Admin)

**Query Parameters:**
- `page` (optional): Page number, default 1
- `page_size` (optional): Page size, default 20

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "identifier": "gameuser123",
        "game_user_id": "12345",
        "house_gid": 123,
        "house_name": "Store A",
        "session_state": "online",
        "last_sync_at": "2025-11-11T10:00:00Z",
        "created_at": "2025-11-01T10:00:00Z"
      }
    ],
    "total": 1,
    "page_no": 1,
    "page_size": 20
  }
}
```

---

### 13. Unbind Game Account from Store (Super Admin)

**Endpoint:** `DELETE /api/game/admin/ctrl-accounts/:id`

**Description:** Unbind a game account from store

**Authentication:** Required (Super Admin)

**Permissions:** `game:ctrl:delete`

**Path Parameters:**
- `id`: Game ctrl account ID

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

---

### 14. Update Store Binding (Super Admin)

**Endpoint:** `PUT /api/game/admin/ctrl-accounts/:id/house`

**Description:** Update store binding for a game account

**Authentication:** Required (Super Admin)

**Permissions:** `game:ctrl:update`

**Path Parameters:**
- `id`: Game ctrl account ID

**Request Body:**
```json
{
  "house_gid": 456  // New store ID
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "house_gid": 456,
    "session_state": "online"
  }
}
```

---

## Session Management APIs

### 15. Start Session

**Endpoint:** `POST /api/game/accounts/sessionStart`

**Description:** Manually start a game session

**Authentication:** Required

**Permissions:** `game:ctrl:create`

**Request Body:**
```json
{
  "id": 1,          // Game ctrl account ID
  "house_gid": 123  // Store ID
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "state": "online"
  }
}
```

---

### 16. Stop Session

**Endpoint:** `POST /api/game/accounts/sessionStop`

**Description:** Stop a game session

**Authentication:** Required

**Permissions:** `game:ctrl:update`

**Request Body:**
```json
{
  "id": 1,          // Game ctrl account ID
  "house_gid": 123  // Store ID
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "state": "offline"
  }
}
```

---

## Game Record APIs

### 17. Get Game Records

**Endpoint:** `GET /api/game/records`

**Description:** Get game records for a store

**Authentication:** Required

**Query Parameters:**
- `house_gid` (required): Store ID
- `start_date` (optional): Start date (YYYY-MM-DD)
- `end_date` (optional): End date (YYYY-MM-DD)
- `page` (optional): Page number, default 1
- `page_size` (optional): Page size, default 20

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "battle_record_id": 100,
        "player_gid": 12345,
        "game_account_id": 1,
        "nickname": "Player Name",
        "score_delta": 100,
        "is_winner": true,
        "battle_at": "2025-11-11T10:00:00Z",
        "kind_id": 1,
        "room_uid": 1001
      }
    ],
    "total": 100,
    "page_no": 1,
    "page_size": 20
  }
}
```

---

### 18. Get Player Statistics

**Endpoint:** `GET /api/game/players/:player_gid/stats`

**Description:** Get statistics for a specific player

**Authentication:** Required

**Path Parameters:**
- `player_gid`: Player game ID

**Query Parameters:**
- `house_gid` (required): Store ID
- `start_date` (optional): Start date
- `end_date` (optional): End date

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "player_gid": 12345,
    "nickname": "Player Name",
    "total_games": 100,
    "total_wins": 60,
    "win_rate": 0.6,
    "total_score": 5000,
    "avg_score": 50,
    "best_score": 500,
    "worst_score": -200
  }
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request / Validation Failed |
| 401 | Unauthorized / Token Invalid |
| 403 | Forbidden / Insufficient Permissions |
| 404 | Not Found |
| 409 | Conflict / Constraint Violation |
| 500 | Internal Server Error |

---

## Rate Limiting

- Authentication endpoints: 10 requests per minute per IP
- Other endpoints: 100 requests per minute per user
- Bulk operations: 10 requests per minute per user

**Rate Limit Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1699999999
```

---

## Webhooks (Future)

### Game Record Sync Complete
```json
{
  "event": "game_record_sync_complete",
  "timestamp": "2025-11-11T10:00:00Z",
  "data": {
    "house_gid": 123,
    "records_synced": 1000,
    "sync_duration_ms": 5000
  }
}
```

---

## SDK Examples

### JavaScript/TypeScript
```typescript
import axios from 'axios';
import md5 from 'md5';

const API_BASE = 'http://localhost:8080/api';

// Register user
async function register(username: string, password: string, gameAccount: string, gamePassword: string) {
  const response = await axios.post(`${API_BASE}/login/register`, {
    username,
    password,
    game_login_mode: 'account',
    game_account: gameAccount,
    game_password_md5: md5(gamePassword)
  });
  
  return response.data;
}

// Login
async function login(username: string, password: string) {
  const response = await axios.post(`${API_BASE}/login/username`, {
    username,
    password
  });
  
  // Save token
  localStorage.setItem('token', response.data.data.token);
  
  return response.data;
}

// Get user info
async function getUserInfo() {
  const token = localStorage.getItem('token');
  const response = await axios.get(`${API_BASE}/basic/user/me`, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  });
  
  return response.data;
}
```

### Go
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type RegisterRequest struct {
    Username        string `json:"username"`
    Password        string `json:"password"`
    GameLoginMode   string `json:"game_login_mode"`
    GameAccount     string `json:"game_account"`
    GamePasswordMD5 string `json:"game_password_md5"`
}

func Register(req RegisterRequest) error {
    body, _ := json.Marshal(req)
    resp, err := http.Post(
        "http://localhost:8080/api/login/register",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Handle response...
    return nil
}
```

