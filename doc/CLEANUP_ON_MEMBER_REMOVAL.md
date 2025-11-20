# 用户被踢出圈子时的关联记录清理机制

## 业务规则

**一个用户只能同时在一个店铺的一个圈子下**

当用户被踢出圈子时，需要清理所有相关联的记录，确保数据一致性。

## 关联表及清理逻辑

### 1. **game_shop_group_member** - 圈子成员表
```sql
-- 用户在平台圈子中的成员关系
DELETE FROM game_shop_group_member 
WHERE group_id = ? AND user_id = ?
```
**作用**：移除用户在圈子中的成员身份

### 2. **game_member** - 游戏成员表
```sql
-- 用户游戏账号在店铺中的成员记录
DELETE FROM game_member 
WHERE house_gid = ? AND game_id = ?
```
**作用**：移除游戏账号在店铺中的记录（game_id 对应 game_account.id）

### 3. **game_account_house** - 账号店铺绑定表
```sql
-- 游戏账号与店铺的绑定关系
DELETE FROM game_account_house 
WHERE game_account_id = ?
```
**作用**：解除游戏账号与店铺的绑定（一个账号只能绑定一个店铺）

## 数据流

```
用户被踢出圈子
    ↓
1. 删除 game_shop_group_member 记录
   (user_id, group_id)
    ↓
2. 查询用户的游戏账号
   game_account WHERE user_id = ?
    ↓
3. 删除 game_member 记录
   (game_id = game_account.id, house_gid)
    ↓
4. 删除 game_account_house 记录
   (game_account_id = game_account.id)
```

## 代码实现

### Repo 层新增方法

#### `/internal/dal/repo/game/game_member.go`
```go
// DeleteByGameID 删除成员（通过 game_id 删除）
DeleteByGameID(ctx context.Context, houseGID int32, gameID int32) error
```

#### `/internal/dal/repo/game/game_account_house.go`
```go
// DeleteByGameAccountID 删除账号店铺绑定（通过 game_account_id 删除）
DeleteByGameAccountID(ctx context.Context, gameAccountID int32) error
```

### UseCase 层修改

#### `/internal/biz/game/shop_group.go`

**修改前**：
```go
func (uc *ShopGroupUseCase) RemoveMemberFromGroup(...) error {
    // 只删除圈子成员记录
    return uc.memberRepo.RemoveMember(ctx, groupID, userID)
}
```

**修改后**：
```go
func (uc *ShopGroupUseCase) RemoveMemberFromGroup(...) error {
    // 1. 删除圈子成员记录
    uc.memberRepo.RemoveMember(ctx, groupID, userID)
    
    // 2. 查询用户的游戏账号
    gameAccount := uc.accountRepo.GetOneByUser(ctx, userID)
    
    // 3. 删除游戏成员记录
    uc.gameMemberRepo.DeleteByGameID(ctx, houseGID, gameAccount.Id)
    
    // 4. 删除账号店铺绑定
    uc.accountHouseRepo.DeleteByGameAccountID(ctx, gameAccount.Id)
}
```

## 错误处理

1. **用户未绑定游戏账号**：只删除圈子成员记录即可
2. **删除关联记录失败**：记录警告日志，但不阻塞主流程

```go
if err == gorm.ErrRecordNotFound {
    uc.log.Infof("用户没有绑定游戏账号，跳过清理关联记录")
    return nil
}
```

## 测试场景

### 场景1：正常用户被踢出
```
初始状态：
- user_id = 7
- game_account.id = 6 (account = "1106162940")
- game_member: game_id = 6, house_gid = 100
- game_account_house: game_account_id = 6, house_gid = 100
- game_shop_group_member: user_id = 7, group_id = 5

踢出后：
- ✅ game_shop_group_member 记录被删除
- ✅ game_member 记录被删除
- ✅ game_account_house 记录被删除
- ✅ game_account 保留（用户仍可重新绑定）
```

### 场景2：未绑定游戏账号的用户被踢出
```
初始状态：
- user_id = 8 (没有绑定游戏账号)
- game_shop_group_member: user_id = 8, group_id = 5

踢出后：
- ✅ game_shop_group_member 记录被删除
- ✅ 跳过游戏账号相关清理
```

## 数据库约束

为确保"一个账号只能绑定一个店铺"的业务规则，添加了唯一约束：

```sql
-- /migrations/20251121_add_unique_game_account_constraint.sql
ALTER TABLE game_account_house 
ADD CONSTRAINT uk_game_account_id_unique UNIQUE (game_account_id);
```

## 相关API

### 踢出成员接口

#### 平台侧（店铺管理）
```
POST /api/shops/members/kick
Body: { "house_gid": 100, "member_id": 123 }
```

#### 圈子侧（圈主管理）
```
POST /api/groups/members/remove
Body: { "group_id": 5, "user_id": 7 }
```

## 日志

成功踢出用户时的日志：
```
[INFO] 成功移除用户 7 并清理关联记录
  - 删除 game_shop_group_member: user_id=7, group_id=5
  - 删除 game_member: game_id=6, house_gid=100
  - 删除 game_account_house: game_account_id=6
```

## 注意事项

1. **事务处理**：建议在同一事务中执行所有删除操作，确保数据一致性
2. **软删除 vs 硬删除**：当前使用硬删除（DELETE），如需软删除需修改实现
3. **级联删除**：数据库层面没有设置 CASCADE，由应用层控制
4. **重新加入**：用户被踢出后可以重新申请加入圈子，此时会重新创建所有关联记录

## 总结

通过在 `RemoveMemberFromGroup` 方法中添加关联记录清理逻辑，确保了：
- ✅ 用户被踢出圈子时，所有相关数据都被清理
- ✅ 符合"一个用户只能在一个店铺的一个圈子"的业务规则
- ✅ 数据一致性得到保证
- ✅ 用户可以重新加入其他圈子
