#!/bin/bash

# 前端构建和部署脚本

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================"
echo "前端构建与部署"
echo "======================================${NC}"

cd /Users/b022mc/project/battle/battle-reusables

# 1. 检查环境变量
echo -e "${YELLOW}步骤 1/4: 检查环境变量配置...${NC}"
if [ -f ".env.production" ]; then
    echo -e "${GREEN}✓ .env.production 已配置${NC}"
    echo ""
    cat .env.production
    echo ""
else
    echo -e "${RED}✗ 缺少 .env.production 文件${NC}"
    exit 1
fi

# 2. 安装依赖
echo ""
echo -e "${YELLOW}步骤 2/4: 安装依赖...${NC}"
npm install

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 依赖安装失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 依赖安装完成${NC}"

# 3. 构建 Web 版本
echo ""
echo -e "${YELLOW}步骤 3/4: 构建 Web 版本...${NC}"
npm run build:web

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 构建失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 构建完成！${NC}"
echo "构建产物位置: web-build/"

# 4. 选择部署平台
echo ""
echo -e "${YELLOW}步骤 4/4: 选择部署方式${NC}"
echo ""
echo "请选择部署平台："
echo "1) Vercel（推荐，免费 CDN）"
echo "2) Netlify（免费静态托管）"
echo "3) 上传到服务器 Nginx"
echo "4) 仅构建，不部署"
read -p "请选择 [1-4]: " -n 1 -r
echo ""

case $REPLY in
    1)
        echo -e "${YELLOW}部署到 Vercel...${NC}"
        
        if ! command -v vercel &> /dev/null; then
            echo "安装 Vercel CLI..."
            npm i -g vercel
        fi
        
        # 创建 vercel.json
        cat > vercel.json << 'EOF'
{
  "buildCommand": "npm run build:web",
  "outputDirectory": "web-build",
  "framework": "expo",
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
EOF
        
        vercel --prod
        ;;
        
    2)
        echo -e "${YELLOW}部署到 Netlify...${NC}"
        
        if ! command -v netlify &> /dev/null; then
            echo "安装 Netlify CLI..."
            npm i -g netlify-cli
        fi
        
        # 创建 netlify.toml
        cat > netlify.toml << 'EOF'
[build]
  command = "npm run build:web"
  publish = "web-build"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
EOF
        
        netlify deploy --prod --dir=web-build
        ;;
        
    3)
        echo -e "${YELLOW}上传到服务器...${NC}"
        
        SERVER_IP="8.137.52.203"
        SERVER_USER="root"
        
        echo "打包构建产物..."
        cd web-build
        tar -czf /tmp/battle-frontend.tar.gz .
        cd ..
        
        echo "上传到服务器..."
        scp /tmp/battle-frontend.tar.gz ${SERVER_USER}@${SERVER_IP}:/tmp/
        
        echo "在服务器上部署..."
        ssh ${SERVER_USER}@${SERVER_IP} << 'ENDSSH'
# 安装 Nginx
if ! command -v nginx &> /dev/null; then
    echo "安装 Nginx..."
    apt update && apt install nginx -y 2>/dev/null || yum install nginx -y
fi

# 解压前端文件
mkdir -p /var/www/battle
cd /var/www/battle
tar -xzf /tmp/battle-frontend.tar.gz

# 配置 Nginx
cat > /etc/nginx/conf.d/battle.conf << 'EOFNGINX'
server {
    listen 80;
    server_name _;
    
    root /var/www/battle;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    gzip on;
    gzip_types text/plain text/css text/javascript application/javascript application/json;
}
EOFNGINX

# 重启 Nginx
nginx -t && systemctl restart nginx

# 开放 80 端口
firewall-cmd --permanent --add-port=80/tcp 2>/dev/null || true
firewall-cmd --reload 2>/dev/null || true
ufw allow 80/tcp 2>/dev/null || true

# 清理
rm -f /tmp/battle-frontend.tar.gz

echo "前端部署完成！访问地址: http://$(hostname -I | awk '{print $1}')"
ENDSSH
        
        rm -f /tmp/battle-frontend.tar.gz
        
        echo -e "${GREEN}✓ 部署完成！${NC}"
        echo "前端地址: http://${SERVER_IP}"
        ;;
        
    4)
        echo -e "${GREEN}构建完成！${NC}"
        echo ""
        echo "构建产物位置: web-build/"
        echo ""
        echo "你可以："
        echo "1. 手动上传到静态托管平台"
        echo "2. 稍后运行 vercel 或 netlify 命令部署"
        ;;
        
    *)
        echo -e "${RED}无效选项${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}======================================"
echo "前端部署完成！"
echo "======================================${NC}"
