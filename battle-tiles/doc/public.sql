create table if not exists basic_user
(
    id            serial
        primary key,
    username      varchar(50)                                                not null,
    password      varchar(255),
    salt          varchar(50),
    wechat_id     varchar(64),
    avatar        varchar(255)             default ''::character varying     not null,
    nick_name     varchar(50)              default ''::character varying     not null,
    game_nickname varchar(64),
    introduction  text,
    role          varchar(20)              default 'user'::character varying not null,
    pinyin_code   varchar(100),
    first_letter  varchar(50),
    last_login_at timestamp with time zone,
    created_at    timestamp with time zone default now()                     not null,
    updated_at    timestamp with time zone default now()                     not null,
    deleted_at    timestamp with time zone,
    is_del        smallint                 default 0                         not null
);

comment on table basic_user is '基础用户表';

comment on column basic_user.id is '用户ID';

comment on column basic_user.username is '用户名/员工工号';

comment on column basic_user.password is '密码（哈希）';

comment on column basic_user.salt is '密码盐值';

comment on column basic_user.wechat_id is '微信号';

comment on column basic_user.avatar is '头像URL';

comment on column basic_user.nick_name is '昵称';

comment on column basic_user.game_nickname is '游戏昵称（注册时从游戏账号获取）';

comment on column basic_user.introduction is '个人介绍';

comment on column basic_user.role is '用户角色: super_admin(超级管理员), store_admin(店铺管理员), user(普通用户)';

comment on column basic_user.pinyin_code is '姓名全拼';

comment on column basic_user.first_letter is '姓名首字母';

comment on column basic_user.last_login_at is '最后登录时间';

comment on column basic_user.created_at is '创建时间';

comment on column basic_user.updated_at is '更新时间';

comment on column basic_user.deleted_at is '删除时间';

comment on column basic_user.is_del is '软删除标记: 0=未删除, 1=已删除';

create unique index if not exists uk_basic_user_username
    on basic_user (username);

create index if not exists idx_basic_user_role
    on basic_user (role);

create table if not exists basic_role
(
    id           serial
        primary key,
    code         varchar(50)                                    not null,
    name         varchar(100)                                   not null,
    parent_id    integer                  default '-1'::integer not null,
    remark       text,
    created_at   timestamp with time zone default now()         not null,
    created_user integer,
    updated_at   timestamp with time zone,
    updated_user integer,
    first_letter varchar(50),
    pinyin_code  varchar(100),
    enable       boolean                  default true          not null,
    is_deleted   boolean                  default false         not null
);

comment on table basic_role is '基础角色表';

comment on column basic_role.id is '角色ID';

comment on column basic_role.code is '角色编码';

comment on column basic_role.name is '角色名称';

comment on column basic_role.parent_id is '父级角色ID，默认为-1表示顶级角色';

comment on column basic_role.remark is '备注';

comment on column basic_role.created_at is '创建时间';

comment on column basic_role.created_user is '创建人用户ID';

comment on column basic_role.updated_at is '更新时间';

comment on column basic_role.updated_user is '更新人用户ID';

comment on column basic_role.first_letter is '名称首字母';

comment on column basic_role.pinyin_code is '名称全拼';

comment on column basic_role.enable is '是否启用';

comment on column basic_role.is_deleted is '是否删除';

create unique index if not exists uk_basic_role_code
    on basic_role (code)
    where (is_deleted = false);

create index if not exists idx_basic_role_enable
    on basic_role (enable);

create index if not exists idx_basic_role_is_deleted
    on basic_role (is_deleted);

create table if not exists basic_user_role_rel
(
    user_id integer not null,
    role_id integer not null,
    primary key (user_id, role_id)
);

comment on table basic_user_role_rel is '用户角色关联表';

comment on column basic_user_role_rel.user_id is '用户ID';

comment on column basic_user_role_rel.role_id is '角色ID';

create table if not exists game_account
(
    id                  serial
        primary key,
    user_id             integer                                                       not null,
    account             varchar(64)                                                   not null,
    pwd_md5             varchar(64)                                                   not null,
    nickname            varchar(64)              default ''::character varying        not null,
    is_default          boolean                  default false                        not null,
    status              integer                  default 1                            not null,
    last_login_at       timestamp with time zone,
    login_mode          varchar(10)              default 'account'::character varying not null,
    ctrl_account_id     integer,
    game_user_id        varchar(32)              default ''::character varying,
    verified_at         timestamp with time zone,
    verification_status varchar(20)              default 'pending'::character varying,
    created_at          timestamp with time zone default now()                        not null,
    updated_at          timestamp with time zone default now()                        not null,
    deleted_at          timestamp with time zone,
    is_del              smallint                 default 0                            not null
);

