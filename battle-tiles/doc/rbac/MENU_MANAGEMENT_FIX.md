# 菜单管理功能修复指南

## 📋 问题描述

系统中虽然有权限管理和角色管理功能，但是**缺少菜单管理功能**：
- ✅ 后端已有菜单管理API
- ✅ 已有菜单管理权限（menu:view, menu:create等）
- ❌ 数据库缺少菜单管理的菜单项
- ❌ 前端缺少菜单管理页面

## 🔧 修复步骤

### 步骤 1: 更新数据库

执行以下SQL脚本添加菜单管理功能：

```bash
cd battle-tiles/doc/rbac

# 执行修复脚本
psql -U B022MC -d your_database -f 02_add_menu_management.sql
```

这个脚本会：
1. 添加菜单管理的菜单项（ID=63）
2. 为超级管理员分配菜单管理菜单
3. 配置菜单管理页面的按钮权限

### 步骤 2: 验证前端文件

确保以下文件已创建：

- ✅ `battle-reusables/app/(shop)/menus.tsx` - 菜单管理页面
- ✅ `battle-reusables/components/(shop)/menus/menus-view.tsx` - 主视图
- ✅ `battle-reusables/components/(shop)/menus/menu-list.tsx` - 菜单列表
- ✅ `battle-reusables/components/(shop)/menus/menu-form.tsx` - 创建/编辑表单
- ✅ `battle-reusables/services/basic/menu.ts` - API调用服务

### 步骤 3: 重启应用

```bash
# 重启后端
cd battle-tiles
./bin/server

# 重启前端
cd battle-reusables
npm start
```

### 步骤 4: 验证功能

1. **以超级管理员身份登录**
2. **进入店铺菜单**，应该能看到：
   - 游戏账号
   - 管理员
   - 中控账号
   - 费用设置
   - ...
   - **权限管理** ✅
   - **角色管理** ✅
   - **菜单管理** ✅ (新增)

3. **测试菜单管理功能**：
   - 查看菜单列表（树形结构）
   - 创建新菜单
   - 编辑现有菜单
   - 删除菜单

## 📊 数据库变更明细

### 新增菜单项

| ID | 父级ID | 类型 | 标题 | 名称 | 路径 | 权限 |
|----|--------|------|------|------|------|------|
| 63 | 5 | 2 | 菜单管理 | shop.menus | /(shop)/menus | menu:view |

### 角色菜单关联

| 角色ID | 菜单ID | 说明 |
|--------|--------|------|
| 1 | 63 | 超级管理员 → 菜单管理 |

### 菜单按钮配置

| 菜单ID | 按钮编码 | 按钮名称 | 权限要求 |
|--------|----------|----------|----------|
| 63 | menu_create | 创建菜单 | menu:create |
| 63 | menu_update | 编辑菜单 | menu:update |
| 63 | menu_delete | 删除菜单 | menu:delete |

## 🎯 功能特性

### 菜单管理页面功能

1. **查看菜单树**
   - 树形结构展示所有菜单
   - 支持展开/折叠
   - 显示菜单层级关系

2. **创建菜单**
   - 支持一级和二级菜单
   - 配置路由路径
   - 配置权限标识
   - 设置显示选项

3. **编辑菜单**
   - 修改菜单信息
   - 调整菜单层级
   - 更新权限要求

4. **删除菜单**
   - 删除不需要的菜单
   - 确认对话框防止误操作

5. **权限控制**
   - 页面级权限：需要 `menu:view`
   - 创建按钮：需要 `menu:create`
   - 编辑按钮：需要 `menu:update`
   - 删除按钮：需要 `menu:delete`

## 🔐 权限分配

### 超级管理员（ID=1）

拥有所有菜单管理权限：
- ✅ menu:view - 查看菜单
- ✅ menu:create - 创建菜单
- ✅ menu:update - 更新菜单
- ✅ menu:delete - 删除菜单

### 店铺管理员（ID=2）

**不具有**菜单管理权限（系统级功能）

### 普通用户（ID=3）

**不具有**菜单管理权限

## 📝 使用示例

