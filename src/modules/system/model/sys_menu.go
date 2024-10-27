package model

// SysMenu 菜单权限表
type SysMenu struct {
	MenuId     string `json:"menu_id" gorm:"column:menu_id;primary_key"`            // 菜单ID
	MenuName   string `json:"menu_name" gorm:"column:menu_name" binding:"required"` // 菜单名称
	ParentId   string `json:"parent_id" gorm:"column:parent_id" binding:"required"` // 父菜单ID 默认0
	MenuSort   int64  `json:"menu_sort" gorm:"column:menu_sort"`                    // 显示顺序
	Path       string `json:"path" gorm:"column:path"`                              // 路由地址
	Component  string `json:"component" gorm:"column:component"`                    // 组件路径
	IsFrame    string `json:"is_frame" gorm:"column:is_frame"`                      // 是否内部跳转（0否 1是）
	IsCache    string `json:"is_cache" gorm:"column:is_cache"`                      // 是否缓存（0不缓存 1缓存）
	MenuType   string `json:"menu_type" gorm:"column:menu_type" binding:"required"` // 菜单类型（D目录 M菜单 B按钮）
	Visible    string `json:"visible" gorm:"column:visible"`                        // 是否显示（0隐藏 1显示）
	Status     string `json:"status" gorm:"column:status"`                          // 菜单状态（0停用 1正常）
	Perms      string `json:"perms" gorm:"column:perms"`                            // 权限标识
	Icon       string `json:"icon" gorm:"column:icon"`                              // 菜单图标（#无图标）
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                          // 备注

	// ====== 非数据库字段属性 ======

	Children []SysMenu `json:"children,omitempty" gorm:"-"` // 子菜单
}

// TableName 表名称
func (*SysMenu) TableName() string {
	return "sys_menu"
}
