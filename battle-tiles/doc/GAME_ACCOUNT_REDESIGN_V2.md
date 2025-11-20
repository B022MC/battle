# 游戏账号系统重新设计方案 V2

## 📋 核心原则

### 1. 游戏账号入圈（不是用户入圈）
- 游戏账号是核心实体
- 圈子关系绑定在游戏账号上
- 用户只是游戏账号的一个可选属性

### 2. 用户通过游戏账号反向查询
- 用户查询圈子：用户 → 游戏账号 → 圈子
- 用户查询战绩：用户 → 游戏账号 → 战绩
- 用户上下分：用户 → 游戏账号 → 操作

---

## 🗄️ 数据库结构（极简版）

### 1. `game_account` 表调整

```sql
-- 修改：user_id 改为可选
ALTER TABLE game_account 
ALTER COLUMN user_id DROP NOT NULL;

COMMENT ON COLUMN game_account.user_id IS '关联的用户ID（可选，用于用户反向查询）';
```

**字段说明：**
- `id` - 游戏账号ID（主键）
- `user_id` - 关联的用户ID（**可为NULL**）
- `game_user_id` - 游戏服务器用户ID（唯一标识）
- `account` - 游戏账号
- `nickname` - 游戏昵称

---

### 2. 新增 `game_account_group` 表

```sql
CREATE TABLE game_account_group (
    id SERIAL PRIMARY KEY,
    game_account_id INTEGER NOT NULL,           -- 游戏账号ID
    house_gid INTEGER NOT NULL,                 -- 店铺GID
    group_id INTEGER NOT NULL,                  -- 圈子ID
    group_name VARCHAR(64) NOT NULL,            -- 圈子名称
    admin_user_id INTEGER NOT NULL,             -- 圈主用户ID
    approved_by_user_id INTEGER NOT NULL,       -- 审批人用户ID
    status VARCHAR(20) DEFAULT 'active' NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 唯一约束：一个游戏账号在一个店铺只能属于一个圈子
    CONSTRAINT uk_game_account_house UNIQUE (game_account_id, house_gid)
);

-- 索引
CREATE INDEX idx_account_group_game_account ON game_account_group(game_account_id);
CREATE INDEX idx_account_group_house ON game_account_group(house_gid);
CREATE INDEX idx_account_group_group ON game_account_group(group_id);
CREATE INDEX idx_account_group_status ON game_account_group(status);
```

---

### 3. `game_battle_record` 表调整

```sql
-- 修改：player_id 改为可选
ALTER TABLE game_battle_record 
ALTER COLUMN player_id DROP NOT NULL;

COMMENT ON COLUMN game_battle_record.player_id IS '玩家用户ID（可选，用于用户反向查询）';
COMMENT ON COLUMN game_battle_record.player_game_id IS '玩家游戏账号ID（game_account.id，必填）';
```

---

## 🔄 核心流程

### 流程 1：游戏内申请 + 管理员审批

```
┌─────────────────────────────┐
│ 1. 玩家在游戏内发起申请      │
│    - applier_gid (游戏ID)   │
│    - house_gid (店铺ID)     │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 2. 后端接收申请（内存队列）  │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 3. 管理员点击"通过"          │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 4. 查找或创建游戏账号                   │
│    SELECT * FROM game_account           │
│    WHERE game_user_id = applier_gid     │
│                                         │
│    如果不存在：                         │
│    INSERT INTO game_account             │
│    (user_id, game_user_id, account,    │
│     nickname, status)                   │
│    VALUES                               │
│    (NULL, applier_gid, applier_gid,    │
│     applier_gname, 1)                   │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 5. 获取或创建管理员的圈子               │
│    SELECT * FROM game_shop_group        │
│    WHERE house_gid = ?                  │
│      AND admin_user_id = current_admin  │
│                                         │
│    如果不存在，创建圈子：               │
│    INSERT INTO game_shop_group          │
│    (house_gid, group_name,              │
│     admin_user_id)                      │
│    VALUES (?, '管理员名_圈', ?)         │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 6. 游戏账号入圈                         │
│    INSERT INTO game_account_group       │
│    (game_account_id, house_gid,         │
│     group_id, group_name,               │
│     admin_user_id, approved_by_user_id) │
│    VALUES (?, ?, ?, ?, ?, ?)            │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 7. 调用游戏API拉入圈子                  │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 8. 记录操作日志                         │
│    INSERT INTO game_shop_application_log│
└─────────────────────────────────────────┘
```

**关键点：**
- 游戏账号的 `user_id` 为 NULL
- 哪个管理员通过，游戏账号就进入哪个管理员的圈子
- 记录在 `game_account_group` 表

---

### 流程 2：用户注册 + 绑定游戏账号

