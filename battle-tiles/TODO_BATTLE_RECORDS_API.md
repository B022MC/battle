# 战绩查询接口实现 TODO

## 📍 位置
`/Users/b022mc/project/battle/battle-tiles/internal/service/game/battle_record.go`

## ⚠️ 待实现接口

### ListGroupBattles - 查询圈子成员战绩

**文件位置**：第 201-203 行

**当前状态**：
```go
func (s *BattleRecordService) ListGroupBattles(context *gin.Context) {
	//TODO
}
```

**需要实现的功能**：
店铺管理员查询自己圈子成员的战绩列表

## 📋 接口规格

### 请求参数

**类型定义**（需要添加）：
```go
type ListGroupBattlesRequest struct {
	HouseGID     int32  `json:"house_gid" binding:"required"`     // 店铺ID
	GroupID      int32  `json:"group_id" binding:"required"`      // 圈子ID
	PlayerGameID *int32 `json:"player_game_id"`                   // 可选：玩家游戏ID，不传则查询所有成员
	StartTime    *int64 `json:"start_time"`                       // 可选：开始时间（Unix timestamp）
	EndTime      *int64 `json:"end_time"`                         // 可选：结束时间（Unix timestamp）
	Page         int32  `json:"page"`                             // 页码
	Size         int32  `json:"size"`                             // 每页数量
}
```

### 响应数据

```go
type ListGroupBattlesResponse struct {
	List  []BattleRecordDTO `json:"list"`
	Total int64             `json:"total"`
}

type BattleRecordDTO struct {
	ID             int64   `json:"id"`
	HouseGID       int32   `json:"house_gid"`
	GroupID        int32   `json:"group_id"`
	GroupName      string  `json:"group_name"`
	RoomUID        int64   `json:"room_uid"`
	KindID         int32   `json:"kind_id"`
	BaseScore      int32   `json:"base_score"`
	BattleAt       string  `json:"battle_at"`        // ISO 8601 格式
	PlayersJSON    string  `json:"players_json"`     // 玩家列表JSON
	PlayerID       *int32  `json:"player_id"`        // 当前查询的玩家ID
	PlayerGameID   *int32  `json:"player_game_id"`   // 当前查询的玩家游戏ID
	Score          int64   `json:"score"`            // 输赢分数（分）
	Fee            int64   `json:"fee"`              // 手续费（分）
	Factor         float64 `json:"factor"`           // 倍率
	PlayerBalance  int64   `json:"player_balance"`   // 对战后余额（分）
	CreatedAt      string  `json:"created_at"`       // ISO 8601 格式
}
```

## 🔧 实现要点

### 1. 权限验证
- ✅ 已添加权限中间件 `middleware.RequirePerm("battles:view")`
- 需要验证管理员是否属于该圈子
- 需要验证圈子是否属于该店铺

### 2. 查询逻辑
```go
func (s *BattleRecordService) ListGroupBattles(c *gin.Context) {
	var in ListGroupBattlesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 从 JWT 中获取用户 ID
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// TODO: 验证管理员权限和圈子归属
	// 1. 验证用户是该店铺的管理员
	// 2. 验证圈子属于该店铺
	// 3. 验证管理员有权限查看该圈子的战绩

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

	// 调用业务层查询战绩
	list, total, err := s.uc.ListGroupBattles(c, &gameBiz.ListGroupBattlesParams{
		HouseGID:     in.HouseGID,
		GroupID:      in.GroupID,
		PlayerGameID: in.PlayerGameID,
		StartTime:    start,
		EndTime:      end,
		Page:         in.Page,
		Size:         in.Size,
	})
	if err != nil {
		response.Fail(c, ecode.ServerError, err)
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
```

### 3. 业务层实现

需要在 `internal/biz/game/battle_record.go` 中添加：

