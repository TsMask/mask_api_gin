package config

// 类型化配置访问器：与 Get 行为统一，键缺失返回默认值，类型不符时安全强转，
// 避免 Get(key).(T) 在配置写错时 panic。
//
//	GetBoolOr(key, def)           // bool
//	GetIntOr(key, def)            // int
//	GetStringOr(key, def)         // string
//	GetStringSliceOr(key, def)    // []string
//	GetStringMapOr(key, def)      // map[string]any

// GetBoolOr 读取布尔配置；键缺失返回 def。
func GetBoolOr(key string, def bool) bool {
	if v := conf.Get(key); v == nil {
		return def
	}
	return conf.GetBool(key)
}

// GetIntOr 读取整数配置；键缺失返回 def。
func GetIntOr(key string, def int) int {
	if v := conf.Get(key); v == nil {
		return def
	}
	return conf.GetInt(key)
}

// GetStringOr 读取字符串配置；键缺失返回 def。
func GetStringOr(key string, def string) string {
	if v := conf.Get(key); v == nil {
		return def
	}
	return conf.GetString(key)
}

// GetStringSliceOr 读取字符串切片配置；键缺失返回 def。
func GetStringSliceOr(key string, def []string) []string {
	if v := conf.Get(key); v == nil {
		return def
	}
	return conf.GetStringSlice(key)
}

// GetStringMapOr 读取字符串映射配置（string -> any，支持嵌套）；键缺失返回 def。
func GetStringMapOr(key string, def map[string]any) map[string]any {
	if v := conf.Get(key); v == nil {
		return def
	}
	return conf.GetStringMap(key)
}
