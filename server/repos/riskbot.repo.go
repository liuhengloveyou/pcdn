package repos

import (
	"arbitrage/common"
	"arbitrage/protos"
	"time"
)

type riskBotRepo struct {
}

func (p *riskBotRepo) Set(req *protos.RiskBotModel) (uint64, error) {
	req.CreateTime = time.Now().UnixMilli()
	req.UpdateTime = req.CreateTime

	tx := common.OrmCli.Create(req)
	if req.Id <= 0 || tx.Error != nil {
		tx = common.OrmCli.Where("uid = ? AND bot_type = ? AND mart_domain = ? AND symbol = ?", req.UserId, req.BotType, req.MartDomain, req.Symbol).Updates(req)
	}

	return req.Id, tx.Error
}

func (p *riskBotRepo) Select(uid int64) (*protos.RiskBotModel, error) {
	m := &protos.RiskBotModel{}

	tx := common.OrmCli.Where("uid = ?", uid)
	tx.Take(m)

	return m, tx.Error
}

func (p *riskBotRepo) Find(uid uint64, botType int64, martDomain, symbol string) (rr []protos.RiskBotModel, err error) {
	tx := common.OrmCli.Table("risk_bots")
	tx = tx.Where("uid = ?", uid)

	if botType > 0 {
		tx = tx.Where("bot_type = ?", botType)
	}

	if len(martDomain) > 0 {
		tx = tx.Where("mart_domain = ?", martDomain)
	}

	if len(symbol) > 0 {
		tx = tx.Where("symbol = ?", symbol)
	}

	err = tx.Find(&rr).Error

	return
}

func (p *riskBotRepo) UpdateIsRun(uid, id uint64, botType protos.BotType, martDomain, symbol string, isRunning protos.BotRunType) (r int64, err error) {
	tx := common.OrmCli.Table("risk_bots")
	tx = tx.Where("id =? AND uid = ? AND bot_type = ? AND mart_domain = ? AND symbol = ?", id, uid, botType, martDomain, symbol)

	var startTime int64
	if isRunning == protos.BotIsRunning {
		startTime = time.Now().UnixMilli()
	}
	err = tx.Updates(map[string]interface{}{"is_running": isRunning, "start_time": startTime, "update_time": time.Now().UnixMilli()}).Error

	return tx.RowsAffected, err
}

func (p *riskBotRepo) Update(uid, id uint64, botType protos.BotType, martDomain, symbol string) (r int64, err error) {
	tx := common.OrmCli.Table("risk_bots")
	tx = tx.Where("id =? AND uid = ?", id, uid)

	err = tx.Updates(map[string]interface{}{"update_time": time.Now().UnixMilli()}).Error

	return tx.RowsAffected, err
}
