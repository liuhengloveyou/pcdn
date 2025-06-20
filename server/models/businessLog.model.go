package models

type BusinessType string

const (
	BUSINESS_TYPE_CREATE_DEVICE BusinessType = "CREATE_DEVICE"
	BUSINESS_TYPE_UPDATE_DEVICE BusinessType = "UPDATE_DEVICE"
	BUSINESS_TYPE_DEL_DEVICE    BusinessType = "DELETE_DEVICE"

	BUSINESS_TYPE_CREATE_TC BusinessType = "TC"

	BUSINESS_TYPE_ERR BusinessType = "ERROR"
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
