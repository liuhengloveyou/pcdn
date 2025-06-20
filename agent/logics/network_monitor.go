package logics

import (
	"context"
	"pcdnagent/common"
	"time"

	"github.com/shirou/gopsutil/v4/net"
	"go.uber.org/zap"
)

// NetworkStats 存储网络接口的流量统计信息
type NetworkStats struct {
	Name        string  `json:"name"`         // 网络接口名称
	Timestamp   int64   `json:"timestamp"`    // 采样时间/秒
	BytesSent   uint64  `json:"bytes_sent"`   // 发送的字节数
	BytesRecv   uint64  `json:"bytes_recv"`   // 接收的字节数
	PacketsSent uint64  `json:"packets_sent"` // 发送的数据包数
	PacketsRecv uint64  `json:"packets_recv"` // 接收的数据包数
	Errin       uint64  `json:"errin"`        // 接收错误数
	Errout      uint64  `json:"errout"`       // 发送错误数
	Dropin      uint64  `json:"dropin"`       // 接收丢包数
	Dropout     uint64  `json:"dropout"`      // 发送丢包数
	SendRate    float64 `json:"send_rate"`    // 发送速率 (bytes/s)
	RecvRate    float64 `json:"recv_rate"`    // 接收速率 (bytes/s)
}

// GetNetworkStatsWithContext 获取所有网络接口的流量统计信息（带上下文）
func GetNetworkStatsWithContext(ctx context.Context) ([]NetworkStats, error) {
	// 获取网络接口的IO计数器
	counters, err := net.IOCountersWithContext(ctx, true) // true表示获取每个网络接口的统计信息
	if err != nil {
		common.Logger.Error("Failed to get network IO counters", zap.Error(err))
		return nil, err
	}

	var stats []NetworkStats
	for _, counter := range counters {

		stats = append(stats, NetworkStats{
			Timestamp:   time.Now().Unix(),
			Name:        counter.Name,
			BytesSent:   counter.BytesSent,
			BytesRecv:   counter.BytesRecv,
			PacketsSent: counter.PacketsSent,
			PacketsRecv: counter.PacketsRecv,
			Errin:       counter.Errin,
			Errout:      counter.Errout,
			Dropin:      counter.Dropin,
			Dropout:     counter.Dropout,
			// 速率需要通过两次采样计算，初始为0
			SendRate: 0,
			RecvRate: 0,
		})
	}

	return stats, nil
}
