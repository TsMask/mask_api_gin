package model

// SysConfig 参数配置对象 sys_config
type SysConfig struct {
	ConfigID    string `json:"configId"`                       // 参数主键
	ConfigName  string `json:"configName" binding:"required"`  // 参数名称
	ConfigKey   string `json:"configKey" binding:"required"`   // 参数键名
	ConfigValue string `json:"configValue" binding:"required"` // 参数键值
	ConfigType  string `json:"configType"`                     // 系统内置（Y是 N否）
	CreateBy    string `json:"createBy"`                       // 创建者
	CreateTime  int64  `json:"createTime"`                     // 创建时间
	UpdateBy    string `json:"updateBy"`                       // 更新者
	UpdateTime  int64  `json:"updateTime"`                     // 更新时间
	Remark      string `json:"remark"`                         // 备注
}
