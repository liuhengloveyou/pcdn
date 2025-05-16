package common

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCode 日期20191025时间戳1571987125435+3位随机数
func GenerateCode() string {
	code := fmt.Sprintf("%s%d%03d", time.Now().Format("0102"), time.Now().UnixMilli(), rand.Intn(1000))
	return code
}
