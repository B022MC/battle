# 店铺管理员首页圈子统计功能

## ✅ 功能完成

在店铺管理员首页添加了圈子战绩统计功能，支持查看：
- 📊 **昨日统计**
- 📈 **本周统计**  
- 📉 **上周统计**

与微信机器人命令功能一致。

---

## 🎯 功能特性

### 1. **三个时间维度**

#### 昨日统计
- 时间范围：昨天 0:00 - 今天 0:00
- 颜色主题：蓝色
- 显示数据：总局数、得分、手续费、活跃人数

#### 本周统计
- 时间范围：本周一 0:00 - 当前时间
- 颜色主题：绿色
- 显示数据：总局数、得分、手续费、活跃人数

#### 上周统计
- 时间范围：上周一 0:00 - 上周日 24:00
- 颜色主题：紫色
- 显示数据：总局数、得分、手续费、活跃人数

### 2. **数据展示**

每个统计卡片显示：
- **总局数**：该时间段内的总对局数
- **得分**：总输赢分数（元），正数显示绿色，负数显示红色
- **手续费**：总手续费（元），橙色显示
- **活跃人数**：参与游戏的成员数量

### 3. **自动刷新**

- 首次进入自动加载
- 下拉刷新重新加载
- 三个时间段并发请求，提升加载速度

---

## 📁 文件修改

### 1. 前端组件

#### `/components/(tabs)/home/group-battles-card.tsx`
**重构战绩卡片组件**：
- 原有设计：单一统计结构（今日、本周、总计）
- 新设计：三个独立时间段（昨日、本周、上周）
- 支持分转元的金额格式化
- 颜色区分不同时间段

**主要变化**：
```typescript
interface PeriodStats {
  total_games: number;      // 总局数
  total_score: number;      // 总得分（分）
  total_fee: number;        // 总手续费（分）
  active_members: number;   // 活跃成员数
}

interface GroupBattlesCardProps {
  groupName?: string;
  yesterdayStats?: PeriodStats;  // 昨日统计
  thisWeekStats?: PeriodStats;   // 本周统计
  lastWeekStats?: PeriodStats;   // 上周统计
  loading?: boolean;
  onViewDetails?: () => void;
}
```

#### `/app/(tabs)/index.tsx`
**首页集成统计功能**：
- 添加时间范围计算函数 `getTimeRange()`
- 添加统计数据加载函数 `loadGroupStats()`
- 并发请求三个时间段数据
- 传递数据给 GroupBattlesCard 组件

**时间范围计算**：
```typescript
function getTimeRange(type: 'yesterday' | 'thisWeek' | 'lastWeek'): {
  start_time: number;
  end_time: number;
}
```

### 2. 后端接口（已有）

**接口**：`POST /battle-query/group/stats`

**请求参数**：
```typescript
{
  house_gid: number;
  group_id: number;
  start_time?: number;  // Unix timestamp（秒）
  end_time?: number;    // Unix timestamp（秒）
}
```

**响应数据**：
```typescript
{
  code: 0,
  msg: "成功",
  data: {
    group_id: number;
    group_name: string;
    total_games: number;       // 总局数
    total_score: number;       // 总得分（分）
    total_fee: number;         // 总手续费（分）
    active_members: number;    // 活跃成员数
  }
}
```

---

## 🔧 技术实现

### 时间计算逻辑

#### 昨日
```typescript
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
const yesterdayStart = new Date(today);
yesterdayStart.setDate(yesterdayStart.getDate() - 1);

// 昨天 0:00 - 今天 0:00
start_time: Math.floor(yesterdayStart.getTime() / 1000)
end_time: Math.floor(today.getTime() / 1000)
```

#### 本周
```typescript
const dayOfWeek = today.getDay() || 7; // 周日转为7
const thisWeekStart = new Date(today);
thisWeekStart.setDate(thisWeekStart.getDate() - dayOfWeek + 1); // 周一

// 本周一 0:00 - 现在
start_time: Math.floor(thisWeekStart.getTime() / 1000)
end_time: Math.floor(now.getTime() / 1000)
```

#### 上周
```typescript
const dayOfWeek = today.getDay() || 7;
const lastWeekStart = new Date(today);
lastWeekStart.setDate(lastWeekStart.getDate() - dayOfWeek - 6); // 上周一
const lastWeekEnd = new Date(lastWeekStart);
lastWeekEnd.setDate(lastWeekEnd.getDate() + 7); // 上周日+1天

// 上周一 0:00 - 上周日 24:00
start_time: Math.floor(lastWeekStart.getTime() / 1000)
end_time: Math.floor(lastWeekEnd.getTime() / 1000)
```

### 并发请求优化

使用 `Promise.all()` 并发请求三个时间段，减少总等待时间：

