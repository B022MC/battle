#!/bin/bash

# Toast 迁移脚本
# 此脚本会批量替换所有 alert 为 toast

# 颜色定义
RED='\033[0:31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始迁移 alert 到 toast...${NC}\n"

# 需要迁移的文件列表
files=(
  "components/(tabs)/profile/profile-ctrl-accounts.tsx"
  "components/(tabs)/profile/profile-ctrl-account.tsx"
  "components/(tabs)/members/members-view.tsx"
  "components/(shop)/battles/my-battles-view.tsx"
  "components/(shop)/battles/my-balances-view.tsx"
  "components/(shop)/battles/group-battles-view.tsx"
  "components/(shop)/battles/group-balances-view.tsx"
  "components/(shop)/admins/admins-assign.tsx"
)

# 备份目录
backup_dir="./backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$backup_dir"

echo -e "${YELLOW}备份原文件到: $backup_dir${NC}\n"

# 迁移每个文件
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    echo -e "处理: ${GREEN}$file${NC}"
    
    # 备份原文件
    cp "$file" "$backup_dir/"
    
    # 1. 替换导入语句
    sed -i '' "s/import { alert } from '@\/utils\/alert';/import { toast } from '@\/utils\/toast';/g" "$file"
    
    # 2. 替换简单的 alert.show (这只是基础替换, 复杂情况需要手动处理)
    # 注意: 这个脚本只做基础替换, 不处理复杂情况
    
    echo -e "  ✓ 已替换导入语句"
    echo -e "  ${YELLOW}⚠️  请手动检查并调整 alert.show 的调用${NC}"
    echo ""
  else
    echo -e "${RED}✗ 文件不存在: $file${NC}\n"
  fi
done

echo -e "${GREEN}迁移完成!${NC}"
echo -e "${YELLOW}请注意:${NC}"
echo -e "1. 已备份原文件到: $backup_dir"
echo -e "2. 只替换了导入语句"
echo -e "3. 请手动检查每个文件中的 alert.show 调用"
echo -e "4. 根据 TOAST_USAGE.md 进行调整"
echo -e "5. 测试每个功能确保正常工作\n"

echo -e "${YELLOW}常见替换模式:${NC}"
echo -e "  alert.show({ title: '成功', description: 'xxx' })"
echo -e "  → toast.success('xxx')"
echo -e ""
echo -e "  alert.show({ title: '错误', description: 'xxx' })"
echo -e "  → toast.error('xxx')"
echo -e ""
echo -e "  alert.show({ ..., onConfirm: ... })"
echo -e "  → toast.confirm({ ..., confirmVariant: 'destructive' })"
echo -e ""

# 显示还有多少 alert.show 需要处理
remaining=$(grep -r "alert\.show" components --include="*.tsx" | wc -l | tr -d ' ')
if [ "$remaining" -gt 0 ]; then
  echo -e "${YELLOW}还有 $remaining 处 alert.show 需要手动替换${NC}"
else
  echo -e "${GREEN}所有 alert 导入已替换完成!${NC}"
fi

