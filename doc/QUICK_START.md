# Quick Start Guide - Developer Setup

## Overview

This guide helps developers quickly set up their environment and start implementing the user management system enhancements.

---

## Prerequisites

### Required Software
- **Go**: 1.21+ ([Download](https://golang.org/dl/))
- **Node.js**: 18+ ([Download](https://nodejs.org/))
- **PostgreSQL**: 14+ ([Download](https://www.postgresql.org/download/))
- **Redis**: 7+ ([Download](https://redis.io/download))
- **Git**: Latest version

### Recommended Tools
- **VS Code** with Go and TypeScript extensions
- **Postman** or **Insomnia** for API testing
- **pgAdmin** or **DBeaver** for database management
- **Docker** (optional, for containerized development)

---

## Environment Setup

### 1. Clone Repository

```bash
cd /Users/b022mc/project/battle
git checkout -b feature/user-management-enhancement
```

### 2. Backend Setup (battle-tiles)

```bash
cd battle-tiles

# Install dependencies
go mod download
go mod verify

# Copy environment config
cp configs/config.example.yaml configs/config.yaml

# Update database connection in config.yaml
# db:
#   host: localhost
#   port: 5432
#   user: your_user
#   password: your_password
#   database: battle_tiles_dev
```

### 3. Database Setup

```bash
# Create development database
createdb battle_tiles_dev

# Apply current schema (if exists)
psql -U your_user -d battle_tiles_dev -f doc/schema_before_migration.sql

# Or let GORM auto-migrate
# (Run the application once to auto-create tables)
```

### 4. Apply New Schema Changes

```bash
# Apply user management DDL
psql -U your_user -d battle_tiles_dev -f doc/user_management.ddl

# Verify changes
psql -U your_user -d battle_tiles_dev -c "\d game_player_record"
```

### 5. Frontend Setup (battle-reusables)

```bash
cd ../battle-reusables

# Install dependencies
npm install

# Install additional dependencies
npm install md5

# Copy environment config
cp .env.example .env.local

# Update API endpoint in .env.local
# NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### 6. Start Services

**Terminal 1 - Redis:**
```bash
redis-server
```

**Terminal 2 - Backend:**
```bash
cd battle-tiles
go run cmd/server/main.go
# Backend should start on http://localhost:8080
```

**Terminal 3 - Frontend:**
```bash
cd battle-reusables
npm run dev
# Frontend should start on http://localhost:3000
```

---

## Development Workflow

### Step 1: Create Feature Branch

```bash
git checkout -b feature/registration-enhancement
```

### Step 2: Implement Backend Changes

#### A. Create Model

```bash
# Create new model file
touch battle-tiles/internal/dal/model/game/game_player_record.go
```

```go
// internal/dal/model/game/game_player_record.go
package game

import "time"

const TableNameGamePlayerRecord = "game_player_record"

type GamePlayerRecord struct {
    Id             int32     `gorm:"primaryKey;column:id" json:"id"`
    CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
    UpdatedAt      time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
    BattleRecordID int32     `gorm:"column:battle_record_id;not null" json:"battle_record_id"`
    HouseGID       int32     `gorm:"column:house_gid;not null" json:"house_gid"`
    PlayerGID      int64     `gorm:"column:player_gid;not null" json:"player_gid"`
    GameAccountID  *int32    `gorm:"column:game_account_id" json:"game_account_id"`
    ScoreDelta     int32     `gorm:"column:score_delta;not null" json:"score_delta"`
    IsWinner       bool      `gorm:"column:is_winner;not null" json:"is_winner"`
    BattleAt       time.Time `gorm:"column:battle_at;not null" json:"battle_at"`
    MetaJSON       string    `gorm:"column:meta_json;type:jsonb" json:"meta_json"`
}

func (GamePlayerRecord) TableName() string { return TableNameGamePlayerRecord }
```

#### B. Update Request Model

```go
// internal/dal/req/basic_user.go
type RegisterRequest struct {
    Username        string `json:"username" binding:"required,min=3,max=50"`
    Password        string `json:"password" binding:"required,min=6,max=30"`
    GameLoginMode   string `json:"game_login_mode" binding:"required,oneof=account mobile"`
    GameAccount     string `json:"game_account" binding:"required"`
    GamePasswordMD5 string `json:"game_password_md5" binding:"required,len=32"`
}
```

#### C. Run Tests

```bash
cd battle-tiles

# Run specific test
go test ./internal/biz/basic -v -run TestRegister

# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover
```

### Step 3: Implement Frontend Changes

#### A. Update Registration Component

```typescript
// app/(auth)/register/page.tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import md5 from 'md5';

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    gameLoginMode: 'account',
    gameAccount: '',
    gamePassword: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const response = await fetch('/api/login/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...formData,
        game_password_md5: md5(formData.gamePassword),
      }),
    });
    
    // Handle response...
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
    </form>
  );
}
```

#### B. Test Frontend

```bash
cd battle-reusables

# Run development server
npm run dev

# Run tests
npm test

# Run linter
npm run lint
```

### Step 4: Test Integration

#### A. Test Registration Flow

```bash
# Using curl
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

#### B. Verify Database

```sql
-- Check user was created
SELECT * FROM basic_user WHERE username = 'testuser';

-- Check game account was created
SELECT * FROM game_account WHERE user_id = (
  SELECT id FROM basic_user WHERE username = 'testuser'
);

-- Check ctrl account was created
SELECT * FROM game_ctrl_account WHERE game_user_id IN (
  SELECT game_user_id FROM game_account WHERE user_id = (
    SELECT id FROM basic_user WHERE username = 'testuser'
  )
);
```

### Step 5: Commit Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: implement game account binding in registration

- Add GamePlayerRecord model
- Update RegisterRequest with game account fields
- Implement game account verification
- Update registration frontend component
- Add integration tests"

# Push to remote
git push origin feature/registration-enhancement
```

### Step 6: Create Pull Request

1. Go to repository on GitHub/GitLab
2. Create pull request from feature branch
3. Fill in PR template with:
   - Description of changes
   - Testing performed
   - Screenshots (if UI changes)
   - Related issues

---

## Common Development Tasks

### Add New API Endpoint

```go
// 1. Define request/response structs
type MyRequest struct {
    Field string `json:"field" binding:"required"`
}

// 2. Add service method
func (s *MyService) MyEndpoint(c *gin.Context) {
    var req MyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, ecode.ParamsFailed, err)
        return
    }
    // Implementation...
    response.Success(c, result)
}

// 3. Register route
func (s *MyService) RegisterRouter(r *gin.RouterGroup) {
    g := r.Group("/my").Use(middleware.JWTAuth())
    g.POST("/endpoint", s.MyEndpoint)
}
```

### Add Database Migration

```sql
-- Create migration file: migrations/001_add_new_field.sql
BEGIN;

ALTER TABLE my_table ADD COLUMN new_field VARCHAR(255);

COMMIT;
```

```bash
# Apply migration
psql -U your_user -d battle_tiles_dev -f migrations/001_add_new_field.sql
```

### Add Background Job

```go
// 1. Define task type
const TypeMyTask = "my:task"

// 2. Define payload
type MyTaskPayload struct {
    Field string `json:"field"`
}

// 3. Add handler
func (b *BizAsynq) HandleMyTask(ctx context.Context, t *asynq.Task) error {
    var p MyTaskPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return err
    }
    // Implementation...
    return nil
}