comment on table game_account is '游戏账号表（用户绑定的游戏账号）';

comment on column game_account.id is '游戏账号ID';

comment on column game_account.user_id is '关联的用户ID';

comment on column game_account.account is '游戏账号';

comment on column game_account.pwd_md5 is '游戏密码MD5';

comment on column game_account.nickname is '游戏昵称';

comment on column game_account.is_default is '是否为默认账号';

comment on column game_account.status is '账号状态: 1=启用, 0=禁用';

comment on column game_account.last_login_at is '最后登录时间';

comment on column game_account.login_mode is '登录方式: account=账号密码, mobile=手机号';

comment on column game_account.ctrl_account_id is '关联的中控账号ID';

comment on column game_account.game_user_id is '游戏服务器返回的用户ID';

comment on column game_account.verified_at is '验证时间';

comment on column game_account.verification_status is '验证状态: pending=待验证, verified=已验证, failed=验证失败';

comment on column game_account.created_at is '创建时间';

comment on column game_account.updated_at is '更新时间';

comment on column game_account.deleted_at is '删除时间';

comment on column game_account.is_del is '软删除标记';

create index if not exists idx_game_account_user_id
    on game_account (user_id);

create index if not exists idx_game_account_game_user_id
    on game_account (game_user_id);

create index if not exists idx_game_account_verification
    on game_account (verification_status);

create table if not exists game_ctrl_account
(
    id             serial
        primary key,
    login_mode     smallint                                               not null,
    identifier     varchar(64)                                            not null,
    pwd_md5        varchar(64)                                            not null,
    game_user_id   varchar(32)              default ''::character varying not null,
    game_id        varchar(32)              default ''::character varying not null,
    status         integer                  default 1                     not null,
    last_verify_at timestamp with time zone,
    created_at     timestamp with time zone default now()                 not null,
    updated_at     timestamp with time zone default now()                 not null,
    deleted_at     timestamp with time zone
);

comment on table game_ctrl_account is '中控账号表（超级管理员管理的游戏账号）';

comment on column game_ctrl_account.id is '中控账号ID';

comment on column game_ctrl_account.login_mode is '登录方式: 1=账号密码, 2=手机号';

comment on column game_ctrl_account.identifier is '账号标识（账号或手机号）';

comment on column game_ctrl_account.pwd_md5 is '密码MD5';

comment on column game_ctrl_account.game_user_id is '游戏服务器返回的用户ID';

comment on column game_ctrl_account.game_id is '游戏ID';

comment on column game_ctrl_account.status is '账号状态: 1=启用, 0=禁用';

comment on column game_ctrl_account.last_verify_at is '最后验证时间';

comment on column game_ctrl_account.created_at is '创建时间';

comment on column game_ctrl_account.updated_at is '更新时间';

comment on column game_ctrl_account.deleted_at is '删除时间';

create unique index if not exists uk_ctrl_account_login_identifier
    on game_ctrl_account (login_mode, identifier);

create index if not exists idx_ctrl_account_identifier
    on game_ctrl_account (identifier);

create index if not exists idx_ctrl_account_status
    on game_ctrl_account (status);

create table if not exists game_account_house
(
    id              serial
        primary key,
    game_account_id integer                                not null,
    house_gid       integer                                not null,
    is_default      boolean                  default false not null,
    status          integer                  default 1     not null,
    created_at      timestamp with time zone default now() not null,
    updated_at      timestamp with time zone default now() not null
);

comment on table game_account_house is '中控账号店铺绑定表';

comment on column game_account_house.id is '绑定ID';

comment on column game_account_house.game_account_id is '中控账号ID';

comment on column game_account_house.house_gid is '店铺GID（游戏茶馆号）';

comment on column game_account_house.is_default is '是否为默认店铺';

comment on column game_account_house.status is '绑定状态: 1=启用, 0=禁用';

comment on column game_account_house.created_at is '创建时间';

comment on column game_account_house.updated_at is '更新时间';

create index if not exists idx_account_house_game_account
    on game_account_house (game_account_id);

create index if not exists idx_account_house_house_gid
    on game_account_house (house_gid);

create unique index if not exists uk_account_house_unique
    on game_account_house (game_account_id, house_gid);

