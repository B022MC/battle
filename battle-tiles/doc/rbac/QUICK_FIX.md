# 菜单管理功能快速修复

## 🎯 问题

管理员权限菜单管理不见了！

## ✅ 解决方案

### 方法一：执行快速修复脚本（推荐）

```bash
# 进入脚本目录
cd battle-tiles/doc/rbac

# 执行修复脚本
psql -U B022MC -d your_database -f 02_add_menu_management.sql
```

### 方法二：手动执行SQL

如果你无法使用脚本，可以直接在数据库中执行以下SQL：

```sql
-- 1. 添加菜单管理菜单项
INSERT INTO "public"."basic_menu" 
("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "is_del") 
VALUES 
(63, 5, 2, '菜单管理', 'shop.menus', '/(shop)/menus', 'shop/menus', NULL, '', '', '', '', '', '', 'menu:view', '', false, false, false, false, true, true, 0)
ON CONFLICT (id) DO UPDATE SET
    "title" = EXCLUDED."title",
    "name" = EXCLUDED."name",
    "path" = EXCLUDED."path",
    "component" = EXCLUDED."component",
    "auths" = EXCLUDED."auths",
    "updated_at" = now();

-- 2. 为超级管理员添加菜单管理权限
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") 
VALUES (1, 63)
ON CONFLICT ("role_id", "menu_id") DO NOTHING;

-- 3. 添加菜单按钮配置
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(63, 'menu_create', '创建菜单', 'menu:create'),
(63, 'menu_update', '编辑菜单', 'menu:update'),
(63, 'menu_delete', '删除菜单', 'menu:delete')
ON CONFLICT ("menu_id", "button_code") DO UPDATE SET
    "button_name" = EXCLUDED."button_name",
    "permission_codes" = EXCLUDED."permission_codes",
    "updated_at" = now();
```

## 🔍 验证

执行完成后，运行以下查询验证：

```sql
-- 检查菜单是否存在
SELECT id, title, name, path, auths 
FROM basic_menu 
WHERE id = 63;

-- 检查超级管理员权限
SELECT COUNT(*) as has_permission
FROM basic_role_menu_rel 
WHERE role_id = 1 AND menu_id = 63;

-- 检查按钮配置
SELECT * FROM basic_menu_button WHERE menu_id = 63;
```

预期结果：
- ✅ 查询到菜单ID=63的记录
- ✅ has_permission = 1
- ✅ 返回3条按钮配置记录

## 🚀 重启应用

数据库更新后，需要重启前端应用：

```bash
cd battle-reusables
npm start
```

或者如果使用Web版本：

```bash
cd battle-reusables
npm run build
```

## ✔️ 测试

1. 以超级管理员身份登录
2. 进入"店铺"菜单
3. 应该能看到以下管理页面：
   - ✅ 权限管理
   - ✅ 角色管理
   - ✅ 菜单管理（新增）

4. 点击"菜单管理"，测试功能：
   - 查看菜单列表
   - 创建新菜单
   - 编辑菜单
   - 删除菜单

## 📝 文件清单

修复涉及的文件：

### 数据库
- ✅ `battle-tiles/doc/rbac/00_init_data.sql` - 已更新
- ✅ `battle-tiles/doc/rbac/01_update_permissions.sql` - 已更新
- ✅ `battle-tiles/doc/rbac/02_add_menu_management.sql` - 新增

### 前端页面
- ✅ `battle-reusables/app/(shop)/menus.tsx` - 新增

### 前端组件
- ✅ `battle-reusables/components/(shop)/menus/menus-view.tsx` - 新增
- ✅ `battle-reusables/components/(shop)/menus/menu-list.tsx` - 新增
- ✅ `battle-reusables/components/(shop)/menus/menu-form.tsx` - 新增

### 服务层
- ✅ `battle-reusables/services/basic/menu.ts` - 已更新

### 文档
- ✅ `battle-tiles/doc/rbac/README.md` - 已更新
- ✅ `battle-tiles/doc/rbac/MENU_MANAGEMENT_FIX.md` - 新增
- ✅ `battle-tiles/doc/rbac/QUICK_FIX.md` - 本文件

## ❓ 如果还是看不到

### 检查1：确认用户角色

```sql
-- 查询你的用户角色
SELECT u.id, u.username, r.id as role_id, r.name as role_name
FROM game_member u
LEFT JOIN basic_user_role_rel urr ON urr.user_id = u.id
LEFT JOIN basic_role r ON r.id = urr.role_id
WHERE u.id = 你的用户ID;
```

确认你的账号有超级管理员角色（role_id = 1）

### 检查2：清理缓存

如果有Redis缓存，清理一下：

```bash
redis-cli FLUSHDB
```

### 检查3：重新登录

退出登录，然后重新登录，让系统重新加载菜单。

## 📞 需要帮助？

如果按照以上步骤操作后仍然有问题，请检查：

1. 后端服务是否正常运行
2. 前端是否已重新构建
3. 浏览器是否有缓存（Ctrl+Shift+R 强制刷新）
4. 查看浏览器控制台是否有错误
5. 查看后端日志是否有错误

---

**修复时间**: 2025-11-18  
**问题**: 管理员权限菜单管理不见了  
**状态**: ✅ 已修复


