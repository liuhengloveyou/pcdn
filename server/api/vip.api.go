package api

import (
	"net/http"
	"pcdn-server/common"
	"pcdn-server/protos"
	"pcdn-server/service"

	gocommon "github.com/liuhengloveyou/go-common"
)

func initVIPApi() {
	Apis["/vip/order/create"] = ApiStruct{
		Handler:   BuyVipCreateOrder,
		NeedLogin: true,
	}
	Apis["/vip/order/payed"] = ApiStruct{
		Handler:   BuyVipPayedOrder,
		NeedLogin: true,
	}

	Apis["/vip/member"] = ApiStruct{
		Handler:   VipMemberInfo,
		NeedLogin: true,
	}
}

func VipMemberInfo(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	info, err := service.VipOrderService.MemberInfo(sessionUser)
	if err != nil {
		Logger.Errorf("VipMemberInfo service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	info.Id = 0
	gocommon.HttpErr(w, http.StatusOK, 0, info)
}

func BuyVipCreateOrder(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.VIPOrderStruct
	if err := common.ReadJsonBodyFromRequest(r, &req, "struct"); err != nil {
		Logger.Errorf("BuyVipCreateOrder param ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	Logger.Infof("BuyVipCreateOrder: %#v\n", req)

	err := service.VipOrderService.Save(sessionUser, &req)
	if err != nil {
		Logger.Errorf("BuyVipCreateOrder service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, req)
}

// 用户提交付款成功。
func BuyVipPayedOrder(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.VIPOrderStruct
	if err := common.ReadJsonBodyFromRequest(r, &req, "struct"); err != nil {
		Logger.Errorf("BuyVipPayedOrder param ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	Logger.Infof("BuyVipPayedOrder: %#v\n", req)

	err := service.VipOrderService.Payed(sessionUser, int64(req.Id), req.OrderId)
	if err != nil {
		Logger.Errorf("BuyVipPayedOrder service ERR: %v\n", err)
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, req)
}
