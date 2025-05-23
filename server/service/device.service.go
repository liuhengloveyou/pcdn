package service

import (
	"context"
	"encoding/json"
	"fmt"
	"pcdn-server/common"
	"pcdn-server/models"
	"pcdn-server/repos"

	passportprotos "github.com/liuhengloveyou/passport/protos"
	"go.uber.org/zap"
)

type deviceService struct {
}

// 新增设备
func (s *deviceService) Create(sessionUser *passportprotos.User, req *models.DeviceModel) (uint64, error) {
	if req == nil || sessionUser == nil || sessionUser.UID <= 0 {
		return 0, common.ErrParam
	}
	req.UserId = sessionUser.UID
	id, err := repos.DeviceRepo.Create(req)
	if err != nil {
		logger.Error("deviceService.Create ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	// 可选：记录业务日志
	log := &models.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: models.BUSINESS_TYPE_CREATE_DEVICE,
		Payload:      fmt.Sprintf("%v", id),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)

	return id, nil
}

// 查询设备
func (s *deviceService) Find(uid uint64, page, pageSize int) ([]models.DeviceModel, int64, error) {

	result, total, err := repos.DeviceRepo.Find(uid, page, pageSize)
	if err != nil {
		logger.Error("deviceService.Find ERR: ", zap.Error(err))
		return nil, 0, common.ErrService
	}

	// 查状态
	keys := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		keys[i] = fmt.Sprintf("agent:%s", result[i].SN)
	}

	status, err := common.RedisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		logger.Error("deviceService.Find redis ERR: ", zap.Error(err))
		return nil, 0, common.ErrService
	}

	for i := 0; i < len(status); i++ {
		if status[i] != nil {
			var agentInfo models.DeviceModel
			if err := json.Unmarshal([]byte(status[i].(string)), &agentInfo); err == nil {
				result[i].LastHeartbear = agentInfo.LastHeartbear
				result[i].Version = agentInfo.Version
				result[i].RemoteAddr = agentInfo.RemoteAddr

			}
		}
	}

	return result, total, nil
}

// 查询单个设备
func (s *deviceService) Take(id, uid uint64) (*models.DeviceModel, error) {
	if id == 0 || uid == 0 {
		return nil, common.ErrParam
	}
	m, err := repos.DeviceRepo.Get(id, uid)
	if err != nil {
		logger.Error("deviceService.Get ERR: ", zap.Error(err))
		return nil, common.ErrService
	}
	return m, nil
}

// 更新设备
func (s *deviceService) Update(sessionUser *passportprotos.User, req *models.DeviceModel) error {
	if req == nil || sessionUser == nil || sessionUser.UID <= 0 {
		return common.ErrParam
	}
	req.UserId = sessionUser.UID
	err := repos.DeviceRepo.Update(req)
	if err != nil {
		logger.Error("deviceService.Update ERR: ", zap.Error(err))
		return common.ErrService
	}
	// 可选：记录业务日志
	log := &models.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: models.BUSINESS_TYPE_UPDATE_DEVICE,
		Payload:      fmt.Sprintf("%v", req.Id),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)
	return nil
}

// 删除设备
func (s *deviceService) Delete(sessionUser *passportprotos.User, id uint64) error {
	if sessionUser == nil || sessionUser.UID <= 0 || id == 0 {
		return common.ErrParam
	}
	err := repos.DeviceRepo.Delete(id, sessionUser.UID)
	if err != nil {
		logger.Error("deviceService.Delete ERR: ", zap.Error(err))
		return common.ErrService
	}
	// 可选：记录业务日志
	log := &models.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: models.BUSINESS_TYPE_DEL_DEVICE,
		Payload:      fmt.Sprintf("%v", id),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)
	return nil
}