create table if not exists game_account_store_binding
(
    id               serial
        primary key,
    game_account_id  integer                                                      not null,
    house_gid        integer                                                      not null,
    bound_by_user_id integer                                                      not null,
    status           varchar(20)              default 'active'::character varying not null,
    created_at       timestamp with time zone default now()                       not null,
    updated_at       timestamp with time zone default now()                       not null
);

comment on table game_account_store_binding is '游戏账号店铺绑定表（业务规则：一个游戏账号只能绑定一个店铺）';

comment on column game_account_store_binding.id is '绑定ID';

comment on column game_account_store_binding.game_account_id is '游戏账号ID';

comment on column game_account_store_binding.house_gid is '店铺GID';

comment on column game_account_store_binding.bound_by_user_id is '绑定操作人用户ID';

comment on column game_account_store_binding.status is '绑定状态: active=激活, inactive=未激活';

comment on column game_account_store_binding.created_at is '创建时间';

comment on column game_account_store_binding.updated_at is '更新时间';

create unique index if not exists uk_game_account_house
    on game_account_store_binding (game_account_id, house_gid);

create index if not exists idx_gasb_house_gid
    on game_account_store_binding (house_gid);

create index if not exists idx_gasb_bound_by_user
    on game_account_store_binding (bound_by_user_id);

create index if not exists idx_gasb_status
    on game_account_store_binding (status);

create table if not exists game_session
(
    id                   serial
        primary key,
    game_ctrl_account_id integer                                                not null,
    user_id              integer                                                not null,
    house_gid            integer                                                not null,
    state                varchar(20)                                            not null,
    device_ip            varchar(64)              default ''::character varying not null,
    error_msg            varchar(255)             default ''::character varying not null,
    end_at               timestamp with time zone,
    auto_sync_enabled    boolean                  default true,
    last_sync_at         timestamp with time zone,
    sync_status          varchar(20)              default 'idle'::character varying,
    game_account_id      integer,
    created_at           timestamp with time zone default now()                 not null,
    updated_at           timestamp with time zone default now()                 not null,
    deleted_at           timestamp with time zone,
    is_del               smallint                 default 0                     not null
);

comment on table game_session is '游戏会话表（记录中控账号的登录会话）';

comment on column game_session.id is '会话ID';

comment on column game_session.game_ctrl_account_id is '中控账号ID';

comment on column game_session.user_id is '创建会话的用户ID';

comment on column game_session.house_gid is '店铺GID';

comment on column game_session.state is '会话状态: active=活跃, inactive=未活跃, error=错误';

comment on column game_session.device_ip is '设备IP地址';

comment on column game_session.error_msg is '错误信息';

comment on column game_session.end_at is '会话结束时间';

comment on column game_session.auto_sync_enabled is '是否启用自动同步';

comment on column game_session.last_sync_at is '最后同步时间';

comment on column game_session.sync_status is '同步状态: idle=空闲, syncing=同步中, error=错误';

comment on column game_session.game_account_id is '关联的游戏账号ID';

comment on column game_session.created_at is '创建时间';

comment on column game_session.updated_at is '更新时间';

comment on column game_session.deleted_at is '删除时间';

comment on column game_session.is_del is '软删除标记';

create index if not exists idx_session_state
    on game_session (state);

create index if not exists idx_session_sync_status
    on game_session (sync_status);

create index if not exists idx_session_game_account
    on game_session (game_account_id);

create index if not exists idx_session_house_gid
    on game_session (house_gid);

create index if not exists idx_session_ctrl_account
    on game_session (game_ctrl_account_id);

create table if not exists game_sync_log
(
    id             serial
        primary key,
    session_id     integer                  not null,
    sync_type      varchar(20)              not null,
    status         varchar(20)              not null,
    records_synced integer default 0        not null,
    error_message  text,
    started_at     timestamp with time zone not null,
    completed_at   timestamp with time zone
);

comment on table game_sync_log is '游戏同步日志表（记录数据同步操作）';

comment on column game_sync_log.id is '日志ID';

comment on column game_sync_log.session_id is '会话ID';

comment on column game_sync_log.sync_type is '同步类型: battle_record=战绩, member_list=成员列表, wallet_update=钱包更新, room_list=房间列表, group_member=圈成员';

comment on column game_sync_log.status is '同步状态: success=成功, failed=失败, partial=部分成功';

comment on column game_sync_log.records_synced is '同步记录数';

comment on column game_sync_log.error_message is '错误信息';

