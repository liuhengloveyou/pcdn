package service

import (
	"context"
	"fmt"
	"pcdn-server/common"
	"pcdn-server/models"
	"pcdn-server/repos"
	"pcdn-server/tcpservice"
	"strings"
	"time"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type tcService struct {
}

// 设置设备的限速规则，每个设备有一条记录
func (s *tcService) TrifficLimit(ctx context.Context, sn, iFaceName string, uploadLimit uint) (val, detail string, err error) {
	if sn == "" {
		return "", "", common.ErrParam
	}
	sn = strings.ToUpper(sn)

	taskId, err := tcpservice.TrifficLimit(sn, iFaceName, uploadLimit)
	if err != nil {
		logger.Error("TrifficLimit ERR: ", zap.Error(err))
		return "", "", err
	}

	// 等待任务响应
	redisKey := fmt.Sprintf("%s%s", common.TASK_RESPONSE_KEY_PREFIX, taskId)
	rstByte, err := common.RedisClient.BRPop(context.Background(), time.Second*10, redisKey).Result()
	if err != nil {
		logger.Error("TrifficLimit redis ERR: ", zap.String("key", redisKey), zap.Error(err))
		return "", "", common.ErrService
	}

	// 下发成功以后，保存到数据库
	m := &models.TcModel{
		TaskID:  taskId,
		SN:      sn,
		UpLimit: uploadLimit,
	}
	m.UserId = ctx.Value("UID").(uint64)
	m.TenantId = ctx.Value("TID").(uint64)
	if err = repos.TcRepo.Save(m); err != nil {
		logger.Error("TrifficLimit DB ERR: ", zap.Error(err))
		return "", "", err
	}

	var task protos.Task
	if err = proto.Unmarshal([]byte(rstByte[1]), &task); err != nil {
		common.Logger.Sugar().Errorf("TrifficLimit msg ERR: ", err)
		return "", "", err
	}
	common.Logger.Sugar().Infof("TrifficLimit: %v %v\n", task.TaskId, task.ErrMsg)

	rate := ""
	if task.Rate != nil {
		rate = *task.Rate
	}

	// 记录业务日志
	businessLog := &models.BusinessLog{
		BusinessType: models.BUSINESS_TYPE_CREATE_TC,
		Payload:      fmt.Sprintf("%s | %s | %v | %s", sn, iFaceName, uploadLimit, rate),
	}
	businessLog.UserId = ctx.Value("UID").(uint64)
	businessLog.TenantId = ctx.Value("TID").(uint64)
	businessLog.UserName = ctx.Value("Nickname").(string)
	_, err = BusinessLogService.Add(businessLog)
	common.Logger.Sugar().Debug("TrifficLimit Add BusinessLog: %v %v\n", businessLog.Id, err)
	if err != nil {
		common.Logger.Sugar().Errorf("TrifficLimit Add BusinessLog ERR: ", err)
	}

	return rate, task.ErrMsg, nil
}

func (s *tcService) TrifficLimitStat(sn, iFaceName string) (val, detail string, err error) {
	if sn == "" || iFaceName == "" {
		return "", "", common.ErrParam
	}
	sn = strings.ToUpper(sn)

	taskId, err := tcpservice.TrifficLimitStat(sn, iFaceName)
	if err != nil {
		logger.Error("TrifficLimitStat ERR: ", zap.Error(err))
		return "", "", err
	}

	redisKey := fmt.Sprintf("%s%s", common.TASK_RESPONSE_KEY_PREFIX, taskId)
	rstByte, err := common.RedisClient.BRPop(context.Background(), time.Second*10, redisKey).Result()
	if err != nil {
		logger.Error("TrifficLimitStat redis ERR: ", zap.String("key", redisKey), zap.Error(err))
		return "", "", common.ErrService
	}

	var task protos.Task
	if err := proto.Unmarshal([]byte(rstByte[1]), &task); err != nil {
		common.Logger.Sugar().Errorf("TrifficLimitStat msg ERR: ", err)
		return "", "", err
	}
	common.Logger.Sugar().Infof("TrifficLimitStat: %v %v\n", task.TaskId, task.ErrMsg)

	rate := ""
	if task.Rate != nil {
		rate = *task.Rate
	}
	return rate, task.ErrMsg, nil
}

// 定时同步数据库中的限速规则
func (s *tcService) SyncAllTrifficLimitToDevice() {
	for page := 1; ; page++ {
		tcList, _, err := repos.TcRepo.List(page, 100)
		if err != nil {
			common.Logger.Sugar().Errorf("SyncAllTrifficLimitToDevice FindAll ERR: ", err)
			return
		}
		if len(tcList) == 0 {
			break
		}

		// 下发限速任务
		for _, tcConf := range tcList {
			if tcConf.UpLimit == 0 {
				// 不限速
				continue
			}
			// 限速
			taskId, err := tcpservice.TrifficLimit(tcConf.SN, "", tcConf.UpLimit)
			if err != nil {
				logger.Error("SyncAllTrifficLimitToDevice ERR: ", zap.Error(err))
				continue
			}

			// 等待任务响应
			redisKey := fmt.Sprintf("%s%s", common.TASK_RESPONSE_KEY_PREFIX, taskId)
			rstByte, err := common.RedisClient.BRPop(context.Background(), time.Second*10, redisKey).Result()
			if err != nil {
				logger.Error("SyncAllTrifficLimitToDevice redis ERR: ", zap.String("key", redisKey), zap.Error(err))
				continue
			}

			// 记录业务日志
			businessLog := &models.BusinessLog{
				BusinessType: models.BUSINESS_TYPE_CREATE_TC,
				Payload:      fmt.Sprintf("%v | %v | %s", tcConf.SN, tcConf.UpLimit, string(rstByte[1])),
			}
			businessLog.UserId = 0
			businessLog.TenantId = 0
			businessLog.UserName = "system"
			if _, err := BusinessLogService.Add(businessLog); err != nil {
				common.Logger.Sugar().Errorf("SyncAllTrifficLimitToDevice Add BusinessLog ERR: ", err)
			}
		}
	}
}
