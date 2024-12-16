package vo

import "mask_api_gin/src/modules/system/model"

// TreeSelect 树结构实体类
type TreeSelect struct {
	ID       string       `json:"id"`       // 节点ID
	Label    string       `json:"label"`    // 节点名称
	Children []TreeSelect `json:"children"` // 子节点
}

// SysMenuTreeSelect 使用给定的 SysMenu 对象解析为 TreeSelect 对象
func SysMenuTreeSelect(sysMenu model.SysMenu) TreeSelect {
	t := TreeSelect{}
	t.ID = sysMenu.MenuId
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
func SysDeptTreeSelect(sysDept model.SysDept) TreeSelect {
	t := TreeSelect{}
	t.ID = sysDept.DeptId
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
