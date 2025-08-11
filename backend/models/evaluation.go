package models

import (
	"time"

	"gorm.io/gorm"
)

type Evaluation struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	OrderID          uint           `json:"order_id" gorm:"uniqueIndex;not null"`
	EvaluatorID      uint           `json:"evaluator_id"`                       // 评估师ID
	AppearanceScore  int            `json:"appearance_score" gorm:"default:0"`  // 外观评分 1-10
	FunctionScore    int            `json:"function_score" gorm:"default:0"`    // 功能评分 1-10
	PerformanceScore int            `json:"performance_score" gorm:"default:0"` // 性能评分 1-10
	OverallScore     float64        `json:"overall_score"`                      // 综合评分
	MarketPrice      float64        `json:"market_price"`                       // 市场参考价
	DepreciationRate float64        `json:"depreciation_rate"`                  // 折旧率
	FinalPrice       float64        `json:"final_price"`                        // 最终评估价格
	EvaluationReport string         `json:"evaluation_report"`                  // 评估报告
	Images           string         `json:"images"`                             // 评估图片
	Status           string         `json:"status" gorm:"default:'pending'"`    // pending, completed
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Order     RecycleOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Evaluator User         `json:"evaluator,omitempty" gorm:"foreignKey:EvaluatorID"`
}

type EvaluationCreateRequest struct {
	OrderID          uint    `json:"order_id" binding:"required"`
	AppearanceScore  int     `json:"appearance_score" binding:"min=1,max=10"`
	FunctionScore    int     `json:"function_score" binding:"min=1,max=10"`
	PerformanceScore int     `json:"performance_score" binding:"min=1,max=10"`
	MarketPrice      float64 `json:"market_price" binding:"min=0"`
	DepreciationRate float64 `json:"depreciation_rate" binding:"min=0,max=1"`
	EvaluationReport string  `json:"evaluation_report"`
	Images           string  `json:"images"`
}

type EvaluationUpdateRequest struct {
	AppearanceScore  int     `json:"appearance_score" binding:"min=1,max=10"`
	FunctionScore    int     `json:"function_score" binding:"min=1,max=10"`
	PerformanceScore int     `json:"performance_score" binding:"min=1,max=10"`
	MarketPrice      float64 `json:"market_price" binding:"min=0"`
	DepreciationRate float64 `json:"depreciation_rate" binding:"min=0,max=1"`
	EvaluationReport string  `json:"evaluation_report"`
	Images           string  `json:"images"`
	Status           string  `json:"status" binding:"oneof=pending completed"`
}

type EvaluationResponse struct {
	ID               uint          `json:"id"`
	OrderID          uint          `json:"order_id"`
	EvaluatorID      uint          `json:"evaluator_id"`
	AppearanceScore  int           `json:"appearance_score"`
	FunctionScore    int           `json:"function_score"`
	PerformanceScore int           `json:"performance_score"`
	OverallScore     float64       `json:"overall_score"`
	MarketPrice      float64       `json:"market_price"`
	DepreciationRate float64       `json:"depreciation_rate"`
	FinalPrice       float64       `json:"final_price"`
	EvaluationReport string        `json:"evaluation_report"`
	Images           string        `json:"images"`
	Status           string        `json:"status"`
	CreatedAt        time.Time     `json:"created_at"`
	Evaluator        *UserResponse `json:"evaluator,omitempty"`
}
