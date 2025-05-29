// 必须运行在root权限下
// setcap 'cap_net_admin=+ep' pcdnagent

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/jsimonetti/rtnetlink"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/go-selfupdate/selfupdate"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

var (
	Version   string
	BuildTime string
	CommitID  string

	Sig string

	showVer      = flag.Bool("version", false, "打印版本号")
	initSys      = flag.Bool("init", false, "初始化系统")
	tcpServer    = flag.String("tcp_server", "101.37.182.58:10001", "tcp服务地址")
	updateServer = flag.String("update_server", "http://pcdn.intelliflyt.com/update", "更新服务器地址")
	DeviceSN     = flag.String("sn", "", "设备SN")
)

// go-selfupdate setup and config
var updater = &selfupdate.Updater{
	CurrentVersion:     Version,       // 手动更新常量，或使用 `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go` 设置
	ApiURL:             *updateServer, // 托管 `$CmdName/$GOOS-$ARCH.json` 的服务器地址，该文件包含二进制文件的校验和
	BinURL:             *updateServer, // 托管二进制应用压缩包的服务器地址，作为补丁方法的备用
	DiffURL:            *updateServer, // 托管二进制补丁差异的服务器地址，用于增量更新
	Dir:                "updated/",    // 应用运行时创建的目录，用于存储 cktime 文件
	CmdName:            "agent",       // 应用名称，会附加到 ApiURL 后用于查找更新
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

	// 启动的时候更新一次
	if err := updater.BackgroundRun(); err != nil {
		fmt.Println("Failed to update app:", err)
	}

	// if err := limitUploadBandwidth("ens33", 1); err != nil {
	// 	fmt.Fprintf(os.Stderr, "限制带宽失败: %v\n", err)
	// }

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

// 限制指定网络接口的上行带宽
func limitUploadBandwidth(ifaceName string, rateMbps uint64) error {
	// 获取网络接口的索引
	rtnl, err := rtnetlink.Dial(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法连接到rtnetlink: %v\n", err)
		return err
	}
	defer rtnl.Close()

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法获取接口 %s: %v\n", ifaceName, err)
		return err
	}
	dev, err := rtnl.Link.Get(uint32(iface.Index))
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法获取接口 %s 的详细信息: %v\n", ifaceName, err)
		return err
	}
	ifaceIndex := dev.Index
	// fmt.Println(">>>>>>>>>>>>", iface.Index, ifaceIndex)

	// 创建 Tc 实例
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法打开rtnetlink套接字: %v\n", err)
		return err
	}
	defer func() {
		if err = tcnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	// For enhanced error messages from the kernel, it is recommended to set
	// option `NETLINK_EXT_ACK`, which is supported since 4.12 kernel.
	//
	// If not supported, `unix.ENOPROTOOPT` is returned.
	err = rtnl.SetOption(netlink.ExtendedAcknowledge, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not set option ExtendedAcknowledge: %v\n", err)
		return err
	}

	// // 捕获信号，退出时清理规则
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-c
	// 	cleanup(tcnl, uint32(ifaceIndex))
	// 	fmt.Println("已清理流量控制规则")
	// 	os.Exit(0)
	// }()

	showAllQdiscs(tcnl)
	// 清理旧的队列规则
	cleanup(tcnl, uint32(ifaceIndex))

	// 创建根 HTB Qdisc
	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(ifaceIndex),
			Handle:  core.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
			Info:    0,
		},
		Attribute: tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Init: &tc.HtbGlob{
					Version:      0x3,
					Rate2Quantum: 0xA,
				},
			},
		},
	}
	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "无法添加HTB qdisc: %v\n", err)
		return err
	}
	fmt.Println("HTB根队列规则已添加")

	// 创建 HTB 类，限速 1Mbps
	class := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(ifaceIndex),
			Handle:  core.BuildHandle(0x1, 0x1),
			Parent:  core.BuildHandle(0x1, 0x0),
		},
		Attribute: tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Parms: &tc.HtbOpt{
					Rate:    tc.RateSpec{Rate: 1000 * 1000}, // 1Mbps，单位 bit/s
					Ceil:    tc.RateSpec{Rate: 1000 * 1000}, // 上限速率也设为1Mbps
					Buffer:  1540,
					Cbuffer: 1540,
					Quantum: 10000,
					Level:   0,
					Prio:    0,
				},
			},
		},
	}
	if err := tcnl.Class().Add(&class); err != nil {
		fmt.Fprintf(os.Stderr, "无法添加HTB类: %v\n", err)
		return err
	}
	fmt.Println("HTB限速类已添加")

	// 添加 SFQ 子队列确保公平排队
	sfq := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(ifaceIndex),
			Handle:  core.BuildHandle(0x2, 0x0),
			Parent:  core.BuildHandle(0x1, 0x1),
		},
		Attribute: tc.Attribute{
			Kind: "sfq",
			Sfq: &tc.Sfq{
				V0: tc.SfqQopt{
					Quantum:       1514,
					PerturbPeriod: 100,
				},
			},
		},
	}
	if err := tcnl.Qdisc().Add(&sfq); err != nil {
		fmt.Fprintf(os.Stderr, "无法添加SFQ队列: %v\n", err)
		return err
	}
	fmt.Println("SFQ公平排队已添加")

	// 添加过滤器将所有流量导向 HTB 类
	classID := core.BuildHandle(0x1, 0x1)
	filter := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(ifaceIndex),
			Parent:  core.BuildHandle(0x1, 0x0),
			Handle:  0,
			Info:    0x300,
		},
		Attribute: tc.Attribute{
			Kind: "u32",
			U32: &tc.U32{
				ClassID: &classID,
				Sel: &tc.U32Sel{
					Flags: 0,
					NKeys: 0, // 不使用任何键，匹配所有流量
				},
			},
		},
	}
	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "无法添加过滤器: %v\n", err)
		return err
	}
	fmt.Println("U32过滤器已添加")

	fmt.Printf("已成功限制接口 %s 的上行带宽为 1 Mbps\n", ifaceName)
	fmt.Println("带宽限制已应用。按Ctrl+C退出。")

	// 保持程序运行，直到收到中断信号
	select {}
}

func showAllQdiscs(tcnl *tc.Tc) {
	// get all the qdiscs from all interfaces
	qdiscs, err := tcnl.Qdisc().Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get qdiscs: %v\n", err)
		return
	}

	for _, qdisc := range qdiscs {
		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get interface from id %d: %v", qdisc.Ifindex, err)
			return
		}
		fmt.Printf("%20s\t%s\n", iface.Name, qdisc.Kind)
	}
}

// 清理现有的 Qdisc
func cleanup(tcnl *tc.Tc, ifaceIndex uint32) {
	qdiscs, err := tcnl.Qdisc().Get()
	if err != nil {
		log.Printf("failed to get qdiscs: %v", err)
		return
	}
	for _, qdisc := range qdiscs {
		if qdisc.Ifindex == ifaceIndex {
			if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
				log.Printf("failed to delete qdisc: %v", err)
			} else {
				fmt.Printf("Deleted qdisc on interface %d\n", ifaceIndex)
			}
		}
	}
}