comment on column game_sync_log.started_at is '开始时间';

comment on column game_sync_log.completed_at is '完成时间';

create index if not exists idx_sync_log_session_started
    on game_sync_log (session_id, started_at);

create index if not exists idx_sync_log_type
    on game_sync_log (sync_type);

create index if not exists idx_sync_log_status
    on game_sync_log (status);

create table if not exists game_shop_admin
(
    id              serial
        primary key,
    house_gid       integer                                not null,
    user_id         integer                                not null,
    role            varchar(20)                            not null,
    game_account_id integer,
    is_exclusive    boolean                  default true,
    created_at      timestamp with time zone default now() not null,
    updated_at      timestamp with time zone default now() not null,
    deleted_at      timestamp with time zone,
    constraint uk_shop_admin_house_user
        unique (house_gid, user_id)
);

comment on table game_shop_admin is '店铺管理员表';

comment on column game_shop_admin.id is '管理员ID';

comment on column game_shop_admin.house_gid is '店铺GID';

comment on column game_shop_admin.user_id is '用户ID';

comment on column game_shop_admin.role is '角色: admin=管理员, operator=操作员';

comment on column game_shop_admin.game_account_id is '关联的游戏账号ID';

comment on column game_shop_admin.is_exclusive is '是否独占（店铺管理员同时只能管理一个店铺）';

comment on column game_shop_admin.created_at is '创建时间';

comment on column game_shop_admin.updated_at is '更新时间';

comment on column game_shop_admin.deleted_at is '删除时间';

create index if not exists idx_shop_admin_house_gid
    on game_shop_admin (house_gid);

create index if not exists idx_shop_admin_user_id
    on game_shop_admin (user_id);

create index if not exists idx_shop_admin_game_account
    on game_shop_admin (game_account_id);

create index if not exists idx_shop_admin_exclusive
    on game_shop_admin (user_id, is_exclusive);

create table if not exists game_house_settings
(
    id          serial
        primary key,
    house_gid   integer                                   not null,
    fees_json   text                     default ''::text not null,
    share_fee   boolean                  default false    not null,
    push_credit integer                  default 0        not null,
    updated_at  timestamp with time zone default now()    not null,
    updated_by  integer                  default 0        not null
);

comment on table game_house_settings is '店铺设置表';

comment on column game_house_settings.id is '设置ID';

comment on column game_house_settings.house_gid is '店铺GID';

comment on column game_house_settings.fees_json is '运费规则JSON';

comment on column game_house_settings.share_fee is '分运开关';

comment on column game_house_settings.push_credit is '推送额度（单位：分）';

comment on column game_house_settings.updated_at is '更新时间';

comment on column game_house_settings.updated_by is '操作人用户ID';

create unique index if not exists uk_house
    on game_house_settings (house_gid);

create table if not exists game_member
(
    id             serial
        primary key,
    house_gid      integer                                                not null,
    game_id        integer                                                not null,
    game_name      varchar(64)              default ''::character varying not null,
    group_name     varchar(64)              default ''::character varying not null,
    balance        integer                  default 0                     not null,
    credit         integer                  default 0                     not null,
    forbid         boolean                  default false                 not null,
    recommender    varchar(64)              default ''::character varying,
    use_multi_gids boolean                  default false,
    active_gid     integer,
    created_at     timestamp with time zone default now()                 not null,
    updated_at     timestamp with time zone default now()                 not null
);

comment on table game_member is '游戏成员表（店铺内的玩家）';

comment on column game_member.id is '成员ID';

comment on column game_member.house_gid is '店铺GID';

comment on column game_member.game_id is '游戏ID';

comment on column game_member.game_name is '游戏昵称';

comment on column game_member.group_name is '圈名';

comment on column game_member.balance is '余额（单位：分）';

comment on column game_member.credit is '信用额度（单位：分）';

comment on column game_member.forbid is '是否禁用';

comment on column game_member.recommender is '推荐人';

comment on column game_member.use_multi_gids is '是否允许多开';

comment on column game_member.active_gid is '当前活跃游戏ID';

comment on column game_member.created_at is '创建时间';

comment on column game_member.updated_at is '更新时间';

create unique index if not exists uk_game_member_house_game
    on game_member (house_gid, game_id);

create index if not exists idx_game_member_game_id
    on game_member (game_id);

create index if not exists idx_game_member_house_group
    on game_member (house_gid, group_name);

create index if not exists idx_game_member_balance
    on game_member (balance);

