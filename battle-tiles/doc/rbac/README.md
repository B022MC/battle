# RBAC 权限管理系统

基于 **battle-reusables** React Native 移动应用的完整权限管理配置。

## 📋 文件说明

### 1. `base_role.sql`
定义系统中的角色：
- **超级管理员** (id=1, code='super_admin')
- **店铺管理员** (id=2, code='shop_admin')
- **普通用户** (id=3, code='user')

### 2. `basic_menu.sql`
定义应用菜单结构，基于 battle-reusables 的实际页面：

#### 一级菜单（底部标签页）
| ID | 名称 | 路径 | 权限 |
|----|------|------|------|
| 1 | 首页 | `/(tabs)/index` | `stats:view` |
| 2 | 桌台 | `/(tabs)/tables` | `shop:table:view` |
| 3 | 成员 | `/(tabs)/members` | `shop:member:view` |
| 4 | 资金 | `/(tabs)/funds` | `fund:wallet:view` |
| 5 | 店铺 | `/(tabs)/shop` | 无 |
| 6 | 我的 | `/(tabs)/profile` | 无 |

#### 二级菜单（店铺子页面）
| ID | 名称 | 路径 | 权限 |
|----|------|------|------|
| 51 | 游戏账号 | `/(shop)/account` | 无 |
| 52 | 管理员 | `/(shop)/admins` | `shop:admin:view,shop:admin:assign,shop:admin:revoke` |
| 53 | 中控账号 | `/(shop)/rooms` | `game:ctrl:view,game:ctrl:create,game:ctrl:update,game:ctrl:delete` |
| 54 | 费用设置 | `/(shop)/fees` | `shop:fees:view` |
| 55 | 余额筛查 | `/(shop)/balances` | `fund:wallet:view` |
| 56 | 成员管理 | `/(shop)/members` | `shop:member:view,shop:member:kick` |
| 57 | 我的战绩 | `/(shop)/my-battles` | 无 |
| 58 | 我的余额 | `/(shop)/my-balances` | 无 |
| 59 | 圈子战绩 | `/(shop)/group-battles` | `shop:member:view` |
| 60 | 圈子余额 | `/(shop)/group-balances` | `shop:member:view` |

### 3. `basic_role_menu_rel.sql`
定义角色与菜单的关联关系：

#### 超级管理员 (role_id=1)
拥有所有菜单权限（菜单 1-6, 51-60）

#### 店铺管理员 (role_id=2)
拥有所有菜单权限（菜单 1-6, 51-60）

#### 普通用户 (role_id=3)
仅拥有基础权限：
- 店铺 (5)
- 我的 (6)
- 游戏账号 (51)
- 我的战绩 (57)
- 我的余额 (58)

## 🔑 权限码说明

### 统计相关
- `stats:view` - 查看统计数据

### 资金相关
- `fund:wallet:view` - 查看钱包/余额
- `fund:ledger:view` - 查看资金流水
- `fund:deposit` - 上分
- `fund:withdraw` - 下分
- `fund:force_withdraw` - 强制下分
- `fund:limit:update` - 更新额度/禁分设置

### 店铺相关
- `shop:table:view` - 查看桌台
- `shop:table:dismiss` - 解散桌台
- `shop:member:view` - 查看成员
- `shop:member:kick` - 踢出成员
- `shop:admin:view` - 查看管理员
- `shop:admin:assign` - 分配管理员
- `shop:admin:revoke` - 撤销管理员
- `shop:apply:view` - 查看入圈申请
- `shop:apply:approve` - 批准入圈申请
- `shop:apply:reject` - 拒绝入圈申请
- `shop:fees:view` - 查看费用设置
- `shop:group:view` - 查看圈子

### 游戏控制相关
- `game:ctrl:view` - 查看中控账号
- `game:ctrl:create` - 创建中控账号
- `game:ctrl:update` - 更新中控账号
- `game:ctrl:delete` - 删除中控账号

### 系统相关
- `menu:view` - 查看菜单
- `menu:create` - 创建菜单
- `menu:update` - 更新菜单
- `menu:delete` - 删除菜单

## 🚀 使用方法

### 1. 初始化数据库

按顺序执行以下 SQL 文件：

```bash
# 1. 创建角色
psql -U B022MC -d your_database -f base_role.sql

# 2. 创建菜单
psql -U B022MC -d your_database -f basic_menu.sql

# 3. 创建角色-菜单关联
psql -U B022MC -d your_database -f basic_role_menu_rel.sql
```

### 2. 分配角色给用户

```sql
-- 将用户设置为普通用户
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (4, 3);

-- 将用户设置为店铺管理员
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (4, 2);

-- 将用户设置为超级管理员
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (4, 1);
```

### 3. 查看用户权限

