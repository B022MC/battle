# 🎈 气泡提示组件使用指南

## 📋 概述

自定义气泡提示组件用于显示操作成功、失败、警告和信息提示。相比普通 toast，气泡提示具有更美观的 UI 和更好的用户体验。

---

## ✨ 特性

- 🎨 **美观的 UI** - 圆角气泡设计，带阴影和边框
- 🌈 **多种类型** - 成功（绿色）、错误（红色）、警告（黄色）、信息（蓝色）
- 🎭 **流畅动画** - 入场和出场动画，支持弹簧效果
- 📱 **响应式** - 自适应不同屏幕尺寸
- 🌓 **深色模式** - 自动适配深色主题
- ⏱️ **自动关闭** - 可配置显示时长
- ✖️ **手动关闭** - 点击关闭按钮立即关闭
- 📚 **多个提示** - 支持同时显示多个气泡

---

## 🎨 视觉效果

### 成功提示（绿色）
```
┌─────────────────────────────────┐
│ ✓  更新成功                      │
│    权限"shop:admin:view"已更新   │
└─────────────────────────────────┘
```

### 错误提示（红色）
```
┌─────────────────────────────────┐
│ ✗  删除失败                      │
│    权限正在被使用，无法删除       │
└─────────────────────────────────┘
```

### 警告提示（黄色）
```
┌─────────────────────────────────┐
│ ⚠  注意                         │
│    此操作将影响所有用户           │
└─────────────────────────────────┘
```

### 信息提示（蓝色）
```
┌─────────────────────────────────┐
│ ℹ  提示                         │
│    数据已保存到本地               │
└─────────────────────────────────┘
```

---

## 📦 安装

气泡提示组件已经集成到应用中，无需额外安装。

---

## 🚀 使用方法

### 1. 导入

```typescript
import { showSuccessBubble, showErrorBubble, showWarningBubble, showInfoBubble } from '@/utils/bubble-toast';
```

或者直接使用 bubbleToast 对象：

```typescript
import { bubbleToast } from '@/utils/bubble-toast';
```

### 2. 基本用法

#### 成功提示
```typescript
showSuccessBubble('操作成功');
// 或带描述
showSuccessBubble('更新成功', '权限"shop:admin:view"已更新');
```

#### 错误提示
```typescript
showErrorBubble('操作失败');
// 或带描述
showErrorBubble('删除失败', '权限正在被使用，无法删除');
```

#### 警告提示
```typescript
showWarningBubble('注意');
// 或带描述
showWarningBubble('注意', '此操作将影响所有用户');
```

#### 信息提示
```typescript
showInfoBubble('提示');
// 或带描述
showInfoBubble('提示', '数据已保存到本地');
```

### 3. 高级用法

使用 `bubbleToast.show()` 自定义更多选项：

```typescript
import { bubbleToast } from '@/utils/bubble-toast';

bubbleToast.show({
  type: 'success',
  title: '操作成功',
  description: '数据已保存',
  duration: 5000, // 显示 5 秒
});
```

---

## 📝 API 参考

### showSuccessBubble(title, description?)

显示成功提示（绿色）

**参数：**
- `title` (string) - 标题，必填
- `description` (string, 可选) - 描述文本

**示例：**
```typescript
showSuccessBubble('创建成功', '角色"管理员"已创建');
```

---

### showErrorBubble(title, description?)

显示错误提示（红色）

**参数：**
- `title` (string) - 标题，必填
- `description` (string, 可选) - 描述文本

**示例：**
```typescript
showErrorBubble('删除失败', '该项目正在被使用');
```

---

### showWarningBubble(title, description?)

显示警告提示（黄色）

**参数：**
- `title` (string) - 标题，必填
- `description` (string, 可选) - 描述文本

**示例：**
```typescript
showWarningBubble('注意', '此操作不可撤销');
```

---

### showInfoBubble(title, description?)

显示信息提示（蓝色）

**参数：**
- `title` (string) - 标题，必填
- `description` (string, 可选) - 描述文本

**示例：**
```typescript
showInfoBubble('提示', '数据已同步');
```

---

### bubbleToast.show(options)

显示自定义气泡提示

**参数：**
- `options` (object)
  - `type` ('success' | 'error' | 'warning' | 'info') - 提示类型，必填
  - `title` (string) - 标题，必填
  - `description` (string, 可选) - 描述文本
  - `duration` (number, 可选) - 显示时长（毫秒），默认值根据类型不同：
    - success: 3000ms
    - error: 4000ms
    - warning: 3500ms
    - info: 3000ms

**示例：**
```typescript
bubbleToast.show({
  type: 'success',
  title: '保存成功',
  description: '您的更改已保存',
  duration: 5000,
});
```

---

