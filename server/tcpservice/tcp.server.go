package tcpservice

import (
	"bufio"
	"bytes"
	"context"
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

var AgentMap map[string]*models.DeviceAgent

func InitTcpService(addr string) {
	AgentMap = make(map[string]*models.DeviceAgent)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			common.Logger.Error("accept failed: ", zap.Error(err))
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
		// common.Logger.Debug("read tcp: ", zap.Any("conn", conn.RemoteAddr()), zap.Any("n", n), zap.Any("buf", string(buf[:n])))

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
		return processTaskRespMsg(conn, msgByte)
	case uint32(protos.MsgType_MSG_TYPE_HTTP_PROXY_RESP):
		return processHttpProxyRespMsg(conn, msgByte)
	default:
		common.Logger.Error("processOneMsg type ERR: ", zap.Uint32("msgType", msgType))
	}

	// sendShellTask(conn, &req)

	return nil
}

func processHeartbeatMsg(conn net.Conn, msgByte []byte) error {
	var heartbeat protos.Heartbeat
	if err := proto.Unmarshal(msgByte, &heartbeat); err != nil {
		common.Logger.Sugar().Errorf("heartbeat err: ", string(msgByte), err)
		return err
	}
	common.Logger.Debug("heartbeat: ",
		zap.Any("heartbeat", heartbeat.Sn),
		zap.Any("ver", heartbeat.Ver),
		zap.Any("timestamp", heartbeat.Timestamp),
		zap.Any("netowrk", len(heartbeat.Monitor.Network)),
		zap.Any("process", len(heartbeat.Monitor.Processes)),
	)

	heartbeat.Sn = strings.ToUpper(heartbeat.Sn)
	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	tmpDevice, ok := AgentMap[heartbeat.Sn]
	if !ok {
		tmpDevice = &models.DeviceAgent{
			SN:         heartbeat.Sn,
			Version:    heartbeat.Ver,
			RemoteAddr: remoteAddr,
		}
		AgentMap[heartbeat.Sn] = tmpDevice
	}

	tmpDevice.SN = heartbeat.Sn
	tmpDevice.Version = heartbeat.Ver
	tmpDevice.RemoteAddr = remoteAddr
	tmpDevice.Timestamp = heartbeat.Timestamp
	tmpDevice.LastHeartbear = time.Now().UnixMilli()
	tmpDevice.ClientTcpConn = conn

	// 更新Redis中的Agent状态
	if err := updateAgentStatusToRedis(tmpDevice); err != nil {
		common.Logger.Error("updateAgentStatusToRedis ERR: ", zap.Error(err))
	}
	// 更新Redis中的Agent进程监控信息
	if err := updateAgentMonitorToRedis(&heartbeat); err != nil {
		common.Logger.Error("updateAgentMonitorToRedis ERR: ", zap.Error(err))
	}

	sendHeartbeat(tmpDevice, &protos.Heartbeat{
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

func sendHeartbeat(device *models.DeviceAgent, msg *protos.Heartbeat) error {
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

	device.MU.Lock()
	defer device.MU.Unlock()

	if n, err := device.ClientTcpConn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		return err
	}

	return nil
}

func SendTaskToDevice(device *models.DeviceAgent, task *protos.Task) error {
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

	device.MU.Lock()
	defer device.MU.Unlock()

	if n, err := device.ClientTcpConn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		device.ClientTcpConn.Close()
		device.ClientTcpConn = nil // 重联
		return err
	}

	return nil
}

func processTaskRespMsg(conn net.Conn, msgByte []byte) error {
	common.Logger.Sugar().Debugf("processTaskRespMsg: %v %v\n", conn.RemoteAddr(), string(msgByte))

	var task protos.Task
	if err := proto.Unmarshal(msgByte, &task); err != nil {
		common.Logger.Sugar().Errorf("processTaskRespMsg msg ERR: ", conn.RemoteAddr(), string(msgByte), err)
		return err
	}
	common.Logger.Sugar().Debugf("processTaskRespMsg: %v %v %s\n", conn.RemoteAddr(), task.TaskId, task.ErrMsg)

	// 把应答信息写到相应的redis队列里
	redisKey := fmt.Sprintf("%s%s", common.TASK_RESPONSE_KEY_PREFIX, task.TaskId)
	_, err := common.RedisClient.LPush(context.Background(), redisKey, msgByte).Result()
	if err != nil {
		common.Logger.Error("processTaskRespMsg redis ERR: ", zap.String("key", redisKey), zap.Error(err))
		return common.ErrService
	}

	// 设置超时清理
	common.RedisClient.Expire(context.Background(), redisKey, time.Minute*10).Err()

	return nil
}
