# RBAC 权限系统文档

本目录包含完整的 RBAC（基于角色的访问控制）系统的数据库脚本和文档。

## 📁 文件说明

### SQL 脚本

1. **00_init_data.sql** - 完整初始化脚本
   - 创建所有 RBAC 相关表结构
   - 插入初始角色、菜单、权限数据
   - 配置角色权限关联
   - 配置菜单按钮权限
   - 适用于首次部署或完全重置

2. **01_update_permissions.sql** - 权限更新脚本
   - 更新现有角色信息
   - 添加新的菜单（权限管理、角色管理、菜单管理）
   - 添加新的权限
   - 更新角色权限关联
   - 清理无效数据
   - 适用于系统升级

3. **02_add_menu_management.sql** - 菜单管理快速修复脚本
   - 添加菜单管理菜单项
   - 为超级管理员分配菜单管理权限
   - 配置菜单管理按钮
   - 适用于修复缺少菜单管理功能的问题

3. **base_role.sql** - 角色表结构和初始数据
   - 超级管理员（id=1）
   - 店铺管理员（id=2）
   - 普通用户（id=3）

4. **basic_menu.sql** - 菜单表结构和数据
   - 一级菜单（底部Tab）
   - 二级菜单（店铺子页面）
   - 包含权限管理和角色管理菜单

5. **basic_role_menu_rel.sql** - 角色菜单关联
   - 定义每个角色可访问的菜单

6. **basic_permission.sql** - 权限表和数据
   - 完整的权限定义
   - 角色权限关联
   - 菜单按钮配置

## 🚀 部署指南

### 首次部署

```bash
# 方式一：执行完整初始化脚本
psql -U B022MC -d your_database -f 00_init_data.sql

# 方式二：分步执行
psql -U B022MC -d your_database -f base_role.sql
psql -U B022MC -d your_database -f basic_menu.sql
psql -U B022MC -d your_database -f basic_role_menu_rel.sql
psql -U B022MC -d your_database -f basic_permission.sql
```

### 系统升级

如果已有旧数据，需要升级到新的 RBAC 系统：

```bash
# 执行更新脚本
psql -U B022MC -d your_database -f 01_update_permissions.sql
```

### 快速修复（缺少菜单管理）

如果发现系统中缺少菜单管理功能：

```bash
# 执行快速修复脚本
psql -U B022MC -d your_database -f 02_add_menu_management.sql
```

详细说明请参考 [MENU_MANAGEMENT_FIX.md](./MENU_MANAGEMENT_FIX.md)

## 📊 数据结构

### 核心表

1. **basic_role** - 角色表
   - 存储系统角色信息
   - 支持角色继承（parent_id）
   - 软删除支持

2. **basic_menu** - 菜单表
   - 存储前端菜单配置
   - 支持二级菜单结构
   - 包含权限标识（auths 字段）

3. **basic_permission** - 权限表
   - 存储细粒度权限定义
   - 按分类组织（stats/fund/shop/game/system）
   - 支持权限描述

4. **basic_role_menu_rel** - 角色菜单关联表
   - 定义角色可访问的菜单
   - 联合主键（role_id, menu_id）

5. **basic_role_permission_rel** - 角色权限关联表
   - 定义角色拥有的权限
   - 联合主键（role_id, permission_id）

6. **basic_menu_button** - 菜单按钮表
   - 定义菜单内的按钮级权限
   - 支持多权限OR逻辑（permission_codes）

7. **basic_user_role_rel** - 用户角色关联表
   - 定义用户的角色分配
   - 支持一个用户拥有多个角色

## 🎯 权限体系

### 权限分类

- **stats** - 统计数据权限
- **fund** - 资金管理权限
- **shop** - 店铺管理权限
- **game** - 游戏控制权限
- **system** - 系统管理权限

### 权限命名规范

格式：`分类:功能:操作`

示例：
- `shop:admin:view` - 查看店铺管理员
- `shop:admin:assign` - 分配店铺管理员
- `fund:deposit` - 上分操作
- `permission:create` - 创建权限

### 预定义角色

#### 1. 超级管理员 (super_admin, id=1)
- 拥有系统所有权限
- 可以管理权限、角色、菜单
- 不能被删除

#### 2. 店铺管理员 (shop_admin, id=2)
- 拥有店铺管理相关权限
- 可以管理成员、资金、桌台等
- 不能管理系统级设置
- 不能被删除