## 💡 使用场景

### 1. CRUD 操作成功

```typescript
// 创建
const handleCreate = async () => {
  const res = await createPermission(data);
  if (res.code === 0) {
    showSuccessBubble('创建成功', `权限"${data.name}"已创建`);
  }
};

// 更新
const handleUpdate = async () => {
  const res = await updateRole(data);
  if (res.code === 0) {
    showSuccessBubble('更新成功', `角色"${data.name}"已更新`);
  }
};

// 删除
const handleDelete = async () => {
  const res = await deleteMenu(id);
  if (res.code === 0) {
    showSuccessBubble('删除成功', `菜单"${menu.title}"已删除`);
  }
};
```

### 2. 上分/下分操作

```typescript
const handleDeposit = async () => {
  const res = await deposit(amount);
  if (res.code === 0) {
    showSuccessBubble('上分成功', `已为${memberName}上分 ${amount} 元`);
  }
};
```

### 3. 账号绑定/解绑

```typescript
const handleBind = async () => {
  const res = await bindAccount(account);
  if (res.code === 0) {
    showSuccessBubble('绑定成功', '游戏账号已成功绑定');
  }
};
```

---

## 🎯 最佳实践

### 1. 标题简洁明了
✅ 好的标题：
- "创建成功"
- "更新成功"
- "删除成功"

❌ 不好的标题：
- "操作已经成功完成"
- "您的请求已被处理"

### 2. 描述提供详细信息
✅ 好的描述：
- `权限"shop:admin:view"已创建`
- `已为张三上分 100 元`
- `菜单"用户管理"已删除`

❌ 不好的描述：
- "成功"
- "完成"
- "OK"

### 3. 只在成功时使用气泡提示
- ✅ **成功操作** - 使用气泡提示
- ❌ **错误提示** - 继续使用 `showToast('错误信息', 'error')`
- ❌ **表单验证** - 使用 `showToast('请填写完整信息', 'error')`

### 4. 合理设置显示时长
- 简单操作（创建、更新）：3000ms（默认）
- 重要操作（删除、转账）：4000-5000ms
- 信息提示：2000-3000ms

---

## 🔧 技术细节

### 组件结构

```
BubbleToastContainer (容器)
  └── BubbleToast (单个气泡)
        ├── 图标
        ├── 内容（标题 + 描述）
        └── 关闭按钮
```

### 动画效果

- **入场动画**：淡入 + 向下滑动 + 缩放
- **出场动画**：淡出 + 向上滑动 + 缩小
- **弹簧效果**：使用 `Animated.spring` 实现自然的弹性动画

### 样式定制

气泡提示支持深色模式，颜色会自动适配：

| 类型 | 浅色模式 | 深色模式 |
|------|---------|---------|
| Success | 绿色 | 深绿色 |
| Error | 红色 | 深红色 |
| Warning | 黄色 | 深黄色 |
| Info | 蓝色 | 深蓝色 |

---

## 📊 已集成的组件

以下组件已经使用气泡提示：

- ✅ 权限管理（创建、更新、删除）
- ✅ 角色管理（创建、更新、删除）
- ✅ 菜单管理（创建、更新、删除）
- ✅ 上分/下分操作
- ✅ 游戏账号绑定/解绑
- ✅ 中控账号管理（创建、更新、删除）

---

## 🐛 故障排除

### 气泡不显示

**原因：** `BubbleToastContainer` 未添加到应用根组件

**解决：** 检查 `app/_layout.tsx` 是否包含：
```typescript
import { BubbleToastContainer } from '@/components/ui/bubble-toast-container';

// 在 return 中
<BubbleToastContainer />
```

### 气泡位置不对

**原因：** z-index 层级问题

**解决：** `BubbleToastContainer` 的 z-index 设置为 100，确保在最上层

### 动画卡顿

**原因：** 使用了非原生动画驱动

**解决：** 所有动画都使用 `useNativeDriver: true`

---

## 🎉 示例代码

完整的使用示例：

```typescript
import React from 'react';
import { View, Button } from 'react-native';
import { showSuccessBubble, showErrorBubble } from '@/utils/bubble-toast';

export function ExampleComponent() {
  const handleSave = async () => {
    try {
      const res = await saveData();
      if (res.code === 0) {
        showSuccessBubble('保存成功', '您的数据已保存');
      } else {
        showErrorBubble('保存失败', res.msg);
      }
    } catch (error) {
      showErrorBubble('保存失败', '网络错误');
    }
  };

  return (
    <View>
      <Button title="保存" onPress={handleSave} />
    </View>
  );
}
```

---

## 📞 支持

如有问题或建议，请联系开发团队。

---

**更新时间：** 2025-11-20
**版本：** 1.0.0

