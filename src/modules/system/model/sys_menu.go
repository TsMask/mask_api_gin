package model

// SysMenu 菜单权限表
type SysMenu struct {
	MenuId      string `json:"menuId" gorm:"column:menu_id;primaryKey;type:int;autoIncrement"` // 菜单ID
	MenuName    string `json:"menuName" gorm:"column:menu_name" binding:"required"`            // 菜单名称
	ParentId    string `json:"parentId" gorm:"column:parent_id"`                               // 父菜单ID 默认0
	MenuSort    int64  `json:"menuSort" gorm:"column:menu_sort"`                               // 显示顺序
	MenuPath    string `json:"menuPath" gorm:"column:menu_path"`                               // 路由地址
	Component   string `json:"component" gorm:"column:component"`                              // 组件路径
	FrameFlag   string `json:"frameFlag" gorm:"column:frame_flag"`                             // 内部跳转标记（0否 1是）
	CacheFlag   string `json:"cacheFlag" gorm:"column:cache_flag"`                             // 缓存标记（0不缓存 1缓存）
	MenuType    string `json:"menuType" gorm:"column:menu_type" binding:"required"`            // 菜单类型（D目录 M菜单 A访问权限）
	VisibleFlag string `json:"visibleFlag" gorm:"column:visible_flag"`                         // 是否显示（0隐藏 1显示）
	StatusFlag  string `json:"statusFlag" gorm:"column:status_flag"`                           // 菜单状态（0停用 1正常）
	Perms       string `json:"perms" gorm:"column:perms"`                                      // 权限标识
	Icon        string `json:"icon" gorm:"column:icon"`                                        // 菜单图标（#无图标）
	DelFlag     string `json:"-" gorm:"column:del_flag"`                                       // 删除标记（0存在 1删除）
	CreateBy    string `json:"createBy" gorm:"column:create_by"`                               // 创建者
	CreateTime  int64  `json:"createTime" gorm:"column:create_time"`                           // 创建时间
	UpdateBy    string `json:"updateBy" gorm:"column:update_by"`                               // 更新者
	UpdateTime  int64  `json:"updateTime" gorm:"column:update_time"`                           // 更新时间
	Remark      string `json:"remark" gorm:"column:remark"`                                    // 备注

	// ====== 非数据库字段属性 ======

	Children []SysMenu `json:"children,omitempty" gorm:"-"` // 子菜单
}

// TableName 表名称
func (*SysMenu) TableName() string {
	return "sys_menu"
}
