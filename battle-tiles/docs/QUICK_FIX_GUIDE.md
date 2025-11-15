# 快速修复指南

## 一句话总结

使用管理员账号调用 API 接口 `POST /game/accounts/fix-empty-game-user-id` 来自动修复所有空的 `game_user_id` 字段。

## 快速步骤

### 1. 编译并启动后端服务

```bash
cd battle-tiles
go build -o battle-tiles ./cmd/main.go
./battle-tiles
```

### 2. 使用 Postman 或 curl 调用修复接口

**获取 token**（使用超级管理员账号）：

```bash
curl -X POST http://localhost:8080/basic/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

响应中会包含 `token`。

**调用修复接口**：

```bash
curl -X POST http://localhost:8080/game/accounts/fix-empty-game-user-id \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 3. 查看修复结果

响应示例：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "fixed": 5,
    "failed": 0
  }
}
```

## 修复前后对比

### 修复前

```sql
SELECT * FROM game_account WHERE game_user_id = '' OR game_user_id IS NULL;
```

结果：有多条记录的 `game_user_id` 为空

### 修复后

```sql
SELECT * FROM game_account WHERE game_user_id = '' OR game_user_id IS NULL;
```

结果：没有记录（或只有修复失败的记录）

## 常见问题

### Q: 修复失败了怎么办？

A: 检查日志中的错误信息。常见原因：
- 游戏账号已被删除
- 游戏服务器连接失败
- 账号密码不匹配

### Q: 修复后用户还是看不到战绩？

A: 
1. 确认 `game_user_id` 已经被填充
2. 让用户重新登录
3. 检查 `game_battle_record` 表中是否有该用户的战绩记录

### Q: 可以只修复某个用户吗？

A: 目前接口会修复所有空的 `game_user_id`。如果需要修复特定用户，可以：
1. 让用户重新绑定游戏账号
2. 或者手动更新数据库

## 相关文件

- 详细文档：`battle-tiles/docs/FIX_EMPTY_GAME_USER_ID.md`
- 代码修改：
  - `battle-tiles/internal/biz/game/game_account.go`
  - `battle-tiles/internal/dal/repo/game/game_account.go`
  - `battle-tiles/internal/service/game/game_account.go`
  - `battle-tiles/pkg/plugin/middleware/rbac.go`