// 4. Enqueue task
task := asynq.NewTask(TypeMyTask, payload)
_, err := client.Enqueue(task)
```

---

## Debugging Tips

### Backend Debugging

```bash
# Enable debug logging
export LOG_LEVEL=debug

# Run with race detector
go run -race cmd/server/main.go

# Use delve debugger
dlv debug cmd/server/main.go
```

### Frontend Debugging

```bash
# Enable verbose logging
export DEBUG=*

# Check console in browser DevTools
# Add console.log statements
console.log('Debug info:', data);

# Use React DevTools extension
```

### Database Debugging

```sql
-- Check table structure
\d+ table_name

-- Check constraints
SELECT conname, contype, conrelid::regclass
FROM pg_constraint
WHERE conrelid = 'table_name'::regclass;

-- Check indexes
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'table_name';

-- Explain query
EXPLAIN ANALYZE SELECT * FROM table_name WHERE condition;
```

---

## Useful Commands

### Backend

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Generate mocks
mockgen -source=interface.go -destination=mock.go

# Update dependencies
go get -u ./...
go mod tidy
```

### Frontend

```bash
# Format code
npm run format

# Type check
npm run type-check

# Build for production
npm run build

# Analyze bundle
npm run analyze
```

### Database

```bash
# Backup database
pg_dump -U user -d database > backup.sql

# Restore database
psql -U user -d database < backup.sql

# Connect to database
psql -U user -d database

# List databases
psql -U user -l
```

---

## Troubleshooting

### Issue: Backend won't start

**Solution:**
```bash
# Check if port is already in use
lsof -i :8080

# Kill process using port
kill -9 <PID>

# Check database connection
psql -U your_user -d battle_tiles_dev -c "SELECT 1"
```

### Issue: Frontend build fails

**Solution:**
```bash
# Clear cache
rm -rf .next node_modules
npm install
npm run build
```

### Issue: Database migration fails

**Solution:**
```bash
# Check current schema
psql -U user -d database -c "\d"

# Rollback migration
psql -U user -d database -f rollback.sql

# Check for locks
SELECT * FROM pg_locks WHERE NOT granted;
```

---

## Resources

### Documentation
- [Implementation Plan](./implementation_plan.md)
- [Technical Specifications](./technical_specifications.md)
- [API Documentation](./api_documentation.md)
- [Migration Guide](./migration_guide.md)

### External Resources
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Next.js Documentation](https://nextjs.org/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

### Team Contacts
- **Backend Lead**: [contact]
- **Frontend Lead**: [contact]
- **DBA**: [contact]
- **DevOps**: [contact]

---

## Next Steps

1. **Complete environment setup** following this guide
2. **Read [Implementation Plan](./implementation_plan.md)** for big picture
3. **Review [Technical Specifications](./technical_specifications.md)** for details
4. **Check [Implementation Checklist](./implementation_checklist.md)** for tasks
5. **Start implementing** assigned features
6. **Ask questions** in team chat or stand-ups

---

Happy coding! 🚀

