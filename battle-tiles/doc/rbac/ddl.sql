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
)
;

ALTER TABLE "public"."basic_menu"
    OWNER TO "B022MC";

CREATE TRIGGER "update_basic_menu_updated_at" BEFORE UPDATE ON "public"."basic_menu"
    FOR EACH ROW
    EXECUTE PROCEDURE "public"."update_updated_at_column"();

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
CREATE TABLE "public"."basic_menu_button" (
                                              "id" int4 NOT NULL DEFAULT nextval('basic_menu_button_id_seq'::regclass),
                                              "menu_id" int4 NOT NULL,
                                              "button_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
                                              "button_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
                                              "permission_codes" varchar(500) COLLATE "pg_catalog"."default" NOT NULL,
                                              "created_at" timestamptz(6) NOT NULL DEFAULT now(),
                                              "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
                                              CONSTRAINT "basic_menu_button_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."basic_menu_button"
    OWNER TO "B022MC";

CREATE INDEX "idx_menu_button_menu_id" ON "public"."basic_menu_button" USING btree (
    "menu_id" "pg_catalog"."int4_ops" ASC NULLS LAST
    );

CREATE UNIQUE INDEX "uk_menu_button_menu_code" ON "public"."basic_menu_button" USING btree (
    "menu_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
    "button_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

COMMENT ON COLUMN "public"."basic_menu_button"."id" IS '按钮ID';

COMMENT ON COLUMN "public"."basic_menu_button"."menu_id" IS '所属菜单ID';

COMMENT ON COLUMN "public"."basic_menu_button"."button_code" IS '按钮编码';

COMMENT ON COLUMN "public"."basic_menu_button"."button_name" IS '按钮名称';

COMMENT ON COLUMN "public"."basic_menu_button"."permission_codes" IS '所需权限码（逗号分隔，满足任一即可）';

COMMENT ON TABLE "public"."basic_menu_button" IS '菜单按钮权限配置表（UI细粒度控制）';
CREATE TABLE "public"."basic_permission" (
                                             "id" int4 NOT NULL DEFAULT nextval('basic_permission_id_seq'::regclass),
                                             "code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
                                             "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
                                             "category" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
                                             "description" text COLLATE "pg_catalog"."default",
                                             "created_at" timestamptz(6) NOT NULL DEFAULT now(),
                                             "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
                                             "is_deleted" bool NOT NULL DEFAULT false,
                                             CONSTRAINT "basic_permission_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."basic_permission"
    OWNER TO "B022MC";

CREATE INDEX "idx_basic_permission_category" ON "public"."basic_permission" USING btree (
    "category" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE UNIQUE INDEX "uk_basic_permission_code" ON "public"."basic_permission" USING btree (
    "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    ) WHERE is_deleted = false;

COMMENT ON COLUMN "public"."basic_permission"."id" IS '权限ID';

COMMENT ON COLUMN "public"."basic_permission"."code" IS '权限编码（唯一标识）';

COMMENT ON COLUMN "public"."basic_permission"."name" IS '权限名称';

COMMENT ON COLUMN "public"."basic_permission"."category" IS '权限分类：stats/fund/shop/game/system';

COMMENT ON COLUMN "public"."basic_permission"."description" IS '权限描述';

COMMENT ON COLUMN "public"."basic_permission"."created_at" IS '创建时间';

COMMENT ON COLUMN "public"."basic_permission"."updated_at" IS '更新时间';

COMMENT ON COLUMN "public"."basic_permission"."is_deleted" IS '是否删除';

COMMENT ON TABLE "public"."basic_permission" IS '基础权限表（细粒度权限定义）';
CREATE TABLE "public"."basic_role" (
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
)
;

ALTER TABLE "public"."basic_role"
    OWNER TO "B022MC";

CREATE INDEX "idx_basic_role_enable" ON "public"."basic_role" USING btree (
    "enable" "pg_catalog"."bool_ops" ASC NULLS LAST
    );

CREATE INDEX "idx_basic_role_is_deleted" ON "public"."basic_role" USING btree (
    "is_deleted" "pg_catalog"."bool_ops" ASC NULLS LAST
    );

CREATE UNIQUE INDEX "uk_basic_role_code" ON "public"."basic_role" USING btree (
    "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    ) WHERE is_deleted = false;

COMMENT ON COLUMN "public"."basic_role"."id" IS '角色ID';

COMMENT ON COLUMN "public"."basic_role"."code" IS '角色编码';

COMMENT ON COLUMN "public"."basic_role"."name" IS '角色名称';

COMMENT ON COLUMN "public"."basic_role"."parent_id" IS '父级角色ID，默认为-1表示顶级角色';

COMMENT ON COLUMN "public"."basic_role"."remark" IS '备注';

COMMENT ON COLUMN "public"."basic_role"."created_at" IS '创建时间';

COMMENT ON COLUMN "public"."basic_role"."created_user" IS '创建人用户ID';

COMMENT ON COLUMN "public"."basic_role"."updated_at" IS '更新时间';

COMMENT ON COLUMN "public"."basic_role"."updated_user" IS '更新人用户ID';

COMMENT ON COLUMN "public"."basic_role"."first_letter" IS '名称首字母';

COMMENT ON COLUMN "public"."basic_role"."pinyin_code" IS '名称全拼';

COMMENT ON COLUMN "public"."basic_role"."enable" IS '是否启用';

COMMENT ON COLUMN "public"."basic_role"."is_deleted" IS '是否删除';

COMMENT ON TABLE "public"."basic_role" IS '基础角色表';
CREATE TABLE "public"."basic_role_menu_rel" (
                                                "role_id" int4 NOT NULL,
                                                "menu_id" int4 NOT NULL,
                                                CONSTRAINT "basic_role_menu_rel_pkey" PRIMARY KEY ("role_id", "menu_id")
)
;

ALTER TABLE "public"."basic_role_menu_rel"
    OWNER TO "B022MC";

COMMENT ON COLUMN "public"."basic_role_menu_rel"."role_id" IS '角色ID';

COMMENT ON COLUMN "public"."basic_role_menu_rel"."menu_id" IS '菜单ID';

COMMENT ON TABLE "public"."basic_role_menu_rel" IS '角色菜单关联表';
CREATE TABLE "public"."basic_role_permission_rel" (
                                                      "role_id" int4 NOT NULL,
                                                      "permission_id" int4 NOT NULL,
                                                      CONSTRAINT "basic_role_permission_rel_pkey" PRIMARY KEY ("role_id", "permission_id")
)
;

ALTER TABLE "public"."basic_role_permission_rel"
    OWNER TO "B022MC";

COMMENT ON COLUMN "public"."basic_role_permission_rel"."role_id" IS '角色ID';

COMMENT ON COLUMN "public"."basic_role_permission_rel"."permission_id" IS '权限ID';

COMMENT ON TABLE "public"."basic_role_permission_rel" IS '角色权限关联表';
CREATE TABLE "public"."basic_user" (
                                       "id" int4 NOT NULL DEFAULT nextval('basic_user_id_seq'::regclass),
                                       "username" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
                                       "password" varchar(255) COLLATE "pg_catalog"."default",
                                       "salt" varchar(50) COLLATE "pg_catalog"."default",
                                       "wechat_id" varchar(64) COLLATE "pg_catalog"."default",
                                       "avatar" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
                                       "nick_name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
                                       "game_nickname" varchar(64) COLLATE "pg_catalog"."default",
                                       "introduction" text COLLATE "pg_catalog"."default",
                                       "role" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'user'::character varying,
                                       "pinyin_code" varchar(100) COLLATE "pg_catalog"."default",
                                       "first_letter" varchar(50) COLLATE "pg_catalog"."default",
                                       "last_login_at" timestamptz(6),
                                       "created_at" timestamptz(6) NOT NULL DEFAULT now(),
                                       "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
                                       "deleted_at" timestamptz(6),
                                       "is_del" int2 NOT NULL DEFAULT 0,
                                       CONSTRAINT "basic_user_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."basic_user"
    OWNER TO "B022MC";

CREATE INDEX "idx_basic_user_role" ON "public"."basic_user" USING btree (
    "role" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE UNIQUE INDEX "uk_basic_user_username" ON "public"."basic_user" USING btree (
    "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE TRIGGER "update_basic_user_updated_at" BEFORE UPDATE ON "public"."basic_user"
    FOR EACH ROW
    EXECUTE PROCEDURE "public"."update_updated_at_column"();

COMMENT ON COLUMN "public"."basic_user"."id" IS '用户ID';

COMMENT ON COLUMN "public"."basic_user"."username" IS '用户名/员工工号';

COMMENT ON COLUMN "public"."basic_user"."password" IS '密码（哈希）';

COMMENT ON COLUMN "public"."basic_user"."salt" IS '密码盐值';

COMMENT ON COLUMN "public"."basic_user"."wechat_id" IS '微信号';

COMMENT ON COLUMN "public"."basic_user"."avatar" IS '头像URL';

COMMENT ON COLUMN "public"."basic_user"."nick_name" IS '昵称';

COMMENT ON COLUMN "public"."basic_user"."game_nickname" IS '游戏昵称（注册时从游戏账号获取）';

COMMENT ON COLUMN "public"."basic_user"."introduction" IS '个人介绍';

COMMENT ON COLUMN "public"."basic_user"."role" IS '用户角色: super_admin(超级管理员), store_admin(店铺管理员), user(普通用户)';

COMMENT ON COLUMN "public"."basic_user"."pinyin_code" IS '姓名全拼';

COMMENT ON COLUMN "public"."basic_user"."first_letter" IS '姓名首字母';

COMMENT ON COLUMN "public"."basic_user"."last_login_at" IS '最后登录时间';

COMMENT ON COLUMN "public"."basic_user"."created_at" IS '创建时间';

COMMENT ON COLUMN "public"."basic_user"."updated_at" IS '更新时间';

COMMENT ON COLUMN "public"."basic_user"."deleted_at" IS '删除时间';

COMMENT ON COLUMN "public"."basic_user"."is_del" IS '软删除标记: 0=未删除, 1=已删除';

COMMENT ON TABLE "public"."basic_user" IS '基础用户表';
CREATE TABLE "public"."basic_user_role_rel" (
                                                "user_id" int4 NOT NULL,
                                                "role_id" int4 NOT NULL,
                                                CONSTRAINT "basic_user_role_rel_pkey" PRIMARY KEY ("user_id", "role_id")
)
;

ALTER TABLE "public"."basic_user_role_rel"
    OWNER TO "B022MC";

COMMENT ON COLUMN "public"."basic_user_role_rel"."user_id" IS '用户ID';

COMMENT ON COLUMN "public"."basic_user_role_rel"."role_id" IS '角色ID';

COMMENT ON TABLE "public"."basic_user_role_rel" IS '用户角色关联表';