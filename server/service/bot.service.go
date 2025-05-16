package service

import (
	"arbitrage/channels"
	"arbitrage/common"
	"arbitrage/marts"
	"arbitrage/protos"
	"arbitrage/repos"
	"fmt"
	"strconv"
	"strings"
	"time"

	passportprotos "github.com/liuhengloveyou/passport/protos"

	"go.uber.org/zap"
)

type botService struct {
}

func (s *botService) Set(sessionUser *passportprotos.User, req *protos.BotModel) (uint64, error) {
	if req == nil {
		return 0, common.ErrParam
	}

	err := s.validate(req)
	if err != nil {
		logger.Error("botService.Set validate  ERR: ", zap.Error(err))
		return 0, err
	}

	id, err := repos.BotRepo.Set(req)
	if err != nil {
		logger.Error("botService.Set ERR: ", zap.Error(err))
		return id, common.ErrService
	}

	// 停掉正在运行的机器人
	if req.Id > 0 {
		channels.PostBotTask(&protos.BotTaskEvent{
			Method: protos.BotTaskStop,
			Uid:    sessionUser.UID,
			BotID:  req.Id,
		})
	}

	log := &protos.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: protos.BUSINESS_TYPE_SET_BOT,
		Payload:      fmt.Sprintf("%v", id),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)

	return id, nil
}