```
┌─────────────────────────────┐
│ 1. 用户注册平台账号          │
│    INSERT INTO basic_user   │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 2. 用户登录后绑定游戏账号    │
│    输入：游戏账号、密码      │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 3. 验证游戏账号密码                     │
│    调用游戏API验证                      │
│    获取 game_user_id                    │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 4. 查找游戏账号                         │
│    SELECT * FROM game_account           │
│    WHERE game_user_id = ?               │
└──────────┬──────────────────────────────┘
           │
       ┌───┴───┐
       │存在？  │
       └───┬───┘
           │
       ┌───┴────┐
       │        │
      是       否
       │        │
       │        ▼
       │   ┌─────────────────────────────┐
       │   │ 创建游戏账号                 │
       │   │ INSERT INTO game_account    │
       │   │ (user_id, game_user_id,     │
       │   │  account, nickname)         │
       │   │ VALUES                      │
       │   │ (current_user_id, ?, ?, ?)  │
       │   └──────┬──────────────────────┘
       │          │
       └──────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 5. 更新游戏账号绑定用户                 │
│    UPDATE game_account                  │
│    SET user_id = current_user_id        │
│    WHERE game_user_id = ?               │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 6. 返回绑定成功                         │
└─────────────────────────────────────────┘
```

**关键点：**
- 游戏账号可能已经存在（已入圈），只需更新 `user_id`
- 游戏账号可能不存在，创建时直接绑定用户

---

### 流程 3：战绩同步

```
┌─────────────────────────────┐
│ 1. 游戏服务器推送战绩        │
│    - room_uid               │
│    - players[]              │
│      - game_user_id         │
│      - nickname             │
│      - score                │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 2. 遍历每个玩家              │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 3. 查询游戏账号及圈子信息               │
│    SELECT                               │
│      ga.id AS game_account_id,          │
│      ga.user_id,                        │
│      gag.group_id,                      │
│      gag.group_name                     │
│    FROM game_account ga                 │
│    LEFT JOIN game_account_group gag     │
│      ON ga.id = gag.game_account_id     │
│      AND gag.house_gid = ?              │
│      AND gag.status = 'active'          │
│    WHERE ga.game_user_id = ?            │
└──────────┬──────────────────────────────┘
           │
       ┌───┴───┐
       │找到？  │
       └───┬───┘
           │
       ┌───┴────┐
       │        │
      是       否
       │        │
       │        ▼
       │   ┌─────────────────────────────┐
       │   │ 跳过该玩家                   │
       │   │ （游戏账号不存在）           │
       │   └─────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│ 4. 检查是否在圈内                       │
│    IF group_id IS NULL THEN             │
│      跳过该玩家（未入圈）               │
│    END IF                               │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 5. 保存战绩记录                         │
│    INSERT INTO game_battle_record       │
│    (house_gid, group_id,                │
│     player_id,          -- 可为NULL     │
│     player_game_id,     -- 必填         │
│     player_game_name,                   │
│     group_name,                         │
│     score, fee, ...)                    │
│    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ...) │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 6. 更新成员余额                         │
│    UPDATE game_member                   │
│    SET balance = balance + score - fee  │
│    WHERE house_gid = ?                  │
│      AND game_id = ?                    │
└─────────────────────────────────────────┘
```

**关键点：**
- 一次查询获取游戏账号和圈子信息
- `player_id` 可为NULL（未绑定用户）
- `player_game_id` 必填（游戏账号ID）
- 只记录已入圈的游戏账号

---

### 流程 4：用户查询圈子（反向查询）

```
┌─────────────────────────────┐
│ 1. 用户请求查询我的圈子      │
│    GET /my/groups           │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 2. 查询用户绑定的游戏账号               │
│    SELECT id FROM game_account          │
│    WHERE user_id = current_user_id      │
│      AND is_del = 0                     │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 3. 查询游戏账号的圈子                   │
│    SELECT                               │
│      gag.*,                             │
│      gsg.description,                   │
│      bu.nick_name AS admin_name         │
│    FROM game_account_group gag          │
│    JOIN game_shop_group gsg             │
│      ON gag.group_id = gsg.id           │
│    JOIN basic_user bu                   │
│      ON gag.admin_user_id = bu.id       │
│    WHERE gag.game_account_id IN (...)   │
│      AND gag.status = 'active'          │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 4. 返回圈子列表                         │
└─────────────────────────────────────────┘
```

**查询路径：** 用户 → 游戏账号 → 圈子

---

### 流程 5：用户查询战绩（反向查询）

```
┌─────────────────────────────┐
│ 1. 用户请求查询我的战绩      │
│    GET /shops/my-battles    │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 2. 查询用户绑定的游戏账号               │
│    SELECT id FROM game_account          │
│    WHERE user_id = current_user_id      │
│      AND is_del = 0                     │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 3. 查询游戏账号的圈子                   │
│    SELECT group_name                    │
│    FROM game_account_group              │
│    WHERE game_account_id IN (...)       │
│      AND status = 'active'              │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 4. 查询战绩记录                         │
│    SELECT * FROM game_battle_record     │
│    WHERE player_game_id IN (...)        │
│      AND group_name IN (...)            │
│    ORDER BY battle_at DESC              │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 5. 返回战绩列表                         │
└─────────────────────────────────────────┘
```

**查询路径：** 用户 → 游戏账号 → 圈子 → 战绩

---

### 流程 6：用户上分/下分

