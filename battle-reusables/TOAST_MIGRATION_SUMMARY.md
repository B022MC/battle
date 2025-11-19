# Toast 系统迁移完成总结

## 已完成工作

### 1. 新系统创建 ✅

- ✅ 创建 `utils/toast.ts` - Toast API 定义和全局实例
- ✅ 创建 `components/shared/toast-provider.tsx` - Toast 提供者组件
- ✅ 创建 `TOAST_USAGE.md` - 完整使用文档
- ✅ 更新 `app/_layout.tsx` - 集成 ToastProvider

### 2. 核心文件迁移 ✅

- ✅ `utils/request.ts` - HTTP 请求错误处理
- ✅ `components/(tabs)/profile/profile-view.tsx` - 个人资料页
- ✅ `components/(tabs)/profile/profile-game-account.tsx` - 游戏账号绑定

### 3. 待迁移文件 📋

以下文件仍在使用旧的 `alert` 系统,需要迁移:

1. `components/(tabs)/profile/profile-ctrl-accounts.tsx`
2. `components/(tabs)/profile/profile-ctrl-account.tsx`
3. `components/(tabs)/members/members-view.tsx`
4. `components/(shop)/battles/my-battles-view.tsx`
5. `components/(shop)/battles/my-balances-view.tsx`
6. `components/(shop)/battles/group-battles-view.tsx`
7. `components/(shop)/battles/group-balances-view.tsx`
8. `components/(shop)/admins/admins-assign.tsx`

## 新系统特性

### 核心功能

```typescript
import { toast } from '@/utils/toast';

// 1. 成功提示
toast.success('操作成功');
toast.success('保存成功', '数据已保存');

// 2. 错误提示
toast.error('操作失败');
toast.error('删除失败', '该项目正在使用中');

// 3. 警告提示
toast.warning('警告', '您的会话即将过期');

// 4. 信息提示
toast.info('提示', '新版本可用');

// 5. 确认对话框
toast.confirm({
  title: '确认删除',
  description: '此操作无法撤销',
  confirmText: '删除',
  cancelText: '取消',
  confirmVariant: 'destructive',
  onConfirm: async () => {
    await deleteItem();
    toast.success('删除成功');
  },
});
```

### 高级配置

```typescript
// 自定义显示时长和位置
toast.show({
  title: '自定义提示',
  description: '这是一个自定义的提示',
  type: 'success',
  duration: 5000, // 5秒
  position: 'center', // top | center | bottom
  showClose: true, // 显示关闭按钮
});

// 永不自动关闭
toast.show({
  title: '重要通知',
  description: '请仔细阅读',
  duration: 0, // 不自动关闭
  showClose: true,
});
```

## 迁移示例

### 示例 1: 简单提示

```typescript
// 旧代码
alert.show({ title: '成功', description: '操作成功' });

// 新代码
toast.success('操作成功');
```

### 示例 2: 错误提示

```typescript
// 旧代码
alert.show({ title: '错误', description: '操作失败' });

// 新代码
toast.error('操作失败');
```

### 示例 3: 确认对话框

```typescript
// 旧代码
alert.show({
  title: '确认删除',
  description: '确定要删除吗?',
  confirmText: '删除',
  cancelText: '取消',
  onConfirm: () => { /* 删除逻辑 */ },
  onCancel: () => { /* 取消逻辑 */ },
});

// 新代码
toast.confirm({
  title: '确认删除',
  description: '确定要删除吗?',
  confirmText: '删除',
  cancelText: '取消',
  confirmVariant: 'destructive', // 新增:危险操作样式
  type: 'warning', // 新增:提示类型
  onConfirm: () => { /* 删除逻辑 */ },
  onCancel: () => { /* 取消逻辑 */ },
});
```

### 示例 4: useRequest 中使用

```typescript
// 旧代码
const { run } = useRequest(api, {
  manual: true,
  onSuccess: () => {
    alert.show({ title: '成功', description: '保存成功' });
  },
});

// 新代码
const { run } = useRequest(api, {
  manual: true,
  onSuccess: () => {
    toast.success('保存成功');
  },
});
```

### 示例 5: 异步确认操作

```typescript
// 新系统支持异步 onConfirm
toast.confirm({
  title: '确认提交',
  description: '确定要提交这些更改吗?',
  confirmText: '提交',
  cancelText: '取消',
  onConfirm: async () => {
    // 执行异步操作,会自动显示"处理中..."
    await submitChanges();
    // 成功后显示提示
    toast.success('提交成功');
  },
});
```

## 迁移步骤

### 对于每个待迁移文件:

1. **更新导入语句**
   ```typescript
   // 替换
   import { alert } from '@/utils/alert';
   // 为
   import { toast } from '@/utils/toast';
   ```

2. **替换简单提示**
   - 成功: `alert.show({ title: 'X', description: 'Y' })` → `toast.success('X', 'Y')`
   - 错误: 同上使用 `toast.error()`
   - 警告: 同上使用 `toast.warning()`

3. **替换确认对话框**
   - 添加 `confirmVariant` 参数(删除操作使用 `'destructive'`)
   - 添加 `type` 参数(通常为 `'warning'`)

4. **测试功能**
   - 验证所有提示正常显示
   - 验证确认对话框逻辑正确
   - 验证异步操作显示"处理中..."状态

## 改进点

### 相比旧系统的改进:

1. ✅ **更丰富的类型** - success/error/warning/info 四种类型,带颜色区分
2. ✅ **更好的动画** - 优雅的滑入/滑出动画
3. ✅ **灵活的位置** - 支持 top/center/bottom 三种位置
4. ✅ **异步支持** - 确认对话框支持异步操作,自动显示加载状态
5. ✅ **类型安全** - 完整的 TypeScript 类型定义
6. ✅ **更好的样式** - 现代化的 UI 设计,支持深色模式
7. ✅ **手动关闭** - 支持永不自动关闭 + 手动关闭按钮
8. ✅ **更简洁的 API** - 快捷方法减少代码量

## 下一步

1. 完成剩余 8 个文件的迁移
2. 测试所有页面的提示功能
3. 删除旧的 `alert` 系统相关文件:
   - `utils/alert.ts`
   - `components/shared/alert-provider.tsx`
   - `components/ui/alert.tsx` (如果只用于 AlertProvider)
4. 更新项目文档

## 注意事项

- ⚠️ 新系统同时只显示一个 Toast/确认框
- ⚠️ 异步 onConfirm 在处理期间会禁用按钮
- ⚠️ duration=0 时必须设置 showClose=true
- ⚠️ 删除操作建议使用 confirmVariant='destructive'

