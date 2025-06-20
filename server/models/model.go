package models

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"

	passport "github.com/liuhengloveyou/passport/protos"
)

type Model struct {
	Id         uint64 `json:"id,omitempty" validate:"omitempty,min=1" db:"id" gorm:"column:id;type:INT;primaryKey;autoIncrement;"`
	UserId     uint64 `json:"uid,omitempty" validate:"omitempty,min=1" db:"uid" gorm:"column:uid;type:INT;not null;index:idx_user_id;"`
	TenantId   uint64 `json:"tenantId,omitempty" validate:"omitempty,min=1" db:"tenant_id" gorm:"column:tenant_id;type:INT;not null;index:idx_tenant_id;"`
	CreateTime int64  `json:"createTime,omitempty" validate:"-" db:"create_time" gorm:"column:create_time;not null;"` // 创建时间 timestamptz
	UpdateTime int64  `json:"updateTime,omitempty" validate:"-" db:"update_time" gorm:"column:update_time;not null;"` // 最后更新时间 timestamptz
}

type PageResponse struct {
	Total int64       `json:"total,omitempty"`
	List  interface{} `json:"list"`
}

type UserLiteArr []passport.User

func (m *UserLiteArr) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m UserLiteArr) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

type Int64Arr []int64

func (m *Int64Arr) Scan(src interface{}) error {
	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m Int64Arr) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

type StringArr []string

func (m *StringArr) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, m)
}
func (m StringArr) Value() (driver.Value, error) {
	return sonic.Marshal(m)
}

type MapStruct map[string]interface{}

func (t *MapStruct) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	if len(src.([]byte)) <= 2 {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, t)
}
func (t MapStruct) Value() (driver.Value, error) {
	return sonic.Marshal(t)
}

type MapStringInt map[string]int

func (t *MapStringInt) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	if len(src.([]byte)) <= 2 {
		return nil
	}

	b, _ := src.([]byte)
	return sonic.Unmarshal(b, t)
}
func (t MapStringInt) Value() (driver.Value, error) {
	return sonic.Marshal(t)
}

type HeartbeatReq struct {
	Name   string `json:"name"`
	Time   string `json:"time"`
	Stat   string `json:"stat"`
	Remark string `json:"remark"`
}

type TrifficLimitReq struct {
	SN          string `json:"sn" validate:"required"`
	IfaceName   string `json:"ifaceName"`
	UploadLimit uint   `json:"uploadLimit" validate:"required"` // mbps
}
