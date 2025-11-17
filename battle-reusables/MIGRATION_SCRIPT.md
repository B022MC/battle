# Alert 到 Toast 迁移脚本

## 需要替换的文件列表

以下文件需要从 `alert` 迁移到 `toast`:

1. ✅ components/(tabs)/profile/profile-view.tsx
2. components/(tabs)/profile/profile-game-account.tsx
3. components/(tabs)/profile/profile-ctrl-accounts.tsx
4. components/(tabs)/profile/profile-ctrl-account.tsx
5. components/(tabs)/members/members-view.tsx
6. components/(shop)/battles/my-battles-view.tsx
7. components/(shop)/battles/my-balances-view.tsx
8. components/(shop)/battles/group-battles-view.tsx
9. components/(shop)/battles/group-balances-view.tsx
10. components/(shop)/admins/admins-assign.tsx
11. ✅ utils/request.ts
12. ✅ app/_layout.tsx

## 替换规则

### 1. 导入语句
```typescript
// 旧
import { alert } from '@/utils/alert';

// 新
import { toast } from '@/utils/toast';
```

### 2. 简单提示
```typescript
// 旧
alert.show({ title: '成功', description: '操作成功' });

// 新
toast.success('操作成功');
// 或
toast.success('成功', '操作成功');
```

### 3. 错误提示
```typescript
// 旧
alert.show({ title: '错误', description: '操作失败' });

// 新
toast.error('操作失败');
// 或
toast.error('错误', '操作失败');
```

### 4. 确认对话框
```typescript
// 旧
alert.show({
  title: '确认删除',
  description: '确定要删除吗?',
  confirmText: '删除',
  cancelText: '取消',
  onConfirm: () => { /* ... */ },
  onCancel: () => { /* ... */ },
});

// 新
toast.confirm({
  title: '确认删除',
  description: '确定要删除吗?',
  confirmText: '删除',
  cancelText: '取消',
  confirmVariant: 'destructive', // 删除操作使用危险样式
  onConfirm: () => { /* ... */ },
  onCancel: () => { /* ... */ },
});
```

## 常见模式

### 模式 1: useRequest 成功回调
```typescript
// 旧
const { run } = useRequest(api, {
  manual: true,
  onSuccess: () => {
    alert.show({ title: '成功', description: '操作成功' });
  },
});

// 新
const { run } = useRequest(api, {
  manual: true,
  onSuccess: () => {
    toast.success('操作成功');
  },
});
```

### 模式 2: 参数验证
```typescript
// 旧
if (!data) {
  return alert.show({ title: '错误', description: '参数错误' });
}

// 新
if (!data) {
  return toast.error('参数错误');
}
```

### 模式 3: try-catch 错误处理
```typescript
// 旧
try {
  await doSomething();
  alert.show({ title: '成功', description: '操作成功' });
} catch (error) {
  alert.show({ title: '错误', description: error.message });
}

// 新
try {
  await doSomething();
  toast.success('操作成功');
} catch (error) {
  toast.error('操作失败', error.message);
}
```

## 自动化脚本 (可选)

如果文件很多,可以使用以下 bash 脚本批量替换:

```bash
#!/bin/bash

# 替换导入语句
find . -name "*.tsx" -o -name "*.ts" | xargs sed -i '' "s/import { alert } from '@\/utils\/alert';/import { toast } from '@\/utils\/toast';/g"

# 替换简单的 alert.show
# 注意:这只是基础替换,复杂情况需要手动处理
# 建议先备份,然后逐个文件检查

echo "已完成基础替换,请手动检查并调整每个文件"
```

## 注意事项

1. **不要直接运行自动化脚本** - 先备份代码
2. **逐个文件检查** - 确保替换正确
3. **测试每个功能** - 确保行为符合预期
4. **调整样式** - 根据操作类型选择合适的 toast 类型
5. **确认框样式** - 删除操作使用 `confirmVariant: 'destructive'`

