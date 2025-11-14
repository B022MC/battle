# Migration Guide - User Management System Enhancement

## Overview

This guide provides step-by-step instructions for migrating from the WeChat-dependent user management system to the new game-account-based system.

---

## Pre-Migration Checklist

### 1. Environment Preparation
- [ ] Backup production database
- [ ] Verify database connection credentials
- [ ] Ensure `pg_dump` is installed for schema export
- [ ] Set up staging environment for testing
- [ ] Review current user data statistics

### 2. Data Analysis
```sql
-- Count total users
SELECT COUNT(*) FROM basic_user WHERE deleted_at IS NULL;

-- Count users with WeChat ID only
SELECT COUNT(*) FROM basic_user 
WHERE wechat_id IS NOT NULL 
  AND wechat_id != '' 
  AND deleted_at IS NULL;

-- Count users with game accounts
SELECT COUNT(DISTINCT user_id) FROM game_account WHERE is_del = 0;

-- Count game accounts with multiple store bindings
SELECT game_account_id, COUNT(*) as binding_count
FROM game_account_house
GROUP BY game_account_id
HAVING COUNT(*) > 1;

-- Count users managing multiple stores
SELECT user_id, COUNT(*) as store_count
FROM game_shop_admin
WHERE deleted_at IS NULL
GROUP BY user_id
HAVING COUNT(*) > 1;
```

### 3. Communication Plan
- [ ] Notify all users about upcoming changes
- [ ] Prepare user documentation for new registration flow
- [ ] Set up support channels for migration assistance
- [ ] Create FAQ document

---

## Migration Steps

### Phase 1: Database Schema Migration

#### Step 1.1: Backup Database
```bash
# Full database backup
pg_dump -U your_db_user -h localhost -p 5432 -d battle_tiles_db \
  -F c -f backup_$(date +%Y%m%d_%H%M%S).dump

# Verify backup
pg_restore --list backup_*.dump | head -20
```

#### Step 1.2: Export Current Schema
```bash
# Export current schema for reference
bash doc/generate_full_schema.sh > doc/schema_before_migration.sql
```

#### Step 1.3: Handle Constraint Violations

**A. Archive Duplicate Game Account House Bindings**
```sql
-- Identify duplicates
SELECT game_account_id, COUNT(*) as count, 
       array_agg(id ORDER BY created_at DESC) as ids
FROM game_account_house
GROUP BY game_account_id
HAVING COUNT(*) > 1;

-- Archive older bindings (keep most recent)
BEGIN;

-- Create archive table if needed
CREATE TABLE IF NOT EXISTS game_account_house_archive (
  LIKE game_account_house INCLUDING ALL
);

-- Move duplicates to archive
WITH ranked_bindings AS (
  SELECT id, game_account_id,
         ROW_NUMBER() OVER (PARTITION BY game_account_id ORDER BY created_at DESC) as rn
  FROM game_account_house
)
INSERT INTO game_account_house_archive
SELECT gah.*
FROM game_account_house gah
INNER JOIN ranked_bindings rb ON gah.id = rb.id
WHERE rb.rn > 1;

-- Delete duplicates from main table
WITH ranked_bindings AS (
  SELECT id, game_account_id,
         ROW_NUMBER() OVER (PARTITION BY game_account_id ORDER BY created_at DESC) as rn
  FROM game_account_house
)
DELETE FROM game_account_house
WHERE id IN (
  SELECT id FROM ranked_bindings WHERE rn > 1
);

COMMIT;
```

**B. Archive Duplicate Store Admin Assignments**
```sql
BEGIN;

-- Create archive table if needed
CREATE TABLE IF NOT EXISTS game_shop_admin_archive (
  LIKE game_shop_admin INCLUDING ALL
);

-- Identify and archive duplicates
WITH ranked_admins AS (
  SELECT id, user_id,
         ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
  FROM game_shop_admin
  WHERE deleted_at IS NULL
)
INSERT INTO game_shop_admin_archive
SELECT gsa.*
FROM game_shop_admin gsa
INNER JOIN ranked_admins ra ON gsa.id = ra.id
WHERE ra.rn > 1;

-- Soft delete duplicates
WITH ranked_admins AS (
  SELECT id, user_id,
         ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
  FROM game_shop_admin
  WHERE deleted_at IS NULL
)
UPDATE game_shop_admin
SET deleted_at = NOW()
WHERE id IN (
  SELECT id FROM ranked_admins WHERE rn > 1
);

COMMIT;
```

