package api

import (
	"net/http"
	"pcdn-server/common"
	"pcdn-server/service"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
)

func initOrderApi() {
	Apis["/order/list"] = ApiStruct{
		Handler:   ListMyOrder,
		NeedLogin: true,
	}

}

func ListMyOrder(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	martDomain := r.FormValue("domain")
	symbol := r.FormValue("symbol")
	// if len(martDomain) == 0 || len(symbol) == 0 {
	// 	gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
	// 	return
	// }
	Logger.Infoln("ListMyOrder: ", martDomain, symbol, sessionUser.UID, int(page))
	rr, err := service.OrderService.Find(sessionUser.UID, "", "", int(page))
	if err != nil {
		Logger.Errorf("ListMyOrder service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	Logger.Infoln("ListMyOrder ok: ", martDomain, symbol, sessionUser.UID, len(rr))
	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}
