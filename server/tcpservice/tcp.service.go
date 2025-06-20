/*
给管理后台的API调用
*/
package tcpservice

import (
	"fmt"
	"strings"
	"time"

	"pcdn-server/common"
	"pcdn-server/models"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
)

func TrifficLimit(sn, iFaceName string, uploadLimit uint) (string, error) {
	if sn == "" {
		return "", common.ErrParam
	}
	sn = strings.ToUpper(sn)
	rate := fmt.Sprintf("%dmbit", uploadLimit)

	agentStat, err := getAgentStatusFromRedis(sn)
	if err != nil {
		return "", err
	}
	common.Logger.Debug("TrifficLimit agentStat: ", zap.Any("stat", agentStat))
	if agentStat.AccessName == "" {
		return "", common.ErrAgentNoAccess
	}

	now := time.Now().UnixMilli()
	task := &protos.Task{
		TaskId:     fmt.Sprintf("%d", now),
		TaskType:   protos.TaskType_TASK_TYPE_TC,
		Timestamp:  now,                  // 当前时间
		Sn:         sn,                   // 设备SN
		AccessName: agentStat.AccessName, // 接入服务名

		// 限速
		IfaceName: &iFaceName,
		Rate:      &rate,
	}

	err = NewTaskToRedis(task)
	if err != nil {
		common.Logger.Error("TrifficLimit NewTaskToRedis ERR: ", zap.Error(err), zap.Any("stat", task), zap.Any("stat", agentStat))
		return "", err
	}

	return task.TaskId, nil
}

func TrifficLimitStat(sn, iFaceName string) (taskId string, err error) {
	if sn == "" || iFaceName == "" {
		return "", common.ErrParam
	}
	sn = strings.ToUpper(sn)

	agentStat, err := getAgentStatusFromRedis(sn)
	if err != nil {
		return "", err
	}
	common.Logger.Debug("TrifficLimitStat agentStat: ", zap.Any("stat", agentStat))
	if agentStat.AccessName == "" {
		return "", common.ErrAgentNoAccess
	}

	now := time.Now().UnixMilli()
	task := &protos.Task{
		TaskId:     fmt.Sprintf("%d", now),
		TaskType:   protos.TaskType_TASK_TYPE_TC_STATUS,
		Timestamp:  now,                  // 当前时间
		Sn:         sn,                   // 设备SN
		AccessName: agentStat.AccessName, // 接入服务名

		// 限速
		IfaceName: &iFaceName,
	}

	err = NewTaskToRedis(task)
	if err != nil {
		common.Logger.Error("TrifficLimit NewTaskToRedis ERR: ", zap.Error(err), zap.Any("stat", task), zap.Any("stat", agentStat))
		return "", err
	}

	return task.TaskId, nil
}

func ResetDevicePWD(sn string) (*protos.Task, error) {
	if sn == "" {
		return nil, common.ErrParam
	}
	sn = strings.ToUpper(sn)

	agentStat, err := getAgentStatusFromRedis(sn)
	if err != nil {
		return nil, err
	}
	common.Logger.Debug("ResetPWD agentStat: ", zap.Any("stat", agentStat))
	if agentStat.AccessName == "" {
		return nil, common.ErrAgentNoAccess
	}

	username := "root"
	pwd := "123456"
	now := time.Now().UnixMilli()
	task := &protos.Task{
		TaskId:     fmt.Sprintf("%d", now),
		TaskType:   protos.TaskType_TASK_TYPE_RESETPWD,
		Timestamp:  now,                  // 当前时间
		Sn:         sn,                   // 设备SN
		AccessName: agentStat.AccessName, // 接入服务名

		// 重置密码字段
		Username: &username,
		Pwd:      &pwd,
	}

	err = NewTaskToRedis(task)
	if err != nil {
		common.Logger.Error("ResetPWD NewTaskToRedis ERR: ", zap.Error(err), zap.Any("stat", task), zap.Any("stat", agentStat))
		return nil, err
	}

	return task, nil
}

