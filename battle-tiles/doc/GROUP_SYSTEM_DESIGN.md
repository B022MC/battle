# 圈子系统设计文档

## 📋 需求概述

### 当前问题
1. 成员列表区分"游戏成员"和"平台成员"，过于复杂
2. 缺少"圈子"（Group）的明确概念
3. 成员管理逻辑不清晰

### 新的业务逻辑

#### 1. **用户角色**
- **超级管理员**：
  - 可以看到所有系统用户
  - 可以指定用户成为某个店铺的管理员
  
- **店铺管理员**：
  - 可以看到所有系统用户
  - 可以创建自己的圈子
  - 可以将用户添加到自己的圈子

#### 2. **圈子（Group）概念**
- 每个店铺管理员 = 一个圈子
- 店铺管理员可以将成员加入自己的圈子
- 一个用户可以加入多个圈子（不同店铺或同一店铺的不同圈子）

#### 3. **成员列表**
- 统一的成员列表 = 系统中的所有用户（`basic_user` 表）
- 不再区分游戏成员和平台成员
- 超级管理员和店铺管理员都能看到所有用户

---

## 🗄️ 数据库设计

### 新增表

#### 1. `game_shop_group` - 店铺圈子表
```sql
CREATE TABLE game_shop_group (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,           -- 店铺GID
    group_name VARCHAR(64) NOT NULL,      -- 圈子名称
    admin_user_id INTEGER NOT NULL,       -- 圈主用户ID（店铺管理员）
    description TEXT DEFAULT '',          -- 圈子描述
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

**说明**：
- 一个店铺管理员在一个店铺下只能有一个圈子
- `admin_user_id` 关联 `basic_user.id`
- `house_gid` 关联店铺

#### 2. `game_shop_group_member` - 圈子成员关系表
```sql
CREATE TABLE game_shop_group_member (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,            -- 圈子ID
    user_id INTEGER NOT NULL,             -- 用户ID
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

**说明**：
- 一个用户可以加入多个圈子
- `group_id` 关联 `game_shop_group.id`
- `user_id` 关联 `basic_user.id`

### 保留的表

#### `game_shop_admin` - 店铺管理员表
- 保持不变
- 用于记录哪些用户是哪个店铺的管理员

#### `game_member` - 游戏成员表
- **可选**：保留用于存储游戏内的玩家数据（余额、信用等）
- 或者：逐步废弃，将余额等信息迁移到新的钱包系统

---

## 🔄 业务流程

### 1. 超级管理员指定店铺管理员

```
超级管理员 → 选择用户 → 指定为某店铺的管理员
                ↓
        插入 game_shop_admin 记录
                ↓
        自动创建该管理员的圈子（game_shop_group）
```

### 2. 店铺管理员创建圈子

```
店铺管理员登录 → 自动检查是否有圈子
                ↓
        如果没有 → 创建默认圈子
                ↓
        圈子名称 = 管理员用户名 + "的圈子"
```

### 3. 店铺管理员添加成员到圈子

```
店铺管理员 → 查看所有用户列表
          ↓
      选择用户 → 添加到自己的圈子
          ↓
  插入 game_shop_group_member 记录
```

### 4. 查看圈子成员

```
店铺管理员 → 查看自己的圈子
          ↓
      显示圈子内的所有成员
```

---

## 📡 API 设计

### 1. 成员管理 API

#### GET `/api/members/list` - 查看所有用户
**权限**：超级管理员、店铺管理员

**请求参数**：
```json
{
  "page": 1,
  "size": 20,
  "keyword": "搜索关键词（可选）"
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "username": "user001",
        "nickname": "张三",
        "phone": "13800138000",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

### 2. 圈子管理 API

#### POST `/api/groups/create` - 创建圈子
**权限**：店铺管理员

**请求参数**：
```json
{
  "house_gid": 60870,
  "group_name": "VIP圈子",
  "description": "高级会员圈子"
}
```

#### GET `/api/groups/my` - 查看我的圈子
**权限**：店铺管理员

**响应**：
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "house_gid": 60870,
    "group_name": "VIP圈子",
    "admin_user_id": 1,
    "member_count": 50,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### POST `/api/groups/members/add` - 添加成员到圈子
**权限**：店铺管理员（只能添加到自己的圈子）

**请求参数**：
```json
{
  "group_id": 1,
  "user_ids": [2, 3, 4]
}
```

#### POST `/api/groups/members/remove` - 从圈子移除成员
**权限**：店铺管理员

**请求参数**：
```json
{
  "group_id": 1,
  "user_id": 2
}
```

#### GET `/api/groups/members/list` - 查看圈子成员列表
**权限**：店铺管理员

**请求参数**：
```json
{
  "group_id": 1,
  "page": 1,
  "size": 20
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 2,
        "username": "user002",
        "nickname": "李四",
        "joined_at": "2025-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 20
  }
}
```

### 3. 店铺管理员管理 API

#### POST `/api/admin/shop-admins/assign` - 指定店铺管理员
**权限**：超级管理员

**请求参数**：
```json
{
  "house_gid": 60870,
  "user_id": 2,
  "role": "admin"
}
```

#### GET `/api/admin/shop-admins/list` - 查看店铺管理员列表
**权限**：超级管理员

**请求参数**：
```json
{
  "house_gid": 60870
}
```

---

## 🔧 实现步骤

### 阶段 1：数据库迁移
1. ✅ 创建 `game_shop_group` 表
2. ✅ 创建 `game_shop_group_member` 表
3. ✅ 创建对应的 Go 模型文件

### 阶段 2：Repository 层
1. 创建 `ShopGroupRepo` 接口和实现
2. 创建 `ShopGroupMemberRepo` 接口和实现
3. 扩展 `UserRepo` 添加成员列表查询方法

### 阶段 3：Use Case 层
1. 创建 `GroupUseCase` - 圈子管理业务逻辑
2. 创建 `MemberUseCase` - 成员管理业务逻辑
3. 扩展 `ShopAdminUseCase` - 店铺管理员管理

### 阶段 4：Service 层
1. 创建 `GroupService` - 圈子管理 HTTP API
2. 创建 `MemberService` - 成员管理 HTTP API
3. 注册路由

### 阶段 5：前端集成
1. 修改成员列表页面，统一显示所有用户
2. 添加圈子管理页面
3. 添加"添加到圈子"功能

---

## 🎯 下一步行动

你希望我：
1. **立即开始实现**：创建 Repository、Use Case、Service 层代码
2. **先执行数据库迁移**：运行 SQL 脚本创建新表
3. **先讨论细节**：确认业务逻辑是否符合你的需求

请告诉我你的选择，我会继续推进！

