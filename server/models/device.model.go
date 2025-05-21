package models

type DeviceModel struct {
	Model

	// 设备SN
	SN         string `json:"sn" gorm:"column:sn;index:idx_sn;type:VARCHAR(45);"`
	Version    string `json:"version" gorm:"-"`
	RemoteAddr string `json:"remote_addr" gorm:"-"`
	// 最后心跳时间
	LastHeartbear int64 `json:"last_heartbear" gorm:"-"`
	// 设备心跳带上来的时间
	Timestamp int64 `json:"timestamp" gorm:"-"`
}

func (DeviceModel) TableName() string {
	return "device"
}
