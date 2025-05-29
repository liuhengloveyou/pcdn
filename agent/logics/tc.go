package logics

import (
	"fmt"
	"log"
	"net"
	"os"
	"pcdnagent/common"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/netlink"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

// 限制指定网络接口的上行带宽
func LimitUploadBandwidth(ifaceName string, rateKbps uint32) error {
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
	common.Logger.Info("iface", zap.Any("iface", iface.Index), zap.Any("ifaceIndex", ifaceIndex))

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
	// If not supported, `unix.ENOPROTOOPT` is returned.
	err = rtnl.SetOption(netlink.ExtendedAcknowledge, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not set option ExtendedAcknowledge: %v\n", err)
		return err
	}

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
		common.Logger.Error("add HTB qdisc ERR: ", zap.Error(err))
		return err
	}
	common.Logger.Info("add HTB qdisc OK")

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
					Rate:    tc.RateSpec{Rate: rateKbps * 1000}, // 1Mbps，单位 bit/s
					Ceil:    tc.RateSpec{Rate: rateKbps * 1000}, // 上限速率也设为1Mbps
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
	common.Logger.Info("add class OK")

	// 添加 SFQ 子队列确保公平排队
	// sfq := tc.Object{
	// 	Msg: tc.Msg{
	// 		Family:  unix.AF_UNSPEC,
	// 		Ifindex: uint32(ifaceIndex),
	// 		Handle:  core.BuildHandle(0x2, 0x0),
	// 		Parent:  core.BuildHandle(0x1, 0x1),
	// 	},
	// 	Attribute: tc.Attribute{
	// 		Kind: "sfq",
	// 		Sfq: &tc.Sfq{
	// 			V0: tc.SfqQopt{
	// 				Quantum:       1514,
	// 				PerturbPeriod: 100,
	// 			},
	// 		},
	// 	},
	// }
	// if err := tcnl.Qdisc().Add(&sfq); err != nil {
	// 	common.Logger.Error("add SFQ ERR", zap.Error(err))
	// 	return err
	// }
	common.Logger.Info("add sfq OK")

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
		common.Logger.Error("add filter ERR", zap.Error(err))
		return err
	}
	common.Logger.Info("add filter OK")

	return nil
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
