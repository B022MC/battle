-- =====================================================
-- 圈子系统迁移脚本 (Group System Migration)
-- =====================================================
-- 说明：重新设计成员和圈子管理系统
-- 1. 简化成员概念：成员就是系统用户
-- 2. 引入圈子概念：每个店铺管理员对应一个圈子
-- 3. 圈子成员关系：用户可以加入多个圈子

-- =====================================================
-- 1. 创建店铺圈子表 (Shop Group Table)
-- =====================================================
CREATE TABLE IF NOT EXISTS game_shop_group (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    group_name VARCHAR(64) NOT NULL,
    admin_user_id INTEGER NOT NULL,
    description TEXT DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_shop_group IS '店铺圈子表（每个店铺管理员对应一个圈子）';
COMMENT ON COLUMN game_shop_group.id IS '圈子ID';
COMMENT ON COLUMN game_shop_group.house_gid IS '店铺GID';
COMMENT ON COLUMN game_shop_group.group_name IS '圈子名称';
COMMENT ON COLUMN game_shop_group.admin_user_id IS '圈主用户ID（店铺管理员）';
COMMENT ON COLUMN game_shop_group.description IS '圈子描述';
COMMENT ON COLUMN game_shop_group.is_active IS '是否激活';
COMMENT ON COLUMN game_shop_group.created_at IS '创建时间';
COMMENT ON COLUMN game_shop_group.updated_at IS '更新时间';

-- 唯一索引：一个店铺下，一个管理员只能有一个圈子
CREATE UNIQUE INDEX uk_shop_group_house_admin ON game_shop_group(house_gid, admin_user_id) WHERE is_active = TRUE;
CREATE INDEX idx_shop_group_house ON game_shop_group(house_gid);
CREATE INDEX idx_shop_group_admin ON game_shop_group(admin_user_id);

-- =====================================================
-- 2. 创建圈子成员关系表 (Group Member Relation Table)
-- =====================================================
CREATE TABLE IF NOT EXISTS game_shop_group_member (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_shop_group_member IS '圈子成员关系表（用户可以加入多个圈子）';
COMMENT ON COLUMN game_shop_group_member.id IS '关系ID';
COMMENT ON COLUMN game_shop_group_member.group_id IS '圈子ID';
COMMENT ON COLUMN game_shop_group_member.user_id IS '用户ID';
COMMENT ON COLUMN game_shop_group_member.joined_at IS '加入时间';
COMMENT ON COLUMN game_shop_group_member.created_at IS '创建时间';

-- 唯一索引：一个用户在一个圈子中只能有一条记录
CREATE UNIQUE INDEX uk_group_member_group_user ON game_shop_group_member(group_id, user_id);
CREATE INDEX idx_group_member_group ON game_shop_group_member(group_id);
CREATE INDEX idx_group_member_user ON game_shop_group_member(user_id);

-- =====================================================
-- 3. 更新触发器 (Update Triggers)
-- =====================================================
CREATE TRIGGER update_game_shop_group_updated_at BEFORE UPDATE ON game_shop_group
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 4. 数据迁移说明 (Data Migration Notes)
-- =====================================================
-- 如果需要从旧的 game_member 表迁移数据到新的圈子系统：
-- 
-- 步骤1：为每个店铺管理员创建圈子
-- INSERT INTO game_shop_group (house_gid, group_name, admin_user_id)
-- SELECT DISTINCT 
--     sa.house_gid,
--     COALESCE(u.username, '默认圈子') as group_name,
--     sa.user_id
-- FROM game_shop_admin sa
-- JOIN basic_user u ON u.id = sa.user_id
-- WHERE sa.role = 'admin';
--
-- 步骤2：将 game_member 中的成员关系迁移到圈子成员表
-- （需要根据实际业务逻辑决定如何映射 group_name 到 group_id）

-- =====================================================
-- 5. 示例查询 (Example Queries)
-- =====================================================

-- 查询某个店铺的所有圈子
-- SELECT * FROM game_shop_group WHERE house_gid = 60870 AND is_active = TRUE;

-- 查询某个圈子的所有成员
-- SELECT u.* 
-- FROM game_shop_group_member gm
-- JOIN basic_user u ON u.id = gm.user_id
-- WHERE gm.group_id = 1;

-- 查询某个用户加入的所有圈子
-- SELECT g.* 
-- FROM game_shop_group_member gm
-- JOIN game_shop_group g ON g.id = gm.group_id
-- WHERE gm.user_id = 1 AND g.is_active = TRUE;

-- 查询某个店铺的所有成员（去重）
-- SELECT DISTINCT u.*
-- FROM game_shop_group g
-- JOIN game_shop_group_member gm ON gm.group_id = g.id
-- JOIN basic_user u ON u.id = gm.user_id
-- WHERE g.house_gid = 60870 AND g.is_active = TRUE;

