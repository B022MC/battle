-- ============================================
-- RBAC 权限系统完整迁移脚本
-- 日期: 2025-11-17
-- 说明: 添加细粒度权限表，支持按钮级别的权限控制
-- ============================================

-- ============================================
-- 1. 创建权限表
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

COMMENT ON TABLE "public"."basic_permission" IS '基础权限表（细粒度权限定义）';
COMMENT ON COLUMN "public"."basic_permission"."id" IS '权限ID';
COMMENT ON COLUMN "public"."basic_permission"."code" IS '权限编码（唯一标识）';
COMMENT ON COLUMN "public"."basic_permission"."name" IS '权限名称';
COMMENT ON COLUMN "public"."basic_permission"."category" IS '权限分类：stats/fund/shop/game/system';
COMMENT ON COLUMN "public"."basic_permission"."description" IS '权限描述';
COMMENT ON COLUMN "public"."basic_permission"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."basic_permission"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."basic_permission"."is_deleted" IS '是否删除';

CREATE UNIQUE INDEX IF NOT EXISTS "uk_basic_permission_code" ON "public"."basic_permission" ("code") WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS "idx_basic_permission_category" ON "public"."basic_permission" ("category");

-- ============================================
-- 2. 创建角色权限关联表
-- ============================================

CREATE TABLE IF NOT EXISTS "public"."basic_role_permission_rel" (
    "role_id" int4 NOT NULL,
    "permission_id" int4 NOT NULL,
    CONSTRAINT "basic_role_permission_rel_pkey" PRIMARY KEY ("role_id", "permission_id")
);

COMMENT ON TABLE "public"."basic_role_permission_rel" IS '角色权限关联表';
COMMENT ON COLUMN "public"."basic_role_permission_rel"."role_id" IS '角色ID';
COMMENT ON COLUMN "public"."basic_role_permission_rel"."permission_id" IS '权限ID';

-- ============================================
-- 3. 创建菜单按钮表
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

COMMENT ON TABLE "public"."basic_menu_button" IS '菜单按钮权限配置表（UI细粒度控制）';
COMMENT ON COLUMN "public"."basic_menu_button"."id" IS '按钮ID';
COMMENT ON COLUMN "public"."basic_menu_button"."menu_id" IS '所属菜单ID';
COMMENT ON COLUMN "public"."basic_menu_button"."button_code" IS '按钮编码';
COMMENT ON COLUMN "public"."basic_menu_button"."button_name" IS '按钮名称';
COMMENT ON COLUMN "public"."basic_menu_button"."permission_codes" IS '所需权限码（逗号分隔，满足任一即可）';

CREATE INDEX IF NOT EXISTS "idx_menu_button_menu_id" ON "public"."basic_menu_button" ("menu_id");
CREATE UNIQUE INDEX IF NOT EXISTS "uk_menu_button_menu_code" ON "public"."basic_menu_button" ("menu_id", "button_code");

-- ============================================
-- 4. 插入权限数据
-- ============================================

-- 统计相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('stats:view', '查看统计', 'stats', '查看统计数据页面和统计数据')
ON CONFLICT (code) WHERE is_deleted = false DO NOTHING;

-- 资金相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('fund:wallet:view', '查看钱包', 'fund', '查看钱包余额信息'),
('fund:ledger:view', '查看账本', 'fund', '查看资金流水记录'),
('fund:deposit', '上分', 'fund', '给成员上分'),
('fund:withdraw', '下分', 'fund', '给成员下分'),
('fund:force_withdraw', '强制下分', 'fund', '强制给成员下分'),
('fund:limit:update', '更新额度', 'fund', '更新成员额度和禁分设置')
ON CONFLICT (code) WHERE is_deleted = false DO NOTHING;

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
('shop:group:view', '查看圈子', 'shop', '查看圈子信息')
ON CONFLICT (code) WHERE is_deleted = false DO NOTHING;

-- 游戏控制相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('game:ctrl:view', '查看中控账号', 'game', '查看中控账号列表'),
('game:ctrl:create', '创建中控账号', 'game', '创建新的中控账号'),
('game:ctrl:update', '更新中控账号', 'game', '更新中控账号信息'),
('game:ctrl:delete', '删除中控账号', 'game', '删除中控账号'),
('game:account:view', '查看游戏账号', 'game', '查看游戏账号信息'),
('game:account:bind', '绑定游戏账号', 'game', '绑定游戏账号'),
('game:account:unbind', '解绑游戏账号', 'game', '解绑游戏账号')
ON CONFLICT (code) WHERE is_deleted = false DO NOTHING;

