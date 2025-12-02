#!/bin/bash

# 后端直接编译并上传（不使用 Docker）

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"

echo -e "${BLUE}======================================"
echo "后端直接编译与上传"
echo "======================================${NC}"

cd /Users/b022mc/project/battle/battle-tiles

# 1. 本地编译（需要本地有 Go 环境）
echo ""
echo -e "${YELLOW}步骤 1/4: 检查 Go 环境...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ 本地未安装 Go${NC}"
    echo ""
    echo "请选择："
    echo "1. 安装 Go: https://go.dev/dl/"
    echo "2. 或使用服务器端编译（稍后会在服务器上编译）"
    echo ""
    read -p "是否跳过本地编译，在服务器编译? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        SERVER_BUILD=true
    else
        exit 1
    fi
else
    echo -e "${GREEN}✓ Go 已安装${NC}"
    go version
    SERVER_BUILD=false
fi

# 2. 本地编译或打包源码
if [ "$SERVER_BUILD" = false ]; then
    echo ""
    echo -e "${YELLOW}步骤 2/4: 本地编译（交叉编译到 Linux）...${NC}"
    
    # 清理旧的编译产物
    rm -rf bin/
    
    # 交叉编译到 Linux
    echo "编译目标: Linux amd64"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo 'unknown')" -o ./bin/ ./...
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ 编译失败${NC}"
        exit 1
    fi
    
    # 检查编译产物
    if [ ! -d "bin" ] || [ -z "$(ls -A bin)" ]; then
        echo -e "${RED}✗ 编译产物不存在${NC}"
        ls -la bin/ 2>/dev/null || echo "bin/ 目录为空"
        exit 1
    fi
    
    # 显示编译产物
    echo "编译产物："
    ls -lh bin/
    
    # 验证是 Linux 可执行文件
    echo ""
    echo "验证平台："
    file bin/*
    
    echo -e "${GREEN}✓ 编译完成 (Linux amd64)${NC}"
    
    # 打包编译产物
    echo "打包..."
    tar -czf /tmp/battle-backend.tar.gz bin/ configs/
else
    echo ""
    echo -e "${YELLOW}步骤 2/4: 打包源码...${NC}"
    
    # 打包源码（在服务器编译）
    tar -czf /tmp/battle-backend.tar.gz \
        --exclude='bin/' \
        --exclude='.git/' \
        --exclude='*.log' \
        .
fi

SIZE=$(du -h /tmp/battle-backend.tar.gz | cut -f1)
echo -e "${GREEN}✓ 打包完成 (${SIZE})${NC}"

# 3. 上传到服务器
echo ""
echo -e "${YELLOW}步骤 3/4: 上传到服务器...${NC}"

scp /tmp/battle-backend.tar.gz ${SERVER_USER}@${SERVER_IP}:/tmp/

echo -e "${GREEN}✓ 上传完成${NC}"

# 4. 在服务器上部署
echo ""
echo -e "${YELLOW}步骤 4/4: 服务器端部署...${NC}"

if [ "$SERVER_BUILD" = true ]; then
    # 服务器端编译
    ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
set -e

echo "解压源码..."
mkdir -p /root/battle-tiles
cd /root/battle-tiles
tar -xzf /tmp/battle-backend.tar.gz

echo "检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "安装 Go..."
    wget -q https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
fi

echo "编译..."
export GOPROXY="https://goproxy.cn,direct"
make build

echo "停止旧进程..."
pkill -f "go-kgin-platform" || true

echo "启动服务..."
mkdir -p logs

# 启动主服务
nohup ./bin/go-kgin-platform > logs/platform.log 2>&1 &
PLATFORM_PID=$!

sleep 3

# 检查服务状态
if ps -p $PLATFORM_PID > /dev/null; then
    echo "✓ Platform 服务启动成功 (PID: $PLATFORM_PID)"
else
    echo "✗ Platform 服务启动失败"
    tail -20 logs/platform.log
    exit 1
fi

rm -f /tmp/battle-backend.tar.gz
ENDSSH
else
    # 直接运行编译好的程序
    ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
set -e

echo "解压..."
mkdir -p /root/battle-tiles
cd /root/battle-tiles
tar -xzf /tmp/battle-backend.tar.gz

echo "停止旧进程..."
pkill -f "go-kgin-platform" || true

echo "启动服务..."
mkdir -p logs

# 启动主服务
nohup ./bin/go-kgin-platform > logs/platform.log 2>&1 &
PLATFORM_PID=$!

sleep 3

# 检查服务状态
if ps -p $PLATFORM_PID > /dev/null; then
    echo "✓ Platform 服务启动成功 (PID: $PLATFORM_PID)"
else
    echo "✗ Platform 服务启动失败"
    tail -20 logs/platform.log
    exit 1
fi

rm -f /tmp/battle-backend.tar.gz
ENDSSH
fi

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "部署成功！"
    echo "======================================${NC}"
    echo ""
    echo "后端地址: http://${SERVER_IP}:8000"
    echo ""
    echo "查看日志: ssh ${SERVER_USER}@${SERVER_IP} 'tail -f /root/battle-tiles/logs/app.log'"
    echo ""
    
    # 清理本地临时文件
    rm -f /tmp/battle-backend.tar.gz
    
    echo -e "${YELLOW}重要: 确保云服务器安全组已开放 8000 端口！${NC}"
else
    echo -e "${RED}部署失败，请检查错误信息${NC}"
    exit 1
fi
