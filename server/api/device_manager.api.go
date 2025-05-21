package api

import (
	"net/http"
	"strconv"
	"time"

	"pcdn-server/common"
	"pcdn-server/models"
	"pcdn-server/service"
	"pcdn-server/tcpservice"

	gocommon "github.com/liuhengloveyou/go-common"
	passportprotos "github.com/liuhengloveyou/passport/protos"
	"go.uber.org/zap"
	"gopkg.in/guregu/null.v4/zero"
)

func initDeviceManagerApi() {
	// 添加设备接口
	Apis["/device/add"] = ApiStruct{
		Handler:   AddDevice,
		Method:    "POST",
		NeedLogin: false,
	}

	Apis["/device/list"] = ApiStruct{
		Handler:   ListDevices,
		Method:    "GET",
		NeedLogin: false,
	}

	Apis["/device/update"] = ApiStruct{
		Handler:   UpdateAgent,
		Method:    "POST",
		NeedLogin: true,
	}

}

// 添加设备
func AddDevice(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)

	/////
	name := zero.StringFrom("admin")
	sessionUser = &passportprotos.User{
		UID:      1,
		Nickname: &name,
	}

	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	req := models.DeviceModel{}
	if err := common.ReadJsonBodyFromRequest(r, &req, ""); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	if req.SN == "" {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	id, err := service.DeviceService.Create(sessionUser, &req)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
}

// 列出设备
func ListDevices(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	/////TODO
	name := zero.StringFrom("admin")
	sessionUser = &passportprotos.User{
		UID:      1,
		Nickname: &name,
	}

	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()

	// 获取分页参数
	page := 1
	pageSize := 30

	if pageStr := r.FormValue("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.FormValue("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// 调用服务层方法获取设备列表
	devices, total, err := service.DeviceService.Find(sessionUser.UID, page, pageSize)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpResponseArray(w, http.StatusOK, 0, devices, total)
}

// 手动更新Agent版本
func UpdateAgent(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	sn := r.FormValue("sn")
	if sn == "" {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	agentOne := tcpservice.AgentMap[sn]
	if agentOne == nil {
		gocommon.HttpErr(w, http.StatusOK, -1, "设备不存在")
		return
	}

	common.Logger.Debug("UpdateAgent", zap.String("sn", sn), zap.Any("agent", agentOne))

	task := tcpservice.UpdateAgent(agentOne)
	if task == nil {
		gocommon.HttpErr(w, http.StatusOK, -1, "请求手机出错")
		return
	}

	select {
	// case resp := <-task.RespChan:
	// 	common.Logger.Sugar().Infof("CsmmProcessList: %v %v %v\n", agentOne, task, resp)
	// 	gocommon.HttpErr(w, http.StatusOK, 0, resp.Resp)

	case <-time.After(time.Second * 15):
		gocommon.HttpErr(w, http.StatusOK, 0, "设备应答超时")
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "OK")
}
