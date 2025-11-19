-- ============================================
-- RBAC 权限系统更新脚本
-- 用于更新现有的权限、角色、菜单数据
-- ============================================

-- ============================================
-- 1. 更新角色信息
-- ============================================

-- 更新超级管理员角色
UPDATE "public"."basic_role" 
SET 
    "name" = '超级管理员',
    "remark" = '拥有系统所有权限，包括权限管理',
    "updated_at" = now(),
    "first_letter" = 'C',
    "pinyin_code" = 'chaojiguanliyuan'
WHERE "code" = 'super_admin' AND "is_deleted" = false;

-- 更新店铺管理员角色
UPDATE "public"."basic_role" 
SET 
    "name" = '店铺管理员',
    "remark" = '拥有店铺管理相关权限',
    "updated_at" = now(),
    "first_letter" = 'D',
    "pinyin_code" = 'dianpuguanliyuan'
WHERE "code" = 'shop_admin' AND "is_deleted" = false;

-- 更新普通用户角色
UPDATE "public"."basic_role" 
SET 
    "name" = '普通用户',
    "remark" = '仅拥有基本查看权限',
    "updated_at" = now(),
    "first_letter" = 'P',
    "pinyin_code" = 'putongyonghu'
WHERE "code" = 'user' AND "is_deleted" = false;

-- ============================================
-- 2. 添加新的菜单（如果不存在）
-- ============================================

