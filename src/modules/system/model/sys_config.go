package model

// SysConfig 参数配置表
type SysConfig struct {
	ID              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`              // 参数ID
	CreateBy        string `json:"createBy" gorm:"column:create_by"`                          // 创建者
	CreateTime      int64  `json:"createTime" gorm:"column:create_time"`                      // 创建时间
	UpdateBy        string `json:"updateBy" gorm:"column:update_by"`                          // 更新者
	UpdateTime      int64  `json:"updateTime" gorm:"column:update_time"`                      // 更新时间
	Remark          string `json:"remark" gorm:"column:remark"`                               // 备注
	ConfigName      string `json:"configName" gorm:"column:config_name" binding:"required"`   // 参数名称
	ConfigKey       string `json:"configKey" gorm:"column:config_key" binding:"required"`     // 参数键名
	ConfigValue     string `json:"configValue" gorm:"column:config_value" binding:"required"` // 参数键值
	ConfigValueType string `json:"configValueType" gorm:"column:config_value_type"`           // 参数键值类型
	ConfigGroup     string `json:"configGroup" gorm:"column:config_group"`                    // 参数组
	ConfigSort      int64  `json:"configSort" gorm:"column:config_sort"`                      // 排序
	ConfigStatus    string `json:"configStatus" gorm:"column:config_status"`                  // 状态（0停用 1正常）
}

// TableName 表名称
func (*SysConfig) TableName() string {
	return "sys_config"
}
