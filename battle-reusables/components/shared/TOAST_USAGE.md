# Toast 提示组件使用指南

## 概述

Toast 是一个统一的消息提示组件系统，用于在应用中显示成功、错误、警告和信息提示。

## 组件结构

### 1. ToastProvider（提供者组件）
位置：`/components/shared/toast-provider.tsx`

这是 Toast 系统的核心组件，需要在应用的根组件中使用。

```tsx
import { ToastProvider } from '@/components/shared/toast-provider';

export default function App() {
  return (
    <ToastProvider>
      {/* 你的应用内容 */}
    </ToastProvider>
  );
}
```

### 2. 全局 Toast API
位置：`/utils/toast.ts`

提供全局访问的 toast 方法。

## 使用方法

### 方法一：使用 `showToast` 函数（推荐用于简单场景）

```tsx
import { showToast } from '@/utils/toast';

// 成功提示
showToast('操作成功', 'success');

// 错误提示
showToast('操作失败', 'error');

// 警告提示
showToast('请注意', 'warning');

// 信息提示
showToast('提示信息', 'info');
```

### 方法二：使用 `toast` 对象（更灵活）

```tsx
import { toast } from '@/utils/toast';

// 成功提示（带描述）
toast.success('保存成功', '数据已成功保存到服务器');

// 错误提示
toast.error('删除失败', '没有权限执行此操作');

// 自定义提示
toast.show({
  title: '自定义提示',
  description: '这是详细描述',
  type: 'info',
  duration: 5000, // 显示时长（毫秒）
  position: 'top', // 位置: top/center/bottom
  showClose: true, // 显示关闭按钮
});
```

### 方法三：使用 Hook（用于组件内部）

```tsx
import { useToast } from '@/components/shared/toast-provider';

function MyComponent() {
  const toast = useToast();

  const handleSubmit = async () => {
    try {
      await submitData();
      toast.success('提交成功');
    } catch (error) {
      toast.error('提交失败', error.message);
    }
  };

  return <button onClick={handleSubmit}>提交</button>;
}
```

## CRUD 操作中的标准用法

### 创建（Create）

```tsx
const handleCreate = async () => {
  try {
    const res = await createItem(data);
    if (res.code === 0) {
      showToast('创建成功', 'success');
      onSuccess();
      onClose();
    } else {
      showToast(res.msg || '创建失败', 'error');
    }
  } catch (error) {
    showToast('创建失败', 'error');
    console.error('Create error:', error);
  }
};
```

### 更新（Update）

```tsx
const handleUpdate = async () => {
  try {
    const res = await updateItem(id, data);
    if (res.code === 0) {
      showToast('更新成功', 'success');
      onSuccess();
      onClose();
    } else {
      showToast(res.msg || '更新失败', 'error');
    }
  } catch (error) {
    showToast('更新失败', 'error');
    console.error('Update error:', error);
  }
};
```

### 删除（Delete）

```tsx
const handleDelete = async (id: number) => {
  try {
    const res = await deleteItem(id);
    if (res.code === 0) {
      showToast('删除成功', 'success');
      loadData(); // 重新加载列表
    } else {
      showToast(res.msg || '删除失败', 'error');
    }
  } catch (error) {
    showToast('删除失败', 'error');
    console.error('Delete error:', error);
  }
};
```

### 表单验证

```tsx
const handleSubmit = async () => {
  if (!name.trim()) {
    showToast('请输入名称', 'error');
    return;
  }
  if (!email.trim()) {
    showToast('请输入邮箱', 'error');
    return;
  }
  // 继续处理...
};
```

## 确认对话框

使用 `toast.confirm` 显示确认对话框（适用于危险操作）：

