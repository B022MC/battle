# RBAC 权限管理系统完整实施指南

## 📋 概述

本文档提供完整的 RBAC（基于角色的访问控制）系统实施指南，包括后端 API、前端页面、数据库配置等所有方面。

## 🎯 系统特性

### 核心功能

1. **权限管理**
   - 创建、编辑、删除权限
   - 按分类管理（stats/fund/shop/game/system）
   - 细粒度到按钮级别

2. **角色管理**
   - 创建、编辑、删除角色
   - 为角色分配权限
   - 为角色分配菜单
   - 启用/禁用角色

3. **菜单管理**
   - 菜单与权限关联
   - 按钮级权限配置
   - 动态菜单显示

4. **用户管理**
   - 用户角色分配
   - 多角色支持
   - 权限继承

## 🗄️ 数据库部署

### 步骤 1: 执行初始化脚本

```bash
cd battle-tiles/doc/rbac

# 首次部署（全新系统）
psql -U B022MC -d your_database -f 00_init_data.sql

# 系统升级（已有旧数据）
psql -U B022MC -d your_database -f 01_update_permissions.sql
```

### 步骤 2: 验证数据

```sql
-- 检查角色数量
SELECT COUNT(*) FROM basic_role WHERE is_deleted = false;
-- 预期: 3 (超级管理员、店铺管理员、普通用户)

-- 检查菜单数量
SELECT COUNT(*) FROM basic_menu WHERE is_del = 0;
-- 预期: 18 (6个一级 + 12个二级)

-- 检查权限数量
SELECT COUNT(*) FROM basic_permission WHERE is_deleted = false;
-- 预期: 43

-- 检查超级管理员权限
SELECT COUNT(*) FROM basic_role_permission_rel WHERE role_id = 1;
-- 预期: 43 (所有权限)
```

### 步骤 3: 为测试用户分配角色

```sql
-- 分配超级管理员角色给用户 ID=1
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (1, 1)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 分配店铺管理员角色给用户 ID=2
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (2, 2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 分配普通用户角色给用户 ID=3
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (3, 3)
ON CONFLICT (user_id, role_id) DO NOTHING;
```

## 🔧 后端配置

### 已实现的 API 接口

#### 权限管理 API

```go
// battle-tiles/internal/service/basic/basic_permission.go

GET  /basic/permission/list          // 查询权限列表
GET  /basic/permission/listAll       // 查询所有权限
POST /basic/permission/create        // 创建权限
POST /basic/permission/update        // 更新权限
POST /basic/permission/delete        // 删除权限
GET  /basic/permission/role/permissions  // 查询角色权限
POST /basic/permission/role/assign   // 为角色分配权限
POST /basic/permission/role/remove   // 从角色移除权限
```

#### 角色管理 API

```go
// battle-tiles/internal/service/basic/basic_role.go

GET  /basic/role/list               // 查询角色列表（分页）
GET  /basic/role/getOne             // 查询单个角色
GET  /basic/role/all                // 查询所有角色
POST /basic/role/create             // 创建角色
POST /basic/role/update             // 更新角色
POST /basic/role/delete             // 删除角色
GET  /basic/role/menus              // 查询角色菜单
POST /basic/role/menus/assign       // 为角色分配菜单
```

### 权限中间件

所有管理接口已添加权限验证：

```go
// 权限管理需要 permission:view 权限
r.GET("/list", middleware.RequireAnyPerm("permission:view"), s.ListPermissions)

// 创建角色需要 role:create 权限
r.POST("/create", middleware.RequirePerm("role:create"), s.Create)

// 分配权限需要 permission:assign 权限
r.POST("/role/assign", middleware.RequirePerm("permission:assign"), s.AssignPermissionsToRole)
```

### 启动后端服务

```bash
cd battle-tiles

# 编译
go build -o bin/server ./cmd/go-kgin-platform

# 运行
./bin/server

# 或使用 make
make run
```

## 💻 前端配置

### 已创建的页面和组件

#### 1. 权限管理页面

**位置**: `battle-reusables/app/(shop)/permissions.tsx`

**组件**:
- `PermissionsView` - 主视图
- `PermissionList` - 权限列表
- `PermissionForm` - 创建/编辑表单

**功能**:
- ✅ 按分类筛选权限
- ✅ 创建新权限
- ✅ 编辑现有权限
- ✅ 删除权限
- ✅ 按钮级权限控制

#### 2. 角色管理页面

**位置**: `battle-reusables/app/(shop)/roles.tsx`

**组件**:
- `RolesView` - 主视图
- `RoleList` - 角色列表
- `RoleForm` - 创建/编辑表单
- `AssignPermissionsModal` - 分配权限弹窗
- `AssignMenusModal` - 分配菜单弹窗

**功能**:
- ✅ 查看所有角色
- ✅ 创建新角色
- ✅ 编辑角色信息
- ✅ 删除自定义角色
- ✅ 为角色分配权限
- ✅ 为角色分配菜单
- ✅ 启用/禁用角色

