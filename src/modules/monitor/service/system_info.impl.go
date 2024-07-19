package service

import (
	"bufio"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
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

// NewSystemInfoService 服务层实例化
var NewSystemInfoService = &SystemInfoServiceImpl{}

// SystemInfoServiceImpl 服务器系统相关信息 服务层处理
type SystemInfoServiceImpl struct{}

// ProjectInfo 程序项目信息
func (s *SystemInfoServiceImpl) ProjectInfo() map[string]any {
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
func (s *SystemInfoServiceImpl) dependencies() map[string]string {
	var requireModules = make(map[string]string)

	// 打开 go.mod 文件
	file, err := os.Open("go.mod")
	if err != nil {
		return requireModules
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Errorf("Close go.mod file error: %s", err.Error())
		}
	}(file)

	// 逐行读取文件内容
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
			requireModules[moduleName] = version
		}
	}

	if err := scanner.Err(); err != nil {
		requireModules["scanner-err"] = err.Error()
	}
	return requireModules
}

// SystemInfo 系统信息
func (s *SystemInfoServiceImpl) SystemInfo() map[string]any {
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
func (s *SystemInfoServiceImpl) TimeInfo() map[string]string {
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
func (s *SystemInfoServiceImpl) MemoryInfo() map[string]any {
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
func (s *SystemInfoServiceImpl) CPUInfo() map[string]any {
	var core = 0
	var speed = "未知"
	var model = "未知"
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
func (s *SystemInfoServiceImpl) NetworkInfo() map[string]string {
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
func (s *SystemInfoServiceImpl) DiskInfo() []map[string]string {
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
			"cent":   fmt.Sprintf("%.1f%%", usage.UsedPercent),
			"target": partition.Device,
		})
	}
	return disks
}
