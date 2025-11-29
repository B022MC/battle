# 超级管理员首页全平台统计功能

## ✅ 功能完成

在超级管理员首页添加了全平台战绩统计功能，显示所有店铺的汇总数据：
- 📊 **昨日统计**（所有店铺汇总）
- 📈 **本周统计**（所有店铺汇总）
- 📉 **上周统计**（所有店铺汇总）

---

## 🎯 功能特性

### 1. **角色区分**

#### 超级管理员首页
- 显示全平台统计卡片
- 汇总所有店铺的数据
- 包含昨日、本周、上周三个时间维度

#### 店铺管理员首页
- 显示圈子战绩统计卡片
- 只显示自己圈子的数据
- 同样包含昨日、本周、上周三个时间维度

### 2. **全平台统计卡片**

显示内容：
- **总局数**：所有店铺的对局总和
- **总得分**：所有店铺的得分总和（元）
- **总手续费**：所有店铺的手续费总和（元）

每个时间段独立显示：
- 昨日统计（蓝色卡片）
- 本周统计（绿色卡片）
- 上周统计（紫色卡片）

### 3. **数据加载逻辑**

#### 超级管理员
```
1. 获取所有店铺列表 (shopsHousesOptions)
2. 为每个店铺并发请求三个时间段的统计
3. 汇总所有店铺的数据
4. 显示在首页卡片中
```

#### 店铺管理员
```
1. 获取管理员信息 (shopsAdminsMe)
2. 获取自己的圈子 (getMyGroupInfo)
3. 为圈子请求三个时间段的统计
4. 显示在首页卡片中
```

---

## 📁 文件清单

### 1. 新增文件

#### `/components/(tabs)/home/house-stats-card.tsx`
**全平台统计卡片组件**：
- 专为超级管理员设计
- 接收多个店铺的统计数据
- 自动汇总计算
- 支持昨日、本周、上周三个维度

**主要属性**：
```typescript
interface HousePeriodStats {
  house_gid: number;
  total_games: number;
  total_score: number;
  total_fee: number;
}

interface HouseStatsCardProps {
  yesterdayStats?: HousePeriodStats[];  // 所有店铺的昨日统计
  thisWeekStats?: HousePeriodStats[];   // 所有店铺的本周统计
  lastWeekStats?: HousePeriodStats[];   // 所有店铺的上周统计
  loading?: boolean;
  onViewDetails?: () => void;
}
```

### 2. 修改文件

#### `/app/(tabs)/index.tsx`
**首页主文件**：

**新增导入**：
```typescript
import { HouseStatsCard } from '@/components/(tabs)/home/house-stats-card';
import { shopsHousesOptions } from '@/services/shops/houses';
import { getHouseStats } from '@/services/battles/query';
import type { HouseStats } from '@/services/battles/query-typing';
```

**新增状态**：
```typescript
const { isSuperAdmin, isStoreAdmin } = usePermission();

// 全平台统计（超级管理员）
const [allHouses, setAllHouses] = useState<number[]>([]);
const [houseYesterdayStats, setHouseYesterdayStats] = useState<HouseStats[]>([]);
const [houseThisWeekStats, setHouseThatsStats] = useState<HouseStats[]>([]);
const [houseLastWeekStats, setHouseLastWeekStats] = useState<HouseStats[]>([]);
const [loadingHouseStats, setLoadingHouseStats] = useState(false);
```

**新增函数**：
- `loadSuperAdminData()` - 加载超级管理员数据
- `loadStoreAdminData()` - 加载店铺管理员数据（重命名自 `loadData`）
- `loadAllHousesStats()` - 加载全平台统计

---

## 🔧 技术实现

### 数据流程

#### 超级管理员视图

```
用户进入首页（超级管理员）
    ↓
检测到 isSuperAdmin = true
    ↓
调用 loadSuperAdminData()
    ↓
获取所有店铺列表 (shopsHousesOptions)
    ↓
并发请求所有店铺的统计数据
    ├─→ 店铺1：昨日、本周、上周
    ├─→ 店铺2：昨日、本周、上周
    └─→ 店铺N：昨日、本周、上周
    ↓
分类汇总结果
    ├─→ yesterdayStats: HouseStats[]
    ├─→ thisWeekStats: HouseStats[]
    └─→ lastWeekStats: HouseStats[]
    ↓
HouseStatsCard 自动计算汇总
    ↓
显示全平台统计卡片
```

#### 店铺管理员视图

```
用户进入首页（店铺管理员）
    ↓
检测到 isStoreAdmin = true
    ↓
调用 loadStoreAdminData()
    ↓
获取管理员信息 → 获取圈子信息
    ↓
并发请求圈子的统计数据
    ├─→ 昨日统计 (getGroupStats)
    ├─→ 本周统计 (getGroupStats)
    └─→ 上周统计 (getGroupStats)
    ↓
显示圈子战绩统计卡片
```

### 并发优化

**超级管理员**：
```typescript
// 假设有3个店铺，总共需要 3 * 3 = 9 个请求
const allRequests = houses.flatMap(houseGid => [
  getHouseStats({ house_gid: houseGid, ...getTimeRange('yesterday') }),
  getHouseStats({ house_gid: houseGid, ...getTimeRange('thisWeek') }),
  getHouseStats({ house_gid: houseGid, ...getTimeRange('lastWeek') })
]);

// 全部并发执行
const results = await Promise.all(allRequests);
```

**性能优势**：
- 9个请求并发执行，而非顺序执行
- 总耗时 ≈ 单个请求耗时（假设网络带宽充足）

### 数据汇总逻辑