#### 3. Service 层

**位置**: `battle-reusables/services/basic/`

- `permission.ts` - 权限相关 API
- `role.ts` - 角色相关 API
- `menu.ts` - 菜单相关 API

### 路由配置

权限管理和角色管理页面已添加到店铺菜单下：

```
/(tabs)/shop
  ├── /(shop)/permissions  ← 权限管理
  └── /(shop)/roles        ← 角色管理
```

### 权限控制

所有页面和按钮都已添加权限控制：

```tsx
// 页面级权限（route guard）
<RouteGuard anyOf={['permission:view']}>
  <PermissionsView />
</RouteGuard>

// 按钮级权限
<PermissionGate anyOf={['permission:create']}>
  <Button>创建权限</Button>
</PermissionGate>
```

### 启动前端应用

```bash
cd battle-reusables

# 安装依赖
npm install

# 启动开发服务器
npm start

# 或构建生产版本
npm run build
```

## 🧪 测试指南

### 1. 后端 API 测试

#### 测试权限管理 API

```bash
# 获取 JWT Token
TOKEN="your_jwt_token"

# 查询所有权限
curl -X GET "http://localhost:8080/api/basic/permission/listAll" \
  -H "Authorization: Bearer $TOKEN"

# 创建新权限
curl -X POST "http://localhost:8080/api/basic/permission/create" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test:permission",
    "name": "测试权限",
    "category": "system",
    "description": "测试用权限"
  }'

# 查询角色权限
curl -X GET "http://localhost:8080/api/basic/permission/role/permissions?role_id=1" \
  -H "Authorization: Bearer $TOKEN"

# 为角色分配权限
curl -X POST "http://localhost:8080/api/basic/permission/role/assign" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 2,
    "permission_ids": [1, 2, 3, 4, 5]
  }'
```

#### 测试角色管理 API

```bash
# 查询所有角色
curl -X GET "http://localhost:8080/api/basic/role/all" \
  -H "Authorization: Bearer $TOKEN"

# 创建新角色
curl -X POST "http://localhost:8080/api/basic/role/create" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test_role",
    "name": "测试角色",
    "remark": "这是一个测试角色"
  }'

# 查询角色菜单
curl -X GET "http://localhost:8080/api/basic/role/menus?role_id=2" \
  -H "Authorization: Bearer $TOKEN"

# 为角色分配菜单
curl -X POST "http://localhost:8080/api/basic/role/menus/assign" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 2,
    "menu_ids": [1, 2, 3, 4, 5, 6]
  }'
```

### 2. 前端功能测试

#### 测试场景 1: 超级管理员

1. **登录**
   - 使用超级管理员账号登录
   - 验证可以看到所有菜单项

2. **权限管理**
   - 进入"店铺" → "权限管理"
   - ✅ 查看所有权限列表
   - ✅ 按分类筛选权限
   - ✅ 创建新权限
   - ✅ 编辑现有权限
   - ✅ 删除权限

3. **角色管理**
   - 进入"店铺" → "角色管理"
   - ✅ 查看所有角色
   - ✅ 创建新角色
   - ✅ 编辑角色信息
   - ✅ 为角色分配权限
   - ✅ 为角色分配菜单
   - ✅ 删除自定义角色
   - ❌ 不能删除系统预定义角色（1,2,3）

#### 测试场景 2: 店铺管理员

1. **登录**
   - 使用店铺管理员账号登录
   - 验证可以看到管理相关菜单

2. **权限验证**
   - ❌ 不应该看到"权限管理"菜单
   - ❌ 不应该看到"角色管理"菜单
   - ✅ 可以看到其他店铺管理功能

#### 测试场景 3: 普通用户

1. **登录**
   - 使用普通用户账号登录
   - 验证只能看到基础菜单

2. **权限验证**
   - ❌ 不应该看到任何管理功能
   - ✅ 可以查看自己的数据
   - ✅ 可以绑定游戏账号

### 3. 按钮级权限测试

访问各个页面，验证按钮是否根据权限正确显示/隐藏：

#### 权限管理页面按钮
- `创建权限` - 需要 `permission:create`
- `编辑` - 需要 `permission:update`
- `删除` - 需要 `permission:delete`

#### 角色管理页面按钮
- `创建角色` - 需要 `role:create`
- `编辑` - 需要 `role:update`
- `删除` - 需要 `role:delete`
- `分配权限` - 需要 `permission:assign`
- `分配菜单` - 需要 `role:update`

#### 其他页面按钮
- 资金页面 - `上分`、`下分`、`设置阈值`
- 桌台页面 - `详情`、`检查`、`解散`
- 成员页面 - `踢出成员`、`成员下线`、`拉入圈子`

### 4. 数据库验证

