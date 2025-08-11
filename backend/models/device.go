package models

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`           // 设备名称
	Brand       string         `json:"brand"`                          // 品牌
	Model       string         `json:"model"`                          // 型号
	Category    string         `json:"category"`                       // 分类：laptop, desktop, tablet, phone
	CPU         string         `json:"cpu"`                            // 处理器
	Memory      string         `json:"memory"`                         // 内存
	Storage     string         `json:"storage"`                        // 存储
	Graphics    string         `json:"graphics"`                       // 显卡
	Screen      string         `json:"screen"`                         // 屏幕
	Condition   string         `json:"condition"`                      // 成色：excellent, good, fair, poor
	YearBought  int            `json:"year_bought"`                    // 购买年份
	BasePrice   float64        `json:"base_price"`                     // 基础回收价格
	Description string         `json:"description"`                    // 描述
	Images      string         `json:"images"`                         // 图片URLs，JSON字符串
	Status      string         `json:"status" gorm:"default:'active'"` // active, inactive
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type DeviceCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Brand       string  `json:"brand" binding:"required"`
	Model       string  `json:"model"`
	Category    string  `json:"category" binding:"required,oneof=laptop desktop tablet phone"`
	CPU         string  `json:"cpu"`
	Memory      string  `json:"memory"`
	Storage     string  `json:"storage"`
	Graphics    string  `json:"graphics"`
	Screen      string  `json:"screen"`
	Condition   string  `json:"condition" binding:"required,oneof=excellent good fair poor"`
	YearBought  int     `json:"year_bought" binding:"min=2000,max=2024"`
	BasePrice   float64 `json:"base_price" binding:"min=0"`
	Description string  `json:"description"`
	Images      string  `json:"images"`
}

type DeviceResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Category    string  `json:"category"`
	CPU         string  `json:"cpu"`
	Memory      string  `json:"memory"`
	Storage     string  `json:"storage"`
	Graphics    string  `json:"graphics"`
	Screen      string  `json:"screen"`
	Condition   string  `json:"condition"`
	YearBought  int     `json:"year_bought"`
	BasePrice   float64 `json:"base_price"`
	Description string  `json:"description"`
	Images      string  `json:"images"`
	Status      string  `json:"status"`
}
