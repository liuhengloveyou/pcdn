package service

import (
	"arbitrage/common"
	"arbitrage/protos"
	"arbitrage/repos"
	"time"

	passportprotos "github.com/liuhengloveyou/passport/protos"

	"go.uber.org/zap"
)

type riskBotService struct {
}

func (s *riskBotService) Set(req *protos.RiskBotModel) (uint64, error) {
	if req == nil {
		return 0, common.ErrParam
	}

	err := s.validate(req)
	if err != nil {
		logger.Error("riskBotService.Set validate  ERR: ", zap.Error(err))
		return 0, err
	}

	id, err := repos.RiskBotRepo.Set(req)
	if err != nil {
		logger.Error("riskBotService.Set ERR: ", zap.Error(err))
		return id, common.ErrService
	}

	return id, nil
}

func (s *riskBotService) Select(uid int64) (m *protos.RiskBotModel, err error) {
	if uid <= 0 {
		return nil, common.ErrParam
	}

	m, err = repos.RiskBotRepo.Select(uid)
	if err != nil {
		logger.Error("select ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return m, nil
}

func (s *riskBotService) Find(uid uint64, botType int64, martDomain, symbol string) (rr []protos.RiskBotModel, err error) {
	if uid <= 0 {
		return nil, nil
	}

	if len(martDomain) == 0 || len(symbol) == 0 {
		return nil, common.ErrParam
	}

	rr, err = repos.RiskBotRepo.Find(uid, botType, martDomain, symbol)
	if err != nil {
		logger.Error("Find ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return
}

func (s *riskBotService) Run(sessionUser *passportprotos.User, botId uint64, botType protos.BotType, martDomain, symbol string) (r int64, err error) {
	if sessionUser.UID <= 0 ||
		botId <= 0 {
		return 0, nil
	}

	r, err = repos.RiskBotRepo.UpdateIsRun(sessionUser.UID, botId, botType, martDomain, symbol, protos.BotIsRunning)
	common.Logger.Info("riskBotService.Run", zap.Any("uid", sessionUser.UID),
		zap.Uint64("botId", botId),
		zap.Uint64("botType", uint64(botType)),
		zap.String("martDomain", martDomain),
		zap.String("symbol", symbol),
		zap.Any("r", r),
		zap.Error(err))
	if err != nil {
		logger.Error("Run ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	log := &protos.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: protos.BUSINESS_TYPE_START_BOT,
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)

	return
}

func (s *riskBotService) Stop(sessionUser *passportprotos.User, botId uint64, botType protos.BotType, martDomain, symbol string) (r int64, err error) {
	if sessionUser.UID <= 0 || botId <= 0 {
		return 0, nil
	}

	r, err = repos.RiskBotRepo.UpdateIsRun(sessionUser.UID, botId, botType, martDomain, symbol, protos.BotIsNotRun)
	if err != nil {
		logger.Error("Stop ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	log := &protos.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: protos.BUSINESS_TYPE_STOP_BOT,
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)
	return
}

func (s *riskBotService) validate(req *protos.RiskBotModel) error {

	if req == nil {
		return common.ErrParam
	}

	if len(req.MartDomain) == 0 {
		return common.ErrBotDomain
	}
	if len(req.Symbol) == 0 {
		return common.ErrBotSymbol
	}

	if req.BotType == protos.BotTypeBalanceMonitor {
		common.Logger.Debug("riskBotService.validate: ", zap.Any("bot", req.BalanceMonitorBot))
		if req.BalanceMonitorBot == nil {
			return common.ErrParam
		}

		if len(req.BalanceMonitorBot.Mail1) == 0 {
			return common.ErrRiskBotParam
		}

		if req.BalanceMonitorBot.LaunchMode == protos.BotImmediateLaunch {
			req.IsRunning = protos.BotIsRunning
			req.StartTime = time.Now().UnixMilli()
		}
	} else if req.BotType == protos.BotTypePriceMonitor {
		common.Logger.Debug("riskBotService.validate: ", zap.Any("bot", req.PriceMonitorBot))
		if req.PriceMonitorBot == nil {
			return common.ErrParam
		}
		if req.PriceMonitorBot.LaunchMode == protos.BotImmediateLaunch {
			req.IsRunning = protos.BotIsRunning
			req.StartTime = time.Now().UnixMilli()
		}
	} else if req.BotType == protos.BotTypeBalanceScheduledReport {
		common.Logger.Debug("riskBotService.validate: ", zap.Any("bot", req.BalanceScheduleReportBot))
		if req.BalanceScheduleReportBot == nil {
			return common.ErrParam
		}

		if req.BalanceScheduleReportBot.LaunchMode == protos.BotImmediateLaunch {
			req.IsRunning = protos.BotIsRunning
			req.StartTime = time.Now().UnixMilli()
		}
	} else if req.BotType == protos.BotTypeEmergencyBrake {
		common.Logger.Debug("riskBotService.validate: ", zap.Any("bot", req.EmergencyBrakeBot))
		if req.EmergencyBrakeBot == nil {
			return common.ErrParam
		}

		if req.EmergencyBrakeBot.LaunchMode == protos.BotImmediateLaunch {
			req.IsRunning = protos.BotIsRunning
			req.StartTime = time.Now().UnixMilli()
		}
	} else if req.BotType == protos.BotTypeTGScheduledReport {
		common.Logger.Debug("riskBotService.validate: ", zap.Any("bot", req.TGScheduledReportBot))
		if req.TGScheduledReportBot == nil {
			return common.ErrParam
		}

		if req.TGScheduledReportBot.LaunchMode == protos.BotImmediateLaunch {
			req.IsRunning = protos.BotIsRunning
			req.StartTime = time.Now().UnixMilli()
		}
	} else if req.BotType == protos.BotTypeTGAlert {
		common.Logger.Debug("riskBotService.validate: ", zap.Any("bot", req.TGAlertBot))
		if req.TGAlertBot == nil {
			return common.ErrParam
		}

		if req.TGAlertBot.LaunchMode == protos.BotImmediateLaunch {
			req.IsRunning = protos.BotIsRunning
			req.StartTime = time.Now().UnixMilli()
		}
	}

	return nil
}
