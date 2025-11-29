# 首页统计功能 - 与微信机器人命令对应

## ✅ 功能完成

根据微信机器人命令实现首页统计功能：

### 店铺管理员
- 📊 **今日统计** - 整个店铺的成绩、笔数、金额
- 📊 **昨日统计** - 整个店铺的成绩、笔数、金额
- 📊 **上周统计** - 整个店铺的成绩、笔数、金额

### 超级管理员
- 📊 **今日统计** - 所有店铺汇总的成绩、笔数、金额
- 📊 **昨日统计** - 所有店铺汇总的成绩、笔数、金额
- 📊 **上周统计** - 所有店铺汇总的成绩、笔数、金额

---

## 🤖 对应微信机器人命令

```go
// 店铺管理员命令
CmdGetTodayStat      = "今日统计"   // 参数：无。整个店铺的成绩，笔数，金额
CmdGetYesterdayStat  = "昨日统计"   // 参数：无。整个店铺的成绩，笔数，金额
CmdGetLastweekStat   = "上周统计"   // 参数：无。整个店铺的成绩，笔数，金额
```

**功能对应**：
- ✅ **笔数** = `total_games`（总局数）
- ✅ **成绩** = `total_score`（总得分，分为单位）
- ✅ **金额** = `total_fee`（总手续费，分为单位）

---

## 🎯 功能特性

### 1. **时间维度**

#### 今日统计
- 时间范围：今天 0:00 - 当前时间
- 颜色主题：绿色
- 显示数据：笔数、得分、手续费

#### 昨日统计
- 时间范围：昨天 0:00 - 今天 0:00
- 颜色主题：蓝色
- 显示数据：笔数、得分、手续费

#### 上周统计
- 时间范围：上周一 0:00 - 上周日 24:00
- 颜色主题：紫色
- 显示数据：笔数、得分、手续费

### 2. **角色区分**

#### 店铺管理员
- 看到：**本店铺**的统计数据
- 数据范围：当前店铺的所有圈子汇总
- 卡片标题：`店铺统计`
- 副标题：`店铺 ID: {house_gid}`

#### 超级管理员
- 看到：**全平台**的统计数据
- 数据范围：所有店铺的汇总
- 卡片标题：`全平台统计`
- 副标题：`所有店铺数据汇总`

---

## 📁 文件清单

### 修改文件

#### 1. `/components/(tabs)/home/house-stats-card.tsx`
**更新统计卡片组件**：
- 添加 `title` 和 `subtitle` 属性支持自定义标题
- 将 `thisWeekStats` 改为 `todayStats`（今日统计）
- 支持三个时间维度：今日、昨日、上周

**新增属性**：
```typescript
interface HouseStatsCardProps {
  title?: string;                       // 卡片标题（默认"店铺统计"）
  subtitle?: string;                    // 卡片副标题
  todayStats?: HousePeriodStats[];      // 今日统计
  yesterdayStats?: HousePeriodStats[];  // 昨日统计
  lastWeekStats?: HousePeriodStats[];   // 上周统计
  loading?: boolean;
  onViewDetails?: () => void;
}
```

#### 2. `/app/(tabs)/index.tsx`
**首页主文件更新**：

**时间范围函数更新**：
```typescript
// 添加 'today' 类型
function getTimeRange(type: 'today' | 'yesterday' | 'thisWeek' | 'lastWeek')
```

**店铺管理员数据加载**：
```typescript
// 加载店铺统计（今日、昨日、上周）
const loadShopStats = async (houseGid: number) => {
  const [today, yesterday, lastWeek] = await Promise.all([
    getHouseStats({ house_gid: houseGid, ...getTimeRange('today') }),
    getHouseStats({ house_gid: houseGid, ...getTimeRange('yesterday') }),
    getHouseStats({ house_gid: houseGid, ...getTimeRange('lastWeek') })
  ]);
};
```