create index if not exists idx_game_member_forbid
    on game_member (forbid);

create table if not exists game_member_wallet
(
    id         serial
        primary key,
    house_gid  integer                                not null,
    member_id  integer                                not null,
    balance    integer                  default 0     not null,
    forbid     boolean                  default false not null,
    limit_min  integer                  default 0     not null,
    updated_at timestamp with time zone default now() not null,
    updated_by integer                  default 0     not null
);

comment on table game_member_wallet is '游戏成员钱包表';

comment on column game_member_wallet.id is '钱包ID';

comment on column game_member_wallet.house_gid is '店铺GID';

comment on column game_member_wallet.member_id is '成员ID';

comment on column game_member_wallet.balance is '余额（单位：分）';

comment on column game_member_wallet.forbid is '是否禁用';

comment on column game_member_wallet.limit_min is '最低限额（单位：分）';

comment on column game_member_wallet.updated_at is '更新时间';

comment on column game_member_wallet.updated_by is '操作人用户ID';

create index if not exists idx_member_wallet_house
    on game_member_wallet (house_gid);

create index if not exists idx_member_wallet_member
    on game_member_wallet (member_id);

create table if not exists game_member_rule
(
    id           serial
        primary key,
    house_gid    integer                                not null,
    member_id    integer                                not null,
    vip          boolean                  default false not null,
    multi_gids   boolean                  default false not null,
    temp_release integer                  default 0     not null,
    expire_at    timestamp with time zone,
    updated_at   timestamp with time zone default now() not null,
    updated_by   integer                  default 0     not null
);

comment on table game_member_rule is '游戏成员规则表（VIP、多号、临时解禁等）';

comment on column game_member_rule.id is '规则ID';

comment on column game_member_rule.house_gid is '店铺GID';

comment on column game_member_rule.member_id is '成员ID';

comment on column game_member_rule.vip is '是否VIP';

comment on column game_member_rule.multi_gids is '是否允许多号';

comment on column game_member_rule.temp_release is '临时解禁上限（单位：分），0表示无限制';

comment on column game_member_rule.expire_at is '规则过期时间';

comment on column game_member_rule.updated_at is '更新时间';

comment on column game_member_rule.updated_by is '操作人用户ID';

create index if not exists idx_member_rule_house
    on game_member_rule (house_gid);

create index if not exists idx_member_rule_member
    on game_member_rule (member_id);

create table if not exists game_battle_record
(
    id               serial
        primary key,
    house_gid        integer                                not null,
    group_id         integer                                not null,
    room_uid         integer                                not null,
    kind_id          integer                                not null,
    base_score       integer                                not null,
    battle_at        timestamp with time zone               not null,
    players_json     text                                   not null,
    player_id        integer,
    player_game_id   integer,
    player_game_name varchar(64)              default ''::character varying,
    group_name       varchar(64)              default ''::character varying,
    score            integer                  default 0,
    fee              integer                  default 0,
    factor           numeric(10, 4)           default 1.0000,
    player_balance   integer                  default 0,
    created_at       timestamp with time zone default now() not null
);

comment on table game_battle_record is '游戏战绩表（本地战绩快照，一行代表一局一桌）';

comment on column game_battle_record.id is '战绩ID';

comment on column game_battle_record.house_gid is '店铺GID';

comment on column game_battle_record.group_id is '圈ID';

comment on column game_battle_record.room_uid is '房间唯一ID';

comment on column game_battle_record.kind_id is '游戏类型ID';

comment on column game_battle_record.base_score is '底分';

comment on column game_battle_record.battle_at is '对战时间';

comment on column game_battle_record.players_json is '玩家列表JSON';

comment on column game_battle_record.player_id is '玩家ID（用于按玩家查询）';

comment on column game_battle_record.player_game_id is '玩家游戏ID';

comment on column game_battle_record.player_game_name is '玩家游戏昵称';

comment on column game_battle_record.group_name is '圈名';

comment on column game_battle_record.score is '得分';

comment on column game_battle_record.fee is '服务费（单位：分）';

comment on column game_battle_record.factor is '结算比例';

comment on column game_battle_record.player_balance is '玩家余额（单位：分）';

comment on column game_battle_record.created_at is '创建时间';

create index if not exists idx_battle_house_gid
    on game_battle_record (house_gid);

create index if not exists idx_battle_room_uid
    on game_battle_record (room_uid);

create index if not exists idx_battle_kind_id
    on game_battle_record (kind_id);

