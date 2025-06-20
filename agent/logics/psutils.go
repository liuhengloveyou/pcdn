package logics

import (
	"context"
	"fmt"
	"pcdnagent/common"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"

	"github.com/shirou/gopsutil/v4/docker"
	"github.com/shirou/gopsutil/v4/process"
)

func FillProcessInfo(heartbeat *protos.Heartbeat) {
	processes, _ := process.ProcessesWithContext(context.Background())
	for _, p := range processes {
		name, _ := p.Name()
		exe, _ := p.Exe()
		CPUPercent, _ := p.CPUPercent()

		heartbeat.Monitor.Processes = append(heartbeat.Monitor.Processes, &protos.SystemMonitorProcess{
			Pid:  p.Pid,
			Name: name,
			Exe:  exe,
			Cpu:  float32(CPUPercent),
		})
	}
}

func FillDockerInfo(heartbeat *protos.Heartbeat) {
	dockerIDList, _ := docker.GetDockerIDList()
	for _, dockerID := range dockerIDList {
		fmt.Println(">>>", dockerID)
	}
}

var prevNewworkStats []NetworkStats

// FillNetworkInfo 填充网络流量信息到心跳消息中
func FillNetworkInfo(heartbeat *protos.Heartbeat) {
	var err error

	if len(prevNewworkStats) == 0 {
		// 获取所有网络接口的初始统计信息
		prevNewworkStats, err = GetNetworkStatsWithContext(context.Background())
		if err != nil {
			common.Logger.Error("GetNetworkStatsWithContext ERR", zap.Error(err))
			return
		}
		return
	}

	// 获取当前统计信息
	currentStats, err := GetNetworkStatsWithContext(context.Background())
	if err != nil {
		common.Logger.Error("GetNetworkStatsWithContext ERR", zap.Error(err))
		return
	}

	// 计算并打印每个接口的速率
	for _, curr := range currentStats {
		for _, prev := range prevNewworkStats {
			if curr.Name == prev.Name {
				// 计算上传和下载速率 (bytes/s)
				uploadRate := float64(0)
				if curr.BytesSent >= prev.BytesSent {
					uploadRate = float64(curr.BytesSent-prev.BytesSent) / float64(curr.Timestamp-prev.Timestamp)
				}

				downloadRate := float64(0)
				if curr.BytesRecv >= prev.BytesRecv {
					downloadRate = float64(curr.BytesRecv-prev.BytesRecv) / float64(curr.Timestamp-prev.Timestamp)
				}

				heartbeat.Monitor.Network = append(heartbeat.Monitor.Network, &protos.SystemMonitorNetwork{
					Name:      curr.Name,
					Timestamp: curr.Timestamp,
					SendRate:  uploadRate,
					RecvRate:  downloadRate,
				})

				// // 转换为更友好的单位
				// uploadRateStr := formatBytesPerSecond(uploadRate)
				// downloadRateStr := formatBytesPerSecond(downloadRate)
				// totalUploadStr := formatBytes(curr.BytesSent)
				// totalDownloadStr := formatBytes(curr.BytesRecv)
				// // 打印结果
				// fmt.Printf("%-10s %-15s %-15s %-15s %-15s\n",
				// 	curr.Name, uploadRateStr, downloadRateStr, totalUploadStr, totalDownloadStr)

				break
			}
		}
	}

	// 更新上一次的统计信息
	prevNewworkStats = currentStats
}

// formatBytes 将字节数格式化为人类可读的形式
func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	var (
		unit  string
		value float64
	)

	switch {
	case bytes >= TB:
		unit = "TB"
		value = float64(bytes) / TB
	case bytes >= GB:
		unit = "GB"
		value = float64(bytes) / GB
	case bytes >= MB:
		unit = "MB"
		value = float64(bytes) / MB
	case bytes >= KB:
		unit = "KB"
		value = float64(bytes) / KB
	default:
		unit = "B"
		value = float64(bytes)
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

// formatBytesPerSecond 将字节/秒格式化为人类可读的形式
func formatBytesPerSecond(bytesPerSecond float64) string {
	return formatBytes(uint64(bytesPerSecond)) + "/s"
}
