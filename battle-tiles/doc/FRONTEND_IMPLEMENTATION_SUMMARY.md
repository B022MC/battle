# 前端实现总结 - 战绩查询功能

## 📋 概述

本文档总结了战绩查询功能的前端实现,包括独立圈子余额查询功能。

## ✅ 已完成的工作

### 1. API 服务层

#### 文件: `battle-reusables/services/battles/query.ts`

实现了 7 个 API 调用函数:

**普通用户接口**:
- `listMyBattles()` - 查询我的战绩
- `getMyBalances()` - 查询我的余额
- `getMyStats()` - 查询我的统计

**管理员接口**:
- `listGroupBattles()` - 查询圈子战绩
- `listGroupMemberBalances()` - 查询圈子成员余额
- `getGroupStats()` - 查询圈子统计

**超级管理员接口**:
- `getHouseStats()` - 查询店铺统计

#### 文件: `battle-reusables/services/battles/query-typing.d.ts`

定义了所有请求和响应的 TypeScript 类型:
- 请求参数类型: `ListMyBattlesParams`, `GetMyBalancesParams`, 等
- 响应数据类型: `BattleRecord`, `MemberBalance`, `BattleStats`, `GroupStats`, `HouseStats`

### 2. 页面路由

创建了 4 个新页面:

| 文件路径 | 路由 | 说明 |
|---------|------|------|
| `app/(shop)/my-battles.tsx` | `/my-battles` | 我的战绩页面 |
| `app/(shop)/my-balances.tsx` | `/my-balances` | 我的余额页面 |
| `app/(shop)/group-battles.tsx` | `/group-battles` | 圈子战绩页面(管理员) |
| `app/(shop)/group-balances.tsx` | `/group-balances` | 圈子余额页面(管理员) |

### 3. UI 组件

创建了 4 个主要视图组件:

#### `components/(shop)/battles/my-battles-view.tsx`

**功能**:
- 查询我的战绩列表
- 查询我的战绩统计
- 支持按圈子筛选
- 支持时间范围筛选
- 分页显示

**UI 特性**:
- 战绩列表卡片展示
- 分数颜色区分(正分绿色,负分红色)
- 统计信息卡片
- 分页控件

#### `components/(shop)/battles/my-balances-view.tsx`

**功能**:
- 查询我的余额
- 支持查询所有圈子或指定圈子
- 显示总余额

**UI 特性**:
- 总余额大卡片展示
- 余额明细列表
- 余额以元和分两种单位显示
- 使用说明提示

#### `components/(shop)/battles/group-battles-view.tsx`

**功能**:
- 管理员查询圈子战绩
- 支持按成员筛选
- 支持时间范围筛选
- 查询圈子统计
- 权限验证(仅管理员可访问)

**UI 特性**:
- 权限检查提示
- 战绩列表展示
- 圈子统计信息
- 分页控件

#### `components/(shop)/battles/group-balances-view.tsx`

**功能**:
- 管理员查询圈子成员余额
- 支持按余额范围筛选
- 显示统计信息(总余额、平均余额)
- 权限验证(仅管理员可访问)

**UI 特性**:
- 权限检查提示
- 统计信息卡片(成员数、总余额、平均余额)
- 成员余额列表
- 分页控件
- 使用说明提示

## 🎨 UI/UX 设计特点

### 1. 一致的表单设计
- 所有查询页面使用统一的表单布局
- 必填字段标注 `*`
- 可选字段有说明文字
- 输入框有占位符提示

### 2. 数据展示
- 使用 Card 组件包裹内容区域
- 列表项使用边框分隔
- 重要数据(余额、分数)使用大字体和颜色强调
- 时间格式统一为 `YYYY-MM-DD HH:mm`

### 3. 交互反馈
- 加载状态显示在按钮上
- 空数据状态有友好提示
- 错误通过 alert 提示
- 分页控件禁用状态处理

### 4. 颜色语义
- 正分数: 绿色 (`text-green-600`)
- 负分数: 红色 (`text-red-600`)
- 主要数据: 主题色 (`text-primary`)
- 次要信息: 灰色 (`text-muted-foreground`)

## 📱 功能特性

### 普通用户功能

#### 1. 查询我的战绩
- **路由**: `/my-battles`
- **筛选条件**:
  - 店铺号(必填)
  - 圈子ID(可选,不填查询所有圈子)
  - 时间范围(可选)
- **显示内容**:
  - 战绩列表(圈子名、时间、分数、费用、房间号、余额)
  - 统计信息(总局数、总分数、总费用、平均分)
  - 分页(每页20条)

#### 2. 查询我的余额
- **路由**: `/my-balances`
- **筛选条件**:
  - 店铺号(必填)
  - 圈子ID(可选,不填查询所有圈子)
- **显示内容**:
  - 总余额(所有圈子汇总)
  - 余额明细(每个圈子的余额)
  - 余额以元和分两种单位显示

### 管理员功能

