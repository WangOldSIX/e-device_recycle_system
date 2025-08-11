#!/bin/bash

echo "================================"
echo "二手电脑回收平台启动脚本"
echo "================================"

# 检查是否安装了必要的工具
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "错误: $1 未安装"
        exit 1
    fi
}

echo "检查环境依赖..."
check_command go
check_command node
check_command npm

echo "✓ 环境检查通过"
echo

# 启动后端服务
echo "1. 启动后端服务..."
cd backend
gnome-terminal --title="Go后端服务" -- bash -c "echo '启动Go后端服务...'; go run main.go; exec bash" &
cd ..

# 启动前端服务
echo "2. 启动前端开发服务..."
cd frontend
gnome-terminal --title="uni-app前端开发" -- bash -c "echo '启动uni-app前端开发...'; npm run dev:h5; exec bash" &
cd ..

echo
echo "服务启动完成！"
echo "后端地址: http://localhost:8080"
echo "前端地址: http://localhost:3000"
echo
echo "请确保已安装并配置："
echo "- Go 1.21+"
echo "- Node.js 16+"
echo "- MySQL 8.0+"
echo
echo "按任意键退出..."
read -n 1
