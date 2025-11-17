# RBAC 权限系统完整实现总结

本文档总结了整个 RBAC（基于角色的访问控制）系统的实现，包括数据库设计、后端 API、前端组件等所有细节。

## 一、数据库改造

### 1.1 新增表结构

#### basic_permission 权限表
存储所有细粒度权限定义：
- `id`: 权限ID（主键）
- `code`: 权限编码（唯一）
- `name`: 权限名称
- `category`: 权限分类（stats/fund/shop/game/system）
- `description`: 权限描述
- `is_deleted`: 软删除标记

#### basic_role_permission_rel 角色权限关联表
定义角色和权限的多对多关系：
- `role_id`: 角色ID
- `permission_id`: 权限ID

#### basic_menu_button 菜单按钮表
定义UI级别的按钮权限配置：
- `id`: 按钮ID
- `menu_id`: 所属菜单ID
- `button_code`: 按钮编码
- `button_name`: 按钮名称
- `permission_codes`: 所需权限码（逗号分隔）

### 1.2 权限码定义

#### 统计相关
- `stats:view`: 查看统计数据

#### 资金相关
- `fund:wallet:view`: 查看钱包
- `fund:ledger:view`: 查看账本
- `fund:deposit`: 上分
- `fund:withdraw`: 下分
- `fund:force_withdraw`: 强制下分
- `fund:limit:update`: 更新额度

#### 店铺相关
- `shop:table:view`: 查看桌台
- `shop:table:detail`: 查看桌台详情
- `shop:table:check`: 检查桌台
- `shop:table:dismiss`: 解散桌台
- `shop:member:view`: 查看成员
- `shop:member:kick`: 踢出成员
- `shop:member:logout`: 成员下线
- `shop:member:update`: 更新成员
- `shop:admin:view`: 查看管理员
- `shop:admin:assign`: 分配管理员
- `shop:admin:revoke`: 撤销管理员
- `shop:apply:view`: 查看申请
- `shop:apply:approve`: 批准申请
- `shop:apply:reject`: 拒绝申请
- `shop:fees:view`: 查看费用
- `shop:fees:update`: 更新费用
- `shop:group:view`: 查看圈子

#### 游戏控制相关
- `game:ctrl:view`: 查看中控账号
- `game:ctrl:create`: 创建中控账号
- `game:ctrl:update`: 更新中控账号
- `game:ctrl:delete`: 删除中控账号
- `game:account:view`: 查看游戏账号
- `game:account:bind`: 绑定游戏账号
- `game:account:unbind`: 解绑游戏账号

#### 战绩相关
- `battles:view`: 查看战绩
- `battles:export`: 导出战绩

#### 系统管理
- `menu:view`: 查看菜单
- `menu:create`: 创建菜单
- `menu:update`: 更新菜单
- `menu:delete`: 删除菜单
- `role:view`: 查看角色
- `role:create`: 创建角色
- `role:update`: 更新角色
- `role:delete`: 删除角色
- `permission:view`: 查看权限
- `permission:assign`: 分配权限

### 1.3 角色权限分配

#### 超级管理员 (role_id=1)
- 拥有所有权限

#### 店铺管理员 (role_id=2)
- 统计、资金、店铺、游戏控制、战绩相关的所有权限

#### 普通用户 (role_id=3)
- 仅拥有游戏账号相关权限（查看、绑定、解绑）

### 1.4 迁移脚本

执行顺序：
```bash
# 1. 执行主迁移脚本
psql -U B022MC -d your_database -f migrations/20251117_rbac_permission_system.sql

# 2. 或者分步执行 RBAC 相关脚本
cd doc/rbac
psql -U B022MC -d your_database -f base_role.sql
psql -U B022MC -d your_database -f basic_menu.sql
psql -U B022MC -d your_database -f basic_role_menu_rel.sql
psql -U B022MC -d your_database -f basic_permission.sql
```

## 二、后端改造

### 2.1 Model 层

#### BasicPermission (`internal/dal/model/basic/basic_permission.go`)
权限模型，对应 `basic_permission` 表

#### BasicRolePermissionRel
角色权限关联模型

#### BasicMenuButton
菜单按钮模型

### 2.2 Repository 层

#### PermissionRepo (`internal/dal/repo/basic/basic_permission.go`)
提供权限相关的数据库操作：
- `Create/Update/Delete`: 权限CRUD
- `GetByID/GetByCode/List/ListAll`: 查询权限
- `AssignPermissionsToRole/RemovePermissionsFromRole`: 角色权限关联
- `GetRolePermissions/GetRolePermissionCodes`: 获取角色的权限
- `GetUserPermissionsFromTable`: 获取用户权限（从权限表）

