package vo

import "mask_api_gin/src/modules/system/model"

// TreeSelect 树结构实体类
type TreeSelect struct {
	ID       int64        `json:"id"`       // 节点ID
	Label    string       `json:"label"`    // 节点名称
	Children []TreeSelect `json:"children"` // 子节点
}

// SysMenuTreeSelect 使用给定的 SysMenu 对象解析为 TreeSelect 对象
func SysMenuTreeSelect(sysMenu model.SysMenu) TreeSelect {
	t := TreeSelect{}
	t.ID = sysMenu.ID
	t.Label = sysMenu.MenuName
	t.Children = []TreeSelect{}
	if len(sysMenu.Children) > 0 {
		for _, menu := range sysMenu.Children {
			child := SysMenuTreeSelect(menu)
			t.Children = append(t.Children, child)
		}
	}
	return t
}

// SysDeptTreeSelect 使用给定的 SysDept 对象解析为 TreeSelect 对象
func SysDeptTreeSelect(sysDept model.SysDept) TreeSelect {
	t := TreeSelect{}
	t.ID = sysDept.ID
	t.Label = sysDept.DeptName
	t.Children = []TreeSelect{}
	if len(sysDept.Children) > 0 {
		for _, dept := range sysDept.Children {
			child := SysDeptTreeSelect(dept)
			t.Children = append(t.Children, child)
		}
	}
	return t
}
