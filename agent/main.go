// 必须运行在root权限下
// setcap 'cap_net_admin=+ep' pcdnagent

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"pcdnagent/common"
	"pcdnagent/upgrade"
	"syscall"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"go.uber.org/zap"
)

var (
	Version   string
	BuildTime string
	CommitID  string

	Sig string

	showVer       = flag.Bool("version", false, "打印版本号")
	initSys       = flag.Bool("init", false, "初始化系统")
	tcpServer     = flag.String("tcp_server", "101.37.182.58:10001", "tcp服务地址")
	updateServer  = flag.String("update_server", "http://update.intelliflyt.com/update/", "更新服务器地址")
	upgradeServer = flag.String("upgrade_server", "http://update.intelliflyt.com/upgrade/", "升级服务器地址")
	DeviceSN      = flag.String("sn", "SN-1234567890", "设备SN")
	dnsServer     = flag.String("dns_server", "", "自定义DNS服务器地址, 如: 8.8.8.8:53")
)

// go-selfupdate setup and config
var updater = &selfupdate.Updater{
	CurrentVersion:     Version,       // 手动更新常量，或使用 `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go` 设置
	ApiURL:             *updateServer, // 托管 `$CmdName/$GOOS-$ARCH.json` 的服务器地址，该文件包含二进制文件的校验和
	BinURL:             *updateServer, // 托管二进制应用压缩包的服务器地址，作为补丁方法的备用
	DiffURL:            *updateServer, // 托管二进制补丁差异的服务器地址，用于增量更新
	Dir:                "updated/",    // 应用运行时创建的目录，用于存储 cktime 文件
	CmdName:            "pcdnagent",   // 应用名称，会附加到 ApiURL 后用于查找更新
	ForceCheck:         true,          // 对于此示例，除非版本为 "dev"，否则始终检查更新
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

	gocommon.SingleInstane("/tmp/pcdnagent.pid")

	// 初始化升级服务
	if err := checkAndUpgrade(); err != nil {
		common.Logger.Error("Failed to initialize upgrade service", zap.Error(err))
		// 不退出程序，升级功能失败不应该影响主要功能
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
				common.Logger.Error("InitTcpClient err: ", zap.Error(err))
			}

			time.Sleep(5 * time.Second)
		}
	}()

	InitTasks()

	select {}
}

// initUpgradeService 初始化升级服务
func checkAndUpgrade() error {
	if *upgradeServer == "" {
		common.Logger.Warn("Update server not configured, upgrade service disabled")
		return nil
	}

	if *DeviceSN == "" {
		common.Logger.Warn("Device SN not configured, upgrade service disabled")
		return nil
	}

	// 如果命令行指定了DNS服务器，设置环境变量
	if dnsServer != nil && *dnsServer != "" {
		os.Setenv("PCDN_DNS_SERVER", *dnsServer)
		common.Logger.Info("Using custom DNS server", zap.String("dns_server", *dnsServer))
	}

	// 创建简单升级器
	upgrader := upgrade.NewSimpleUpgrader(*upgradeServer, Version, *DeviceSN)

	// 启动时检查一次更新
	if err := upgrader.CheckAndUpgrade(); err != nil {
		common.Logger.Error("upgrader.CheckAndUpgrade ERR", zap.Error(err))
	}

	return nil
}
