# 战绩查询功能实现方案

## 📋 功能概述

实现用户战绩查询功能,支持三种角色:
1. **普通用户**: 查询自己的战绩和余额
2. **店铺管理员**: 查询所管理圈子内用户的战绩和余额
3. **超级管理员**: 查询整个店铺的统计数据

## 🎯 实现目标

### 第一阶段 (核心功能)
- [x] 数据模型修改 (支持独立圈子余额)
- [ ] Repository 层修改
- [ ] UseCase 层实现
- [ ] Service 层实现
- [ ] API 接口实现

### 第二阶段 (统计功能)
- [ ] 用户统计数据 (今日/昨日/本周)
- [ ] 管理员圈子统计
- [ ] 超级管理员店铺统计

## 📊 API 设计

### 1. 普通用户 API

#### 1.1 查询自己的战绩列表
```
GET /api/v1/members/battles/my
参数:
  - group_id (可选): 圈子ID,不传则查询所有圈子
  - page: 页码,默认1
  - page_size: 每页数量,默认20,最大100
  - start_date (可选): 开始日期 YYYY-MM-DD
  - end_date (可选): 结束日期 YYYY-MM-DD

响应:
{
  "code": 0,
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "battle_at": "2025-11-15T10:30:00Z",
        "room_uid": 12345,
        "kind_id": 1,
        "base_score": 5,
        "group_id": 10,
        "group_name": "VIP圈",
        "score": 200,
        "fee": 50,
        "player_balance": 1500
      }
    ]
  }
}
```

#### 1.2 查询自己的余额
```
GET /api/v1/members/balance/my
参数:
  - group_id (可选): 圈子ID,不传则查询所有圈子

响应:
{
  "code": 0,
  "data": {
    "balances": [
      {
        "group_id": 10,
        "group_name": "VIP圈",
        "balance": 1500,
        "balance_yuan": 15.00
      },
      {
        "group_id": 11,
        "group_name": "普通圈",
        "balance": 500,
        "balance_yuan": 5.00
      }
    ],
    "total_balance": 2000,
    "total_balance_yuan": 20.00
  }
}
```

#### 1.3 查询自己的统计数据
```
GET /api/v1/members/stats/my
参数:
  - group_id (可选): 圈子ID
  - period: 统计周期 (today/yesterday/this_week/last_week)

响应:
{
  "code": 0,
  "data": {
    "period": "today",
    "group_id": 10,
    "group_name": "VIP圈",
    "total_games": 50,
    "total_score": 500,
    "total_fee": 100,
    "avg_score": 10.0,
    "win_rate": 0.6
  }
}
```

### 2. 店铺管理员 API

#### 2.1 查询圈子成员余额列表
```
GET /api/v1/shops/groups/:group_id/members/balances
参数:
  - page: 页码,默认1
  - page_size: 每页数量,默认20,最大100
  - min_balance (可选): 最小余额筛选
  - max_balance (可选): 最大余额筛选
  - sort: 排序方式 (balance_desc/balance_asc/updated_desc)

响应:
{
  "code": 0,
  "data": {
    "total": 50,
    "list": [
      {
        "member_id": 1,
        "game_id": 12345,
        "game_name": "张三",
        "balance": 1500,
        "balance_yuan": 15.00,
        "updated_at": "2025-11-15T10:30:00Z"
      }
    ]
  }
}
```

#### 2.2 查询圈子成员战绩列表
```
GET /api/v1/shops/groups/:group_id/members/battles
参数:
  - game_id (可选): 指定成员的游戏ID
  - page: 页码,默认1
  - page_size: 每页数量,默认20,最大100
  - start_date (可选): 开始日期
  - end_date (可选): 结束日期

响应:
{
  "code": 0,
  "data": {
    "total": 200,
    "list": [
      {
        "id": 1,
        "battle_at": "2025-11-15T10:30:00Z",
        "player_game_id": 12345,
        "player_game_name": "张三",
        "score": 200,
        "fee": 50,
        "player_balance": 1500
      }
    ]
  }
}
```

#### 2.3 查询圈子统计数据
```
GET /api/v1/shops/groups/:group_id/stats
参数:
  - period: 统计周期 (today/yesterday/this_week/last_week)

响应:
{
  "code": 0,
  "data": {
    "period": "today",
    "group_id": 10,
    "group_name": "VIP圈",
    "total_members": 50,
    "total_games": 500,
    "total_score": 5000,
    "total_fee": 1000,
    "active_members": 30
  }
}
```

