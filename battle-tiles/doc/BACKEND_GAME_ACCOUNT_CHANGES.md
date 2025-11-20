# 后端游戏账号系统改造总结

## 📋 改造目标

实现游戏账号入圈机制，解耦用户和游戏账号的强绑定关系，支持以下场景：
1. 游戏内玩家可以在未注册平台账号的情况下申请入圈
2. 店铺管理员审批后，游戏账号直接绑定到管理员的圈子
3. 用户通过游戏账号反向查询圈子和战绩
4. 用户绑定游戏账号时，如果游戏账号已存在则直接关联

---

## 🗄️ 数据库变更

### 1. 修改 `game_account` 表
- **变更**: `user_id` 字段改为可选（nullable）
- **原因**: 允许游戏账号在未绑定用户的情况下存在
- **迁移**: 已在 `migrations/20251120_game_account_group.sql` 中实现

### 2. 新增 `game_account_group` 表
- **作用**: 存储游戏账号与圈子的关系
- **核心字段**:
  - `game_account_id`: 游戏账号ID
  - `house_gid`: 店铺GID
  - `group_id`: 圈子ID
  - `admin_user_id`: 圈主用户ID
  - `approved_by_user_id`: 审批人用户ID
  - `status`: 状态（active/inactive）
- **唯一约束**: `(game_account_id, house_gid)` - 一个游戏账号在一个店铺只能属于一个圈子

### 3. 修改 `game_battle_record` 表
- **变更**: `player_id` 字段改为可选（nullable）
- **原因**: 战绩记录可能属于未注册用户的游戏账号

---

## 🔧 代码变更

### 1. 新增模型和仓库

#### `internal/dal/model/game/game_account_group.go`
```go
type GameAccountGroup struct {
    Id                int32
    GameAccountID     int32
    HouseGID          int32
    GroupID           int32
    GroupName         string
    AdminUserID       int32
    ApprovedByUserID  int32
    Status            string
    JoinedAt          time.Time
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

#### `internal/dal/repo/game/game_account_group.go`
提供以下方法：
- `Create`: 创建游戏账号圈子关系
- `GetByGameAccountAndHouse`: 查询游戏账号在某店铺的圈子
- `ListByGameAccount`: 查询游戏账号的所有圈子
- `ListByHouse`: 查询店铺的所有游戏账号关系
- `ListByGroup`: 查询圈子的所有游戏账号
- `UpdateStatus`: 更新状态
- `Delete`: 删除关系
- `ExistsByGameAccountAndHouse`: 检查游戏账号是否已在圈子中

### 2. 新增业务逻辑层

#### `internal/biz/game/game_account_group.go`
核心方法：
- `FindOrCreateGameAccount`: 根据游戏用户ID查找或创建游戏账号
- `EnsureGroupForAdmin`: 确保管理员有圈子
- `AddGameAccountToGroup`: 将游戏账号加入圈子
- `RemoveGameAccountFromGroup`: 将游戏账号从圈子移除
- `ListGroupsByUser`: 用户反向查询圈子（用户 → 游戏账号 → 圈子）

### 3. 修改现有业务逻辑

#### `internal/biz/game/game_account.go`
**修改 `BindSingle` 方法**:
```go
// 检查游戏账号是否已存在
existingAccount, err := uc.accRepo.GetByGameUserID(ctx, gameUserID)
if err == nil && existingAccount != nil {
    if existingAccount.UserID != nil {
        return nil, errors.New("游戏账号已被其他用户绑定")
    }
    // 游戏账号存在但未绑定用户，更新绑定
    existingAccount.UserID = &userID
    if err := uc.accRepo.Update(ctx, existingAccount); err != nil {
        return nil, err
    }
    return existingAccount, nil
}
// 否则创建新账号...
```

#### `internal/service/game/game_shop_application.go`
**修改 `respondGameApplication` 方法**:
- 在审批通过后调用 `handleGameAccountJoinGroup`
- 实现游戏账号入圈流程：
  1. 查找或创建游戏账号
  2. 确保管理员有圈子
  3. 将游戏账号加入圈子

#### `internal/service/game/shop_group.go`
**修改 `ListMyGroups` 方法**:
- 使用 `accountGroupUC.ListGroupsByUser` 进行反向查询
- 用户 → 游戏账号 → 圈子

#### `internal/biz/game/battle_record.go`
**修改 `buildPlayerGroupMapping` 方法**:
- 处理 `account.UserID` 为指针类型的情况
- 允许 `player_id` 为 NULL

#### `internal/biz/game/game_ctrl_account.go`
**修改中控账号创建逻辑**:
- 将 `UserID` 改为指针类型

### 4. 依赖注入配置

#### `internal/biz/biz.go`
```go
game.NewGameAccountGroupUseCase, // 新增
```

#### `internal/dal/repo/repo.go`
```go
game.NewGameAccountGroupRepo, // 新增
```

---

## 🔄 核心流程

### 流程 1: 游戏内申请 + 管理员审批

```
1. 玩家在游戏内发起申请
   ↓
