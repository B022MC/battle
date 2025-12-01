#!/bin/bash

# 服务器端启动脚本
# 加载 Docker 镜像并启动前后端服务

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================"
echo "服务器端部署启动"
echo "======================================${NC}"

# 1. 检查临时文件
echo -e "${YELLOW}检查上传的文件...${NC}"

if [ ! -f "/tmp/battle-tiles.tar" ]; then
    echo -e "${RED}✗ 后端镜像文件不存在${NC}"
    exit 1
fi

if [ ! -f "/tmp/battle-frontend.tar.gz" ]; then
    echo -e "${RED}✗ 前端文件不存在${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 文件检查通过${NC}"

# 2. 加载后端 Docker 镜像
echo ""
echo -e "${YELLOW}[1/4] 加载后端 Docker 镜像...${NC}"

docker load -i /tmp/battle-tiles.tar

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 镜像加载失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 镜像加载完成${NC}"

# 3. 部署后端
echo ""
echo -e "${YELLOW}[2/4] 部署后端服务...${NC}"

mkdir -p /root/battle-tiles/logs

cd /root/battle-tiles

# 停止旧容器
docker stop battle-tiles 2>/dev/null || true
docker rm battle-tiles 2>/dev/null || true

# 启动新容器
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
    echo -e "${GREEN}✓ 后端启动成功${NC}"
    docker ps | grep battle-tiles
else
    echo -e "${RED}✗ 后端启动失败${NC}"
    docker logs battle-tiles
    exit 1
fi

# 4. 部署前端
echo ""
echo -e "${YELLOW}[3/4] 部署前端...${NC}"

# 安装 Nginx
if ! command -v nginx &> /dev/null; then
    echo "安装 Nginx..."
    apt update && apt install nginx -y 2>/dev/null || yum install nginx -y
fi

# 解压前端
mkdir -p /var/www/battle
cd /var/www/battle
tar -xzf /tmp/battle-frontend.tar.gz

# 配置 Nginx
cat > /etc/nginx/conf.d/battle.conf << 'EOF'
server {
    listen 80;
    server_name _;
    
    root /var/www/battle;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API 代理（可选）
    location /api/ {
        proxy_pass http://localhost:8000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css text/javascript application/javascript application/json;
}
EOF

# 测试并重启 Nginx
nginx -t && systemctl restart nginx

echo -e "${GREEN}✓ 前端部署完成${NC}"

# 5. 配置防火墙
echo ""
echo -e "${YELLOW}[4/4] 配置防火墙...${NC}"

# 开放 8000 端口
if command -v firewall-cmd &> /dev/null; then
    firewall-cmd --permanent --add-port=8000/tcp 2>/dev/null || true
    firewall-cmd --permanent --add-port=80/tcp 2>/dev/null || true
    firewall-cmd --reload 2>/dev/null || true
    echo -e "${GREEN}✓ 防火墙配置完成 (firewalld)${NC}"
elif command -v ufw &> /dev/null; then
    ufw allow 8000/tcp 2>/dev/null || true
    ufw allow 80/tcp 2>/dev/null || true
    echo -e "${GREEN}✓ 防火墙配置完成 (ufw)${NC}"
fi

echo -e "${YELLOW}⚠️  请在云服务器控制台开放 8000 和 80 端口${NC}"

# 6. 清理临时文件
echo ""
echo "清理临时文件..."
rm -f /tmp/battle-tiles.tar
rm -f /tmp/battle-frontend.tar.gz
rm -f /tmp/battle-configs.tar.gz

# 7. 测试服务
echo ""
echo -e "${YELLOW}测试服务...${NC}"

sleep 2

if curl -s http://localhost:8000 > /dev/null; then
    echo -e "${GREEN}✓ 后端服务正常${NC}"
else
    echo -e "${YELLOW}⚠️  后端服务可能未就绪${NC}"
fi

if curl -s http://localhost > /dev/null; then
    echo -e "${GREEN}✓ 前端服务正常${NC}"
else
    echo -e "${YELLOW}⚠️  前端服务可能未就绪${NC}"
fi

# 完成
echo ""
echo -e "${GREEN}======================================"
echo "部署完成！"
echo "======================================${NC}"
echo ""
echo "访问地址："
SERVER_IP=$(hostname -I | awk '{print $1}')
echo "  后端 API: http://${SERVER_IP}:8000"
echo "  前端应用: http://${SERVER_IP}"
echo ""
echo "管理命令："
echo "  查看后端日志: docker logs -f battle-tiles"
echo "  查看前端日志: tail -f /var/log/nginx/access.log"
echo "  重启后端: docker restart battle-tiles"
echo "  重启前端: systemctl restart nginx"
