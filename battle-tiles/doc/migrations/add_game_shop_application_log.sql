-- ============================================
-- 游戏店铺申请操作日志表
-- 创建时间: 2025-11-20
-- 用途: 记录管理员对玩家申请的处理操作（通过/拒绝）
-- 说明: 申请数据本身存储在 Plaza Session 内存中，此表仅记录操作日志
-- ============================================

CREATE TABLE IF NOT EXISTS game_shop_application_log (
    id BIGSERIAL PRIMARY KEY,
    house_gid INT NOT NULL,
    applier_gid INT NOT NULL,
    applier_gname VARCHAR(100) NOT NULL,
    action VARCHAR(20) NOT NULL,
    admin_user_id INT NOT NULL,
    admin_game_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_app_log_house ON game_shop_application_log(house_gid);
CREATE INDEX IF NOT EXISTS idx_app_log_admin ON game_shop_application_log(admin_user_id);
CREATE INDEX IF NOT EXISTS idx_app_log_created ON game_shop_application_log(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_app_log_action ON game_shop_application_log(action);

-- 注释
COMMENT ON TABLE game_shop_application_log IS '店铺申请操作日志（记录通过/拒绝操作，申请数据存在内存中）';
COMMENT ON COLUMN game_shop_application_log.id IS '主键ID';
COMMENT ON COLUMN game_shop_application_log.house_gid IS '店铺游戏ID';
COMMENT ON COLUMN game_shop_application_log.applier_gid IS '申请人游戏ID';
COMMENT ON COLUMN game_shop_application_log.applier_gname IS '申请人游戏昵称';
COMMENT ON COLUMN game_shop_application_log.action IS '操作类型：approved=通过，rejected=拒绝';
COMMENT ON COLUMN game_shop_application_log.admin_user_id IS '处理的管理员系统用户ID';
COMMENT ON COLUMN game_shop_application_log.admin_game_id IS '管理员游戏ID（可选）';
COMMENT ON COLUMN game_shop_application_log.created_at IS '操作时间';
