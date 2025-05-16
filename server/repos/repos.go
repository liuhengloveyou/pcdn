package repos

import (
	"arbitrage/protos"

	"gorm.io/gorm"
)

var (
	BotRepo         = &botRepo{}
	RiskBotRepo     = &riskBotRepo{}
	OrderRepo       = &orderRepo{}
	MartParamRepo   = &martParamRepo{}
	BusinessLogRepo = &businessLogRepo{}
	VipRepo         = &vipRepo{}
)

func InitModels(db *gorm.DB) error {
	if err := db.AutoMigrate(protos.BotModel{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(protos.RiskBotModel{}); err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE `risk_bots` ADD UNIQUE INDEX `idx_uid_type_mart` (`uid` ASC, `bot_type` ASC, `mart_domain` ASC, `symbol` ASC) VISIBLE;").Error; err != nil {
		// return err VISIBLE
	}

	if err := db.AutoMigrate(protos.MartParamModel{}); err != nil {
		return err
	}

	if err := db.Exec("ALTER TABLE `martparams` ADD UNIQUE INDEX `index5` (`uid` ASC, `domain` ASC) VISIBLE;").Error; err != nil {
		// return err
	}

	if err := db.AutoMigrate(protos.OrderModel{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(protos.BusinessLog{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(protos.VIPOrderStruct{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(protos.VIPMemberStruct{}); err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE `vip_members` ADD UNIQUE INDEX `index5` (`uid` ASC) VISIBLE;").Error; err != nil {
		// return err
	}

	return nil
}
