package service

import (
	"context"
	"encoding/json"
	"fmt"
	"pcdn-server/common"
	"pcdn-server/models"
	"pcdn-server/repos"
	"pcdn-server/tcpservice"
	"strings"
	"time"

	passportprotos "github.com/liuhengloveyou/passport/protos"
	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type deviceService struct {
}

// 新增设备
func (s *deviceService) Create(sessionUser *passportprotos.User, req *models.DeviceModel) (uint64, error) {
	if req == nil || sessionUser == nil || sessionUser.UID <= 0 {
		return 0, common.ErrParam
	}
	req.UserId = sessionUser.UID
	req.SN = strings.ToUpper(req.SN)
	req.CreateTime = time.Now().UnixMilli()
	common.Logger.Debug("deviceService.Create", zap.Any("sess", sessionUser), zap.Any("req", req))

	id, err := repos.DeviceRepo.Create(req)
	if err != nil {
		logger.Error("deviceService.Create ERR: ", zap.Error(err))
		// 检查是否是唯一键约束错误
		if strings.Contains(err.Error(), "duplicate key") {
			return 0, common.ErrAgentSNExists
		}
		return 0, common.ErrService
	}

	// 可选：记录业务日志
	log := &models.BusinessLog{
		UserName:     sessionUser.Cellphone.String,
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
		keys[i] = fmt.Sprintf("agent/%s", strings.ToUpper(result[i].SN))
	}

	if len(keys) > 0 {
		status, err := common.RedisClient.MGet(context.Background(), keys...).Result()
		// logger.Debug("deviceService.Find redis: ", zap.Any("status", status), zap.Error(err))
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

// 获取设备监控信息
func (s *deviceService) GetMonitorInfo(sn string) (*protos.SystemMonitorData, error) {
	if sn == "" {
		return nil, common.ErrParam
	}

	// TODO 当前用户有没有权限查看？

	// 从Redis获取监控信息
	key := fmt.Sprintf("%s%s", common.AGENT_MONITOR_KEY_PREFIX, strings.ToUpper(sn))
	monitorData, err := common.RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error("deviceService.GetMonitorInfo redis ERR: ", zap.String("key", key), zap.Error(err))
		return nil, common.ErrService
	}

	// 解析监控信息
	var monitorJson protos.SystemMonitorData
	if err := json.Unmarshal([]byte(monitorData), &monitorJson); err != nil {
		logger.Error("deviceService.GetMonitorInfo Unmarshal ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return &monitorJson, nil
}

// GetRouterAdminURL 获取路由器管理界面URL
func (s *deviceService) GetRouterAdminURL(sn string) (string, error) {
	if sn == "" {
		return "", common.ErrParam
	}
	sn = strings.ToUpper(sn)

	// 获取设备状态，确保设备在线
	agentStat, err := tcpservice.GetAgentStatusFromRedis(sn)
	if err != nil {
		logger.Error("GetRouterAdminURL GetAgentStatusFromRedis ERR: ", zap.Error(err))
		return "", err
	}

	if agentStat.AccessName == "" {
		return "", common.ErrAgentNoAccess
	}

	// 创建路由器管理任务
	taskId, err := tcpservice.CreateRouterAdminTask(sn, agentStat.AccessName)
	if err != nil {
		logger.Error("GetRouterAdminURL CreateRouterAdminTask ERR: ", zap.Error(err))
		return "", err
	}

	// 等待任务响应
	redisKey := fmt.Sprintf("%s%s", common.TASK_RESPONSE_KEY_PREFIX, taskId)
	rstByte, err := common.RedisClient.BRPop(context.Background(), time.Second*10, redisKey).Result()
	if err != nil {
		logger.Error("GetRouterAdminURL redis ERR: ", zap.String("key", redisKey), zap.Error(err))
		return "", common.ErrService
	}

	// 解析任务响应
	var task protos.Task
	if err := proto.Unmarshal([]byte(rstByte[1]), &task); err != nil {
		common.Logger.Sugar().Errorf("GetRouterAdminURL msg ERR: %v", err)
		return "", err
	}

	// 检查任务是否成功
	if task.ErrMsg != "" {
		common.Logger.Sugar().Errorf("GetRouterAdminURL task error: %v", task.ErrMsg)
		return "", fmt.Errorf(task.ErrMsg)
	}

	// 返回路由器管理URL
	if task.Url == nil || *task.Url == "" {
		return "", fmt.Errorf("获取路由器管理URL失败")
	}

	return *task.Url, nil
}
