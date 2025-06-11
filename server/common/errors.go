package common

import (
	"github.com/liuhengloveyou/go-errors"
)

var (
	ErrOK      = &errors.Error{Code: 0, Message: "OK"}
	ErrParam   = &errors.Error{Code: 1, Message: "参数错误"}
	ErrService = &errors.Error{Code: 2, Message: "操作失败"}
	ErrNoAuth  = &errors.Error{Code: 3, Message: "权限错误"}
	ErrSession = &errors.Error{Code: 4, Message: "Session error."}

	ErrAgentNoAccess = errors.NewError(-10000, "AGENT接入点错误")
	ErrAgentNoId     = errors.NewError(-10001, "任务没有ID")
	ErrAgentSNExists = errors.NewError(-10002, "设备SN已存在")
)
