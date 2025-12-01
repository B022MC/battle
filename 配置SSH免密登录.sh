#!/bin/bash

# 配置 SSH 免密登录，避免每次都输入密码

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"

echo -e "${YELLOW}配置 SSH 免密登录...${NC}"
echo ""

# 检查是否已有密钥
if [ ! -f ~/.ssh/id_rsa ]; then
    echo "生成 SSH 密钥..."
    ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""
    echo -e "${GREEN}✓ 密钥已生成${NC}"
else
    echo -e "${GREEN}✓ SSH 密钥已存在${NC}"
fi

echo ""
echo "复制公钥到服务器..."
echo "（需要输入一次服务器密码）"
echo ""

ssh-copy-id ${SERVER_USER}@${SERVER_IP}

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "配置成功！"
    echo "======================================${NC}"
    echo ""
    echo "测试免密登录："
    echo "  ssh ${SERVER_USER}@${SERVER_IP}"
    echo ""
    echo "现在部署脚本不需要输入密码了！"
else
    echo ""
    echo -e "${YELLOW}配置失败，请检查网络或服务器设置${NC}"
fi
