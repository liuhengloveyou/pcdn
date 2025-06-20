package api

import (
	"context"
	"net/http"

	"pcdn-server/common"
	"pcdn-server/models"
	"pcdn-server/service"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func initTcApi() {

	// 网卡限速
	Apis["/device/tc"] = ApiStruct{
		Handler:   TrifficLimit,
		Method:    "POST",
		NeedLogin: true,
	}

	// 网卡限速状态
	Apis["/device/tc/stat"] = ApiStruct{
		Handler:   TrifficLimitStatus,
		Method:    "POST",
		NeedLogin: true,
	}

}

func TrifficLimit(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	req := models.TrifficLimitReq{}
	if err := common.ReadJsonBodyFromRequest(r, &req, ""); err != nil {
		common.Logger.Error("TrifficLimit param ERR: ", zap.Any("req", req), zap.Any("sess", sessionUser))
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	if req.SN == "" {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	if req.UploadLimit < 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	common.Logger.Debug("TrifficLimit", zap.Any("req", req), zap.Any("sess", sessionUser))

	ctx := context.WithValue(r.Context(), "UID", sessionUser.UID)
	ctx = context.WithValue(ctx, "Nickname", sessionUser.Cellphone.String)
	ctx = context.WithValue(ctx, "TID", sessionUser.TenantID)
	val, detail, err := service.TcService.TrifficLimit(ctx, req.SN, req.IfaceName, req.UploadLimit)
	if err != nil {
		common.Logger.Error("TrifficLimit", zap.Any("req", req), zap.Error(err))
		gocommon.HttpErr(w, http.StatusOK, -1, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, map[string]string{
		"val":    val,
		"detail": detail,
	})
}

func TrifficLimitStatus(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	req := models.TrifficLimitReq{}
	if err := common.ReadJsonBodyFromRequest(r, &req, ""); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	if req.SN == "" || req.IfaceName == "" {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}
	common.Logger.Debug("TrifficLimitStatus", zap.Any("device", req), zap.Any("sess", sessionUser))

	val, detail, err := service.TcService.TrifficLimitStat(req.SN, req.IfaceName)
	common.Logger.Debug("TrifficLimitStatus", zap.Any("req", req), zap.Any("val", val), zap.Any("detail", detail))
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, map[string]string{
		"val":    val,
		"detail": detail,
	})
}