```sql
-- 查看用户的角色
SELECT r.* 
FROM basic_role r
JOIN basic_user_role_rel urr ON r.id = urr.role_id
WHERE urr.user_id = 4;

-- 查看用户可访问的菜单
SELECT m.* 
FROM basic_menu m
JOIN basic_role_menu_rel rmr ON m.id = rmr.menu_id
JOIN basic_user_role_rel urr ON rmr.role_id = urr.role_id
WHERE urr.user_id = 4
ORDER BY m.parent_id, m.rank;
```

### 4. 修改用户角色

```sql
-- 删除用户的所有角色
DELETE FROM basic_user_role_rel WHERE user_id = 4;

-- 重新分配角色
INSERT INTO basic_user_role_rel (user_id, role_id) VALUES (4, 3);
```

## 📊 角色权限对比

| 功能 | 超级管理员 | 店铺管理员 | 普通用户 |
|------|-----------|-----------|---------|
| 首页统计 | ✅ | ✅ | ❌ |
| 桌台管理 | ✅ | ✅ | ❌ |
| 成员管理 | ✅ | ✅ | ❌ |
| 资金管理 | ✅ | ✅ | ❌ |
| 店铺入口 | ✅ | ✅ | ✅ |
| 我的页面 | ✅ | ✅ | ✅ |
| 游戏账号 | ✅ | ✅ | ✅ |
| 管理员设置 | ✅ | ✅ | ❌ |
| 中控账号 | ✅ | ✅ | ❌ |
| 费用设置 | ✅ | ✅ | ❌ |
| 余额筛查 | ✅ | ✅ | ❌ |
| 成员管理 | ✅ | ✅ | ❌ |
| 我的战绩 | ✅ | ✅ | ✅ |
| 我的余额 | ✅ | ✅ | ✅ |
| 圈子战绩 | ✅ | ✅ | ❌ |
| 圈子余额 | ✅ | ✅ | ❌ |

## 🔧 自定义配置

### 添加新菜单

```sql
-- 添加一级菜单
INSERT INTO basic_menu (id, parent_id, menu_type, title, name, path, component, rank, redirect, icon, extra_icon, enter_transition, leave_transition, active_path, auths, frame_src, frame_loading, keep_alive, hidden_tag, fixed_tag, show_link, show_parent) 
VALUES (7, -1, 1, '新菜单', 'new_menu', '/(tabs)/new', 'tabs/new', '7', '', 'icon-name', '', '', '', '', 'new:view', '', false, false, false, false, true, true);

-- 添加二级菜单
INSERT INTO basic_menu (id, parent_id, menu_type, title, name, path, component, rank, redirect, icon, extra_icon, enter_transition, leave_transition, active_path, auths, frame_src, frame_loading, keep_alive, hidden_tag, fixed_tag, show_link, show_parent) 
VALUES (71, 7, 2, '子菜单', 'new_menu.sub', '/(new)/sub', 'new/sub', NULL, '', '', '', '', '', '', 'new:sub:view', '', false, false, false, false, true, true);
```

### 为角色分配新菜单

```sql
-- 为超级管理员分配新菜单
INSERT INTO basic_role_menu_rel (role_id, menu_id) VALUES (1, 7);
INSERT INTO basic_role_menu_rel (role_id, menu_id) VALUES (1, 71);
```

## 📝 注意事项

1. **菜单 ID 规则**：
   - 一级菜单：1-9
   - 二级菜单：父菜单ID * 10 + 序号（如 51, 52, 53...）

2. **权限检查**：
   - 菜单的 `auths` 字段为空时，所有用户都可以访问
   - 菜单的 `auths` 字段有值时，用户必须拥有其中任一权限才能访问

3. **角色继承**：
   - 当前系统不支持角色继承
   - 每个角色的权限需要单独配置

4. **数据一致性**：
   - 删除菜单前，先删除 `basic_role_menu_rel` 中的关联记录
   - 删除角色前，先删除 `basic_user_role_rel` 和 `basic_role_menu_rel` 中的关联记录

## 🐛 故障排查

### 用户看不到某个菜单

1. 检查用户是否有对应角色：
```sql
SELECT * FROM basic_user_role_rel WHERE user_id = 4;
```

2. 检查角色是否有菜单权限：
```sql
SELECT * FROM basic_role_menu_rel WHERE role_id = 3 AND menu_id = 1;
```

3. 检查用户是否有菜单要求的权限：
```sql
SELECT p.code 
FROM basic_permission p
JOIN basic_role_permission_rel rpr ON p.id = rpr.permission_id
JOIN basic_user_role_rel urr ON rpr.role_id = urr.role_id
WHERE urr.user_id = 4;
```

### 重置所有权限

```sql
-- 清空所有角色-菜单关联
TRUNCATE TABLE basic_role_menu_rel;

-- 重新执行 basic_role_menu_rel.sql
\i basic_role_menu_rel.sql
```

## 📚 相关文档

- [React Native 应用结构](../../battle-reusables/README.md)
- [权限系统设计](./permission-design.md)
- [API 权限控制](../api/permission-control.md)

