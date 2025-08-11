# 部署指南

本文档描述如何部署二手电脑回收平台到生产环境。

## 部署方式

### 1. Docker部署（推荐）

使用Docker Compose一键部署所有服务：

```bash
# 1. 克隆项目
git clone <repository-url>
cd e-device_recycle

# 2. 构建前端
cd frontend
npm install
npm run build:h5

# 3. 启动所有服务
cd ..
docker-compose up -d

# 4. 查看日志
docker-compose logs -f
```

服务访问地址：
- 前端：http://localhost:80
- 后端API：http://localhost:80/api
- MySQL：localhost:3306
- Redis：localhost:6379

### 2. 手动部署

#### 2.1 环境准备

**服务器要求：**
- CPU: 2核心以上
- 内存: 4GB以上
- 存储: 50GB以上
- 操作系统: Ubuntu 20.04+ / CentOS 8+

**软件依赖：**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y golang-go nodejs npm mysql-server nginx

# CentOS/RHEL
sudo yum install -y golang nodejs npm mysql-server nginx
```

#### 2.2 数据库配置

```bash
# 1. 启动MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# 2. 安全设置
sudo mysql_secure_installation

# 3. 创建数据库和用户
mysql -u root -p
```

```sql
CREATE DATABASE device_recycle CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'app'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON device_recycle.* TO 'app'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

```bash
# 4. 导入初始数据
mysql -u app -p device_recycle < backend/scripts/init_db.sql
```

#### 2.3 后端部署

```bash
# 1. 进入后端目录
cd backend

# 2. 安装依赖
go mod tidy

# 3. 配置环境变量
cp config.example .env
vim .env
```

编辑.env文件：
```bash
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=app
DB_PASS=your_password
DB_NAME=device_recycle
JWT_SECRET=your-secret-key
```

```bash
# 4. 构建应用
go build -o device-recycle-server main.go

# 5. 创建systemd服务
sudo vim /etc/systemd/system/device-recycle.service
```

systemd服务配置：
```ini
[Unit]
Description=Device Recycle Backend
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/backend
ExecStart=/path/to/backend/device-recycle-server
Restart=always
RestartSec=5
Environment=PATH=/usr/bin:/usr/local/bin
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

```bash
# 6. 启动服务
sudo systemctl daemon-reload
sudo systemctl start device-recycle
sudo systemctl enable device-recycle
```

#### 2.4 前端部署

```bash
# 1. 进入前端目录
cd frontend

# 2. 安装依赖
npm install

# 3. 构建生产版本
npm run build:h5

# 4. 复制到nginx目录
sudo cp -r dist/build/h5/* /var/www/html/
```

#### 2.5 Nginx配置

```bash
# 1. 配置nginx
sudo vim /etc/nginx/sites-available/device-recycle
```

Nginx配置：
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

```bash
# 2. 启用站点
sudo ln -s /etc/nginx/sites-available/device-recycle /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 3. SSL证书配置

使用Let's Encrypt免费SSL证书：

```bash
# 1. 安装certbot
sudo apt install certbot python3-certbot-nginx

# 2. 获取证书
sudo certbot --nginx -d your-domain.com

# 3. 自动续期
sudo crontab -e
# 添加：0 12 * * * /usr/bin/certbot renew --quiet
```

## 监控和维护

### 1. 日志查看

```bash
# 应用日志
sudo journalctl -u device-recycle -f

# Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# MySQL日志
sudo tail -f /var/log/mysql/error.log
```

### 2. 性能监控

```bash
# 系统资源
htop
df -h
free -h

# 应用状态
sudo systemctl status device-recycle
sudo systemctl status nginx
sudo systemctl status mysql
```

### 3. 备份策略

**数据库备份：**
```bash
# 创建备份脚本
vim backup.sh
```

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/mysql"
mkdir -p $BACKUP_DIR

mysqldump -u app -p device_recycle > $BACKUP_DIR/device_recycle_$DATE.sql

# 保留最近7天的备份
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
```

```bash
# 设置定时备份
chmod +x backup.sh
crontab -e
# 添加：0 2 * * * /path/to/backup.sh
```

**文件备份：**
```bash
# 备份上传的文件
rsync -av /path/to/uploads/ /backup/uploads/
```

### 4. 更新部署

```bash
# 1. 备份数据
./backup.sh

# 2. 更新代码
git pull origin main

# 3. 重新构建
cd backend
go build -o device-recycle-server main.go

cd ../frontend
npm run build:h5
sudo cp -r dist/build/h5/* /var/www/html/

# 4. 重启服务
sudo systemctl restart device-recycle
sudo systemctl reload nginx
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接配置
   - 查看防火墙设置

2. **API请求失败**
   - 检查后端服务状态
   - 查看应用日志
   - 验证Nginx配置

3. **前端页面无法加载**
   - 检查静态文件路径
   - 验证Nginx配置
   - 查看浏览器控制台错误

### 性能优化

1. **数据库优化**
   - 添加索引
   - 配置连接池
   - 定期清理日志

2. **应用优化**
   - 启用Gzip压缩
   - 配置缓存策略
   - 使用CDN加速静态资源

3. **服务器优化**
   - 调整Nginx工作进程数
   - 配置系统资源限制
   - 监控系统性能指标

## 安全建议

1. **网络安全**
   - 配置防火墙规则
   - 使用SSL/TLS加密
   - 定期更新系统补丁

2. **应用安全**
   - 强化JWT密钥
   - 配置请求速率限制
   - 验证用户输入

3. **数据安全**
   - 定期备份数据
   - 加密敏感信息
   - 控制数据库访问权限

---

如有问题，请查看项目文档或联系技术支持。
