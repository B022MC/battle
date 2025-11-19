-- ============================================
-- 添加菜单管理功能
-- 用于修复缺少的菜单管理页面
-- ============================================

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

-- 2. 为超级管理员添加菜单管理菜单权限
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") 
VALUES (1, 63)
ON CONFLICT ("role_id", "menu_id") DO NOTHING;

-- 3. 添加菜单管理页面按钮配置
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(63, 'menu_create', '创建菜单', 'menu:create'),
(63, 'menu_update', '编辑菜单', 'menu:update'),
(63, 'menu_delete', '删除菜单', 'menu:delete')
ON CONFLICT ("menu_id", "button_code") DO UPDATE SET
    "button_name" = EXCLUDED."button_name",
    "permission_codes" = EXCLUDED."permission_codes",
    "updated_at" = now();

-- 4. 验证数据
SELECT '✅ 菜单管理功能添加完成！' AS status;

-- 查看菜单管理菜单
SELECT id, title, name, path, auths 
FROM "public"."basic_menu" 
WHERE id = 63;

-- 查看超级管理员是否有菜单管理权限
SELECT COUNT(*) as has_permission
FROM "public"."basic_role_menu_rel" 
WHERE role_id = 1 AND menu_id = 63;

-- 查看菜单管理按钮配置
SELECT * 
FROM "public"."basic_menu_button" 
WHERE menu_id = 63;


