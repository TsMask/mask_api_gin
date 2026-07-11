package model

// SysMenu 菜单权限表（树形: parent_id 自关联；menu_type: D目录 M菜单 A权限）
// menu_perms 权限标识需与 controller 路由中间件 hasPerms 一致
type SysMenu struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`        // 菜单ID
	CreateBy      string `json:"createBy" gorm:"column:create_by"`                    // 创建者
	CreateTime    int64  `json:"createTime" gorm:"column:create_time"`                // 创建时间
	UpdateBy      string `json:"updateBy" gorm:"column:update_by"`                    // 更新者
	UpdateTime    int64  `json:"updateTime" gorm:"column:update_time"`                // 更新时间
	Remark        string `json:"remark" gorm:"column:remark"`                         // 备注
	ParentId      int64  `json:"parentId" gorm:"column:parent_id"`                    // 父菜单ID 默认0
	MenuName      string `json:"menuName" gorm:"column:menu_name" binding:"required"` // 菜单名称
	MenuSort      int64  `json:"menuSort" gorm:"column:menu_sort"`                    // 显示顺序
	MenuPath      string `json:"menuPath" gorm:"column:menu_path"`                    // 路由地址
	MenuComponent string `json:"menuComponent" gorm:"column:menu_component"`          // 组件路径
	MenuPerms     string `json:"menuPerms" gorm:"column:menu_perms"`                  // 权限标识
	MenuIcon      string `json:"menuIcon" gorm:"column:menu_icon"`                    // 菜单图标（#无图标）
	MenuExternal  string `json:"menuExternal" gorm:"column:menu_external"`            // 外部跳转标记（0否 1是）
	MenuCache     string `json:"menuCache" gorm:"column:menu_cache"`                  // 缓存标记（0不缓存 1缓存）
	MenuType      string `json:"menuType" gorm:"column:menu_type" binding:"required"` // 菜单类型（D目录 M菜单 A权限）
	MenuVisible   string `json:"menuVisible" gorm:"column:menu_visible"`              // 是否显示（0隐藏 1显示）
	MenuStatus    string `json:"menuStatus" gorm:"column:menu_status"`                // 菜单状态（0停用 1正常）
	MenuDel       string `json:"-" gorm:"column:menu_del"`                            // 删除标记（0存在 1删除）

	// ====== 非数据库字段属性 ======

	Children []SysMenu `json:"children,omitempty" gorm:"-"` // 子菜单
}

// TableName 表名称
func (*SysMenu) TableName() string {
	return "sys_menu"
}
