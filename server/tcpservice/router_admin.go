package tcpservice

import (
	"fmt"
	"strings"
	"time"

	"pcdn-server/common"
	"pcdn-server/models"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
)

// CreateRouterAdminTask 创建路由器管理任务
func CreateRouterAdminTask(sn, accessName string) (string, error) {
	if sn == "" || accessName == "" {
		return "", common.ErrParam
	}
	sn = strings.ToUpper(sn)

	// 创建任务
	now := time.Now().UnixMilli()
	task := &protos.Task{
		TaskId:     fmt.Sprintf("%d", now),
		TaskType:   protos.TaskType_TASK_TYPE_ROUTER_ADMIN,
		Timestamp:  now,
		Sn:         sn,
		AccessName: accessName,
	}

	// 将任务保存到Redis
	err := NewTaskToRedis(task)
	if err != nil {
		common.Logger.Error("CreateRouterAdminTask NewTaskToRedis ERR: ", zap.Error(err), zap.Any("task", task))
		return "", err
	}

	return task.TaskId, nil
}

// GetAgentStatusFromRedis 从Redis获取设备状态
func GetAgentStatusFromRedis(sn string) (*models.DeviceAgent, error) {
	return getAgentStatusFromRedis(sn)
}