```
┌─────────────────────────────┐
│ 1. 管理员选择游戏账号        │
│    - game_account_id        │
│    - amount                 │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 2. 查询游戏账号的圈子关系               │
│    SELECT * FROM game_account_group     │
│    WHERE game_account_id = ?            │
│      AND house_gid = ?                  │
│      AND status = 'active'              │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 3. 验证权限                             │
│    IF admin_user_id != current_user_id  │
│       AND NOT super_admin THEN          │
│      返回错误：无权限                   │
│    END IF                               │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 4. 调用游戏API执行上分/下分             │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 5. 更新本地余额                         │
│    UPDATE game_member                   │
│    SET balance = balance + amount       │
└──────────┬──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ 6. 记录充值记录                         │
│    INSERT INTO game_recharge_record     │
└─────────────────────────────────────────┘
```

---

## 📊 数据关系（简化版）

```
┌──────────────┐
│  basic_user  │ (用户表)
│  - id        │
└──────┬───────┘
       │ 1
       │ (可选绑定)
       │ 0..n
       ▼
┌──────────────────┐
│  game_account    │ (游戏账号表) ⭐核心
│  - id            │
│  - user_id       │ (可为NULL)
│  - game_user_id  │ (唯一标识)
└──────┬───────────┘
       │ 1
       │ (入圈)
       │ 0..n
       ▼
┌──────────────────────┐
│ game_account_group   │ (游戏账号圈子关系) ⭐核心
│  - game_account_id   │
│  - house_gid         │
│  - group_id          │
│  - group_name        │
│  - admin_user_id     │
└──────────────────────┘
       │
       │ n:1
       ▼
┌──────────────────┐
│ game_shop_group  │ (圈子表)
│  - id            │
│  - group_name    │
│  - admin_user_id │
└──────────────────┘
```

---

## 🔍 常用查询示例

### 1. 战绩同步时查询游戏账号和圈子
```sql
SELECT 
    ga.id AS game_account_id,
    ga.user_id,
    gag.group_id,
    gag.group_name
FROM game_account ga
LEFT JOIN game_account_group gag 
    ON ga.id = gag.game_account_id 
    AND gag.house_gid = 58959
    AND gag.status = 'active'
WHERE ga.game_user_id = '22805688'
  AND ga.is_del = 0;
```

### 2. 用户查询自己的圈子
```sql
-- 第一步：获取用户的游戏账号
SELECT id FROM game_account 
WHERE user_id = 123 AND is_del = 0;

-- 第二步：获取游戏账号的圈子
SELECT 
    gag.*,
    gsg.description,
    bu.nick_name AS admin_name
FROM game_account_group gag
JOIN game_shop_group gsg ON gag.group_id = gsg.id
JOIN basic_user bu ON gag.admin_user_id = bu.id
WHERE gag.game_account_id IN (456, 789)
  AND gag.status = 'active';
```

### 3. 用户查询自己的战绩
```sql
-- 第一步：获取用户的游戏账号
SELECT id FROM game_account 
WHERE user_id = 123 AND is_del = 0;

-- 第二步：查询战绩
SELECT * FROM game_battle_record
WHERE player_game_id IN (456, 789)
ORDER BY battle_at DESC
LIMIT 50;
```

### 4. 管理员查询圈内成员的游戏账号
```sql
SELECT 
    ga.*,
    gag.joined_at,
    gag.status
FROM game_account_group gag
JOIN game_account ga ON gag.game_account_id = ga.id
WHERE gag.group_id = 1
  AND gag.status = 'active'
ORDER BY gag.joined_at DESC;
```

---

## 🎯 核心要点

### 1. 游戏账号是核心
- ✅ 游戏账号可以独立存在（user_id 可为NULL）
- ✅ 圈子关系绑定在游戏账号上
- ✅ 战绩记录绑定在游戏账号上

### 2. 用户是可选属性
- ✅ 用户绑定游戏账号只是为了方便查询
- ✅ 用户通过游戏账号反向查询圈子和战绩
- ✅ 未绑定用户的游戏账号也能正常游戏

### 3. 简洁的设计
- ✅ 不需要视图
- ✅ 不需要存储过程
- ✅ 只需要两次简单查询

### 4. 查询路径清晰
- 用户查圈子：用户 → 游戏账号 → 圈子
- 用户查战绩：用户 → 游戏账号 → 战绩
- 战绩同步：游戏ID → 游戏账号 → 圈子

---

## 📝 总结

这个简化设计实现了：
1. ✅ 游戏账号入圈（不是用户入圈）
2. ✅ 用户通过游戏账号反向查询
3. ✅ 管理员审批自动分配圈子
4. ✅ 支持未注册用户的游戏账号
5. ✅ 极简的数据库结构
6. ✅ 清晰的查询路径

**核心表：**
- `game_account` - 游戏账号（核心）
- `game_account_group` - 游戏账号圈子关系（核心）
- `game_battle_record` - 战绩记录

**核心字段调整：**
- `game_account.user_id` - 可为NULL
- `game_battle_record.player_id` - 可为NULL

