-- ============================================
-- 角色菜单关联表 - 基于 battle-reusables React Native 应用
-- ============================================

CREATE TABLE "public"."basic_role_menu_rel" (
    "role_id" int4 NOT NULL,
    "menu_id" int4 NOT NULL,
    CONSTRAINT "basic_role_menu_rel_pkey" PRIMARY KEY ("role_id", "menu_id")
);

ALTER TABLE "public"."basic_role_menu_rel" OWNER TO "B022MC";

COMMENT ON COLUMN "public"."basic_role_menu_rel"."role_id" IS '角色ID';
COMMENT ON COLUMN "public"."basic_role_menu_rel"."menu_id" IS '菜单ID';
COMMENT ON TABLE "public"."basic_role_menu_rel" IS '角色菜单关联表';

-- ============================================
-- 角色1：超级管理员 - 拥有所有菜单权限
-- ============================================
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 1);  -- 首页
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 2);  -- 桌台
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 3);  -- 成员
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 4);  -- 资金
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 5);  -- 店铺
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 6);  -- 我的
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 51); -- 游戏账号
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 52); -- 管理员
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 53); -- 中控账号
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 54); -- 费用设置
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 55); -- 余额筛查
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 56); -- 成员管理
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 57); -- 我的战绩
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 58); -- 我的余额
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 59); -- 圈子战绩
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (1, 60); -- 圈子余额

-- ============================================
-- 角色2：店铺管理员 - 拥有店铺管理相关权限
-- ============================================
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 1);  -- 首页
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 2);  -- 桌台
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 3);  -- 成员
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 4);  -- 资金
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 5);  -- 店铺
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 6);  -- 我的
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 51); -- 游戏账号
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 52); -- 管理员
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 53); -- 中控账号
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 54); -- 费用设置
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 55); -- 余额筛查
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 56); -- 成员管理
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 57); -- 我的战绩
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 58); -- 我的余额
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 59); -- 圈子战绩
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (2, 60); -- 圈子余额

-- ============================================
-- 角色3：普通用户 - 仅拥有基础查看权限
-- ============================================
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (3, 5);  -- 店铺
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (3, 6);  -- 我的
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (3, 51); -- 游戏账号
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (3, 57); -- 我的战绩
INSERT INTO "public"."basic_role_menu_rel" ("role_id", "menu_id") VALUES (3, 58); -- 我的余额