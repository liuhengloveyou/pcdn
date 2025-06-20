package repos

import (
	"pcdn-server/common"
	"pcdn-server/models"
)

type businessLogRepo struct {
}

func (p *businessLogRepo) Add(m *models.BusinessLog) (uint64, error) {
	tx := common.OrmCli.Create(m)

	return m.Id, tx.Error
}

func (p *businessLogRepo) Find(uid uint64, martDomain string, page int) (rr []models.BusinessLog, err error) {
	tx := common.OrmCli.Table("business_log")
	tx = tx.Where("uid = ?", uid)

	if len(martDomain) > 0 {
		tx = tx.Where("mart_domain = ?", martDomain)
	}
	// if len(symbol) > 0 {
	// 	tx = tx.Where("symbol = ?", symbol)
	// }

	err = tx.Order("id desc").Offset((page - 1) * 20).Limit(20).Find(&rr).Error

	return
}