2. 后端接收申请（内存队列）
   ↓
3. 管理员点击"通过"
   ↓
4. 查找或创建游戏账号
   - SELECT * FROM game_account WHERE game_user_id = ?
   - 如果不存在，INSERT INTO game_account (user_id=NULL, ...)
   ↓
5. 获取或创建管理员的圈子
   - SELECT * FROM game_shop_group WHERE house_gid = ? AND admin_user_id = ?
   - 如果不存在，创建圈子
   ↓
6. 游戏账号入圈
   - INSERT INTO game_account_group (game_account_id, house_gid, group_id, ...)
   ↓
7. 调用游戏API拉入圈子
   ↓
8. 记录操作日志
```

### 流程 2: 用户绑定游戏账号

```
1. 用户输入游戏账号和密码
   ↓
2. 验证游戏账号（探活）
   ↓
3. 检查游戏账号是否已存在
   - 如果存在且未绑定用户 → 更新 user_id
   - 如果存在且已绑定其他用户 → 报错
   - 如果不存在 → 创建新账号
   ↓
4. 返回绑定结果
```

### 流程 3: 用户查询圈子（反向查询）

```
1. 根据用户ID查询游戏账号
   - SELECT * FROM game_account WHERE user_id = ?
   ↓
2. 根据游戏账号ID查询圈子关系
   - SELECT * FROM game_account_group WHERE game_account_id = ?
   ↓
3. 返回圈子列表
```

### 流程 4: 用户查询战绩（反向查询）

```
1. 根据用户ID查询游戏账号
   - SELECT * FROM game_account WHERE user_id = ?
   ↓
2. 根据游戏账号查询战绩
   - SELECT * FROM game_battle_record WHERE player_game_id = ?
   ↓
3. 返回战绩列表
```

---

## ✅ 测试建议

### 1. 数据库迁移测试
```bash
psql -U postgres -d battle_db -f migrations/20251120_game_account_group.sql
```

### 2. 功能测试

#### 测试用例 1: 游戏内申请入圈
1. 未注册用户在游戏内申请
2. 管理员审批通过
3. 验证 `game_account` 表中创建了记录（`user_id` 为 NULL）
4. 验证 `game_account_group` 表中创建了关系

#### 测试用例 2: 用户绑定已存在的游戏账号
1. 游戏账号已通过申请入圈（`user_id` 为 NULL）
2. 用户注册并绑定该游戏账号
3. 验证 `game_account.user_id` 被更新

#### 测试用例 3: 用户查询圈子
1. 用户绑定游戏账号
2. 游戏账号已入圈
3. 调用 `/api/groups/my/list`
4. 验证返回正确的圈子列表

#### 测试用例 4: 用户查询战绩
1. 用户绑定游戏账号
2. 游戏账号有战绩记录
3. 调用 `/api/battle-query/my/battles`
4. 验证返回正确的战绩列表

---

## 🚨 注意事项

### 1. 数据一致性
- 确保 `game_account_group` 的唯一约束生效
- 游戏账号在一个店铺只能属于一个圈子

### 2. 兼容性
- 旧的战绩记录可能有 `player_id`，新的可能为 NULL
- 查询时需要同时支持两种情况

### 3. 性能优化
- `game_account_group` 表的索引已创建
- 反向查询需要两次数据库查询，考虑缓存优化

### 4. 错误处理
- 游戏账号不存在时的处理
- 用户未绑定游戏账号时的提示
- 游戏账号已被其他用户绑定的错误提示

---

## 📝 API 变更

### 无需修改的 API
- `/api/game/accounts/verify` - 验证游戏账号
- `/api/game/accounts` - 绑定游戏账号（逻辑已更新）
- `/api/game/accounts/me` - 查询我的游戏账号
- `/api/game/accounts/me/houses` - 查询我的店铺
- `/api/battle-query/my/battles` - 查询我的战绩（逻辑已更新）
- `/api/groups/my/list` - 查询我的圈子（逻辑已更新）

### 内部逻辑变更的 API
- `/api/shops/game-applications/approve` - 审批通过（新增游戏账号入圈逻辑）
- `/api/shops/game-applications/reject` - 审批拒绝

---

## 🎯 完成状态

- ✅ 数据库迁移脚本
- ✅ 模型和仓库层
- ✅ 业务逻辑层
- ✅ 服务层
- ✅ 依赖注入配置
- ✅ 编译通过

---

## 📚 相关文档

- [游戏账号系统重新设计方案 V2](./GAME_ACCOUNT_REDESIGN_V2.md)
- [游戏账号简明指南](./GAME_ACCOUNT_SIMPLE_GUIDE.md)
- [数据库迁移指南](./GAME_ACCOUNT_MIGRATION_GUIDE.md)
- [完整数据库结构](./public.sql)

