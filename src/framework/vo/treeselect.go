package vo

import systemModel "mask_api_gin/src/modules/system/model"

// TreeSelect 树结构实体类
type TreeSelect struct {
	// ID 节点ID
	ID string `json:"id"`

	// Label 节点名称
	Label string `json:"label"`

	// Children 子节点
	Children []TreeSelect `json:"children"`
}

// SysMenuTreeSelect 使用给定的 SysMenu 对象解析为 TreeSelect 对象
func SysMenuTreeSelect(sysMenu systemModel.SysMenu) TreeSelect {
	t := TreeSelect{}
	t.ID = sysMenu.MenuID
	t.Label = sysMenu.MenuName

	if len(sysMenu.Children) > 0 {
		for _, menu := range sysMenu.Children {
			child := SysMenuTreeSelect(menu)
			t.Children = append(t.Children, child)
		}
	} else {
		t.Children = []TreeSelect{}
	}

	return t
}

// SysDeptTreeSelect 使用给定的 SysDept 对象解析为 TreeSelect 对象
func SysDeptTreeSelect(sysDept systemModel.SysDept) TreeSelect {
	t := TreeSelect{}
	t.ID = sysDept.DeptID
	t.Label = sysDept.DeptName

	if len(sysDept.Children) > 0 {
		for _, dept := range sysDept.Children {
			child := SysDeptTreeSelect(dept)
			t.Children = append(t.Children, child)
		}
	} else {
		t.Children = []TreeSelect{}
	}

	return t
}
