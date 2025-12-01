#!/bin/bash

# 后端 Docker 镜像构建并上传到服务器

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"
IMAGE_NAME="battle-tiles"
IMAGE_TAG="latest"

echo -e "${BLUE}======================================"
echo "后端 Docker 镜像构建与上传"
echo "======================================${NC}"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    exit 1
fi

cd /Users/b022mc/project/battle/battle-tiles

# 1. 构建镜像
echo ""
echo -e "${YELLOW}步骤 1/5: 构建 Docker 镜像...${NC}"
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 镜像构建失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 镜像构建完成${NC}"

# 2. 保存镜像为 tar 文件
echo ""
echo -e "${YELLOW}步骤 2/5: 保存镜像为 tar 文件...${NC}"

IMAGE_FILE="/tmp/battle-tiles.tar"
docker save -o ${IMAGE_FILE} ${IMAGE_NAME}:${IMAGE_TAG}

SIZE=$(du -h ${IMAGE_FILE} | cut -f1)
echo -e "${GREEN}✓ 镜像已保存: ${IMAGE_FILE} (${SIZE})${NC}"

# 3. 压缩配置文件
echo ""
echo -e "${YELLOW}步骤 3/5: 打包配置文件...${NC}"

CONFIG_FILE="/tmp/battle-configs.tar.gz"
tar -czf ${CONFIG_FILE} configs/

echo -e "${GREEN}✓ 配置文件已打包${NC}"

# 4. 上传到服务器
echo ""
echo -e "${YELLOW}步骤 4/5: 上传到服务器...${NC}"

echo "上传镜像文件..."
scp ${IMAGE_FILE} ${SERVER_USER}@${SERVER_IP}:/tmp/

echo "上传配置文件..."
scp ${CONFIG_FILE} ${SERVER_USER}@${SERVER_IP}:/tmp/

echo -e "${GREEN}✓ 上传完成${NC}"

# 5. 在服务器上加载镜像并启动
echo ""
echo -e "${YELLOW}步骤 5/5: 在服务器上部署...${NC}"

ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
set -e

echo "加载 Docker 镜像..."
docker load -i /tmp/battle-tiles.tar

echo "解压配置文件..."
mkdir -p /root/battle-tiles
cd /root/battle-tiles
tar -xzf /tmp/battle-configs.tar.gz

echo "创建日志目录..."
mkdir -p logs

echo "停止旧容器..."
docker stop battle-tiles 2>/dev/null || true
docker rm battle-tiles 2>/dev/null || true

echo "启动新容器..."
docker run -d \
  --name battle-tiles \
  --restart always \
  -p 8000:8000 \
  -v /root/battle-tiles/configs:/work/configs \
  -v /root/battle-tiles/logs:/work/logs \
  -e TZ=Asia/Shanghai \
  battle-tiles:latest

sleep 3

# 检查状态
if docker ps | grep -q battle-tiles; then
    echo "✓ 容器启动成功！"
    docker ps | grep battle-tiles
else
    echo "✗ 容器启动失败"
    docker logs battle-tiles
    exit 1
fi

# 清理临时文件
rm -f /tmp/battle-tiles.tar
rm -f /tmp/battle-configs.tar.gz

echo ""
echo "后端部署完成！"
ENDSSH

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "部署成功！"
    echo "======================================${NC}"
    echo ""
    echo "后端地址: http://${SERVER_IP}:8000"
    echo ""
    echo "查看日志: ssh ${SERVER_USER}@${SERVER_IP} 'docker logs -f battle-tiles'"
    
    # 清理本地临时文件
    rm -f ${IMAGE_FILE}
    rm -f ${CONFIG_FILE}
    
    echo -e "${YELLOW}重要: 确保云服务器安全组已开放 8000 端口！${NC}"
else
    echo -e "${RED}部署失败，请检查错误信息${NC}"
    exit 1
fi
