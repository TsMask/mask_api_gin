package service

import (
	"bufio"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/utils/parse"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// 实例化服务层 SystemInfoImpl 结构体
var NewSystemInfoImpl = &SystemInfoImpl{}

// SystemInfoImpl 服务器系统相关信息 服务层处理
type SystemInfoImpl struct{}

// ProjectInfo 程序项目信息
func (s *SystemInfoImpl) ProjectInfo() map[string]any {
	// 获取工作目录
	appDir, err := os.Getwd()
	if err != nil {
		appDir = ""
	}
	// 项目依赖
	dependencies := s.dependencies()
	return map[string]any{
		"appDir":       appDir,
		"env":          config.Env(),
		"name":         config.Get("framework.name"),
		"version":      config.Get("framework.version"),
		"dependencies": dependencies,
	}
}

// dependencies 读取mod内项目包依赖
func (s *SystemInfoImpl) dependencies() map[string]string {
	var pkgs = make(map[string]string)

	// 打开 go.mod 文件
	file, err := os.Open("go.mod")
	if err != nil {
		return pkgs
	}
	defer file.Close()

	// 使用 bufio.Scanner 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// 行不为空，不以module\require开头，不带有 // indirect 注释，则解析包名和版本
		prefixLine := strings.HasPrefix(line, "module") || strings.HasPrefix(line, "require") || strings.HasPrefix(line, "go ")
		suffixLine := strings.HasSuffix(line, ")") || strings.HasSuffix(line, "// indirect")
		if line == "" || prefixLine || suffixLine {
			continue
		}

		modInfo := strings.Split(line, " ")
		if len(modInfo) >= 2 {
			moduleName := strings.TrimSpace(modInfo[0])
			version := strings.TrimSpace(modInfo[1])
			pkgs[moduleName] = version
		}
	}

	if err := scanner.Err(); err != nil {
		pkgs["scanner-err"] = err.Error()
	}
	return pkgs
}

// SystemInfo 系统信息
func (s *SystemInfoImpl) SystemInfo() map[string]any {
	info, err := host.Info()
	if err != nil {
		info.Platform = err.Error()
	}
	// 用户目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}
	cmd, err := os.Executable()
	if err != nil {
		cmd = ""
	}
	return map[string]any{
		"platform":    info.Platform,
		"go":          runtime.Version(),
		"processId":   os.Getpid(),
		"arch":        info.KernelArch,
		"uname":       runtime.GOARCH,
		"release":     info.OS,
		"hostname":    info.Hostname,
		"homeDir":     homeDir,
		"cmd":         cmd,
		"execCommand": strings.Join(os.Args, " "),
	}
}

// TimeInfo 系统时间信息
func (s *SystemInfoImpl) TimeInfo() map[string]string {
	// 获取当前时间
	current := time.Now().Format("2006-01-02 15:04:05")
	// 获取程序运行时间
	uptime := time.Since(config.RunTime()).String()
	// 获取时区
	timezone := time.Now().Format("-0700 MST")
	// 获取时区名称
	timezoneName := time.Now().Format("MST")

	return map[string]string{
		"current":      current,
		"uptime":       uptime,
		"timezone":     timezone,
		"timezoneName": timezoneName,
	}
}

// MemoryInfo 内存信息
func (s *SystemInfoImpl) MemoryInfo() map[string]any {
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
func (s *SystemInfoImpl) CPUInfo() map[string]any {
	var core int = 0
	var speed string = "未知"
	var model string = "未知"
	cpuInfo, err := cpu.Info()
	if err == nil {
		core = runtime.NumCPU()
		speed = fmt.Sprintf("%.0fMHz", cpuInfo[0].Mhz)
		model = strings.TrimSpace(cpuInfo[0].ModelName)
	}

	useds := []string{}
	cpuPercent, err := cpu.Percent(0, true)
	if err == nil {
		for _, v := range cpuPercent {
			useds = append(useds, fmt.Sprintf("%.2f", v))
		}
	}

	return map[string]any{
		"model":    model,
		"speed":    speed,
		"core":     core,
		"coreUsed": useds,
	}
}

// NetworkInfo 网络信息
func (s *SystemInfoImpl) NetworkInfo() map[string]string {
	ipAddrs := make(map[string]string)
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			name := iface.Name
			if name[len(name)-1] == '0' {
				name = name[0 : len(name)-1]
				name = strings.Trim(name, "")
			}
			// ignore localhost
			if name == "lo" {
				continue
			}
			var addrs []string
			for _, v := range iface.Addrs {
				prefix := strings.Split(v.Addr, "/")[0]
				if strings.Contains(prefix, "::") {
					addrs = append(addrs, fmt.Sprintf("IPv6 %s", prefix))
				}
				if strings.Contains(prefix, ".") {
					addrs = append(addrs, fmt.Sprintf("IPv4 %s", prefix))
				}
			}
			ipAddrs[name] = strings.Join(addrs, " / ")
		}
	}
	return ipAddrs
}

// DiskInfo 磁盘信息
func (s *SystemInfoImpl) DiskInfo() []map[string]string {
	disks := make([]map[string]string, 0)

	partitions, err := disk.Partitions(false)
	if err != nil {
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
			"pcent":  fmt.Sprintf("%.1f%%", usage.UsedPercent),
			"target": partition.Device,
		})
	}
	return disks
}
