package tcpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"pcdn-server/common"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
)

func updateAgentMonitorToRedis(heartbeat *protos.Heartbeat) error {
	if heartbeat.Sn == "" {
		return fmt.Errorf("heartbeat.Sn is empty")
	}
	if heartbeat.Monitor == nil {
		return fmt.Errorf("heartbeat.MonitorInfo is empty")
	}

	monitorJson, err := json.Marshal(heartbeat.Monitor)
	if err != nil {
		return fmt.Errorf("序列化Agent信息失败: %v", err)
	}

	// 设置数据到Redis，并设置过期时间（例如30分钟）
	key := fmt.Sprintf("%s%s", common.AGENT_MONITOR_KEY_PREFIX, strings.ToUpper(heartbeat.Sn))
	err = common.RedisClient.Set(context.Background(), key, monitorJson, 30*time.Minute).Err()
	if err != nil {
		common.Logger.Error("updateAgentMonitorToRedis ERR", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("%s: %v", key, err)
	}

	return nil
}
