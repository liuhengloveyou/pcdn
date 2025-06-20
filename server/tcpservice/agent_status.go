package tcpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"pcdn-server/common"
	"pcdn-server/models"
)

func getAgentStatusFromRedis(sn string) (agent *models.DeviceAgent, err error) {
	ctx := context.Background()
	rstStr, err := common.RedisClient.Get(ctx, snToKey(sn)).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rstStr, &agent)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func updateAgentStatusToRedis(agent *models.DeviceAgent) error {
	agent.AccessName = common.ServConfig.AccessName

	// 将Agent信息序列化为JSON
	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("序列化Agent信息失败: %v", err)
	}

	// 使用Redis客户端将数据保存到Redis
	// 使用agent的SN作为键
	key := snToKey(agent.SN)

	// 设置数据到Redis，并设置过期时间（例如30分钟）
	ctx := context.Background()
	err = common.RedisClient.Set(ctx, key, string(agentJSON), 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("保存Agent信息到Redis失败: %v", err)
	}

	// 同时更新一个集合，用于列出所有在线的Agent
	// TODO 过期清理
	err = common.RedisClient.SAdd(ctx, "agents:online", agent.SN).Err()
	if err != nil {
		common.Logger.Sugar().Warnf("更新在线Agent集合失败: %v", err)
	}

	return nil
}

func snToKey(sn string) string {
	return fmt.Sprintf("%s%s", common.AGENT_KEY_PREFIX, strings.ToUpper(sn))
}
