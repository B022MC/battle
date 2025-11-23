# 房间额度限制 - 菜单权限配置说明

## 数据库表结构

### 1. `basic_menu` - 菜单表
用于定义前端路由和页面菜单

| 字段 | 类型 | 说明 |
|------|------|------|
| `title` | varchar | 菜单标题（显示名称） |
| `name` | varchar | 菜单名称（唯一标识，用于路由） |
| `path` | varchar | 路由路径 |
| `component` | varchar | 组件路径 |
| `menu_type` | int | 1=一级菜单, 2=二级菜单 |
| `parent_id` | int | 父级菜单ID，-1表示顶级 |
| `auths` | varchar | 权限标识（逗号分隔） |
| `rank` | varchar | 排序 |

### 2. `basic_menu_button` - 菜单按钮权限表
用于页面内按钮的权限控制

| 字段 | 类型 | 说明 |
|------|------|------|
| `menu_id` | int | 所属菜单ID |
| `button_code` | varchar | 按钮编码 |
| `button_name` | varchar | 按钮名称 |
| `permission_codes` | varchar | 所需权限码（逗号分隔，满足任一即可） |

### 3. `basic_permission` - 权限表
定义细粒度权限

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | varchar | 权限编码（唯一） |
| `name` | varchar | 权限名称 |
| `category` | varchar | 权限分类：game/shop/fund/stats/system |
| `description` | text | 权限描述 |

## 执行步骤

### 步骤1：执行SQL脚本

```bash
psql -U username -d database_name -f doc/rbac/add_room_credit_limit_menu.sql
```

执行后会创建：
- ✅ 4个权限记录 (`basic_permission`)
- ✅ 1个菜单记录 (`basic_menu`)
- ✅ 3个按钮权限记录 (`basic_menu_button`)

### 步骤2：查询新增权限的ID

执行SQL查看权限ID（用于后续分配给角色）：

```sql
SELECT id, code, name FROM basic_permission WHERE code LIKE 'game:room_credit:%' ORDER BY id;
```

结果示例：
```
 id  |          code             |      name
-----+---------------------------+-----------------
 47  | game:room_credit:view     | 查看房间额度
 48  | game:room_credit:set      | 设置房间额度
 49  | game:room_credit:delete   | 删除房间额度
 50  | game:room_credit:check    | 检查玩家额度
```

### 步骤3：为角色分配菜单

**方式1：通过管理后台分配**（推荐）
1. 登录管理后台
2. 进入"角色管理"
3. 编辑需要的角色
4. 勾选"房间额度"菜单
5. 保存

**方式2：通过SQL直接分配**

为超级管理员（role_id=1）分配菜单：
```sql
INSERT INTO basic_role_menu_rel (role_id, menu_id) 
VALUES (1, 64)  -- 64是房间额度菜单的ID
ON CONFLICT DO NOTHING;
```

为店铺管理员（role_id=2）分配菜单：
```sql
INSERT INTO basic_role_menu_rel (role_id, menu_id) 
VALUES (2, 64)
ON CONFLICT DO NOTHING;
```

### 步骤4：为角色分配权限

**方式1：通过管理后台分配**（推荐）
1. 登录管理后台
2. 进入"角色管理"
3. 编辑需要的角色
4. 点击"分配权限"
5. 勾选相关权限（查看/设置/删除/检查）
6. 保存

**方式2：通过SQL直接分配**

```sql
-- 为超级管理员分配所有房间额度权限
INSERT INTO basic_role_permission_rel (role_id, permission_id) 
VALUES 
(1, 50),  -- game:room_credit:view
(1, 51),  -- game:room_credit:set
(1, 52),  -- game:room_credit:delete
(1, 53)   -- game:room_credit:check
ON CONFLICT DO NOTHING;

-- 为店铺管理员分配所有权限
INSERT INTO basic_role_permission_rel (role_id, permission_id) 
VALUES 
(2, 50),  -- game:room_credit:view
(2, 51),  -- game:room_credit:set
(2, 52),  -- game:room_credit:delete
(2, 53)   -- game:room_credit:check
ON CONFLICT DO NOTHING;
```

## 前端路由配置

### 菜单路径
- **路由名称**: `shop.room_credits`
- **路由路径**: `/(shop)/room-credits`
- **组件路径**: `shop/room-credits`

### 前端文件结构

```
app/(shop)/room-credits/
├── index.tsx              # 主页面
├── _layout.tsx            # 布局（如果需要）
└── components/
    ├── CreditForm.tsx     # 设置表单
    └── CheckDialog.tsx    # 检查对话框
```

## 权限使用

### 权限标识列表

| 权限码 | 说明 | 使用场景 |
|--------|------|----------|
| `game:room_credit:view` | 查看房间额度 | 菜单访问、列表查询 |
| `game:room_credit:set` | 设置房间额度 | 创建/编辑额度规则 |
| `game:room_credit:delete` | 删除房间额度 | 删除额度规则 |
| `game:room_credit:check` | 检查玩家额度 | 检查玩家是否满足要求 |

