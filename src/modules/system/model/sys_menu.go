package model

// SysMenu 菜单权限对象 sys_menu
type SysMenu struct {
	MenuID     string `json:"menuId"`                      // 菜单ID
	MenuName   string `json:"menuName" binding:"required"` // 菜单名称
	ParentID   string `json:"parentId" binding:"required"` // 父菜单ID 默认0
	MenuSort   int    `json:"menuSort"`                    // 显示顺序
	Path       string `json:"path"`                        // 路由地址
	Component  string `json:"component"`                   // 组件路径
	IsFrame    string `json:"isFrame"`                     // 是否内部跳转（0否 1是）
	IsCache    string `json:"isCache"`                     // 是否缓存（0不缓存 1缓存）
	MenuType   string `json:"menuType" binding:"required"` // 菜单类型（D目录 M菜单 B按钮）
	Visible    string `json:"visible"`                     // 是否显示（0隐藏 1显示）
	Status     string `json:"status"`                      // 菜单状态（0停用 1正常）
	Perms      string `json:"perms"`                       // 权限标识
	Icon       string `json:"icon"`                        // 菜单图标（#无图标）
	CreateBy   string `json:"createBy"`                    // 创建者
	CreateTime int64  `json:"createTime"`                  // 创建时间
	UpdateBy   string `json:"updateBy"`                    // 更新者
	UpdateTime int64  `json:"updateTime"`                  // 更新时间
	Remark     string `json:"remark"`                      // 备注

	// ====== 非数据库字段属性 ======

	Children []SysMenu `json:"children,omitempty"` // 子菜单
}
