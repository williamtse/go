package orderutil

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderNo() string {
	// 获取当前时间戳
	timestamp := time.Now().UnixNano()

	// 生成随机数
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100000) // 生成0到99999之间的随机数

	// 组合时间戳和随机数
	orderID := fmt.Sprintf("%d%d", timestamp, randomNumber)
	return orderID
}
