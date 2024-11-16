package model

// SysConfig 参数配置表
type SysConfig struct {
	ConfigId    int64  `json:"configId" gorm:"column:config_id;primary_key"`              // 参数ID
	ConfigName  string `json:"configName" gorm:"column:config_name" binding:"required"`   // 参数名称
	ConfigKey   string `json:"configKey" gorm:"column:config_key" binding:"required"`     // 参数键名
	ConfigValue string `json:"configValue" gorm:"column:config_value" binding:"required"` // 参数键值
	ConfigType  string `json:"configType" gorm:"column:config_type"`                      // 系统内置（Y是 N否）
	DelFlag     string `json:"-" gorm:"column:del_flag"`                                  // 删除标记（0存在 1删除）
	CreateBy    string `json:"createBy" gorm:"column:create_by"`                          // 创建者
	CreateTime  int64  `json:"createTime" gorm:"column:create_time"`                      // 创建时间
	UpdateBy    string `json:"updateBy" gorm:"column:update_by"`                          // 更新者
	UpdateTime  int64  `json:"updateTime" gorm:"column:update_time"`                      // 更新时间
	Remark      string `json:"remark" gorm:"column:remark"`                               // 备注
}

// TableName 表名称
func (*SysConfig) TableName() string {
	return "sys_config"
}
