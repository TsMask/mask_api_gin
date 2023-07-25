package controller

import (
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"

	"github.com/gin-gonic/gin"
)

// 缓存监控信息
//
// PATH /monitor/cache
var SysCache = &sysCache{}

type sysCache struct{}

// Redis信息
//
// GET /
func (s *sysCache) Info(c *gin.Context) {
	c.JSON(200, result.OkData(map[string]interface{}{
		"info":         redis.Info(),
		"dbSize":       redis.KeySize(),
		"commandStats": redis.CommandStats(),
	}))
}

// 缓存名称列表
//
// GET /getNames
func (s *sysCache) Names(c *gin.Context) {
	caches := []model.SysCache{
		model.NewSysCacheNames("用户信息", cachekey.LOGIN_TOKEN_KEY),
		model.NewSysCacheNames("配置信息", cachekey.SYS_CONFIG_KEY),
		model.NewSysCacheNames("数据字典", cachekey.SYS_DICT_KEY),
		model.NewSysCacheNames("验证码", cachekey.CAPTCHA_CODE_KEY),
		model.NewSysCacheNames("防重提交", cachekey.REPEAT_SUBMIT_KEY),
		model.NewSysCacheNames("限流处理", cachekey.RATE_LIMIT_KEY),
		model.NewSysCacheNames("密码错误次数", cachekey.PWD_ERR_CNT_KEY),
	}
	c.JSON(200, result.OkData(caches))
}

// 缓存名称下键名列表
//
// GET /getKeys/:cacheName
func (s *sysCache) Keys(c *gin.Context) {
	cacheName := c.Param("cacheName")
	if cacheName == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	caches := []model.SysCache{}

	// 遍历组装
	cacheKeys := redis.GetKeys(cacheName + ":*")
	for _, key := range cacheKeys {
		caches = append(caches, model.NewSysCacheKeys(cacheName, key))
	}

	c.JSON(200, result.OkData(caches))
}

// 缓存内容
//
// GET /getValue/:cacheName/:cacheKey
func (s *sysCache) Value(c *gin.Context) {
	cacheName := c.Param("cacheName")
	cacheKey := c.Param("cacheKey")
	if cacheName == "" || cacheKey == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	cacheValue := redis.Get(cacheName + ":" + cacheKey)
	sysCache := model.NewSysCacheValue(cacheName, cacheKey, cacheValue)
	c.JSON(200, result.OkData(sysCache))
}

// 删除缓存名称下键名列表
//
// DELETE /clearCacheName/:cacheName
func (s *sysCache) ClearCacheName(c *gin.Context) {
	cacheName := c.Param("cacheName")
	if cacheName == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	cacheKeys := redis.GetKeys(cacheName + ":*")
	ok := redis.DelKeys(cacheKeys)
	if ok {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 删除缓存键名
//
// DELETE /clearCacheKey/:cacheName/:cacheKey
func (s *sysCache) ClearCacheKey(c *gin.Context) {
	cacheName := c.Param("cacheName")
	cacheKey := c.Param("cacheKey")
	if cacheName == "" || cacheKey == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	ok := redis.Del(cacheName + ":" + cacheKey)
	if ok {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 安全清理缓存名称
//
// DELETE /clearCacheSafe
func (s *sysCache) ClearCacheSafe(c *gin.Context) {
	caches := []model.SysCache{
		model.NewSysCacheNames("配置信息", cachekey.SYS_CONFIG_KEY),
		model.NewSysCacheNames("数据字典", cachekey.SYS_DICT_KEY),
		model.NewSysCacheNames("验证码", cachekey.CAPTCHA_CODE_KEY),
		model.NewSysCacheNames("防重提交", cachekey.REPEAT_SUBMIT_KEY),
		model.NewSysCacheNames("限流处理", cachekey.RATE_LIMIT_KEY),
		model.NewSysCacheNames("密码错误次数", cachekey.PWD_ERR_CNT_KEY),
	}
	for _, v := range caches {
		cacheKeys := redis.GetKeys(v.CacheKey + ":*")
		redis.DelKeys(cacheKeys)
	}
	c.JSON(200, result.Ok(nil))
}