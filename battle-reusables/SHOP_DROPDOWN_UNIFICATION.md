# 店铺号下拉框统一实现

## 需求概述
在应用的多个页面中统一使用 `/shops/houses/options` 接口提供的下拉框来选择店铺号，替代原有的文本输入框。

## 实现范围

### 1. 桌台战绩查询 (tables-search.tsx) ✅
- **状态**: 已完成（之前已实现）
- **位置**: `battle-reusables/components/(tabs)/tables/tables-search.tsx`
- **实现**: 使用 `shopsHousesOptions` 获取店铺列表，通过 Select 组件展示下拉框

### 2. 资金页面 (funds-search.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(tabs)/funds/funds-search.tsx`
- **变更**:
  - 添加导入: `shopsHousesOptions`, `Select` 组件, `TriggerRef`, `isWeb`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框
  - 保留历史记录功能的兼容性

### 3. 首页统计 (stats-search.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(tabs)/index/admin/stats-search.tsx`
- **变更**:
  - 添加导入: `useRequest`, `shopsHousesOptions`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框
  - 保留时间选择器和圈子ID输入框

### 4. 成员圈子 (members-view.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(tabs)/members/members-view.tsx`
- **变更**:
  - 添加导入: `useRef`, `shopsHousesOptions`, `Select` 组件, `TriggerRef`, `isWeb`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框（仅在超级管理员模式下显示）
  - 店铺管理员自动加载，无需手动输入

### 5. 我的战绩 (my-battles-view.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(shop)/battles/my-battles-view.tsx`
- **变更**:
  - 添加导入: `shopsHousesOptions`, `Select` 组件, `TriggerRef`, `isWeb`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框
  - 移除时间选择输入框，改为自动查询最近7天的数据
  - 添加 `getTimeRange()` 函数自动计算最近7天的时间戳

### 6. 圈子战绩 (group-battles-view.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(shop)/battles/group-battles-view.tsx`
- **变更**:
  - 添加导入: `shopsHousesOptions`, `Select` 组件, `TriggerRef`, `isWeb`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框

### 7. 我的余额 (my-balances-view.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(shop)/battles/my-balances-view.tsx`
- **变更**:
  - 添加导入: `shopsHousesOptions`, `Select` 组件, `TriggerRef`, `isWeb`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框

### 8. 圈子成员余额 (group-balances-view.tsx) ✅
- **状态**: 已完成
- **位置**: `battle-reusables/components/(shop)/battles/group-balances-view.tsx`
- **变更**:
  - 添加导入: `shopsHousesOptions`, `Select` 组件, `TriggerRef`, `isWeb`
  - 添加 hooks: `useRequest(shopsHousesOptions)`, `useRef<TriggerRef>`
  - 将店铺号输入框替换为下拉框

## 技术细节

### 共同的实现模式
```typescript
// 1. 导入必要的组件和 hooks
import { useRequest } from '@/hooks/use-request';
import { shopsHousesOptions } from '@/services/shops/houses';
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from '@/components/ui/select';
import { TriggerRef } from '@rn-primitives/select';
import { isWeb } from '@/utils/platform';

// 2. 获取店铺选项
const { data: houseOptions } = useRequest(shopsHousesOptions);
const ref = useRef<TriggerRef>(null);

// 3. Web 平台的触摸处理
function onTouchStart() {
  isWeb && ref.current?.open();
}

// 4. 使用 Select 组件
<Select
  value={value ? ({ label: `店铺 ${value}`, value } as any) : undefined}
  onValueChange={(opt) => onChange(String(opt?.value ?? ''))}
>
  <SelectTrigger ref={ref} onTouchStart={onTouchStart} className="min-w-[160px]">
    <SelectValue placeholder={value ? `店铺 ${value}` : '选择店铺号'} />
  </SelectTrigger>
  <SelectContent>
    <SelectGroup>
      <SelectLabel>店铺号</SelectLabel>
      {(houseOptions ?? []).map((gid) => (
        <SelectItem key={String(gid)} label={`店铺 ${gid}`} value={String(gid)}>
          店铺 {gid}
        </SelectItem>
      ))}
    </SelectGroup>
  </SelectContent>
</Select>
```

## 后端接口
- **端点**: `/shops/houses/options`
- **方法**: GET
- **返回**: `number[]` - 店铺号列表
- **实现**: `battle-reusables/services/shops/houses/index.ts`

## 测试建议
1. 验证下拉框能正确加载店铺列表
2. 测试选择不同店铺号的功能
3. 验证 Web 和移动端的触摸交互
4. 确保表单验证仍然正常工作
5. 检查历史记录功能（资金页面）的兼容性

## 相关文件清单
- `battle-reusables/components/(tabs)/tables/tables-search.tsx`
- `battle-reusables/components/(tabs)/funds/funds-search.tsx`
- `battle-reusables/components/(tabs)/index/admin/stats-search.tsx`
- `battle-reusables/components/(tabs)/members/members-view.tsx`
- `battle-reusables/components/(shop)/battles/my-battles-view.tsx`
- `battle-reusables/components/(shop)/battles/group-battles-view.tsx`
- `battle-reusables/components/(shop)/battles/my-balances-view.tsx`
- `battle-reusables/components/(shop)/battles/group-balances-view.tsx`
- `battle-reusables/services/shops/houses/index.ts`