**超级管理员数据加载**：
```typescript
// 为每个店铺并发请求三个时间段的统计
const allRequests = houses.flatMap(houseGid => [
  getHouseStats({ house_gid: houseGid, ...getTimeRange('today') }),
  getHouseStats({ house_gid: houseGid, ...getTimeRange('yesterday') }),
  getHouseStats({ house_gid: houseGid, ...getTimeRange('lastWeek') })
]);
```

---

## 🔧 技术实现

### 数据流程

#### 店铺管理员

```
用户进入首页（店铺管理员）
    ↓
获取管理员信息 (shopsAdminsMe)
    ↓
获取 house_gid
    ↓
并发请求店铺的三个时间段统计
    ├─→ 今日统计 (getHouseStats)
    ├─→ 昨日统计 (getHouseStats)
    └─→ 上周统计 (getHouseStats)
    ↓
显示店铺统计卡片
```

#### 超级管理员

```
用户进入首页（超级管理员）
    ↓
获取所有店铺列表 (shopsHousesOptions)
    ↓
为每个店铺并发请求三个时间段
    ├─→ 店铺1：今日、昨日、上周
    ├─→ 店铺2：今日、昨日、上周
    └─→ 店铺N：今日、昨日、上周
    ↓
分类汇总所有店铺数据
    ├─→ todayStats: HouseStats[]
    ├─→ yesterdayStats: HouseStats[]
    └─→ lastWeekStats: HouseStats[]
    ↓
HouseStatsCard 自动计算总和
    ↓
显示全平台统计卡片
```

### 时间范围计算

#### 今日
```typescript
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
return {
  start_time: Math.floor(today.getTime() / 1000),  // 今天0点
  end_time: Math.floor(now.getTime() / 1000)       // 现在
};
```

#### 昨日
```typescript
const yesterdayStart = new Date(today);
yesterdayStart.setDate(yesterdayStart.getDate() - 1);
return {
  start_time: Math.floor(yesterdayStart.getTime() / 1000),  // 昨天0点
  end_time: Math.floor(today.getTime() / 1000)              // 今天0点
};
```

#### 上周
```typescript
const dayOfWeek = today.getDay() || 7;
const lastWeekStart = new Date(today);
lastWeekStart.setDate(lastWeekStart.getDate() - dayOfWeek - 6);  // 上周一
const lastWeekEnd = new Date(lastWeekStart);
lastWeekEnd.setDate(lastWeekEnd.getDate() + 7);  // 上周日+1天
return {
  start_time: Math.floor(lastWeekStart.getTime() / 1000),
  end_time: Math.floor(lastWeekEnd.getTime() / 1000)
};
```

---

## 🎨 UI 展示

### 店铺管理员首页

```
┌─────────────────────────────────┐
│ 🏆 欢迎回来！                   │
│ 用户: shop_admin                │
│ 认证状态: 已登录                │
│ 店铺: 60870                     │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 📢 系统公告                     │
│ ...                             │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 🏢 店铺统计        查看详情 >   │
│ 店铺 ID: 60870                  │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 今日统计           125 场│ │  ← 绿色
│ │ 总得分         总手续费     │ │
│ │ +1,250.50      125.50       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 昨日统计           112 场│ │  ← 蓝色
│ │ 总得分         总手续费     │ │
│ │ +1,120.00      112.00       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 上周统计           520 场│ │  ← 紫色
│ │ 总得分         总手续费     │ │
│ │ -1,255.50      520.00       │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

### 超级管理员首页

```
┌─────────────────────────────────┐
│ 🏢 全平台统计      查看详情 >   │
│ 所有店铺数据汇总                │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 今日统计           350 场│ │  ← 绿色
│ │ 总得分         总手续费     │ │
│ │ +3,500.00      350.00       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 昨日统计           320 场│ │  ← 蓝色
│ │ 总得分         总手续费     │ │
│ │ +3,200.00      320.00       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 上周统计         1,500 场│ │  ← 紫色
│ │ 总得分         总手续费     │ │
│ │ +15,000.00    1,500.00      │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

