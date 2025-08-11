package main

import (
	"e-device-recycle-backend/config"
	"e-device-recycle-backend/models"
	"e-device-recycle-backend/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	models.InitDB()

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := config.GetConfig().Port
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口: %s", port)
	log.Fatal(r.Run(":" + port))
}
