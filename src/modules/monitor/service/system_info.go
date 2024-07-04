package service

// ISystemInfoService 服务器系统相关信息 服务层接口
type ISystemInfoService interface {
	// ProjectInfo 程序项目信息
	ProjectInfo() map[string]any

	// SystemInfo 系统信息
	SystemInfo() map[string]any

	// TimeInfo 系统时间信息
	TimeInfo() map[string]string

	// MemoryInfo 内存信息
	MemoryInfo() map[string]any

	// CPUInfo CPU信息
	CPUInfo() map[string]any

	// NetworkInfo 网络信息
	NetworkInfo() map[string]string

	// DiskInfo 磁盘信息
	DiskInfo() []map[string]string
}
