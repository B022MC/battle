-- ========================================
-- 修复房间额度限制功能的权限配置
-- 执行此脚本以确保角色拥有正确的菜单和权限
-- ========================================

-- 1. 确认菜单64存在（如果不存在则创建）
INSERT INTO "basic_menu" (
    "id", "parent_id", "menu_type", "title", "name", "path", "component", 
    "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", 
    "active_path", "auths", "frame_src", "frame_loading", "keep_alive", 
    "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "is_del"
) 
VALUES (
    64, 5, 2, '房间额度', 'shop.room_credits', '/(shop)/room-credits', 'shop/room-credits',
    NULL, '', '', '', '', '', '', 'game:room_credit:view', '', false, false, 
    false, false, true, true, NOW(), NOW(), 0
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    name = EXCLUDED.name,
    path = EXCLUDED.path,
    auths = EXCLUDED.auths,
    updated_at = NOW();

-- 2. 为超级管理员(role_id=1)分配菜单64
INSERT INTO basic_role_menu_rel (role_id, menu_id) 
VALUES (1, 64)
ON CONFLICT DO NOTHING;

-- 3. 为店铺管理员(role_id=2)分配菜单64
INSERT INTO basic_role_menu_rel (role_id, menu_id) 
VALUES (2, 64)
ON CONFLICT DO NOTHING;

-- 4. 删除可能的错误权限关系（旧的47-49）
DELETE FROM basic_role_permission_rel 
WHERE role_id IN (1, 2) 
  AND permission_id IN (47, 48, 49);

-- 5. 为超级管理员(role_id=1)添加所有4个房间额度权限
INSERT INTO basic_role_permission_rel (role_id, permission_id) 
VALUES 
(1, 50),  -- game:room_credit:view
(1, 51),  -- game:room_credit:set
(1, 52),  -- game:room_credit:delete
(1, 53)   -- game:room_credit:check
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 6. 为店铺管理员(role_id=2)添加所有4个权限
INSERT INTO basic_role_permission_rel (role_id, permission_id) 
VALUES 
(2, 50),  -- game:room_credit:view
(2, 51),  -- game:room_credit:set
(2, 52),  -- game:room_credit:delete
(2, 53)   -- game:room_credit:check
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- 验证结果
-- ========================================

-- 查看权限配置
SELECT '=== 权限列表 ===' AS info;
SELECT id, code, name FROM basic_permission WHERE code LIKE 'game:room_credit:%' ORDER BY id;

-- 查看菜单配置
SELECT '=== 菜单配置 ===' AS info;
SELECT id, title, name, path, auths FROM basic_menu WHERE id = 64;

-- 查看超级管理员的权限
SELECT '=== 超级管理员的房间额度权限 ===' AS info;
SELECT p.id, p.code, p.name 
FROM basic_permission p
JOIN basic_role_permission_rel rpr ON p.id = rpr.permission_id
WHERE rpr.role_id = 1 AND p.code LIKE 'game:room_credit:%'
ORDER BY p.id;

-- 查看超级管理员的菜单
SELECT '=== 超级管理员的房间额度菜单 ===' AS info;
SELECT m.id, m.title, m.name, m.path
FROM basic_menu m
JOIN basic_role_menu_rel rmr ON m.id = rmr.menu_id
WHERE rmr.role_id = 1 AND m.id = 64;

-- 查看店铺管理员的权限
SELECT '=== 店铺管理员的房间额度权限 ===' AS info;
SELECT p.id, p.code, p.name 
FROM basic_permission p
JOIN basic_role_permission_rel rpr ON p.id = rpr.permission_id
WHERE rpr.role_id = 2 AND p.code LIKE 'game:room_credit:%'
ORDER BY p.id;

-- 查看店铺管理员的菜单
SELECT '=== 店铺管理员的房间额度菜单 ===' AS info;
SELECT m.id, m.title, m.name, m.path
FROM basic_menu m
JOIN basic_role_menu_rel rmr ON m.id = rmr.menu_id
WHERE rmr.role_id = 2 AND m.id = 64;

-- 完成提示
SELECT '✅ 配置完成！请退出登录后重新登录查看菜单。' AS result;
