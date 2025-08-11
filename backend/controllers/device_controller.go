package controllers

import (
	"e-device-recycle-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeviceController struct{}

// 获取设备列表
func (dc *DeviceController) GetDevices(c *gin.Context) {
	var devices []models.Device

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 过滤参数
	category := c.Query("category")
	brand := c.Query("brand")
	condition := c.Query("condition")

	query := models.DB.Where("status = ?", "active")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if brand != "" {
		query = query.Where("brand LIKE ?", "%"+brand+"%")
	}
	if condition != "" {
		query = query.Where("condition = ?", condition)
	}

	// 获取总数
	var total int64
	query.Model(&models.Device{}).Count(&total)

	// 获取数据
	if err := query.Offset(offset).Limit(pageSize).Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设备列表失败"})
		return
	}

	// 转换响应
	var deviceResponses []models.DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, models.DeviceResponse{
			ID:          device.ID,
			Name:        device.Name,
			Brand:       device.Brand,
			Model:       device.Model,
			Category:    device.Category,
			CPU:         device.CPU,
			Memory:      device.Memory,
			Storage:     device.Storage,
			Graphics:    device.Graphics,
			Screen:      device.Screen,
			Condition:   device.Condition,
			YearBought:  device.YearBought,
			BasePrice:   device.BasePrice,
			Description: device.Description,
			Images:      device.Images,
			Status:      device.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": deviceResponses,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// 获取设备详情
func (dc *DeviceController) GetDevice(c *gin.Context) {
	id := c.Param("id")

	var device models.Device
	if err := models.DB.Where("id = ? AND status = ?", id, "active").First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设备信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device": models.DeviceResponse{
			ID:          device.ID,
			Name:        device.Name,
			Brand:       device.Brand,
			Model:       device.Model,
			Category:    device.Category,
			CPU:         device.CPU,
			Memory:      device.Memory,
			Storage:     device.Storage,
			Graphics:    device.Graphics,
			Screen:      device.Screen,
			Condition:   device.Condition,
			YearBought:  device.YearBought,
			BasePrice:   device.BasePrice,
			Description: device.Description,
			Images:      device.Images,
			Status:      device.Status,
		},
	})
}

// 创建设备（管理员功能）
func (dc *DeviceController) CreateDevice(c *gin.Context) {
	var req models.DeviceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device := models.Device{
		Name:        req.Name,
		Brand:       req.Brand,
		Model:       req.Model,
		Category:    req.Category,
		CPU:         req.CPU,
		Memory:      req.Memory,
		Storage:     req.Storage,
		Graphics:    req.Graphics,
		Screen:      req.Screen,
		Condition:   req.Condition,
		YearBought:  req.YearBought,
		BasePrice:   req.BasePrice,
		Description: req.Description,
		Images:      req.Images,
		Status:      "active",
	}

	if err := models.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建设备失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "设备创建成功",
		"device": models.DeviceResponse{
			ID:          device.ID,
			Name:        device.Name,
			Brand:       device.Brand,
			Model:       device.Model,
			Category:    device.Category,
			CPU:         device.CPU,
			Memory:      device.Memory,
			Storage:     device.Storage,
			Graphics:    device.Graphics,
			Screen:      device.Screen,
			Condition:   device.Condition,
			YearBought:  device.YearBought,
			BasePrice:   device.BasePrice,
			Description: device.Description,
			Images:      device.Images,
			Status:      device.Status,
		},
	})
}

// 更新设备（管理员功能）
func (dc *DeviceController) UpdateDevice(c *gin.Context) {
	id := c.Param("id")

	var device models.Device
	if err := models.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Model(&device).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "设备更新成功",
		"device": models.DeviceResponse{
			ID:          device.ID,
			Name:        device.Name,
			Brand:       device.Brand,
			Model:       device.Model,
			Category:    device.Category,
			CPU:         device.CPU,
			Memory:      device.Memory,
			Storage:     device.Storage,
			Graphics:    device.Graphics,
			Screen:      device.Screen,
			Condition:   device.Condition,
			YearBought:  device.YearBought,
			BasePrice:   device.BasePrice,
			Description: device.Description,
			Images:      device.Images,
			Status:      device.Status,
		},
	})
}

// 删除设备（管理员功能）
func (dc *DeviceController) DeleteDevice(c *gin.Context) {
	id := c.Param("id")

	var device models.Device
	if err := models.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	// 软删除
	if err := models.DB.Model(&device).Update("status", "inactive").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备删除成功"})
}
