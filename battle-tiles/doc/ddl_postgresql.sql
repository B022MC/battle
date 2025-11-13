-- =====================================================
-- Battle Tiles 数据库 DDL (PostgreSQL)
-- 创建日期: 2025-11-12
-- 说明: 完整的数据库表结构定义，包含基础用户、游戏账号、店铺管理等模块
-- =====================================================

-- =====================================================
-- 基础用户模块 (Basic User Module)
-- =====================================================

-- 基础用户表
CREATE TABLE IF NOT EXISTS basic_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255),
    salt VARCHAR(50),
    wechat_id VARCHAR(64),
    avatar VARCHAR(255) NOT NULL DEFAULT '',
    nick_name VARCHAR(50) NOT NULL DEFAULT '',
    game_nickname VARCHAR(64),
    introduction TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    pinyin_code VARCHAR(100),
    first_letter VARCHAR(50),
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_del SMALLINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE basic_user IS '基础用户表';
COMMENT ON COLUMN basic_user.id IS '用户ID';
COMMENT ON COLUMN basic_user.username IS '用户名/员工工号';
COMMENT ON COLUMN basic_user.password IS '密码（哈希）';
COMMENT ON COLUMN basic_user.salt IS '密码盐值';
COMMENT ON COLUMN basic_user.wechat_id IS '微信号';
COMMENT ON COLUMN basic_user.avatar IS '头像URL';
COMMENT ON COLUMN basic_user.nick_name IS '昵称';
COMMENT ON COLUMN basic_user.game_nickname IS '游戏昵称（注册时从游戏账号获取）';
COMMENT ON COLUMN basic_user.introduction IS '个人介绍';
COMMENT ON COLUMN basic_user.role IS '用户角色: super_admin(超级管理员), store_admin(店铺管理员), user(普通用户)';
COMMENT ON COLUMN basic_user.pinyin_code IS '姓名全拼';
COMMENT ON COLUMN basic_user.first_letter IS '姓名首字母';
COMMENT ON COLUMN basic_user.last_login_at IS '最后登录时间';
COMMENT ON COLUMN basic_user.created_at IS '创建时间';
COMMENT ON COLUMN basic_user.updated_at IS '更新时间';
COMMENT ON COLUMN basic_user.deleted_at IS '删除时间';
COMMENT ON COLUMN basic_user.is_del IS '软删除标记: 0=未删除, 1=已删除';

CREATE UNIQUE INDEX uk_basic_user_username ON basic_user(username);
CREATE INDEX idx_basic_user_role ON basic_user(role);

