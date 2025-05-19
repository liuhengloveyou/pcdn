package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

var conn net.Conn

func InitTcpClient(addr string) (err error) {
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败", err)
		return
	}
	defer conn.Close()

	go processWrite(conn)

	processRead(conn)

	return nil
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

// 处理写
func processWrite(conn net.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		_, err := conn.Write([]byte("heartbeat\n"))
		fmt.Println("心跳发送:", err)

		if err != nil {
			return
		}
		<-ticker.C
	}
}

func processOneMsg(conn net.Conn, msgType uint32, msgByte []byte) error {
	switch msgType {
	// case MSGTYPE_HEARTBEAT:
	// 	return processHeartbeatMsg(conn, msgByte)
	// case MSGTYPE_GETTASK:
	// 	return processGetTaskMsg(conn, msgByte)
	// case MSGTYPE_TASKRESP:
	// 	processGetTaskRespMsg(conn, msgByte)
	default:
		Logger.Sugar().Debugf("processOneMsg type ERR: %v\n", msgType, string(msgByte))
	}

	// sendShellTask(conn, &req)

	return nil
}
