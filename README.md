# 二手电脑回收平台

一个基于Go后端和uni-app前端的二手电脑回收平台，支持多端发布（H5、小程序、APP等）。

## 项目特点

### 🚀 技术栈
- **后端**: Go + Gin + GORM + MySQL
- **前端**: uni-app (Vue3) + Pinia
- **数据库**: MySQL
- **认证**: JWT
- **部署**: 支持多平台部署

### 📱 多端支持
- H5网页版
- 微信小程序
- 支付宝小程序  
- Android/iOS APP
- 快应用等

### 🎯 核心功能
- 用户注册登录
- 设备信息管理
- 在线评估报价
- 上门回收预约
- 订单状态跟踪
- 专业设备评估
- 管理后台

## 项目结构

```
e-device_recycle/
├── backend/                 # Go后端
│   ├── main.go             # 应用入口
│   ├── go.mod              # Go模块文件
│   ├── config/             # 配置文件
│   ├── models/             # 数据模型
│   ├── controllers/        # 控制器
│   ├── middleware/         # 中间件
│   ├── routes/             # 路由
│   ├── utils/              # 工具函数
│   └── scripts/            # 数据库脚本
└── frontend/               # uni-app前端
    ├── App.vue             # 应用根组件
    ├── main.js             # 应用入口
    ├── manifest.json       # 应用配置
    ├── pages.json          # 页面配置
    ├── package.json        # 依赖配置
    ├── store/              # 状态管理
    └── pages/              # 页面文件
        ├── index/          # 首页
        ├── devices/        # 设备相关
        ├── order/          # 订单相关
        └── user/           # 用户相关
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- MySQL 8.0+
- HBuilderX (推荐) 或 CLI

### 后端启动

1. **配置数据库**
   ```bash
   # 创建数据库
   mysql -u root -p
   CREATE DATABASE device_recycle CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. **安装依赖**
   ```bash
   cd backend
   go mod tidy
   ```

3. **配置环境变量**
   ```bash
   # 复制配置文件
   cp config.example .env
   
   # 编辑配置文件，设置数据库连接信息
   vim .env
   ```

4. **初始化数据库**
   ```bash
   # 导入初始化脚本
   mysql -u root -p device_recycle < scripts/init_db.sql
   ```

5. **启动服务**
   ```bash
   go run main.go
   ```

   服务将在 `http://localhost:8080` 启动

### 前端启动

1. **安装依赖**
   ```bash
   cd frontend
   npm install
   ```

2. **启动开发服务器**
   ```bash
   # H5开发
   npm run dev:h5
   
   # 微信小程序开发
   npm run dev:mp-weixin
   
   # 支付宝小程序开发  
   npm run dev:mp-alipay
   ```

3. **使用HBuilderX开发**
   - 用HBuilderX打开frontend目录
   - 点击"运行" -> 选择对应平台
   - 或点击"发行" -> 选择对应平台

## API接口文档

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 设备相关
- `GET /api/v1/devices` - 获取设备列表
- `GET /api/v1/devices/:id` - 获取设备详情

### 订单相关  
- `POST /api/v1/orders` - 创建回收订单
- `GET /api/v1/orders` - 获取用户订单列表
- `GET /api/v1/orders/:id` - 获取订单详情
- `PUT /api/v1/orders/:id/cancel` - 取消订单

### 用户相关
- `GET /api/v1/user/profile` - 获取用户信息
- `PUT /api/v1/user/profile` - 更新用户信息

### 管理员接口
- `POST /api/v1/admin/devices` - 创建设备
- `PUT /api/v1/admin/devices/:id` - 更新设备
- `DELETE /api/v1/admin/devices/:id` - 删除设备
- `GET /api/v1/admin/orders` - 获取所有订单
- `PUT /api/v1/admin/orders/:id` - 更新订单状态
- `POST /api/v1/admin/evaluations` - 创建评估
- `GET /api/v1/admin/evaluations` - 获取评估列表

## 数据库设计

### 用户表 (users)
- 用户基本信息
- 角色权限管理
- 登录认证信息

### 设备表 (devices)  
- 设备型号信息
- 技术规格参数
- 基础回收价格

### 回收订单表 (recycle_orders)
- 订单基本信息
- 联系人信息
- 上门地址时间
- 订单状态流转

### 评估表 (evaluations)
- 专业评估结果
- 各项评分详情
- 最终回收价格

## 业务流程

### 用户回收流程
1. **浏览设备** - 查看支持回收的设备型号
2. **提交申请** - 填写设备信息和联系方式
3. **预约上门** - 选择合适的上门时间
4. **现场评估** - 专业人员上门检测评估
5. **确认交易** - 当面验收，即时付款
6. **完成回收** - 数据清除，完成交易

### 管理员工作流程
1. **订单管理** - 查看和处理回收申请
2. **派单调度** - 安排评估师上门服务  
3. **设备评估** - 专业检测设备状况
4. **价格确认** - 根据评估结果定价
5. **完成交易** - 确认回收完成

## 部署说明

### 后端部署
```bash
# 构建可执行文件
go build -o device-recycle-server main.go

# 运行服务
./device-recycle-server
```

### 前端部署

1. **H5版本**
   ```bash
   npm run build:h5
   # 部署 dist/build/h5 目录到web服务器
   ```

2. **小程序版本**  
   ```bash
   # 微信小程序
   npm run build:mp-weixin
   # 使用微信开发者工具打开 dist/build/mp-weixin
   
   # 支付宝小程序
   npm run build:mp-alipay  
   # 使用支付宝开发者工具打开 dist/build/mp-alipay
   ```

3. **APP版本**
   ```bash
   npm run build:app
   # 使用HBuilderX云打包或本地打包
   ```

## 开发说明

### 添加新页面
1. 在 `frontend/pages/` 下创建页面文件
2. 在 `pages.json` 中配置页面路由
3. 更新导航和链接

### 添加新API
1. 在 `backend/controllers/` 创建控制器
2. 在 `backend/routes/` 配置路由
3. 更新API文档

### 数据库变更
1. 修改 `backend/models/` 中的模型
2. 使用GORM自动迁移
3. 或编写迁移脚本

## 注意事项

1. **安全性**
   - 所有API接口都有JWT认证保护
   - 用户密码使用bcrypt加密存储
   - 敏感操作需要管理员权限

2. **性能优化**
   - 数据库查询使用索引优化
   - 图片上传建议使用CDN
   - 考虑使用Redis缓存热点数据

3. **跨平台兼容**
   - 使用uni-app标准组件和API
   - 避免平台特定的代码
   - 测试各个目标平台

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交改动 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目地址:(https://github.com/WangOldSIX/e-device_recycle_system)
- 邮箱: 1447895999@qq.com

---

**二手电脑回收平台** - 让闲置设备重新焕发价值！💻♻️
