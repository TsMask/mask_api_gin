package vo

// Router 路由信息对象
type Router struct {
	// 路由名字 英文首字母大写
	Name string `json:"name"`
	// 路由地址
	Path string `json:"path"`
	// 其他元素
	Meta RouterMeta `json:"meta"`
	// 组件地址
	Component string `json:"component"`
	// 重定向地址
	Redirect string `json:"redirect"`
	// 子路由
	Children []Router `json:"children"`
}
