@echo off
echo ================================
echo 二手电脑回收平台启动脚本
echo ================================

echo.
echo 1. 启动后端服务...
cd backend
start cmd /k "echo 启动Go后端服务 && go run main.go"

echo.
echo 2. 启动前端开发服务...
cd ../frontend
start cmd /k "echo 启动uni-app前端开发 && npm run dev:h5"

echo.
echo 服务启动完成！
echo 后端地址: http://localhost:8080
echo 前端地址: http://localhost:3000
echo.
echo 请确保已安装并配置：
echo - Go 1.21+
echo - Node.js 16+
echo - MySQL 8.0+
echo.
pause
