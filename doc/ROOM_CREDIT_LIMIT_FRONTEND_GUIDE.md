# 房间额度限制功能 - 前端配置指南

## 一、菜单配置

### 1.1 菜单层级结构

建议将"房间额度限制"菜单放在"游戏管理"或"店铺管理"下：

```
游戏管理
  ├── 店铺管理
  ├── 圈子管理
  ├── 成员管理
  ├── 房间额度限制  ⭐ 新增
  └── 战绩查询
```

### 1.2 菜单数据（SQL）

需要在 `basic_menu` 表中插入菜单数据：

```sql
-- 1. 插入父级菜单（如果没有"游戏管理"菜单）
INSERT INTO basic_menu (id, parent_id, menu_type, menu_name, route_name, route_path, component, icon, sort, visible, status, perms, created_at, updated_at)
VALUES 
(200, -1, 'M', '游戏管理', '', '/game', 'Layout', 'game', 200, true, true, '', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. 插入"房间额度限制"菜单
INSERT INTO basic_menu (id, parent_id, menu_type, menu_name, route_name, route_path, component, icon, sort, visible, status, perms, created_at, updated_at)
VALUES 
(210, 200, 'C', '房间额度限制', 'RoomCreditLimit', '/game/room-credit-limit', '/game/room-credit-limit/index', 'credit-card', 210, true, true, 'room:credit:view', NOW(), NOW());

-- 3. 插入按钮权限
INSERT INTO basic_menu (id, parent_id, menu_type, menu_name, route_name, route_path, component, icon, sort, visible, status, perms, created_at, updated_at)
VALUES 
-- 查看权限
(2101, 210, 'F', '查看额度', '', '', '', '', 1, true, true, 'room:credit:view', NOW(), NOW()),
-- 设置权限
(2102, 210, 'F', '设置额度', '', '', '', '', 2, true, true, 'room:credit:set', NOW(), NOW()),
-- 删除权限
(2103, 210, 'F', '删除额度', '', '', '', '', 3, true, true, 'room:credit:delete', NOW(), NOW()),
-- 检查权限
(2104, 210, 'F', '检查玩家额度', '', '', '', '', 4, true, true, 'room:credit:check', NOW(), NOW());
```

### 1.3 字段说明

| 字段 | 值 | 说明 |
|------|-----|------|
| `menu_type` | `M` | Menu - 目录菜单 |
| `menu_type` | `C` | Component - 页面菜单 |
| `menu_type` | `F` | Function - 按钮权限 |
| `route_name` | `RoomCreditLimit` | Vue Router 路由名称 |
| `route_path` | `/game/room-credit-limit` | URL 路径 |
| `component` | `/game/room-credit-limit/index` | 组件路径 |
| `icon` | `credit-card` | 图标（根据你的图标库调整） |
| `perms` | `room:credit:view` | 权限标识 |

## 二、权限标识列表

### 2.1 后端权限定义

| 权限标识 | 说明 | 对应API |
|----------|------|---------|
| `room:credit:view` | 查看房间额度 | `GET /room-credit/list`, `GET /room-credit/get` |
| `room:credit:set` | 设置房间额度 | `POST /room-credit/set` |
| `room:credit:delete` | 删除房间额度 | `POST /room-credit/delete` |
| `room:credit:check` | 检查玩家额度 | `POST /room-credit/check` |

### 2.2 前端权限使用

在 Vue 组件中使用权限指令：

```vue
<template>
  <div>
    <!-- 设置按钮 -->
    <el-button 
      v-has-perms="'room:credit:set'"
      type="primary"
      @click="handleAdd">
      设置额度
    </el-button>

    <!-- 删除按钮 -->
    <el-button
      v-has-perms="'room:credit:delete'"
      type="danger"
      @click="handleDelete">
      删除
    </el-button>

    <!-- 检查按钮 -->
    <el-button
      v-has-perms="'room:credit:check'"
      @click="handleCheck">
      检查玩家
    </el-button>
  </div>
</template>
```

## 三、前端路由配置

### 3.1 路由定义（Vue Router）

在前端路由文件中添加：

```typescript
// src/router/modules/game.ts

export default {
  path: '/game',
  component: Layout,
  meta: { title: '游戏管理', icon: 'game' },
  children: [
    // ... 其他路由
    {
      path: 'room-credit-limit',
      name: 'RoomCreditLimit',
      component: () => import('@/views/game/room-credit-limit/index.vue'),
      meta: { 
        title: '房间额度限制', 
        icon: 'credit-card',
        perms: ['room:credit:view'] 
      }
    }
  ]
}
```

## 四、前端页面组件结构

