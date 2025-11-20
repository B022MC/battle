# 前端游戏账号系统改造说明

## 📋 概述

后端已完成游戏账号系统的重大改造，实现了游戏账号入圈机制。前端基本**无需修改**，因为所有改动都在后端透明处理。

---

## ✅ 已完成的前端优化

### 1. 游戏账号绑定组件优化

**文件**: `battle-reusables/components/(tabs)/profile/profile-game-account.tsx`

**改动**:
- ✅ 添加二次确认对话框
- ✅ 优化错误提示（特别是"游戏账号已被其他用户绑定"的情况）
- ✅ 使用气泡提示显示成功消息

**新增功能**:
```typescript
// 绑定时会自动处理以下情况：
// 1. 游戏账号不存在 → 创建新账号并绑定
// 2. 游戏账号存在但未绑定用户 → 直接绑定到当前用户
// 3. 游戏账号已被其他用户绑定 → 显示错误提示
```

---

## 🔄 后端透明改动（前端无需修改）

### 1. 游戏内申请审批

**API**: `POST /api/shops/game-applications/approve`

**后端新增逻辑**:
```
审批通过 → 自动查找/创建游戏账号 → 确保管理员圈子 → 游戏账号入圈
```

**前端影响**: ❌ 无需修改
- 前端只需调用审批 API
- 后端自动处理游戏账号入圈
- 用户体验无变化

---

### 2. 用户查询圈子

**API**: `POST /api/groups/my/list`

**后端新增逻辑**:
```
用户 → 查询绑定的游戏账号 → 查询游戏账号的圈子 → 返回圈子列表
```

**前端影响**: ❌ 无需修改
- API 接口不变
- 返回数据结构不变
- 查询逻辑在后端透明处理

---

### 3. 用户查询战绩

**API**: `POST /api/battle-query/my/battles`

**后端新增逻辑**:
```
用户 → 查询绑定的游戏账号 → 查询游戏账号的战绩 → 返回战绩列表
```

**前端影响**: ❌ 无需修改
- API 接口不变
- 返回数据结构不变
- 查询逻辑在后端透明处理

---

## 📝 API 接口说明

### 游戏账号相关 API

#### 1. 验证游戏账号
```typescript
POST /api/game/accounts/verify
{
  mode: 'account' | 'mobile',
  account: string,
  pwd_md5: string
}
```

#### 2. 绑定游戏账号
```typescript
POST /api/game/accounts
{
  mode: 'account' | 'mobile',
  account: string,
  pwd_md5: string,
  nickname?: string
}
```

**后端处理**:
- 如果游戏账号不存在 → 创建新账号
- 如果游戏账号存在但未绑定 → 更新 `user_id`
- 如果游戏账号已被其他用户绑定 → 返回错误

#### 3. 查询我的游戏账号
```typescript
GET /api/game/accounts/me
```

#### 4. 解绑游戏账号
```typescript
DELETE /api/game/accounts/me
```

---

### 圈子相关 API

#### 1. 查询我的圈子
```typescript
POST /api/groups/my/list
{}
```

**返回数据**:
```typescript
{
  id: number;
  house_gid: number;
  group_name: string;
  admin_user_id: number;
  status: string;
  joined_at: string;
}[]
```

---

### 战绩相关 API

#### 1. 查询我的战绩
```typescript
POST /api/battle-query/my/battles
{
  house_gid?: number;
  group_id?: number;
  start_time?: number;
  end_time?: number;
  page: number;
  size: number;
}
```

---

### 申请审批 API

#### 1. 查询游戏内申请列表
```typescript
POST /api/shops/game-applications/list
{
  house_gid: number;
}
```

#### 2. 通过申请
```typescript
POST /api/shops/game-applications/approve
{
  house_gid: number;
  message_id: number;
}
```

**后端自动处理**:
1. 查找或创建游戏账号（`user_id` 为 NULL）
2. 确保管理员有圈子
3. 将游戏账号加入圈子
4. 调用游戏 API 拉入圈子

#### 3. 拒绝申请
```typescript
POST /api/shops/game-applications/reject
{
  house_gid: number;
  message_id: number;
}
```

---

## 🎯 用户使用流程

### 流程 1: 普通用户绑定游戏账号

```
1. 用户注册平台账号
   ↓
2. 进入"个人中心" → "我的游戏账号"
   ↓
3. 输入游戏账号和密码
   ↓
4. 点击"绑定"
   ↓
5. 系统验证游戏账号
   ↓
6. 绑定成功（如果游戏账号已存在且未绑定，则直接关联）
```

