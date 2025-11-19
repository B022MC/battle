# 店铺导航更新说明

## 问题
超级管理员看不到"权限管理"和"角色管理"页面。

## 原因
店铺页面（`app/(tabs)/shop.tsx`）的导航列表中缺少这两个菜单项的入口。

## 解决方案
已在 `app/(tabs)/shop.tsx` 中添加完整的店铺子页面导航，包括：

### 店铺管理功能
1. **游戏账号** - `/(shop)/account` (所有用户)
2. **管理员** - `/(shop)/admins` (需要 `shop:admin:view`)
3. **中控账号** - `/(shop)/rooms` (需要 `game:ctrl:view`)
4. **费用设置** - `/(shop)/fees` (需要 `shop:fees:view`)
5. **余额筛查** - `/(shop)/balances` (需要 `fund:wallet:view`)
6. **成员管理** - `/(shop)/members` (需要 `shop:member:view`)

### 战绩查询功能
1. **我的战绩** - `/(shop)/my-battles` (所有用户)
2. **我的余额** - `/(shop)/my-balances` (所有用户)
3. **圈子战绩(管理员)** - `/(shop)/group-battles` (需要 `shop:member:view`)
4. **圈子余额(管理员)** - `/(shop)/group-balances` (需要 `shop:member:view`)

### 系统管理功能 ✨ 新增
1. **权限管理** - `/(shop)/permissions` (需要 `permission:view`)
2. **角色管理** - `/(shop)/roles` (需要 `role:view`)

## 权限控制
所有管理功能都使用 `<PermissionGate>` 组件包裹，只有拥有相应权限的用户才能看到相应的菜单项。

## 测试
1. 使用超级管理员账号登录
2. 进入"店铺"标签页
3. 应该能看到"系统管理"分组
4. 点击"权限管理"或"角色管理"可以进入相应页面

## 注意事项
- 超级管理员拥有所有权限，会看到所有菜单项
- 店铺管理员会看到大部分管理功能，但**不会看到**"权限管理"和"角色管理"
- 普通用户只会看到"游戏账号"、"我的战绩"、"我的余额"三个菜单项

