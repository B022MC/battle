-- ============================================
-- RBAC 系统完整初始化数据
-- 用于首次部署或重置RBAC数据
-- ============================================

-- ============================================
-- 1. 基础角色表初始化
-- ============================================

-- 确保角色表存在
CREATE TABLE IF NOT EXISTS "public"."basic_role" (
    "id" int4 NOT NULL DEFAULT nextval('basic_role_id_seq'::regclass),
    "code" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
    "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "parent_id" int4 NOT NULL DEFAULT '-1'::integer,
    "remark" text COLLATE "pg_catalog"."default",
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "created_user" int4,
    "updated_at" timestamptz(6),
    "updated_user" int4,
    "first_letter" varchar(50) COLLATE "pg_catalog"."default",
    "pinyin_code" varchar(100) COLLATE "pg_catalog"."default",
    "enable" bool NOT NULL DEFAULT true,
    "is_deleted" bool NOT NULL DEFAULT false,
    CONSTRAINT "basic_role_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."basic_role" OWNER TO "B022MC";

-- 清空并插入初始角色数据
TRUNCATE TABLE "public"."basic_role" RESTART IDENTITY CASCADE;

INSERT INTO "public"."basic_role" 
("id", "code", "name", "parent_id", "remark", "created_at", "created_user", "updated_at", "updated_user", "first_letter", "pinyin_code", "enable", "is_deleted") 
VALUES 
(1, 'super_admin', '超级管理员', -1, '拥有系统所有权限，包括权限管理', now(), NULL, now(), NULL, 'C', 'chaojiguanliyuan', true, false),
(2, 'shop_admin', '店铺管理员', -1, '拥有店铺管理相关权限', now(), NULL, now(), NULL, 'D', 'dianpuguanliyuan', true, false),
(3, 'user', '普通用户', -1, '仅拥有基本查看权限', now(), NULL, now(), NULL, 'P', 'putongyonghu', true, false);

-- 重置序列
SELECT setval('basic_role_id_seq', (SELECT MAX(id) FROM "public"."basic_role"));

-- ============================================
-- 2. 基础菜单表初始化
-- ============================================

-- 确保菜单表存在
CREATE TABLE IF NOT EXISTS "public"."basic_menu" (
    "id" int4 NOT NULL DEFAULT nextval('basic_menu_id_seq'::regclass),
    "parent_id" int4 NOT NULL DEFAULT '-1'::integer,
    "menu_type" int4 NOT NULL,
    "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "path" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "component" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "rank" varchar(255) COLLATE "pg_catalog"."default",
    "redirect" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "icon" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "extra_icon" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "enter_transition" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "leave_transition" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "active_path" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "auths" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "frame_src" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "frame_loading" bool NOT NULL DEFAULT false,
    "keep_alive" bool NOT NULL DEFAULT false,
    "hidden_tag" bool NOT NULL DEFAULT false,
    "fixed_tag" bool NOT NULL DEFAULT false,
    "show_link" bool NOT NULL DEFAULT true,
    "show_parent" bool NOT NULL DEFAULT true,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
    "deleted_at" timestamptz(6),
    "is_del" int2 NOT NULL DEFAULT 0,
    CONSTRAINT "basic_menu_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."basic_menu" OWNER TO "B022MC";

-- 清空并插入初始菜单数据
TRUNCATE TABLE "public"."basic_menu" RESTART IDENTITY CASCADE;

-- 一级菜单
INSERT INTO "public"."basic_menu" 
("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "is_del") 
VALUES 
(1, -1, 1, '首页', 'home', '/(tabs)/index', 'tabs/index', '1', '', 'home', '', '', '', '', 'stats:view', '', false, false, false, false, true, true, 0),
(2, -1, 1, '桌台', 'tables', '/(tabs)/tables', 'tabs/tables', '2', '', 'layout-grid', '', '', '', '', 'shop:table:view', '', false, false, false, false, true, true, 0),
(3, -1, 1, '成员', 'members', '/(tabs)/members', 'tabs/members', '3', '', 'users', '', '', '', '', 'shop:member:view', '', false, false, false, false, true, true, 0),
(4, -1, 1, '资金', 'funds', '/(tabs)/funds', 'tabs/funds', '4', '', 'wallet', '', '', '', '', 'fund:wallet:view', '', false, false, false, false, true, true, 0),
(5, -1, 1, '店铺', 'shop', '/(tabs)/shop', 'tabs/shop', '5', '', 'store', '', '', '', '', '', '', false, false, false, false, true, true, 0),
(6, -1, 1, '我的', 'profile', '/(tabs)/profile', 'tabs/profile', '6', '', 'user-circle', '', '', '', '', '', '', false, false, false, false, true, true, 0);

-- 二级菜单（店铺子页面）
INSERT INTO "public"."basic_menu" 
("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "is_del") 
VALUES 
(51, 5, 2, '游戏账号', 'shop.account', '/(shop)/account', 'shop/account', NULL, '', '', '', '', '', '', '', '', false, false, false, false, true, true, 0),
(52, 5, 2, '管理员', 'shop.admins', '/(shop)/admins', 'shop/admins', NULL, '', '', '', '', '', '', 'shop:admin:view,shop:admin:assign,shop:admin:revoke', '', false, false, false, false, true, true, 0),
(53, 5, 2, '中控账号', 'shop.rooms', '/(shop)/rooms', 'shop/rooms', NULL, '', '', '', '', '', '', 'game:ctrl:view,game:ctrl:create,game:ctrl:update,game:ctrl:delete', '', false, false, false, false, true, true, 0),
(54, 5, 2, '费用设置', 'shop.fees', '/(shop)/fees', 'shop/fees', NULL, '', '', '', '', '', '', 'shop:fees:view', '', false, false, false, false, true, true, 0),
(55, 5, 2, '余额筛查', 'shop.balances', '/(shop)/balances', 'shop/balances', NULL, '', '', '', '', '', '', 'fund:wallet:view', '', false, false, false, false, true, true, 0),
(56, 5, 2, '成员管理', 'shop.members', '/(shop)/members', 'shop/members', NULL, '', '', '', '', '', '', 'shop:member:view,shop:member:kick', '', false, false, false, false, true, true, 0),
(57, 5, 2, '我的战绩', 'shop.my_battles', '/(shop)/my-battles', 'shop/my-battles', NULL, '', '', '', '', '', '', '', '', false, false, false, false, true, true, 0),
(58, 5, 2, '我的余额', 'shop.my_balances', '/(shop)/my-balances', 'shop/my-balances', NULL, '', '', '', '', '', '', '', '', false, false, false, false, true, true, 0),
(59, 5, 2, '圈子战绩', 'shop.group_battles', '/(shop)/group-battles', 'shop/group-battles', NULL, '', '', '', '', '', '', 'shop:member:view', '', false, false, false, false, true, true, 0),
(60, 5, 2, '圈子余额', 'shop.group_balances', '/(shop)/group-balances', 'shop/group-balances', NULL, '', '', '', '', '', '', 'shop:member:view', '', false, false, false, false, true, true, 0),
(61, 5, 2, '权限管理', 'shop.permissions', '/(shop)/permissions', 'shop/permissions', NULL, '', '', '', '', '', '', 'permission:view', '', false, false, false, false, true, true, 0),
(62, 5, 2, '角色管理', 'shop.roles', '/(shop)/roles', 'shop/roles', NULL, '', '', '', '', '', '', 'role:view', '', false, false, false, false, true, true, 0),
(63, 5, 2, '菜单管理', 'shop.menus', '/(shop)/menus', 'shop/menus', NULL, '', '', '', '', '', '', 'menu:view', '', false, false, false, false, true, true, 0);

-- 重置序列
SELECT setval('basic_menu_id_seq', (SELECT MAX(id) FROM "public"."basic_menu"));

-- ============================================
-- 3. 角色菜单关联表初始化
-- ============================================

CREATE TABLE IF NOT EXISTS "public"."basic_role_menu_rel" (
    "role_id" int4 NOT NULL,
    "menu_id" int4 NOT NULL,
    CONSTRAINT "basic_role_menu_rel_pkey" PRIMARY KEY ("role_id", "menu_id")
);

ALTER TABLE "public"."basic_role_menu_rel" OWNER TO "B022MC";

-- 清空并插入角色菜单关联
TRUNCATE TABLE "public"."basic_role_menu_rel";

-- 超级管理员（所有菜单）
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES 
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6),
(1, 51), (1, 52), (1, 53), (1, 54), (1, 55), (1, 56), (1, 57), (1, 58), (1, 59), (1, 60), (1, 61), (1, 62), (1, 63);

-- 店铺管理员（店铺管理相关菜单）
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES 
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6),
(2, 51), (2, 52), (2, 53), (2, 54), (2, 55), (2, 56), (2, 57), (2, 58), (2, 59), (2, 60);

-- 普通用户（基础菜单）
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES 
(3, 5), (3, 6), (3, 51), (3, 57), (3, 58);

-- ============================================
-- 4. 基础权限表初始化
-- ============================================

CREATE TABLE IF NOT EXISTS "public"."basic_permission" (
    "id" SERIAL PRIMARY KEY,
    "code" varchar(100) NOT NULL,
    "name" varchar(255) NOT NULL,
    "category" varchar(50) NOT NULL,
    "description" text,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
    "is_deleted" bool NOT NULL DEFAULT false
);

ALTER TABLE "public"."basic_permission" OWNER TO "B022MC";

CREATE UNIQUE INDEX IF NOT EXISTS "uk_basic_permission_code" ON "public"."basic_permission" USING btree ("code") WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS "idx_basic_permission_category" ON "public"."basic_permission" USING btree ("category");

-- 清空并插入权限数据
TRUNCATE TABLE "public"."basic_permission" RESTART IDENTITY CASCADE;

-- 统计相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('stats:view', '查看统计', 'stats', '查看统计数据页面和统计数据');

-- 资金相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('fund:wallet:view', '查看钱包', 'fund', '查看钱包余额信息'),
('fund:ledger:view', '查看账本', 'fund', '查看资金流水记录'),
('fund:deposit', '上分', 'fund', '给成员上分'),
('fund:withdraw', '下分', 'fund', '给成员下分'),
('fund:force_withdraw', '强制下分', 'fund', '强制给成员下分'),
('fund:limit:update', '更新额度', 'fund', '更新成员额度和禁分设置');

-- 店铺相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('shop:table:view', '查看桌台', 'shop', '查看桌台列表'),
('shop:table:detail', '查看桌台详情', 'shop', '查看桌台详细信息'),
('shop:table:check', '检查桌台', 'shop', '检查桌台状态'),
('shop:table:dismiss', '解散桌台', 'shop', '解散游戏桌台'),
('shop:member:view', '查看成员', 'shop', '查看店铺成员列表'),
('shop:member:kick', '踢出成员', 'shop', '将成员踢出店铺'),
('shop:member:logout', '成员下线', 'shop', '强制成员下线'),
('shop:member:update', '更新成员', 'shop', '更新成员信息（拉入/踢出圈子等）'),
('shop:admin:view', '查看管理员', 'shop', '查看店铺管理员列表'),
('shop:admin:assign', '分配管理员', 'shop', '分配用户为管理员'),
('shop:admin:revoke', '撤销管理员', 'shop', '撤销用户的管理员权限'),
('shop:apply:view', '查看申请', 'shop', '查看入圈申请列表'),
('shop:apply:approve', '批准申请', 'shop', '批准入圈申请'),
('shop:apply:reject', '拒绝申请', 'shop', '拒绝入圈申请'),
('shop:fees:view', '查看费用', 'shop', '查看费用设置'),
('shop:fees:update', '更新费用', 'shop', '更新费用设置'),
('shop:group:view', '查看圈子', 'shop', '查看圈子信息');

-- 游戏控制相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('game:ctrl:view', '查看中控账号', 'game', '查看中控账号列表'),
('game:ctrl:create', '创建中控账号', 'game', '创建新的中控账号'),
('game:ctrl:update', '更新中控账号', 'game', '更新中控账号信息'),
('game:ctrl:delete', '删除中控账号', 'game', '删除中控账号'),
('game:account:view', '查看游戏账号', 'game', '查看游戏账号信息'),
('game:account:bind', '绑定游戏账号', 'game', '绑定游戏账号'),
('game:account:unbind', '解绑游戏账号', 'game', '解绑游戏账号');

-- 战绩相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('battles:view', '查看战绩', 'game', '查看圈子或店铺战绩'),
('battles:export', '导出战绩', 'game', '导出战绩数据');

-- 系统管理权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('menu:view', '查看菜单', 'system', '查看菜单列表'),
('menu:create', '创建菜单', 'system', '创建新菜单'),
('menu:update', '更新菜单', 'system', '更新菜单信息'),
('menu:delete', '删除菜单', 'system', '删除菜单'),
('role:view', '查看角色', 'system', '查看角色列表'),
('role:create', '创建角色', 'system', '创建新角色'),
('role:update', '更新角色', 'system', '更新角色信息'),
('role:delete', '删除角色', 'system', '删除角色'),
('permission:view', '查看权限', 'system', '查看权限列表'),
('permission:create', '创建权限', 'system', '创建新权限'),
('permission:update', '更新权限', 'system', '更新权限信息'),
('permission:delete', '删除权限', 'system', '删除权限'),
('permission:assign', '分配权限', 'system', '为角色分配权限');

-- ============================================
-- 5. 角色权限关联表初始化
-- ============================================

CREATE TABLE IF NOT EXISTS "public"."basic_role_permission_rel" (
    "role_id" int4 NOT NULL,
    "permission_id" int4 NOT NULL,
    CONSTRAINT "basic_role_permission_rel_pkey" PRIMARY KEY ("role_id", "permission_id")
);

ALTER TABLE "public"."basic_role_permission_rel" OWNER TO "B022MC";

-- 清空并插入角色权限关联
TRUNCATE TABLE "public"."basic_role_permission_rel";

-- 超级管理员（所有权限）
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 1, id FROM "public"."basic_permission" WHERE is_deleted = false;

-- 店铺管理员（店铺管理相关权限）
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 2, id FROM "public"."basic_permission" 
WHERE code IN (
    'stats:view',
    'fund:wallet:view', 'fund:ledger:view', 'fund:deposit', 'fund:withdraw', 'fund:limit:update',
    'shop:table:view', 'shop:table:detail', 'shop:table:check', 'shop:table:dismiss',
    'shop:member:view', 'shop:member:kick', 'shop:member:logout', 'shop:member:update',
    'shop:admin:view', 'shop:admin:assign', 'shop:admin:revoke',
    'shop:apply:view', 'shop:apply:approve', 'shop:apply:reject',
    'shop:fees:view', 'shop:fees:update',
    'shop:group:view',
    'game:ctrl:view', 'game:ctrl:create', 'game:ctrl:update', 'game:ctrl:delete',
    'game:account:view',
    'battles:view', 'battles:export'
) AND is_deleted = false;

-- 普通用户（基础权限）
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 3, id FROM "public"."basic_permission" 
WHERE code IN (
    'game:account:view', 'game:account:bind', 'game:account:unbind'
) AND is_deleted = false;

-- ============================================
-- 6. 菜单按钮表初始化
-- ============================================

CREATE TABLE IF NOT EXISTS "public"."basic_menu_button" (
    "id" SERIAL PRIMARY KEY,
    "menu_id" int4 NOT NULL,
    "button_code" varchar(100) NOT NULL,
    "button_name" varchar(255) NOT NULL,
    "permission_codes" varchar(500) NOT NULL,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now()
);

ALTER TABLE "public"."basic_menu_button" OWNER TO "B022MC";

CREATE INDEX IF NOT EXISTS "idx_menu_button_menu_id" ON "public"."basic_menu_button" USING btree ("menu_id");
CREATE UNIQUE INDEX IF NOT EXISTS "uk_menu_button_menu_code" ON "public"."basic_menu_button" USING btree ("menu_id", "button_code");

-- 清空并插入菜单按钮配置
TRUNCATE TABLE "public"."basic_menu_button" RESTART IDENTITY;

-- 桌台页面按钮 (menu_id=2)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(2, 'table_detail', '查看详情', 'shop:table:detail,shop:table:view'),
(2, 'table_check', '检查桌台', 'shop:table:check,shop:table:view'),
(2, 'table_dismiss', '解散桌台', 'shop:table:dismiss');

-- 成员页面按钮 (menu_id=3)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(3, 'member_kick', '踢出成员', 'shop:member:kick'),
(3, 'member_logout', '成员下线', 'shop:member:logout'),
(3, 'member_add_group', '拉入圈子', 'shop:member:update'),
(3, 'member_remove_group', '踢出圈子', 'shop:member:kick');

-- 资金页面按钮 (menu_id=4)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(4, 'fund_deposit', '上分', 'fund:deposit'),
(4, 'fund_withdraw', '下分', 'fund:withdraw'),
(4, 'fund_force_withdraw', '强制下分', 'fund:force_withdraw'),
(4, 'fund_update_limit', '更新额度', 'fund:limit:update'),
(4, 'fund_view_ledger', '查看流水', 'fund:ledger:view');

-- 管理员页面按钮 (menu_id=52)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(52, 'admin_assign', '分配管理员', 'shop:admin:assign'),
(52, 'admin_revoke', '撤销管理员', 'shop:admin:revoke');

-- 中控账号页面按钮 (menu_id=53)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(53, 'ctrl_create', '创建账号', 'game:ctrl:create'),
(53, 'ctrl_update', '更新账号', 'game:ctrl:update'),
(53, 'ctrl_delete', '删除账号', 'game:ctrl:delete'),
(53, 'ctrl_login', '登录游戏', 'game:ctrl:view');

-- 费用设置页面按钮 (menu_id=54)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(54, 'fees_update', '更新费用', 'shop:fees:update');

-- 权限管理页面按钮 (menu_id=61)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(61, 'permission_create', '创建权限', 'permission:create'),
(61, 'permission_update', '编辑权限', 'permission:update'),
(61, 'permission_delete', '删除权限', 'permission:delete');

-- 角色管理页面按钮 (menu_id=62)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(62, 'role_create', '创建角色', 'role:create'),
(62, 'role_update', '编辑角色', 'role:update'),
(62, 'role_delete', '删除角色', 'role:delete'),
(62, 'role_assign_menu', '分配菜单', 'role:update'),
(62, 'role_assign_permission', '分配权限', 'permission:assign');

-- 菜单管理页面按钮 (menu_id=63)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(63, 'menu_create', '创建菜单', 'menu:create'),
(63, 'menu_update', '编辑菜单', 'menu:update'),
(63, 'menu_delete', '删除菜单', 'menu:delete');

-- ============================================
-- 7. 用户角色关联表（如果需要）
-- ============================================

CREATE TABLE IF NOT EXISTS "public"."basic_user_role_rel" (
    "user_id" int4 NOT NULL,
    "role_id" int4 NOT NULL,
    CONSTRAINT "basic_user_role_rel_pkey" PRIMARY KEY ("user_id", "role_id")
);

ALTER TABLE "public"."basic_user_role_rel" OWNER TO "B022MC";

COMMENT ON COLUMN "public"."basic_user_role_rel"."user_id" IS '用户ID';
COMMENT ON COLUMN "public"."basic_user_role_rel"."role_id" IS '角色ID';
COMMENT ON TABLE "public"."basic_user_role_rel" IS '用户角色关联表';

-- ============================================
-- 完成
-- ============================================

SELECT '✅ RBAC 系统初始化完成！' AS status;
SELECT COUNT(*) AS role_count FROM "public"."basic_role" WHERE is_deleted = false;
SELECT COUNT(*) AS menu_count FROM "public"."basic_menu" WHERE is_del = 0;
SELECT COUNT(*) AS permission_count FROM "public"."basic_permission" WHERE is_deleted = false;
SELECT COUNT(*) AS role_menu_count FROM "public"."basic_role_menu_rel";
SELECT COUNT(*) AS role_permission_count FROM "public"."basic_role_permission_rel";
SELECT COUNT(*) AS menu_button_count FROM "public"."basic_menu_button";

