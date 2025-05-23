package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liuhengloveyou/go-selfupdate/selfupdate"
)

var (
	Version   string
	BuildTime string
	CommitID  string

	Sig string

	showVer   = flag.Bool("version", false, "打印版本号")
	initSys   = flag.Bool("init", false, "初始化系统")
	tcpServer = flag.String("tcp_server", "127.0.0.1:10001", "tcp服务地址")
)

// go-selfupdate setup and config
var updater = &selfupdate.Updater{
	CurrentVersion:     Version,                         // 手动更新常量，或使用 `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go` 设置
	ApiURL:             "http://localhost:10000/update", // 托管 `$CmdName/$GOOS-$ARCH.json` 的服务器地址，该文件包含二进制文件的校验和
	BinURL:             "http://localhost:10000/update", // 托管二进制应用压缩包的服务器地址，作为补丁方法的备用
	DiffURL:            "http://localhost:10000/update", // 托管二进制补丁差异的服务器地址，用于增量更新
	Dir:                "updated/",                      // 应用运行时创建的目录，用于存储 cktime 文件
	CmdName:            "agent",                         // 应用名称，会附加到 ApiURL 后用于查找更新
	ForceCheck:         true,                            // 对于此示例，除非版本为 "dev"，否则始终检查更新
	OnSuccessfulUpdate: onUpdated,
}

func onUpdated() {
	fmt.Println("Successfully updated. please restart...")

	os.Exit(0)
}

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)

		// Perform any necessary cleanup here
		// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()
		// if err := srv.Shutdown(ctx); err != nil {
		// 	fmt.Printf("Server Shutdown: %v\n", err)
		// }
		os.Exit(0)
		fmt.Println("Server exiting")
	}()
}

func main() {
	flag.Parse()
	sigHandler()

	if *showVer {
		fmt.Printf("%s\t%s\n", Version, BuildTime)
		return
	}

	if *initSys {
		// InitSystemEnv()
		return
	}

	// 启动的时候更新一次
	if err := updater.BackgroundRun(); err != nil {
		fmt.Println("Failed to update app:", err)
	}

	if tcpServer == nil || *tcpServer == "" {
		fmt.Println("tcp_server is nil")
		return
	}

	go func() {
		for {
			if err := InitTcpClient(*tcpServer); err != nil {
				fmt.Println(err)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	InitTasks()

	select {}
}