### 4.1 页面文件结构

```
src/views/game/room-credit-limit/
├── index.vue              # 主页面
├── components/
│   ├── CreditLimitForm.vue    # 设置/编辑表单
│   ├── CreditLimitTable.vue   # 列表表格
│   └── PlayerCheckDialog.vue  # 检查玩家对话框
└── types.ts               # TypeScript 类型定义
```

### 4.2 主页面示例 (index.vue)

```vue
<template>
  <div class="room-credit-limit-container">
    <el-card shadow="never">
      <!-- 查询表单 -->
      <el-form :inline="true" :model="queryForm">
        <el-form-item label="店铺">
          <el-select v-model="queryForm.house_gid" placeholder="请选择店铺">
            <el-option 
              v-for="house in houseList" 
              :key="house.id" 
              :label="house.name" 
              :value="house.id" 
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="圈子">
          <el-input v-model="queryForm.group_name" placeholder="留空查询全部" clearable />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作按钮 -->
    <el-card shadow="never" style="margin-top: 16px;">
      <el-button 
        v-has-perms="'room:credit:set'"
        type="primary" 
        icon="Plus"
        @click="handleAdd">
        新增额度规则
      </el-button>
      
      <el-button 
        v-has-perms="'room:credit:check'"
        icon="User"
        @click="handleCheckPlayer">
        检查玩家额度
      </el-button>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" style="margin-top: 16px;">
      <el-table :data="tableData" border stripe>
        <el-table-column prop="house_gid" label="店铺GID" width="100" />
        <el-table-column prop="group_name" label="圈子名称" width="120">
          <template #default="{ row }">
            {{ row.group_name || '全局' }}
          </template>
        </el-table-column>
        <el-table-column prop="game_kind_name" label="游戏类型" width="100">
          <template #default="{ row }">
            {{ row.game_kind === 0 ? '全部' : row.game_kind_name }}
          </template>
        </el-table-column>
        <el-table-column prop="base_score" label="底分" width="80">
          <template #default="{ row }">
            {{ row.base_score === 0 ? '全部' : row.base_score }}
          </template>
        </el-table-column>
        <el-table-column prop="credit_yuan" label="额度（元）" width="120">
          <template #default="{ row }">
            <el-tag type="warning">{{ row.credit_yuan }}元</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              v-has-perms="'room:credit:set'"
              type="primary"
              link
              size="small"
              @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              v-has-perms="'room:credit:delete'"
              type="danger"
              link
              size="small"
              @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleQuery"
        @current-change="handleQuery"
      />
    </el-card>

    <!-- 设置/编辑对话框 -->
    <CreditLimitForm
      v-model:visible="formVisible"
      :form-data="formData"
      :is-edit="isEdit"
      @success="handleQuery"
    />

    <!-- 检查玩家对话框 -->
    <PlayerCheckDialog
      v-model:visible="checkDialogVisible"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listRoomCreditLimits, deleteRoomCreditLimit } from '@/api/game/roomCredit'
import CreditLimitForm from './components/CreditLimitForm.vue'
import PlayerCheckDialog from './components/PlayerCheckDialog.vue'

// 查询表单
const queryForm = reactive({
  house_gid: null,
  group_name: ''
})

// 表格数据
const tableData = ref([])

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 表单对话框
const formVisible = ref(false)
const formData = ref({})
const isEdit = ref(false)

// 检查玩家对话框
const checkDialogVisible = ref(false)

// 店铺列表（从其他接口获取）
const houseList = ref([])

// 查询
const handleQuery = async () => {
  try {
    const { data } = await listRoomCreditLimits({
      house_gid: queryForm.house_gid,
      group_name: queryForm.group_name
    })
    tableData.value = data.items
    pagination.total = data.total
  } catch (error) {
    ElMessage.error('查询失败')
  }
}

// 重置
const handleReset = () => {
  queryForm.house_gid = null
  queryForm.group_name = ''
  handleQuery()
}

// 新增
const handleAdd = () => {
  formData.value = {}
  isEdit.value = false
  formVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  formData.value = { ...row }
  isEdit.value = true
  formVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条额度规则吗？', '提示', {
      type: 'warning'
    })
    
    await deleteRoomCreditLimit({
      house_gid: row.house_gid,
      group_name: row.group_name,
      game_kind: row.game_kind,
      base_score: row.base_score
    })
    
    ElMessage.success('删除成功')
    handleQuery()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 检查玩家
const handleCheckPlayer = () => {
  checkDialogVisible.value = true
}

onMounted(() => {
  handleQuery()
})
</script>

<style scoped lang="scss">
.room-credit-limit-container {
  padding: 16px;
}
</style>
```

