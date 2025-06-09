package tcpservice

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"pcdn-server/common"
	"pcdn-server/models"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var AgentMap map[string]*models.DeviceModel = make(map[string]*models.DeviceModel)

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
		fmt.Printf("read tcp: %v %v %s %v\n", conn.RemoteAddr(), n, string(buf[:n]), err)

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

func processOneMsg(conn net.Conn, msgType uint32, msgByte []byte) error {
	switch msgType {
	case uint32(protos.MsgType_MSG_TYPE_HEARTBEAT):
		return processHeartbeatMsg(conn, msgByte)
	case uint32(protos.MsgType_MSG_TYPE_TASKRESP):
		return processGetTaskRespMsg(conn, msgByte)
	default:
		common.Logger.Error("processOneMsg type ERR: ", zap.Uint32("msgType", msgType))
	}

	// sendShellTask(conn, &req)

	return nil
}

func processHeartbeatMsg(conn net.Conn, msgByte []byte) error {
	var req protos.Heartbeat
	if err := proto.Unmarshal(msgByte, &req); err != nil {
		common.Logger.Sugar().Errorf("heartbeat err: ", string(msgByte), err)
		return err
	}
	common.Logger.Debug("heartbeat: ", zap.Any("req", req))

	req.Sn = strings.ToUpper(req.Sn)
	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	tmpDevice := AgentMap[req.Sn]
	if tmpDevice == nil {
		tmpDevice = &models.DeviceModel{
			SN:         req.Sn,
			Version:    req.Ver,
			RemoteAddr: remoteAddr,
		}
		AgentMap[req.Sn] = tmpDevice
	}

	tmpDevice.SN = req.Sn
	tmpDevice.Version = req.Ver
	tmpDevice.RemoteAddr = remoteAddr
	tmpDevice.Timestamp = req.Timestamp
	tmpDevice.LastHeartbear = time.Now().UnixMilli()
	tmpDevice.ClientTcpConn = conn

	// 更新Redis中的Agent状态
	if err := updateAgentStatusToRedis(tmpDevice); err != nil {
		common.Logger.Sugar().Errorf("更新Agent状态到Redis失败: %v", err)
	}

	sendHeartbeat(conn, &protos.Heartbeat{
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

func sendHeartbeat(conn net.Conn, msg *protos.Heartbeat) error {
	if msg == nil {
		return nil
	}

	msgByte, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer([]byte("\r\n"))
	if err := binary.Write(buff, binary.LittleEndian, uint32(protos.MsgType_MSG_TYPE_HEARTBEAT)); err != nil {
		return err
	}
	if err := binary.Write(buff, binary.LittleEndian, uint32(len(msgByte))); err != nil {
		return err
	}
	if n, err := buff.Write(msgByte); n != len(msgByte) || err != nil {
		return err
	}

	if n, err := conn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		return err
	}

	return nil
}

func SendTaskToDevice(device *models.DeviceModel, task *protos.Task) error {
	for i := 0; i < 10; i++ {
		if device.ClientTcpConn != nil {
			break
		}

		time.Sleep(time.Second)
		continue // 等设备联上来取任务
	}

	if device.ClientTcpConn == nil {
		return fmt.Errorf("下发任务超时")
	}

	taskByte, err := proto.Marshal(task)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer([]byte("\r\n"))
	if err := binary.Write(buff, binary.LittleEndian, uint32(protos.MsgType_MSG_TYPE_TASK)); err != nil {
		return err
	}
	if err := binary.Write(buff, binary.LittleEndian, uint32(len(taskByte))); err != nil {
		return err
	}
	if n, err := buff.Write(taskByte); n != len(taskByte) || err != nil {
		return err
	}

	if n, err := device.ClientTcpConn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		device.ClientTcpConn.Close()
		device.ClientTcpConn = nil // 重联
		return err
	}

	return nil
}

func processGetTaskMsg(conn net.Conn, msgByte []byte) error {
	// remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]

	// var req protos.TaskReq
	// if err := json.Unmarshal(msgByte, &req); err != nil {
	// 	common.Logger.Sugar().Errorf("processGetTaskMsg msg ERR: ", conn.RemoteAddr(), string(msgByte), err)
	// 	return err
	// }

	// csmm := AgentMap[req.Name]
	// if csmm == nil {
	// 	common.Logger.Sugar().Errorf("processGetTaskMsg ERR: %v %#v\n", remoteAddr, req, AgentMap)
	// 	return nil
	// }

	// common.Logger.Sugar().Info("processGetTaskMsg msg: ", conn.RemoteAddr(), string(msgByte))
	// csmm.NcConn = &conn

	return nil
}

func processGetTaskRespMsg(conn net.Conn, msgByte []byte) error {
	common.Logger.Sugar().Debugf("processGetTaskRespMsg: %v %v\n", conn.RemoteAddr(), string(msgByte))

	// var resp protos.TaskResp
	// if err := json.Unmarshal(msgByte, &resp); err != nil {
	// 	common.Logger.Sugar().Errorf("processGetTaskRespMsg msg ERR: ", conn.RemoteAddr(), string(msgByte), err)
	// 	return
	// }

	// // remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	// csmm := AgentMap[resp.Name]
	// if csmm == nil {
	// 	common.Logger.Sugar().Errorf("processGetTaskRespMsg csmm ERR: %v %v\n", conn.RemoteAddr(), string(msgByte))
	// 	return
	// }

	// task := csmm.Tasks[resp.TaskId]
	// if task == nil {
	// 	common.Logger.Sugar().Errorf("processGetTaskRespMsg task ERR: %v %v\n", conn.RemoteAddr(), string(msgByte))
	// 	return
	// }

	// taskId, _ := strconv.Atoi(task.TaskId)
	// respTaskId, _ := strconv.Atoi(resp.TaskId)
	// if respTaskId >= taskId {
	// 	go func() {
	// 		task.RespChan <- &resp
	// 		close(task.RespChan)
	// 		delete(csmm.Tasks, task.TaskId)
	// 	}()
	// } else {
	// 	// TODO
	// }

	return nil
}
