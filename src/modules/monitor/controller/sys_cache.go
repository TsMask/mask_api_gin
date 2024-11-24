package controller

import (
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/modules/monitor/model"

	"github.com/gin-gonic/gin"
)

// NewSysCache 实例化控制层
var NewSysCache = &SysCacheController{}

// SysCacheController 缓存信息 控制层处理
//
// PATH /monitor/cache
type SysCacheController struct{}

// Info Redis信息
//
// GET /
func (s SysCacheController) Info(c *gin.Context) {
	c.JSON(200, response.OkData(map[string]any{
		"info":         redis.Info(""),
		"dbSize":       redis.KeySize(""),
		"commandStats": redis.CommandStats(""),
	}))
}

// Names 缓存名称列表
//
// GET /names
func (s SysCacheController) Names(c *gin.Context) {
	caches := []model.SysCache{
		model.NewSysCacheNames("用户信息", constCacheKey.LOGIN_TOKEN_KEY),
		model.NewSysCacheNames("配置信息", constCacheKey.SYS_CONFIG_KEY),
		model.NewSysCacheNames("数据字典", constCacheKey.SYS_DICT_KEY),
		model.NewSysCacheNames("验证码", constCacheKey.CAPTCHA_CODE_KEY),
		model.NewSysCacheNames("防重提交", constCacheKey.REPEAT_SUBMIT_KEY),
		model.NewSysCacheNames("限流处理", constCacheKey.RATE_LIMIT_KEY),
		model.NewSysCacheNames("密码错误次数", constCacheKey.PWD_ERR_COUNT_KEY),
	}
	c.JSON(200, response.OkData(caches))
}

// Keys 缓存名称下键名列表
//
// GET /keys?cacheName=xxx
func (s SysCacheController) Keys(c *gin.Context) {
	var query struct {
		CacheName string `form:"cacheName" binding:"required"` // 键名列表中得到的缓存名称
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	caches := []model.SysCache{}
	cacheKeys, _ := redis.GetKeys("", query.CacheName+":*")
	for _, key := range cacheKeys {
		caches = append(caches, model.NewSysCacheKeys(query.CacheName, key))
	}

	c.JSON(200, response.OkData(caches))
}

// Value 缓存内容信息
//
// GET /value?cacheName=xxx&cacheKey=xxx
func (s SysCacheController) Value(c *gin.Context) {
	var query struct {
		CacheName string `form:"cacheName" binding:"required"` // 键名列表中得到的缓存名称
		CacheKey  string `form:"cacheKey" binding:"required"`  // 键名列表中得到的缓存键名
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	cacheValue, err := redis.Get("", query.CacheName+":"+query.CacheKey)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	sysCache := model.NewSysCacheValue(query.CacheName, query.CacheKey, cacheValue)
	c.JSON(200, response.OkData(sysCache))
}

// CleanNames 缓存名称列表安全删除
//
// DELETE /clean/names
func (s SysCacheController) CleanNames(c *gin.Context) {
	caches := []model.SysCache{
		model.NewSysCacheNames("配置信息", constCacheKey.SYS_CONFIG_KEY),
		model.NewSysCacheNames("数据字典", constCacheKey.SYS_DICT_KEY),
		model.NewSysCacheNames("验证码", constCacheKey.CAPTCHA_CODE_KEY),
		model.NewSysCacheNames("防重提交", constCacheKey.REPEAT_SUBMIT_KEY),
		model.NewSysCacheNames("限流处理", constCacheKey.RATE_LIMIT_KEY),
		model.NewSysCacheNames("密码错误次数", constCacheKey.PWD_ERR_COUNT_KEY),
	}
	for _, v := range caches {
		cacheKeys, err := redis.GetKeys("", v.CacheName+":*")
		if err != nil {
			continue
		}
		_ = redis.DelKeys("", cacheKeys)
	}
	c.JSON(200, response.Ok(nil))
}

// CleanKeys 缓存名称下键名删除
//
// DELETE /clean/keys?cacheName=xxx
func (s SysCacheController) CleanKeys(c *gin.Context) {
	var query struct {
		CacheName string `form:"cacheName" binding:"required"` // 键名列表中得到的缓存名称
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}
	if constCacheKey.LOGIN_TOKEN_KEY == query.CacheName {
		c.JSON(200, response.ErrMsg("不能删除用户信息缓存"))
		return
	}

	cacheKeys, err := redis.GetKeys("", query.CacheName+":*")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	if err = redis.DelKeys("", cacheKeys); err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, response.Ok(nil))
}

// CleanValue 缓存内容删除
//
// DELETE /value?cacheName=xxx&cacheKey=xxx
func (s SysCacheController) CleanValue(c *gin.Context) {
	var query struct {
		CacheName string `form:"cacheName" binding:"required"` // 键名列表中得到的缓存名称
		CacheKey  string `form:"cacheKey" binding:"required"`  // 键名列表中得到的缓存键名
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	if err := redis.Del("", query.CacheName+":"+query.CacheKey); err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, response.Ok(nil))
}