#### BasicMenuRepo 扩展
新增方法：
- `GetUserMenus`: 获取用户可访问的菜单列表
- `GetMenuButtons`: 获取菜单的按钮配置

### 2.3 RBAC Store 增强 (`internal/dal/repo/rbac/rbac.go`)

`queryUserPermCodesDB` 方法现在同时从两个来源获取权限：
1. **菜单权限**（兼容旧系统）：从 `basic_menu.auths` 字段获取
2. **权限表**（新系统）：从 `basic_permission` 表获取

两个来源的权限会合并去重后返回。

### 2.4 Service 层

#### BasicPermissionService (`internal/service/basic/basic_permission.go`)
提供权限管理 API：
- `GET /basic/permission/list`: 查询权限列表
- `GET /basic/permission/listAll`: 查询所有权限
- `POST /basic/permission/create`: 创建权限
- `POST /basic/permission/update`: 更新权限
- `POST /basic/permission/delete`: 删除权限
- `GET /basic/permission/role/permissions`: 查询角色权限
- `POST /basic/permission/role/assign`: 为角色分配权限
- `POST /basic/permission/role/remove`: 从角色移除权限

#### BasicMenuService 已有 API
- `GET /basic/menu/me/tree`: 获取当前用户可访问的菜单树
- `GET /basic/menu/me/buttons`: 获取用户在某菜单下的按钮权限

### 2.5 权限中间件

已有的权限中间件（`pkg/plugin/middleware/rbac.go`）：
- `RequirePerm(...perms)`: 要求全部权限（AND逻辑）
- `RequireAnyPerm(...perms)`: 要求任一权限（OR逻辑）
- `AdminOnly()`: 仅限超级管理员

所有需要权限控制的API都已添加相应的中间件：
- 资金相关API：`fund:deposit`, `fund:withdraw`, `fund:force_withdraw`, `fund:limit:update`
- 桌台相关API：`shop:table:view`, `shop:table:dismiss`
- 成员相关API：已有权限控制
- 申请相关API：`shop:apply:view`, `shop:apply:approve`, `shop:apply:reject`
- 战绩相关API：`battles:view`
- 菜单管理API：`menu:view`, `menu:create`, `menu:update`, `menu:delete`

## 三、前端改造

### 3.1 权限 Hooks

#### usePermission (`hooks/use-permission.ts`)
提供权限检查功能：
- `isSuperAdmin`: 是否超级管理员
- `isStoreAdmin`: 是否店铺管理员
- `hasPerm(code)`: 检查单个权限
- `hasAny(codes)`: 检查是否拥有任一权限
- `hasAll(codes)`: 检查是否拥有所有权限

### 3.2 权限组件

#### PermissionGate (`components/auth/PermissionGate.tsx`)
按钮级权限控制组件：
```tsx
<PermissionGate anyOf={['fund:deposit']}>
  <Button>上分</Button>
</PermissionGate>
```

#### RouteGuard (`components/auth/RouteGuard.tsx`)
页面级权限控制组件：
```tsx
<RouteGuard anyOf={['shop:admin:view']}>
  <AdminsView />
</RouteGuard>
```

### 3.3 页面权限控制

所有店铺子页面都已添加 `RouteGuard`：

#### /(shop)/admins.tsx
```tsx
<RouteGuard anyOf={['shop:admin:view', 'shop:admin:assign', 'shop:admin:revoke']}>
  <AdminsView />
</RouteGuard>
```

#### /(shop)/rooms.tsx
```tsx
<RouteGuard anyOf={['game:ctrl:view', 'game:ctrl:create', 'game:ctrl:update', 'game:ctrl:delete']}>
  <CtrlAccountsView />
</RouteGuard>
```

#### /(shop)/fees.tsx
```tsx
<RouteGuard anyOf={['shop:fees:view']}>
  <FeesView />
</RouteGuard>
```

#### /(shop)/balances.tsx
```tsx
<RouteGuard anyOf={['fund:wallet:view']}>
  <BalancesContent />
</RouteGuard>
```

#### /(shop)/members.tsx
```tsx
<RouteGuard anyOf={['shop:member:view', 'shop:member:kick']}>
  <MembersView />
</RouteGuard>
```

#### /(shop)/group-battles.tsx
```tsx
<RouteGuard anyOf={['shop:member:view', 'battles:view']}>
  <GroupBattlesView />
</RouteGuard>
```

