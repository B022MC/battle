# 游戏账号系统简明指南

## 🎯 核心原则

### 游戏账号入圈，用户反向查询

```
旧：用户 → 圈子
新：游戏账号 → 圈子
    用户 → 游戏账号 → 圈子
```

---

## 📊 核心表

### 1. game_account (游戏账号)
```
- id              主键
- user_id         用户ID (可为NULL) ⭐
- game_user_id    游戏服务器ID (唯一)
- account         游戏账号
- nickname        游戏昵称
```

### 2. game_account_group (游戏账号圈子关系) ⭐新增
```
- id                    主键
- game_account_id       游戏账号ID
- house_gid             店铺GID
- group_id              圈子ID
- group_name            圈子名称
- admin_user_id         圈主ID
- approved_by_user_id   审批人ID
- status                状态
```

### 3. game_battle_record (战绩)
```
- player_id         用户ID (可为NULL) ⭐
- player_game_id    游戏账号ID (必填)
- group_name        圈名
- score             得分
```

---

## 🔄 关键流程

### 1. 申请入圈
```
游戏内申请 
  ↓
管理员审批通过
  ↓
查找/创建游戏账号 (user_id=NULL)
  ↓
获取管理员的圈子
  ↓
INSERT INTO game_account_group
  ↓
完成！游戏账号已入圈
```

### 2. 用户绑定
```
用户注册
  ↓
输入游戏账号密码
  ↓
验证游戏账号
  ↓
UPDATE game_account SET user_id = ?
  ↓
完成！用户已绑定游戏账号
```

### 3. 战绩同步
```sql
-- 一次查询获取游戏账号和圈子
SELECT 
    ga.id AS game_account_id,
    ga.user_id,
    gag.group_id,
    gag.group_name
FROM game_account ga
LEFT JOIN game_account_group gag 
    ON ga.id = gag.game_account_id 
    AND gag.house_gid = ?
    AND gag.status = 'active'
WHERE ga.game_user_id = ?;

-- 如果 group_id 不为空，保存战绩
INSERT INTO game_battle_record
  (player_id, player_game_id, group_name, ...)
VALUES
  (ga.user_id, ga.id, gag.group_name, ...);
```

### 4. 用户查圈子
```sql
-- 第一步：获取用户的游戏账号
SELECT id FROM game_account 
WHERE user_id = ? AND is_del = 0;

-- 第二步：获取圈子
SELECT * FROM game_account_group
WHERE game_account_id IN (...);
```

### 5. 用户查战绩
```sql
-- 第一步：获取用户的游戏账号
SELECT id FROM game_account 
WHERE user_id = ? AND is_del = 0;

-- 第二步：查询战绩
SELECT * FROM game_battle_record
WHERE player_game_id IN (...)
ORDER BY battle_at DESC;
```

---

## 🔑 关键点

1. **游戏账号是核心**
   - 圈子关系绑定在游戏账号上
   - 战绩记录绑定在游戏账号上

2. **用户是可选属性**
   - `game_account.user_id` 可为NULL
   - `game_battle_record.player_id` 可为NULL

3. **反向查询**
   - 用户查圈子：用户 → 游戏账号 → 圈子
   - 用户查战绩：用户 → 游戏账号 → 战绩

4. **管理员审批决定圈子**
   - 哪个管理员通过申请
   - 游戏账号就进入哪个管理员的圈子

---

## 📝 数据库迁移

```bash
psql -U postgres -d battle_db -f migrations/20251120_game_account_redesign_v2.sql
```

---

## ✅ 完成

- ✅ 游戏账号入圈
- ✅ 用户反向查询
- ✅ 支持未注册用户
- ✅ 极简设计
- ✅ 无视图，无函数

