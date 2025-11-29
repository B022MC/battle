# 实现圈子战绩查询接口

## 📊 现状分析

✅ **已完成部分**：
- 数据表 `game_battle_record` 已创建并有索引
- Repo 层查询方法已实现：`ListByPlayer()`
- 其他战绩查询接口已正常工作（我的战绩、统计等）

⚠️ **待实现**：
- Service 层 `ListGroupBattles` 接口（目前标记为 TODO）
- Biz 层 `ListGroupBattles` 业务逻辑

## 🔧 实现步骤

### 第 1 步：添加 Biz 层方法

**文件**：`/internal/biz/game/battle_record.go`

在 `BattleRecordUseCase` 中添加：

```go
// ListGroupBattles 查询圈子成员战绩（管理员）
func (uc *BattleRecordUseCase) ListGroupBattles(
	ctx context.Context,
	houseGID int32,
	groupID int32,
	playerGameID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	// 如果指定了玩家ID，查询该玩家的战绩
	if playerGameID != nil && *playerGameID > 0 {
		return uc.repo.ListByPlayer(ctx, houseGID, *playerGameID, &groupID, start, end, page, size)
	}
	
	// 否则查询整个圈子的战绩
	return uc.repo.List(ctx, houseGID, &groupID, nil, start, end, page, size)
}
```

### 第 2 步：实现 Service 层接口

**文件**：`/internal/service/game/battle_record.go`

将第 201-203 行的 TODO 实现替换为：

```go
// ListGroupBattlesRequest 查询圈子战绩请求
type ListGroupBattlesRequest struct {
	HouseGID     int32  `json:"house_gid" binding:"required"`
	GroupID      int32  `json:"group_id" binding:"required"`
	PlayerGameID *int32 `json:"player_game_id"` // 可选：指定玩家
	StartTime    *int64 `json:"start_time"`     // Unix timestamp
	EndTime      *int64 `json:"end_time"`       // Unix timestamp
	Page         int32  `json:"page"`
	Size         int32  `json:"size"`
}

// ListGroupBattles 查询圈子成员战绩
// @Summary      查询圈子战绩（管理员）
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body ListGroupBattlesRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battle-query/group/battles [post]
func (s *BattleRecordService) ListGroupBattles(c *gin.Context) {
	var in ListGroupBattlesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 从 JWT 中获取用户信息
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// TODO: 可选 - 验证管理员是否属于该圈子（根据业务需求）
	// 当前已经通过权限中间件验证了 battles:view 权限

	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 设置默认分页参数
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 || in.Size > 100 {
		in.Size = 10
	}

	// 调用业务层查询战绩
	list, total, err := s.uc.ListGroupBattles(
		c,
		in.HouseGID,
		in.GroupID,
		in.PlayerGameID,
		start,
		end,
		in.Page,
		in.Size,
	)
	if err != nil {
		response.Fail(c, ecode.ServerError, err)
		return
	}

	// 返回结果
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
```

## 🧪 测试

### 1. 测试查询指定玩家战绩

```bash
curl -X POST http://localhost:8000/battle-query/group/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "house_gid": 1,
    "group_id": 9,
    "player_game_id": 22953243,
    "page": 1,
    "size": 10
  }'
```

### 2. 测试查询整个圈子战绩

```bash
curl -X POST http://localhost:8000/battle-query/group/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "house_gid": 1,
    "group_id": 9,
    "page": 1,
    "size": 10
  }'
```

### 3. 测试带时间范围

```bash
curl -X POST http://localhost:8000/battle-query/group/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "house_gid": 1,
    "group_id": 9,
    "start_time": 1732752000,
    "end_time": 1732838400,
    "page": 1,
    "size": 10
  }'
```