### 流程 2: 游戏内玩家申请入圈

```
1. 玩家在游戏内发起申请（未注册平台账号）
   ↓
2. 店铺管理员在"申请管理"中看到申请
   ↓
3. 管理员点击"通过"
   ↓
4. 后端自动：
   - 创建游戏账号（user_id = NULL）
   - 将游戏账号加入管理员的圈子
   ↓
5. 玩家可以开始游戏
   ↓
6. 玩家注册平台账号后，绑定游戏账号
   ↓
7. 系统自动关联：游戏账号.user_id = 用户ID
   ↓
8. 用户可以查询战绩和圈子
```

### 流程 3: 用户查询圈子和战绩

```
1. 用户登录平台
   ↓
2. 查询"我的圈子"
   - 后端：用户 → 游戏账号 → 圈子
   ↓
3. 查询"我的战绩"
   - 后端：用户 → 游戏账号 → 战绩
   ↓
4. 显示结果
```

---

## 🚨 注意事项

### 1. 错误处理

前端需要处理以下错误情况：

#### 游戏账号绑定错误
```typescript
// 错误类型
- "you have already bound a game account" - 用户已绑定游戏账号
- "this game account is already bound to another user" - 游戏账号已被其他用户绑定
- "账号或密码不正确" - 游戏账号验证失败
```

#### 查询错误
```typescript
// 错误类型
- "未找到绑定的游戏账号" - 用户未绑定游戏账号
- "会话不存在，请先登录游戏" - 管理员未启动游戏会话
```

### 2. 用户提示

建议在以下场景添加友好提示：

#### 未绑定游戏账号时
```
提示：您还未绑定游戏账号，无法查询战绩和圈子。
请前往"个人中心"绑定游戏账号。
```

#### 游戏账号已被绑定时
```
提示：该游戏账号已被其他用户绑定。
如果这是您的账号，请联系管理员处理。
```

---

## 📊 数据流向图

### 游戏账号入圈流程

```
游戏内申请
    ↓
[内存队列]
    ↓
管理员审批 ← 前端调用 API
    ↓
后端处理：
├─ 查找/创建游戏账号 (user_id = NULL)
├─ 确保管理员圈子
├─ 创建 game_account_group 记录
└─ 调用游戏 API
    ↓
[game_account_group 表]
    ↓
用户绑定游戏账号 ← 前端调用 API
    ↓
更新 game_account.user_id
    ↓
用户可查询圈子和战绩
```

### 用户查询流程

```
前端调用 API
    ↓
后端反向查询：
├─ 根据 user_id 查询 game_account
├─ 根据 game_account_id 查询 game_account_group
└─ 返回圈子/战绩列表
    ↓
前端显示结果
```

---

## ✅ 测试检查清单

### 1. 游戏账号绑定
- [ ] 绑定新的游戏账号（账号不存在）
- [ ] 绑定已存在的游戏账号（未绑定用户）
- [ ] 尝试绑定已被其他用户绑定的账号（应显示错误）
- [ ] 解绑游戏账号
- [ ] 重新绑定游戏账号

### 2. 圈子查询
- [ ] 未绑定游戏账号时查询圈子（应提示绑定）
- [ ] 绑定游戏账号后查询圈子
- [ ] 游戏账号在多个店铺的圈子中（应显示所有圈子）

### 3. 战绩查询
- [ ] 未绑定游戏账号时查询战绩（应提示绑定）
- [ ] 绑定游戏账号后查询战绩
- [ ] 筛选不同店铺/圈子的战绩

### 4. 申请审批
- [ ] 查看游戏内申请列表
- [ ] 通过申请（检查是否自动创建游戏账号和圈子关系）
- [ ] 拒绝申请
- [ ] 通过申请后，玩家注册并绑定游戏账号

---

## 📚 相关文档

- [后端改造总结](../../battle-tiles/doc/BACKEND_GAME_ACCOUNT_CHANGES.md)
- [游戏账号系统重新设计方案 V2](../../battle-tiles/doc/GAME_ACCOUNT_REDESIGN_V2.md)
- [游戏账号简明指南](../../battle-tiles/doc/GAME_ACCOUNT_SIMPLE_GUIDE.md)

---

## 🎉 总结

前端改动非常小，主要是：
1. ✅ 优化游戏账号绑定的错误提示
2. ✅ 添加二次确认对话框
3. ✅ 使用气泡提示显示成功消息

其他所有逻辑都在后端透明处理，前端无需修改！🎊

