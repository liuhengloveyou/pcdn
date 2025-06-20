package api

import (
	"net/http"
	"strconv"
	"strings"

	"pcdn-server/common"
	"pcdn-server/models"
	"pcdn-server/service"
	"pcdn-server/tcpservice"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func initDeviceManagerApi() {
	// 添加设备接口
	Apis["/device/add"] = ApiStruct{
		Handler:   AddDevice,
		Method:    "POST",
		NeedLogin: true,
	}

	// 查询设备列表
	Apis["/device/list"] = ApiStruct{
		Handler:   ListDevices,
		Method:    "GET",
		NeedLogin: true,
	}

	// 更新设备信息
	Apis["/device/update"] = ApiStruct{
		Handler:   UpdateAgent,
		Method:    "POST",
		NeedLogin: true,
	}

	// 重置密码
	Apis["/device/resetpwd"] = ApiStruct{
		Handler:   ResetDevicePWD,
		Method:    "GET",
		NeedLogin: true,
	}

	// 查询设备监控信息
	Apis["/device/monitor"] = ApiStruct{
		Handler:   GetDeviceMonitorInfo,
		Method:    "GET",
		NeedLogin: true,
	}

	// 获取路由器管理界面URL
	Apis["/device/router-admin"] = ApiStruct{
		Handler:   GetRouterAdmin,
		Method:    "GET",
		NeedLogin: true,
	}
}

// 添加设备
func AddDevice(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
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
		common.Logger.Error("AddDevice", zap.Error(err))
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
}

// 列出设备
func ListDevices(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		common.Logger.Error("session ERR:", zap.Any("sess", sessionUser))
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
	sn = strings.ToUpper(sn)

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

	// select {
	// // case resp := <-task.RespChan:
	// // 	common.Logger.Sugar().Infof("CsmmProcessList: %v %v %v\n", agentOne, task, resp)
	// // 	gocommon.HttpErr(w, http.StatusOK, 0, resp.Resp)

	// case <-time.After(time.Second * 15):
	// 	gocommon.HttpErr(w, http.StatusOK, 0, "设备应答超时")
	// }

	gocommon.HttpErr(w, http.StatusOK, 0, "OK")
}

// 重置密码
func ResetDevicePWD(w http.ResponseWriter, r *http.Request) {
	sessionUser := ReadSessionFromRequest(r)
	if sessionUser == nil || sessionUser.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	r.ParseForm()
	sn := r.FormValue("sn")
	if sn == "" {
		common.Logger.Error("ResetDevicePWD param ERR: ", zap.String("sn", sn), zap.Any("sess", sessionUser))
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		return
	}

	common.Logger.Debug("ResetDevicePWD", zap.String("sn", sn), zap.Any("sess", sessionUser))

	task, err := tcpservice.ResetDevicePWD(sn)
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, "请求手机出错")
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, task.GetTaskId())
}

// 获取设备监控信息
func GetDeviceMonitorInfo(w http.ResponseWriter, r *http.Request) {
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

	// 调用服务层方法获取设备监控信息
	monitorInfo, err := service.DeviceService.GetMonitorInfo(sn)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, monitorInfo)
}

// 获取路由器管理界面URL
func GetRouterAdmin(w http.ResponseWriter, r *http.Request) {
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

	// 调用服务层方法获取路由器管理界面URL
	url, err := service.DeviceService.GetRouterAdminURL(sn)
	if err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, map[string]string{"url": url})
}
