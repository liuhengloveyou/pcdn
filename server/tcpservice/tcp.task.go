package tcpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"pcdn-server/common"
	"time"

	"github.com/liuhengloveyou/pcdn/protos"
	redis "github.com/redis/go-redis/v9"
)

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