```typescript
const calculateTotal = (stats?: HousePeriodStats[]) => {
  if (!stats || stats.length === 0) {
    return { total_games: 0, total_score: 0, total_fee: 0 };
  }
  return stats.reduce(
    (acc, item) => ({
      total_games: acc.total_games + item.total_games,
      total_score: acc.total_score + item.total_score,
      total_fee: acc.total_fee + item.total_fee,
    }),
    { total_games: 0, total_score: 0, total_fee: 0 }
  );
};
```

---

## 🎨 UI 展示

### 超级管理员首页

```
┌─────────────────────────────────┐
│ 🏆 欢迎回来！                   │
│ 用户: admin                     │
│ 认证状态: 已登录                │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 📢 系统公告                     │
│ ...                             │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 🏢 全平台统计      查看详情 >   │
│ 所有店铺数据汇总                │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 昨日统计           125 场│ │
│ │ 总得分         总手续费     │ │
│ │ +1,250.50      125.50       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 本周统计           450 场│ │
│ │ 总得分         总手续费     │ │
│ │ +3,800.00      450.00       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 上周统计           520 场│ │
│ │ 总得分         总手续费     │ │
│ │ -1,255.50      520.00       │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

### 店铺管理员首页

```
┌─────────────────────────────────┐
│ 🏆 欢迎回来！                   │
│ 用户: shop_admin                │
│ 认证状态: 已登录                │
│ 店铺: 60870                     │
│ 我的圈子: b022mc的圈子          │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 📢 系统公告                     │
│ ...                             │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 🏆 圈子战绩统计    查看详情 >   │
│ 圈子：b022mc的圈子              │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 📅 昨日统计            12 场│ │
│ │ 得分    手续费    活跃人数  │ │
│ │ +125.50  12.50      8       │ │
│ └─────────────────────────────┘ │
│ ... (本周、上周)                │
└─────────────────────────────────┘
```

---

## 📊 API 使用

### 1. 获取店铺列表

**接口**：`GET /shops/houses/options`

**响应**：
```json
{
  "code": 0,
  "msg": "成功",
  "data": [60870, 60871, 60872]
}
```

### 2. 获取店铺统计

**接口**：`POST /battle-query/house/stats`

**请求**：
```json
{
  "house_gid": 60870,
  "start_time": 1701302400,  // 昨天0点
  "end_time": 1701388800     // 今天0点
}
```

**响应**：
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "house_gid": 60870,
    "total_games": 125,
    "total_score": 125050,  // 分
    "total_fee": 12550      // 分
  }
}
```

### 3. 获取圈子统计（店铺管理员）

**接口**：`POST /battle-query/group/stats`

**请求**：
```json
{
  "house_gid": 60870,
  "group_id": 1,
  "start_time": 1701302400,
  "end_time": 1701388800
}
```

**响应**：
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "group_id": 1,
    "group_name": "b022mc的圈子",
    "total_games": 12,
    "total_score": 12550,
    "total_fee": 1250,
    "active_members": 8
  }
}
```

---

## 🔐 权限控制

### 超级管理员
- ✅ 可以查看全平台统计
- ✅ 可以查看所有店铺列表
- ✅ 可以查看任意店铺的详细统计
- ✅ 具有最高权限，可以访问所有数据

### 店铺管理员
- ✅ 只能查看自己圈子的统计
- ❌ 不能查看其他圈子的数据
- ❌ 不能查看全平台统计
- ✅ 需要 `battles:view` 权限

---

## 🧪 测试建议

### 功能测试

#### 超级管理员
- [ ] 首次进入加载全平台统计
- [ ] 下拉刷新重新加载数据
- [ ] 显示所有店铺的汇总数据
- [ ] 昨日、本周、上周数据正确
- [ ] 多个店铺数据正确汇总
- [ ] 无数据时显示"暂无统计数据"

#### 店铺管理员
- [ ] 首次进入加载圈子统计
- [ ] 下拉刷新重新加载数据
- [ ] 只显示自己圈子的数据
- [ ] 昨日、本周、上周数据正确
- [ ] 活跃成员数显示正确

### 性能测试
- [ ] 多店铺并发请求性能
- [ ] 大数据量渲染性能
- [ ] 刷新操作流畅度

### 边界测试
- [ ] 无店铺时的处理
- [ ] 单个店铺的情况
- [ ] 店铺统计数据为空
- [ ] 网络请求失败处理

---

## 🔄 对比：超级管理员 vs 店铺管理员

| 功能 | 超级管理员 | 店铺管理员 |
|-----|----------|----------|
| 首页统计卡片 | 全平台统计 | 圈子战绩统计 |
| 数据范围 | 所有店铺 | 单个圈子 |
| 统计维度 | 总局数、总得分、总手续费 | 总局数、得分、手续费、活跃人数 |
| 时间维度 | 昨日、本周、上周 | 昨日、本周、上周 |
| 数据来源 | `getHouseStats` | `getGroupStats` |
| 权限要求 | 超级管理员角色 | 店铺管理员角色 + `battles:view` |

---

## 💡 后续优化建议

1. **店铺筛选**：超级管理员可以选择查看特定店铺的详细数据
2. **图表展示**：添加趋势图表和对比图
3. **排行榜**：显示店铺排行榜
4. **导出功能**：支持导出统计报表
5. **实时更新**：WebSocket 推送最新数据
6. **详细视图**：点击卡片查看更详细的统计信息

---

## ✅ 完成状态

- ✅ 后端 API 已有（复用 `getHouseStats`）
- ✅ 前端组件已创建（`HouseStatsCard`）
- ✅ 首页集成已完成
- ✅ 超级管理员视图已实现
- ✅ 店铺管理员视图保持不变
- ✅ 角色区分逻辑已完成
- ✅ 并发请求优化已完成
- ✅ 数据汇总逻辑已实现

功能现已上线可用！🎉