```sql
-- 验证角色权限分配
SELECT r.name, COUNT(*) as perm_count
FROM basic_role r
JOIN basic_role_permission_rel rpr ON rpr.role_id = r.id
WHERE r.is_deleted = false
GROUP BY r.id, r.name
ORDER BY r.id;

-- 验证角色菜单分配
SELECT r.name, COUNT(*) as menu_count
FROM basic_role r
JOIN basic_role_menu_rel rmr ON rmr.role_id = r.id
WHERE r.is_deleted = false
GROUP BY r.id, r.name
ORDER BY r.id;

-- 验证用户角色分配
SELECT u.id as user_id, r.name as role_name
FROM basic_user_role_rel urr
JOIN basic_role r ON r.id = urr.role_id
WHERE r.is_deleted = false
ORDER BY u.id;
```

## 📝 使用文档

### 创建新权限

1. 进入"店铺" → "权限管理"
2. 点击"创建权限"按钮
3. 填写权限信息：
   - **权限编码**: `module:feature:action` 格式
   - **权限名称**: 中文描述
   - **权限分类**: 选择合适的分类
   - **权限描述**: 详细说明
4. 点击"创建"

### 创建新角色

1. 进入"店铺" → "角色管理"
2. 点击"创建角色"按钮
3. 填写角色信息：
   - **角色编码**: 英文字母和下划线
   - **角色名称**: 中文名称
   - **角色备注**: 用途说明
4. 点击"创建"
5. 为角色分配权限和菜单

### 为角色分配权限

1. 在角色列表中找到目标角色
2. 点击"分配权限"按钮
3. 勾选需要的权限：
   - 可以按分类全选
   - 可以单独勾选
4. 点击"保存"

### 为角色分配菜单

1. 在角色列表中找到目标角色
2. 点击"分配菜单"按钮
3. 勾选需要的菜单：
   - 选中一级菜单会自动选中所有子菜单
   - 也可以单独选择子菜单
4. 点击"保存"

### 为用户分配角色

```sql
-- 通过 SQL 分配（暂无前端界面）
INSERT INTO basic_user_role_rel (user_id, role_id) 
VALUES (用户ID, 角色ID)
ON CONFLICT (user_id, role_id) DO NOTHING;
```

## 🚨 常见问题

### Q1: 权限修改后不生效？

**A**: RBAC 系统使用 Redis 缓存（TTL 10分钟），修改权限后需要：
- 等待缓存过期（10分钟）
- 或手动清理 Redis 缓存
- 或重新登录

```bash
# 清理用户权限缓存
redis-cli DEL "rbac:perms:用户ID"
```

### Q2: 不能删除系统角色？

**A**: 为了系统安全，以下角色不能删除：
- 超级管理员 (id=1)
- 店铺管理员 (id=2)
- 普通用户 (id=3)

如需修改，请直接更新数据库。

### Q3: 创建权限时编码重复？

**A**: 权限编码必须唯一。建议命名规范：
- `module:feature:action`
- 例如：`shop:member:view`、`fund:deposit`

### Q4: 按钮权限配置在哪里？

**A**: 按钮权限在 `basic_menu_button` 表中配置：

```sql
INSERT INTO basic_menu_button (menu_id, button_code, button_name, permission_codes)
VALUES (菜单ID, '按钮编码', '按钮名称', '权限1,权限2');
```

### Q5: 如何调试权限问题？

**A**: 使用以下 SQL 查询用户的完整权限：

```sql
-- 查询用户所有权限
SELECT DISTINCT p.code, p.name 
FROM basic_user_role_rel urr
JOIN basic_role_permission_rel rpr ON rpr.role_id = urr.role_id
JOIN basic_permission p ON p.id = rpr.permission_id
WHERE urr.user_id = ? AND p.is_deleted = false;
```

## 📊 性能优化

1. **索引优化**
   - 所有关联表已添加适当索引
   - 权限表使用唯一索引

2. **缓存策略**
   - 用户权限缓存 10 分钟
   - 菜单数据登录时一次性获取

3. **查询优化**
   - 使用 JOIN 减少查询次数
   - 分页查询避免大数据量

## 🔐 安全建议

1. **最小权限原则**: 只分配必要的权限
2. **定期审计**: 检查用户权限分配情况
3. **密码安全**: 超级管理员账号使用强密码
4. **操作日志**: 记录权限变更操作
5. **备份策略**: 定期备份权限数据

## 📚 相关文档

- [RBAC 系统总结](./RBAC_IMPLEMENTATION_SUMMARY.md)
- [数据库脚本说明](./README.md)
- [API 文档](../../docs/swagger.yaml)

## ✅ 部署检查清单

- [ ] 数据库脚本执行完成
- [ ] 后端服务启动正常
- [ ] 前端应用启动正常
- [ ] 超级管理员账号可以登录
- [ ] 权限管理页面可以访问
- [ ] 角色管理页面可以访问
- [ ] 按钮权限正常显示/隐藏
- [ ] 页面权限路由守卫生效
- [ ] 测试用户权限验证通过

---

**文档版本**: 1.0  
**创建日期**: 2025-11-18  
**维护者**: Development Team

