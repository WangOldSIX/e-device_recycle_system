package models

import (
	"e-device-recycle-backend/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.GetConfig()

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&User{},
		&Device{},
		&RecycleOrder{},
		&Evaluation{},
	)

	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	log.Println("数据库连接成功")
}
