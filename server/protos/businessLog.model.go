package protos

type BusinessType string

const (
	BUSINESS_TYPE_SET_BOT   BusinessType = "Set Bot"
	BUSINESS_TYPE_START_BOT BusinessType = "Start Bot"
	BUSINESS_TYPE_STOP_BOT  BusinessType = "Stop Bot"
	BUSINESS_TYPE_DEL_BOT   BusinessType = "Delete Bot"
	BUSINESS_TYPE_ERR       BusinessType = "ERROR"
)

type BusinessLog struct {
	Model

	UserName     string       `json:"userName" validate:"-" gorm:"column:user_name;type:varchar(128);"`
	BusinessType BusinessType `json:"businessType"  gorm:"column:business_type;type:varchar(100);not null;"`
	Payload      string       `json:"payload"  gorm:"column:payload;"`
}

func (BusinessLog) TableName() string {
	return "business_log"
}
