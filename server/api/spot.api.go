package api

import (
	"arbitrage/common"
	"arbitrage/crawler"
	"arbitrage/protos"
	"arbitrage/service"
	"net/http"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func initSpotApi() {
	Apis["/spot/neworder"] = ApiStruct{
		Handler:   NewSpotOrder,
		NeedLogin: true,
	}

	Apis["/spot/tickers"] = ApiStruct{
		Handler:   LoadTickersData,
		NeedLogin: false,
	}

}

func LoadTickersData(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	tickers := crawler.DumpTickCacheToTable()

	common.Logger.Debug("LoadTickersData: ", zap.Any("tickers", len(tickers)))

	gocommon.HttpErr(w, http.StatusOK, 0, tickers)
}

func NewSpotOrder(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.OrderProp
	if err := common.ReadJsonBodyFromRequest(r, &req, ""); err != nil {
		Logger.Errorf("NewSpotOrder param ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	if len(req.MartDomain) == 0 {
		Logger.Errorf("NewSpotOrder Domain ERR: %v\n", req)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOrderParam)
		return
	}
	if len(req.Side) == 0 {
		Logger.Errorf("NewSpotOrder Side ERR: %v\n", req)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOrderParam)
		return
	}

	if len(req.Price) == 0 {
		Logger.Errorf("NewSpotOrder key ERR: %v\n", req)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOrderParam)
		return
	}
	if len(req.Quantity) == 0 {
		Logger.Errorf("NewSpotOrder key ERR: %v\n", req)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOrderParam)
		return
	}
	req.Type = protos.ORDER_TYPE_LIMIT // 限价单

	common.Logger.Info("NewSpotOrder: ", zap.Uint64("uid", sessionUser.UID), zap.Any("req", req))

	err := service.MartService.NewOrder(sessionUser, &req)
	if err != nil {
		Logger.Errorf("NewSpotOrder service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	common.Logger.Info("NewSpotOrder: ", zap.Uint64("uid", sessionUser.UID), zap.Any("req", req), zap.Error(err))

	gocommon.HttpErr(w, http.StatusOK, 0, "OK")
}
