package api

import (
	"arbitrage/common"
	martsws "arbitrage/marts-ws"
	"arbitrage/service"
	"net/http"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
)

func initAdminApi() {
	Apis["/admin/deep"] = ApiStruct{
		Handler:   ListDeep,
		NeedLogin: true,
	}

	Apis["/admin/order"] = ApiStruct{
		Handler:   ListAllOrder,
		NeedLogin: true,
	}
}
func ListDeep(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}
	if sessionUser.UID != uint64(common.ServConfig.AdminUID) {
		return
	}

	deeps := martsws.DumpBookTicker()
	gocommon.HttpErr(w, http.StatusOK, 0, deeps)
}

func ListAllOrder(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}
	if sessionUser.UID != uint64(common.ServConfig.AdminUID) {
		return
	}

	r.ParseForm()
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	martDomain := r.FormValue("domain")
	symbol := r.FormValue("symbol")
	rr, err := service.OrderService.Find(0, "", "", int(page))
	if err != nil {
		Logger.Errorf("ListAllOrder service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	Logger.Infoln("ListAllOrder ok: ", martDomain, symbol, sessionUser.UID, len(rr))
	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}
