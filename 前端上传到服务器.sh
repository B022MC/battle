#!/bin/bash

# 直接上传前端到服务器 Nginx

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVER_IP="8.137.52.203"
SERVER_USER="root"

echo -e "${BLUE}======================================"
echo "前端部署到服务器 Nginx"
echo "======================================${NC}"

cd /Users/b022mc/project/battle/battle-reusables

# 1. 强制重新构建（确保部署最新代码）
echo -e "${YELLOW}开始构建最新代码...${NC}"
rm -rf web-build
npm run build:web

echo -e "${GREEN}✓ 构建完成${NC}"

# 2. 打包
echo ""
echo -e "${YELLOW}打包前端文件...${NC}"
cd web-build
tar -czf /tmp/battle-frontend.tar.gz .
cd ..

SIZE=$(du -h /tmp/battle-frontend.tar.gz | cut -f1)
echo -e "${GREEN}✓ 打包完成 (${SIZE})${NC}"

# 3. 上传
echo ""
echo -e "${YELLOW}上传到服务器...${NC}"
scp /tmp/battle-frontend.tar.gz ${SERVER_USER}@${SERVER_IP}:/tmp/

echo -e "${GREEN}✓ 上传完成${NC}"

# 4. 部署
echo ""
echo -e "${YELLOW}在服务器上部署...${NC}"

ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
set -e

# 安装 Nginx
if ! command -v nginx &> /dev/null; then
    echo "安装 Nginx..."
    dnf install nginx -y 2>/dev/null || yum install nginx -y 2>/dev/null || apt install nginx -y
fi

# 解压前端文件
echo "解压前端文件..."
mkdir -p /var/www/battle
cd /var/www/battle
tar -xzf /tmp/battle-frontend.tar.gz

# 配置 Nginx
echo "配置 Nginx..."
cat > /etc/nginx/conf.d/battle.conf << 'EOF'
server {
    listen 80;
    server_name _;
    
    root /var/www/battle;
    index index.html;
    
    # SPA 路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_types text/plain text/css text/javascript application/javascript application/json;
    gzip_min_length 1024;
}
EOF

# 测试配置
echo "测试 Nginx 配置..."
nginx -t

# 重启 Nginx
echo "重启 Nginx..."
systemctl enable nginx
systemctl restart nginx

# 开放 80 端口
echo "配置防火墙..."
firewall-cmd --permanent --add-port=80/tcp 2>/dev/null || true
firewall-cmd --reload 2>/dev/null || true
ufw allow 80/tcp 2>/dev/null || true

# 清理
rm -f /tmp/battle-frontend.tar.gz

echo ""
echo "部署完成！"
SERVER_IP=$(hostname -I | awk '{print $1}')
echo "访问地址: http://${SERVER_IP}"
ENDSSH

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "前端部署成功！"
    echo "======================================${NC}"
    echo ""
    echo "访问地址："
    echo "  http://${SERVER_IP}"
    echo ""
    echo "后端 API: http://${SERVER_IP}:8000"
    echo ""
    echo -e "${YELLOW}重要: 在云服务器控制台开放 80 端口！${NC}"
    
    # 清理本地临时文件
    rm -f /tmp/battle-frontend.tar.gz
else
    echo -e "${RED}部署失败${NC}"
    exit 1
fi
