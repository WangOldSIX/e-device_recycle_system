package controllers

import (
	"e-device-recycle-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EvaluationController struct{}

// 创建评估（管理员/评估师）
func (ec *EvaluationController) CreateEvaluation(c *gin.Context) {
	evaluatorID, _ := c.Get("user_id")

	var req models.EvaluationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查订单是否存在
	var order models.RecycleOrder
	if err := models.DB.First(&order, req.OrderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单不存在"})
		return
	}

	// 检查是否已有评估
	var existingEvaluation models.Evaluation
	if err := models.DB.Where("order_id = ?", req.OrderID).First(&existingEvaluation).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该订单已有评估记录"})
		return
	}

	// 计算综合评分
	overallScore := float64(req.AppearanceScore+req.FunctionScore+req.PerformanceScore) / 3.0

	// 计算最终价格
	finalPrice := req.MarketPrice * (1 - req.DepreciationRate) * (overallScore / 10.0)

	// 创建评估
	evaluation := models.Evaluation{
		OrderID:          req.OrderID,
		EvaluatorID:      evaluatorID.(uint),
		AppearanceScore:  req.AppearanceScore,
		FunctionScore:    req.FunctionScore,
		PerformanceScore: req.PerformanceScore,
		OverallScore:     overallScore,
		MarketPrice:      req.MarketPrice,
		DepreciationRate: req.DepreciationRate,
		FinalPrice:       finalPrice,
		EvaluationReport: req.EvaluationReport,
		Images:           req.Images,
		Status:           "completed",
	}

	if err := models.DB.Create(&evaluation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评估失败"})
		return
	}

	// 更新订单状态和最终价格
	models.DB.Model(&order).Updates(map[string]interface{}{
		"status":      "evaluated",
		"final_price": finalPrice,
	})

	// 预加载关联数据
	models.DB.Preload("Evaluator").First(&evaluation, evaluation.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "评估创建成功",
		"evaluation": ec.convertToResponse(evaluation),
	})
}

// 获取评估列表（管理员）
func (ec *EvaluationController) GetEvaluations(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 过滤参数
	status := c.Query("status")
	evaluatorID := c.Query("evaluator_id")

	query := models.DB.Model(&models.Evaluation{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if evaluatorID != "" {
		query = query.Where("evaluator_id = ?", evaluatorID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取评估数据
	var evaluations []models.Evaluation
	if err := query.Preload("Order").Preload("Evaluator").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&evaluations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评估列表失败"})
		return
	}

	// 转换响应
	var evaluationResponses []models.EvaluationResponse
	for _, evaluation := range evaluations {
		evaluationResponses = append(evaluationResponses, ec.convertToResponse(evaluation))
	}

	c.JSON(http.StatusOK, gin.H{
		"evaluations": evaluationResponses,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// 获取评估详情
func (ec *EvaluationController) GetEvaluation(c *gin.Context) {
	id := c.Param("id")

	var evaluation models.Evaluation
	if err := models.DB.Preload("Order").Preload("Evaluator").
		First(&evaluation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "评估不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评估信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"evaluation": ec.convertToResponse(evaluation),
	})
}

// 根据订单ID获取评估
func (ec *EvaluationController) GetEvaluationByOrder(c *gin.Context) {
	orderID := c.Param("order_id")

	var evaluation models.Evaluation
	if err := models.DB.Preload("Order").Preload("Evaluator").
		Where("order_id = ?", orderID).First(&evaluation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "该订单暂无评估"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评估信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"evaluation": ec.convertToResponse(evaluation),
	})
}

// 更新评估
func (ec *EvaluationController) UpdateEvaluation(c *gin.Context) {
	id := c.Param("id")
	evaluatorID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var evaluation models.Evaluation
	query := models.DB.Where("id = ?", id)

	// 普通评估师只能更新自己的评估
	if role != "admin" {
		query = query.Where("evaluator_id = ?", evaluatorID)
	}

	if err := query.First(&evaluation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评估不存在或无权限"})
		return
	}

	var req models.EvaluationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新计算综合评分和最终价格
	overallScore := float64(req.AppearanceScore+req.FunctionScore+req.PerformanceScore) / 3.0
	finalPrice := req.MarketPrice * (1 - req.DepreciationRate) * (overallScore / 10.0)

	// 更新评估
	updates := map[string]interface{}{
		"appearance_score":  req.AppearanceScore,
		"function_score":    req.FunctionScore,
		"performance_score": req.PerformanceScore,
		"overall_score":     overallScore,
		"market_price":      req.MarketPrice,
		"depreciation_rate": req.DepreciationRate,
		"final_price":       finalPrice,
		"evaluation_report": req.EvaluationReport,
		"images":            req.Images,
		"status":            req.Status,
	}

	if err := models.DB.Model(&evaluation).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评估失败"})
		return
	}

	// 同时更新订单的最终价格
	models.DB.Model(&models.RecycleOrder{}).
		Where("id = ?", evaluation.OrderID).
		Update("final_price", finalPrice)

	// 重新加载评估数据
	models.DB.Preload("Order").Preload("Evaluator").First(&evaluation, evaluation.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":    "评估更新成功",
		"evaluation": ec.convertToResponse(evaluation),
	})
}

// 转换为响应格式
func (ec *EvaluationController) convertToResponse(evaluation models.Evaluation) models.EvaluationResponse {
	response := models.EvaluationResponse{
		ID:               evaluation.ID,
		OrderID:          evaluation.OrderID,
		EvaluatorID:      evaluation.EvaluatorID,
		AppearanceScore:  evaluation.AppearanceScore,
		FunctionScore:    evaluation.FunctionScore,
		PerformanceScore: evaluation.PerformanceScore,
		OverallScore:     evaluation.OverallScore,
		MarketPrice:      evaluation.MarketPrice,
		DepreciationRate: evaluation.DepreciationRate,
		FinalPrice:       evaluation.FinalPrice,
		EvaluationReport: evaluation.EvaluationReport,
		Images:           evaluation.Images,
		Status:           evaluation.Status,
		CreatedAt:        evaluation.CreatedAt,
	}

	// 添加评估师信息
	if evaluation.Evaluator.ID != 0 {
		response.Evaluator = &models.UserResponse{
			ID:       evaluation.Evaluator.ID,
			Username: evaluation.Evaluator.Username,
			Phone:    evaluation.Evaluator.Phone,
			Email:    evaluation.Evaluator.Email,
			RealName: evaluation.Evaluator.RealName,
			Avatar:   evaluation.Evaluator.Avatar,
			Role:     evaluation.Evaluator.Role,
			Status:   evaluation.Evaluator.Status,
		}
	}

	return response
}