create index if not exists idx_battle_at
    on game_battle_record (battle_at);

create index if not exists idx_battle_player_id
    on game_battle_record (player_id);

create index if not exists idx_battle_player_game_id
    on game_battle_record (player_game_id);

create index if not exists idx_battle_group
    on game_battle_record (house_gid, group_name);

create table if not exists game_recharge_record
(
    id               serial
        primary key,
    house_gid        integer                                                not null,
    player_id        integer                                                not null,
    group_name       varchar(64)              default ''::character varying not null,
    amount           integer                                                not null,
    balance_before   integer                                                not null,
    balance_after    integer                                                not null,
    operator_user_id integer,
    recharged_at     timestamp with time zone                               not null,
    created_at       timestamp with time zone default now()                 not null
);

comment on table game_recharge_record is '充值记录表（上下分记录）';

comment on column game_recharge_record.id is '记录ID';

comment on column game_recharge_record.house_gid is '店铺GID';

comment on column game_recharge_record.player_id is '玩家ID';

comment on column game_recharge_record.group_name is '圈名';

comment on column game_recharge_record.amount is '金额（单位：分），正数=充值，负数=提现';

comment on column game_recharge_record.balance_before is '操作前余额（单位：分）';

comment on column game_recharge_record.balance_after is '操作后余额（单位：分）';

comment on column game_recharge_record.operator_user_id is '操作人用户ID';

comment on column game_recharge_record.recharged_at is '充值时间';

comment on column game_recharge_record.created_at is '创建时间';

create index if not exists idx_recharge_house_gid
    on game_recharge_record (house_gid);

create index if not exists idx_recharge_player
    on game_recharge_record (player_id);

create index if not exists idx_recharge_recharged_at
    on game_recharge_record (recharged_at);

create index if not exists idx_recharge_house_group
    on game_recharge_record (house_gid, group_name);

create table if not exists game_fee_settle
(
    id         serial
        primary key,
    house_gid  integer                                                not null,
    play_group varchar(32)              default ''::character varying not null,
    amount     integer                                                not null,
    feed_at    timestamp with time zone                               not null,
    created_at timestamp with time zone default now()                 not null
);

comment on table game_fee_settle is '费用结算表（本地费用结算记录）';

comment on column game_fee_settle.id is '结算ID';

comment on column game_fee_settle.house_gid is '店铺GID';

comment on column game_fee_settle.play_group is '圈名';

comment on column game_fee_settle.amount is '金额（单位：分）';

comment on column game_fee_settle.feed_at is '结算时间';

comment on column game_fee_settle.created_at is '创建时间';

create index if not exists idx_fee_settle_house
    on game_fee_settle (house_gid);

create index if not exists idx_fee_settle_feed_at
    on game_fee_settle (feed_at);

create table if not exists game_shop_group
(
    id            serial
        primary key,
    house_gid     integer                                not null,
    group_name    varchar(64)                            not null,
    admin_user_id integer                                not null,
    description   text                     default ''::text,
    is_active     boolean                  default true  not null,
    created_at    timestamp with time zone default now() not null,
    updated_at    timestamp with time zone default now() not null
);

comment on table game_shop_group is '店铺圈子表（每个店铺管理员对应一个圈子）';

comment on column game_shop_group.id is '圈子ID';

comment on column game_shop_group.house_gid is '店铺GID';

comment on column game_shop_group.group_name is '圈子名称';

comment on column game_shop_group.admin_user_id is '圈主用户ID（店铺管理员）';

comment on column game_shop_group.description is '圈子描述';

comment on column game_shop_group.is_active is '是否激活';

comment on column game_shop_group.created_at is '创建时间';

comment on column game_shop_group.updated_at is '更新时间';

create unique index if not exists uk_shop_group_house_admin
    on game_shop_group (house_gid, admin_user_id)
    where (is_active = true);

create index if not exists idx_shop_group_house
    on game_shop_group (house_gid);

create index if not exists idx_shop_group_admin
    on game_shop_group (admin_user_id);

create table if not exists game_shop_group_member
(
    id         serial
        primary key,
    group_id   integer                                not null,
    user_id    integer                                not null,
    joined_at  timestamp with time zone default now() not null,
    created_at timestamp with time zone default now() not null
);

comment on table game_shop_group_member is '圈子成员关系表（用户可以加入多个圈子）';

comment on column game_shop_group_member.id is '关系ID';

comment on column game_shop_group_member.group_id is '圈子ID';