### 4.3 API 接口定义 (TypeScript)

```typescript
// src/api/game/roomCredit.ts

import request from '@/utils/request'

// 设置房间额度限制
export function setRoomCreditLimit(data: {
  house_gid: number
  group_name?: string
  game_kind?: number
  base_score?: number
  credit_limit: number
}) {
  return request({
    url: '/room-credit/set',
    method: 'post',
    data
  })
}

// 查询房间额度限制列表
export function listRoomCreditLimits(params: {
  house_gid: number
  group_name?: string
}) {
  return request({
    url: '/room-credit/list',
    method: 'post',
    data: params
  })
}

// 删除房间额度限制
export function deleteRoomCreditLimit(data: {
  house_gid: number
  group_name: string
  game_kind: number
  base_score: number
}) {
  return request({
    url: '/room-credit/delete',
    method: 'post',
    data
  })
}

// 检查玩家额度
export function checkPlayerCredit(data: {
  house_gid: number
  game_id: number
  group_name?: string
  game_kind: number
  base_score: number
}) {
  return request({
    url: '/room-credit/check',
    method: 'post',
    data
  })
}
```

## 五、角色权限配置

### 5.1 为角色分配权限（SQL）

```sql
-- 示例：为"店铺管理员"角色分配房间额度管理权限
-- 假设角色ID为 5

-- 分配菜单权限
INSERT INTO basic_role_menu_rel (role_id, menu_id)
VALUES 
(5, 210),  -- 房间额度限制页面
(5, 2101), -- 查看权限
(5, 2102), -- 设置权限
(5, 2103), -- 删除权限
(5, 2104); -- 检查权限

-- 或者通过管理后台界面分配
```

### 5.2 推荐的角色权限组合

| 角色 | 推荐权限 |
|------|---------|
| 超级管理员 | 全部权限 |
| 店铺管理员 | `room:credit:view`, `room:credit:set`, `room:credit:delete`, `room:credit:check` |
| 圈主 | `room:credit:view`, `room:credit:set` (仅限自己圈子) |
| 普通用户 | `room:credit:view` (只读) |

## 六、使用示例

### 6.1 设置全局默认额度

```json
{
  "house_gid": 123456,
  "group_name": "",
  "game_kind": 0,
  "base_score": 0,
  "credit_limit": 5000  // 50元
}
```

### 6.2 设置特定游戏额度

```json
{
  "house_gid": 123456,
  "group_name": "",
  "game_kind": 5,      // 红中
  "base_score": 1,
  "credit_limit": 10000 // 100元
}
```

### 6.3 设置圈子默认额度

```json
{
  "house_gid": 123456,
  "group_name": "VIP圈",
  "game_kind": 0,
  "base_score": 0,
  "credit_limit": 20000 // 200元
}
```

## 七、前端显示格式

### 7.1 额度显示格式

参考 passing-dragonfly 的显示格式：

- 全局默认：`🈲 50`
- 特定游戏：`🈲 100/红中/1`
- 圈子默认：`🈲 200/VIP圈`
- 圈子+游戏：`🈲 300/VIP圈/红中/2`

### 7.2 格式化函数

```typescript
// utils/format.ts

export function formatCreditDisplay(
  creditLimit: number,
  gameKindName?: string,
  baseScore?: number,
  groupName?: string
): string {
  const yuan = creditLimit / 100
  
  if (groupName && gameKindName && baseScore) {
    return `🈲 ${yuan}/${groupName}/${gameKindName}/${baseScore}`
  } else if (groupName && (!gameKindName || !baseScore)) {
    return `🈲 ${yuan}/${groupName}`
  } else if (gameKindName && baseScore) {
    return `🈲 ${yuan}/${gameKindName}/${baseScore}`
  } else {
    return `🈲 ${yuan}`
  }
}
```

## 八、注意事项

1. **权限标识**必须与后端定义一致
2. **菜单ID**不要与现有菜单冲突
3. 前端需要实现**权限指令** (`v-has-perms`)
4. 建议在设置额度时提供**游戏类型选择器**
5. 建议添加**额度规则说明**，帮助用户理解优先级

## 九、游戏类型映射表

需要维护游戏类型ID与名称的映射关系：

```typescript
// constants/gameKinds.ts

export const GAME_KINDS = {
  0: '全部',
  5: '红中',
  6: '跑得快',
  7: '二人麻将',
  // ... 其他游戏类型
}

export const GAME_KIND_OPTIONS = Object.entries(GAME_KINDS).map(([value, label]) => ({
  value: Number(value),
  label
}))
```
