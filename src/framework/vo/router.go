package vo

// Router 路由信息对象
type Router struct {
	Name      string     `json:"name"`               // 路由名字 英文首字母大写
	Path      string     `json:"path"`               // 路由地址
	Meta      RouterMeta `json:"meta"`               // 其他元素
	Component string     `json:"component"`          // 组件地址
	Redirect  string     `json:"redirect"`           // 重定向地址
	Children  []Router   `json:"children,omitempty"` // 子路由
}
