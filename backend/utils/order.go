package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成订单号
func GenerateOrderNo() string {
	now := time.Now()
	// 格式: RC + 年月日时分秒 + 4位随机数
	// 例如: RC20240115143000001
	return fmt.Sprintf("RC%s%04d",
		now.Format("20060102150405"),
		rand.Intn(10000))
}

// 计算设备评估价格
func CalculateDevicePrice(basePrice float64, condition string, yearBought int) float64 {
	currentYear := time.Now().Year()
	yearsSinceBoought := currentYear - yearBought

	// 基础折旧率（每年10%）
	depreciationRate := 0.1 * float64(yearsSinceBoought)
	if depreciationRate > 0.8 {
		depreciationRate = 0.8 // 最大折旧80%
	}

	// 成色系数
	conditionMultiplier := map[string]float64{
		"excellent": 1.0,
		"good":      0.8,
		"fair":      0.6,
		"poor":      0.4,
	}

	multiplier, exists := conditionMultiplier[condition]
	if !exists {
		multiplier = 0.5
	}

	// 计算最终价格
	finalPrice := basePrice * (1 - depreciationRate) * multiplier

	// 确保价格不低于基础价格的10%
	minPrice := basePrice * 0.1
	if finalPrice < minPrice {
		finalPrice = minPrice
	}

	return finalPrice
}
