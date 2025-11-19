-- 修复 game_account 表中空的 game_user_id 字段
-- 这个脚本用于修复在修复代码之前注册的用户的 game_user_id

-- 注意：这个脚本需要手动执行，因为需要从游戏服务器获取用户ID
-- 目前只能通过以下方式修复：
-- 1. 让用户重新绑定游戏账号（推荐）
-- 2. 或者通过管理员界面手动验证账号

-- 查看所有 game_user_id 为空的记录
SELECT 
    ga.id,
    ga.user_id,
    ga.account,
    ga.nickname,
    bu.username,
    bu.nick_name,
    ga.created_at
FROM game_account ga
LEFT JOIN basic_user bu ON ga.user_id = bu.id
WHERE ga.game_user_id = '' OR ga.game_user_id IS NULL
ORDER BY ga.created_at DESC;

-- 如果需要删除这些空的记录（谨慎操作）：
-- DELETE FROM game_account 
-- WHERE (game_user_id = '' OR game_user_id IS NULL)
-- AND created_at < '2025-11-15'::timestamp;

-- 建议：通过前端提示用户重新绑定游戏账号，这样可以获取正确的 game_user_id