## ✅ 预期响应

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "list": [
      {
        "id": 12345,
        "house_gid": 1,
        "group_id": 9,
        "room_uid": 1001,
        "kind_id": 1,
        "base_score": 1,
        "battle_at": "2025-11-29T14:30:00+08:00",
        "players_json": "[...]",
        "player_game_id": 22953243,
        "player_game_name": "1106162940",
        "group_name": "b022mc的圈子",
        "score": 1250,
        "fee": 50,
        "factor": 0,
        "player_balance": 0,
        "created_at": "2025-11-29T14:30:05+08:00"
      }
    ],
    "total": 1
  }
}
```

## 📝 数据说明

### 字段含义
- `score`: 输赢分数（分），正数=赢，负数=输
- `fee`: 手续费（分）
- `player_balance`: 玩家余额（分，可能为0如果未实现）
- `battle_at`: 对战时间（ISO 8601格式）
- `player_game_id`: 游戏玩家ID（如：22953243）
- `player_game_name`: 游戏账号名称（如：1106162940）

### 重要提示
1. **所有金额单位都是分**，前端需要除以100转为元
2. 时间戳使用 Unix 秒级时间戳
3. 返回数据按 `battle_at` 降序排列（最新的在前）
4. 默认每页10条，最多100条

## 🔐 权限要求

接口已配置权限验证：
- 需要 JWT 认证
- 需要 `battles:view` 权限

确保你的角色有这个权限：

```sql
-- 查看权限
SELECT * FROM basic_permission WHERE perm_code = 'battles:view';

-- 如果没有，创建权限
INSERT INTO basic_permission (perm_name, perm_code, description)
VALUES ('查看战绩', 'battles:view', '查看战绩记录和统计');

-- 为角色分配权限（假设店铺管理员角色ID是2）
INSERT INTO basic_role_permission_rel (role_id, permission_id)
VALUES (2, (SELECT id FROM basic_permission WHERE perm_code = 'battles:view'));
```

## 🚀 启用前端

实现完成后，前端启用步骤：

1. **删除禁用代码**：
   ```typescript
   // /components/(tabs)/tables/members-list.tsx 第274行
   // 将这行：
   {false && item.game_player_id && item.game_id && myGroupId && item.current_group_id === myGroupId && (
   
   // 改为：
   {item.game_player_id && item.game_id && myGroupId && item.current_group_id === myGroupId && (
   ```

2. **测试功能**：
   - 点击成员的"📊 查看战绩"按钮
   - 确认能正常加载战绩列表
   - 检查时间、金额显示是否正确

## 📊 数据流程

```
用户点击"查看战绩"
    ↓
前端调用 listGroupBattles({
  house_gid: 1,
  group_id: 9,
  player_game_id: 22953243,
  page: 1,
  size: 10
})
    ↓
后端 Service 层验证权限和参数
    ↓
后端 Biz 层调用 Repo 查询
    ↓
Repo 查询 game_battle_record 表
    ↓
返回战绩列表给前端
    ↓
前端渲染战绩卡片
```

## 🔍 调试提示

如果遇到问题：

1. **检查数据表**：
   ```sql
   -- 查看是否有战绩数据
   SELECT * FROM game_battle_record 
   WHERE house_gid = 1 AND group_id = 9 
   ORDER BY battle_at DESC LIMIT 10;
   ```

2. **检查权限**：
   ```sql
   -- 查看用户是否有 battles:view 权限
   SELECT u.username, r.role_name, p.perm_code
   FROM basic_user u
   JOIN basic_user_role_rel urr ON u.id = urr.user_id
   JOIN basic_role r ON urr.role_id = r.id
   JOIN basic_role_permission_rel rpr ON r.id = rpr.role_id
   JOIN basic_permission p ON rpr.permission_id = p.id
   WHERE u.id = <your_user_id> AND p.perm_code = 'battles:view';
   ```

3. **查看日志**：检查后端日志是否有错误信息

4. **API 测试**：使用 Postman 或 curl 直接测试接口

## ⏱️ 预计完成时间

- 添加 Biz 层方法：5 分钟
- 实现 Service 层接口：10 分钟
- 测试和调试：10-15 分钟
- **总计：25-30 分钟**

## 📚 相关文件

- Repo 层：`/internal/dal/repo/game/battle_record.go`（已实现）
- Biz 层：`/internal/biz/game/battle_record.go`（需添加方法）
- Service 层：`/internal/service/game/battle_record.go`（需实现）
- 前端API：`/battle-reusables/services/battles/query.ts`（已实现）
- 前端组件：`/battle-reusables/components/(tabs)/tables/members-list.tsx`（已实现但禁用）
