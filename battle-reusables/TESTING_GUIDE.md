# Testing and Integration Guide

This document provides comprehensive testing procedures for the game account integration features.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Test Environment Setup](#test-environment-setup)
3. [Feature Testing Checklist](#feature-testing-checklist)
4. [Integration Testing](#integration-testing)
5. [Known Issues and Troubleshooting](#known-issues-and-troubleshooting)

## Prerequisites

### Backend Requirements
- ✅ Backend server running at configured API endpoint
- ✅ Database with proper schema (see `battle-tiles/doc/ddl.sql`)
- ✅ Game server accessible at `androidsc.foxuc.com:8200`
- ✅ Test game accounts available for validation

### Frontend Requirements
- ✅ Node.js and npm installed
- ✅ Expo CLI installed (`npm install -g expo-cli`)
- ✅ React Native development environment set up
- ✅ Mobile device or emulator for testing

### Test Accounts
Prepare the following test accounts:

1. **Super Admin Account**
   - Username: `super_admin_test`
   - Role: `super_admin`
   - Purpose: Test control account management

2. **Store Admin Account**
   - Username: `store_admin_test`
   - Role: `store_admin`
   - Purpose: Test store dashboard

3. **Regular User Account**
   - Username: `user_test`
   - Role: `user`
   - Purpose: Test basic registration and account binding

4. **Game Accounts** (for testing)
   - Valid game account with correct password
   - Valid game account with incorrect password
   - Invalid game account

## Test Environment Setup

### 1. Start Backend Server

```bash
cd battle-tiles
go run cmd/main.go
```

Verify backend is running:
- Check logs for successful startup
- Verify database connection
- Confirm API endpoints are accessible

### 2. Start Frontend Development Server

```bash
cd battle-reusables
npm install
npm run dev
```

Choose platform:
- Press `w` for web
- Press `a` for Android
- Press `i` for iOS

### 3. Configure API Endpoint

Verify `utils/request.ts` has correct backend URL:
```typescript
const BASE_URL = 'http://localhost:8080/api'; // Adjust as needed
```

## Feature Testing Checklist

### Phase 1: User Registration with Game Account Validation ✅

#### Test Case 1.1: Valid Game Account Registration
**Steps:**
1. Navigate to registration page (`/auth`)
2. Select "游戏账号" (Game Account) as account type
3. Enter valid game account credentials:
   - Account: `[valid_game_account]`
   - Password: `[correct_password]`
4. Wait for validation (should show green checkmark)
5. Fill in other required fields (username, password)
6. Click "注册" (Register)

**Expected Results:**
- ✅ Validation shows loading spinner during check
- ✅ Green checkmark appears when validation succeeds
- ✅ Game nickname is automatically fetched
- ✅ Register button becomes enabled
- ✅ Registration completes successfully
- ✅ User is redirected to main app
- ✅ User profile shows game nickname

**Verification:**
```sql
-- Check database
SELECT * FROM basic_users WHERE username = 'test_user';
SELECT * FROM game_accounts WHERE user_id = [user_id];
```

#### Test Case 1.2: Invalid Game Account Registration
**Steps:**
1. Navigate to registration page
2. Select "游戏账号" (Game Account)
3. Enter invalid credentials:
   - Account: `invalid_account_12345`
   - Password: `wrong_password`
4. Wait for validation

**Expected Results:**
- ✅ Validation shows loading spinner
- ✅ Red X appears when validation fails
- ✅ Error message displayed: "游戏账号验证失败"
- ✅ Register button remains disabled
- ✅ User cannot proceed with registration

#### Test Case 1.3: Network Error Handling
**Steps:**
1. Disconnect network or stop backend
2. Navigate to registration page
3. Enter game account credentials
4. Wait for validation

**Expected Results:**
- ✅ Loading spinner appears
- ✅ Error message displayed after timeout
- ✅ User can retry validation
- ✅ Register button remains disabled

### Phase 2: Super Admin Control Account Management ✅

#### Test Case 2.1: Create Control Account
**Steps:**
1. Login as super admin
2. Navigate to Profile tab
3. Scroll to "中控账号管理" section
4. Click "添加中控账号" button
5. Fill in form:
   - Login Mode: "游戏账号"
   - Account: `[test_ctrl_account]`
   - Password: `[password]`
6. Click "验证并添加"

**Expected Results:**
- ✅ Validation occurs before creation
- ✅ Success toast appears
- ✅ New account appears in list
- ✅ Account shows status badge (green = enabled)

**Verification:**
```sql
SELECT * FROM game_ctrl_accounts WHERE identifier = 'test_ctrl_account';
```

#### Test Case 2.2: Bind Control Account to Store
**Steps:**
1. In control account list, find the account
2. Click "绑定店铺" button
3. Enter store ID: `12345`
4. Click "绑定"

**Expected Results:**
- ✅ Success toast appears
- ✅ Store ID appears in account's store list
- ✅ Store badge shows "启用" status

**Verification:**
```sql
SELECT * FROM game_account_houses 
WHERE game_ctrl_account_id = [ctrl_id] AND house_gid = 12345;
```

#### Test Case 2.3: Start Session
**Steps:**
1. Find account with bound store
2. Click "启动会话" button for that store

**Expected Results:**
- ✅ Loading indicator appears
- ✅ Success toast: "会话已启动"
- ✅ Session status updates to "活跃"
- ✅ Backend creates game_session record
- ✅ Auto-sync worker starts

**Verification:**
```sql
SELECT * FROM game_sessions 
WHERE game_ctrl_account_id = [ctrl_id] AND house_gid = 12345;

-- Check sync logs after a few seconds
SELECT * FROM game_sync_logs WHERE session_id = [session_id] ORDER BY created_at DESC LIMIT 10;
```

#### Test Case 2.4: Stop Session
**Steps:**
1. Find active session
2. Click "停止会话" button

**Expected Results:**
- ✅ Success toast: "会话已停止"
- ✅ Session status updates to "非活跃"
- ✅ Auto-sync worker stops
- ✅ No new sync logs created

#### Test Case 2.5: Unbind Control Account
**Steps:**
1. Find account with bound store
2. Click "解绑" button
3. Confirm action

**Expected Results:**
- ✅ Confirmation dialog appears
- ✅ Success toast after confirmation
- ✅ Store removed from account's store list
- ✅ Session automatically stopped if active

**Verification:**
```sql
-- Should return no rows
SELECT * FROM game_account_houses 
WHERE game_ctrl_account_id = [ctrl_id] AND house_gid = 12345;
```

### Phase 3: Store Admin Dashboard ✅

#### Test Case 3.1: View Store Sessions
**Steps:**
1. Login as store admin
2. Navigate to "会话管理" (Session Management)
3. Page should auto-load with store admin's store

**Expected Results:**
- ✅ Store ID auto-populated
- ✅ Store info card displays:
  - Store ID
  - Number of control accounts
  - Number of active sessions
- ✅ Control account cards show:
  - Account ID and identifier
  - Login mode
  - Status badge (active/inactive)
  - Session status
  - Sync information

#### Test Case 3.2: Pull to Refresh
**Steps:**
1. On session management page
2. Pull down to refresh (mobile) or click refresh button

**Expected Results:**
- ✅ Loading indicator appears
- ✅ Data refreshes
- ✅ Updated session statuses displayed
- ✅ Sync times updated

#### Test Case 3.3: View Sync Information
**Steps:**
1. Find an active session
2. Expand sync information section

**Expected Results:**
- ✅ Shows sync types:
  - Battle records (5s interval)
  - Member list (30s interval)
  - Wallet transactions (10s interval)
- ✅ Shows last sync time
- ✅ Shows sync status (idle/syncing/error)

### Phase 4: Auto-Sync Verification ✅

#### Test Case 4.1: Battle Record Sync
**Steps:**
1. Start a session for a store
2. Play a game in that store (or simulate)
3. Wait 5-10 seconds
4. Check database

**Expected Results:**
- ✅ New battle records appear in database
- ✅ Sync logs show successful syncs
- ✅ Records synced every ~5 seconds

**Verification:**
```sql
-- Check battle records
SELECT COUNT(*) FROM game_battle_records WHERE house_gid = 12345;

-- Check sync logs
SELECT * FROM game_sync_logs 
WHERE session_id = [session_id] 
AND sync_type = 'battle_record' 
ORDER BY created_at DESC LIMIT 5;
```

#### Test Case 4.2: Member List Sync
**Steps:**
1. Ensure session is active
2. Wait 30-60 seconds
3. Check database

**Expected Results:**
- ✅ Member list updated in database
- ✅ Sync logs show successful syncs every ~30 seconds

**Verification:**
```sql
SELECT * FROM game_sync_logs 
WHERE session_id = [session_id] 
AND sync_type = 'member_list' 
ORDER BY created_at DESC LIMIT 3;
```

#### Test Case 4.3: Wallet Transaction Sync
**Steps:**
1. Ensure session is active
2. Perform wallet transaction in game
3. Wait 10-20 seconds
4. Check database

**Expected Results:**
- ✅ Wallet transactions appear in database
- ✅ Sync logs show successful syncs every ~10 seconds

## Integration Testing

### End-to-End Flow 1: New User Registration to Game Play

1. **Register New User**
   - Complete registration with game account
   - Verify game nickname appears in profile

2. **Login**
   - Login with new credentials
   - Verify redirect to main app

3. **View Profile**
   - Check game account is bound
   - Verify nickname is correct

### End-to-End Flow 2: Super Admin Complete Workflow

1. **Login as Super Admin**
   - Verify super admin permissions

2. **Create Control Account**
   - Add new control account
   - Verify validation works

3. **Bind to Store**
   - Bind account to test store
   - Verify binding successful

4. **Start Session**
   - Start session for store
   - Verify session becomes active

5. **Monitor Sync**
   - Wait 1-2 minutes
   - Check sync logs in database
   - Verify all sync types working

6. **Stop Session**
   - Stop the session
   - Verify sync stops

7. **Unbind Account**
   - Unbind from store
   - Verify cleanup

### End-to-End Flow 3: Store Admin Monitoring

1. **Login as Store Admin**
   - Verify store admin permissions

2. **View Dashboard**
   - Navigate to session management
   - Verify auto-load works

3. **Monitor Sessions**
   - Check session statuses
   - Verify sync information displays

4. **Refresh Data**
   - Pull to refresh
   - Verify updates

## Performance Testing

### Load Testing Scenarios

1. **Multiple Sessions**
   - Start 5-10 sessions simultaneously
   - Monitor CPU and memory usage
   - Verify all sessions sync correctly

2. **High Frequency Sync**
   - Monitor sync performance over 1 hour
   - Check for memory leaks
   - Verify sync intervals remain consistent

3. **Large Data Sets**
   - Test with stores having 100+ members
   - Test with 1000+ battle records
   - Verify sync performance

## Known Issues and Troubleshooting

### Issue 1: Game Account Validation Timeout
**Symptom:** Validation takes too long or times out
**Solution:**
- Check game server connectivity
- Verify firewall rules
- Increase timeout in request configuration

### Issue 2: Session Won't Start
**Symptom:** Session start fails with error
**Possible Causes:**
- Game account credentials invalid
- Store ID doesn't exist
- Network connectivity issues
**Solution:**
- Verify credentials
- Check backend logs
- Test game server connection

### Issue 3: Sync Not Working
**Symptom:** No sync logs appearing
**Possible Causes:**
- Session not actually started
- Worker not running
- Database connection issues
**Solution:**
- Check session state in database
- Verify worker is running (check backend logs)
- Test database connectivity

### Issue 4: UI Not Updating
**Symptom:** Status doesn't update after actions
**Solution:**
- Pull to refresh
- Check network requests in browser dev tools
- Verify API responses

## Test Results Template

Use this template to document test results:

```markdown
## Test Run: [Date]
**Tester:** [Name]
**Environment:** [Dev/Staging/Production]
**Platform:** [Web/iOS/Android]

### Phase 1: Registration
- [ ] Test Case 1.1: Valid Account - PASS/FAIL
- [ ] Test Case 1.2: Invalid Account - PASS/FAIL
- [ ] Test Case 1.3: Network Error - PASS/FAIL

### Phase 2: Control Account Management
- [ ] Test Case 2.1: Create Account - PASS/FAIL
- [ ] Test Case 2.2: Bind to Store - PASS/FAIL
- [ ] Test Case 2.3: Start Session - PASS/FAIL
- [ ] Test Case 2.4: Stop Session - PASS/FAIL
- [ ] Test Case 2.5: Unbind Account - PASS/FAIL

### Phase 3: Store Admin Dashboard
- [ ] Test Case 3.1: View Sessions - PASS/FAIL
- [ ] Test Case 3.2: Pull to Refresh - PASS/FAIL
- [ ] Test Case 3.3: View Sync Info - PASS/FAIL

### Phase 4: Auto-Sync
- [ ] Test Case 4.1: Battle Records - PASS/FAIL
- [ ] Test Case 4.2: Member List - PASS/FAIL
- [ ] Test Case 4.3: Wallet Transactions - PASS/FAIL

### Issues Found:
1. [Description]
2. [Description]

### Notes:
[Additional observations]
```

## Automated Testing (Future Enhancement)

For future implementation, consider adding:

1. **Unit Tests** (Jest + React Native Testing Library)
   - Component rendering tests
   - Hook behavior tests
   - Utility function tests

2. **Integration Tests** (Detox)
   - E2E user flows
   - Navigation tests
   - API integration tests

3. **API Tests** (Postman/Newman)
   - Endpoint validation
   - Response format verification
   - Error handling tests

4. **Performance Tests** (Lighthouse/React Native Performance)
   - Load time measurements
   - Memory usage monitoring
   - Render performance

## Conclusion

This testing guide covers all major features implemented in the game account integration project. Follow the test cases systematically to ensure all functionality works as expected. Document any issues found and report them to the development team.

