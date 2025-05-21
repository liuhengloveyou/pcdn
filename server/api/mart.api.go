package api

import (
	"net/http"
	"pcdn-server/common"
	"pcdn-server/service"

	gocommon "github.com/liuhengloveyou/go-common"
)

func initKlineApi() {
	Apis["/spot/kline"] = ApiStruct{
		Handler:   LoadLatestKLine,
		NeedLogin: true,
	}

	Apis["/spot/wallet"] = ApiStruct{
		Handler:   LoadSpotWallet,
		NeedLogin: true,
	}

	Apis["/spot/depth"] = ApiStruct{
		Handler:   LoadDepth,
		NeedLogin: true,
	}

}

func LoadLatestKLine(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	mart := r.FormValue("mart")
	symbol := r.FormValue("symbol")
	if len(symbol) == 0 || len(mart) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	Logger.Infoln("LoadLatestKLine: ", symbol, sessionUser.UID)
	rr, err := service.MartService.LoadLatestKLine(mart, symbol)
	if err != nil {
		Logger.Errorf("LoadLatestKLine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	Logger.Infoln("LoadLatestKLine ok: ", symbol, sessionUser.UID, len(rr))
	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}

func LoadSpotWallet(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()

	martDomain := r.FormValue("domain")
	if len(martDomain) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	Logger.Infoln("LoadSpotWallet: ", sessionUser.UID, martDomain)
	rr, err := service.MartService.LoadSpotWallet(sessionUser.UID, martDomain)
	if err != nil {
		Logger.Errorf("LoadLatestKLine service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	Logger.Infoln("LoadSpotWallet ok: ", martDomain, sessionUser.UID, len(rr))
	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}

func LoadDepth(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	mart := r.FormValue("mart")
	symbol := r.FormValue("symbol")
	if len(symbol) == 0 || len(mart) == 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	Logger.Infoln("LoadDepth: ", symbol, sessionUser.UID)
	rr, err := service.MartService.LoadDepth(mart, symbol)
	if err != nil {
		Logger.Errorf("LoadDepth service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	Logger.Infoln("LoadDepth ok: ", symbol, sessionUser.UID)
	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}
