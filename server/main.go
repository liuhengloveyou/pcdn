package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"pcdn-server/api"
	"pcdn-server/common"
	"pcdn-server/repos"
	"pcdn-server/tcpservice"
)

var (
	BuildTime string
	CommitID  string

	Sig string
)

var (
	showVer = flag.Bool("version", false, "打印版本号")
	initSys = flag.Bool("init", false, "初始化系统环境")
)

// func init() {
// 	os.Setenv("HTTP_PROXY", "http://127.0.0.1:10808")
// 	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:10808")
// 	os.Setenv("NO_PROXY", "localhost,127.0.0.1")
// }

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
}

func main() {
	flag.Parse()
	sigHandler()

	if *showVer {
		fmt.Printf("%s\t%s\n", BuildTime, CommitID)
		return
	}

	if *initSys {
		InitSystemEnv()
		return
	}

	go tcpservice.InitTcpService(common.ServConfig.TcpServerAddr)

	if err := api.InitAndRunHttpApi(common.ServConfig.HttpServerAddr); err != nil {
		panic("HTTPAPI: " + err.Error())
	}
}

func InitSystemEnv() {
	if err := repos.InitModels(common.OrmCli); err != nil {
		panic(err)
	}
}