#### /(shop)/group-balances.tsx
```tsx
<RouteGuard anyOf={['shop:member:view', 'fund:wallet:view']}>
  <GroupBalancesView />
</RouteGuard>
```

### 3.4 按钮权限控制

#### 资金页面按钮
```tsx
// components/(tabs)/funds/funds-item.tsx
<PermissionGate anyOf={['fund:deposit']}>
  <Button onPress={handleDeposit}>上分</Button>
</PermissionGate>

<PermissionGate anyOf={['fund:withdraw']}>
  <Button onPress={handleWithdraw}>下分</Button>
</PermissionGate>

<PermissionGate anyOf={['fund:limit:update']}>
  <Button onPress={handleUpdateLimit}>设置阈值</Button>
</PermissionGate>
```

#### 桌台页面按钮
```tsx
// components/(tabs)/tables/tables-item.tsx
<PermissionGate anyOf={['shop:table:detail', 'shop:table:view']}>
  <Button>详情</Button>
</PermissionGate>

<PermissionGate anyOf={['shop:table:check', 'shop:table:view']}>
  <Button>检查</Button>
</PermissionGate>

<PermissionGate anyOf={['shop:table:dismiss']}>
  <Button>解散</Button>
</PermissionGate>
```

#### 成员页面按钮
```tsx
// components/(tabs)/members/members-item.tsx
<PermissionGate anyOf={['shop:member:kick']}>
  <Button>踢出成员</Button>
</PermissionGate>

<PermissionGate anyOf={['shop:member:logout']}>
  <Button>成员下线</Button>
</PermissionGate>

<PermissionGate anyOf={['shop:member:update']}>
  <Button>拉入圈子</Button>
</PermissionGate>
```

#### 管理员页面按钮
```tsx
// components/(shop)/admins/admins-item.tsx
<PermissionGate anyOf={['shop:admin:revoke']}>
  <Button>撤销</Button>
</PermissionGate>

// components/(shop)/admins/admins-view.tsx
<PermissionGate anyOf={['shop:admin:assign']}>
  <AdminsAssign />
</PermissionGate>
```

#### 中控账号页面按钮
```tsx
// components/(shop)/ctrl-accounts/ctrl-accounts-item.tsx
<PermissionGate anyOf={['game:ctrl:delete']}>
  <Button>解绑</Button>
</PermissionGate>
```

### 3.5 Tab 导航权限控制

Tab 导航已有基础权限控制：
```tsx
// app/(tabs)/_layout.tsx
<Tabs.Screen
  name="tables"
  options={{
    href: canViewTables ? undefined : null, // 无权限时隐藏Tab
  }}
/>
```

## 四、使用指南

### 4.1 为用户分配角色

```sql
-- 将用户设置为普通用户
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (4, 3);

-- 将用户设置为店铺管理员
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (5, 2);

-- 将用户设置为超级管理员
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (1, 1);
```

### 4.2 查看用户权限

```sql
-- 查看用户的所有权限（从两个来源）
SELECT DISTINCT p.code, p.name, p.category
FROM basic_user_role_rel urr
JOIN basic_role_permission_rel rpr ON rpr.role_id = urr.role_id
JOIN basic_permission p ON p.id = rpr.permission_id
WHERE urr.user_id = 4 AND p.is_deleted = false
UNION
SELECT DISTINCT unnest(string_to_array(m.auths, ',')) as code, '' as name, '' as category
FROM basic_user_role_rel urr
JOIN basic_role_menu_rel rmr ON rmr.role_id = urr.role_id
JOIN basic_menu m ON m.id = rmr.menu_id
WHERE urr.user_id = 4 AND m.auths IS NOT NULL AND m.auths != ''
ORDER BY category, code;
```

### 4.3 为角色添加新权限

```sql
-- 1. 先创建权限（如果不存在）
INSERT INTO basic_permission (code, name, category, description)
VALUES ('shop:report:view', '查看报表', 'shop', '查看店铺报表数据')
ON CONFLICT (code) DO NOTHING;

-- 2. 为角色分配权限
INSERT INTO basic_role_permission_rel (role_id, permission_id)
SELECT 2, id FROM basic_permission WHERE code = 'shop:report:view'
ON CONFLICT DO NOTHING;
```

### 4.4 添加新页面的权限控制

#### 后端
```go
// 在路由注册时添加权限中间件
g.GET("/reports/list", middleware.RequirePerm("shop:report:view"), s.ListReports)
```

#### 前端页面
```tsx
// app/(shop)/reports.tsx
import { ReportsView } from '@/components/(shop)/reports/reports-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function ReportsScreen() {
  return (
    <RouteGuard anyOf={['shop:report:view']}>
      <ReportsView />
    </RouteGuard>
  );
}
```

