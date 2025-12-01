#!/bin/bash

# 完整构建：前端 + 后端一起构建并上传到服务器

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"

echo -e "${BLUE}======================================"
echo "完整构建与上传"
echo "======================================${NC}"

# 1. 构建后端 Docker 镜像
echo ""
echo -e "${YELLOW}[1/4] 构建后端 Docker 镜像...${NC}"

cd /Users/b022mc/project/battle/battle-tiles
docker build -t battle-tiles:latest .

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 后端镜像构建失败${NC}"
    exit 1
fi

echo "保存镜像..."
docker save -o /tmp/battle-tiles.tar battle-tiles:latest

SIZE=$(du -h /tmp/battle-tiles.tar | cut -f1)
echo -e "${GREEN}✓ 后端镜像已保存 (${SIZE})${NC}"

# 2. 构建前端
echo ""
echo -e "${YELLOW}[2/4] 构建前端...${NC}"

cd /Users/b022mc/project/battle/battle-reusables

if [ ! -f ".env.production" ]; then
    echo -e "${RED}✗ 缺少 .env.production 文件${NC}"
    exit 1
fi

npm install
npm run build:web

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 前端构建失败${NC}"
    exit 1
fi

echo "打包前端..."
cd web-build
tar -czf /tmp/battle-frontend.tar.gz .
cd ..

echo -e "${GREEN}✓ 前端构建完成${NC}"

# 3. 打包配置和部署脚本
echo ""
echo -e "${YELLOW}[3/4] 打包配置文件...${NC}"

cd /Users/b022mc/project/battle

tar -czf /tmp/battle-configs.tar.gz \
  battle-tiles/configs/ \
  服务器启动.sh

echo -e "${GREEN}✓ 配置文件已打包${NC}"

# 4. 上传到服务器
echo ""
echo -e "${YELLOW}[4/4] 上传到服务器...${NC}"

echo "上传后端镜像..."
scp /tmp/battle-tiles.tar ${SERVER_USER}@${SERVER_IP}:/tmp/

echo "上传前端..."
scp /tmp/battle-frontend.tar.gz ${SERVER_USER}@${SERVER_IP}:/tmp/

echo "上传配置..."
scp /tmp/battle-configs.tar.gz ${SERVER_USER}@${SERVER_IP}:/tmp/

echo -e "${GREEN}✓ 上传完成${NC}"

# 清理本地临时文件
echo ""
echo "清理临时文件..."
rm -f /tmp/battle-tiles.tar
rm -f /tmp/battle-frontend.tar.gz
rm -f /tmp/battle-configs.tar.gz

echo ""
echo -e "${GREEN}======================================"
echo "构建上传完成！"
echo "======================================${NC}"
echo ""
echo "下一步："
echo "1. SSH 登录服务器: ssh ${SERVER_USER}@${SERVER_IP}"
echo "2. 解压配置: cd /root && tar -xzf /tmp/battle-configs.tar.gz"
echo "3. 运行部署脚本: ./服务器启动.sh"
echo ""
echo "或运行快捷命令:"
echo "  ssh ${SERVER_USER}@${SERVER_IP} 'cd /root && tar -xzf /tmp/battle-configs.tar.gz && ./服务器启动.sh'"