comment on column game_shop_group_member.user_id is '用户ID';

comment on column game_shop_group_member.joined_at is '加入时间';

comment on column game_shop_group_member.created_at is '创建时间';

create unique index if not exists uk_group_member_group_user
    on game_shop_group_member (group_id, user_id);

create index if not exists idx_group_member_group
    on game_shop_group_member (group_id);

create index if not exists idx_group_member_user
    on game_shop_group_member (user_id);

create table if not exists basic_menu
(
    id               serial
        primary key,
    parent_id        integer                     default '-1'::integer not null,
    menu_type        integer                                           not null,
    title            varchar(255)                                      not null,
    name             varchar(255)                                      not null,
    path             varchar(255)                                      not null,
    component        varchar(255)                                      not null,
    rank             varchar(255),
    redirect         varchar(255)                                      not null,
    icon             varchar(255)                                      not null,
    extra_icon       varchar(255)                                      not null,
    enter_transition varchar(255)                                      not null,
    leave_transition varchar(255)                                      not null,
    active_path      varchar(255)                                      not null,
    auths            varchar(255)                                      not null,
    frame_src        varchar(255)                                      not null,
    frame_loading    boolean                     default false         not null,
    keep_alive       boolean                     default false         not null,
    hidden_tag       boolean                     default false         not null,
    fixed_tag        boolean                     default false         not null,
    show_link        boolean                     default true          not null,
    show_parent      boolean                     default true          not null,
    created_at       timestamp(6) with time zone default now()         not null,
    updated_at       timestamp(6) with time zone default now()         not null,
    deleted_at       timestamp(6) with time zone,
    is_del           smallint                    default 0             not null
);

comment on table basic_menu is '基础菜单表';

comment on column basic_menu.id is '菜单ID';

comment on column basic_menu.parent_id is '父级菜单ID，默认为-1表示顶级菜单';

comment on column basic_menu.menu_type is '菜单类型：1=一级菜单，2=二级菜单';

comment on column basic_menu.title is '菜单标题';

comment on column basic_menu.name is '菜单名称（唯一标识）';

comment on column basic_menu.path is '路由路径';

comment on column basic_menu.component is '组件路径';

comment on column basic_menu.rank is '排序';

comment on column basic_menu.redirect is '重定向路径';

comment on column basic_menu.icon is '图标';

comment on column basic_menu.extra_icon is '额外图标';

comment on column basic_menu.enter_transition is '进入动画';

comment on column basic_menu.leave_transition is '离开动画';

comment on column basic_menu.active_path is '激活路径';

comment on column basic_menu.auths is '权限标识（逗号分隔）';

comment on column basic_menu.frame_src is '内嵌iframe地址';

comment on column basic_menu.frame_loading is '是否显示加载动画';

comment on column basic_menu.keep_alive is '是否缓存页面';

comment on column basic_menu.hidden_tag is '是否隐藏标签';

comment on column basic_menu.fixed_tag is '是否固定标签';

comment on column basic_menu.show_link is '是否显示链接';

comment on column basic_menu.show_parent is '是否显示父级';

create table if not exists basic_role_menu_rel
(
    role_id integer not null,
    menu_id integer not null,
    primary key (role_id, menu_id)
);

comment on table basic_role_menu_rel is '角色菜单关联表';

comment on column basic_role_menu_rel.role_id is '角色ID';

comment on column basic_role_menu_rel.menu_id is '菜单ID';

create table if not exists basic_permission
(
    id          serial
        primary key,
    code        varchar(100)                              not null,
    name        varchar(255)                              not null,
    category    varchar(50)                               not null,
    description text,
    created_at  timestamp(6) with time zone default now() not null,
    updated_at  timestamp(6) with time zone default now() not null,
    is_deleted  boolean                     default false not null
);

comment on table basic_permission is '基础权限表（细粒度权限定义）';

comment on column basic_permission.id is '权限ID';

comment on column basic_permission.code is '权限编码（唯一标识）';

comment on column basic_permission.name is '权限名称';

comment on column basic_permission.category is '权限分类：stats/fund/shop/game/system';

comment on column basic_permission.description is '权限描述';

comment on column basic_permission.created_at is '创建时间';

comment on column basic_permission.updated_at is '更新时间';

comment on column basic_permission.is_deleted is '是否删除';

create unique index if not exists uk_basic_permission_code
    on basic_permission (code)
    where (is_deleted = false);

