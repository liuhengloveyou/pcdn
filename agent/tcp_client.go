package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"pcdnagent/common"
	"pcdnagent/logics"
	"pcdnagent/proxy"

	"github.com/liuhengloveyou/pcdn/protos"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var taskCh = make(chan *protos.Task, 100)

func InitTcpClient(addr string) (err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		common.Logger.Error("net.Dial ", zap.Error(err))
		return
	}
	defer conn.Close()

	go processRead(conn)
	processWrite(conn)

	return
}

// 处理读
func processRead(conn net.Conn) {
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
		common.Logger.Debug("read tcp: ", zap.String("addr", conn.RemoteAddr().String()), zap.Any("n", n), zap.Error(err))

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
				common.Logger.Sugar().Infof("read tcp msg: %v %v %v %v\n", conn.RemoteAddr(), msgType, msgLen, data.Len())
				if data.Len()-(i+10) >= int(msgLen) {
					processOneMsg(conn, msgType, dataByte[i+10:i+10+int(msgLen)])

					data = bytes.NewBuffer(dataByte[i+6+int(msgLen):]) // 重置数据
				}
			}
		}
	}
}

// 处理写
func processWrite(conn net.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case task := <-taskCh:
			if err := processTaskReal(conn, task); err != nil {
				common.Logger.Error("processTaskReal ERR: ", zap.Error(err))
			}
		case <-ticker.C:
			if err := sendHeartbeat(conn); err != nil {
				return
			}
		}
	}
}

func processOneMsg(conn net.Conn, msgType uint32, msgByte []byte) error {
	switch msgType {
	case uint32(protos.MsgType_MSG_TYPE_HEARTBEAT):
		return processHeartbeatMsg(msgByte)
	case uint32(protos.MsgType_MSG_TYPE_TASK):
		return processTaskMsg(msgByte)
	case uint32(protos.MsgType_MSG_TYPE_HTTP_PROXY_REQ):
		return proxy.ProcessHttpProxyReqMsg(conn, msgByte)
	default:
		common.Logger.Sugar().Debugf("processOneMsg type ERR: %v\n", msgType, string(msgByte))
	}

	// sendShellTask(conn, &req)

	return nil
}

// 发送心跳
func sendHeartbeat(conn net.Conn) error {
	// 创建心跳包
	heartbeat := &protos.Heartbeat{
		Sn:        "SN-00001",
		Ver:       Version,
		Timestamp: time.Now().UnixMilli(),
		Monitor:   &protos.SystemMonitorData{},
	}

	if DeviceSN != nil && *DeviceSN != "" {
		heartbeat.Sn = strings.ToUpper(*DeviceSN)
	}

	// PS进程信息
	logics.FillProcessInfo(heartbeat)

	// 网络流量信息
	logics.FillNetworkInfo(heartbeat)

	// 序列化为二进制数据
	data, err := proto.Marshal(heartbeat)
	if err != nil {
		common.Logger.Error("心跳包序列化失败: ", zap.Error(err))
		return err
	}

	// 构建消息头
	buf := new(bytes.Buffer)
	buf.Write([]byte("\r\n"))

	// 写入消息类型 (假设心跳消息类型为1)
	msgType := uint32(protos.MsgType_MSG_TYPE_HEARTBEAT)
	binary.Write(buf, binary.LittleEndian, msgType)

	// 写入消息长度
	msgLen := uint32(len(data))
	binary.Write(buf, binary.LittleEndian, msgLen)

	// 写入消息体
	buf.Write(data)

	// 发送消息
	if conn == nil {
		return fmt.Errorf("连接未建立，无法发送心跳包")
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		common.Logger.Error("发送心跳包失败: %v", zap.Error(err))
		return err
	}

	common.Logger.Debug("sendHeartbeat OK: ", zap.Any("sn", heartbeat.Sn), zap.Any("ver", heartbeat.Ver), zap.Any("ts", heartbeat.Timestamp), zap.Any("msgLen", msgLen))

	return nil
}

// 发送心跳
func sendTaskResp(conn net.Conn, task *protos.Task) error {
	// 发送消息
	if conn == nil {
		return fmt.Errorf("sendTaskResp ERR: conn nil")
	}
	if task == nil {
		return fmt.Errorf("sendTaskResp ERR: task nil")
	}

	// 序列化为二进制数据
	data, err := proto.Marshal(task)
	if err != nil {
		common.Logger.Error("序列化失败: ", zap.Error(err))
		return err
	}
	common.Logger.Sugar().Debug("sendTaskResp: ", task)

	// 构建消息头
	buf := new(bytes.Buffer)
	buf.Write([]byte("\r\n"))

	// 写入消息类型 (假设心跳消息类型为1)
	msgType := uint32(protos.MsgType_MSG_TYPE_TASKRESP)
	binary.Write(buf, binary.LittleEndian, msgType)

	// 写入消息长度
	msgLen := uint32(len(data))
	binary.Write(buf, binary.LittleEndian, msgLen)

	// 写入消息体
	buf.Write(data)

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		common.Logger.Error("conn.Write ERR: ", zap.Error(err))
		return err

	}

	return nil
}

