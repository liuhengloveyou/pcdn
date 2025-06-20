package models

import (
	"net"
	"sync"
)

type DeviceAgent struct {
	// 设备SN
	SN string `json:"sn" gorm:"column:sn;uniqueIndex:idx_sn;type:VARCHAR(45);"`
	// agent 版本
	Version string `json:"version" gorm:"-"`
	// 设备IP
	RemoteAddr string `json:"remoteAddr" gorm:"-"`
	// 最后心跳时间
	LastHeartbear int64 `json:"lastHeartbear" gorm:"-"`
	// 设备心跳带上来的时间
	Timestamp int64 `json:"timestamp" gorm:"-"`
	// 接入点名
	AccessName string `json:"accessName" gorm:"-"`

	// 设备tcp长连接
	ClientTcpConn net.Conn `json:"-" gorm:"-"`
	// 同时只能有一个任务在执行
	MU sync.Mutex
}

type DeviceModel struct {
	Model

	// 设备SN
	SN string `json:"sn" gorm:"column:sn;uniqueIndex:idx_sn;type:VARCHAR(45);"`
	// agent 版本
	Version string `json:"version" gorm:"-"`
	// 设备IP
	RemoteAddr string `json:"remoteAddr" gorm:"-"`
	// 最后心跳时间
	LastHeartbear int64 `json:"lastHeartbear" gorm:"-"`
	// 设备心跳带上来的时间
	Timestamp int64 `json:"timestamp" gorm:"-"`
	// 接入点名
	AccessName string `json:"accessName" gorm:"-"`
}

func (DeviceModel) TableName() string {
	return "device"
}
