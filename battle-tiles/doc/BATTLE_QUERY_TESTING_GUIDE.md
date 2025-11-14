# 战绩查询功能测试指南

## 前置准备

### 1. 执行数据库迁移

```bash
# 连接到 PostgreSQL 数据库
psql -U your_user -d battle_tiles

# 执行迁移脚本
\i doc/migration_independent_group_balance.sql

# 验证迁移结果
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'game_member' AND column_name = 'group_id';

SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'game_member_wallet' AND column_name = 'group_id';
```

### 2. 重新生成 Wire 依赖注入代码

```bash
cd battle-tiles/cmd/go-kgin-platform
wire
```

### 3. 编译并运行后端

```bash
cd battle-tiles
go build -o bin/platform cmd/go-kgin-platform/main.go cmd/go-kgin-platform/wire_gen.go
./bin/platform -conf configs/config.yaml
```

## 测试用例

### 测试 1: 查询我的战绩

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/my/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "page": 1,
    "size": 10
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "house_gid": 1,
        "group_id": 10,
        "group_name": "测试圈子",
        "player_game_id": 12345,
        "score": 100,
        "fee": 10,
        "player_balance": 1000,
        "battle_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 50
  }
}
```

### 测试 2: 查询我的余额(所有圈子)

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/my/balances \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "house_gid": 1
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "balances": [
      {
        "member_id": 1,
        "game_id": 12345,
        "game_name": "测试用户",
        "group_id": 10,
        "group_name": "圈子A",
        "balance": 100000,
        "balance_yuan": 1000.00,
        "updated_at": "2024-01-01 10:00:00"
      },
      {
        "member_id": 2,
        "game_id": 12345,
        "game_name": "测试用户",
        "group_id": 20,
        "group_name": "圈子B",
        "balance": 50000,
        "balance_yuan": 500.00,
        "updated_at": "2024-01-01 11:00:00"
      }
    ]
  }
}
```

### 测试 3: 查询我的余额(指定圈子)

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/my/balances \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "group_id": 10
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "balances": [
      {
        "member_id": 1,
        "game_id": 12345,
        "game_name": "测试用户",
        "group_id": 10,
        "group_name": "圈子A",
        "balance": 100000,
        "balance_yuan": 1000.00,
        "updated_at": "2024-01-01 10:00:00"
      }
    ]
  }
}
```

### 测试 4: 查询我的统计数据

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/my/stats \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "group_id": 10,
    "start_time": 1704067200,
    "end_time": 1704153600
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total_games": 100,
    "total_score": 5000,
    "total_fee": 500,
    "avg_score": 50.0,
    "group_id": 10,
    "group_name": "圈子A"
  }
}
```

### 测试 5: 管理员查询圈子战绩

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/group/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "group_id": 10,
    "page": 1,
    "size": 20
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "house_gid": 1,
        "group_id": 10,
        "group_name": "圈子A",
        "player_game_id": 12345,
        "score": 100,
        "fee": 10,
        "player_balance": 1000,
        "battle_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 200
  }
}
```

### 测试 6: 管理员查询圈子成员余额

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/group/balances \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "group_id": 10,
    "min_yuan": 100,
    "max_yuan": 10000,
    "page": 1,
    "size": 20
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "member_id": 1,
        "game_id": 12345,
        "game_name": "用户A",
        "group_id": 10,
        "group_name": "圈子A",
        "balance": 100000,
        "balance_yuan": 1000.00,
        "updated_at": "2024-01-01 10:00:00"
      }
    ],
    "total": 50
  }
}
```

### 测试 7: 管理员查询圈子统计

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/group/stats \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "group_id": 10,
    "start_time": 1704067200,
    "end_time": 1704153600
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "group_id": 10,
    "group_name": "圈子A",
    "total_games": 1000,
    "total_score": 50000,
    "total_fee": 5000,
    "active_members": 50
  }
}
```

### 测试 8: 超级管理员查询店铺统计

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/house/stats \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "start_time": 1704067200,
    "end_time": 1704153600
  }'
```

**预期响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "house_gid": 1,
    "total_games": 5000,
    "total_score": 250000,
    "total_fee": 25000
  }
}
```

## 错误处理测试

### 测试 9: 无权限访问圈子

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/group/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer NON_ADMIN_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "group_id": 10,
    "page": 1,
    "size": 20
  }'
```

**预期响应**:
```json
{
  "code": 403,
  "msg": "无权限访问该圈子",
  "data": null
}
```

### 测试 10: 参数验证失败

**请求**:
```bash
curl -X POST http://localhost:8000/battle-query/my/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "page": 1,
    "size": 10
  }'
```

**预期响应**:
```json
{
  "code": 400,
  "msg": "参数错误",
  "data": null
}
```

## 性能测试

### 测试大数据量查询

```bash
# 查询大量战绩记录
curl -X POST http://localhost:8000/battle-query/my/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "house_gid": 1,
    "page": 1,
    "size": 200
  }'
```

**验证点**:
- 响应时间应该在 1 秒以内
- 数据库查询应该使用索引
- 分页功能正常工作

## 数据一致性测试

### 验证独立圈子余额

1. 创建测试数据:
```sql
-- 插入同一用户在不同圈子的记录
INSERT INTO game_member (house_gid, game_id, game_name, group_id, group_name, balance)
VALUES 
  (1, 12345, '测试用户', 10, '圈子A', 100000),
  (1, 12345, '测试用户', 20, '圈子B', 50000);

INSERT INTO game_member_wallet (house_gid, member_id, group_id, balance)
VALUES 
  (1, 1, 10, 100000),
  (1, 1, 20, 50000);
```

2. 查询验证:
```bash
# 查询所有圈子余额
curl -X POST http://localhost:8000/battle-query/my/balances \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"house_gid": 1}'
```

3. 验证结果:
- 应该返回两条记录
- 每条记录的 `group_id` 和 `balance` 应该不同

## 常见问题排查

### 1. Wire 生成失败

**错误**: `wire: no provider found for ...`

**解决方案**:
- 检查 `repo.go`、`biz.go`、`service.go` 中是否正确添加了新的 Provider
- 确保所有 `New*` 函数的签名正确
- 重新运行 `wire` 命令

### 2. 路由未注册

**错误**: `404 Not Found`

**解决方案**:
- 检查 `game_router.go` 中是否正确添加了 `battleQueryService`
- 确保在 `InitRouter` 方法中调用了 `RegisterRouter`
- 重启后端服务

### 3. 数据库查询失败

**错误**: `column "group_id" does not exist`

**解决方案**:
- 确认数据库迁移脚本已执行
- 检查表结构是否正确
- 验证 GORM 模型定义是否正确

## 总结

完成以上测试后,应该验证:
- ✅ 所有 API 接口正常工作
- ✅ 权限控制正确
- ✅ 数据查询准确
- ✅ 分页功能正常
- ✅ 错误处理完善
- ✅ 性能满足要求

