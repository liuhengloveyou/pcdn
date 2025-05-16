package common

import (
	"math"
	"time"
)

// XYs implements the XYer interface.
type XYs []XY

// XY is an x and y value.
type XY struct {
	X int64
	Y float64
}

/*
points: 点数
startY: 起始值
endY: 结束值
peaks: 波峰数量
amplitude: 振幅
*/
func MakeSinPoints(minutes int, startY, endY float64, peaks int, amplitude float64) XYs {
	now := time.Now().Unix()

	// 每分钟一个点
	points := make(XYs, minutes)

	// 计算斜率
	slope := (endY - startY) / float64(minutes) // 使用时间来计算斜率

	// 计算波形频率系数，直接控制波峰数量
	waveFreq := float64(peaks) * math.Pi / float64(minutes/2) // 调整系数使peaks为波峰数量

	for i := 0; i < minutes; i++ {
		x := float64(i)           // x轴表示分钟数
		baseY := startY + slope*x // 基准线

		// 使用振幅直接控制波形高度
		points[i].X = now + int64(i*60)
		// amplitude直接表示相对于baseY的振幅百分比
		points[i].Y = baseY * (1 + amplitude*math.Sin(waveFreq*x))
	}

	return points
}
