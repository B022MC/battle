# 店铺费用规则验证功能

## 功能概述

当玩家加入桌子（坐下）时，后台会自动检测该玩家是否符合店铺费用规则的资格要求。如果玩家余额不足以支付费用，系统会自动将玩家从桌子上踢出（解散桌子）。

此功能与 `passing-dragonfly` 项目的逻辑保持一致。

## 核心逻辑

### 1. 费用规则匹配

费用规则的匹配优先级（与 `passing-dragonfly` 一致）：

1. **全局规则优先**：如果存在全局规则（`kind=0` 且 `base=0`），优先使用全局规则
2. **特定规则**：如果没有全局规则，匹配特定游戏类型和底分的规则
   - `kind`：游戏类型ID（0 表示所有游戏类型）
   - `base`：底分（0 表示所有底分）

### 2. 玩家资格检查流程

玩家坐下时会进行以下检查：

1. **玩家身份验证**
   - 检查玩家是否在店铺成员列表中
   - 如果玩家未录入，自动解散桌子

2. **房间额度检查**
   - 获取房间的额度限制配置
   - 计算有效额度 = 房间额度 + 玩家个人额度调整
   - 检查玩家余额是否满足额度要求

3. **店铺费用规则检查**（新增功能）
   - 从店铺设置中获取费用规则配置
   - 根据当前房间的游戏类型和底分，计算需要的最低费用
   - 检查玩家余额是否足以支付费用
   - 如果余额不足，自动解散桌子

## 实现细节

### 核心文件

1. **`internal/biz/game/room_credit_limit.go`**
   - `HandlePlayerSitDown`: 处理玩家坐下事件的主逻辑
   - `calculateMinRequiredFee`: 计算玩家需要的最低费用

2. **`internal/biz/game/room_credit_event_handler.go`**
   - `RoomCreditEventHandler`: 房间额度事件处理器
   - `sessionCreditHandler`: 实现 `plaza.IPlazaHandler` 接口
   - `OnUserSitDown`: 当玩家坐下时触发检查

3. **`internal/biz/game/ctrl_session.go`**
   - `CtrlSessionUseCase`: 集成费用检查处理器
   - `compositeHandler`: 组合多个事件处理器

### 费用计算逻辑

```go
func calculateMinRequiredFee(config *FeesConfig, kindID int, baseScore int) int32 {
    // 1. 优先查找全局规则 (kind=0 && base=0)
    for _, rule := range config.Rules {
        if rule.Kind == 0 && rule.Base == 0 {
            return int32(rule.Fee)
        }
    }

    // 2. 查找特定游戏类型和底分的规则
    for _, rule := range config.Rules {
        if rule.Kind == 0 && rule.Base == 0 {
            continue // 跳过全局规则
        }
        
        kindMatches := (rule.Kind == 0 || rule.Kind == kindID)
        baseMatches := (rule.Base == 0 || rule.Base == baseScore)
        
        if kindMatches && baseMatches {
            return int32(rule.Fee)
        }
    }

    return 0 // 未匹配到任何规则
}
```

## 配置示例

### 店铺费用规则配置

费用规则存储在 `game_house_settings` 表的 `fees_json` 字段中，格式为 JSON：

```json
{
  "rules": [
    {
      "threshold": 100,
      "fee": 1000,
      "kind": 0,
      "base": 0
    }
  ]
}
```

字段说明：
- `threshold`: 分数阈值（达到该分数时收取费用）
- `fee`: 费用（单位：分，1000 = 10元）
- `kind`: 游戏类型ID（0 表示所有游戏）
- `base`: 底分（0 表示不限）

### 前端配置界面

前端使用 `battle-reusables/components/(shop)/shop-fees/shop-fees-view.tsx` 组件进行配置：

1. 选择店铺
2. 添加费用规则
   - 游戏类型（可选）
   - 底分（可选）
   - 分数阈值（必填）
   - 费用金额（必填）
3. 保存配置

## 技术亮点

1. **与 passing-dragonfly 逻辑一致**
   - 费用规则匹配逻辑完全对齐
   - 全局规则优先，特定规则其次
   
2. **实时验证**
   - 玩家坐下时立即检查
   - 不符合条件的玩家无法进入游戏

3. **优雅的事件处理**
   - 使用组合模式集成多个事件处理器
   - 支持灵活扩展

4. **容错设计**
   - 配置错误不影响游戏运行
   - 详细的日志记录便于排查问题

## 依赖注入

Wire 自动注入依赖：

```go
game2.NewRoomCreditLimitUseCase(
    roomCreditLimitRepo,
    memberRepo,
    houseSettingsRepo,
    plazaMgr,
    logger,
)

game2.NewRoomCreditEventHandler(
    roomCreditLimitUseCase,
    plazaMgr,
    logger,
)

game2.NewCtrlSessionUseCase(
    ctrlRepo,
    linkRepo,
    sessRepo,
    plazaMgr,
    syncMgr,
    logger,
    creditHandler, // 自动注入
)
```

## 日志输出

系统会记录详细的检查日志：

- **玩家坐下**: `HandlePlayerSitDown: gameID=xxx, balance=xxx, ...`
- **额度检查**: `checking fee eligibility, gameID=xxx, balance=xxx, requiredFee=xxx`
- **余额不足**: `insufficient balance for fee, gameID=xxx, balance=xxx < fee=xxx, dismissing table`
- **检查通过**: `all checks passed, gameID=xxx, balance=xxx >= credit=xxx`

## 测试建议

1. **正常场景**
   - 玩家余额充足，能正常进入桌子
   
2. **余额不足场景**
   - 玩家余额低于费用要求，自动被踢出
   
3. **无规则场景**
   - 店铺未配置费用规则，不影响玩家进入
   
4. **规则匹配场景**
   - 测试全局规则和特定规则的优先级
   - 测试游戏类型和底分的匹配逻辑

## 注意事项

1. **费用单位**：配置时注意费用单位是"分"（100分 = 1元）
2. **规则优先级**：全局规则始终优先于特定规则
3. **余额检查**：同时检查房间额度和费用规则
4. **自动解散**：不符合条件时会立即解散桌子，玩家会收到通知
