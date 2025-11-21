# 战绩同步费用规则说明

## 概述

战绩同步时，系统会根据店铺配置的费用规则自动计算和分配运费。费用规则基于战局分数、玩法类型等条件进行匹配。

---

## 费用配置结构

### 数据库表
- **表名**: `game_house_settings`
- **关键字段**:
  - `house_gid` (int32): 店铺GID
  - `fees_json` (text): 费用规则配置（JSON格式）
  - `share_fee` (boolean): 是否启用分运费模式

### 费用规则JSON格式

```json
{
  "rules": [
    {
      "threshold": 5000,
      "fee": 500,
      "kind": "",
      "base": 0
    },
    {
      "threshold": 3000,
      "fee": 300
    },
    {
      "threshold": 0,
      "fee": 100
    }
  ]
}
```

#### 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `threshold` | int | ✅ | 分数阈值（最高分需达到此值） |
| `fee` | int | ✅ | 运费金额（单位：分） |
| `kind` | string | ❌ | 玩法类型（预留，暂未实现） |
| `base` | int | ❌ | 底分（预留，暂未实现） |

---

## 费用计算流程

### 1. 计算总运费

**函数**: `CalculateFee(feesJSON string, battle *BattleInfo) int32`

**逻辑**:
1. 解析费用规则配置
2. 找出本局最高分
3. **从上到下**匹配规则（规则应按阈值从高到低排列）
4. 返回第一个匹配的费用金额

**示例**:
```go
// 配置规则
rules = [
  {threshold: 5000, fee: 500},
  {threshold: 3000, fee: 300},
  {threshold: 0, fee: 100}
]

// 战局最高分 = 4200
// 匹配结果: 4200 >= 3000 → 返回 300 分
```

---

### 2. 构建玩家圈子映射

**函数**: `buildPlayerGroupMapping(houseGID, players)`

**流程**:
```
game_player_id → game_member.game_id → game_member.group_id
```

**条件**:
- 玩家必须在 `game_member` 表中
- 必须有有效的 `group_id`（不为 NULL 且不为 0）

**结果**:
- `playerGroups`: 玩家ID → 圈子ID 的映射
- `validPlayers`: 有效玩家标记（只有在圈子中的玩家才计费）

---

### 3. 费用分配

**函数**: `CalculateFeeDistribution(players, playerGroups, totalFee, shareFee)`

#### 模式1: 不分运费 (share_fee = false)

**规则**: 赢家承担全部运费

```
1. 找出最高分玩家（赢家）
2. 赢家平分总运费
3. 赢家圈子内的玩家再平分圈子费用
```

**示例**:
```
战局: A圈子(玩家1, 玩家2), B圈子(玩家3)
分数: 玩家1=5000, 玩家2=3000, 玩家3=5000
总运费: 500分

赢家: 玩家1和玩家3（都是5000分）
赢家圈子: A圈子, B圈子

费用分配:
- A圈子总费用: 500 / 2 = 250分
  - 玩家1: 250 / 2 = 125分
  - 玩家2: 250 / 2 = 125分
- B圈子总费用: 500 / 2 = 250分
  - 玩家3: 250 / 1 = 250分
```

#### 模式2: 分运费 (share_fee = true)

**规则**: 所有圈子平分运费

```
1. 统计参与的圈子数量
2. 总运费 ÷ 圈子数 = 每圈子应付
3. 圈子内的玩家再平分圈子费用
```

**示例**:
```
战局: A圈子(玩家1, 玩家2), B圈子(玩家3), C圈子(玩家4)
总运费: 600分

费用分配:
- A圈子总费用: 600 / 3 = 200分
  - 玩家1: 200 / 2 = 100分
  - 玩家2: 200 / 2 = 100分
- B圈子总费用: 600 / 3 = 200分
  - 玩家3: 200 / 1 = 200分
- C圈子总费用: 600 / 3 = 200分
  - 玩家4: 200 / 1 = 200分
```

---

## 数据记录

### 1. 战绩记录 (game_battle_record)

每个**有效玩家**都会生成一条战绩记录：

| 字段 | 说明 |
|------|------|
| `house_gid` | 店铺GID |
| `group_id` | 玩家所属圈子ID |
| `room_uid` | 房间ID |
| `kind_id` | 玩法ID |
| `base_score` | 底分 |
| `battle_at` | 战局时间 |
| `player_game_id` | 游戏玩家ID |
| `player_game_name` | 游戏玩家名称 |
| `players_json` | 完整玩家数据（JSON） |
| `score` | 玩家得分 |
| `fee` | **玩家应付运费** |

