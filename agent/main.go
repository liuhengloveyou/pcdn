package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"
	"strings"
	"net/http"
)

func main() {
	serverAddr := "127.0.0.1:9000" // 替换为你的服务器地址
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()

	go sendHeartbeat(conn)

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("服务器关闭连接")
			} else {
				fmt.Println("读取数据失败:", err)
			}
			break
		}
		msg = strings.TrimSpace(msg)
		fmt.Println("收到指令:", msg)
		if strings.HasPrefix(msg, "upgrade ") {
			url := strings.TrimPrefix(msg, "upgrade ")
			upgradeSelf(url)
		}
	}
}

func sendHeartbeat(conn net.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		_, err := conn.Write([]byte("heartbeat\n"))
		if err != nil {
			fmt.Println("心跳发送失败:", err)
			return
		}
		<-ticker.C
	}
}

func upgradeSelf(url string) {
	fmt.Println("开始升级，下载:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("下载失败:", err)
		return
	}
	defer resp.Body.Close()

	file, err := os.Create("agent_new.exe")
	if err != nil {
		fmt.Println("创建新文件失败:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("写入新文件失败:", err)
		return
	}

	fmt.Println("下载完成，准备替换并重启...")
	// Windows下替换自身可执行文件较为复杂，通常需要先启动新进程再退出当前进程
	exec.Command("cmd", "/C", "timeout 1 && move /Y agent_new.exe agent.exe && start agent.exe && exit").Start()
	os.Exit(0)
}