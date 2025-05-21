package service

import (
	"pcdn-server/common"
	"pcdn-server/protos"
	"strings"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	BotService         = &botService{}
	RiskBotService     = &riskBotService{}
	MartParamService   = &martParamService{}
	OrderService       = &orderService{}
	MartService        = &martService{}
	BusinessLogService = &businessLogService{}
	VipOrderService    = &vipOrderService{}
)

func init() {
	logger = common.Logger
}

// 判断是否是合法交易所域名
func IsMartDomainLegal(martDomain string) bool {
	for i := 0; i < len(protos.MartsInSystem); i++ {
		if strings.Compare(martDomain, protos.MartsInSystem[i]) == 0 {
			return true
		}
	}

	return false
}