```go
type ListGroupBattlesParams struct {
	HouseGID     int32
	GroupID      int32
	PlayerGameID *int32
	StartTime    *time.Time
	EndTime      *time.Time
	Page         int32
	Size         int32
}

func (uc *BattleRecordUseCase) ListGroupBattles(
	ctx context.Context,
	params *ListGroupBattlesParams,
) ([]*BattleRecordDTO, int64, error) {
	// TODO: 实现业务逻辑
	// 1. 查询圈子成员的战绩记录
	// 2. 如果指定了 PlayerGameID，只查询该玩家的战绩
	// 3. 按时间倒序排列
	// 4. 分页返回
	return nil, 0, nil
}
```

### 4. 数据库查询

需要在 `internal/dal/repo/game/battle_record.go` 中添加：

```go
func (r *BattleRecordRepo) ListByGroup(
	ctx context.Context,
	houseGID int32,
	groupID int32,
	playerGameID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	// TODO: 实现数据库查询
	// 使用 GORM 查询 game_battle_records 表
	// 条件：house_gid = ? AND group_id = ?
	// 如果 playerGameID 不为空，添加条件：player_game_id = ?
	// 如果有时间范围，添加条件：battle_at BETWEEN ? AND ?
	// 排序：ORDER BY battle_at DESC
	// 分页：LIMIT ? OFFSET ?
	return nil, 0, nil
}
```

## 📊 数据表结构

假设使用的表是 `game_battle_records`，需要包含以下字段：

```sql
CREATE TABLE game_battle_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL,
    group_id INT NOT NULL,
    room_uid BIGINT NOT NULL,
    kind_id INT NOT NULL,
    base_score INT NOT NULL,
    battle_at TIMESTAMP NOT NULL,
    players_json TEXT,
    player_id INT,
    player_game_id INT,
    score BIGINT,
    fee BIGINT,
    factor DECIMAL(10,2),
    player_balance BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_house_group (house_gid, group_id),
    INDEX idx_player_game_id (player_game_id),
    INDEX idx_battle_at (battle_at)
);
```

## ✅ 前端准备情况

前端代码已经准备好：
- ✅ API 调用函数：`listGroupBattles()` - `/services/battles/query.ts`
- ✅ UI 组件：战绩列表展示 - `/components/(tabs)/tables/members-list.tsx`
- ✅ 状态管理：展开/收起、加载状态
- ✅ 数据格式化：时间、金额格式化
- ⏸️ **已暂时禁用**：使用 `false &&` 条件禁用显示

## 🚀 启用步骤

后端实现完成后：

1. 确保接口返回正确的数据格式
2. 测试接口是否正常工作
3. 前端删除 `members-list.tsx` 第 274 行的 `false &&` 条件
4. 测试前端功能是否正常

## 🔍 测试用例

### 测试请求
```bash
curl -X POST http://localhost:8000/battle-query/group/battles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "house_gid": 1,
    "group_id": 9,
    "player_game_id": 22953243,
    "page": 1,
    "size": 10
  }'
```

### 预期响应
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
        "group_name": "b022mc的圈子",
        "room_uid": 1001,
        "kind_id": 1,
        "base_score": 1,
        "battle_at": "2025-11-29T14:30:00Z",
        "player_game_id": 22953243,
        "score": 1250,
        "fee": 50,
        "factor": 1.0,
        "player_balance": 15000,
        "created_at": "2025-11-29T14:30:05Z"
      }
    ],
    "total": 1
  }
}
```

## 📝 注意事项

1. **金额单位**：所有金额字段（score、fee、player_balance）都以**分**为单位
2. **时间格式**：统一使用 ISO 8601 格式（`2006-01-02T15:04:05Z07:00`）
3. **权限验证**：必须验证管理员是否有权查看该圈子的战绩
4. **性能优化**：添加适当的索引，确保查询性能
5. **数据一致性**：确保战绩数据与实际游戏记录同步

## 📚 相关文件

- 前端API：`/battle-reusables/services/battles/query.ts`
- 前端类型：`/battle-reusables/services/battles/query-typing.d.ts`
- 前端组件：`/battle-reusables/components/(tabs)/tables/members-list.tsx`
- 后端服务：`/battle-tiles/internal/service/game/battle_record.go`
- 功能文档：`/battle-reusables/MEMBER_BATTLE_RECORDS_FEATURE.md`