#### 3. 普通用户 (user, id=3)
- 仅拥有基础查看权限
- 可以查看自己的数据
- 不能管理他人数据
- 不能被删除

## 🔧 使用示例

### 创建新角色

```sql
-- 1. 创建角色
INSERT INTO basic_role (code, name, remark, enable) 
VALUES ('custom_role', '自定义角色', '描述信息', true);

-- 2. 为角色分配权限
INSERT INTO basic_role_permission_rel (role_id, permission_id)
SELECT 4, id FROM basic_permission 
WHERE code IN ('shop:member:view', 'fund:wallet:view');

-- 3. 为角色分配菜单
INSERT INTO basic_role_menu_rel (role_id, menu_id)
VALUES (4, 3), (4, 4), (4, 5), (4, 6);
```

### 为用户分配角色

```sql
-- 将用户 ID=100 设置为店铺管理员
INSERT INTO basic_user_role_rel (user_id, role_id) 
VALUES (100, 2)
ON CONFLICT (user_id, role_id) DO NOTHING;
```

### 查询用户权限

```sql
-- 查询用户的所有权限
SELECT DISTINCT p.code, p.name, p.category
FROM basic_user_role_rel urr
JOIN basic_role_permission_rel rpr ON rpr.role_id = urr.role_id
JOIN basic_permission p ON p.id = rpr.permission_id
WHERE urr.user_id = 100 AND p.is_deleted = false
ORDER BY p.category, p.code;
```

### 查询角色的菜单

```sql
-- 查询角色可访问的菜单
SELECT m.*
FROM basic_role_menu_rel rmr
JOIN basic_menu m ON m.id = rmr.menu_id
WHERE rmr.role_id = 2 AND m.is_del = 0
ORDER BY m.menu_type, m.id;
```

## 📝 前端集成

### 路由配置

RBAC管理页面已添加到系统：

- `/(shop)/permissions` - 权限管理页面
- `/(shop)/roles` - 角色管理页面
- `/(shop)/menus` - 菜单管理页面

### 权限控制组件

```tsx
// 页面级权限控制
<RouteGuard anyOf={['permission:view']}>
  <PermissionsView />
</RouteGuard>

// 按钮级权限控制
<PermissionGate anyOf={['permission:create']}>
  <Button>创建权限</Button>
</PermissionGate>
```

### API 接口

#### 权限管理
- GET `/basic/permission/list` - 查询权限列表
- GET `/basic/permission/listAll` - 查询所有权限
- POST `/basic/permission/create` - 创建权限
- POST `/basic/permission/update` - 更新权限
- POST `/basic/permission/delete` - 删除权限
- GET `/basic/permission/role/permissions` - 查询角色权限
- POST `/basic/permission/role/assign` - 为角色分配权限
- POST `/basic/permission/role/remove` - 从角色移除权限

#### 角色管理
- GET `/basic/role/list` - 查询角色列表
- GET `/basic/role/getOne` - 查询单个角色
- GET `/basic/role/all` - 查询所有角色
- POST `/basic/role/create` - 创建角色
- POST `/basic/role/update` - 更新角色
- POST `/basic/role/delete` - 删除角色
- GET `/basic/role/menus` - 查询角色菜单
- POST `/basic/role/menus/assign` - 为角色分配菜单

## 🔒 安全注意事项

1. **不要删除系统预定义角色**（id=1,2,3）
2. **超级管理员账号需要严格保护**
3. **权限变更后需要清理 Redis 缓存**
4. **定期审计权限分配情况**
5. **最小权限原则**：只分配必要的权限

## 📊 数据统计

执行初始化脚本后的数据量：

- 角色：3 个（超级管理员、店铺管理员、普通用户）
- 菜单：19 个（6个一级菜单 + 13个二级菜单，包含菜单管理）
- 权限：43 个（覆盖所有功能模块）
- 菜单按钮：26 个（细粒度按钮控制）

## 🔄 维护建议

1. **定期备份**：在修改权限数据前进行备份
2. **版本控制**：所有 SQL 脚本纳入版本控制
3. **文档更新**：添加新权限时同步更新文档
4. **测试验证**：权限变更后进行完整测试
5. **日志记录**：记录权限变更操作

## 📞 问题反馈

如有问题或建议，请联系开发团队。

---

**最后更新**: 2025-11-18  
**版本**: 2.0  
**维护者**: Development Team