-- 添加权限管理菜单
INSERT INTO "public"."basic_menu" 
("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "is_del") 
VALUES 
(61, 5, 2, '权限管理', 'shop.permissions', '/(shop)/permissions', 'shop/permissions', NULL, '', '', '', '', '', '', 'permission:view', '', false, false, false, false, true, true, 0)
ON CONFLICT (id) DO UPDATE SET
    "title" = EXCLUDED."title",
    "auths" = EXCLUDED."auths",
    "updated_at" = now();

-- 添加角色管理菜单
INSERT INTO "public"."basic_menu" 
("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "is_del") 
VALUES 
(62, 5, 2, '角色管理', 'shop.roles', '/(shop)/roles', 'shop/roles', NULL, '', '', '', '', '', '', 'role:view', '', false, false, false, false, true, true, 0)
ON CONFLICT (id) DO UPDATE SET
    "title" = EXCLUDED."title",
    "auths" = EXCLUDED."auths",
    "updated_at" = now();

-- 添加菜单管理菜单
INSERT INTO "public"."basic_menu" 
("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "is_del") 
VALUES 
(63, 5, 2, '菜单管理', 'shop.menus', '/(shop)/menus', 'shop/menus', NULL, '', '', '', '', '', '', 'menu:view', '', false, false, false, false, true, true, 0)
ON CONFLICT (id) DO UPDATE SET
    "title" = EXCLUDED."title",
    "auths" = EXCLUDED."auths",
    "updated_at" = now();

-- ============================================
-- 3. 为超级管理员添加新菜单权限
-- ============================================

-- 添加权限管理菜单权限
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") 
VALUES (1, 61)
ON CONFLICT ("role_id", "menu_id") DO NOTHING;

-- 添加角色管理菜单权限
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") 
VALUES (1, 62)
ON CONFLICT ("role_id", "menu_id") DO NOTHING;

-- 添加菜单管理菜单权限
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") 
VALUES (1, 63)
ON CONFLICT ("role_id", "menu_id") DO NOTHING;

-- ============================================
-- 4. 添加新的系统管理权限（如果不存在）
-- ============================================

-- 添加权限管理相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('permission:create', '创建权限', 'system', '创建新权限'),
('permission:update', '更新权限', 'system', '更新权限信息'),
('permission:delete', '删除权限', 'system', '删除权限')
ON CONFLICT (code) WHERE is_deleted = false DO UPDATE SET
    "name" = EXCLUDED."name",
    "description" = EXCLUDED."description",
    "updated_at" = now();

-- ============================================
-- 5. 为超级管理员添加新权限
-- ============================================

-- 为超级管理员添加所有新权限
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 1, id FROM "public"."basic_permission" 
WHERE code IN ('permission:create', 'permission:update', 'permission:delete')
AND is_deleted = false
ON CONFLICT ("role_id", "permission_id") DO NOTHING;

-- ============================================
-- 6. 更新现有权限的描述信息
-- ============================================

-- 更新统计权限
UPDATE "public"."basic_permission" 
SET "description" = '查看统计数据页面和统计数据', "updated_at" = now()
WHERE "code" = 'stats:view' AND "is_deleted" = false;

-- 更新资金权限
UPDATE "public"."basic_permission" 
SET "description" = '查看钱包余额信息', "updated_at" = now()
WHERE "code" = 'fund:wallet:view' AND "is_deleted" = false;

UPDATE "public"."basic_permission" 
SET "description" = '查看资金流水记录', "updated_at" = now()
WHERE "code" = 'fund:ledger:view' AND "is_deleted" = false;

UPDATE "public"."basic_permission" 
SET "description" = '给成员上分', "updated_at" = now()
WHERE "code" = 'fund:deposit' AND "is_deleted" = false;

UPDATE "public"."basic_permission" 
SET "description" = '给成员下分', "updated_at" = now()
WHERE "code" = 'fund:withdraw' AND "is_deleted" = false;

UPDATE "public"."basic_permission" 
SET "description" = '强制给成员下分', "updated_at" = now()
WHERE "code" = 'fund:force_withdraw' AND "is_deleted" = false;

UPDATE "public"."basic_permission" 
SET "description" = '更新成员额度和禁分设置', "updated_at" = now()
WHERE "code" = 'fund:limit:update' AND "is_deleted" = false;

-- ============================================
-- 7. 添加或更新菜单按钮配置
-- ============================================

-- 权限管理页面按钮
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(61, 'permission_create', '创建权限', 'permission:create'),
(61, 'permission_update', '编辑权限', 'permission:update'),
(61, 'permission_delete', '删除权限', 'permission:delete')
ON CONFLICT ("menu_id", "button_code") DO UPDATE SET
    "button_name" = EXCLUDED."button_name",
    "permission_codes" = EXCLUDED."permission_codes",
    "updated_at" = now();

-- 角色管理页面按钮
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(62, 'role_create', '创建角色', 'role:create'),
(62, 'role_update', '编辑角色', 'role:update'),
(62, 'role_delete', '删除角色', 'role:delete'),
(62, 'role_assign_menu', '分配菜单', 'role:update'),
(62, 'role_assign_permission', '分配权限', 'permission:assign')
ON CONFLICT ("menu_id", "button_code") DO UPDATE SET
    "button_name" = EXCLUDED."button_name",
    "permission_codes" = EXCLUDED."permission_codes",
    "updated_at" = now();

-- 菜单管理页面按钮
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(63, 'menu_create', '创建菜单', 'menu:create'),
(63, 'menu_update', '编辑菜单', 'menu:update'),
(63, 'menu_delete', '删除菜单', 'menu:delete')
ON CONFLICT ("menu_id", "button_code") DO UPDATE SET
    "button_name" = EXCLUDED."button_name",
    "permission_codes" = EXCLUDED."permission_codes",
    "updated_at" = now();

-- ============================================
-- 8. 清理无效的关联数据
-- ============================================

-- 删除已删除角色的菜单关联
DELETE FROM "public"."basic_role_menu_rel" 
WHERE "role_id" IN (
    SELECT id FROM "public"."basic_role" WHERE "is_deleted" = true
);

-- 删除已删除角色的权限关联
DELETE FROM "public"."basic_role_permission_rel" 
WHERE "role_id" IN (
    SELECT id FROM "public"."basic_role" WHERE "is_deleted" = true
);

-- 删除已删除菜单的关联
DELETE FROM "public"."basic_role_menu_rel" 
WHERE "menu_id" IN (
    SELECT id FROM "public"."basic_menu" WHERE "is_del" != 0
);

-- 删除已删除权限的关联
DELETE FROM "public"."basic_role_permission_rel" 
WHERE "permission_id" IN (
    SELECT id FROM "public"."basic_permission" WHERE "is_deleted" = true
);

-- 删除已删除菜单的按钮配置
DELETE FROM "public"."basic_menu_button" 
WHERE "menu_id" IN (
    SELECT id FROM "public"."basic_menu" WHERE "is_del" != 0
);

-- ============================================
-- 9. 刷新超级管理员的所有权限
-- ============================================

-- 删除超级管理员现有的权限关联
DELETE FROM "public"."basic_role_permission_rel" 
WHERE "role_id" = 1;

-- 重新为超级管理员分配所有权限
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 1, id FROM "public"."basic_permission" 
WHERE "is_deleted" = false;

-- ============================================
-- 10. 验证数据完整性
-- ============================================

-- 检查角色数据
SELECT '角色数据:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_role" 
WHERE "is_deleted" = false;

-- 检查菜单数据
SELECT '菜单数据:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_menu" 
WHERE "is_del" = 0;

-- 检查权限数据
SELECT '权限数据:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_permission" 
WHERE "is_deleted" = false;

-- 检查角色菜单关联
SELECT '角色菜单关联:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_role_menu_rel";

-- 检查角色权限关联
SELECT '角色权限关联:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_role_permission_rel";

-- 检查菜单按钮配置
SELECT '菜单按钮配置:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_menu_button";

-- 检查超级管理员权限数量
SELECT '超级管理员权限数:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_role_permission_rel" 
WHERE "role_id" = 1;

-- 检查店铺管理员权限数量
SELECT '店铺管理员权限数:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_role_permission_rel" 
WHERE "role_id" = 2;

-- 检查普通用户权限数量
SELECT '普通用户权限数:' AS check_type, COUNT(*) AS count 
FROM "public"."basic_role_permission_rel" 
WHERE "role_id" = 3;

-- ============================================
-- 完成
-- ============================================

SELECT '✅ RBAC 系统更新完成！' AS status;

