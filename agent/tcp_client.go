package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/liuhengloveyou/pcdn/protos"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// var conn net.Conn

func InitTcpClient(addr string) (err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		Logger.Error("连接服务端失败", zap.Error(err))
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
			Logger.Sugar().Errorf("read client ERR: %v %v\n", conn.RemoteAddr(), err)
			if strings.Contains(err.Error(), "i/o timeout") {
				continue
			}

			break
		}
		Logger.Sugar().Debug("read tcp: %v %v %s %v\n", conn.RemoteAddr(), n, string(buf[:n]), err)

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
				Logger.Sugar().Infoln("read tcp msg: %v %v %v %v\n", conn.RemoteAddr(), msgType, msgLen, data.Len())
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
		// _, err := conn.Write([]byte("heartbeat\n"))

		if err := sendHeartbeat(conn); err != nil {
			return
		}
		<-ticker.C
	}
}

func processOneMsg(conn net.Conn, msgType uint32, msgByte []byte) error {
	switch msgType {
	case uint32(protos.MsgType_MSG_TYPE_HEARTBEAT):
		return processHeartbeatMsg(conn, msgByte)
	// case uint32(protos.MsgType_MSG_TYPE_GET_TASK):
	// 	return processGetTaskMsg(conn, msgByte)
	// case uint32(protos.MsgType_MSG_TYPE_TASK_RESP):
	// 	processGetTaskRespMsg(conn, msgByte)
	default:
		Logger.Sugar().Debugf("processOneMsg type ERR: %v\n", msgType, string(msgByte))
	}

	// sendShellTask(conn, &req)

	return nil
}

// 发送心跳
func sendHeartbeat(conn net.Conn) error {
	// 创建心跳包
	heartbeat := &protos.Heartbeat{
		Sn:        "sn-000000001",
		Ver:       Version,
		Timestamp: time.Now().UnixMilli(),
	}

	if DeviceSN != nil && *DeviceSN != "" {
		heartbeat.Sn = *DeviceSN
	}

	// 序列化为二进制数据
	data, err := proto.Marshal(heartbeat)
	if err != nil {
		Logger.Error("心跳包序列化失败: ", zap.Error(err))
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
		Logger.Error("发送心跳包失败: %v", zap.Error(err))
		return err

	}

	Logger.Debug("sendHeartbeat OK: ", zap.Any("heartbeat", heartbeat))

	return nil
}

func processHeartbeatMsg(conn net.Conn, msgByte []byte) error {
	var req protos.Heartbeat
	if err := proto.Unmarshal(msgByte, &req); err != nil {
		Logger.Sugar().Errorf("heartbeat err: ", string(msgByte), err)
		return err
	}
	Logger.Debug("heartbeat: ", zap.Any("Timestamp", req.Timestamp))

	return nil
}