#### 3. 查询圈子战绩
- **路由**: `/group-battles`
- **权限**: 仅店铺管理员
- **筛选条件**:
  - 店铺号(必填)
  - 圈子ID(必填)
  - 玩家游戏ID(可选,不填查询所有成员)
  - 时间范围(可选)
- **显示内容**:
  - 战绩列表(玩家ID、时间、分数、费用、房间号、余额)
  - 圈子统计(总局数、总分数、总费用、活跃成员数)
  - 分页(每页20条)

#### 4. 查询圈子成员余额
- **路由**: `/group-balances`
- **权限**: 仅店铺管理员
- **筛选条件**:
  - 店铺号(必填)
  - 圈子ID(必填)
  - 余额范围(可选,最小/最大余额)
- **显示内容**:
  - 统计信息(成员数、总余额、平均余额)
  - 成员余额列表(昵称、游戏ID、成员ID、余额)
  - 分页(每页20条)

## 🔧 技术实现细节

### 1. 时间处理
```typescript
// 将用户输入的时间字符串转换为 Unix timestamp
const timestamp = new Date(startTime).getTime() / 1000;
params.start_time = Math.floor(timestamp);
```

### 2. 余额单位转换
```typescript
// 后端返回分(int32),前端显示元(float64)
balance_yuan: number; // 元
balance: number;      // 分
```

### 3. 权限验证
```typescript
const { isStoreAdmin } = usePermission();

if (!isStoreAdmin) {
  return <Text>仅店铺管理员可访问此页面</Text>;
}
```

### 4. 分页处理
```typescript
const [page, setPage] = useState(1);

// 上一页
if (page > 1) {
  setPage(page - 1);
  handleLoadBattles();
}

// 下一页
if (page * 20 < battlesData.total) {
  setPage(page + 1);
  handleLoadBattles();
}
```

## 📝 待完成的工作

### 1. 导航菜单集成
需要在主导航菜单中添加新页面的入口:
- 在 `app/(shop)/_layout.tsx` 中添加路由配置
- 在侧边栏或底部导航中添加菜单项

### 2. 超级管理员功能
虽然后端已实现,但前端尚未创建超级管理员查询店铺统计的页面:
- 创建 `app/(shop)/house-stats.tsx`
- 创建 `components/(shop)/battles/house-stats-view.tsx`

### 3. 日期选择器
当前使用文本输入框输入时间,可以改进为:
- 集成日期时间选择器组件
- 提供快捷时间范围选择(今天、本周、本月等)

### 4. 数据导出
可以添加导出功能:
- 导出战绩为 CSV/Excel
- 导出余额为 CSV/Excel

### 5. 图表展示
可以添加数据可视化:
- 战绩趋势图
- 余额分布图
- 圈子对比图

### 6. 实时刷新
可以添加自动刷新功能:
- 定时刷新数据
- 下拉刷新

## 🧪 测试建议

### 1. 功能测试
- [ ] 测试所有筛选条件组合
- [ ] 测试分页功能
- [ ] 测试空数据状态
- [ ] 测试错误处理

### 2. 权限测试
- [ ] 测试普通用户访问管理员页面
- [ ] 测试管理员访问自己的圈子
- [ ] 测试管理员访问其他圈子

### 3. UI 测试
- [ ] 测试不同屏幕尺寸
- [ ] 测试长文本显示
- [ ] 测试加载状态
- [ ] 测试交互反馈

### 4. 性能测试
- [ ] 测试大数据量列表渲染
- [ ] 测试快速切换页面
- [ ] 测试并发请求

## 📚 相关文档

- [后端实现总结](./BATTLE_QUERY_IMPLEMENTATION_SUMMARY.md)
- [后端测试指南](./BATTLE_QUERY_TESTING_GUIDE.md)
- [独立圈子余额设计](./INDEPENDENT_GROUP_BALANCE_DESIGN.md)
- [数据库迁移脚本](./migration_independent_group_balance.sql)

## 🎯 使用示例

### 普通用户查询自己的战绩
1. 打开 `/my-battles` 页面
2. 输入店铺号
3. (可选)输入圈子ID
4. (可选)选择时间范围
5. 点击"查询战绩"按钮

### 管理员查询圈子成员余额
1. 打开 `/group-balances` 页面
2. 输入店铺号
3. 输入圈子ID
4. (可选)设置余额范围
5. 点击"查询余额"按钮

## 🔗 API 端点映射

| 前端函数 | 后端端点 | 说明 |
|---------|---------|------|
| `listMyBattles()` | `POST /battle-query/my/battles` | 查询我的战绩 |
| `getMyBalances()` | `POST /battle-query/my/balances` | 查询我的余额 |
| `getMyStats()` | `POST /battle-query/my/stats` | 查询我的统计 |
| `listGroupBattles()` | `POST /battle-query/group/battles` | 查询圈子战绩 |
| `listGroupMemberBalances()` | `POST /battle-query/group/balances` | 查询圈子成员余额 |
| `getGroupStats()` | `POST /battle-query/group/stats` | 查询圈子统计 |
| `getHouseStats()` | `POST /battle-query/house/stats` | 查询店铺统计 |