---

## 📊 数据说明

### 字段对应

| 微信机器人术语 | API 字段 | 说明 | 单位 |
|--------------|----------|------|------|
| 笔数 | `total_games` | 总对局数 | 场 |
| 成绩 | `total_score` | 总输赢分数 | 分（前端转元） |
| 金额 | `total_fee` | 总手续费 | 分（前端转元） |

### 金额转换

```typescript
// 后端存储：分
total_score: 125050  // 1,250.50 元
total_fee: 12550     // 125.50 元

// 前端显示：元
formatScore(125050) = "1250.50"
formatFee(12550) = "125.50"
```

---

## 📊 API 使用

### 接口：`POST /battle-query/house/stats`

#### 今日统计请求
```json
{
  "house_gid": 60870,
  "start_time": 1701388800,  // 今天0点
  "end_time": 1701430000     // 当前时间
}
```

#### 昨日统计请求
```json
{
  "house_gid": 60870,
  "start_time": 1701302400,  // 昨天0点
  "end_time": 1701388800     // 今天0点
}
```

#### 上周统计请求
```json
{
  "house_gid": 60870,
  "start_time": 1700438400,  // 上周一0点
  "end_time": 1701043200     // 上周日24点
}
```

#### 响应
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "house_gid": 60870,
    "total_games": 125,      // 笔数
    "total_score": 125050,   // 成绩（分）
    "total_fee": 12550       // 金额（分）
  }
}
```

---

## 🔄 对比：新实现 vs 旧实现

| 项目 | 旧实现 | 新实现 |
|-----|-------|-------|
| 店铺管理员看什么 | 自己圈子的统计 | 整个店铺的统计 |
| 时间维度 | 昨日、本周、上周 | **今日**、昨日、上周 |
| 数据字段 | 局数、得分、手续费、**活跃人数** | 笔数、成绩、金额 |
| 使用的 API | `getGroupStats` | `getHouseStats` |
| 与机器人对应 | ❌ 不对应 | ✅ 完全对应 |

---

## 🧪 测试建议

### 功能测试

#### 店铺管理员
- [ ] 首次进入显示店铺统计
- [ ] 今日统计时间范围正确（今天0点到现在）
- [ ] 昨日统计时间范围正确（昨天0点到今天0点）
- [ ] 上周统计时间范围正确（上周一到上周日）
- [ ] 金额显示正确（分转元）
- [ ] 正负分颜色正确（正数绿色，负数红色）
- [ ] 下拉刷新重新加载

#### 超级管理员
- [ ] 显示全平台统计（所有店铺汇总）
- [ ] 三个时间段数据正确汇总
- [ ] 多个店铺数据正确累加

### 边界测试
- [ ] 今天刚开始（0点）时的今日统计
- [ ] 周一查看上周统计
- [ ] 新店铺无数据
- [ ] 单个店铺有数据，其他店铺无数据
- [ ] 网络请求失败处理

---

## ✅ 完成状态

- ✅ 后端 API 已有（`getHouseStats`）
- ✅ 前端组件已更新（`HouseStatsCard`）
- ✅ 首页集成已完成
- ✅ 店铺管理员视图已实现（店铺统计）
- ✅ 超级管理员视图已实现（全平台统计）
- ✅ 时间维度已更新（今日、昨日、上周）
- ✅ 与微信机器人命令完全对应

功能现已上线可用！🎉

---

## 📝 重要说明

### 1. 数据范围变更
**旧实现**：店铺管理员看自己圈子的统计
**新实现**：店铺管理员看整个店铺的统计（所有圈子汇总）

### 2. 时间维度变更
**旧实现**：昨日、本周、上周
**新实现**：**今日**、昨日、上周（更符合机器人命令）

### 3. 字段名称
- **笔数** = `total_games`
- **成绩** = `total_score`
- **金额** = `total_fee`

### 4. 金额单位
后端统一使用"分"，前端显示时转换为"元"，保留2位小数。
