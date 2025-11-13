# API Services Documentation

This directory contains all frontend API service functions that communicate with the backend.

## Directory Structure

```
services/
├── index.ts                    # Centralized exports
├── README.md                   # This file
├── basic/                      # Basic system services
│   ├── user/                   # User management
│   ├── role/                   # Role management
│   └── menu/                   # Menu management
├── game/                       # Game-related services
│   ├── account/                # Game account verification & binding
│   ├── ctrl-account/           # Control account management (Super Admin)
│   ├── session/                # Session management
│   └── funds/                  # Funds management
├── shops/                      # Shop management services
│   ├── houses/                 # House/Store options
│   ├── ctrlAccounts/           # Control accounts (legacy)
│   ├── admins/                 # Shop admins
│   ├── applications/           # Shop applications
│   ├── fees/                   # Fee management
│   ├── groups/                 # Group management
│   ├── members/                # Shop members
│   └── tables/                 # Table management
├── members/                    # Member-related services
│   ├── battle/                 # Battle records
│   ├── ledger/                 # Ledger management
│   └── wallet/                 # Wallet management
├── platforms/                  # Platform services
├── applications/               # Application services
├── stats/                      # Statistics services
└── login/                      # Login services
```

## Core Services

### 1. Authentication & User Management

#### Login (`services/login`)
- User authentication
- Token management
- Session handling

#### User Management (`services/basic/user`)
- User CRUD operations
- Profile management
- User role assignment

### 2. Game Account Management

#### Game Account (`services/game/account`)
**Purpose**: Regular user game account verification and binding

**Key Functions**:
- `gameAccountVerify(data)` - Verify game account credentials
- `gameAccountBind(data)` - Bind game account to user
- `gameAccountMe()` - Get current user's game account
- `gameAccountDelete()` - Delete game account binding

**Required During**:
- User registration (must bind game account)
- Profile management

#### Control Account (`services/game/ctrl-account`)
**Purpose**: Super admin control account management for store operations

**Key Functions**:
- `createCtrlAccount(data)` - Create/update control account
- `bindCtrlAccount(data)` - Bind control account to store
- `unbindCtrlAccount(data)` - Unbind control account from store
- `listCtrlAccountsByHouse(data)` - List accounts by store
- `listAllCtrlAccounts(data)` - List all control accounts with filters
- `getHouseOptions()` - Get available store options
- `startSession(data)` - Start session for control account
- `stopSession(data)` - Stop session for control account

**Permissions Required**:
- `game:ctrl:create` - Create accounts and start sessions
- `game:ctrl:update` - Update accounts, bind/unbind, stop sessions
- `game:ctrl:view` - View account lists

**Use Cases**:
- Super admin managing multiple game accounts
- Binding game accounts to stores
- Starting/stopping store sessions for auto-sync

#### Session Management (`services/game/session`)
**Purpose**: Manage game sessions for control accounts

**Key Functions**:
- `gameSessionStart(data)` - Start a session
- `gameSessionStop(data)` - Stop a session
- `gameSessionQuery(params)` - Query session details (TODO: backend)
- `gameSessionDetail(sessionId)` - Get session detail (TODO: backend)
- `gameSessionSyncLogs(params)` - Get sync logs (TODO: backend)

**Session States**:
- `active` - Session is running
- `inactive` - Session is stopped
- `error` - Session encountered an error

**Sync Status**:
- `idle` - Not currently syncing
- `syncing` - Sync in progress
- `error` - Sync error occurred

**Sync Types**:
- `battle_record` - Battle records (syncs every 5 seconds)
- `member_list` - Member list (syncs every 30 seconds)
- `wallet_update` - Wallet transactions (syncs every 10 seconds)
- `room_list` - Room list
- `group_member` - Group members

### 3. Shop Management

#### Houses (`services/shops/houses`)
- `shopsHousesOptions()` - Get list of available store IDs

#### Control Accounts (Legacy) (`services/shops/ctrlAccounts`)
**Note**: This is the older version. Use `services/game/ctrl-account` for new code.

- `shopsCtrlAccountsCreate(data)` - Create control account
- `shopsCtrlAccountsBind(data)` - Bind to store
- `shopsCtrlAccountsUnbind(data)` - Unbind from store
- `shopsCtrlAccountsList(data)` - List by store
- `shopsCtrlAccountsListAll(data)` - List all

## Usage Examples

### Example 1: User Registration with Game Account

```typescript
import { gameAccountVerify, gameAccountBind } from '@/services/game/account';
import { md5Upper } from '@/utils/crypto';

// Step 1: Verify game account
const verifyResult = await gameAccountVerify({
  mode: 'account',
  account: 'player123',
  pwd_md5: md5Upper('password123'),
});

if (verifyResult.ok) {
  // Step 2: Bind to user during registration
  await gameAccountBind({
    mode: 'account',
    account: 'player123',
    pwd_md5: md5Upper('password123'),
    nickname: 'Player Name', // Retrieved from game server
  });
}
```

### Example 2: Super Admin Managing Control Accounts

```typescript
import {
  createCtrlAccount,
  bindCtrlAccount,
  listAllCtrlAccounts,
  startSession,
} from '@/services/game/ctrl-account';
import { md5Upper } from '@/utils/crypto';

// Step 1: Create control account
const account = await createCtrlAccount({
  login_mode: 'account',
  identifier: 'ctrl_account_1',
  pwd_md5: md5Upper('password'),
  status: 1,
});

// Step 2: Bind to store
await bindCtrlAccount({
  ctrl_id: account.id,
  house_gid: 12345,
  status: 1,
});

// Step 3: Start session (auto-sync begins)
await startSession({
  id: account.id,
  house_gid: 12345,
});

// Step 4: List all accounts
const accounts = await listAllCtrlAccounts({
  status: 1,
  page: 1,
  size: 20,
});
```

### Example 3: Store Admin Viewing Sessions

```typescript
import { listCtrlAccountsByHouse } from '@/services/game/ctrl-account';
import { gameSessionStart, gameSessionStop } from '@/services/game/session';

// Get control accounts for store
const accounts = await listCtrlAccountsByHouse({
  house_gid: 12345,
});

// Start session for an account
await gameSessionStart({
  id: accounts[0].id,
  house_gid: 12345,
});

// Stop session
await gameSessionStop({
  id: accounts[0].id,
  house_gid: 12345,
});
```

## API Response Format

All API responses follow this format:

```typescript
{
  code: number;      // 0 = success, non-zero = error
  message: string;   // Error message or success message
  data: T;          // Response data (type varies by endpoint)
}
```

## Error Handling

Use the `useRequest` hook for automatic error handling:

```typescript
import { useRequest } from '@/hooks/use-request';
import { listAllCtrlAccounts } from '@/services/game/ctrl-account';

const { data, loading, error, run } = useRequest(listAllCtrlAccounts, {
  manual: true,
});

// Call the API
await run({ status: 1 });
```

## Future Enhancements

### Planned Backend Endpoints (Not Yet Implemented)

1. **Session Detail Query**
   - `GET /game/sessions/:id` - Get session details
   - `POST /game/sessions/query` - Query sessions with filters

2. **Sync Log Query**
   - `GET /game/sessions/:id/logs` - Get sync logs for a session
   - Support filtering by sync_type and status

3. **Real-time Updates**
   - WebSocket support for live session status
   - Push notifications for sync errors

## Notes

- All endpoints require JWT authentication (Bearer token)
- Permission checks are enforced on the backend
- Use `md5Upper()` utility for password hashing (32-char uppercase MD5)
- Session auto-sync starts immediately after session start
- Control accounts can be bound to multiple stores (one at a time per account)

