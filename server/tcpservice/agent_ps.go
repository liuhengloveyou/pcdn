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

func updateAgentProcessToRedis(sn string, info []*protos.ProcessInfo) error {
	// 将Agent信息序列化为JSON
	infoJSON, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化Agent信息失败: %v", err)
	}

	// 设置数据到Redis，并设置过期时间（例如30分钟）
	key := fmt.Sprintf("%sprocess/%s", AGENT_KEY_PREFIX, strings.ToUpper(sn))
	err = common.RedisClient.Set(context.Background(), key, string(infoJSON), 30*time.Minute).Err()
	if err != nil {
		common.Logger.Error("updateAgentProcessToRedis ERR", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("%s: %v", key, err)
	}

	return nil
}
