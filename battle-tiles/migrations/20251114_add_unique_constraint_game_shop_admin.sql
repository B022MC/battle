-- 添加 game_shop_admin 表的唯一约束
-- 确保一个用户在一个店铺只能有一个管理员记录

-- 先删除可能存在的重复数据（保留最新的记录）
DELETE FROM game_shop_admin a
USING game_shop_admin b
WHERE a.id < b.id
  AND a.house_gid = b.house_gid
  AND a.user_id = b.user_id;

-- 添加唯一约束
ALTER TABLE game_shop_admin
ADD CONSTRAINT uk_shop_admin_house_user UNIQUE (house_gid, user_id);

