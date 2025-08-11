package models

import (
	"time"

	"gorm.io/gorm"
)

type RecycleOrder struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id" gorm:"not null"`
	DeviceID       uint           `json:"device_id" gorm:"not null"`
	OrderNo        string         `json:"order_no" gorm:"uniqueIndex;not null"` // 订单号
	ContactName    string         `json:"contact_name" gorm:"not null"`         // 联系人姓名
	ContactPhone   string         `json:"contact_phone" gorm:"not null"`        // 联系电话
	PickupAddress  string         `json:"pickup_address" gorm:"not null"`       // 上门地址
	PickupTime     *time.Time     `json:"pickup_time"`                          // 预约上门时间
	DeviceInfo     string         `json:"device_info"`                          // 设备详细信息，JSON字符串
	Images         string         `json:"images"`                               // 设备图片
	EstimatedPrice float64        `json:"estimated_price"`                      // 预估价格
	FinalPrice     *float64       `json:"final_price"`                          // 最终价格
	Status         string         `json:"status" gorm:"default:'pending'"`      // pending, confirmed, picked_up, evaluated, completed, cancelled
	Remark         string         `json:"remark"`                               // 备注
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Device     Device      `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	Evaluation *Evaluation `json:"evaluation,omitempty" gorm:"foreignKey:OrderID"`
}

type RecycleOrderCreateRequest struct {
	DeviceID      uint       `json:"device_id" binding:"required"`
	ContactName   string     `json:"contact_name" binding:"required"`
	ContactPhone  string     `json:"contact_phone" binding:"required"`
	PickupAddress string     `json:"pickup_address" binding:"required"`
	PickupTime    *time.Time `json:"pickup_time"`
	DeviceInfo    string     `json:"device_info"`
	Images        string     `json:"images"`
	Remark        string     `json:"remark"`
}

type RecycleOrderUpdateRequest struct {
	Status     string     `json:"status" binding:"oneof=pending confirmed picked_up evaluated completed cancelled"`
	FinalPrice *float64   `json:"final_price"`
	Remark     string     `json:"remark"`
	PickupTime *time.Time `json:"pickup_time"`
}

type RecycleOrderResponse struct {
	ID             uint                `json:"id"`
	UserID         uint                `json:"user_id"`
	DeviceID       uint                `json:"device_id"`
	OrderNo        string              `json:"order_no"`
	ContactName    string              `json:"contact_name"`
	ContactPhone   string              `json:"contact_phone"`
	PickupAddress  string              `json:"pickup_address"`
	PickupTime     *time.Time          `json:"pickup_time"`
	DeviceInfo     string              `json:"device_info"`
	Images         string              `json:"images"`
	EstimatedPrice float64             `json:"estimated_price"`
	FinalPrice     *float64            `json:"final_price"`
	Status         string              `json:"status"`
	Remark         string              `json:"remark"`
	CreatedAt      time.Time           `json:"created_at"`
	User           *UserResponse       `json:"user,omitempty"`
	Device         *DeviceResponse     `json:"device,omitempty"`
	Evaluation     *EvaluationResponse `json:"evaluation,omitempty"`
}
