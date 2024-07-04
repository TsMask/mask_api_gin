package controller

import (
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"

	"github.com/gin-gonic/gin"
)

// NewSysCache 实例化控制层
var NewSysCache = &SysCacheController{}

// SysCacheController 缓存监控信息 控制层处理
//
// PATH /monitor/cache
type SysCacheController struct{}

// Info Redis信息
//
// GET /
func (s *SysCacheController) Info(c *gin.Context) {
	c.JSON(200, result.OkData(map[string]any{
		"info":         redis.Info(""),
		"dbSize":       redis.KeySize(""),
		"commandStats": redis.CommandStats(""),
	}))
}

// Names 缓存名称列表
//
// GET /getNames
func (s *SysCacheController) Names(c *gin.Context) {
	caches := []model.SysCache{
		model.NewSysCacheNames("用户信息", constCacheKey.LoginTokenKey),
		model.NewSysCacheNames("配置信息", constCacheKey.SysConfigKey),
		model.NewSysCacheNames("数据字典", constCacheKey.SysDictKey),
		model.NewSysCacheNames("验证码", constCacheKey.CaptchaCodeKey),
		model.NewSysCacheNames("防重提交", constCacheKey.RepeatSubmitKey),
		model.NewSysCacheNames("限流处理", constCacheKey.RateLimitKey),
		model.NewSysCacheNames("密码错误次数", constCacheKey.PwdErrCntKey),
	}
	c.JSON(200, result.OkData(caches))
}

// Keys 缓存名称下键名列表
//
// GET /getKeys/:cacheName
func (s *SysCacheController) Keys(c *gin.Context) {
	cacheName := c.Param("cacheName")
	if cacheName == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	caches := []model.SysCache{}

	// 遍历组装
	cacheKeys, _ := redis.GetKeys("", cacheName+":*")
	for _, key := range cacheKeys {
		caches = append(caches, model.NewSysCacheKeys(cacheName, key))
	}

	c.JSON(200, result.OkData(caches))
}

// Value 缓存内容
//
// GET /getValue/:cacheName/:cacheKey
func (s *SysCacheController) Value(c *gin.Context) {
	cacheName := c.Param("cacheName")
	cacheKey := c.Param("cacheKey")
	if cacheName == "" || cacheKey == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	cacheValue, err := redis.Get("", cacheName+":"+cacheKey)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	sysCache := model.NewSysCacheValue(cacheName, cacheKey, cacheValue)
	c.JSON(200, result.OkData(sysCache))
}

// CleanCacheName 删除缓存名称下键名列表
//
// DELETE /cleanCacheName/:cacheName
func (s *SysCacheController) CleanCacheName(c *gin.Context) {
	cacheName := c.Param("cacheName")
	if cacheName == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	cacheKeys, err := redis.GetKeys("", cacheName+":*")
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	err = redis.DelKeys("", cacheKeys)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// CleanCacheKey 删除缓存键名
//
// DELETE /cleanCacheKey/:cacheName/:cacheKey
func (s *SysCacheController) CleanCacheKey(c *gin.Context) {
	cacheName := c.Param("cacheName")
	cacheKey := c.Param("cacheKey")
	if cacheName == "" || cacheKey == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	err := redis.Del("", cacheName+":"+cacheKey)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// CleanCacheSafe 安全清理缓存名称
//
// DELETE /cleanCacheSafe
func (s *SysCacheController) CleanCacheSafe(c *gin.Context) {
	caches := []model.SysCache{
		model.NewSysCacheNames("配置信息", constCacheKey.SysConfigKey),
		model.NewSysCacheNames("数据字典", constCacheKey.SysDictKey),
		model.NewSysCacheNames("验证码", constCacheKey.CaptchaCodeKey),
		model.NewSysCacheNames("防重提交", constCacheKey.RepeatSubmitKey),
		model.NewSysCacheNames("限流处理", constCacheKey.RateLimitKey),
		model.NewSysCacheNames("密码错误次数", constCacheKey.PwdErrCntKey),
	}
	for _, v := range caches {
		cacheKeys, err := redis.GetKeys("", v.CacheName+":*")
		if err != nil {
			continue
		}
		_ = redis.DelKeys("", cacheKeys)
	}
	c.JSON(200, result.Ok(nil))
}
