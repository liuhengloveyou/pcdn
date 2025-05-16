package repos

import (
	"arbitrage/common"
	"arbitrage/protos"
	"time"
)

type botRepo struct {
}

func (p *botRepo) Set(req *protos.BotModel) (uint64, error) {
	req.UpdateTime = time.Now().UnixMilli()
	if req.Id <= 0 {
		req.CreateTime = req.UpdateTime
	}

	tx := common.OrmCli.Create(req)
	if req.Id <= 0 || tx.Error != nil {
		tx = common.OrmCli.Where("uid = ? AND bot_type = ?", req.UserId, req.BotType).Omit("bot_type", "is_running", "start_time").Updates(req)
	}

	return req.Id, tx.Error
}

func (p *botRepo) Select(id, uid int64) (*protos.BotModel, error) {
	m := &protos.BotModel{}

	tx := common.OrmCli.Where("id = ? AND uid = ?", id, uid)
	tx.Take(m)

	return m, tx.Error
}

func (p *botRepo) Find(id, uid uint64, botType int64) (rr []protos.BotModel, err error) {
	tx := common.OrmCli.Table("bots")
	tx = tx.Where("uid = ?", uid)

	if id > 0 {
		tx = tx.Where("id = ?", id)
	} else if botType > 0 {
		tx = tx.Where("bot_type = ?", botType)
	}

	err = tx.Find(&rr).Error

	return
}

func (p *botRepo) Delete(id, uid uint64) error {
	tx := common.OrmCli
	tx = tx.Where("id = ? AND uid = ?", id, uid)

	return tx.Delete(&protos.BotModel{}).Error
}

func (p *botRepo) UpdateIsRun(uid, id uint64, isRunning protos.BotRunType) (r int64, err error) {
	tx := common.OrmCli.Table("bots")
	tx = tx.Where("id =? AND uid = ?", id, uid)

	var startTime int64
	if isRunning == protos.BotIsRunning {
		startTime = time.Now().UnixMilli()
	}
	err = tx.Updates(map[string]interface{}{"is_running": isRunning, "start_time": startTime, "update_time": time.Now().UnixMilli()}).Error

	return tx.RowsAffected, err
}