### 2. 费用结算记录 (game_fee_settle)

记录费用结算明细（主要用于分运费模式）：

| 字段 | 说明 |
|------|------|
| `house_gid` | 店铺GID |
| `room_uid` | 房间ID |
| `group_id` | 圈子ID |
| `amount` | 结算金额（正数=支出，负数=收入） |
| `settle_type` | 结算类型（fee=运费，payoff=补偿） |
| `players_json` | 玩家列表（JSON） |
| `settle_at` | 结算时间 |

---

## 关键逻辑点

### ✅ 只对有效玩家计费

**有效玩家**的定义：
1. 在 `game_member` 表中有记录
2. 有明确的 `group_id`（不为 NULL 且不为 0）

**跳过的玩家**：
- 不在 `game_member` 表中
- `group_id` 为 NULL 或 0
- 这些玩家的战绩**不会被记录**

### ✅ 零除保护

所有除法运算都有零除保护：
```go
// 示例
if numGroups > 0 {
    sharedFee = totalFee / int32(numGroups)
}
```

### ✅ 赢家判定

- 找出**所有**最高分玩家
- 支持**多个赢家**（平分情况）
- 一次遍历完成，性能优化

---

## 配置示例

### 示例1: 简单运费规则

```sql
UPDATE game_house_settings
SET 
  fees_json = '{
    "rules": [
      {"threshold": 5000, "fee": 500},
      {"threshold": 3000, "fee": 300},
      {"threshold": 0, "fee": 100}
    ]
  }',
  share_fee = false
WHERE house_gid = 60870;
```

**说明**:
- 最高分 >= 5000: 收费 500分
- 最高分 >= 3000: 收费 300分
- 最高分 >= 0: 收费 100分
- 不分运费（赢家承担）

### 示例2: 分运费模式

```sql
UPDATE game_house_settings
SET 
  fees_json = '{
    "rules": [
      {"threshold": 0, "fee": 400}
    ]
  }',
  share_fee = true
WHERE house_gid = 60870;
```

**说明**:
- 固定收费 400分
- 分运费（所有圈子平分）

### 示例3: 不收费

```sql
UPDATE game_house_settings
SET 
  fees_json = '{"rules": []}',
  share_fee = false
WHERE house_gid = 60870;
```

或者设置为 NULL：
```sql
UPDATE game_house_settings
SET 
  fees_json = NULL,
  share_fee = false
WHERE house_gid = 60870;
```

---

## 常见问题

### Q1: 为什么有些玩家不计费？

**A**: 只有满足以下条件的玩家才会计费：
1. 在 `game_member` 表中有记录
2. 有有效的 `group_id`（不为 NULL 且不为 0）

**检查方法**:
```sql
-- 查看玩家是否在圈子中
SELECT game_id, game_name, group_id 
FROM game_member 
WHERE house_gid = 60870 AND game_id = 21309263;
```

### Q2: 如何查看某局的费用分配？

```sql
-- 查看战绩记录中的费用
SELECT 
    player_game_id,
    player_game_name,
    group_id,
    score,
    fee
FROM game_battle_record
WHERE house_gid = 60870 AND room_uid = 524605;

-- 查看费用结算记录
SELECT * FROM game_fee_settle
WHERE house_gid = 60870 AND room_uid = 524605;
```

### Q3: 分运费和不分运费的区别？

| 模式 | 谁承担运费 | 适用场景 |
|------|-----------|----------|
| **不分运费** | 赢家圈子承担全部 | 赢家支付场地费 |
| **分运费** | 所有圈子平分 | 大家共同承担 |

### Q4: 费用规则的匹配顺序？

**从上到下匹配**，返回第一个满足条件的规则：
```json
{
  "rules": [
    {"threshold": 5000, "fee": 500},  // 优先匹配
    {"threshold": 3000, "fee": 300},  // 其次
    {"threshold": 0, "fee": 100}      // 最后
  ]
}
```

**建议**: 规则按阈值从高到低排列

---

## 代码位置

| 文件 | 说明 |
|------|------|
| `internal/biz/game/battle_record.go` | 战绩同步主流程 |
| `internal/biz/game/fee_calculator.go` | 费用计算核心逻辑 |
| `internal/dal/model/game/house_settings.go` | 配置数据模型 |
| `internal/dal/repo/game/house_settings.go` | 配置数据访问 |

---

**文档版本**: 2025-11-22  
**作者**: AI Assistant  
**状态**: ✅ 完整
