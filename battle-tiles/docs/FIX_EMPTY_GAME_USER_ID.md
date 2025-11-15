# 修复空的 game_user_id 字段

## 问题描述

在修复代码之前注册的用户，他们的 `game_account` 表中的 `game_user_id` 字段为空。这导致用户无法查询自己的战绩。

## 根本原因

原来的 `BindSingle` 方法只调用了 `Verify()` 来验证账号密码，但没有获取游戏服务器返回的用户 ID。

## 修复方案

### 1. 代码修复（已完成）

修改了 `battle-tiles/internal/biz/game/game_account.go` 中的 `BindSingle` 方法：

- 改为调用 `ProbeLoginWithInfo()` 而不是 `Verify()`
- 获取游戏服务器返回的 `UserLogonInfo`
- 将 `info.UserID` 保存到 `game_user_id` 字段

### 2. 历史数据修复

#### 方式一：使用 API 接口（推荐）

**接口信息**：
- 路径：`POST /game/accounts/fix-empty-game-user-id`
- 认证：需要超级管理员权限
- 请求体：无

**使用步骤**：

1. 使用超级管理员账号登录获取 token
2. 调用接口：

```bash
curl -X POST http://localhost:8080/game/accounts/fix-empty-game-user-id \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

**响应示例**：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "fixed": 10,
    "failed": 2
  }
}
```

- `fixed`：成功修复的记录数
- `failed`：修复失败的记录数

#### 方式二：手动 SQL 查询

查看所有 `game_user_id` 为空的记录：

```sql
SELECT 
    ga.id,
    ga.user_id,
    ga.account,
    ga.nickname,
    bu.username,
    bu.nick_name,
    ga.created_at
FROM game_account ga
LEFT JOIN basic_user bu ON ga.user_id = bu.id
WHERE ga.game_user_id = '' OR ga.game_user_id IS NULL
ORDER BY ga.created_at DESC;
```

#### 方式三：让用户重新绑定

1. 用户登录系统
2. 进入账号设置页面
3. 解绑现有的游戏账号
4. 重新绑定游戏账号

新绑定时，系统会自动获取并保存 `game_user_id`。

## 修复失败的原因

如果某些记录修复失败，可能的原因包括：

1. **游戏账号已被删除或禁用** - 游戏服务器无法验证该账号
2. **游戏服务器连接失败** - 无法连接到游戏服务器
3. **账号密码错误** - 存储的密码与游戏服务器不匹配

对于这些失败的记录，建议：
- 联系用户确认账号是否仍然有效
- 如果账号无效，可以删除该记录
- 如果账号有效，让用户重新绑定

## 修改的文件

1. `battle-tiles/internal/biz/game/game_account.go`
   - 添加 `FixEmptyGameUserID` 方法

2. `battle-tiles/internal/dal/repo/game/game_account.go`
   - 添加 `ListByGameUserIDEmpty` 方法
   - 添加 `Update` 方法

3. `battle-tiles/internal/service/game/game_account.go`
   - 添加 `FixEmptyGameUserID` 处理函数
   - 添加管理员路由

4. `battle-tiles/pkg/plugin/middleware/rbac.go`
   - 添加 `AdminOnly` 中间件

## 验证修复结果

修复完成后，可以验证：

```sql
-- 查看修复后的记录
SELECT 
    ga.id,
    ga.user_id,
    ga.account,
    ga.game_user_id,
    bu.username
FROM game_account ga
LEFT JOIN basic_user bu ON ga.user_id = bu.id
WHERE ga.game_user_id != '' AND ga.game_user_id IS NOT NULL
ORDER BY ga.updated_at DESC
LIMIT 10;
```

用户应该能够正常查询自己的战绩。

