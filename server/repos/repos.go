package repos

import (
	"pcdn-server/models"

	"gorm.io/gorm"
)

var (
	DeviceRepo      = &deviceRepo{}
	BusinessLogRepo = &businessLogRepo{}
)

func InitModels(db *gorm.DB) error {
	if err := db.AutoMigrate(models.DeviceModel{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(models.BusinessLog{}); err != nil {
		return err
	}

	return nil
}
