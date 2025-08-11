package controllers

import (
	"e-device-recycle-backend/models"
	"e-device-recycle-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct{}

// 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := models.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if err := models.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "手机号已存在"})
			return
		}
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Role:     "user",
		Status:   "active",
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"token":   token,
		"user": models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Phone:    user.Phone,
			Email:    user.Email,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}

// 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "账户已被禁用"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Phone:    user.Phone,
			Email:    user.Email,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}

// 获取用户信息
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Phone:    user.Phone,
			Email:    user.Email,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}

// 更新用户信息
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 过滤可更新的字段
	allowedFields := []string{"phone", "email", "real_name", "avatar"}
	filteredUpdates := make(map[string]interface{})
	for _, field := range allowedFields {
		if value, exists := updates[field]; exists {
			filteredUpdates[field] = value
		}
	}

	if err := models.DB.Model(&user).Updates(filteredUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user": models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Phone:    user.Phone,
			Email:    user.Email,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}
