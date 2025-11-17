# Toast 提示系统使用指南

## 概述

项目现在使用统一的 `toast` 提示系统替代了之前的 `alert`。新系统提供了更丰富的功能和更好的用户体验。

## 特性

- ✅ **多种类型**: success, error, warning, info
- ✅ **确认对话框**: 支持异步操作
- ✅ **自定义样式**: 可自定义图标、位置、持续时间
- ✅ **响应式动画**: 优雅的进入/退出动画
- ✅ **全局调用**: 在任何地方都可以调用
- ✅ **类型安全**: 完整的 TypeScript 支持

## 基础用法

### 1. 成功提示

```typescript
import { toast } from '@/utils/toast';

// 简单用法
toast.success('操作成功');

// 带描述
toast.success('保存成功', '您的更改已保存');
```

### 2. 错误提示

```typescript
// 简单用法
toast.error('操作失败');

// 带详细描述
toast.error('删除失败', '该项目正在使用中，无法删除');
```

### 3. 警告提示

```typescript
toast.warning('警告', '您的会话即将过期');
```

### 4. 信息提示

```typescript
toast.info('提示', '新版本可用');
```

## 确认对话框

### 基础确认框

```typescript
import { toast } from '@/utils/toast';

toast.confirm({
  title: '确认删除',
  description: '此操作无法撤销，确定要删除吗？',
  confirmText: '删除',
  cancelText: '取消',
  confirmVariant: 'destructive',
  onConfirm: () => {
    // 执行删除操作
    console.log('已删除');
  },
  onCancel: () => {
    console.log('已取消');
  },
});
```

### 异步确认框

```typescript
toast.confirm({
  title: '确认提交',
  description: '确定要提交这些更改吗？',
  confirmText: '提交',
  cancelText: '取消',
  onConfirm: async () => {
    // 执行异步操作
    await submitChanges();
    toast.success('提交成功');
  },
});
```

## 高级用法

### 自定义配置

```typescript
import { CheckCircle } from 'lucide-react-native';

toast.show({
  title: '自定义提示',
  description: '这是一个自定义的提示',
  type: 'success',
  icon: CheckCircle,
  duration: 5000, // 5秒后自动关闭
  position: 'top', // 'top' | 'center' | 'bottom'
  showClose: true, // 显示关闭按钮
});
```

### 不自动关闭

```typescript
toast.show({
  title: '重要通知',
  description: '请仔细阅读此信息',
  duration: 0, // 0 表示不自动关闭
  showClose: true, // 必须手动关闭
});
```

### 手动关闭

```typescript
// 关闭当前显示的 toast
toast.close();
```

## 在组件中使用

### 使用 Hook

```typescript
import { useToast } from '@/components/shared/toast-provider';

function MyComponent() {
  const toast = useToast();

  const handleSave = () => {
    // 执行保存操作
    toast.success('保存成功');
  };

  return <Button onPress={handleSave}>保存</Button>;
}
```

## 实际示例

### 示例 1: 表单提交

```typescript
const handleSubmit = async () => {
  try {
    await submitForm(data);
    toast.success('提交成功', '您的表单已成功提交');
  } catch (error) {
    toast.error('提交失败', error.message);
  }
};
```

### 示例 2: 删除确认

```typescript
const handleDelete = (id: number) => {
  toast.confirm({
    title: '确认删除',
    description: '删除后无法恢复，确定要继续吗？',
    type: 'warning',
    confirmText: '删除',
    confirmVariant: 'destructive',
    cancelText: '取消',
    onConfirm: async () => {
      try {
        await deleteItem(id);
        toast.success('删除成功');
      } catch (error) {
        toast.error('删除失败', error.message);
      }
    },
  });
};
```

### 示例 3: 网络请求错误处理

```typescript
try {
  const response = await fetch('/api/data');
  if (!response.ok) {
    throw new Error('请求失败');
  }
  const data = await response.json();
  toast.success('加载成功');
  return data;
} catch (error) {
  toast.error('网络错误', '无法连接到服务器，请检查网络连接');
  throw error;
}
```

## API 参考

### ToastOptions

| 属性 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `title` | `string` | - | 标题（必填） |
| `description` | `string` | - | 描述文本 |
| `type` | `'success' \| 'error' \| 'warning' \| 'info'` | `'info'` | 提示类型 |
| `icon` | `LucideIcon` | - | 自定义图标 |
| `duration` | `number` | `3000` | 显示时长(毫秒)，0表示不自动关闭 |
| `position` | `'top' \| 'center' \| 'bottom'` | `'top'` | 显示位置 |
| `showClose` | `boolean` | `false` | 是否显示关闭按钮 |

### ConfirmOptions

| 属性 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `title` | `string` | - | 标题（必填） |
| `description` | `string` | - | 描述文本 |
| `type` | `'success' \| 'error' \| 'warning' \| 'info'` | `'warning'` | 提示类型 |
| `icon` | `LucideIcon` | - | 自定义图标 |
| `confirmText` | `string` | `'确定'` | 确认按钮文本 |
| `cancelText` | `string` | `'取消'` | 取消按钮文本 |
| `confirmVariant` | `'default' \| 'destructive'` | `'default'` | 确认按钮样式 |
| `onConfirm` | `() => void \| Promise<void>` | - | 确认回调 |
| `onCancel` | `() => void` | - | 取消回调 |

## 迁移指南

### 从 alert 迁移到 toast

**旧代码:**
```typescript
import { alert } from '@/utils/alert';

alert.show({ title: '成功', description: '操作成功' });
```

**新代码:**
```typescript
import { toast } from '@/utils/toast';

toast.success('操作成功');
```

**确认框迁移:**

**旧代码:**
```typescript
alert.show({
  title: '确认删除',
  description: '确定要删除吗？',
  confirmText: '删除',
  cancelText: '取消',
  onConfirm: () => { /* ... */ },
  onCancel: () => { /* ... */ },
});
```

**新代码:**
```typescript
toast.confirm({
  title: '确认删除',
  description: '确定要删除吗？',
  confirmText: '删除',
  cancelText: '取消',
  confirmVariant: 'destructive',
  onConfirm: () => { /* ... */ },
  onCancel: () => { /* ... */ },
});
```

## 最佳实践

1. **选择合适的类型**: 根据操作结果选择 success/error/warning/info
2. **简洁的文案**: 标题简短，描述清晰
3. **异步操作**: 确认框的 `onConfirm` 可以是异步函数
4. **错误处理**: 在 try-catch 中使用 toast.error
5. **用户反馈**: 重要操作后总是给用户反馈
6. **避免滥用**: 不要为每个小操作都显示提示

## 注意事项

- Toast 同时只显示一个，新的会替换旧的
- 确认框会阻止用户操作，直到做出选择
- 异步 onConfirm 会显示 "处理中..." 状态
- duration 为 0 时必须手动关闭（showClose: true）

