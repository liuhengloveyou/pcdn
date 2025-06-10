package logics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

func PS() {
	// host
	hostInfo, err := host.Info()
	if err != nil {
		fmt.Println("get host info fail, error: ", err)
	}
	fmt.Printf("host is: %v \n", hostInfo)

	// 系统负载
	loadInfo, err := load.Avg()
	if err != nil {
		fmt.Println("get average load fail. err: ", err)
	}
	fmt.Printf("系统平均负载: %v \n", loadInfo)

	//内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("get memory info fail. err： ", err)
	}
	// 获取总内存大小，单位GB
	memTotal := memInfo.Total / 1024 / 1024 / 1024
	// 获取已用内存大小，单位MB
	memUsed := memInfo.Used / 1024 / 1024
	// 可用内存大小
	memAva := memInfo.Available / 1024 / 1024
	// 内存可用率
	memUsedPercent := memInfo.UsedPercent
	fmt.Printf("总内存: %v GB, 已用内存: %v MB, 可用内存: %v MB, 内存使用率: %.3f %% \n", memTotal, memUsed, memAva, memUsedPercent)

	// CPU信息
	cpuPercent, _ := cpu.Percent(time.Second, true)
	fmt.Printf("CPU使用率: %.3f%% \n", cpuPercent[0])
	cpuNumber, _ := cpu.Counts(true)
	fmt.Printf("CPU核心数: %v \n", cpuNumber)

	// disk
	diskPart, err := disk.Partitions(false)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(diskPart)
	for _, dp := range diskPart {
		fmt.Println(dp)
		diskUsed, _ := disk.Usage(dp.Mountpoint)
		fmt.Printf("分区总大小: %d MB \n", diskUsed.Total/1024/1024)
		fmt.Printf("分区使用率: %.3f %% \n", diskUsed.UsedPercent)
		fmt.Printf("分区inode使用率: %.3f %% \n", diskUsed.InodesUsedPercent)
	}

	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)
}
