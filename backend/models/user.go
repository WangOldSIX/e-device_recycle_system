package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // 密码不返回给前端
	Phone     string         `json:"phone" gorm:"unique"`
	Email     string         `json:"email" gorm:"unique"`
	RealName  string         `json:"real_name"`
	Avatar    string         `json:"avatar"`
	Role      string         `json:"role" gorm:"default:'user'"`     // user, admin
	Status    string         `json:"status" gorm:"default:'active'"` // active, banned
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	RecycleOrders []RecycleOrder `json:"recycle_orders,omitempty" gorm:"foreignKey:UserID"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"email"`
	RealName string `json:"real_name"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}
