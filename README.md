# Battle 项目部署

## 🎉 部署完成！

- ✅ **后端**: http://8.137.52.203:8000
- ✅ **前端**: http://8.137.52.203

---

## 📋 快速部署命令

### 首次部署

```bash
# 1. 配置 SSH 免密登录（可选）
./配置SSH免密登录.sh

# 2. 部署后端
./后端直接编译上传.sh

# 3. 部署前端
./前端上传到服务器.sh
```

### 更新部署

```bash
# 更新后端
./后端直接编译上传.sh

# 更新前端
./前端上传到服务器.sh
```

---

## 🔧 可用脚本

- **`后端直接编译上传.sh`** - 后端部署脚本
  - 交叉编译到 Linux
  - 上传并启动服务
  
- **`前端上传到服务器.sh`** - 前端部署脚本
  - 构建并上传前端
  - 配置 Nginx
  
- **`配置SSH免密登录.sh`** - SSH 免密配置
  - 避免每次输入密码

---

## 📁 项目配置

### 前端环境变量

`battle-reusables/.env.production`：

```bash
EXPO_PUBLIC_DEV_API_URL=http://8.137.52.203:8000
EXPO_PUBLIC_API_HOST=8.137.52.203
EXPO_PUBLIC_BYPASS_PROXY=true
```

### 后端配置

`battle-tiles/configs/config.yaml` - 数据库和 Redis 配置

---

## 🌐 访问地址

- **前端应用**: http://8.137.52.203
- **后端 API**: http://8.137.52.203:8000

---

## 🔧 服务管理

### 后端

```bash
# 查看日志
ssh root@8.137.52.203 'tail -f /root/battle-tiles/logs/platform.log'

# 查看进程
ssh root@8.137.52.203 'ps aux | grep go-kgin'

# 重启服务
./后端直接编译上传.sh
```

### 前端

```bash
# 查看 Nginx 状态
ssh root@8.137.52.203 'systemctl status nginx'

# 查看 Nginx 日志
ssh root@8.137.52.203 'tail -f /var/log/nginx/access.log'

# 重启 Nginx
ssh root@8.137.52.203 'systemctl restart nginx'

# 更新前端
./前端上传到服务器.sh
```

---

## ⚠️ 防火墙配置

确保在阿里云控制台的**安全组规则**中开放：

- **80 端口** - HTTP (前端)
- **8000 端口** - 后端 API

---

## 📖 详细文档

查看 `部署说明.md` 了解更多细节和故障排查。

---

## 🚀 技术栈

- **前端**: Expo + React Native Web
- **后端**: Go + Kratos
- **服务器**: Nginx + Linux
- **数据库**: PostgreSQL
- **缓存**: Redis