#### Step 1.4: Apply DDL Changes
```bash
# Apply user management DDL
psql -U your_db_user -h localhost -p 5432 -d battle_tiles_db \
  -f doc/user_management.ddl

# Verify constraints were added
psql -U your_db_user -h localhost -p 5432 -d battle_tiles_db -c "
  SELECT conname, contype, conrelid::regclass
  FROM pg_constraint
  WHERE conname LIKE '%game_account%' OR conname LIKE '%shop_admin%'
  ORDER BY conrelid::regclass, conname;
"
```

#### Step 1.5: Verify Schema Changes
```sql
-- Verify game_player_record table exists
SELECT table_name, column_name, data_type
FROM information_schema.columns
WHERE table_name = 'game_player_record'
ORDER BY ordinal_position;

-- Verify constraints
SELECT constraint_name, constraint_type
FROM information_schema.table_constraints
WHERE table_name IN ('game_account', 'game_account_house', 'game_shop_admin')
ORDER BY table_name, constraint_name;

-- Verify indexes
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'game_player_record'
ORDER BY indexname;
```

#### Step 1.6: Export Updated Schema
```bash
# Export new schema for documentation
bash doc/generate_full_schema.sh > doc/schema_after_migration.sql
```

---

### Phase 2: Backend Code Deployment

#### Step 2.1: Update Dependencies
```bash
cd battle-tiles

# Ensure all Go dependencies are up to date
go mod tidy
go mod verify

# Run tests
go test ./...
```

#### Step 2.2: Create New Model Files

**A. Create GamePlayerRecord Model**
```bash
# File: internal/dal/model/game/game_player_record.go
cat > internal/dal/model/game/game_player_record.go << 'EOF'
package game

import "time"

const TableNameGamePlayerRecord = "game_player_record"

type GamePlayerRecord struct {
	Id             int32     `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	BattleRecordID int32     `gorm:"column:battle_record_id;not null" json:"battle_record_id"`
	HouseGID       int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	PlayerGID      int64     `gorm:"column:player_gid;not null" json:"player_gid"`
	GameAccountID  *int32    `gorm:"column:game_account_id" json:"game_account_id"`
	GroupID        int32     `gorm:"column:group_id;not null" json:"group_id"`
	RoomUID        int32     `gorm:"column:room_uid;not null" json:"room_uid"`
	KindID         int32     `gorm:"column:kind_id;not null" json:"kind_id"`
	BaseScore      int32     `gorm:"column:base_score;not null" json:"base_score"`
	ScoreDelta     int32     `gorm:"column:score_delta;not null;default:0" json:"score_delta"`
	IsWinner       bool      `gorm:"column:is_winner;not null;default:false" json:"is_winner"`
	BattleAt       time.Time `gorm:"column:battle_at;type:timestamp with time zone;not null" json:"battle_at"`
	MetaJSON       string    `gorm:"column:meta_json;type:jsonb;not null;default:'{}'::jsonb" json:"meta_json"`
}

func (GamePlayerRecord) TableName() string { return TableNameGamePlayerRecord }
EOF
```

#### Step 2.3: Update Registration Request Model
```bash
# Backup original file
cp internal/dal/req/basic_user.go internal/dal/req/basic_user.go.bak

# Update RegisterRequest struct
# (Manual edit required - see technical_specifications.md section 2.1.1)
```

#### Step 2.4: Update Registration Business Logic
```bash
# Backup original file
cp internal/biz/basic/basic_login.go internal/biz/basic/basic_login.go.bak

# Update Register method
# (Manual edit required - see technical_specifications.md section 2.1.2)
```

#### Step 2.5: Build and Test
```bash
# Build the application
go build -o battle-tiles ./cmd/server

# Run unit tests
go test ./internal/biz/basic/... -v
go test ./internal/service/basic/... -v

