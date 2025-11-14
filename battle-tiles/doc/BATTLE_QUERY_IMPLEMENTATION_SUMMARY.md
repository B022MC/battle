# 战绩查询功能实现总结

## 概述

本文档总结了战绩查询和余额查询功能的完整实现,包括数据模型修改、Repository 层、UseCase 层、Service 层和依赖注入配置。

## 实现的功能

### 1. 普通用户功能
- ✅ 查询自己的战绩列表(可按圈子筛选)
- ✅ 查询自己的余额(可查询所有圈子或指定圈子)
- ✅ 查询自己的统计数据(总局数、总分数、总费用、平均分)

### 2. 管理员功能
- ✅ 查询圈子战绩列表(可按成员筛选)
- ✅ 查询圈子成员余额列表
- ✅ 查询圈子统计数据(总局数、总分数、总费用、活跃成员数)

### 3. 超级管理员功能
- ✅ 查询店铺统计数据(总局数、总分数、总费用)

## 数据模型修改

### 1. GameMember 表
**文件**: `internal/dal/model/game/game_member.go`

**修改内容**:
- 增加 `GroupID *int32` 字段
- 唯一索引从 `(house_gid, game_id)` 改为 `(house_gid, game_id, group_id)`

**影响**: 支持同一用户在不同圈子有独立的余额记录

### 2. GameMemberWallet 表
**文件**: `internal/dal/model/game/game_member_wallet.go`

**修改内容**:
- 增加 `GroupID *int32` 字段
- 唯一索引改为 `(house_gid, member_id, group_id)`

**影响**: 支持同一成员在不同圈子有独立的钱包

### 3. GameBattleRecord 表
**文件**: `internal/dal/model/game/game_battle_record.go`

**状态**: 无需修改,已有 `GroupID` 和 `GroupName` 字段

## Repository 层实现

### 1. BattleRecordRepo
**文件**: `internal/dal/repo/game/battle_record.go`

**修改的方法**:
- `ListByPlayer`: 增加 `groupID *int32` 参数,支持按圈子筛选
- `GetPlayerStats`: 增加 `groupID *int32` 参数,支持按圈子统计

**新增的方法**:
- `GetGroupStats`: 查询圈子统计数据(总局数、总分数、总费用、活跃成员数)
- `GetHouseStats`: 查询店铺统计数据(总局数、总分数、总费用)

### 2. GameMemberRepo (新建)
**文件**: `internal/dal/repo/game/game_member.go`

**新增的方法**:
- `GetByGameIDAndGroup`: 根据 game_id 和 group_id 查询成员
- `ListByGameID`: 查询成员在所有圈子的记录
- `ListByGroup`: 查询圈子的所有成员(分页)
- `GetByID`: 根据 member_id 查询成员

### 3. WalletReadRepo
**文件**: `internal/dal/repo/game/game_wallet_read.go`

**修改的方法**:
- `Get`: 增加 `groupID *int32` 参数,支持按圈子查询
- `ListWallets`: 增加 `groupID *int32` 参数,支持按圈子筛选

**新增的方法**:
- `ListMemberBalances`: 查询成员在所有圈子的余额

## UseCase 层实现

### 1. BattleQueryUseCase (新建)
**文件**: `internal/biz/game/battle_query.go`

**实现的方法**:
- `ListMyBattles`: 用户查询自己的战绩
- `GetMyStats`: 用户查询自己的统计数据
- `ListGroupBattles`: 管理员查询圈子战绩
- `GetGroupStats`: 管理员查询圈子统计
- `GetHouseStats`: 超级管理员查询店铺统计

**权限控制**:
- 管理员接口会验证 `adminUserID` 是否是圈子的管理员
- 超级管理员接口需要额外的权限验证(TODO)

### 2. BalanceQueryUseCase (新建)
**文件**: `internal/biz/game/balance_query.go`

**实现的方法**:
- `GetMyBalances`: 用户查询自己的余额(可查询所有圈子或指定圈子)
- `ListGroupMemberBalances`: 管理员查询圈子成员余额

**数据结构**:
- `MemberBalance`: 包含余额(分)和余额(元)两个字段,方便前端使用

## Service 层实现

### BattleQueryService (新建)
**文件**: `internal/service/game/battle_query_service.go`

