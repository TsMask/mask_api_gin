package model

// SysConfig 参数配置对象 sys_config
type SysConfig struct {
	ConfigId    string `json:"config_id" gorm:"column:config_id;primary_key;AUTO_INCREMENT"` // 参数主键
	ConfigName  string `json:"config_name" gorm:"column:config_name" binding:"required"`     // 参数名称
	ConfigKey   string `json:"config_key" gorm:"column:config_key" binding:"required"`       // 参数键名
	ConfigValue string `json:"config_value" gorm:"column:config_value" binding:"required"`   // 参数键值
	ConfigType  string `json:"config_type" gorm:"column:config_type"`                        // 系统内置（Y是 N否）
	CreateBy    string `json:"create_by" gorm:"column:create_by"`                            // 创建者
	CreateTime  int64  `json:"create_time" gorm:"column:create_time"`                        // 创建时间
	UpdateBy    string `json:"update_by" gorm:"column:update_by"`                            // 更新者
	UpdateTime  int64  `json:"update_time" gorm:"column:update_time"`                        // 更新时间
	Remark      string `json:"remark" gorm:"column:remark"`                                  // 备注
}

// TableName 表名称
func (*SysConfig) TableName() string {
	return "sys_config"
}