### 3. 超级管理员 API

#### 3.1 查询店铺统计数据
```
GET /api/v1/shops/:house_gid/stats
参数:
  - period: 统计周期 (today/yesterday/this_week/last_week)

响应:
{
  "code": 0,
  "data": {
    "period": "today",
    "house_gid": 100,
    "total_groups": 5,
    "total_members": 200,
    "total_games": 2000,
    "total_score": 20000,
    "total_fee": 4000,
    "groups": [
      {
        "group_id": 10,
        "group_name": "VIP圈",
        "total_games": 500,
        "total_fee": 1000
      }
    ]
  }
}
```

## 🔧 Repository 层修改

### 1. BattleRecordRepo 接口修改

```go
type BattleRecordRepo interface {
    // 现有方法
    SaveBatch(ctx context.Context, list []*model.GameBattleRecord) error
    SaveBatchWithDedup(ctx context.Context, list []*model.GameBattleRecord) (int, error)
    List(ctx context.Context, houseGID int32, groupID *int32, gameID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)
    
    // 修改: 增加 groupID 参数
    ListByPlayer(ctx context.Context, houseGID int32, playerGameID int32, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)
    
    // 修改: 增加 groupID 参数
    GetPlayerStats(ctx context.Context, houseGID int32, playerGameID int32, groupID *int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)
    
    // 新增: 查询圈子统计
    GetGroupStats(ctx context.Context, houseGID int32, groupID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, activeMembers int64, err error)
    
    // 新增: 查询店铺统计
    GetHouseStats(ctx context.Context, houseGID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)
}
```

### 2. WalletReadRepo 接口修改

```go
type WalletReadRepo interface {
    // 修改: 增加 groupID 参数
    Get(ctx context.Context, houseGID, memberID int32, groupID *int32) (*model.GameMemberWallet, error)
    
    // 修改: 增加 groupID 参数
    ListWallets(ctx context.Context, houseGID int32, groupID *int32, min, max *int32, hasCustomLimit *bool, page, size int32) ([]*model.GameMemberWallet, int64, error)
    
    // 新增: 查询成员在所有圈子的余额
    ListMemberBalances(ctx context.Context, houseGID int32, memberID int32) ([]*model.GameMemberWallet, error)
    
    // 现有方法保持不变
    ListLedger(ctx context.Context, houseGID int32, memberID *int32, tp *int32, start, end time.Time, page, size int32) ([]*model.GameWalletLedger, int64, error)
    ListWalletsByMembers(ctx context.Context, houseGID int32, memberIDs []int32, min, max *int32, page, size int32) ([]*model.GameMemberWallet, int64, error)
}
```

### 3. GameMemberRepo 接口新增

```go
type GameMemberRepo interface {
    // 根据 game_id 和 group_id 查询成员
    GetByGameIDAndGroup(ctx context.Context, houseGID int32, gameID int32, groupID *int32) (*model.GameMember, error)
    
    // 查询成员在所有圈子的记录
    ListByGameID(ctx context.Context, houseGID int32, gameID int32) ([]*model.GameMember, error)
    
    // 查询圈子的所有成员
    ListByGroup(ctx context.Context, houseGID int32, groupID int32, page, size int32) ([]*model.GameMember, int64, error)
}
```

## 💼 UseCase 层实现

### 1. BattleRecordUseCase

```go
type BattleRecordUseCase struct {
    repo BattleRecordRepo
    memberRepo GameMemberRepo
    walletRepo WalletReadRepo
    log *log.Helper
}

// 用户查询自己的战绩
func (uc *BattleRecordUseCase) ListMyBattles(ctx context.Context, userID int32, houseGID int32, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)

// 用户查询自己的统计
func (uc *BattleRecordUseCase) GetMyStats(ctx context.Context, userID int32, houseGID int32, groupID *int32, start, end *time.Time) (*BattleStats, error)

// 管理员查询圈子战绩
func (uc *BattleRecordUseCase) ListGroupBattles(ctx context.Context, adminUserID int32, houseGID int32, groupID int32, playerGameID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)

// 管理员查询圈子统计
func (uc *BattleRecordUseCase) GetGroupStats(ctx context.Context, adminUserID int32, houseGID int32, groupID int32, start, end *time.Time) (*GroupStats, error)

// 超级管理员查询店铺统计
func (uc *BattleRecordUseCase) GetHouseStats(ctx context.Context, superAdminUserID int32, houseGID int32, start, end *time.Time) (*HouseStats, error)
}
```

