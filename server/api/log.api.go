package api

import (
	"net/http"
	"pcdn-server/common"
	"pcdn-server/service"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
)

func initBusinessLogApi() {
	Apis["/log/list"] = ApiStruct{
		Handler:   ListMyBusinessLog,
		NeedLogin: true,
	}

}
func ListMyBusinessLog(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	martDomain := r.FormValue("domain")
	// symbol := r.FormValue("symbol")
	// if len(martDomain) == 0 || len(symbol) == 0 {
	// 	gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
	// 	return
	// }
	Logger.Infoln("ListMyBusinessLog: ", martDomain, sessionUser.UID)
	rr, err := service.BusinessLogService.Find(sessionUser.UID, martDomain, int(page))
	if err != nil {
		Logger.Errorf("ListMyBusinessLog service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	Logger.Infoln("ListMyBusinessLog ok: ", martDomain, sessionUser.UID, len(rr))
	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}
