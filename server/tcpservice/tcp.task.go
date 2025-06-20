package tcpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"pcdn-server/common"
	"strings"
	"time"

	"github.com/liuhengloveyou/pcdn/protos"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func RunTcpTasks() {
	go sendTaskToDeviceTask()
}

func sendTaskToDeviceTask() {
	if common.ServConfig.AccessName == "" {
		panic("no access name")
	}

	key := fmt.Sprintf("%s%s", common.AGENT_TASK_KEY_PREFIX, common.ServConfig.AccessName)

	for {
		// 同时更新一个队列
		rst, err := common.RedisClient.BRPop(context.Background(), time.Second*3, key).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				continue
			}
			common.Logger.Warn("SendTaskToDevice redis ERR: ", zap.Error(err))
		}
		if len(rst) < 2 {
			continue
		}
		common.Logger.Warn("SendTaskToDevice task: ", zap.Any("task", rst))
		if rst[0] != key {
			continue
		}

		var taskJson protos.Task
		if err = json.Unmarshal([]byte(rst[1]), &taskJson); err != nil {
			common.Logger.Error("sendTaskToDeviceTask json ERR ", zap.Error(err))
			continue
		}

		// 找TCP连接
		taskJson.Sn = strings.ToUpper(taskJson.Sn)
		tmpAgent, ok := AgentMap[taskJson.Sn]
		common.Logger.Debug("sendTaskToDeviceTask ", zap.Any("task", AgentMap))
		if !ok {
			common.Logger.Error("sendTaskToDeviceTask AgentMap ERR ", zap.Any("task", taskJson.String()))

			// TODO 把错误信息返回给前端

			continue
		}

		// 发送任务到设备
		if err = SendTaskToDevice(tmpAgent, &taskJson); err != nil {
			common.Logger.Error("sendTaskToDeviceTask tcp ERR ", zap.Error(err))
		}

		common.Logger.Info("sendTaskToDeviceTask OK: ", zap.Any("task", taskJson.String()))
	}
}

// TODO
// 在适当的地方（如main函数或初始化函数中）添加定时任务
func startAgentCleanupTask() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			cleanupOfflineAgents()
		}
	}()
}

// 清理离线Agent
func cleanupOfflineAgents() {
	ctx := context.Background()

	// 获取所有在线Agent
	agentSNs, err := common.RedisClient.SMembers(ctx, "agents:online").Result()
	if err != nil {
		common.Logger.Sugar().Errorf("获取在线Agent列表失败: %v", err)
		return
	}

	now := time.Now().UnixMilli()
	for _, sn := range agentSNs {
		// 获取Agent详情
		agentJSON, err := common.RedisClient.Get(ctx, fmt.Sprintf("agent:%s", sn)).Result()
		if err != nil {
			if err == redis.Nil {
				// Agent详情已过期，从在线列表中删除
				common.RedisClient.SRem(ctx, "agents:online", sn)
			}
			continue
		}

		var agent protos.DeviceAgent
		if err := json.Unmarshal([]byte(agentJSON), &agent); err != nil {
			common.Logger.Sugar().Errorf("解析Agent信息失败: %v", err)
			continue
		}

		// 检查最后心跳时间，如果超过5分钟没有心跳，则认为离线
		if now-agent.LastHeartbear > 5*60*1000 {
			common.RedisClient.SRem(ctx, "agents:online", sn)
			common.Logger.Sugar().Infof("Agent %s 已超时离线，从在线列表中删除", sn)
		}
	}
}

// 管理后台往设备发任务，都先放到redis
func NewTaskToRedis(task *protos.Task) error {
	if task.AccessName == "" {
		return common.ErrAgentNoAccess
	}
	if task.GetTaskId() == "" {
		return common.ErrAgentNoAccess
	}

	// 将Agent信息序列化为JSON
	taskJson, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("SendTaskToDevice json ERR: %v", err)
	}

	// 设置数据到Redis，并设置过期时间（例如30分钟）
	if err = common.RedisClient.Set(context.Background(), fmt.Sprintf("task/%s", task.GetTaskId()), string(taskJson), 30*time.Minute).Err(); err != nil {
		return fmt.Errorf("SendTaskToDevice redis ERR: %v", err)
	}

	// 同时更新一个队列. 接入点会从这里拉任务下发
	err = common.RedisClient.LPush(context.Background(), fmt.Sprintf("%s%s", common.AGENT_TASK_KEY_PREFIX, task.AccessName), taskJson).Err()
	if err != nil {
		common.Logger.Sugar().Warnf("SendTaskToDevice redis ERR: %v", err)
	}

	return nil
}