# Run integration tests
go test ./internal/... -integration -v
```

#### Step 2.6: Deploy Backend
```bash
# Stop current service
systemctl stop battle-tiles

# Backup current binary
cp /usr/local/bin/battle-tiles /usr/local/bin/battle-tiles.bak

# Deploy new binary
cp battle-tiles /usr/local/bin/battle-tiles

# Start service
systemctl start battle-tiles

# Check logs
journalctl -u battle-tiles -f
```

---

### Phase 3: Frontend Deployment

#### Step 3.1: Update Frontend Code
```bash
cd battle-reusables

# Install dependencies
npm install md5

# Update registration component
# (See technical_specifications.md section 3.1)
```

#### Step 3.2: Build Frontend
```bash
# Build for production
npm run build

# Test build locally
npm run start
```

#### Step 3.3: Deploy Frontend
```bash
# Deploy to hosting (example for Vercel)
vercel --prod

# Or deploy to custom server
rsync -avz --delete .next/ user@server:/var/www/battle-reusables/
ssh user@server 'pm2 restart battle-reusables'
```

---

### Phase 4: Data Migration for Existing Users

#### Step 4.1: Identify Users Needing Migration
```sql
-- Users with WeChat ID but no game account
CREATE TEMP TABLE users_to_migrate AS
SELECT u.id, u.username, u.nick_name, u.wechat_id
FROM basic_user u
LEFT JOIN game_account ga ON ga.user_id = u.id
WHERE u.wechat_id IS NOT NULL 
  AND u.wechat_id != ''
  AND ga.id IS NULL
  AND u.deleted_at IS NULL;

-- Export list for notification
\copy users_to_migrate TO 'users_to_migrate.csv' CSV HEADER;
```

#### Step 4.2: Notify Users
```bash
# Send email notifications to users
# Include instructions for binding game account
# Provide deadline for migration
```

#### Step 4.3: Create Migration Endpoint (Temporary)
```go
// internal/service/basic/basic_user.go
// Add temporary endpoint for existing users to bind game account