func GetAppList(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_APPLIST,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("GetAppList: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func GetProcessList(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_PROCLIST,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("GetProcessList: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func GetDir(DeviceAgent *protos.DeviceAgent, path string) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_DIR,
		// Payload:  path,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("GetDir: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func ChatMsg(DeviceAgent *protos.DeviceAgent, chat string) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_CHAT,
		// Payload:  chat,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("ChatMsg: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("ChatMsg SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func Contact(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_CONTACT,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("GetDir: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func Calllog(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_CALLLOG,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("Calllog: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func MessageLog(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_MESSAGE,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("MessageLog: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

// 更新agent版本
func UpdateAgent(DeviceAgent *models.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_UPDATE,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("Calendar: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func Internet(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: protos.TaskType_TASK_TYPE_RESETPWD,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("Internet: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func Gps(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := NewTask(task)
	// common.Logger.Sugar().Debugf("Gps: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("Gps SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func NetLink(DeviceAgent *protos.DeviceAgent) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_NetLink,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("NetLink: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("NetLink SendTask ERR: ", err)
	// 	return nil
	// }

	return task
}

func ScreenLive(DeviceAgent *protos.DeviceAgent, sessionId string) (task *protos.Task, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())

	task = &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_SCREENLIVE,
		// Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err = SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("SendTask ScreenLive: %#v %v %v\n", DeviceAgent, task, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SendTask ScreenLive ERR: ", err)
	// 	return
	// }

	return
}

func VideoLive(DeviceAgent *protos.DeviceAgent, videoNum, sessionId string) (task *protos.Task, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.Task{
		TaskId: taskId,
		// TaskType: videoNum,
		// Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err = SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("VideoLive SendTask: %#v %v %v\n", DeviceAgent, task, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("VideoLive SendTask ERR: ", err)
	// 	return
	// }

	return
}

func SwitchCamera(DeviceAgent *protos.DeviceAgent, sessionId string) (task *protos.Task, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_SWITCHCAMERA,
		// Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err = SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("SwitchCamera SendTask: %#v %v %v\n", DeviceAgent, task, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("SwitchCamera SendTask ERR: ", err)
	// 	return
	// }

	return
}

func AudioLive(DeviceAgent *protos.DeviceAgent, sessionId string) (task *protos.Task, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_AUDIOLIVE,
		// Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err = SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("AudioLive SendTask: %#v %v %v\n", DeviceAgent, task, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("AudioLive SendTask ERR: ", err)
	// 	return
	// }

	return
}

func Remark(DeviceAgent *protos.DeviceAgent, remark string) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_REMARK,
		// Payload:  remark,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("Remark: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("Remark ERR: ", DeviceAgent, err)
	// 	return nil
	// }

	return task
}

func ShellCmd(DeviceAgent *protos.DeviceAgent, cmd string) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_SHELL,
		// Payload:  cmd,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("ShellCmd SendTask: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("ShellCmd ERR: %v\n", err)
	// 	return nil
	// }

	return task
}

func Download(DeviceAgent *protos.DeviceAgent, path string) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_DOWNLOAD,
		// HostName: fmt.Sprintf("https://%s/api/upload", common.ServConfig.Host),
		// Payload:  path,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
	// 	return nil
	// }

	return task
}

func Upload(DeviceAgent *protos.DeviceAgent, path string) *protos.Task {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.Task{
		TaskId: taskId,
		// TaskType: TASKTYPE_UPLOAD,
		// HostName: fmt.Sprintf("http://%s/upload/", common.ServConfig.Host),
		// Payload:  path,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	// err := SendTask(DeviceAgent, task)
	// common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", DeviceAgent, err)
	// if err != nil {
	// 	common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
	// 	return nil
	// }

	return task
}

func addShellTask(DeviceAgent *protos.DeviceAgent) {
	// has := false

	// for i := 0; i < len(DeviceAgent.Tasks); i++ {
	// 	if strings.Compare(DeviceAgent.Tasks[i].TaskType, "shell") == 0 {
	// 		has = true
	// 		break
	// 	}
	// }

	// if has == false {
	// 	DeviceAgent.Tasks = append(DeviceAgent.Tasks, protos.Task{
	// 		TaskId:   time.Now().Format("2006-01-02 15:04:05"),
	// 		TaskType: "shell",
	// 		HostName: "csmm.feitian.link",
	// 		Port:     "20002",
	// 	})
	// }
}
