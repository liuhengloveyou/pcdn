package models

type TcModel struct {
	Model

	// 设备SN，每个设备记录最新状态
	SN string `json:"sn" gorm:"column:sn;uniqueIndex:idx_sn;type:VARCHAR(45);"`

	// 任务ID
	TaskID string `json:"taskId" gorm:"column:task_id;uniqueIndex:idx_task_id;type:VARCHAR(45);"`

	// 上行限速 mbps mbit
	UpLimit uint `json:"upLimit" gorm:"column:up_limit;type:int;default:0;"`
	// 下行限速 mbps mbit
	DownLimit uint `json:"downLimit" gorm:"column:down_limit;type:int;default:0;"`

	// 任务状态
	Status int `json:"status" gorm:"column:status;type:int;default:0;"`
}

func (TcModel) TableName() string {
	return "tc"
}
