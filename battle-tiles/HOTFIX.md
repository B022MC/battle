# 紧急修复：game_player_id 字段缺失

## 问题
代码已更新为使用 `game_player_id`，但数据库表 `game_account_group` 还没有此字段。

## 错误日志
```
ERROR: column "game_player_id" does not exist (SQLSTATE 42703)
```

## 快速修复步骤

### 1. 停止服务（可选，建议）
```bash
# 查找进程
lsof -ti:8000

# 停止服务
lsof -ti:8000 | xargs kill -9
```

### 2. 执行迁移脚本
```bash
# 连接数据库
psql -h 8.137.52.203 -p 26655 -U B022MC -d battle-tiles-dev

# 执行迁移
\i /Users/b022mc/project/battle/battle-tiles/migrations/20251122_add_game_player_id_to_account_group.sql

# 验证
\d game_account_group

# 查看数据
SELECT 
    COUNT(*) as total,
    COUNT(game_player_id) as with_player_id,
    COUNT(game_account_id) as with_account_id
FROM game_account_group;
```

### 3. 重启服务
```bash
cd /Users/b022mc/project/battle/battle-tiles
./battle-tiles &

# 查看日志
tail -f logs/battle-tiles.log
```

## 一行命令修复（推荐）
```bash
psql -h 8.137.52.203 -p 26655 -U B022MC -d battle-tiles-dev -f /Users/b022mc/project/battle/battle-tiles/migrations/20251122_add_game_player_id_to_account_group.sql && lsof -ti:8000 | xargs kill -9 && cd /Users/b022mc/project/battle/battle-tiles && ./battle-tiles &
```

## 验证
等待几分钟后检查日志，应该不再有 `column "game_player_id" does not exist` 错误。

---
**创建时间**: 2025-11-22 01:32
