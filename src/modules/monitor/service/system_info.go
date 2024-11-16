package service

import (
	"context"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/utils/parse"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// NewSystemInfo 服务层实例化
var NewSystemInfo = &SystemInfo{}

// SystemInfo 服务器系统相关信息 服务层处理
type SystemInfo struct{}

// ProjectInfo 程序项目信息
func (s SystemInfo) ProjectInfo() map[string]any {
	// 获取工作目录
	appDir, err := os.Getwd()
	if err != nil {
		appDir = ""
	}
	return map[string]any{
		"appDir":  appDir,
		"env":     config.Env(),
		"name":    config.Name,
		"version": config.Version,
	}
}

// SystemInfo 系统信息
func (s SystemInfo) SystemInfo() map[string]any {
	info, err := host.Info()
	if err != nil {
		info.Platform = err.Error()
	}
	// 用户目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}
	// 执行命令
	cmd, err := os.Executable()
	if err != nil {
		cmd = ""
	}
	// 获取主机运行时间
	bootTime := time.Since(time.Unix(int64(info.BootTime), 0)).Seconds()
	// 获取程序运行时间
	runTime := time.Since(config.RunTime()).Abs().Seconds()
	return map[string]any{
		"platform":        info.Platform,
		"platformVersion": info.PlatformVersion,
		"arch":            info.KernelArch,
		"archVersion":     info.KernelVersion,
		"os":              info.OS,
		"hostname":        info.Hostname,
		"bootTime":        int64(bootTime),
		"processId":       os.Getpid(),
		"runArch":         runtime.GOARCH,
		"runVersion":      runtime.Version(),
		"runTime":         int64(runTime),
		"homeDir":         homeDir,
		"cmd":             cmd,
		"execCommand":     strings.Join(os.Args, " "),
	}
}

// TimeInfo 系统时间信息
func (s SystemInfo) TimeInfo() map[string]string {
	now := time.Now()
	// 获取当前时间
	current := now.Format("2006-01-02 15:04:05")
	// 获取时区
	timezone := now.Format("-0700 MST")
	// 获取时区名称
	timezoneName := now.Format("MST")

	return map[string]string{
		"current":      current,
		"timezone":     timezone,
		"timezoneName": timezoneName,
	}
}

// MemoryInfo 内存信息
func (s SystemInfo) MemoryInfo() map[string]any {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		memInfo.UsedPercent = 0
		memInfo.Available = 0
		memInfo.Total = 0
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]any{
		"usage":     fmt.Sprintf("%.2f", memInfo.UsedPercent),            // 内存利用率
		"freemem":   parse.Bit(float64(memInfo.Available)),               // 可用内存大小（GB）
		"totalmem":  parse.Bit(float64(memInfo.Total)),                   // 总内存大小（GB）
		"rss":       parse.Bit(float64(memStats.Sys)),                    // 常驻内存大小（RSS）
		"heapTotal": parse.Bit(float64(memStats.HeapSys)),                // 堆总大小
		"heapUsed":  parse.Bit(float64(memStats.HeapAlloc)),              // 堆已使用大小
		"external":  parse.Bit(float64(memStats.Sys - memStats.HeapSys)), // 外部内存大小（非堆）
	}
}

// CPUInfo CPU信息
func (s SystemInfo) CPUInfo() map[string]any {
	var core = 0
	var speed = "-"
	var model = "-"
	cpuInfo, err := cpu.Info()
	if err == nil {
		core = runtime.NumCPU()
		speed = fmt.Sprintf("%.0fMHz", cpuInfo[0].Mhz)
		model = strings.TrimSpace(cpuInfo[0].ModelName)
	}

	var used []string
	cpuPercent, err := cpu.Percent(0, true)
	if err == nil {
		for _, v := range cpuPercent {
			used = append(used, fmt.Sprintf("%.2f", v))
		}
	}

	return map[string]any{
		"model":    model,
		"speed":    speed,
		"core":     core,
		"coreUsed": used,
	}
}

// NetworkInfo 网络信息
func (s SystemInfo) NetworkInfo() map[string]string {
	ipAdders := make(map[string]string)
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, v := range interfaces {
			name := v.Name
			if name[len(name)-1] == '0' {
				name = name[0 : len(name)-1]
				name = strings.Trim(name, "")
			}
			// ignore localhost
			if name == "lo" {
				continue
			}
			var adders []string
			for _, v := range v.Addrs {
				prefix := strings.Split(v.Addr, "/")[0]
				if strings.Contains(prefix, "::") {
					adders = append(adders, fmt.Sprintf("IPv6 %s", prefix))
				}
				if strings.Contains(prefix, ".") {
					adders = append(adders, fmt.Sprintf("IPv4 %s", prefix))
				}
			}
			ipAdders[name] = strings.Join(adders, " / ")
		}
	}
	return ipAdders
}

// DiskInfo 磁盘信息
func (s SystemInfo) DiskInfo() []map[string]string {
	disks := make([]map[string]string, 0)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil && err != context.DeadlineExceeded {
		return disks
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		disks = append(disks, map[string]string{
			"size":   parse.Bit(float64(usage.Total)),
			"used":   parse.Bit(float64(usage.Used)),
			"avail":  parse.Bit(float64(usage.Free)),
			"cent":   fmt.Sprintf("%.1f%%", usage.UsedPercent),
			"target": partition.Device,
		})
	}
	return disks
}
