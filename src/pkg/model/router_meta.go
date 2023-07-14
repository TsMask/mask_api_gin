package model

// RouterMeta 路由元信息对象
type RouterMeta struct {
	// 设置该菜单在侧边栏和面包屑中展示的名字
	Title string `json:"title"`
	// 设置该菜单的图标
	Icon string `json:"icon"`
	// 设置为true，则不会被 <keep-alive>缓存
	Cache bool `json:"cache"`
	// 内链地址（http(s)://开头）, 打开目标位置 '_blank' | '_self' | null | undefined
	Target string `json:"target"`
	// 在菜单中隐藏自己和子节点
	Hide bool `json:"hide"`
}
