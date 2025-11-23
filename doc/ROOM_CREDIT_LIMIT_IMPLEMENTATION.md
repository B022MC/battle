# 房间额度限制功能 - 正确实现方式

## 问题背景

在 passing-dragonfly 项目中，有"额度100/1/红中"指令，表示：**低于100分的玩家不能进入1分底的红中房间**。

## 关键理解

**游戏额度房间规则应该从游戏服务器实时获取和检查**，而不是存储在 battle-tiles 的数据库中。

## passing-dragonfly 的实现方式

### 1. 额度规则存储

在 passing-dragonfly 中，额度规则存储在 `t_house` 表的 `game_credits` 字段（JSON格式）：

```go
type THouse struct {
    GameCredits  string `json:"game_credits"`  // 各种游戏额度，JSON结构
    // ...
}

type GameCredit struct {
    Group     string `json:"group"`      // 圈子名称
    Kind      int    `json:"code"`       // 游戏类型
    BaseScore int    `json:"base_score"` // 底分
    Credit    int    `json:"credit"`     // 额度限制（分）
}
```

### 2. 额度检查时机

**关键点：额度检查发生在玩家坐下（SitDown）时**

```go
func (that *House) _handlePlayerSitDown(playerGid, kindId, baseScore, mappedNum int) {
    // 1. 获取玩家ID
    playerId, ok := that._getPlayerIdByGameId(playerGid)
    
    // 2. 获取玩家余额
    bal, ok := that._getPlayerBalance(playerId)
    
    // 3. 获取该房间的额度要求（根据 kindId 和 baseScore）
    credit, _, ok := that._getPlayerGameCredit(playerId, kindId, baseScore)
    
    // 4. 检查余额是否满足额度要求
    if bal < credit {
        // 余额不足，解散桌子（踢出玩家）
        that.session.DismissTable(kindId, mappedNum)
        // 发送通知
        that.pushWxNotice("", playerId, playerGroup, model.NOTICE_TYPE_CREDIT, "🤝🤝🤝🤝🤝")
    }
}
```

### 3. 额度查找优先级

```go
func (that *House) _getGameCredit(group string, kind int, baseScore int) (int, bool, bool) {
    // 1. 全局默认 (0-0)
    if v, ok := that.gameCredits.Load("0-0"); ok {
        return gc.Credit, false, true
    }
    
    // 2. 圈子默认 (group-0-0)
    if v, ok := that.gameCredits.Load(fmt.Sprintf("%s-0-0", group)); ok {
        return gc.Credit, false, true
    }
    
    // 3. 圈子精确 (group-kind-base)
    if v, ok := that.gameCredits.Load(fmt.Sprintf("%s-%d-%d", group, kind, baseScore)); ok {
        return gc.Credit, true, true
    }
    
    // 4. 全局精确 (-kind-base)
    if v, ok := that.gameCredits.Load(fmt.Sprintf("-%d-%d", kind, baseScore)); ok {
        return gc.Credit, true, true
    }
    
    // 5. 兜底额度
    return 99999 * 100, true, true
}
```

## battle-tiles 的实现方案

### 方案一：由游戏服务器完全处理（推荐）

**如果游戏服务器（plaza）已经支持额度检查**，那么：

1. **battle-tiles 不需要实现额度检查逻辑**
2. 游戏服务器会在玩家坐下时自动检查额度
3. 如果余额不足，游戏服务器会自动解散桌子
4. battle-tiles 只需要监听游戏服务器的事件即可

```go
// battle-tiles 不需要额外代码
// 游戏服务器会自动处理额度检查
```

**优点：**
- 逻辑集中在游戏服务器，减少重复代码
- 规则统一，所有客户端行为一致
- battle-tiles 无需维护额度规则

### 方案二：battle-tiles 主动查询（如需要）

如果需要在 battle-tiles 中**提前检查**或**显示**玩家是否满足额度要求：

#### 步骤1：从游戏服务器获取额度配置

```go
// 通过 plaza.Manager 接口查询
// 注意：需要游戏服务器提供相应的API

type RoomCreditChecker struct {
    plazaMgr plaza.Manager
}

func (c *RoomCreditChecker) CheckCredit(ctx context.Context, userID, houseGID, kindID, baseScore int) (canEnter bool, err error) {
    // 1. 获取该用户在该店铺的会话
    session, ok := c.plazaMgr.Get(userID, houseGID)
    if !ok {
        return false, fmt.Errorf("session not found")
    }
    
    // 2. 通过会话向游戏服务器查询额度配置
    // 注意：这需要游戏服务器提供查询接口
    // creditLimit := session.GetCreditLimit(kindID, baseScore)
    
    // 3. 查询玩家余额
    // balance := session.GetBalance()
    
    // 4. 比较
    // return balance >= creditLimit, nil
    
    return true, nil
}
```

#### 步骤2：在需要的地方调用检查

```go
// 例如：在房间列表中显示玩家是否可以进入
func (s *RoomService) ListRooms(c *gin.Context) {
    rooms := getRooms()
    
    for _, room := range rooms {
        canEnter, _ := s.creditChecker.CheckCredit(
            ctx, 
            userID, 
            room.HouseGID, 
            room.KindID, 
            room.BaseScore,
        )
        room.CanEnter = canEnter
    }
    
    response.Success(c, rooms)
}
```

**缺点：**
- 需要游戏服务器提供额度查询API
- 增加网络调用
- 可能存在时差（查询时可以进入，实际坐下时余额已不足）

### 方案三：在 battle-tiles 中缓存规则（不推荐）

如果游戏服务器不支持实时查询，可以：

1. 通过某种方式将游戏服务器的额度配置同步到 battle-tiles
2. 在 battle-tiles 中维护一份缓存
3. 用于前端显示和提前检查

**不推荐原因：**
- 需要保持两边数据一致性
- 增加维护成本
- 可能出现数据不同步问题

## 实际操作建议

### 确认游戏服务器能力

首先需要确认游戏服务器（plaza）是否已经实现了：

1. ✅ **额度检查**：玩家坐下时自动检查余额
2. ✅ **自动踢人**：余额不足时自动解散桌子
3. ❓ **额度查询API**：提供查询接口给 battle-tiles

### 推荐做法

**如果游戏服务器已经实现了额度检查和自动踢人：**

✅ **battle-tiles 不需要做任何事情**

游戏服务器会自动处理，battle-tiles 只需要：
- 保持会话连接
- 监听游戏服务器的事件通知
- 显示通知给管理员

**如果需要在前端提前显示玩家是否可以进入：**

1. 向游戏服务器团队申请提供额度查询API
2. 在 battle-tiles 中调用该API获取规则
3. 在前端显示提示信息

**如果游戏服务器不支持额度检查：**

那就需要在 battle-tiles 中实现完整的额度管理功能（类似我之前创建的方案），但这不是推荐的架构。

## 总结

**核心原则：额度规则应该由游戏服务器管理和执行**

- ✅ 游戏服务器负责存储规则
- ✅ 游戏服务器负责实时检查
- ✅ 游戏服务器负责执行踢人
- ❌ battle-tiles 不应该存储规则
- ❌ battle-tiles 不应该执行检查
- ✅ battle-tiles 可以查询规则（如果游戏服务器提供API）
- ✅ battle-tiles 可以显示提示信息

这样的架构更加清晰，责任分明，减少了数据同步问题。
