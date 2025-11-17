# Toast 系统测试指南

## 问题修复

### 修复的问题
✅ **无限循环错误** - "Maximum update depth exceeded"

### 修复内容
1. **移除 Button 组件依赖** - 改用 `Pressable` 避免循环依赖
2. **修复定时器管理** - 使用 `useRef` 正确管理定时器，避免状态更新循环
3. **添加清理逻辑** - 确保组件卸载时清理定时器

### 修改的文件
- `components/shared/toast-provider.tsx`

## 测试步骤

### 1. 测试成功提示
在任意组件中调用：
```typescript
import { toast } from '@/utils/toast';

toast.success('测试成功');
toast.success('操作成功', '这是详细描述');
```

**预期结果**：
- 显示绿色提示框
- 有成功图标
- 3秒后自动消失

### 2. 测试错误提示
```typescript
toast.error('测试错误');
toast.error('操作失败', '这是错误详情');
```

**预期结果**：
- 显示红色提示框
- 有错误图标
- 4秒后自动消失

### 3. 测试警告提示
```typescript
toast.warning('测试警告');
toast.warning('警告', '这是警告信息');
```

**预期结果**：
- 显示黄色提示框
- 有警告图标
- 3.5秒后自动消失

### 4. 测试信息提示
```typescript
toast.info('测试信息');
toast.info('提示', '这是信息内容');
```

**预期结果**：
- 显示蓝色提示框
- 有信息图标
- 3秒后自动消失

### 5. 测试确认对话框
```typescript
toast.confirm({
  title: '确认删除',
  description: '此操作无法撤销，确定要继续吗？',
  confirmText: '删除',
  cancelText: '取消',
  confirmVariant: 'destructive',
  onConfirm: () => {
    console.log('确认删除');
    toast.success('删除成功');
  },
  onCancel: () => {
    console.log('取消删除');
  },
});
```

**预期结果**：
- 显示模态对话框
- 背景半透明遮罩
- 确认按钮是红色(destructive)
- 取消按钮是灰色边框
- 点击确认/取消后对话框关闭

### 6. 测试异步确认
```typescript
toast.confirm({
  title: '确认提交',
  description: '确定要提交这些更改吗？',
  confirmText: '提交',
  cancelText: '取消',
  onConfirm: async () => {
    // 模拟异步操作
    await new Promise(resolve => setTimeout(resolve, 2000));
    toast.success('提交成功');
  },
});
```

**预期结果**：
- 点击确认后按钮显示"处理中..."
- 按钮禁用状态
- 2秒后对话框关闭
- 显示成功提示

### 7. 测试不同位置
```typescript
// 顶部
toast.show({
  title: '顶部提示',
  position: 'top',
  type: 'info',
});

// 居中
toast.show({
  title: '居中提示',
  position: 'center',
  type: 'info',
});

// 底部
toast.show({
  title: '底部提示',
  position: 'bottom',
  type: 'info',
});
```

**预期结果**：
- 提示显示在相应位置
- 动画效果正确

### 8. 测试手动关闭
```typescript
toast.show({
  title: '永不自动关闭',
  description: '需要手动关闭',
  duration: 0,
  showClose: true,
  type: 'warning',
});
```

**预期结果**：
- 提示不会自动消失
- 显示关闭按钮
- 点击关闭按钮后提示消失

### 9. 测试快速连续调用
```typescript
toast.success('第一个提示');
setTimeout(() => toast.error('第二个提示'), 100);
setTimeout(() => toast.warning('第三个提示'), 200);
```

**预期结果**：
- 只显示最后一个提示(第三个)
- 之前的提示被替换
- 没有错误或卡顿

## 常见问题排查

### 问题 1: 提示不显示
**检查**：
- [ ] `ToastProvider` 是否已添加到 `app/_layout.tsx`
- [ ] 是否在 `PortalHost` 之前添加
- [ ] 导入路径是否正确

### 问题 2: 样式异常
**检查**：
- [ ] Tailwind CSS 配置是否正确
- [ ] `global.css` 是否已导入
- [ ] 颜色主题是否定义

### 问题 3: 动画不流畅
**检查**：
- [ ] `react-native-reanimated` 是否正确安装
- [ ] 是否需要重启开发服务器

### 问题 4: 确认对话框点击无响应
**检查**：
- [ ] `onConfirm` 函数是否正确定义
- [ ] 是否有 JavaScript 错误
- [ ] 检查浏览器控制台

## 性能测试

### 测试 1: 内存泄漏
```typescript
// 连续显示100个提示
for (let i = 0; i < 100; i++) {
  setTimeout(() => {
    toast.info(`提示 ${i + 1}`);
  }, i * 100);
}
```

**预期**：
- 不应该有内存泄漏
- 性能保持稳定
- 最后一个提示正常显示

### 测试 2: 长文本处理
```typescript
toast.info(
  '这是一个非常长的标题'.repeat(10),
  '这是一个非常长的描述'.repeat(20)
);
```

**预期**：
- 文本正确换行
- 布局不崩溃
- 可以正常关闭

## 集成测试

### 测试现有功能
1. **登录页面** - 测试错误提示
2. **个人资料页** - 测试成功提示
3. **删除操作** - 测试确认对话框
4. **表单提交** - 测试异步确认

## 测试清单

- [ ] 成功提示正常显示
- [ ] 错误提示正常显示
- [ ] 警告提示正常显示
- [ ] 信息提示正常显示
- [ ] 确认对话框正常显示
- [ ] 异步确认正常工作
- [ ] 不同位置显示正确
- [ ] 手动关闭功能正常
- [ ] 快速连续调用无问题
- [ ] 深色模式下样式正确
- [ ] Web 和 Native 平台都正常
- [ ] 无内存泄漏
- [ ] 无性能问题
- [ ] 与现有功能集成无冲突

## 回归测试

确保以下页面的提示功能正常：
- [ ] 登录/注册页面
- [ ] 个人资料页面
- [ ] 游戏账号绑定页面
- [ ] 表单提交页面
- [ ] 删除确认操作

## 修复确认

✅ 无限循环错误已修复
✅ 定时器管理正确
✅ 组件依赖优化
✅ 性能优化完成