func processHeartbeatMsg(msgByte []byte) error {
	var req protos.Heartbeat
	if err := proto.Unmarshal(msgByte, &req); err != nil {
		common.Logger.Sugar().Errorf("heartbeat err: ", string(msgByte), err)
		return err
	}
	common.Logger.Debug("heartbeat: ", zap.Any("sn", req.Sn), zap.Any("ver", req.Ver), zap.Any("Timestamp", req.Timestamp))

	return nil
}

func processTaskMsg(msgByte []byte) error {
	var task protos.Task
	if err := proto.Unmarshal(msgByte, &task); err != nil {
		common.Logger.Sugar().Errorf("processTaskMsg err: ", string(msgByte), err)
		return err
	}
	common.Logger.Debug("processTaskMsg: ", zap.Any("task", task.String()))

	taskCh <- &task

	return nil
}

func processTaskReal(conn net.Conn, task *protos.Task) error {
	common.Logger.Debug("processTaskReal: ", zap.Any("task", task.String()))

	var resp string = "OK"
	var err error
	if task.TaskType == protos.TaskType_TASK_TYPE_RESETPWD {
		// 重置密码
		err = logics.ResetRootPWD(task.Username, task.Pwd)
	} else if task.TaskType == protos.TaskType_TASK_TYPE_TC {
		if task.Rate == nil {
			return fmt.Errorf("rate or targetIP or ifaceName is nil")
		}
		targetIp := strings.Split(*tcpServer, ":")[0]
		if targetIp == "" && task.TargetIp != nil {
			targetIp = *task.TargetIp
		}
		// 网卡限速
		err = logics.ApplyLimitUploadBandwidthRules(*task.IfaceName, *task.Rate, targetIp)
	} else if task.TaskType == protos.TaskType_TASK_TYPE_TC_CLEAN {
		// 清除限速
		logics.ClearAllLimitUploadBandwidthRules()
	} else if task.TaskType == protos.TaskType_TASK_TYPE_TC_STATUS {
		if task.IfaceName == nil {
			return fmt.Errorf("rate or targetIP or ifaceName is nil")
		}
		var rate string
		rate, resp, err = logics.GetTCStatus(*task.IfaceName)
		task.Rate = &rate
		task.ErrMsg = resp
	} else if task.TaskType == protos.TaskType_TASK_TYPE_ROUTER_ADMIN {
		// 路由器管理功能
		err = logics.HandleRouterAdmin(task)
		// TODO: 代理功能需要在protobuf中定义相应的字段和任务类型
		// } else if task.TaskType == protos.TaskType_TASK_TYPE_PROXY_CREATE {
		// 	// 创建代理连接
		// 	if task.ProxyId == nil || task.TargetHost == nil || task.TargetPort == nil {
		// 		return fmt.Errorf("proxy_id, target_host or target_port is nil")
		// 	}
		// 	err = proxy.CreateProxyConnection(*task.ProxyId, *task.TargetHost, int(*task.TargetPort))
		// } else if task.TaskType == protos.TaskType_TASK_TYPE_PROXY_REMOVE {
		// 	// 移除代理连接
		// 	if task.ProxyId == nil {
		// 		return fmt.Errorf("proxy_id is nil")
		// 	}
		// 	proxy.RemoveProxyConnection(*task.ProxyId)
	}

	if err != nil {
		sendTaskResp(conn, task)
	} else {
		sendTaskResp(conn, task)
	}

	return nil
}

// TODO: 代理请求消息处理函数，需要在protobuf中定义相应的消息类型后启用
// func processProxyRequestMsg(msgByte []byte) error {
// 	var req protos.ProxyRequest
// 	if err := proto.Unmarshal(msgByte, &req); err != nil {
// 		common.Logger.Sugar().Errorf("processProxyRequestMsg err: ", string(msgByte), err)
// 		return err
// 	}
// 	common.Logger.Debug("processProxyRequestMsg: ", zap.Any("req", req.String()))
//
// 	// 这里可以处理来自服务器的代理请求
// 	// 例如：建立反向隧道连接等
//
// 	return nil
// }
