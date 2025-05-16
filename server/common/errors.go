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

	ErrBotTypeErr      = errors.NewError(100, "机器人类型错误")
	ErrBotNumPeerOrder = errors.NewError(101, "单笔交易量设置错误")
	ErrBotTimeInterval = errors.NewError(102, "时间间隔设置错误")
	ErrBotOrderMode    = errors.NewError(103, "下单模式设置错误")
	ErrBotPriceMode    = errors.NewError(104, "价格模式设置错误")
	ErrBotQuantMode    = errors.NewError(105, "数量分布类型设置错误")
	ErrBotDomain       = errors.NewError(106, "机器人域名错误")
	ErrBotSymbol       = errors.NewError(107, "机器人交易对错误")
	ErrBotCost         = errors.NewError(108, "机器人Cost错误")
	ErrBotTotalTime    = errors.NewError(109, "机器人TotalTime错误")
	ErrBotOrderSize    = errors.NewError(110, "机器人OrderSize错误")
	ErrBotMail         = errors.NewError(111, "机器人Mail错误")
	ErrBotTGChatID     = errors.NewError(112, "机器人TGChat配置错误")

	ErrMartParamDomain = errors.NewError(200, "交易所域名错误")
	ErrMartParamSymbol = errors.NewError(201, "交易对错误")
	ErrMartParamKey    = errors.NewError(202, "交易所APIKEY错误")
	ErrMartParamMemo   = errors.NewError(203, "交易所MEMO错误")

	ErrOrderParam = errors.NewError(300, "下单参数错误")

	ErrRiskBotParam = errors.NewError(400, "风控机器人参数错误")

	ErrVIPPaying = errors.NewError(500, "存在正在支付的订单")
	ErrNoVIP     = errors.NewError(501, "非VIP用户")
)
