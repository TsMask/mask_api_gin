package model

// 参数配置对象 sys_config
type SysConfig struct {
	// 参数主键
	ConfigID string `json:"configId"`
	// 参数名称
	ConfigName string `json:"configName"  binding:"required"`
	// 参数键名
	ConfigKey string `json:"configKey"  binding:"required"`
	// 参数键值
	ConfigValue string `json:"configValue"  binding:"required"`
	// 系统内置（Y是 N否）
	ConfigType string `json:"configType"`
	// 创建者
	CreateBy string `json:"createBy"`
	// 创建时间
	CreateTime int64 `json:"createTime"`
	// 更新者
	UpdateBy string `json:"updateBy"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
	// 备注
	Remark string `json:"remark"`
}