create index if not exists idx_basic_permission_category
    on basic_permission (category);

create table if not exists basic_role_permission_rel
(
    role_id       integer not null,
    permission_id integer not null,
    primary key (role_id, permission_id)
);

comment on table basic_role_permission_rel is '角色权限关联表';

comment on column basic_role_permission_rel.role_id is '角色ID';

comment on column basic_role_permission_rel.permission_id is '权限ID';

create table if not exists basic_menu_button
(
    id               serial
        primary key,
    menu_id          integer                                   not null,
    button_code      varchar(100)                              not null,
    button_name      varchar(255)                              not null,
    permission_codes varchar(500)                              not null,
    created_at       timestamp(6) with time zone default now() not null,
    updated_at       timestamp(6) with time zone default now() not null
);

comment on table basic_menu_button is '菜单按钮权限配置表（UI细粒度控制）';

comment on column basic_menu_button.id is '按钮ID';

comment on column basic_menu_button.menu_id is '所属菜单ID';

comment on column basic_menu_button.button_code is '按钮编码';

comment on column basic_menu_button.button_name is '按钮名称';

comment on column basic_menu_button.permission_codes is '所需权限码（逗号分隔，满足任一即可）';

create index if not exists idx_menu_button_menu_id
    on basic_menu_button (menu_id);

create unique index if not exists uk_menu_button_menu_code
    on basic_menu_button (menu_id, button_code);

create table if not exists game_shop_application_log
(
    id            bigserial
        primary key,
    house_gid     integer                                not null,
    applier_gid   integer                                not null,
    applier_gname varchar(100)                           not null,
    action        varchar(20)                            not null,
    admin_user_id integer                                not null,
    admin_game_id integer,
    created_at    timestamp with time zone default now() not null
);

comment on table game_shop_application_log is '店铺申请操作日志（记录通过/拒绝操作，申请数据存在内存中）';

comment on column game_shop_application_log.id is '主键ID';

comment on column game_shop_application_log.house_gid is '店铺游戏ID';

comment on column game_shop_application_log.applier_gid is '申请人游戏ID';

comment on column game_shop_application_log.applier_gname is '申请人游戏昵称';

comment on column game_shop_application_log.action is '操作类型：approved=通过，rejected=拒绝';

comment on column game_shop_application_log.admin_user_id is '处理的管理员系统用户ID';

comment on column game_shop_application_log.admin_game_id is '管理员游戏ID（可选）';

comment on column game_shop_application_log.created_at is '操作时间';

create index if not exists idx_app_log_house
    on game_shop_application_log (house_gid);

create index if not exists idx_app_log_admin
    on game_shop_application_log (admin_user_id);

create index if not exists idx_app_log_created
    on game_shop_application_log (created_at desc);

create index if not exists idx_app_log_action
    on game_shop_application_log (action);

create or replace function update_updated_at_column() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

create trigger update_basic_user_updated_at
    before update
    on basic_user
    for each row
execute procedure update_updated_at_column();

create trigger update_game_account_updated_at
    before update
    on game_account
    for each row
execute procedure update_updated_at_column();

create trigger update_game_ctrl_account_updated_at
    before update
    on game_ctrl_account
    for each row
execute procedure update_updated_at_column();

create trigger update_game_account_house_updated_at
    before update
    on game_account_house
    for each row
execute procedure update_updated_at_column();

create trigger update_game_account_store_binding_updated_at
    before update
    on game_account_store_binding
    for each row
execute procedure update_updated_at_column();

create trigger update_game_session_updated_at
    before update
    on game_session
    for each row
execute procedure update_updated_at_column();

create trigger update_game_shop_admin_updated_at
    before update
    on game_shop_admin
    for each row
execute procedure update_updated_at_column();

create trigger update_game_house_settings_updated_at
    before update
    on game_house_settings
    for each row
execute procedure update_updated_at_column();

create trigger update_game_member_updated_at
    before update
    on game_member
    for each row
execute procedure update_updated_at_column();

create trigger update_game_member_wallet_updated_at
    before update
    on game_member_wallet
    for each row
execute procedure update_updated_at_column();

create trigger update_game_member_rule_updated_at
    before update
    on game_member_rule
    for each row
execute procedure update_updated_at_column();

create trigger update_game_shop_group_updated_at
    before update
    on game_shop_group
    for each row
execute procedure update_updated_at_column();

create trigger update_basic_menu_updated_at
    before update
    on basic_menu
    for each row
execute procedure update_updated_at_column();


