package api

import (
	"arbitrage/common"
	"arbitrage/protos"
	"arbitrage/service"
	"net/http"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func SetOneBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.BotModel
	if err := common.ReadJsonBodyFromRequest(r, &req, ""); err != nil {
		Logger.Errorf("SetOneBot param ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	req.TenantId = sessionUser.TenantID
	req.UserId = sessionUser.UID
	Logger.Infof("SetOneBot: %#v\n", req)

	id, err := service.BotService.Set(sessionUser, &req)
	if err != nil {
		Logger.Errorf("SetOneBot service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
}

func LoadMyBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	rr, err := service.BotService.Find(0, sessionUser.UID, 0)
	if err != nil {
		Logger.Errorf("ListMyMachine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rr)

}

func LoadOneMyBot(w http.ResponseWriter, r *http.Request) {
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
	id, _ := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if id <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rr, err := service.BotService.Find(id, sessionUser.UID, botType)
	if err != nil {
		Logger.Errorf("ListMyMachine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	if len(rr) != 1 {
		gocommon.HttpErr(w, http.StatusOK, 0, nil)
	} else {
		gocommon.HttpErr(w, http.StatusOK, 0, rr[0])
	}
}

func DeleteOneMyBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	id, _ := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if id <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	err := service.BotService.Delete(sessionUser, id)
	if err != nil {
		Logger.Errorf("DeleteOneMyBot service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "ok")
}

func RunMyBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	botId, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	botType, err := strconv.ParseInt(r.FormValue("t"), 10, 64)
	Logger.Infof("StopMyBot: %v %v\n", botId, err)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.BotService.Run(sessionUser, botId, botType)
	if err != nil {
		Logger.Errorf("ListMyMachine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)
}

func StopMyBot(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	botId, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	Logger.Infof("StopMyBot: %v %v %v %v\n", botId, err)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.BotService.Stop(sessionUser, botId)
	common.Logger.Info("StopMyBot: ",
		zap.Uint64("id", botId),
		zap.Any("rst", rst),
		zap.Error(err))

	if err != nil {
		Logger.Errorf("ListMyMachine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)
}
