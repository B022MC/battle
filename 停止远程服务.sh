#!/bin/bash

# 停止远程服务

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"

echo -e "${BLUE}======================================"
echo "停止远程服务"
echo "======================================${NC}"

echo ""
echo -e "${YELLOW}连接到服务器并停止服务...${NC}"

ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
set -e

echo "正在停止后端服务..."

# 停止 go-kgin-platform
if pgrep -f "go-kgin-platform" > /dev/null; then
    pkill -f "go-kgin-platform"
    echo "✓ go-kgin-platform 已停止"
else
    echo "○ go-kgin-platform 未运行"
fi

# 停止 go-kgin-asynq（如果存在）
if pgrep -f "go-kgin-asynq" > /dev/null; then
    pkill -f "go-kgin-asynq"
    echo "✓ go-kgin-asynq 已停止"
else
    echo "○ go-kgin-asynq 未运行"
fi

echo ""
echo "正在停止 Nginx..."
if systemctl is-active --quiet nginx; then
    systemctl stop nginx
    echo "✓ Nginx 已停止"
else
    echo "○ Nginx 未运行"
fi

echo ""
echo "检查进程状态..."
if pgrep -f "go-kgin" > /dev/null; then
    echo "⚠ 警告: 仍有 go-kgin 进程在运行"
    pgrep -af "go-kgin"
else
    echo "✓ 所有后端进程已停止"
fi

if systemctl is-active --quiet nginx; then
    echo "⚠ 警告: Nginx 仍在运行"
else
    echo "✓ Nginx 已停止"
fi

ENDSSH

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "服务停止成功！"
    echo "======================================${NC}"
    echo ""
    echo "已停止的服务："
    echo "  • go-kgin-platform (后端主服务)"
    echo "  • go-kgin-asynq (异步任务服务)"
    echo "  • Nginx (前端服务器)"
    echo ""
    echo -e "${YELLOW}提示: 使用相应的部署脚本可以重新启动服务${NC}"
else
    echo -e "${RED}停止服务失败，请检查错误信息${NC}"
    exit 1
fi