**实现的接口**:

#### 普通用户接口
- `POST /battle-query/my/battles`: 查询我的战绩
- `POST /battle-query/my/balances`: 查询我的余额
- `POST /battle-query/my/stats`: 查询我的统计

#### 管理员接口
- `POST /battle-query/group/battles`: 查询圈子战绩
- `POST /battle-query/group/balances`: 查询圈子成员余额
- `POST /battle-query/group/stats`: 查询圈子统计

#### 超级管理员接口
- `POST /battle-query/house/stats`: 查询店铺统计

**请求参数**:
- 所有接口都需要 JWT 认证
- 时间参数使用 Unix timestamp
- 余额参数使用元(自动转换为分)
- 支持分页(page, size)

## 依赖注入配置

### 1. Repository 层
**文件**: `internal/dal/repo/repo.go`

**新增**:
```go
game.NewGameMemberRepo,
```

### 2. UseCase 层
**文件**: `internal/biz/biz.go`

**新增**:
```go
game.NewBattleQueryUseCase,
game.NewBalanceQueryUseCase,
```

### 3. Service 层
**文件**: `internal/service/service.go`

**新增**:
```go
game.NewBattleQueryService,
```

### 4. 路由配置
**文件**: `internal/router/game_router.go`

**新增**:
- 在 `GameRouter` 结构体中添加 `battleQueryService` 字段
- 在 `InitRouter` 方法中注册路由
- 在 `NewGameRouter` 构造函数中添加参数

## 数据库迁移

### 迁移脚本
**文件**: `doc/migration_independent_group_balance.sql`

**包含内容**:
1. 数据备份
2. 添加 `group_id` 字段到 `game_member` 和 `game_member_wallet`
3. 更新唯一索引
4. 数据迁移(为现有数据填充 `group_id`)
5. 数据验证
6. 回滚脚本

**执行方式**:
```bash
psql -U your_user -d battle_tiles -f doc/migration_independent_group_balance.sql
```

## API 使用示例

### 1. 查询我的战绩
```bash
POST /battle-query/my/battles
Content-Type: application/json
Authorization: Bearer <token>

{
  "house_gid": 1,
  "group_id": 10,  // 可选,不传则查询所有圈子
  "start_time": 1609459200,  // 可选
  "end_time": 1612137600,    // 可选
  "page": 1,
  "size": 20
}
```

### 2. 查询我的余额
```bash
POST /battle-query/my/balances
Content-Type: application/json
Authorization: Bearer <token>

{
  "house_gid": 1,
  "group_id": 10  // 可选,不传则查询所有圈子
}
```

### 3. 管理员查询圈子成员余额
```bash
POST /battle-query/group/balances
Content-Type: application/json
Authorization: Bearer <token>

{
  "house_gid": 1,
  "group_id": 10,
  "min_yuan": 100,  // 可选,最小余额(元)
  "max_yuan": 1000, // 可选,最大余额(元)
  "page": 1,
  "size": 20
}
```

## 待完成的工作

### 1. 权限控制 (TODO)
- [ ] 实现超级管理员权限验证
- [ ] 添加权限检查中间件
- [ ] 完善管理员权限验证逻辑

### 2. 测试
- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 测试数据库迁移脚本

### 3. 文档
- [ ] 更新 API 文档
- [ ] 添加 Swagger 注释
- [ ] 更新用户手册

### 4. 优化
- [ ] 添加缓存机制(如果需要)
- [ ] 优化查询性能
- [ ] 添加日志记录

## 注意事项

1. **数据迁移**: 在执行数据库迁移前,请务必备份数据
2. **权限验证**: 当前管理员权限验证基于 `GameShopGroup.AdminUserID`,超级管理员权限验证需要补充
3. **余额单位**: 数据库存储使用分(int32),API 返回同时提供分和元两个字段
4. **时间参数**: 所有时间参数使用 Unix timestamp
5. **分页**: 默认 page=1, size=20,最大 size=200

## 相关文档

- [独立圈子余额设计文档](./INDEPENDENT_GROUP_BALANCE_DESIGN.md)
- [战绩查询功能实现方案](./BATTLE_RECORD_QUERY_IMPLEMENTATION.md)
- [数据库迁移脚本](./migration_independent_group_balance.sql)