### 2. MemberBalanceUseCase

```go
type MemberBalanceUseCase struct {
    memberRepo GameMemberRepo
    walletRepo WalletReadRepo
    groupRepo GameShopGroupRepo
    log *log.Helper
}

// 用户查询自己的余额
func (uc *MemberBalanceUseCase) GetMyBalances(ctx context.Context, userID int32, houseGID int32, groupID *int32) ([]*MemberBalance, error)

// 管理员查询圈子成员余额
func (uc *MemberBalanceUseCase) ListGroupMemberBalances(ctx context.Context, adminUserID int32, houseGID int32, groupID int32, min, max *int32, page, size int32) ([]*MemberBalance, int64, error)
}
```

## 🌐 Service 层实现

### 1. HTTP 路由注册

```go
// internal/service/game/battle_record.go
func (s *BattleRecordService) RegisterRoutes(r *gin.RouterGroup) {
    // 普通用户路由
    members := r.Group("/members")
    {
        members.GET("/battles/my", s.ListMyBattles)
        members.GET("/balance/my", s.GetMyBalances)
        members.GET("/stats/my", s.GetMyStats)
    }
    
    // 店铺管理员路由
    shops := r.Group("/shops")
    {
        groups := shops.Group("/groups/:group_id")
        {
            groups.GET("/members/balances", s.ListGroupMemberBalances)
            groups.GET("/members/battles", s.ListGroupMemberBattles)
            groups.GET("/stats", s.GetGroupStats)
        }
    }
    
    // 超级管理员路由
    admin := r.Group("/admin/shops/:house_gid")
    {
        admin.GET("/stats", s.GetHouseStats)
    }
}
```

## ✅ 实现检查清单

### Repository 层
- [ ] 修改 `BattleRecordRepo.ListByPlayer` 增加 `groupID` 参数
- [ ] 修改 `BattleRecordRepo.GetPlayerStats` 增加 `groupID` 参数
- [ ] 新增 `BattleRecordRepo.GetGroupStats` 方法
- [ ] 新增 `BattleRecordRepo.GetHouseStats` 方法
- [ ] 修改 `WalletReadRepo.Get` 增加 `groupID` 参数
- [ ] 修改 `WalletReadRepo.ListWallets` 增加 `groupID` 参数
- [ ] 新增 `WalletReadRepo.ListMemberBalances` 方法
- [ ] 新增 `GameMemberRepo` 接口和实现

### UseCase 层
- [ ] 实现 `BattleRecordUseCase.ListMyBattles`
- [ ] 实现 `BattleRecordUseCase.GetMyStats`
- [ ] 实现 `BattleRecordUseCase.ListGroupBattles`
- [ ] 实现 `BattleRecordUseCase.GetGroupStats`
- [ ] 实现 `BattleRecordUseCase.GetHouseStats`
- [ ] 实现 `MemberBalanceUseCase.GetMyBalances`
- [ ] 实现 `MemberBalanceUseCase.ListGroupMemberBalances`

### Service 层
- [ ] 实现 `BattleRecordService.ListMyBattles` HTTP handler
- [ ] 实现 `BattleRecordService.GetMyBalances` HTTP handler
- [ ] 实现 `BattleRecordService.GetMyStats` HTTP handler
- [ ] 实现 `BattleRecordService.ListGroupMemberBalances` HTTP handler
- [ ] 实现 `BattleRecordService.ListGroupMemberBattles` HTTP handler
- [ ] 实现 `BattleRecordService.GetGroupStats` HTTP handler
- [ ] 实现 `BattleRecordService.GetHouseStats` HTTP handler

### 权限控制
- [ ] 实现普通用户权限检查中间件
- [ ] 实现店铺管理员权限检查中间件
- [ ] 实现超级管理员权限检查中间件

## 📝 注意事项

1. **数据库迁移**: 在实现功能前,必须先执行 `migration_independent_group_balance.sql`
2. **权限控制**: 所有 API 都需要进行权限检查
3. **性能优化**: 使用索引优化查询性能
4. **错误处理**: 统一错误处理和返回格式
5. **日志记录**: 记录关键操作日志

## 🔗 相关文档

- [独立圈子余额设计文档](./INDEPENDENT_GROUP_BALANCE_DESIGN.md)
- [数据库迁移脚本](./migration_independent_group_balance.sql)

