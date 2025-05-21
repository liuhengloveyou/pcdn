package api

import (
	"net/http"
	"pcdn-server/common"
	"pcdn-server/protos"
	"pcdn-server/service"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func initRiskBotApi() {
	// 配置机器人
	Apis["/riskman/set"] = ApiStruct{
		Handler:   SetOneRiskBot,
		NeedLogin: true,
	}
	Apis["/riskman/load"] = ApiStruct{
		Handler:   LoadMyRiskBot,
		NeedLogin: true,
	}
	Apis["/riskman/take"] = ApiStruct{
		Handler:   LoadOneMyRiskBot,
		NeedLogin: true,
	}
	Apis["/riskman/run"] = ApiStruct{
		Handler:   RunMyRiskBot,
		NeedLogin: true,
	}
	Apis["/riskman/stop"] = ApiStruct{
		Handler:   StopMyRiskBot,
		NeedLogin: true,
	}
}

func SetOneRiskBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.RiskBotModel
	if err := common.ReadJsonBodyFromRequest(r, &req, ""); err != nil {
		Logger.Errorf("SetOneRiskBot param ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	req.TenantId = sessionUser.TenantID
	req.UserId = sessionUser.UID

	id, err := service.RiskBotService.Set(&req)
	if err != nil {
		Logger.Errorf("SetOneRiskBot service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
}

func LoadMyRiskBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	martDomain := r.FormValue("domain")
	symbol := r.FormValue("symbol")
	if len(martDomain) == 0 || len(symbol) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rr, err := service.RiskBotService.Find(sessionUser.UID, 0, martDomain, symbol)
	if err != nil {
		Logger.Errorf("ListMyMachine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rr)

}

func LoadOneMyRiskBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	botType, err := strconv.ParseInt(r.FormValue("t"), 10, 64)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	r.ParseForm()
	martDomain := r.FormValue("domain")
	symbol := r.FormValue("symbol")
	if len(martDomain) == 0 || len(symbol) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	common.Logger.Debug("LoadOneMyRiskBot: ", zap.String("domain", martDomain), zap.String("symbol", symbol), zap.Int64("botType", botType))

	rr, err := service.RiskBotService.Find(sessionUser.UID, botType, martDomain, symbol)
	if err != nil {
		Logger.Errorf("RiskBotService service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}
	common.Logger.Debug("LoadOneMyRiskBot: ", zap.String("domain", martDomain), zap.String("symbol", symbol), zap.Int64("botType", botType), zap.Any("rst", rr))

	if len(rr) != 1 {
		gocommon.HttpErr(w, http.StatusOK, 0, nil)
	} else {
		gocommon.HttpErr(w, http.StatusOK, 0, rr[0])
	}
}

func RunMyRiskBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	botId, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	botType, err1 := strconv.ParseInt(r.FormValue("t"), 10, 64)
	Logger.Infof("StopMyBot: %v %v %v %v\n", botId, botType, err, err1)
	if err != nil || err1 != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	martDomain := r.FormValue("domain")
	symbol := r.FormValue("symbol")
	Logger.Infof("StopMyBot: %v %v %v %v\n", botId, botType, martDomain, symbol)
	if len(martDomain) == 0 || len(symbol) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.RiskBotService.Run(sessionUser, botId, protos.BotType(botType), martDomain, symbol)
	if err != nil {
		Logger.Errorf("RiskBotService service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)

}

func StopMyRiskBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	botId, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	botType, err1 := strconv.ParseInt(r.FormValue("t"), 10, 64)
	Logger.Infof("StopMyBot: %v %v %v %v\n", botId, botType, err, err1)
	if err != nil || err1 != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	martDomain := r.FormValue("domain")
	symbol := r.FormValue("symbol")
	if len(martDomain) == 0 || len(symbol) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.RiskBotService.Stop(sessionUser, botId, protos.BotType(botType), martDomain, symbol)
	common.Logger.Info("StopMyBot: ",
		zap.Uint64("id", botId),
		zap.Int64("botType", botType),
		zap.String("martDomain", martDomain),
		zap.String("symbol", symbol),
		zap.Any("rst", rst),
		zap.Error(err))

	if err != nil {
		Logger.Errorf("ListMyMachine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)
}
