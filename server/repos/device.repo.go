package repos

import (
	"time"

	"pcdn-server/common"
	"pcdn-server/models"
)

type deviceRepo struct {
}

// 新增设备
func (p *deviceRepo) Create(req *models.DeviceModel) (uint64, error) {
	req.CreateTime = time.Now().UnixMilli()
	req.UpdateTime = req.CreateTime
	tx := common.OrmCli.Create(req)
	return req.Id, tx.Error
}

// 查询设备（根据 id 和 uid）
func (p *deviceRepo) Get(id, uid uint64) (*models.DeviceModel, error) {
	m := &models.DeviceModel{}
	tx := common.OrmCli.Where("id = ? AND uid = ?", id, uid).Take(m)
	return m, tx.Error
}

// Find 方法用于根据条件查询设备信息
func (p *deviceRepo) Find(uid uint64, page, pageSize int) ([]models.DeviceModel, int64, error) {
	var devices []models.DeviceModel
	var total int64

	tx := common.OrmCli.Model(&models.DeviceModel{}) // 从 DeviceModel 表开始查询

	// 添加用户ID条件
	if uid > 0 {
		tx = tx.Where("uid = ?", uid)
	}

	// 获取总记录数
	tx.Count(&total)

	// 添加分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		tx = tx.Offset(offset).Limit(pageSize)
	}

	// 执行查询
	tx = tx.Find(&devices)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}

	return devices, total, nil
}

// 更新设备
func (p *deviceRepo) Update(req *models.DeviceModel) error {
	req.UpdateTime = time.Now().UnixMilli()
	tx := common.OrmCli.Model(&models.DeviceModel{}).Where("id = ? AND uid = ?", req.Id, req.UserId).Updates(req)
	return tx.Error
}

// 删除设备
func (p *deviceRepo) Delete(id, uid uint64) error {
	tx := common.OrmCli.Where("id = ? AND uid = ?", id, uid).Delete(&models.DeviceModel{})
	return tx.Error
}
