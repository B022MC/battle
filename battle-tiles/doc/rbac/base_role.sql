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

INSERT INTO "public"."basic_role" ("id", "code", "name", "parent_id", "remark", "created_at", "created_user", "updated_at", "updated_user", "first_letter", "pinyin_code", "enable", "is_deleted") VALUES (2, 'shop_admin', '店铺管理员', -1, NULL, '2025-09-08 15:35:27.037738+00', NULL, '2025-09-08 15:35:27.037738+00', NULL, '', '', 't', 'f');
INSERT INTO "public"."basic_role" ("id", "code", "name", "parent_id", "remark", "created_at", "created_user", "updated_at", "updated_user", "first_letter", "pinyin_code", "enable", "is_deleted") VALUES (1, 'super_admin', '超级管理员', -1, NULL, '2025-09-08 15:35:27.037738+00', NULL, '2025-09-08 15:35:27.037738+00', NULL, '', '', 't', 'f');
INSERT INTO "public"."basic_role" ("id", "code", "name", "parent_id", "remark", "created_at", "created_user", "updated_at", "updated_user", "first_letter", "pinyin_code", "enable", "is_deleted") VALUES (3, 'user', '普通用户', -1, NULL, '2025-10-23 08:08:59.184065+00', NULL, '2025-10-23 08:08:59.184065+00', NULL, '', '', 't', 'f');