-- 基础角色表
CREATE TABLE IF NOT EXISTS basic_role (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    parent_id INTEGER NOT NULL DEFAULT -1,
    remark TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_user INTEGER,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_user INTEGER,
    first_letter VARCHAR(50),
    pinyin_code VARCHAR(100),
    enable BOOLEAN NOT NULL DEFAULT TRUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE basic_role IS '基础角色表';
COMMENT ON COLUMN basic_role.id IS '角色ID';
COMMENT ON COLUMN basic_role.code IS '角色编码';
COMMENT ON COLUMN basic_role.name IS '角色名称';
COMMENT ON COLUMN basic_role.parent_id IS '父级角色ID，默认为-1表示顶级角色';
COMMENT ON COLUMN basic_role.remark IS '备注';
COMMENT ON COLUMN basic_role.created_at IS '创建时间';
COMMENT ON COLUMN basic_role.created_user IS '创建人用户ID';
COMMENT ON COLUMN basic_role.updated_at IS '更新时间';
COMMENT ON COLUMN basic_role.updated_user IS '更新人用户ID';
COMMENT ON COLUMN basic_role.first_letter IS '名称首字母';
COMMENT ON COLUMN basic_role.pinyin_code IS '名称全拼';
COMMENT ON COLUMN basic_role.enable IS '是否启用';
COMMENT ON COLUMN basic_role.is_deleted IS '是否删除';

CREATE UNIQUE INDEX uk_basic_role_code ON basic_role(code) WHERE is_deleted = FALSE;
CREATE INDEX idx_basic_role_enable ON basic_role(enable);
CREATE INDEX idx_basic_role_is_deleted ON basic_role(is_deleted);

-- 基础菜单表
CREATE TABLE IF NOT EXISTS basic_menu (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER NOT NULL DEFAULT -1,
    menu_type INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    component VARCHAR(255) NOT NULL,
    rank VARCHAR(255),
    redirect VARCHAR(255) NOT NULL,
    icon VARCHAR(255) NOT NULL,
    extra_icon VARCHAR(255) NOT NULL,
    enter_transition VARCHAR(255) NOT NULL,
    leave_transition VARCHAR(255) NOT NULL,
    active_path VARCHAR(255) NOT NULL,
    auths VARCHAR(255) NOT NULL,
    frame_src VARCHAR(255) NOT NULL,
    frame_loading BOOLEAN NOT NULL DEFAULT FALSE,
    keep_alive BOOLEAN NOT NULL DEFAULT FALSE,
    hidden_tag BOOLEAN NOT NULL DEFAULT FALSE,
    fixed_tag BOOLEAN NOT NULL DEFAULT FALSE,
    show_link BOOLEAN NOT NULL DEFAULT TRUE,
    show_parent BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_del SMALLINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE basic_menu IS '基础菜单表';
COMMENT ON COLUMN basic_menu.id IS '菜单ID';
COMMENT ON COLUMN basic_menu.parent_id IS '父级菜单ID，默认为-1表示顶级菜单';
COMMENT ON COLUMN basic_menu.menu_type IS '菜单类型';
COMMENT ON COLUMN basic_menu.title IS '菜单标题';
COMMENT ON COLUMN basic_menu.name IS '菜单名称';
COMMENT ON COLUMN basic_menu.path IS '路由路径';
COMMENT ON COLUMN basic_menu.component IS '组件路径';
COMMENT ON COLUMN basic_menu.rank IS '排序';
COMMENT ON COLUMN basic_menu.redirect IS '重定向路径';
COMMENT ON COLUMN basic_menu.icon IS '图标';
COMMENT ON COLUMN basic_menu.extra_icon IS '额外图标';
COMMENT ON COLUMN basic_menu.enter_transition IS '进入动画';
COMMENT ON COLUMN basic_menu.leave_transition IS '离开动画';
COMMENT ON COLUMN basic_menu.active_path IS '激活路径';
COMMENT ON COLUMN basic_menu.auths IS '权限标识';
COMMENT ON COLUMN basic_menu.frame_src IS '内嵌iframe地址';
COMMENT ON COLUMN basic_menu.frame_loading IS '是否显示加载动画';
COMMENT ON COLUMN basic_menu.keep_alive IS '是否缓存页面';
COMMENT ON COLUMN basic_menu.hidden_tag IS '是否隐藏标签';
COMMENT ON COLUMN basic_menu.fixed_tag IS '是否固定标签';
COMMENT ON COLUMN basic_menu.show_link IS '是否显示链接';
COMMENT ON COLUMN basic_menu.show_parent IS '是否显示父级';

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS basic_role_menu_rel (
    role_id INTEGER NOT NULL,
    menu_id INTEGER NOT NULL,
    PRIMARY KEY (role_id, menu_id)
);

COMMENT ON TABLE basic_role_menu_rel IS '角色菜单关联表';
COMMENT ON COLUMN basic_role_menu_rel.role_id IS '角色ID';
COMMENT ON COLUMN basic_role_menu_rel.menu_id IS '菜单ID';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS basic_user_role_rel (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

COMMENT ON TABLE basic_user_role_rel IS '用户角色关联表';
COMMENT ON COLUMN basic_user_role_rel.user_id IS '用户ID';
COMMENT ON COLUMN basic_user_role_rel.role_id IS '角色ID';

-- =====================================================
-- 云平台模块 (Cloud Platform Module)
-- =====================================================

-- 平台表
CREATE TABLE IF NOT EXISTS base_platform (
    platform VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    db_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE base_platform IS '云平台表';
COMMENT ON COLUMN base_platform.platform IS '平台标识';
COMMENT ON COLUMN base_platform.name IS '平台名称';
COMMENT ON COLUMN base_platform.db_name IS '数据库名称';
COMMENT ON COLUMN base_platform.created_at IS '创建时间';

-- =====================================================
-- 游戏账号模块 (Game Account Module)
-- =====================================================

-- 游戏账号表（用户绑定的游戏账号）
CREATE TABLE IF NOT EXISTS game_account (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    account VARCHAR(64) NOT NULL,
    pwd_md5 VARCHAR(64) NOT NULL,
    nickname VARCHAR(64) NOT NULL DEFAULT '',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    status INTEGER NOT NULL DEFAULT 1,
    last_login_at TIMESTAMP WITH TIME ZONE,
    login_mode VARCHAR(10) NOT NULL DEFAULT 'account',
    ctrl_account_id INTEGER,
    game_user_id VARCHAR(32) DEFAULT '',
    verified_at TIMESTAMP WITH TIME ZONE,
    verification_status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_del SMALLINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE game_account IS '游戏账号表（用户绑定的游戏账号）';
COMMENT ON COLUMN game_account.id IS '游戏账号ID';
COMMENT ON COLUMN game_account.user_id IS '关联的用户ID';
COMMENT ON COLUMN game_account.account IS '游戏账号';
COMMENT ON COLUMN game_account.pwd_md5 IS '游戏密码MD5';
COMMENT ON COLUMN game_account.nickname IS '游戏昵称';
COMMENT ON COLUMN game_account.is_default IS '是否为默认账号';
COMMENT ON COLUMN game_account.status IS '账号状态: 1=启用, 0=禁用';
COMMENT ON COLUMN game_account.last_login_at IS '最后登录时间';
COMMENT ON COLUMN game_account.login_mode IS '登录方式: account=账号密码, mobile=手机号';
COMMENT ON COLUMN game_account.ctrl_account_id IS '关联的中控账号ID';
COMMENT ON COLUMN game_account.game_user_id IS '游戏服务器返回的用户ID';
COMMENT ON COLUMN game_account.verified_at IS '验证时间';
COMMENT ON COLUMN game_account.verification_status IS '验证状态: pending=待验证, verified=已验证, failed=验证失败';
COMMENT ON COLUMN game_account.created_at IS '创建时间';
COMMENT ON COLUMN game_account.updated_at IS '更新时间';
COMMENT ON COLUMN game_account.deleted_at IS '删除时间';
COMMENT ON COLUMN game_account.is_del IS '软删除标记';

CREATE INDEX idx_game_account_user_id ON game_account(user_id);
CREATE INDEX idx_game_account_game_user_id ON game_account(game_user_id);
CREATE INDEX idx_game_account_verification ON game_account(verification_status);

-- 中控账号表（超级管理员管理的游戏账号）
CREATE TABLE IF NOT EXISTS game_ctrl_account (
    id SERIAL PRIMARY KEY,
    login_mode SMALLINT NOT NULL,
    identifier VARCHAR(64) NOT NULL,
    pwd_md5 VARCHAR(64) NOT NULL,
    game_user_id VARCHAR(32) NOT NULL DEFAULT '',
    game_id VARCHAR(32) NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 1,
    last_verify_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

COMMENT ON TABLE game_ctrl_account IS '中控账号表（超级管理员管理的游戏账号）';
COMMENT ON COLUMN game_ctrl_account.id IS '中控账号ID';
COMMENT ON COLUMN game_ctrl_account.login_mode IS '登录方式: 1=账号密码, 2=手机号';
COMMENT ON COLUMN game_ctrl_account.identifier IS '账号标识（账号或手机号）';
COMMENT ON COLUMN game_ctrl_account.pwd_md5 IS '密码MD5';
COMMENT ON COLUMN game_ctrl_account.game_user_id IS '游戏服务器返回的用户ID';
COMMENT ON COLUMN game_ctrl_account.game_id IS '游戏ID';
COMMENT ON COLUMN game_ctrl_account.status IS '账号状态: 1=启用, 0=禁用';
COMMENT ON COLUMN game_ctrl_account.last_verify_at IS '最后验证时间';
COMMENT ON COLUMN game_ctrl_account.created_at IS '创建时间';
COMMENT ON COLUMN game_ctrl_account.updated_at IS '更新时间';
COMMENT ON COLUMN game_ctrl_account.deleted_at IS '删除时间';

CREATE UNIQUE INDEX uk_ctrl_account_login_identifier ON game_ctrl_account(login_mode, identifier);
CREATE INDEX idx_ctrl_account_identifier ON game_ctrl_account(identifier);
CREATE INDEX idx_ctrl_account_status ON game_ctrl_account(status);

-- 中控账号店铺绑定表
CREATE TABLE IF NOT EXISTS game_account_house (
    id SERIAL PRIMARY KEY,
    game_account_id INTEGER NOT NULL,
    house_gid INTEGER NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    status INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_account_house IS '中控账号店铺绑定表';
COMMENT ON COLUMN game_account_house.id IS '绑定ID';
COMMENT ON COLUMN game_account_house.game_account_id IS '中控账号ID';
COMMENT ON COLUMN game_account_house.house_gid IS '店铺GID（游戏茶馆号）';
COMMENT ON COLUMN game_account_house.is_default IS '是否为默认店铺';
COMMENT ON COLUMN game_account_house.status IS '绑定状态: 1=启用, 0=禁用';
COMMENT ON COLUMN game_account_house.created_at IS '创建时间';
COMMENT ON COLUMN game_account_house.updated_at IS '更新时间';

CREATE INDEX idx_account_house_game_account ON game_account_house(game_account_id);
CREATE INDEX idx_account_house_house_gid ON game_account_house(house_gid);
CREATE UNIQUE INDEX uk_account_house_unique ON game_account_house(game_account_id, house_gid);

-- 游戏账号店铺绑定表（用户游戏账号与店铺的绑定关系）
CREATE TABLE IF NOT EXISTS game_account_store_binding (
    id SERIAL PRIMARY KEY,
    game_account_id INTEGER NOT NULL,
    house_gid INTEGER NOT NULL,
    bound_by_user_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_account_store_binding IS '游戏账号店铺绑定表（业务规则：一个游戏账号只能绑定一个店铺）';
COMMENT ON COLUMN game_account_store_binding.id IS '绑定ID';
COMMENT ON COLUMN game_account_store_binding.game_account_id IS '游戏账号ID';
COMMENT ON COLUMN game_account_store_binding.house_gid IS '店铺GID';
COMMENT ON COLUMN game_account_store_binding.bound_by_user_id IS '绑定操作人用户ID';
COMMENT ON COLUMN game_account_store_binding.status IS '绑定状态: active=激活, inactive=未激活';
COMMENT ON COLUMN game_account_store_binding.created_at IS '创建时间';
COMMENT ON COLUMN game_account_store_binding.updated_at IS '更新时间';

CREATE UNIQUE INDEX uk_game_account_house ON game_account_store_binding(game_account_id, house_gid);
CREATE INDEX idx_gasb_house_gid ON game_account_store_binding(house_gid);
CREATE INDEX idx_gasb_bound_by_user ON game_account_store_binding(bound_by_user_id);
CREATE INDEX idx_gasb_status ON game_account_store_binding(status);

-- =====================================================
-- 游戏会话模块 (Game Session Module)
-- =====================================================

-- 游戏会话表
CREATE TABLE IF NOT EXISTS game_session (
    id SERIAL PRIMARY KEY,
    game_ctrl_account_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    house_gid INTEGER NOT NULL,
    state VARCHAR(20) NOT NULL,
    device_ip VARCHAR(64) NOT NULL DEFAULT '',
    error_msg VARCHAR(255) NOT NULL DEFAULT '',
    end_at TIMESTAMP WITH TIME ZONE,
    auto_sync_enabled BOOLEAN DEFAULT TRUE,
    last_sync_at TIMESTAMP WITH TIME ZONE,
    sync_status VARCHAR(20) DEFAULT 'idle',
    game_account_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_del SMALLINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE game_session IS '游戏会话表（记录中控账号的登录会话）';
COMMENT ON COLUMN game_session.id IS '会话ID';
COMMENT ON COLUMN game_session.game_ctrl_account_id IS '中控账号ID';
COMMENT ON COLUMN game_session.user_id IS '创建会话的用户ID';
COMMENT ON COLUMN game_session.house_gid IS '店铺GID';
COMMENT ON COLUMN game_session.state IS '会话状态: active=活跃, inactive=未活跃, error=错误';
COMMENT ON COLUMN game_session.device_ip IS '设备IP地址';
COMMENT ON COLUMN game_session.error_msg IS '错误信息';
COMMENT ON COLUMN game_session.end_at IS '会话结束时间';
COMMENT ON COLUMN game_session.auto_sync_enabled IS '是否启用自动同步';
COMMENT ON COLUMN game_session.last_sync_at IS '最后同步时间';
COMMENT ON COLUMN game_session.sync_status IS '同步状态: idle=空闲, syncing=同步中, error=错误';
COMMENT ON COLUMN game_session.game_account_id IS '关联的游戏账号ID';
COMMENT ON COLUMN game_session.created_at IS '创建时间';
COMMENT ON COLUMN game_session.updated_at IS '更新时间';
COMMENT ON COLUMN game_session.deleted_at IS '删除时间';
COMMENT ON COLUMN game_session.is_del IS '软删除标记';

CREATE INDEX idx_session_state ON game_session(state);
CREATE INDEX idx_session_sync_status ON game_session(sync_status);
CREATE INDEX idx_session_game_account ON game_session(game_account_id);
CREATE INDEX idx_session_house_gid ON game_session(house_gid);
CREATE INDEX idx_session_ctrl_account ON game_session(game_ctrl_account_id);

-- 游戏同步日志表
CREATE TABLE IF NOT EXISTS game_sync_log (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL,
    sync_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    records_synced INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE
);

COMMENT ON TABLE game_sync_log IS '游戏同步日志表（记录数据同步操作）';
COMMENT ON COLUMN game_sync_log.id IS '日志ID';
COMMENT ON COLUMN game_sync_log.session_id IS '会话ID';
COMMENT ON COLUMN game_sync_log.sync_type IS '同步类型: battle_record=战绩, member_list=成员列表, wallet_update=钱包更新, room_list=房间列表, group_member=圈成员';
COMMENT ON COLUMN game_sync_log.status IS '同步状态: success=成功, failed=失败, partial=部分成功';
COMMENT ON COLUMN game_sync_log.records_synced IS '同步记录数';
COMMENT ON COLUMN game_sync_log.error_message IS '错误信息';
COMMENT ON COLUMN game_sync_log.started_at IS '开始时间';
COMMENT ON COLUMN game_sync_log.completed_at IS '完成时间';

CREATE INDEX idx_sync_log_session_started ON game_sync_log(session_id, started_at);
CREATE INDEX idx_sync_log_type ON game_sync_log(sync_type);
CREATE INDEX idx_sync_log_status ON game_sync_log(status);

-- =====================================================
-- 店铺管理模块 (Shop Management Module)
-- =====================================================

-- 店铺管理员表
CREATE TABLE IF NOT EXISTS game_shop_admin (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL,
    game_account_id INTEGER,
    is_exclusive BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

COMMENT ON TABLE game_shop_admin IS '店铺管理员表';
COMMENT ON COLUMN game_shop_admin.id IS '管理员ID';
COMMENT ON COLUMN game_shop_admin.house_gid IS '店铺GID';
COMMENT ON COLUMN game_shop_admin.user_id IS '用户ID';
COMMENT ON COLUMN game_shop_admin.role IS '角色: admin=管理员, operator=操作员';
COMMENT ON COLUMN game_shop_admin.game_account_id IS '关联的游戏账号ID';
COMMENT ON COLUMN game_shop_admin.is_exclusive IS '是否独占（店铺管理员同时只能管理一个店铺）';
COMMENT ON COLUMN game_shop_admin.created_at IS '创建时间';
COMMENT ON COLUMN game_shop_admin.updated_at IS '更新时间';
COMMENT ON COLUMN game_shop_admin.deleted_at IS '删除时间';

CREATE INDEX idx_shop_admin_house_gid ON game_shop_admin(house_gid);
CREATE INDEX idx_shop_admin_user_id ON game_shop_admin(user_id);
CREATE INDEX idx_shop_admin_game_account ON game_shop_admin(game_account_id);
CREATE INDEX idx_shop_admin_exclusive ON game_shop_admin(user_id, is_exclusive);

-- 店铺设置表
CREATE TABLE IF NOT EXISTS game_house_settings (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    fees_json TEXT NOT NULL DEFAULT '',
    share_fee BOOLEAN NOT NULL DEFAULT FALSE,
    push_credit INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by INTEGER NOT NULL DEFAULT 0
);

COMMENT ON TABLE game_house_settings IS '店铺设置表';
COMMENT ON COLUMN game_house_settings.id IS '设置ID';
COMMENT ON COLUMN game_house_settings.house_gid IS '店铺GID';
COMMENT ON COLUMN game_house_settings.fees_json IS '运费规则JSON';
COMMENT ON COLUMN game_house_settings.share_fee IS '分运开关';
COMMENT ON COLUMN game_house_settings.push_credit IS '推送额度（单位：分）';
COMMENT ON COLUMN game_house_settings.updated_at IS '更新时间';
COMMENT ON COLUMN game_house_settings.updated_by IS '操作人用户ID';

CREATE UNIQUE INDEX uk_house ON game_house_settings(house_gid);

-- =====================================================
-- 游戏成员模块 (Game Member Module)
-- =====================================================

-- 游戏成员表
CREATE TABLE IF NOT EXISTS game_member (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    game_name VARCHAR(64) NOT NULL DEFAULT '',
    group_name VARCHAR(64) NOT NULL DEFAULT '',
    balance INTEGER NOT NULL DEFAULT 0,
    credit INTEGER NOT NULL DEFAULT 0,
    forbid BOOLEAN NOT NULL DEFAULT FALSE,
    recommender VARCHAR(64) DEFAULT '',
    use_multi_gids BOOLEAN DEFAULT FALSE,
    active_gid INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_member IS '游戏成员表（店铺内的玩家）';
COMMENT ON COLUMN game_member.id IS '成员ID';
COMMENT ON COLUMN game_member.house_gid IS '店铺GID';
COMMENT ON COLUMN game_member.game_id IS '游戏ID';
COMMENT ON COLUMN game_member.game_name IS '游戏昵称';
COMMENT ON COLUMN game_member.group_name IS '圈名';
COMMENT ON COLUMN game_member.balance IS '余额（单位：分）';
COMMENT ON COLUMN game_member.credit IS '信用额度（单位：分）';
COMMENT ON COLUMN game_member.forbid IS '是否禁用';
COMMENT ON COLUMN game_member.recommender IS '推荐人';
COMMENT ON COLUMN game_member.use_multi_gids IS '是否允许多开';
COMMENT ON COLUMN game_member.active_gid IS '当前活跃游戏ID';
COMMENT ON COLUMN game_member.created_at IS '创建时间';
COMMENT ON COLUMN game_member.updated_at IS '更新时间';

CREATE UNIQUE INDEX uk_game_member_house_game ON game_member(house_gid, game_id);
CREATE INDEX idx_game_member_game_id ON game_member(game_id);
CREATE INDEX idx_game_member_house_group ON game_member(house_gid, group_name);
CREATE INDEX idx_game_member_balance ON game_member(balance);
CREATE INDEX idx_game_member_forbid ON game_member(forbid);

-- 游戏成员钱包表
CREATE TABLE IF NOT EXISTS game_member_wallet (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    balance INTEGER NOT NULL DEFAULT 0,
    forbid BOOLEAN NOT NULL DEFAULT FALSE,
    limit_min INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by INTEGER NOT NULL DEFAULT 0
);

COMMENT ON TABLE game_member_wallet IS '游戏成员钱包表';
COMMENT ON COLUMN game_member_wallet.id IS '钱包ID';
COMMENT ON COLUMN game_member_wallet.house_gid IS '店铺GID';
COMMENT ON COLUMN game_member_wallet.member_id IS '成员ID';
COMMENT ON COLUMN game_member_wallet.balance IS '余额（单位：分）';
COMMENT ON COLUMN game_member_wallet.forbid IS '是否禁用';
COMMENT ON COLUMN game_member_wallet.limit_min IS '最低限额（单位：分）';
COMMENT ON COLUMN game_member_wallet.updated_at IS '更新时间';
COMMENT ON COLUMN game_member_wallet.updated_by IS '操作人用户ID';

CREATE INDEX idx_member_wallet_house ON game_member_wallet(house_gid);
CREATE INDEX idx_member_wallet_member ON game_member_wallet(member_id);

-- 游戏成员规则表
CREATE TABLE IF NOT EXISTS game_member_rule (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    vip BOOLEAN NOT NULL DEFAULT FALSE,
    multi_gids BOOLEAN NOT NULL DEFAULT FALSE,
    temp_release INTEGER NOT NULL DEFAULT 0,
    expire_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by INTEGER NOT NULL DEFAULT 0
);

COMMENT ON TABLE game_member_rule IS '游戏成员规则表（VIP、多号、临时解禁等）';
COMMENT ON COLUMN game_member_rule.id IS '规则ID';
COMMENT ON COLUMN game_member_rule.house_gid IS '店铺GID';
COMMENT ON COLUMN game_member_rule.member_id IS '成员ID';
COMMENT ON COLUMN game_member_rule.vip IS '是否VIP';
COMMENT ON COLUMN game_member_rule.multi_gids IS '是否允许多号';
COMMENT ON COLUMN game_member_rule.temp_release IS '临时解禁上限（单位：分），0表示无限制';
COMMENT ON COLUMN game_member_rule.expire_at IS '规则过期时间';
COMMENT ON COLUMN game_member_rule.updated_at IS '更新时间';
COMMENT ON COLUMN game_member_rule.updated_by IS '操作人用户ID';

CREATE INDEX idx_member_rule_house ON game_member_rule(house_gid);
CREATE INDEX idx_member_rule_member ON game_member_rule(member_id);

-- =====================================================
-- 游戏战绩模块 (Game Battle Record Module)
-- =====================================================

-- 游戏战绩表
CREATE TABLE IF NOT EXISTS game_battle_record (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    room_uid INTEGER NOT NULL,
    kind_id INTEGER NOT NULL,
    base_score INTEGER NOT NULL,
    battle_at TIMESTAMP WITH TIME ZONE NOT NULL,
    players_json TEXT NOT NULL,
    player_id INTEGER,
    player_game_id INTEGER,
    player_game_name VARCHAR(64) DEFAULT '',
    group_name VARCHAR(64) DEFAULT '',
    score INTEGER DEFAULT 0,
    fee INTEGER DEFAULT 0,
    factor DECIMAL(10,4) DEFAULT 1.0000,
    player_balance INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_battle_record IS '游戏战绩表（本地战绩快照，一行代表一局一桌）';
COMMENT ON COLUMN game_battle_record.id IS '战绩ID';
COMMENT ON COLUMN game_battle_record.house_gid IS '店铺GID';
COMMENT ON COLUMN game_battle_record.group_id IS '圈ID';
COMMENT ON COLUMN game_battle_record.room_uid IS '房间唯一ID';
COMMENT ON COLUMN game_battle_record.kind_id IS '游戏类型ID';
COMMENT ON COLUMN game_battle_record.base_score IS '底分';
COMMENT ON COLUMN game_battle_record.battle_at IS '对战时间';
COMMENT ON COLUMN game_battle_record.players_json IS '玩家列表JSON';
COMMENT ON COLUMN game_battle_record.player_id IS '玩家ID（用于按玩家查询）';
COMMENT ON COLUMN game_battle_record.player_game_id IS '玩家游戏ID';
COMMENT ON COLUMN game_battle_record.player_game_name IS '玩家游戏昵称';
COMMENT ON COLUMN game_battle_record.group_name IS '圈名';
COMMENT ON COLUMN game_battle_record.score IS '得分';
COMMENT ON COLUMN game_battle_record.fee IS '服务费（单位：分）';
COMMENT ON COLUMN game_battle_record.factor IS '结算比例';
COMMENT ON COLUMN game_battle_record.player_balance IS '玩家余额（单位：分）';
COMMENT ON COLUMN game_battle_record.created_at IS '创建时间';

CREATE INDEX idx_battle_house_gid ON game_battle_record(house_gid);
CREATE INDEX idx_battle_room_uid ON game_battle_record(room_uid);
CREATE INDEX idx_battle_kind_id ON game_battle_record(kind_id);
CREATE INDEX idx_battle_at ON game_battle_record(battle_at);
CREATE INDEX idx_battle_player_id ON game_battle_record(player_id);
CREATE INDEX idx_battle_player_game_id ON game_battle_record(player_game_id);
CREATE INDEX idx_battle_group ON game_battle_record(house_gid, group_name);

-- =====================================================
-- 充值记录模块 (Recharge Record Module)
-- =====================================================

-- 充值记录表
CREATE TABLE IF NOT EXISTS game_recharge_record (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    group_name VARCHAR(64) NOT NULL DEFAULT '',
    amount INTEGER NOT NULL,
    balance_before INTEGER NOT NULL,
    balance_after INTEGER NOT NULL,
    operator_user_id INTEGER,
    recharged_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_recharge_record IS '充值记录表（上下分记录）';
COMMENT ON COLUMN game_recharge_record.id IS '记录ID';
COMMENT ON COLUMN game_recharge_record.house_gid IS '店铺GID';
COMMENT ON COLUMN game_recharge_record.player_id IS '玩家ID';
COMMENT ON COLUMN game_recharge_record.group_name IS '圈名';
COMMENT ON COLUMN game_recharge_record.amount IS '金额（单位：分），正数=充值，负数=提现';
COMMENT ON COLUMN game_recharge_record.balance_before IS '操作前余额（单位：分）';
COMMENT ON COLUMN game_recharge_record.balance_after IS '操作后余额（单位：分）';
COMMENT ON COLUMN game_recharge_record.operator_user_id IS '操作人用户ID';
COMMENT ON COLUMN game_recharge_record.recharged_at IS '充值时间';
COMMENT ON COLUMN game_recharge_record.created_at IS '创建时间';

CREATE INDEX idx_recharge_house_gid ON game_recharge_record(house_gid);
CREATE INDEX idx_recharge_player ON game_recharge_record(player_id);
CREATE INDEX idx_recharge_recharged_at ON game_recharge_record(recharged_at);
CREATE INDEX idx_recharge_house_group ON game_recharge_record(house_gid, group_name);

-- =====================================================
-- 费用结算模块 (Fee Settlement Module)
-- =====================================================

-- 费用结算表
CREATE TABLE IF NOT EXISTS game_fee_settle (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    play_group VARCHAR(32) NOT NULL DEFAULT '',
    amount INTEGER NOT NULL,
    feed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE game_fee_settle IS '费用结算表（本地费用结算记录）';
COMMENT ON COLUMN game_fee_settle.id IS '结算ID';
COMMENT ON COLUMN game_fee_settle.house_gid IS '店铺GID';
COMMENT ON COLUMN game_fee_settle.play_group IS '圈名';
COMMENT ON COLUMN game_fee_settle.amount IS '金额（单位：分）';
COMMENT ON COLUMN game_fee_settle.feed_at IS '结算时间';
COMMENT ON COLUMN game_fee_settle.created_at IS '创建时间';

CREATE INDEX idx_fee_settle_house ON game_fee_settle(house_gid);
CREATE INDEX idx_fee_settle_feed_at ON game_fee_settle(feed_at);

-- =====================================================
-- 钱包账本模块 (Wallet Ledger Module)
-- =====================================================

-- 钱包账本表（如果存在）
-- 注意：根据实际需求，此表可能需要进一步定义
-- CREATE TABLE IF NOT EXISTS game_wallet_ledger (
--     id SERIAL PRIMARY KEY,
--     house_gid INTEGER NOT NULL,
--     member_id INTEGER NOT NULL,
--     transaction_type VARCHAR(20) NOT NULL,
--     amount INTEGER NOT NULL,
--     balance_before INTEGER NOT NULL,
--     balance_after INTEGER NOT NULL,
--     reference_id INTEGER,
--     reference_type VARCHAR(50),
--     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- );

-- =====================================================
-- 其他辅助表 (Additional Tables)
-- =====================================================

-- 店铺圈管理员表（如果需要）
-- CREATE TABLE IF NOT EXISTS game_shop_group_admin (
--     id SERIAL PRIMARY KEY,
--     house_gid INTEGER NOT NULL,
--     group_name VARCHAR(64) NOT NULL,
--     user_id INTEGER NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- );

-- 用户申请表（如果需要）
-- CREATE TABLE IF NOT EXISTS user_application (
--     id SERIAL PRIMARY KEY,
--     user_id INTEGER NOT NULL,
--     application_type VARCHAR(50) NOT NULL,
--     status VARCHAR(20) NOT NULL DEFAULT 'pending',
--     content TEXT,
--     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- );

-- =====================================================
-- 触发器和函数 (Triggers and Functions)
-- =====================================================

-- 自动更新 updated_at 字段的触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为所有需要的表创建触发器
CREATE TRIGGER update_basic_user_updated_at BEFORE UPDATE ON basic_user
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_basic_menu_updated_at BEFORE UPDATE ON basic_menu
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_account_updated_at BEFORE UPDATE ON game_account
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_ctrl_account_updated_at BEFORE UPDATE ON game_ctrl_account
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_account_house_updated_at BEFORE UPDATE ON game_account_house
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_account_store_binding_updated_at BEFORE UPDATE ON game_account_store_binding
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_session_updated_at BEFORE UPDATE ON game_session
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_shop_admin_updated_at BEFORE UPDATE ON game_shop_admin
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_house_settings_updated_at BEFORE UPDATE ON game_house_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_member_updated_at BEFORE UPDATE ON game_member
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_member_wallet_updated_at BEFORE UPDATE ON game_member_wallet
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_member_rule_updated_at BEFORE UPDATE ON game_member_rule
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 初始数据 (Initial Data)
-- =====================================================

-- 插入默认超级管理员（可选）
-- INSERT INTO basic_user (username, password, salt, role, nick_name, created_at, updated_at)
-- VALUES ('admin', 'hashed_password_here', 'salt_here', 'super_admin', '系统管理员', NOW(), NOW())
-- ON CONFLICT (username) DO NOTHING;

-- =====================================================
-- 权限和安全 (Permissions and Security)
-- =====================================================

-- 创建只读用户（可选）
-- CREATE USER battle_readonly WITH PASSWORD 'your_password_here';
-- GRANT CONNECT ON DATABASE your_database_name TO battle_readonly;
-- GRANT USAGE ON SCHEMA public TO battle_readonly;
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO battle_readonly;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO battle_readonly;

-- 创建读写用户（可选）
-- CREATE USER battle_readwrite WITH PASSWORD 'your_password_here';
-- GRANT CONNECT ON DATABASE your_database_name TO battle_readwrite;
-- GRANT USAGE ON SCHEMA public TO battle_readwrite;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO battle_readwrite;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO battle_readwrite;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO battle_readwrite;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO battle_readwrite;

-- =====================================================
-- 索引优化建议 (Index Optimization Suggestions)
-- =====================================================

-- 根据实际查询模式，可能需要添加以下复合索引：
-- CREATE INDEX idx_game_account_user_status ON game_account(user_id, status) WHERE is_del = 0;
-- CREATE INDEX idx_game_session_house_state ON game_session(house_gid, state) WHERE is_del = 0;
-- CREATE INDEX idx_game_member_house_forbid ON game_member(house_gid, forbid);
-- CREATE INDEX idx_battle_record_house_time ON game_battle_record(house_gid, battle_at DESC);

-- =====================================================
-- 维护脚本 (Maintenance Scripts)
-- =====================================================

-- 清理软删除数据（定期执行）
-- DELETE FROM basic_user WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';
-- DELETE FROM game_account WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';
-- DELETE FROM game_session WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';

-- 分析表以优化查询计划
-- ANALYZE basic_user;
-- ANALYZE game_account;
-- ANALYZE game_session;
-- ANALYZE game_battle_record;
-- ANALYZE game_member;

-- =====================================================
-- 备份建议 (Backup Recommendations)
-- =====================================================

-- 1. 每日全量备份
-- pg_dump -h localhost -U postgres -d battle_tiles -F c -f backup_$(date +%Y%m%d).dump

-- 2. 实时 WAL 归档（用于时间点恢复）
-- 在 postgresql.conf 中配置:
-- wal_level = replica
-- archive_mode = on
-- archive_command = 'cp %p /path/to/archive/%f'

-- 3. 定期测试恢复流程

-- =====================================================
-- 结束 (End of DDL)
-- =====================================================