func (s *botService) Select(id, uid int64) (m *protos.BotModel, err error) {
	if uid <= 0 || id <= 0 {
		return nil, common.ErrParam
	}

	m, err = repos.BotRepo.Select(id, uid)
	if err != nil {
		logger.Error("select ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return m, nil
}

func (s *botService) Find(id, uid uint64, botType int64) (rr []protos.BotModel, err error) {
	if uid <= 0 {
		return nil, nil
	}

	rr, err = repos.BotRepo.Find(id, uid, botType)
	if err != nil {
		logger.Error("Find ERR: ", zap.Error(err))
		return nil, common.ErrService
	}

	return
}

func (s *botService) Delete(sessionUser *passportprotos.User, id uint64) error {
	if sessionUser == nil || sessionUser.UID <= 0 {
		return nil
	}

	err := repos.BotRepo.Delete(id, sessionUser.UID)
	if err != nil {
		logger.Error("Delete ERR: ", zap.Error(err))
		return common.ErrService
	}

	// 停掉正在运行的机器人
	channels.PostBotTask(&protos.BotTaskEvent{
		Method: protos.BotTaskStop,
		Uid:    sessionUser.UID,
		BotID:  id,
	})

	log := &protos.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: protos.BUSINESS_TYPE_DEL_BOT,
		Payload:      fmt.Sprintf("%v", id),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)

	return nil
}

func (s *botService) Run(sessionUser *passportprotos.User, botId uint64, botType int64) (r int64, err error) {
	if sessionUser.UID <= 0 || botId <= 0 {
		return 0, nil
	}

	if botType != int64(protos.BotTypeTimerReport) {
		// 只有会员能启机器人
		memberInfo, err := repos.VipRepo.LoadMemberInfo(sessionUser.UID)
		if err != nil || memberInfo == nil {
			logger.Error("LoadMemberInfo ERR: ", zap.Any("uid", sessionUser.UID), zap.Error(err))
			return 0, common.ErrNoVIP
		}
		if memberInfo.Id <= 0 || memberInfo.UserId != sessionUser.UID {
			logger.Error("LoadMemberInfo ERR: ", zap.Error(err), zap.Any("memberInfo", memberInfo))
			return 0, common.ErrNoVIP
		}
		if memberInfo.EndTime <= time.Now().UnixMilli() {
			logger.Error("MemberInfo ERR: ", zap.Error(err), zap.Any("memberInfo", memberInfo))
			return 0, common.ErrNoVIP
		}
	}

	// 是否有API_KEY, 没有API_KEY不能启动机器人
	m, err := s.Select(int64(botId), int64(sessionUser.UID))
	if err != nil {
		logger.Error("Select ERR: ", zap.Error(err))
		return 0, common.ErrMartParamKey
	}
	if m == nil {
		return 0, common.ErrMartParamKey
	}

	// 建一个交易市场API客户端
	var martA, martB marts.Mart
	if botType == int64(protos.BotTypeTimerReport) {
		martA = marts.GetMartByName(m.TimerReportBot.MartA)
		martB = marts.GetMartByName(m.TimerReportBot.MartB)
	} else if botType == int64(protos.BotTypeObserveAndReport) {
		martA = marts.GetMartByName(m.ObserveAndReportBot.MartA)
		martB = marts.GetMartByName(m.ObserveAndReportBot.MartB)
	} else if botType == int64(protos.BotTypeArbirageDepth1) {
		martA = marts.GetMartByName(m.ArbirageDepth1Bot.MartA)
		martB = marts.GetMartByName(m.ArbirageDepth1Bot.MartB)
	}
	if martA == nil || martB == nil {
		common.Logger.Error("martA.Init ERR: ", zap.Any("bot", m))
		return 0, common.ErrMartParamKey
	}
	// 查询用户的账号配置
	martParam, err := MartParamService.Select(sessionUser.UID, 0, martA.GetName())
	if err != nil || martParam == nil {
		common.Logger.Error("martA.Init ERR: ", zap.Any("bot", m))
		return 0, common.ErrMartParamKey
	}

	if err = martA.Init(&protos.MartParamModel{
		MartDomain: martParam.MartDomain,
		AccessKey:  martParam.AccessKey,
		SecretKey:  martParam.SecretKey,
		Passphrase: martParam.Passphrase,
		Memo:       martParam.Memo,
	}); err != nil {
		common.Logger.Error("martA.Init ERR: ", zap.Any("bot", m))
		return 0, common.ErrMartParamKey
	}

	// 建一个交易市场API客户端
	// 查询用户的账号配置
	martParam, err = MartParamService.Select(sessionUser.UID, 0, martB.GetName())
	if err != nil || martParam == nil {
		common.Logger.Error("martA.Init ERR: ", zap.Any("bot", m))
		return 0, common.ErrMartParamKey
	}

	if err = martB.Init(&protos.MartParamModel{
		MartDomain: martParam.MartDomain,
		AccessKey:  martParam.AccessKey,
		SecretKey:  martParam.SecretKey,
		Passphrase: martParam.Passphrase,
		Memo:       martParam.Memo,
	}); err != nil {
		common.Logger.Error("martA.Init ERR: ", zap.Any("bot", m))
		return 0, common.ErrMartParamKey
	}

	r, err = repos.BotRepo.UpdateIsRun(sessionUser.UID, botId, protos.BotIsRunning)
	if err != nil {
		logger.Error("Run ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	// 停掉正在运行的机器人
	channels.PostBotTask(&protos.BotTaskEvent{
		Method: protos.BotTaskRun,
		Uid:    sessionUser.UID,
		BotID:  botId,
	})

	log := &protos.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: protos.BUSINESS_TYPE_START_BOT,
		Payload:      fmt.Sprintf("%v", botId),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)

	return
}

func (s *botService) Stop(sessionUser *passportprotos.User, botId uint64) (r int64, err error) {
	if sessionUser.UID <= 0 || botId <= 0 {
		return 0, nil
	}

	r, err = repos.BotRepo.UpdateIsRun(sessionUser.UID, botId, protos.BotIsNotRun)
	if err != nil {
		logger.Error("Stop ERR: ", zap.Error(err))
		return 0, common.ErrService
	}

	// 停掉正在运行的机器人
	channels.PostBotTask(&protos.BotTaskEvent{
		Method: protos.BotTaskStop,
		Uid:    sessionUser.UID,
		BotID:  botId,
	})

	log := &protos.BusinessLog{
		UserName:     sessionUser.Nickname.String,
		BusinessType: protos.BUSINESS_TYPE_STOP_BOT,
		Payload:      fmt.Sprintf("%v", botId),
	}
	log.UserId = sessionUser.UID
	log.TenantId = sessionUser.TenantID
	BusinessLogService.Add(log)
	return
}

func (s *botService) validate(req *protos.BotModel) error {

	if req == nil {
		return common.ErrParam
	}

	if req.BotType == protos.BotTypeAutoArbirageDepth1 {
		common.Logger.Debug("botService.validate: ", zap.Any("bot", req.AutoArbirageDepth1Bot))
		if req.AutoArbirageDepth1Bot == nil {
			return common.ErrParam
		}

		if len(req.AutoArbirageDepth1Bot.SymbolA) == 0 || len(req.AutoArbirageDepth1Bot.SymbolB) == 0 {
			return common.ErrBotSymbol
		}

		if len(req.AutoArbirageDepth1Bot.MartA) == 0 || len(req.AutoArbirageDepth1Bot.MartB) == 0 {
			return common.ErrBotDomain
		}

		if strings.Compare(req.AutoArbirageDepth1Bot.MartA, req.AutoArbirageDepth1Bot.MartB) == 0 {
			return common.ErrBotDomain
		}

		minVolume, _ := strconv.ParseFloat(req.AutoArbirageDepth1Bot.MinVolume, 64)
		if minVolume <= 0 {
			return common.ErrParam
		}
		minSpreadRatio, _ := strconv.ParseFloat(req.AutoArbirageDepth1Bot.MinSpreadRatio, 64)
		if minSpreadRatio <= 0 {
			return common.ErrParam
		}

	} else if req.BotType == protos.BotTypeArbirageDepth1 {
		common.Logger.Debug("botService.validate: ", zap.Any("bot", req.ArbirageDepth1Bot))
		if req.ArbirageDepth1Bot == nil {
			return common.ErrParam
		}

		if len(req.ArbirageDepth1Bot.SymbolA) == 0 || len(req.ArbirageDepth1Bot.SymbolB) == 0 {
			return common.ErrBotSymbol
		}

		if len(req.ArbirageDepth1Bot.MartA) == 0 || len(req.ArbirageDepth1Bot.MartB) == 0 {
			return common.ErrBotDomain
		}

		if strings.Compare(req.ArbirageDepth1Bot.MartA, req.ArbirageDepth1Bot.MartB) == 0 {
			return common.ErrBotDomain
		}

		minVolume, _ := strconv.ParseFloat(req.ArbirageDepth1Bot.MinVolume, 64)
		if minVolume <= 0 {
			return common.ErrParam
		}
		minSpreadRatio, _ := strconv.ParseFloat(req.ArbirageDepth1Bot.MinSpreadRatio, 64)
		if minSpreadRatio <= 0 {
			return common.ErrParam
		}

		if req.ArbirageDepth1Bot.MinOrderSize < 0 {
			return common.ErrParam
		}

		if req.ArbirageDepth1Bot.MaxOrderSize < 0 {
			return common.ErrParam
		}

		if req.ArbirageDepth1Bot.MinOrderSize > req.ArbirageDepth1Bot.MaxOrderSize {
			return common.ErrParam
		}
	} else if req.BotType == protos.BotTypeObserveAndReport {
		common.Logger.Debug("botService.validate: ", zap.Any("bot", req.ObserveAndReportBot))
		if req.ObserveAndReportBot == nil {
			return common.ErrParam
		}

		if len(req.ObserveAndReportBot.TGChatId) == 0 {
			return common.ErrBotTGChatID
		}

		if len(req.ObserveAndReportBot.SymbolA) == 0 || len(req.ObserveAndReportBot.SymbolB) == 0 {
			return common.ErrBotSymbol
		}

		if len(req.ObserveAndReportBot.MartA) == 0 || len(req.ObserveAndReportBot.MartB) == 0 {
			return common.ErrBotDomain
		}

		if strings.Compare(req.ObserveAndReportBot.MartA, req.ObserveAndReportBot.MartB) == 0 {
			return common.ErrBotDomain
		}

		minVolume, _ := strconv.ParseFloat(req.ObserveAndReportBot.MinVolume, 64)
		if minVolume <= 0 {
			return common.ErrParam
		}
		minSpreadRatio, _ := strconv.ParseFloat(req.ObserveAndReportBot.MinSpreadRatio, 64)
		if minSpreadRatio <= 0 {
			return common.ErrParam
		}
	}
	return nil

}
