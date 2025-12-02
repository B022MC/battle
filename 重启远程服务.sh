#!/bin/bash

# 重启远程服务

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"

echo -e "${BLUE}======================================"
echo "重启远程服务"
echo "======================================${NC}"

echo ""
echo -e "${YELLOW}连接到服务器并重启服务...${NC}"

ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
set -e

echo "========================================="
echo "步骤 1/2: 停止现有服务"
echo "========================================="

# 停止后端服务
echo "停止后端服务..."
pkill -f "go-kgin-platform" || echo "go-kgin-platform 未运行"
pkill -f "go-kgin-asynq" || echo "go-kgin-asynq 未运行"

# 停止 Nginx
echo "停止 Nginx..."
systemctl stop nginx || echo "Nginx 未运行"

echo "✓ 所有服务已停止"

sleep 2

echo ""
echo "========================================="
echo "步骤 2/2: 启动服务"
echo "========================================="

# 进入后端目录
cd /root/battle-tiles

# 启动后端服务
echo "启动后端服务..."
mkdir -p logs

nohup ./bin/go-kgin-platform > logs/platform.log 2>&1 &
PLATFORM_PID=$!

sleep 3

# 检查后端服务
if ps -p $PLATFORM_PID > /dev/null; then
    echo "✓ Platform 服务启动成功 (PID: $PLATFORM_PID)"
else
    echo "✗ Platform 服务启动失败"
    tail -20 logs/platform.log
    exit 1
fi

# 启动 Nginx
echo ""
echo "启动 Nginx..."
systemctl start nginx

if systemctl is-active --quiet nginx; then
    echo "✓ Nginx 启动成功"
else
    echo "✗ Nginx 启动失败"
    systemctl status nginx
    exit 1
fi

echo ""
echo "========================================="
echo "服务状态检查"
echo "========================================="

# 检查后端进程
echo "后端进程："
pgrep -af "go-kgin" || echo "无后端进程"

echo ""
echo "Nginx 状态："
systemctl status nginx | grep Active

echo ""
echo "监听端口："
netstat -tlnp | grep -E ':(80|8000) ' || echo "无相关端口监听"

ENDSSH

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "服务重启成功！"
    echo "======================================${NC}"
    echo ""
    echo "访问地址："
    echo "  前端: http://${SERVER_IP}"
    echo "  后端: http://${SERVER_IP}:8000"
    echo ""
    echo "查看日志："
    echo "  后端: ssh ${SERVER_USER}@${SERVER_IP} 'tail -f /root/battle-tiles/logs/platform.log'"
    echo "  Nginx: ssh ${SERVER_USER}@${SERVER_IP} 'tail -f /var/log/nginx/error.log'"
else
    echo -e "${RED}服务重启失败，请检查错误信息${NC}"
    exit 1
fi
