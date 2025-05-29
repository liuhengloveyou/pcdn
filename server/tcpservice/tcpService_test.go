package tcpservice

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"pcdn-server/models"

	"github.com/liuhengloveyou/pcdn/protos"
)

// go test -v -count=1 ./... -run TestTcpTask
func ReadTcpTask(conn net.Conn) {
	// conn, err := net.Dial("tcp", "127.0.0.1:20001")
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	return
	// }
	// defer conn.Close()
	time.Sleep(time.Second * 3)

	taskRep := protos.Task{}
	taskRepByte, _ := json.Marshal(taskRep)

	buff := bytes.NewBuffer([]byte("\r\n"))
	if err := binary.Write(buff, binary.LittleEndian, uint32(protos.MsgType_MSG_TYPE_HEARTBEAT)); err != nil {
		panic(err)
	}
	if err := binary.Write(buff, binary.LittleEndian, uint32(len(taskRepByte))); err != nil {
		panic(err)
	}
	if n, err := buff.Write(taskRepByte); n != len(taskRepByte) || err != nil {
		panic(err)
	}

	if n, err := conn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		panic(err)
	}
	fmt.Println("TestTcpTask req OK")

	for {
		buf := [1024]byte{}
		read, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed,err:", err)
			return
		}
		fmt.Println(string(buf[:read]))
	}
}

// go test -v -count=1 -run TestTcpHeartbeat csmm/tcpservice
func TestTcpHeartbeat(t *testing.T) {
	// conn, err := net.Dial("tcp", "127.0.0.1:20001")
	conn, err := net.Dial("tcp", "43.128.140.17:20001")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer conn.Close()

	go ReadTcpTask(conn)

	req := models.HeartbeatReq{
		Name: "testing",
		Stat: "idel",
	}

	for {
		req.Time = time.Now().Format("2006-01-02 15:04:05")
		reqByte, _ := json.Marshal(req)

		buff := bytes.NewBuffer([]byte("\r\n"))
		if err = binary.Write(buff, binary.LittleEndian, uint32(protos.MsgType_MSG_TYPE_HEARTBEAT)); err != nil {
			panic(err)
		}
		if err = binary.Write(buff, binary.LittleEndian, uint32(len(reqByte))); err != nil {
			panic(err)
		}
		if n, err := buff.Write(reqByte); n != len(reqByte) || err != nil {
			panic(err)
		}

		if n, err := conn.Write(buff.Bytes()); n != buff.Len() || err != nil {
			panic(err)
		}

		// 读应答
		// respBuf := [1024]byte{}
		// read, err := conn.Read(respBuf[:])
		// if err != nil {
		// 	fmt.Println("recv err:", err)
		// 	return
		// }

		// msgType := binary.LittleEndian.Uint32(respBuf[2:6])
		// msgLen := binary.LittleEndian.Uint32(respBuf[6:10])

		fmt.Println("TestTcpHeartbeat: ", buff.String())

		time.Sleep(time.Second)
	}
}
