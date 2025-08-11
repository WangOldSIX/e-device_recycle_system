package controllers

import (
	"e-device-recycle-backend/models"
	"e-device-recycle-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecycleOrderController struct{}

// 创建回收订单
func (roc *RecycleOrderController) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.RecycleOrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := models.DB.Where("id = ? AND status = ?", req.DeviceID, "active").First(&device).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备不存在"})
		return
	}

	// 计算预估价格
	estimatedPrice := utils.CalculateDevicePrice(device.BasePrice, device.Condition, device.YearBought)

	// 创建订单
	order := models.RecycleOrder{
		UserID:         userID.(uint),
		DeviceID:       req.DeviceID,
		OrderNo:        utils.GenerateOrderNo(),
		ContactName:    req.ContactName,
		ContactPhone:   req.ContactPhone,
		PickupAddress:  req.PickupAddress,
		PickupTime:     req.PickupTime,
		DeviceInfo:     req.DeviceInfo,
		Images:         req.Images,
		EstimatedPrice: estimatedPrice,
		Status:         "pending",
		Remark:         req.Remark,
	}

	if err := models.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 预加载关联数据
	models.DB.Preload("User").Preload("Device").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "订单创建成功",
		"order":   roc.convertToResponse(order),
	})
}

// 获取用户订单列表
func (roc *RecycleOrderController) GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 状态过滤
	status := c.Query("status")

	query := models.DB.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Model(&models.RecycleOrder{}).Count(&total)

	// 获取订单数据
	var orders []models.RecycleOrder
	if err := query.Preload("Device").Preload("Evaluation").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单列表失败"})
		return
	}

	// 转换响应
	var orderResponses []models.RecycleOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, roc.convertToResponse(order))
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderResponses,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// 获取所有订单列表（管理员）
func (roc *RecycleOrderController) GetAllOrders(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 过滤参数
	status := c.Query("status")
	userID := c.Query("user_id")

	query := models.DB.Model(&models.RecycleOrder{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取订单数据
	var orders []models.RecycleOrder
	if err := query.Preload("User").Preload("Device").Preload("Evaluation").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单列表失败"})
		return
	}

	// 转换响应
	var orderResponses []models.RecycleOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, roc.convertToResponse(order))
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderResponses,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// 获取订单详情
func (roc *RecycleOrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var order models.RecycleOrder
	query := models.DB.Preload("User").Preload("Device").Preload("Evaluation")

	// 普通用户只能查看自己的订单
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": roc.convertToResponse(order),
	})
}

// 更新订单状态（管理员）
func (roc *RecycleOrderController) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.RecycleOrder
	if err := models.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var req models.RecycleOrderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新订单
	updates := map[string]interface{}{
		"status": req.Status,
		"remark": req.Remark,
	}

	if req.FinalPrice != nil {
		updates["final_price"] = *req.FinalPrice
	}

	if req.PickupTime != nil {
		updates["pickup_time"] = *req.PickupTime
	}

	if err := models.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单失败"})
		return
	}

	// 重新加载订单数据
	models.DB.Preload("User").Preload("Device").Preload("Evaluation").First(&order, order.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "订单更新成功",
		"order":   roc.convertToResponse(order),
	})
}

// 取消订单
func (roc *RecycleOrderController) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var order models.RecycleOrder
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 只有pending状态的订单可以取消
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许取消"})
		return
	}

	if err := models.DB.Model(&order).Update("status", "cancelled").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单取消成功"})
}

// 转换为响应格式
func (roc *RecycleOrderController) convertToResponse(order models.RecycleOrder) models.RecycleOrderResponse {
	response := models.RecycleOrderResponse{
		ID:             order.ID,
		UserID:         order.UserID,
		DeviceID:       order.DeviceID,
		OrderNo:        order.OrderNo,
		ContactName:    order.ContactName,
		ContactPhone:   order.ContactPhone,
		PickupAddress:  order.PickupAddress,
		PickupTime:     order.PickupTime,
		DeviceInfo:     order.DeviceInfo,
		Images:         order.Images,
		EstimatedPrice: order.EstimatedPrice,
		FinalPrice:     order.FinalPrice,
		Status:         order.Status,
		Remark:         order.Remark,
		CreatedAt:      order.CreatedAt,
	}

	// 添加用户信息
	if order.User.ID != 0 {
		response.User = &models.UserResponse{
			ID:       order.User.ID,
			Username: order.User.Username,
			Phone:    order.User.Phone,
			Email:    order.User.Email,
			RealName: order.User.RealName,
			Avatar:   order.User.Avatar,
			Role:     order.User.Role,
			Status:   order.User.Status,
		}
	}

	// 添加设备信息
	if order.Device.ID != 0 {
		response.Device = &models.DeviceResponse{
			ID:          order.Device.ID,
			Name:        order.Device.Name,
			Brand:       order.Device.Brand,
			Model:       order.Device.Model,
			Category:    order.Device.Category,
			CPU:         order.Device.CPU,
			Memory:      order.Device.Memory,
			Storage:     order.Device.Storage,
			Graphics:    order.Device.Graphics,
			Screen:      order.Device.Screen,
			Condition:   order.Device.Condition,
			YearBought:  order.Device.YearBought,
			BasePrice:   order.Device.BasePrice,
			Description: order.Device.Description,
			Images:      order.Device.Images,
			Status:      order.Device.Status,
		}
	}

	// 添加评估信息
	if order.Evaluation != nil && order.Evaluation.ID != 0 {
		response.Evaluation = &models.EvaluationResponse{
			ID:               order.Evaluation.ID,
			OrderID:          order.Evaluation.OrderID,
			EvaluatorID:      order.Evaluation.EvaluatorID,
			AppearanceScore:  order.Evaluation.AppearanceScore,
			FunctionScore:    order.Evaluation.FunctionScore,
			PerformanceScore: order.Evaluation.PerformanceScore,
			OverallScore:     order.Evaluation.OverallScore,
			MarketPrice:      order.Evaluation.MarketPrice,
			DepreciationRate: order.Evaluation.DepreciationRate,
			FinalPrice:       order.Evaluation.FinalPrice,
			EvaluationReport: order.Evaluation.EvaluationReport,
			Images:           order.Evaluation.Images,
			Status:           order.Evaluation.Status,
			CreatedAt:        order.Evaluation.CreatedAt,
		}
	}

	return response
}
