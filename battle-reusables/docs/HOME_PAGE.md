# 首页功能说明

## 概述

首页已完成基础功能集成，可以显示店铺管理员的圈子信息、成员列表和系统公告。

---

## ✅ 已实现功能

### 1. 用户欢迎卡片
- ✅ 显示当前登录用户信息
- ✅ 显示认证状态
- ✅ 显示店铺信息（店铺管理员）
- ✅ 显示圈子名称（店铺管理员）

### 2. 系统公告
- ✅ 展示系统公告列表
- ✅ 支持重要程度标记（普通/重要/紧急）
- ✅ 只显示前3条，可点击查看全部
- ⏳ 当前使用模拟数据（需要后端公告管理系统）

### 3. 圈子战绩统计
- ✅ 显示圈子名称
- ✅ 战绩卡片UI已完成
- ⏳ 今日战绩统计（需要后端API）
- ⏳ 本周战绩统计（需要后端API）
- ⏳ 总对局数统计（需要后端API）
- **仅店铺管理员可见**

### 4. 成员在线情况
- ✅ 显示圈子总成员数
- ✅ 显示成员列表
- ✅ 成员昵称显示
- ⏳ 在线状态检测（需要后端WebSocket或轮询API）
- ⏳ 在线率统计（需要在线状态数据）
- **仅店铺管理员可见**

### 5. 通用功能
- ✅ 下拉刷新
- ✅ 自动加载数据
- ✅ 加载状态提示
- ✅ 权限控制

---

## 📁 文件结构

```
battle-reusables/
├── app/(tabs)/index.tsx              # 首页主文件
├── services/home/index.ts            # 首页数据服务
└── components/(tabs)/home/
    ├── announcements-card.tsx        # 系统公告卡片
    ├── group-battles-card.tsx        # 圈子战绩卡片
    └── members-online-card.tsx       # 成员在线卡片
```

---

## 🔌 使用的API

### 复用现有API

| API | 说明 | 状态 |
|-----|------|------|
| `POST /shops/admins/me` | 获取管理员信息 | ✅ 已集成 |
| `POST /groups/my` | 获取我的圈子 | ✅ 已集成 |
| `POST /groups/members/list` | 获取圈子成员列表 | ✅ 已集成 |

---

## 📊 数据流

```
用户登录 → 首页加载
    ↓
1. 检查是否为店铺管理员
    ↓
2. 获取管理员信息 (shopsAdminsMe)
    ↓
3. 获取我的圈子 (getMyGroupInfo)
    ↓
4. 获取圈子成员 (getGroupMembers)
    ↓
5. 渲染页面
```

---

## ⏳ 待完善功能

### 高优先级

**1. 圈子战绩统计API**
- 需要创建后端API：`GET /api/battle/group-stats`
- 返回数据：
  ```json
  {
    "today_battles": 12,
    "today_wins": 8,
    "today_losses": 4,
    "week_battles": 45,
    "week_wins": 30,
    "week_winrate": 66.7,
    "total_battles": 156
  }
  ```

**2. 成员在线状态**
- 方案A：WebSocket实时推送
- 方案B：轮询API `GET /api/members/online-status`
- 需要后端记录成员最后活跃时间

### 中优先级

**3. 系统公告管理**
- 创建公告管理后台
- API：`GET /api/announcements`
- 支持CRUD操作
- 支持优先级设置

**4. 快捷操作**
- 添加"快速上分"按钮
- 添加"查看战绩"跳转
- 添加"管理成员"跳转

### 低优先级

**5. 数据可视化**
- 胜率趋势图表
- 战绩热力图
- 成员活跃度分析

---

## 🚀 使用方法

### 店铺管理员

1. 登录系统
2. 进入首页（默认首页）
3. 查看内容：
   - ✅ 欢迎信息
   - ✅ 系统公告
   - ✅ 圈子战绩（基础UI）
   - ✅ 成员列表

4. 下拉刷新获取最新数据

### 普通用户

1. 登录系统
2. 进入首页
3. 查看内容：
   - ✅ 欢迎信息
   - ✅ 系统公告

---

## 🔧 配置说明

### 数据刷新策略

```typescript
// 自动加载：登录后且为店铺管理员时
useEffect(() => {
  if (isAuthenticated && isStoreAdmin) {
    loadData();
  }
}, [isAuthenticated, isStoreAdmin]);

// 手动刷新：下拉刷新
const handleRefresh = async () => {
  setRefreshing(true);
  await loadData();
  setRefreshing(false);
};
```

### 权限控制

```typescript
// 圈子战绩 - 仅店铺管理员
{isStoreAdmin && <GroupBattlesCard />}

// 成员在线 - 需要权限
{isStoreAdmin && (
  <PermissionGate anyOf={['shop:member:view']}>
    <MembersOnlineCard />
  </PermissionGate>
)}
```

---

## 📝 开发建议

### 完善战绩统计

1. 创建后端Use Case：
```go
// internal/biz/game/battle_stats.go
func (uc *BattleStatsUseCase) GetGroupDailyStats(
    ctx context.Context,
    houseGID int32,
    groupID int32,
) (*GroupDailyStats, error) {
    // 统计今日战绩
    // 统计本周战绩
    // 统计总对局
}
```

2. 创建Service层：
```go
// internal/service/game/battle_stats.go
func (s *BattleStatsService) GetGroupStats(c *gin.Context) {
    // 调用Use Case
    // 返回统计数据
}
```

3. 前端集成：
```typescript
// services/battle/index.ts
export function getGroupStats(data: { house_gid: number; group_id: number }) {
  return post('/api/battle/group-stats', data);
}
```

### 添加在线状态

1. 方案一：WebSocket（推荐）
```go
// 服务端推送在线状态变更
type OnlineStatusUpdate struct {
    MemberID  int32 `json:"member_id"`
    IsOnline  bool  `json:"is_online"`
    Timestamp int64 `json:"timestamp"`
}
```

2. 方案二：定时轮询
```typescript
// 每30秒查询一次在线状态
setInterval(async () => {
  const status = await getOnlineStatus({ group_id: myGroup.id });
  updateMemberStatus(status);
}, 30000);
```

---

## 🐛 已知问题

1. ⚠️ 战绩数据全部显示为0（需要后端API）
2. ⚠️ 成员在线状态全部为离线（需要后端支持）
3. ⚠️ 系统公告使用硬编码数据（需要公告管理系统）

---

## ✨ 未来优化

1. **性能优化**
   - 实现数据缓存
   - 减少不必要的API调用
   - 使用虚拟列表优化长列表

2. **用户体验**
   - 添加骨架屏
   - 优化加载动画
   - 添加错误重试机制

3. **功能扩展**
   - 添加消息通知
   - 支持多圈子切换
   - 个性化首页配置

---

**文档版本**: v1.0  
**更新时间**: 2025-11-22  
**状态**: ✅ 基础功能已完成，部分高级功能待实现