```tsx
import { toast } from '@/utils/toast';

const handleDelete = (item) => {
  toast.confirm({
    title: '确认删除',
    description: `确定要删除"${item.name}"吗？此操作无法撤销。`,
    type: 'warning',
    confirmText: '删除',
    cancelText: '取消',
    confirmVariant: 'destructive',
    onConfirm: async () => {
      await deleteItem(item.id);
      showToast('删除成功', 'success');
    },
    onCancel: () => {
      console.log('取消删除');
    },
  });
};
```

## Toast 类型说明

| 类型 | 说明 | 默认持续时间 | 图标 | 颜色 |
|------|------|-------------|------|------|
| `success` | 成功提示 | 3000ms | ✓ | 绿色 |
| `error` | 错误提示 | 4000ms | ✕ | 红色 |
| `warning` | 警告提示 | 3500ms | ⚠ | 黄色 |
| `info` | 信息提示 | 3000ms | ℹ | 蓝色 |

## 配置选项

### ToastOptions

```typescript
{
  title: string;           // 标题（必需）
  description?: string;    // 描述（可选）
  type?: ToastType;        // 类型（可选，默认 'info'）
  icon?: LucideIcon;       // 自定义图标（可选）
  duration?: number;       // 持续时间，0 为不自动关闭（可选）
  showClose?: boolean;     // 显示关闭按钮（可选）
  position?: 'top' | 'center' | 'bottom'; // 位置（可选，默认 'top'）
}
```

## 最佳实践

### ✅ 推荐做法

1. **操作成功后显示提示**
   ```tsx
   showToast('保存成功', 'success');
   ```

2. **操作失败时显示错误**
   ```tsx
   showToast(res.msg || '操作失败', 'error');
   ```

3. **表单验证失败时提示**
   ```tsx
   showToast('请填写完整信息', 'error');
   ```

4. **异步操作的完整处理**
   ```tsx
   try {
     const res = await operation();
     if (res.code === 0) {
       showToast('成功', 'success');
     } else {
       showToast(res.msg || '失败', 'error');
     }
   } catch (error) {
     showToast('操作失败', 'error');
   }
   ```

### ❌ 避免的做法

1. **不要重复显示提示**
   ```tsx
   // 错误：重复提示
   showToast('保存中...', 'info');
   await save();
   showToast('保存完成', 'info');
   showToast('保存成功', 'success'); // 多余
   ```

2. **不要使用过长的文本**
   ```tsx
   // 错误：文本过长
   showToast('这是一个非常非常非常长的提示消息，用户可能看不完...', 'info');
   ```

3. **不要忘记检查响应状态**
   ```tsx
   // 错误：没有检查 code
   const res = await updateItem(data);
   showToast('更新成功', 'success'); // 可能实际失败了
   ```

## 注意事项

1. **API 响应字段名**：使用 `res.msg` 而非 `res.message`
2. **响应状态码**：检查 `res.code === 0` 表示成功
3. **异常处理**：始终使用 try-catch 包裹异步操作
4. **日志记录**：错误时使用 `console.error` 记录详细信息

## 示例：完整的表单组件

```tsx
import React, { useState } from 'react';
import { showToast } from '@/utils/toast';
import { createItem, updateItem } from '@/services/api';

function ItemForm({ item, onSuccess }) {
  const [name, setName] = useState(item?.name || '');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    // 表单验证
    if (!name.trim()) {
      showToast('请输入名称', 'error');
      return;
    }

    setLoading(true);
    try {
      const data = { name };
      const res = item
        ? await updateItem(item.id, data)
        : await createItem(data);

      if (res.code === 0) {
        showToast(item ? '更新成功' : '创建成功', 'success');
        onSuccess();
      } else {
        showToast(res.msg || '操作失败', 'error');
      }
    } catch (error) {
      showToast(item ? '更新失败' : '创建失败', 'error');
      console.error('Form submit error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <View>
      <TextInput value={name} onChangeText={setName} />
      <Button onPress={handleSubmit} disabled={loading}>
        {loading ? '提交中...' : '提交'}
      </Button>
    </View>
  );
}
```

## 支持与贡献

如有问题或建议，请联系开发团队。
