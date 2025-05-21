package api

import (
	"net/http"
	"strconv"

	"pcdn-server/common"
	"pcdn-server/protos"
	"pcdn-server/service"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func SetMartParam(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.MartParamModel
	if err := common.ReadJsonBodyFromRequest(r, &req, "struct"); err != nil {
		Logger.Errorf("SetMartParam param ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	if len(req.MartDomain) == 0 {
		Logger.Errorf("SetMartParam Domain ERR: %v\n", req)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrMartParamDomain)
		return
	}

	if len(req.AccessKey) == 0 || len(req.SecretKey) == 0 {
		Logger.Errorf("SetMartParam key ERR: %v\n", req)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrMartParamKey)
		return
	}

	req.TenantId = sessionUser.TenantID
	req.UserId = sessionUser.UID
	Logger.Infof("SetMartParam: %#v\n", req)

	id, err := service.MartParamService.Set(&req)
	if err != nil {
		Logger.Errorf("SetOneBot SetMartParam ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
}

func LoadMyMartParam(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	rr, err := service.MartParamService.Find(sessionUser.UID)
	if err != nil {
		Logger.Errorf("LoadMyMartParam service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rr)
}

func LoadOneMyMartParam(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.MartParamService.Select(sessionUser.UID, id, "")
	if err != nil {
		Logger.Errorf("LoadOneMyMartParam service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)

}

func LoadMyMartParamLite(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	rr, err := service.MartParamService.Find(sessionUser.UID)
	if err != nil {
		Logger.Errorf("LoadMyMartParam service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	for i := 0; i < len(rr); i++ {
		rr[i].AccessKey = ""
		rr[i].SecretKey = ""
		rr[i].Memo = ""
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rr)

}

func ActiveMyMartParam(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.MartParamService.Active(id, sessionUser.UID)
	if err != nil {
		Logger.Errorf("ActiveMyMartParam service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)

}

func DeleteMyMartParam(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	rst, err := service.MartParamService.Delete(id, sessionUser.UID)
	common.Logger.Info("DeleteMyMartParam: %v\n", zap.Error(err), zap.Int64("rst", rst), zap.Uint64("id", id), zap.Uint64("uid", sessionUser.UID))
	if err != nil {
		Logger.Errorf("DeleteMyMartParam service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)

}