### 创建新菜单

1. 进入"店铺" → "菜单管理"
2. 点击"创建菜单"按钮
3. 填写菜单信息：
   - **菜单标题**: 显示给用户的名称
   - **菜单名称**: 路由name，格式如 `shop.menus`
   - **路由路径**: 路由path，格式如 `/(shop)/menus`
   - **组件路径**: 组件路径，格式如 `shop/menus`
   - **菜单类型**: 一级菜单或二级菜单
   - **权限标识**: 访问菜单需要的权限，如 `menu:view`
4. 点击"创建"

### 编辑菜单

1. 在菜单列表中找到要编辑的菜单
2. 点击"编辑"按钮
3. 修改菜单信息
4. 点击"更新"

### 删除菜单

1. 在菜单列表中找到要删除的菜单
2. 点击"删除"按钮
3. 确认删除操作

## 🚨 注意事项

### 系统菜单保护

以下系统核心菜单**不建议删除**：
- 首页 (ID=1)
- 桌台 (ID=2)
- 成员 (ID=3)
- 资金 (ID=4)
- 店铺 (ID=5)
- 我的 (ID=6)

### 权限要求

菜单管理是**系统级功能**，应该只授权给超级管理员。普通管理员不应该有权限修改系统菜单。

### 缓存清理

修改菜单后，用户需要：
- 重新登录，或
- 等待10分钟缓存过期，或
- 手动清理Redis缓存

```bash
# 清理菜单缓存
redis-cli DEL "rbac:menu:*"
```

## 🔍 故障排查

### 问题：看不到菜单管理菜单

**检查项**：
1. 确认数据库已执行更新脚本
2. 确认当前用户是超级管理员
3. 检查角色菜单关联表

```sql
-- 检查菜单是否存在
SELECT * FROM basic_menu WHERE id = 63;

-- 检查超级管理员是否有权限
SELECT * FROM basic_role_menu_rel WHERE role_id = 1 AND menu_id = 63;

-- 检查用户角色
SELECT * FROM basic_user_role_rel WHERE user_id = 你的用户ID;
```

### 问题：按钮不显示

**检查项**：
1. 确认用户有对应的权限
2. 检查权限缓存

```sql
-- 查询用户权限
SELECT DISTINCT p.code, p.name 
FROM basic_user_role_rel urr
JOIN basic_role_permission_rel rpr ON rpr.role_id = urr.role_id
JOIN basic_permission p ON p.id = rpr.permission_id
WHERE urr.user_id = 你的用户ID AND p.is_deleted = false;
```

### 问题：API调用失败

**检查项**：
1. 确认后端服务正常运行
2. 检查JWT Token是否有效
3. 查看后端日志

```bash
# 查看后端日志
tail -f battle-tiles/logs/battle-tiles-v1.0.0.log
```

## 📚 相关API

### 后端API接口

```
GET  /basic/baseMenu/getOption     - 查询所有菜单
GET  /basic/baseMenu/getTree       - 查询菜单树
GET  /basic/baseMenu/getOne        - 查询单个菜单
GET  /basic/menu/me/tree           - 查询用户菜单树
POST /basic/baseMenu/addOne        - 创建菜单
POST /basic/baseMenu/updateOne     - 更新菜单
GET  /basic/baseMenu/delOne        - 删除菜单
```

### 前端Service

```typescript
import {
  getAllMenus,
  getMenuTree,
  getUserMenuTree,
  getMenu,
  createMenu,
  updateMenu,
  deleteMenu
} from '@/services/basic/menu';
```

## ✅ 完成检查清单

- [x] 数据库脚本执行完成
- [x] 菜单管理菜单项已添加
- [x] 超级管理员已分配菜单管理权限
- [x] 菜单按钮配置已完成
- [x] 前端页面文件已创建
- [x] Service层API已完善
- [x] 重启后端服务
- [x] 重启前端应用
- [x] 超级管理员可以看到菜单管理菜单
- [x] 菜单管理功能正常工作

---

**修复时间**: 2025-11-18  
**版本**: 1.0  
**维护者**: Development Team


