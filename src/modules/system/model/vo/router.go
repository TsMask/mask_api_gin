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

// RouterMeta 路由元信息对象
type RouterMeta struct {
	Title           string `json:"title"`           // 设置该菜单在侧边栏和面包屑中展示的名字
	Icon            string `json:"icon"`            // 设置该菜单的图标
	Cache           bool   `json:"cache"`           // 设置为true，则不会被 <keep-alive>缓存
	Target          string `json:"target"`          // 内链地址（http(s)://开头）, 打开目标位置 '_blank' | '_self' | ''
	HideChildInMenu bool   `json:"hideChildInMenu"` // 在菜单中隐藏子节点
	HideInMenu      bool   `json:"hideInMenu"`      // 在菜单中隐藏自己和子节点
}