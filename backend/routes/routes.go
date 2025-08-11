package routes

import (
	"e-device-recycle-backend/controllers"
	"e-device-recycle-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 创建控制器实例
	userController := &controllers.UserController{}
	deviceController := &controllers.DeviceController{}
	recycleOrderController := &controllers.RecycleOrderController{}
	evaluationController := &controllers.EvaluationController{}

	// API版本分组
	v1 := r.Group("/api/v1")

	// 公开路由
	{
		// 用户认证
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// 设备信息（公开查看）
		devices := v1.Group("/devices")
		{
			devices.GET("/", deviceController.GetDevices)
			devices.GET("/:id", deviceController.GetDevice)
		}
	}

	// 需要认证的路由
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关
		user := protected.Group("/user")
		{
			user.GET("/profile", userController.GetProfile)
			user.PUT("/profile", userController.UpdateProfile)
		}

		// 回收订单
		orders := protected.Group("/orders")
		{
			orders.POST("/", recycleOrderController.CreateOrder)
			orders.GET("/", recycleOrderController.GetUserOrders)
			orders.GET("/:id", recycleOrderController.GetOrder)
			orders.PUT("/:id/cancel", recycleOrderController.CancelOrder)
		}

		// 评估相关（普通用户可查看）
		evaluations := protected.Group("/evaluations")
		{
			evaluations.GET("/order/:order_id", evaluationController.GetEvaluationByOrder)
		}
	}

	// 管理员路由
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.AdminAuth())
	{
		// 设备管理
		devices := admin.Group("/devices")
		{
			devices.POST("/", deviceController.CreateDevice)
			devices.PUT("/:id", deviceController.UpdateDevice)
			devices.DELETE("/:id", deviceController.DeleteDevice)
		}

		// 订单管理
		orders := admin.Group("/orders")
		{
			orders.GET("/", recycleOrderController.GetAllOrders)
			orders.PUT("/:id", recycleOrderController.UpdateOrder)
		}

		// 评估管理
		evaluations := admin.Group("/evaluations")
		{
			evaluations.POST("/", evaluationController.CreateEvaluation)
			evaluations.GET("/", evaluationController.GetEvaluations)
			evaluations.GET("/:id", evaluationController.GetEvaluation)
			evaluations.PUT("/:id", evaluationController.UpdateEvaluation)
		}
	}
}
