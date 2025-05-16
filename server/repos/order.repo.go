package repos

import (
	"arbitrage/common"
	"arbitrage/protos"
	"time"
)

type orderRepo struct {
}

func (p *orderRepo) New(req *protos.OrderModel) (uint64, error) {
	req.CreateTime = time.Now().UnixMilli()
	req.UpdateTime = req.CreateTime

	tx := common.OrmCli.Create(req)

	return req.Id, tx.Error
}

func (p *orderRepo) Select(uid int64) (*protos.OrderModel, error) {
	m := &protos.OrderModel{}

	tx := common.OrmCli.Where("uid = ?", uid)
	tx.Take(m)

	return m, tx.Error
}

func (p *orderRepo) Find(uid uint64, martDomain, symbol string, page int) (rr []protos.OrderModel, err error) {
	tx := common.OrmCli.Table("orders")
	if uid > 0 {
		tx = tx.Where("uid = ?", uid)
	}

	if len(martDomain) > 0 {
		tx = tx.Where("mart_domain = ?", martDomain)
	}

	if len(symbol) > 0 {
		tx = tx.Where("symbol = ?", symbol)
	}

	err = tx.Order("id desc").Offset((page - 1) * 20).Limit(20).Find(&rr).Error

	return
}
