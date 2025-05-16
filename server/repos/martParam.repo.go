package repos

import (
	"arbitrage/common"
	"arbitrage/protos"
	"time"
)

type martParamRepo struct {
}

func (p *martParamRepo) Set(req *protos.MartParamModel) (uint64, error) {
	req.CreateTime = time.Now().UnixMilli()
	req.UpdateTime = req.CreateTime

	tx := common.OrmCli.Create(req)
	// if req.Id <= 0 || tx.Error != nil {
	// 	tx = common.OrmCli.Where("uid = ? AND domain = ?", req.UserId, req.MartDomain).Updates(req)
	// }

	return req.Id, tx.Error
}

func (p *martParamRepo) Select(uid, id uint64, martDomain string) (*protos.MartParamModel, error) {
	if uid <= 0 {
		return nil, common.ErrParam
	}

	m := &protos.MartParamModel{}

	tx := common.OrmCli.Where("uid = ?", uid)
	if id > 0 {
		tx = tx.Where("id = ?", id)
	}

	if len(martDomain) > 0 {
		tx = tx.Where("domain = ?", martDomain)
	}

	tx.Take(m)

	return m, tx.Error
}

func (p *martParamRepo) Find(uid uint64, martDomain string) (rr []protos.MartParamModel, err error) {
	tx := common.OrmCli.Table("martparams")
	tx = tx.Where("uid = ?", uid)

	if len(martDomain) > 0 {
		tx = tx.Where("domain = ?", martDomain)
	}

	err = tx.Find(&rr).Error

	return
}

func (p *martParamRepo) UpdateActive(id, uid uint64) (r int64, err error) {
	tx := common.OrmCli.Table("martparams")
	tx = tx.Where("id =? AND uid = ?", id, uid)

	err = tx.Updates(map[string]interface{}{"active": time.Now().UnixMilli(), "update_time": time.Now().UnixMilli()}).Error

	return tx.RowsAffected, err
}

func (p *martParamRepo) Delete(id, uid uint64) (r int64, err error) {
	tx := common.OrmCli.Table("martparams")
	tx = tx.Where("id =? AND uid = ?", id, uid)

	err = tx.Delete(&protos.MartParamModel{}).Error

	return tx.RowsAffected, err
}