func (s *BasicUserService) BindGameAccount(c *gin.Context) {
    var in struct {
        GameLoginMode   string `json:"game_login_mode" binding:"required"`
        GameAccount     string `json:"game_account" binding:"required"`
        GamePasswordMD5 string `json:"game_password_md5" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&in); err != nil {
        response.Fail(c, ecode.ParamsFailed, err)
        return
    }
    
    claims, err := utils.GetClaims(c)
    if err != nil {
        response.Fail(c, ecode.TokenValidateFailed, err)
        return
    }
    
    // Verify game account and create binding
    // (Implementation similar to registration flow)
}
```

#### Step 4.4: Monitor Migration Progress
```sql
-- Check migration progress
SELECT 
  COUNT(*) FILTER (WHERE ga.id IS NOT NULL) as migrated,
  COUNT(*) FILTER (WHERE ga.id IS NULL) as pending,
  COUNT(*) as total
FROM users_to_migrate utm
LEFT JOIN game_account ga ON ga.user_id = utm.id;
```

---

### Phase 5: Post-Migration Cleanup

#### Step 5.1: Verify Data Integrity
```sql
-- Verify all active users have game accounts
SELECT COUNT(*) FROM basic_user u
LEFT JOIN game_account ga ON ga.user_id = u.id
WHERE u.deleted_at IS NULL
  AND ga.id IS NULL;
-- Should return 0

-- Verify no constraint violations
SELECT game_account_id, COUNT(*) 
FROM game_account_house
GROUP BY game_account_id
HAVING COUNT(*) > 1;
-- Should return 0 rows

-- Verify store admin constraints
SELECT user_id, COUNT(*)
FROM game_shop_admin
WHERE deleted_at IS NULL
GROUP BY user_id
HAVING COUNT(*) > 1;
-- Should return 0 rows
```

#### Step 5.2: Performance Testing
```bash
# Run load tests on new registration endpoint
ab -n 1000 -c 10 -p register.json -T application/json \
  http://localhost:8080/api/login/register

# Monitor database performance
psql -U your_db_user -d battle_tiles_db -c "
  SELECT schemaname, tablename, 
         pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
  FROM pg_tables
  WHERE schemaname = 'public'
  ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
"
```

#### Step 5.3: Deprecate WeChat Fields (Future)
```sql
-- After all users migrated, mark wechat_id as deprecated
COMMENT ON COLUMN basic_user.wechat_id IS 'DEPRECATED - Will be removed in future version';

-- Schedule for removal in next major version
-- ALTER TABLE basic_user DROP COLUMN wechat_id;
```

#### Step 5.4: Update Documentation
- [ ] Update API documentation
- [ ] Update user guides
- [ ] Update developer documentation
- [ ] Archive old documentation

---

## Rollback Plan

### If Migration Fails

#### Step 1: Stop Services
```bash
systemctl stop battle-tiles
pm2 stop battle-reusables
```

#### Step 2: Restore Database
```bash
# Drop new constraints
psql -U your_db_user -d battle_tiles_db << EOF
BEGIN;
ALTER TABLE game_account_house DROP CONSTRAINT IF EXISTS uk_game_account_house_account;
DROP INDEX IF EXISTS uk_game_shop_admin_user_active;
DROP TABLE IF EXISTS game_player_record;
COMMIT;
EOF

# Or full database restore if needed
pg_restore -U your_db_user -d battle_tiles_db -c backup_*.dump
```

#### Step 3: Restore Code
```bash
# Backend
cp /usr/local/bin/battle-tiles.bak /usr/local/bin/battle-tiles

# Frontend
git checkout HEAD~1
npm run build
vercel --prod
```

#### Step 4: Restart Services
```bash
systemctl start battle-tiles
pm2 start battle-reusables
```

---

## Monitoring and Alerts

### Key Metrics to Monitor

1. **Registration Success Rate**
   - Target: >95% success rate
   - Alert if <90% for 5 minutes

2. **Game Account Verification Time**
   - Target: <2 seconds average
   - Alert if >5 seconds average

3. **Database Constraint Violations**
   - Target: 0 violations
   - Alert on any violation

4. **User Migration Progress**
   - Track daily migration count
   - Alert if stalled for 3 days

### Monitoring Queries

```sql
-- Registration attempts today
SELECT DATE(created_at), COUNT(*)
FROM basic_user
WHERE created_at >= CURRENT_DATE
GROUP BY DATE(created_at);

-- Game account binding success rate
SELECT 
  COUNT(*) FILTER (WHERE ga.id IS NOT NULL) * 100.0 / COUNT(*) as success_rate
FROM basic_user u
LEFT JOIN game_account ga ON ga.user_id = u.id
WHERE u.created_at >= CURRENT_DATE;
```

---

## Support and Troubleshooting

### Common Issues

**Issue 1: Game account verification fails**
- Check game API connectivity
- Verify credentials are correct
- Check rate limiting

**Issue 2: Constraint violation on game_account_house**
- User trying to bind account already bound to another store
- Solution: Unbind from previous store first

**Issue 3: Registration fails with "username already exists"**
- Username collision
- Solution: Choose different username

**Issue 4: Slow game record synchronization**
- Large data volume
- Solution: Implement pagination and batch processing

### Support Contacts
- Technical Lead: [contact info]
- Database Admin: [contact info]
- DevOps: [contact info]

---

## Success Criteria

- [ ] All database constraints applied successfully
- [ ] Zero constraint violations in production
- [ ] New registration flow working correctly
- [ ] >90% of existing users migrated within 30 days
- [ ] No critical bugs reported
- [ ] Performance metrics within acceptable range
- [ ] All documentation updated
- [ ] Team trained on new system

---

## Timeline

| Phase | Duration | Start Date | End Date |
|-------|----------|------------|----------|
| Pre-Migration Prep | 1 week | TBD | TBD |
| Database Migration | 1 day | TBD | TBD |
| Backend Deployment | 1 day | TBD | TBD |
| Frontend Deployment | 1 day | TBD | TBD |
| User Migration Period | 30 days | TBD | TBD |
| Post-Migration Cleanup | 1 week | TBD | TBD |

**Total Estimated Time:** 6 weeks

