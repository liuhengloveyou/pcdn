package api

import (
	"net/http"
	"pcdn-server/common"
	"pcdn-server/protos"
	"pcdn-server/tcpservice"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func ListDevices(w http.ResponseWriter, r *http.Request) {
	rst := make([]*protos.AgentClient, 0)

	for _, v := range tcpservice.AgentMap {
		if v != nil {
			rst = append(rst, v)
		}

	}

	gocommon.HttpErr(w, http.StatusOK, 0, rst)
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
	case resp := <-task.RespChan:
		common.Logger.Sugar().Infof("CsmmProcessList: %v %v %v\n", agentOne, task, resp)
		gocommon.HttpErr(w, http.StatusOK, 0, resp.Resp)

	case <-time.After(time.Second * 15):
		gocommon.HttpErr(w, http.StatusOK, 0, "设备应答超时")
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "OK")
}