```typescript
const [yesterday, thisWeek, lastWeek] = await Promise.all([
  getGroupStats({ house_gid, group_id, ...getTimeRange('yesterday') }),
  getGroupStats({ house_gid, group_id, ...getTimeRange('thisWeek') }),
  getGroupStats({ house_gid, group_id, ...getTimeRange('lastWeek') })
]);
```

---

## 🎨 UI 设计

### 卡片样式

每个统计卡片采用不同颜色主题：

1. **昨日统计**（蓝色）
   - 背景：`bg-blue-50 dark:bg-blue-950`
   - 文字：`text-blue-900 dark:text-blue-100`
   - 强调色：`text-blue-700 dark:text-blue-300`

2. **本周统计**（绿色）
   - 背景：`bg-green-50 dark:bg-green-950`
   - 文字：`text-green-900 dark:text-green-100`
   - 强调色：`text-green-700 dark:text-green-300`

3. **上周统计**（紫色）
   - 背景：`bg-purple-50 dark:bg-purple-950`
   - 文字：`text-purple-900 dark:text-purple-100`
   - 强调色：`text-purple-700 dark:text-purple-300`

### 数据格式化

#### 金额显示
```typescript
// 分转元，保留2位小数
const formatScore = (score: number) => {
  return (score / 100).toFixed(2);
};

// 得分显示（带正负号）
{stats.total_score >= 0 ? '+' : ''}{formatScore(stats.total_score)}
```

#### 颜色指示
- 正分（赢）：绿色 `text-green-600`
- 负分（输）：红色 `text-red-600`
- 手续费：橙色 `text-orange-600`
- 活跃人数：蓝色 `text-blue-600`

---

## 📊 数据流程

```
用户进入首页
    ↓
加载店铺管理员信息 (shopsAdminsMe)
    ↓
加载我的圈子信息 (getMyGroupInfo)
    ↓
并发请求三个时间段统计
    ├─→ 昨日统计 (getGroupStats)
    ├─→ 本周统计 (getGroupStats)
    └─→ 上周统计 (getGroupStats)
    ↓
更新状态并渲染卡片
```

---

## 🔐 权限控制

- **显示条件**：`isStoreAdmin` 为 true
- **API 权限**：需要 `battles:view` 权限
- **数据隔离**：只能查看自己圈子的数据

---

## 🧪 测试建议

### 功能测试
- [ ] 首次进入首页加载统计数据
- [ ] 下拉刷新重新加载数据
- [ ] 昨日统计时间范围正确
- [ ] 本周统计时间范围正确（周一到今天）
- [ ] 上周统计时间范围正确（上周一到上周日）
- [ ] 金额显示正确（分转元）
- [ ] 正负分颜色显示正确
- [ ] 无数据时显示"暂无战绩数据"

### 边界测试
- [ ] 周一查看本周统计（应该只有今天的数据）
- [ ] 周日查看上周统计
- [ ] 新圈子无战绩数据
- [ ] 加载失败后的错误处理

### 性能测试
- [ ] 并发请求响应时间
- [ ] 大数据量渲染性能
- [ ] 刷新操作的流畅度

---

## 📱 界面示例

```
┌─────────────────────────────────┐
│ 🏆 圈子战绩统计      查看详情 > │
│ 圈子：b022mc的圈子              │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 昨日统计            12 场│ │
│ │ 得分      手续费    活跃人数│ │
│ │ +125.50   12.50       8     │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 本周统计            45 场│ │
│ │ 得分      手续费    活跃人数│ │
│ │ +380.00   45.00      15     │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 上周统计            52 场│ │
│ │ 得分      手续费    活跃人数│ │
│ │ -125.50   52.00      18     │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

---

## 🎯 与机器人命令对应

| 机器人命令 | 首页统计 | 时间范围 |
|----------|---------|---------|
| `昨日统计` | 昨日统计 | 昨天 0:00 - 今天 0:00 |
| `本周统计` | 本周统计 | 本周一 0:00 - 现在 |
| `上周统计` | 上周统计 | 上周一 0:00 - 上周日 24:00 |

**相同数据**：
- 总局数
- 总得分
- 总手续费  
- 活跃成员数

---

## 💡 后续优化建议

1. **缓存优化**：缓存统计数据，减少 API 调用
2. **图表展示**：添加趋势图表
3. **对比功能**：本周与上周数据对比
4. **导出功能**：支持导出统计报表
5. **实时更新**：WebSocket 推送最新数据
6. **更多维度**：支持自定义时间范围

---

## ✅ 完成状态

- ✅ 后端 API 已完成（复用已有接口）
- ✅ 前端组件已重构
- ✅ 首页集成已完成
- ✅ 时间计算逻辑已实现
- ✅ 并发请求优化已完成
- ✅ UI 设计已完成
- ✅ 数据格式化已完成

功能现已上线可用！🎉
