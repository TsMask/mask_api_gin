package vo

// AuthMenu 当前登录用户可用的语义化菜单。
type AuthMenu struct {
	ID        int64      `json:"id"`                  // 菜单ID
	Kind      string     `json:"kind"`                // 菜单类型 directory/page/iframe/external
	Title     string     `json:"title"`               // 展示名称
	Path      string     `json:"path"`                // 完整访问路径或网络地址
	Icon      string     `json:"icon,omitempty"`      // 菜单图标
	Hidden    bool       `json:"hidden,omitempty"`    // 是否在菜单中隐藏
	Cache     bool       `json:"cache,omitempty"`     // 页面是否缓存
	Target    string     `json:"target,omitempty"`    // 外链打开方式
	Component string     `json:"component,omitempty"` // 前端页面组件路径
	Children  []AuthMenu `json:"children,omitempty"`  // 子菜单
}
