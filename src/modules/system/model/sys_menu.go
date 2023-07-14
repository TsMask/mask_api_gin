package model

import frameworkModel "mask_api_gin/src/framework/model"

// SysMenu 菜单权限对象 sys_menu
type SysMenu struct {
	// 菜单ID
	MenuID string `json:"menuId"`
	// 菜单名称
	MenuName string `json:"menuName"`
	// 父菜单ID 默认0
	ParentID string `json:"parentId"`
	// 显示顺序
	MenuSort int `json:"menuSort"`
	// 路由地址
	Path string `json:"path"`
	// 组件路径
	Component string `json:"component"`
	// 是否内部跳转（0否 1是）
	IsFrame string `json:"isFrame"`
	// 是否缓存（0不缓存 1缓存）
	IsCache string `json:"isCache"`
	// 菜单类型（D目录 M菜单 B按钮）
	MenuType string `json:"menuType"`
	// 是否显示（0隐藏 1显示）
	Visible string `json:"visible"`
	// 菜单状态（0停用 1正常）
	Status string `json:"status"`
	// 权限标识
	Perms string `json:"perms"`
	// 菜单图标（#无图标）
	Icon string `json:"icon"`
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

	// ====== 非数据库字段属性 ======

	// 子菜单
	Children []SysMenu `json:"children,omitempty"`
}

// SysMenuTreeSelect 使用给定的 SysMenu 对象解析为 TreeSelect 对象
func SysMenuTreeSelect(sysMenu SysMenu) frameworkModel.TreeSelect {
	t := frameworkModel.TreeSelect{}
	t.ID = sysMenu.MenuID
	t.Label = sysMenu.MenuName

	if len(sysMenu.Children) > 0 {
		for _, menu := range sysMenu.Children {
			child := SysMenuTreeSelect(menu)
			t.Children = append(t.Children, child)
		}
	}

	return t
}