-- 战绩相关权限
INSERT INTO "public"."basic_permission" ("code", "name", "category", "description") VALUES
('battles:view', '查看战绩', 'game', '查看圈子或店铺战绩'),
('battles:export', '导出战绩', 'game', '导出战绩数据')
ON CONFLICT (code) WHERE is_deleted = false DO NOTHING;

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
('permission:assign', '分配权限', 'system', '为角色分配权限')
ON CONFLICT (code) WHERE is_deleted = false DO NOTHING;

-- ============================================
-- 5. 为角色分配权限
-- ============================================

-- 超级管理员（role_id=1）拥有所有权限
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 1, id FROM "public"."basic_permission" WHERE is_deleted = false
ON CONFLICT DO NOTHING;

-- 店铺管理员（role_id=2）的权限
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
) AND is_deleted = false
ON CONFLICT DO NOTHING;

-- 普通用户（role_id=3）的权限 - 仅查看自己的数据
INSERT INTO "public"."basic_role_permission_rel" ("role_id", "permission_id") 
SELECT 3, id FROM "public"."basic_permission" 
WHERE code IN (
    'game:account:view', 'game:account:bind', 'game:account:unbind'
) AND is_deleted = false
ON CONFLICT DO NOTHING;

-- ============================================
-- 6. 插入菜单按钮配置
-- ============================================

-- 桌台页面按钮 (menu_id=2)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(2, 'table_detail', '查看详情', 'shop:table:detail,shop:table:view'),
(2, 'table_check', '检查桌台', 'shop:table:check,shop:table:view'),
(2, 'table_dismiss', '解散桌台', 'shop:table:dismiss')
ON CONFLICT (menu_id, button_code) DO NOTHING;

-- 成员页面按钮 (menu_id=3)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(3, 'member_kick', '踢出成员', 'shop:member:kick'),
(3, 'member_logout', '成员下线', 'shop:member:logout'),
(3, 'member_add_group', '拉入圈子', 'shop:member:update'),
(3, 'member_remove_group', '踢出圈子', 'shop:member:kick')
ON CONFLICT (menu_id, button_code) DO NOTHING;

-- 资金页面按钮 (menu_id=4)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(4, 'fund_deposit', '上分', 'fund:deposit'),
(4, 'fund_withdraw', '下分', 'fund:withdraw'),
(4, 'fund_force_withdraw', '强制下分', 'fund:force_withdraw'),
(4, 'fund_update_limit', '更新额度', 'fund:limit:update'),
(4, 'fund_view_ledger', '查看流水', 'fund:ledger:view')
ON CONFLICT (menu_id, button_code) DO NOTHING;

-- 管理员页面按钮 (menu_id=52)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(52, 'admin_assign', '分配管理员', 'shop:admin:assign'),
(52, 'admin_revoke', '撤销管理员', 'shop:admin:revoke')
ON CONFLICT (menu_id, button_code) DO NOTHING;

-- 中控账号页面按钮 (menu_id=53)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(53, 'ctrl_create', '创建账号', 'game:ctrl:create'),
(53, 'ctrl_update', '更新账号', 'game:ctrl:update'),
(53, 'ctrl_delete', '删除账号', 'game:ctrl:delete'),
(53, 'ctrl_login', '登录游戏', 'game:ctrl:view')
ON CONFLICT (menu_id, button_code) DO NOTHING;

-- 费用设置页面按钮 (menu_id=54)
INSERT INTO "public"."basic_menu_button" ("menu_id", "button_code", "button_name", "permission_codes") VALUES 
(54, 'fees_update', '更新费用', 'shop:fees:update')
ON CONFLICT (menu_id, button_code) DO NOTHING;

-- ============================================
-- 7. 更新RBAC查询逻辑 - 同时支持菜单权限和独立权限
-- ============================================

-- 说明：原有的 RBAC 系统从 basic_menu.auths 字段获取权限
-- 新系统从 basic_permission 表获取权限
-- 两者将并存，权限检查时会合并两个来源的权限

-- 验证脚本：查看某个用户的所有权限（从两个来源）
-- SELECT DISTINCT p.code 
-- FROM basic_user_role_rel urr
-- JOIN basic_role_permission_rel rpr ON rpr.role_id = urr.role_id
-- JOIN basic_permission p ON p.id = rpr.permission_id
-- WHERE urr.user_id = ? AND p.is_deleted = false
-- UNION
-- SELECT DISTINCT unnest(string_to_array(m.auths, ',')) as code
-- FROM basic_user_role_rel urr
-- JOIN basic_role_menu_rel rmr ON rmr.role_id = urr.role_id
-- JOIN basic_menu m ON m.id = rmr.menu_id
-- WHERE urr.user_id = ? AND m.auths IS NOT NULL AND m.auths != ''

