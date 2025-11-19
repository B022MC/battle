-- ============================================
-- 上分/下分功能 - 权限分配配置
-- ============================================
-- 功能说明：店铺管理员可以为成员充值（上分）或提现（下分）
-- 前端路径：/members 页面中的操作按钮
-- 后端API：/members/credit/deposit, /members/credit/withdraw, /members/credit/force_withdraw
-- ============================================
-- 注意：权限已在 dml.sql 中定义：
--   - id=4: fund:deposit（上分）
--   - id=5: fund:withdraw（下分）
--   - id=6: fund:force_withdraw（强制下分）
-- 本SQL仅为角色分配这些已存在的权限
-- ============================================

-- 1. 为超级管理员角色分配权限
INSERT INTO "basic_role_permission_rel" ("role_id", "permission_id")
SELECT 1, id
FROM "basic_permission"
WHERE code IN ('fund:deposit', 'fund:withdraw', 'fund:force_withdraw')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 2. 为店铺管理员角色分配权限
INSERT INTO "basic_role_permission_rel" ("role_id", "permission_id")
SELECT 2, id
FROM "basic_permission"
WHERE code IN ('fund:deposit', 'fund:withdraw', 'fund:force_withdraw')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================
-- 说明：
-- 1. 上分/下分功能集成在成员管理页面中，无需单独菜单
-- 2. 在成员列表中，每个成员旁边会显示"上分"和"下分"按钮
-- 3. 需要相应权限才能看到和使用这些按钮
-- 4. fund:force_withdraw 用于特殊情况的强制提现
-- ============================================

-- 验证SQL
SELECT 
    p.id,
    p.code,
    p.name,
    p.description,
    CASE 
        WHEN EXISTS(SELECT 1 FROM basic_role_permission_rel WHERE role_id = 1 AND permission_id = p.id) 
        THEN '是' ELSE '否' 
    END as "超级管理员",
    CASE 
        WHEN EXISTS(SELECT 1 FROM basic_role_permission_rel WHERE role_id = 2 AND permission_id = p.id) 
        THEN '是' ELSE '否' 
    END as "店铺管理员"
FROM basic_permission p
WHERE p.code IN ('fund:deposit', 'fund:withdraw', 'fund:force_withdraw')
ORDER BY p.id;