#### 前端按钮
```tsx
<PermissionGate anyOf={['shop:report:export']}>
  <Button onPress={handleExport}>导出报表</Button>
</PermissionGate>
```

## 五、测试建议

### 5.1 数据库测试
1. 验证权限表数据是否正确插入
2. 测试角色权限关联是否正确
3. 验证权限查询性能

### 5.2 后端API测试
1. 测试无权限用户访问受保护API（应返回403）
2. 测试超级管理员访问所有API（应全部成功）
3. 测试普通用户只能访问被授权的API
4. 测试权限管理API的CRUD功能

### 5.3 前端测试
1. 测试不同角色用户登录后看到的Tab和页面
2. 测试按钮在无权限时是否隐藏
3. 测试路由守卫对无权限访问的拦截
4. 测试权限变更后的实时更新

### 5.4 集成测试场景

#### 场景1：普通用户
1. 登录后应只看到"店铺"和"我的"Tab
2. 可以查看自己的游戏账号、战绩、余额
3. 不能查看其他成员信息
4. 不能进行上下分等操作

#### 场景2：店铺管理员
1. 登录后看到所有Tab
2. 可以查看和管理店铺成员
3. 可以进行上下分操作
4. 可以解散桌台
5. 可以分配/撤销其他管理员
6. 可以查看圈子战绩和余额

#### 场景3：超级管理员
1. 拥有所有权限
2. 可以管理权限和角色
3. 可以查看系统统计数据

## 六、注意事项

1. **权限缓存**：RBAC Store 使用 Redis 缓存用户权限（TTL 10分钟），权限变更后需等待缓存过期或手动清理

2. **权限检查逻辑**：
   - 超级管理员（role_id=1）会自动通过所有权限检查
   - 前端和后端都会检查权限，确保安全性
   - 权限码统一转换为小写进行比较

3. **兼容性**：
   - 新系统同时支持菜单权限（`basic_menu.auths`）和权限表（`basic_permission`）
   - 两个来源的权限会自动合并

4. **数据一致性**：
   - 删除角色前需先删除该角色的所有关联（用户角色关联、角色权限关联、角色菜单关联）
   - 删除权限使用软删除（`is_deleted = true`）

5. **性能优化**：
   - 权限查询已添加适当的索引
   - 使用Redis缓存减少数据库查询
   - 前端权限数据在登录时一次性获取并持久化

## 七、文件清单

### 数据库
- `/migrations/20251117_rbac_permission_system.sql` - 主迁移脚本
- `/doc/rbac/basic_permission.sql` - 权限表定义和数据
- `/doc/rbac/base_role.sql` - 角色定义
- `/doc/rbac/basic_menu.sql` - 菜单定义
- `/doc/rbac/basic_role_menu_rel.sql` - 角色菜单关联
- `/doc/ddl_postgresql.sql` - 完整DDL（已更新）

### 后端
- `/internal/dal/model/basic/basic_permission.go` - 权限模型
- `/internal/dal/repo/basic/basic_permission.go` - 权限仓库
- `/internal/dal/repo/basic/basic_menu.go` - 菜单仓库（已扩展）
- `/internal/dal/repo/rbac/rbac.go` - RBAC存储（已增强）
- `/internal/service/basic/basic_permission.go` - 权限服务

### 前端
- `/components/auth/RouteGuard.tsx` - 路由守卫组件（新增）
- `/components/auth/PermissionGate.tsx` - 权限门控组件（已有）
- `/hooks/use-permission.ts` - 权限钩子
- `/hooks/use-auth-store.ts` - 认证状态管理

### 文档
- `/doc/RBAC_IMPLEMENTATION_SUMMARY.md` - 本文档
- `/doc/rbac/README.md` - RBAC系统说明

## 八、后续优化建议

1. **动态权限管理界面**：
   - 前端添加权限管理页面（CRUD权限）
   - 前端添加角色权限分配界面

2. **权限审计**：
   - 记录权限变更日志
   - 记录用户权限使用情况

3. **权限组**：
   - 支持权限分组，简化权限分配
   - 支持权限模板，快速创建角色

4. **更细粒度的控制**：
   - 数据级权限（如只能查看自己管理的店铺）
   - 字段级权限（如隐藏敏感字段）

5. **权限继承**：
   - 支持角色继承，减少配置重复

---

**文档版本**: 1.0  
**创建日期**: 2025-11-17  
**维护人员**: AI Assistant

