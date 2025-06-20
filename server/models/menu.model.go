package models

import (
	"database/sql/driver"
	"encoding/json"
)

type MenuItem struct {
	Title    string        `json:"title"`
	Icon     string        `json:"icon"`
	Link     string        `json:"link"`
	Children MenuItemArray `json:"children"`
}

type MenuItemArray []MenuItem

func (m *MenuItemArray) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return json.Unmarshal(b, m)
}
func (m MenuItemArray) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type QuicMenu struct {
	Model
	Payload MenuItemArray `json:"payload" gorm:"column:payload;type:JSON;"`
}

func (QuicMenu) TableName() string {
	return "quic_menu"
}
