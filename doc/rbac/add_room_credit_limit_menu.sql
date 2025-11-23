-- ========================================
-- 房间额度限制功能 - 菜单和权限配置
-- ========================================

-- 1. 添加权限到 basic_permission 表
-- category 使用 'game' 因为这是游戏管理相关的功能
INSERT INTO "basic_permission" ("code", "name", "category", "description", "created_at", "updated_at", "is_deleted") 
VALUES 
('game:room_credit:view', '查看房间额度', 'game', '查看房间额度限制配置', NOW(), NOW(), false),
('game:room_credit:set', '设置房间额度', 'game', '设置/更新房间额度限制', NOW(), NOW(), false),
('game:room_credit:delete', '删除房间额度', 'game', '删除房间额度限制配置', NOW(), NOW(), false),
('game:room_credit:check', '检查玩家额度', 'game', '检查玩家是否满足房间额度要求', NOW(), NOW(), false);

-- 2. 添加菜单到 basic_menu 表
-- 作为"店铺"(id=5)的二级菜单 (menu_type=2)
INSERT INTO "basic_menu" (
    "id", 
    "parent_id", 
    "menu_type", 
    "title", 
    "name", 
    "path", 
    "component", 
    "rank", 
    "redirect", 
    "icon", 
    "extra_icon", 
    "enter_transition", 
    "leave_transition", 
    "active_path", 
    "auths", 
    "frame_src", 
    "frame_loading", 
    "keep_alive", 
    "hidden_tag", 
    "fixed_tag", 
    "show_link", 
    "show_parent", 
    "created_at", 
    "updated_at", 
    "deleted_at", 
    "is_del"
) 
VALUES (
    64,                              -- id (63后面的下一个)
    5,                               -- parent_id (店铺菜单)
    2,                               -- menu_type (2=二级菜单)
    '房间额度',                       -- title
    'shop.room_credits',             -- name
    '/(shop)/room-credits',          -- path
    'shop/room-credits',             -- component
    NULL,                            -- rank (排序，可以后续调整)
    '',                              -- redirect
    '',                              -- icon (二级菜单通常不需要)
    '',                              -- extra_icon
    '',                              -- enter_transition
    '',                              -- leave_transition
    '',                              -- active_path
    'game:room_credit:view',         -- auths (查看权限)
    '',                              -- frame_src
    false,                           -- frame_loading
    false,                           -- keep_alive
    false,                           -- hidden_tag
    false,                           -- fixed_tag
    true,                            -- show_link
    true,                            -- show_parent
    NOW(),                           -- created_at
    NOW(),                           -- updated_at
    NULL,                            -- deleted_at
    0                                -- is_del
);

-- 3. 添加菜单按钮到 basic_menu_button 表
INSERT INTO "basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes", "created_at", "updated_at") 
VALUES 
(64, 'credit_set', '设置额度', 'game:room_credit:set', NOW(), NOW()),
(64, 'credit_delete', '删除额度', 'game:room_credit:delete', NOW(), NOW()),
(64, 'credit_check', '检查玩家', 'game:room_credit:check', NOW(), NOW());

-- 4. (可选) 为超级管理员角色分配菜单权限
-- 取消以下注释可自动为超级管理员分配此菜单
-- INSERT INTO "basic_role_menu_rel" ("role_id", "menu_id") 
-- VALUES (1, 64);

-- 5. 为超级管理员和店铺管理员角色分配权限
-- 超级管理员(role_id=1)拥有所有4个权限
INSERT INTO basic_role_permission_rel (role_id, permission_id) 
SELECT 1, id FROM basic_permission WHERE code IN (
    'game:room_credit:view',
    'game:room_credit:set',
    'game:room_credit:delete',
    'game:room_credit:check'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 店铺管理员(role_id=2)拥有所有4个权限
INSERT INTO basic_role_permission_rel (role_id, permission_id) 
SELECT 2, id FROM basic_permission WHERE code IN (
    'game:room_credit:view',
    'game:room_credit:set',
    'game:room_credit:delete',
    'game:room_credit:check'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 查询新增的权限ID（验证）
SELECT id, code, name FROM basic_permission WHERE code LIKE 'game:room_credit:%' ORDER BY id;
