package service

import (
	"pcdn-server/common"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	TcService          = &tcService{}
	DeviceService      = &deviceService{}
	BusinessLogService = &businessLogService{}
)

func init() {
	logger = common.Logger
}
