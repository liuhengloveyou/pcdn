package repos

import (
	"pcdn-server/common"
	"pcdn-server/models"

	"gorm.io/gorm"
)

var (
	DeviceRepo      = &deviceRepo{}
	BusinessLogRepo = &businessLogRepo{}
	TcRepo          *tcRepo
)

func InitRepos() {
	TcRepo = NewTcRepo(common.OrmCli)
}

func InitModels(db *gorm.DB) error {
	if err := db.AutoMigrate(models.DeviceModel{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(models.BusinessLog{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(models.TcModel{}); err != nil {
		return err
	}

	return nil
}
