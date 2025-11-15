-- ============================================
-- 基础菜单表 - 基于 battle-reusables React Native 应用
-- ============================================

CREATE TABLE "public"."basic_menu" (
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

CREATE TRIGGER "update_basic_menu_updated_at" BEFORE UPDATE ON "public"."basic_menu"
    FOR EACH ROW EXECUTE PROCEDURE "public"."update_updated_at_column"();

COMMENT ON COLUMN "public"."basic_menu"."id" IS '菜单ID';
COMMENT ON COLUMN "public"."basic_menu"."parent_id" IS '父级菜单ID，默认为-1表示顶级菜单';
COMMENT ON COLUMN "public"."basic_menu"."menu_type" IS '菜单类型：1=一级菜单，2=二级菜单';
COMMENT ON COLUMN "public"."basic_menu"."title" IS '菜单标题';
COMMENT ON COLUMN "public"."basic_menu"."name" IS '菜单名称（唯一标识）';
COMMENT ON COLUMN "public"."basic_menu"."path" IS '路由路径';
COMMENT ON COLUMN "public"."basic_menu"."component" IS '组件路径';
COMMENT ON COLUMN "public"."basic_menu"."rank" IS '排序';
COMMENT ON COLUMN "public"."basic_menu"."redirect" IS '重定向路径';
COMMENT ON COLUMN "public"."basic_menu"."icon" IS '图标';
COMMENT ON COLUMN "public"."basic_menu"."extra_icon" IS '额外图标';
COMMENT ON COLUMN "public"."basic_menu"."enter_transition" IS '进入动画';
COMMENT ON COLUMN "public"."basic_menu"."leave_transition" IS '离开动画';
COMMENT ON COLUMN "public"."basic_menu"."active_path" IS '激活路径';
COMMENT ON COLUMN "public"."basic_menu"."auths" IS '权限标识（逗号分隔）';
COMMENT ON COLUMN "public"."basic_menu"."frame_src" IS '内嵌iframe地址';
COMMENT ON COLUMN "public"."basic_menu"."frame_loading" IS '是否显示加载动画';
COMMENT ON COLUMN "public"."basic_menu"."keep_alive" IS '是否缓存页面';
COMMENT ON COLUMN "public"."basic_menu"."hidden_tag" IS '是否隐藏标签';
COMMENT ON COLUMN "public"."basic_menu"."fixed_tag" IS '是否固定标签';
COMMENT ON COLUMN "public"."basic_menu"."show_link" IS '是否显示链接';
COMMENT ON COLUMN "public"."basic_menu"."show_parent" IS '是否显示父级';
COMMENT ON TABLE "public"."basic_menu" IS '基础菜单表';

-- ============================================
-- 一级菜单（底部标签页）
-- ============================================

-- 1. 首页
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (1, -1, 1, '首页', 'home', '/(tabs)/index', 'tabs/index', '1', '', 'home', '', '', '', '', 'stats:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 2. 桌台
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (2, -1, 1, '桌台', 'tables', '/(tabs)/tables', 'tabs/tables', '2', '', 'layout-grid', '', '', '', '', 'shop:table:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 3. 成员
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (3, -1, 1, '成员', 'members', '/(tabs)/members', 'tabs/members', '3', '', 'users', '', '', '', '', 'shop:member:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 4. 资金
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (4, -1, 1, '资金', 'funds', '/(tabs)/funds', 'tabs/funds', '4', '', 'wallet', '', '', '', '', 'fund:wallet:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5. 店铺（父菜单）
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (5, -1, 1, '店铺', 'shop', '/(tabs)/shop', 'tabs/shop', '5', '', 'store', '', '', '', '', '', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 6. 我的
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (6, -1, 1, '我的', 'profile', '/(tabs)/profile', 'tabs/profile', '6', '', 'user-circle', '', '', '', '', '', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- ============================================
-- 二级菜单（店铺子页面）
-- ============================================

-- 5.1 游戏账号
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (51, 5, 2, '游戏账号', 'shop.account', '/(shop)/account', 'shop/account', NULL, '', '', '', '', '', '', '', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.2 管理员
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (52, 5, 2, '管理员', 'shop.admins', '/(shop)/admins', 'shop/admins', NULL, '', '', '', '', '', '', 'shop:admin:view,shop:admin:assign,shop:admin:revoke', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.3 中控账号
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (53, 5, 2, '中控账号', 'shop.rooms', '/(shop)/rooms', 'shop/rooms', NULL, '', '', '', '', '', '', 'game:ctrl:view,game:ctrl:create,game:ctrl:update,game:ctrl:delete', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.4 费用设置
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (54, 5, 2, '费用设置', 'shop.fees', '/(shop)/fees', 'shop/fees', NULL, '', '', '', '', '', '', 'shop:fees:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.5 余额筛查
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (55, 5, 2, '余额筛查', 'shop.balances', '/(shop)/balances', 'shop/balances', NULL, '', '', '', '', '', '', 'fund:wallet:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.6 成员管理
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (56, 5, 2, '成员管理', 'shop.members', '/(shop)/members', 'shop/members', NULL, '', '', '', '', '', '', 'shop:member:view,shop:member:kick', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.7 我的战绩
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (57, 5, 2, '我的战绩', 'shop.my_battles', '/(shop)/my-battles', 'shop/my-battles', NULL, '', '', '', '', '', '', '', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.8 我的余额
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (58, 5, 2, '我的余额', 'shop.my_balances', '/(shop)/my-balances', 'shop/my-balances', NULL, '', '', '', '', '', '', '', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.9 圈子战绩
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (59, 5, 2, '圈子战绩', 'shop.group_battles', '/(shop)/group-battles', 'shop/group-battles', NULL, '', '', '', '', '', '', 'shop:member:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);

-- 5.10 圈子余额
INSERT INTO "public"."basic_menu" ("id", "parent_id", "menu_type", "title", "name", "path", "component", "rank", "redirect", "icon", "extra_icon", "enter_transition", "leave_transition", "active_path", "auths", "frame_src", "frame_loading", "keep_alive", "hidden_tag", "fixed_tag", "show_link", "show_parent", "created_at", "updated_at", "deleted_at", "is_del")
VALUES (60, 5, 2, '圈子余额', 'shop.group_balances', '/(shop)/group-balances', 'shop/group-balances', NULL, '', '', '', '', '', '', 'shop:member:view', '', 'f', 'f', 'f', 'f', 't', 't', now(), now(), NULL, 0);