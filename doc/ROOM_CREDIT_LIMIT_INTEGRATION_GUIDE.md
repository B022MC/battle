# 房间额度限制功能集成指南

## 功能说明

实现 "额度100/1/红中" 指令功能：**低于100分的玩家不能进入1分底的红中房间**

## 实现方式

按照 passing-dragonfly 的方式：
1. 额度规则存储在数据库中
2. 监听游戏服务器的 SitDown 事件
3. 在 battle-tiles 中检查余额是否满足额度
4. 如果不足，调用 plaza 接口解散桌子

## 已创建的文件

### 1. 数据库迁移
- `doc/migrations/add_game_room_credit_limit.sql`

### 2. 数据模型
- `internal/dal/model/game/game_room_credit_limit.go`

### 3. Repository层
- `internal/dal/repo/game/room_credit_limit.go`

### 4. Business层
- `internal/biz/game/room_credit_limit.go`
- `internal/biz/game/room_credit_event_handler.go` (需要手动修复)

### 5. Service层
- `internal/service/game/room_credit_limit.go`

### 6. DTO
- `internal/dal/req/room_credit_limit.go`
- `internal/dal/resp/room_credit_limit.go`

## 需要完成的步骤

### 步骤1：修复事件处理器

编辑 `internal/biz/game/room_credit_event_handler.go`，确保导入正确：

```go
package game

import (
	"battle-tiles/internal/infra/plaza"
	utilsplaza "battle-tiles/internal/utils/plaza"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)
```

### 步骤2：添加到依赖注入

#### 2.1 更新 `internal/dal/repo/repo.go`

```go
game.NewRoomCreditLimitRepo,
```

#### 2.2 更新 `internal/biz/biz.go`

```go
game.NewRoomCreditLimitUseCase,
game.NewRoomCreditEventHandler,
```

#### 2.3 更新 `internal/service/service.go`

```go
game.NewRoomCreditLimitService,
```

#### 2.4 更新 `internal/router/game_router.go`

在 `GameRouter` 结构体中添加：
```go
roomCreditLimitService *game.RoomCreditLimitService
```

在 `InitRouter` 方法中添加：
```go
r.roomCreditLimitService.RegisterRouter(root)
```

在 `NewGameRouter` 函数参数中添加：
```go
roomCreditLimitService *game.RoomCreditLimitService,
```

在返回值中添加：
```go
roomCreditLimitService: roomCreditLimitService,
```

### 步骤3：在启动会话时注册事件处理器

当调用 `plaza.Manager.StartUser` 启动会话时，需要传入事件处理器。

找到启动会话的地方（可能在 `internal/biz/game/ctrl_session.go` 或类似文件），修改为：

```go
// 注入 RoomCreditEventHandler
eventHandler := roomCreditEventHandler.CreateHandler(userID, houseGID)

// 启动会话时传入 handler
err := plazaMgr.StartUser(ctx, userID, houseGID, mode, identifier, pwdMD5, gameUserID, eventHandler)
```

### 步骤4：执行数据库迁移

```bash
psql -U username -d database_name -f doc/migrations/add_game_room_credit_limit.sql
```

### 步骤5：重新生成 wire

```bash
cd cmd/go-kgin-platform
wire
```

## API接口

### 1. 设置房间额度限制

```
POST /room-credit/set
```

请求示例：
```json
{
  "house_gid": 123456,
  "group_name": "",     // 空表示全局
  "game_kind": 5,        // 游戏类型ID（5=红中）
  "base_score": 1,       // 底分
  "credit_limit": 10000  // 100元 = 10000分
}
```

### 2. 查询额度限制

```
POST /room-credit/list
```

请求示例：
```json
{
  "house_gid": 123456,
  "group_name": ""  // 空表示查询所有
}
```

### 3. 删除额度限制

```
POST /room-credit/delete
```

### 4. 检查玩家是否满足额度

```
POST /room-credit/check
```

## 额度查找优先级

系统按以下优先级查找额度限制：

1. **圈子+游戏类型+底分** (`group-kind-base`)
2. **圈子默认** (`group-0-0`)
3. **全局+游戏类型+底分** (`''-kind-base`)
4. **全局默认** (`''-0-0`)
5. **兜底值**：99999元

## 工作流程

1. **管理员设置额度**
   - 调用 `/room-credit/set` API 设置规则
   - 规则保存到 `game_room_credit_limit` 表

2. **玩家坐下房间**
   - 游戏服务器发送 `OnUserSitDown` 事件
   - `sessionCreditHandler` 接收事件
   - 从 Session 获取桌子信息（kindID、baseScore）
   - 调用 `RoomCreditLimitUseCase.HandlePlayerSitDown`
   - 检查玩家余额是否满足额度要求

3. **余额不足处理**
   - 调用 `plaza.Manager.DismissTable` 解散桌子
   - 玩家被踢出房间

## 玩家个人额度调整

`game_member.credit` 字段可用于调整单个玩家的额度：
- **正值**：提高额度要求（更严格）
- **负值**：降低额度要求（更宽松）
- **0**：使用房间基础额度

有效额度 = 房间额度 + 玩家个人额度

## 注意事项

1. 确保 `plaza.Manager` 已正确注入到 `RoomCreditEventHandler`
2. 每个会话需要创建独立的 `sessionCreditHandler` 实例
3. 事件处理器需要在启动会话时注册
4. 桌子信息从 Session 的缓存中获取，确保 `ListTables()` 有数据

## 权限配置

需要在权限系统中添加以下权限：
- `room:credit:set` - 设置额度限制
- `room:credit:view` - 查看额度限制
- `room:credit:delete` - 删除额度限制
- `room:credit:check` - 检查玩家额度
