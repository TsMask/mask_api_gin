package model

import "strings"

// SysCache 缓存信息对象
type SysCache struct {
	CacheName  string `json:"cacheName"`  // 缓存名称
	CacheKey   string `json:"cacheKey"`   // 缓存键名
	CacheValue string `json:"cacheValue"` // 缓存内容
	Remark     string `json:"remark"`     // 备注
}

// NewNames 创建新的缓存名称列表项实例
func NewNames(cacheName string, cacheKey string) *SysCache {
	return &SysCache{
		CacheName:  cacheKey[:len(cacheKey)-1],
		CacheKey:   "",
		CacheValue: "",
		Remark:     cacheName,
	}
}

// NewKeys 创建新的缓存键名列表项实例
func NewKeys(cacheName string, cacheKey string) *SysCache {
	return &SysCache{
		CacheName:  cacheName,
		CacheKey:   strings.Replace(cacheKey, cacheName+":", "", 1),
		CacheValue: "",
		Remark:     "",
	}
}

// NewValue 创建新的缓存键名内容项实例
func NewValue(cacheName string, cacheKey string, cacheValue string) *SysCache {
	return &SysCache{
		CacheName:  cacheName,
		CacheKey:   cacheKey,
		CacheValue: cacheValue,
		Remark:     "",
	}
}
