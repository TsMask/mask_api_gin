package model

// TreeSelect 树结构实体类
type TreeSelect struct {
	// ID 节点ID
	ID string `json:"id"`

	// Label 节点名称
	Label string `json:"label"`

	// Children 子节点
	Children []TreeSelect `json:"children"`
}
