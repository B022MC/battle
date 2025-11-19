-- ============================================
-- 游戏内申请功能 - 菜单和权限配置
-- 创建时间: 2025-11-20
-- 用途: 为游戏内申请功能添加菜单、权限和按钮
-- ============================================

-- 1. 添加权限（游戏内申请功能）
INSERT INTO "basic_permission" ("id", "code", "name", "category", "description", "created_at", "updated_at", "is_deleted") 
VALUES 
  (47, 'shop:applications:view', '查看游戏申请', 'shop', '查看游戏内申请列表（从Plaza内存读取）', now(), now(), 'f'),
  (48, 'shop:applications:approve', '通过游戏申请', 'shop', '通过游戏内申请', now(), now(), 'f'),
  (49, 'shop:applications:reject', '拒绝游戏申请', 'shop', '拒绝游戏内申请', now(), now(), 'f');

-- 2. 添加菜单（游戏申请管理）
INSERT INTO "basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del") 
VALUES 
  (64, 5, 2, '游戏申请', 'shop.game_applications', '/(shop)/game-applications', 'shop/game-applications', NULL, '', '', '', '', '', '', 'shop:applications:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 3. 添加菜单按钮
INSERT INTO "basic_menu_button" ("id", "menu_id", "button_code", "button_name", "permission_codes", "created_at", "updated_at") 
VALUES 
  (31, 64, 'application_approve', '通过申请', 'shop:applications:approve', now(), now()),
  (32, 64, 'application_reject', '拒绝申请', 'shop:applications:reject', now(), now());

-- 4. 为超级管理员角色分配菜单
INSERT INTO "basic_role_menu_rel" ("role_id", "menu_id") 
VALUES 
  (1, 64);  -- 超级管理员

-- 5. 为店铺管理员角色分配菜单
INSERT INTO "basic_role_menu_rel" ("role_id", "menu_id") 
VALUES 
  (2, 64);  -- 店铺管理员

-- 6. 为超级管理员角色分配权限
INSERT INTO "basic_role_permission_rel" ("role_id", "permission_id") 
VALUES 
  (1, 47),  -- 查看游戏申请
  (1, 48),  -- 通过游戏申请
  (1, 49);  -- 拒绝游戏申请

-- 7. 为店铺管理员角色分配权限
INSERT INTO "basic_role_permission_rel" ("role_id", "permission_id") 
VALUES 
  (2, 47),  -- 查看游戏申请
  (2, 48),  -- 通过游戏申请
  (2, 49);  -- 拒绝游戏申请

-- ============================================
-- 说明：
-- 1. 新增3个权限用于游戏内申请功能
-- 2. 新增1个菜单项 "游戏申请"，挂在 "店铺" 菜单下
-- 3. 新增2个菜单按钮：通过申请、拒绝申请
-- 4. 为超级管理员和店铺管理员分配相关菜单和权限
-- 5. 分运费功能无需额外配置，已使用现有权限：
--    - shop:fees:view（查看费用）
--    - shop:fees:update（更新费用）
-- ============================================