### 前端权限判断示例

```typescript
import { usePermission } from '@/hooks/usePermission'

const { hasPermission } = usePermission()

// 检查是否有设置权限
if (hasPermission('game:room_credit:set')) {
  // 显示设置按钮
}

// 检查是否有删除权限
if (hasPermission('game:room_credit:delete')) {
  // 显示删除按钮
}
```

### React Native / Expo 按钮示例

```tsx
import { usePermissions } from '@/hooks/usePermissions'

export default function RoomCreditsPage() {
  const { hasPermission } = usePermissions()

  return (
    <View>
      {/* 设置按钮 - 需要 game:room_credit:set 权限 */}
      {hasPermission('game:room_credit:set') && (
        <Button onPress={handleSetCredit}>
          设置额度
        </Button>
      )}

      {/* 删除按钮 - 需要 game:room_credit:delete 权限 */}
      {hasPermission('game:room_credit:delete') && (
        <Button onPress={handleDelete} color="red">
          删除
        </Button>
      )}

      {/* 检查按钮 - 需要 game:room_credit:check 权限 */}
      {hasPermission('game:room_credit:check') && (
        <Button onPress={handleCheck}>
          检查玩家
        </Button>
      )}
    </View>
  )
}
```

## API接口对应权限

| API接口 | 方法 | 需要的权限 |
|---------|------|------------|
| `/room-credit/list` | POST | `game:room_credit:view` |
| `/room-credit/get` | POST | `game:room_credit:view` |
| `/room-credit/set` | POST | `game:room_credit:set` |
| `/room-credit/delete` | POST | `game:room_credit:delete` |
| `/room-credit/check` | POST | `game:room_credit:check` |

## 推荐角色权限组合

### 超级管理员
```sql
-- 拥有所有权限
game:room_credit:view
game:room_credit:set
game:room_credit:delete
game:room_credit:check
```

### 店铺管理员
```sql
-- 拥有所有权限
game:room_credit:view
game:room_credit:set
game:room_credit:delete
game:room_credit:check
```

### 普通用户
```sql
-- 只能查看
game:room_credit:view
```

## 验证配置是否成功

### 1. 检查权限是否创建
```sql
SELECT * FROM basic_permission WHERE code LIKE 'game:room_credit:%';
```

### 2. 检查菜单是否创建
```sql
SELECT id, title, name, path, auths FROM basic_menu WHERE name = 'shop.room_credits';
```

### 3. 检查按钮权限是否创建
```sql
SELECT * FROM basic_menu_button WHERE menu_id = 64;
```

### 4. 检查角色是否有权限
```sql
-- 查看角色1（超级管理员）的权限
SELECT p.code, p.name 
FROM basic_permission p
JOIN basic_role_permission_rel rpr ON p.id = rpr.permission_id
WHERE rpr.role_id = 1 AND p.code LIKE 'game:room_credit:%';
```

### 5. 检查角色是否有菜单
```sql
-- 查看角色1的菜单
SELECT m.id, m.title, m.name 
FROM basic_menu m
JOIN basic_role_menu_rel rmr ON m.id = rmr.menu_id
WHERE rmr.role_id = 1 AND m.name = 'shop.room_credits';
```

## 前端获取菜单API

根据系统内存，前端调用 `/basic/baseMenu/getOption` 获取菜单列表时需要注意：

```typescript
// 获取菜单列表
const response = await fetch('/basic/baseMenu/getOption', {
  method: 'POST',
  body: JSON.stringify({
    page: 1,
    page_size: 1000  // 足够大以获取所有菜单
  })
})
```

**注意**：参数使用 `snake_case` 格式（`page_size`），不是 `camelCase`。

## 常见问题

### Q: 为什么我看不到"房间额度"菜单？
A: 检查以下几点：
1. 菜单是否已创建（查询 `basic_menu`）
2. 你的角色是否有此菜单权限（查询 `basic_role_menu_rel`）
3. 你的账号是否属于该角色（查询 `basic_user_role_rel`）
4. 前端是否已刷新菜单数据

### Q: 为什么按钮不显示？
A: 检查：
1. 按钮权限是否已创建（查询 `basic_menu_button`）
2. 你的角色是否有对应的权限（查询 `basic_role_permission_rel`）
3. 前端权限判断逻辑是否正确

### Q: 如何临时给用户某个权限？
A: 
1. 创建临时角色，分配该权限
2. 将用户分配到该角色
3. 或者直接修改角色的权限配置

## 总结

执行顺序：
1. ✅ 运行 `add_room_credit_limit_menu.sql`
2. ✅ 查询权限ID
3. ✅ 为角色分配菜单（通过后台或SQL）
4. ✅ 为角色分配权限（通过后台或SQL）
5. ✅ 前端开发对应页面
6. ✅ 测试权限控制

配置完成后，用户就可以根据其角色权限访问"房间额度"功能了！
