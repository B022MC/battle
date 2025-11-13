# 前端集成指南 - 圈子系统

## 概述

新的圈子系统已经完成后端开发和前端基础集成。本文档说明如何使用新的 API 和前端组件。

## 后端 API 端点

### 成员管理 API

#### 1. 列出所有用户
```
POST /api/members/list
```

**请求参数：**
```typescript
{
  page?: number;      // 页码，默认 1
  size?: number;      // 每页数量，默认 20
  keyword?: string;   // 搜索关键词（用户名或昵称）
}
```

**响应：**
```typescript
{
  items: BasicUser[];  // 用户列表
  total: number;       // 总数
}
```

#### 2. 获取用户信息
```
POST /api/members/get
```

**请求参数：**
```typescript
{
  user_id: number;
}
```

**响应：**
```typescript
BasicUser
```

#### 3. 列出店铺管理员
```
POST /api/members/shop-admins
```

**请求参数：**
```typescript
{
  house_gid: number;
}
```

**响应：**
```typescript
BasicUser[]
```

### 圈子管理 API

#### 1. 创建圈子
```
POST /api/groups/create
```

**请求参数：**
```typescript
{
  house_gid: number;
  group_name: string;
  description?: string;
}
```

**响应：**
```typescript
ShopGroup
```

#### 2. 获取我的圈子
```
POST /api/groups/my
```

**请求参数：**
```typescript
{
  house_gid: number;
}
```

**响应：**
```typescript
ShopGroup
```

#### 3. 列出店铺下的所有圈子（超级管理员）
```
POST /api/groups/list
```

**请求参数：**
```typescript
{
  house_gid: number;
}
```

**响应：**
```typescript
ShopGroup[]
```

#### 4. 添加成员到圈子
```
POST /api/groups/members/add
```

**请求参数：**
```typescript
{
  group_id: number;
  user_ids: number[];
}
```

**响应：**
```typescript
{
  success: boolean;
}
```

#### 5. 从圈子移除成员
```
POST /api/groups/members/remove
```

**请求参数：**
```typescript
{
  group_id: number;
  user_id: number;
}
```

**响应：**
```typescript
{
  success: boolean;
}
```

#### 6. 列出圈子成员
```
POST /api/groups/members/list
```

**请求参数：**
```typescript
{
  group_id: number;
  page?: number;
  size?: number;
}
```

**响应：**
```typescript
{
  items: BasicUser[];
  total: number;
}
```

#### 7. 列出我加入的所有圈子
```
POST /api/groups/my/list
```

**请求参数：**
```typescript
{}
```

**响应：**
```typescript
ShopGroup[]
```

## 前端使用

### 1. 导入服务

```typescript
import { listAllUsers, getUser, listShopAdmins } from '@/services/members';
import { 
  createGroup,
  getMyGroup, 
  listGroupsByHouse, 
  addMembersToGroup, 
  removeMemberFromGroup, 
  listGroupMembers,
  listMyGroups
} from '@/services/shops/groups';
```

### 2. 使用示例

#### 列出所有用户
```typescript
const { data, loading, run } = useRequest(listAllUsers, { manual: true });

// 调用
await run({ page: 1, size: 20, keyword: '张三' });
```

#### 获取我的圈子
```typescript
const { data: myGroup, run: runGetMyGroup } = useRequest(getMyGroup, { manual: true });

// 调用
const group = await runGetMyGroup({ house_gid: 60870 });
```

#### 添加成员到圈子
```typescript
const { run: runAddMembers } = useRequest(addMembersToGroup, { manual: true });

// 调用
await runAddMembers({ 
  group_id: 1, 
  user_ids: [10, 20, 30] 
});
```

## 数据类型定义

### BasicUser
```typescript
type BasicUser = {
  id: number;
  username: string;
  nick_name: string;
  phone?: string;
  email?: string;
  created_at: string;
};
```

### ShopGroup
```typescript
type ShopGroup = {
  id: number;
  house_gid: number;
  group_name: string;
  admin_user_id: number;
  description: string;
  member_count?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};
```

## 权限说明

### 超级管理员
- 可以查看所有用户
- 可以查看所有店铺的所有圈子
- 可以指定用户为店铺管理员

### 店铺管理员
- 可以查看所有用户
- 可以创建自己的圈子（每个店铺管理员一个圈子）
- 可以管理自己圈子的成员（添加/移除）
- 可以查看自己圈子的成员列表

### 普通用户
- 可以查看自己加入的圈子
- 不能管理圈子

## 业务流程

### 1. 店铺管理员创建圈子
1. 用户成为店铺管理员（由超级管理员指定）
2. 店铺管理员首次访问圈子功能时，系统自动创建圈子
3. 圈子名称默认为 "店铺{house_gid}的圈子"

### 2. 添加成员到圈子
1. 店铺管理员查看所有用户列表
2. 选择要添加的用户
3. 点击"添加到圈子"按钮
4. 系统将用户添加到圈子

### 3. 移除圈子成员
1. 店铺管理员查看圈子成员列表
2. 点击成员旁边的"移除"按钮
3. 系统将用户从圈子移除

## 注意事项

1. **一个用户可以加入多个圈子**
2. **每个店铺管理员只能创建一个圈子**
3. **圈子与店铺绑定**（house_gid）
4. **成员列表不再区分游戏端和平台端**，统一为系统用户
5. **所有 API 都需要 JWT 认证**

## 测试建议

1. 测试超级管理员查看所有用户
2. 测试店铺管理员创建圈子
3. 测试添加成员到圈子
4. 测试从圈子移除成员
5. 测试分页和搜索功能
6. 测试权限控制（普通用户不能管理圈子）

