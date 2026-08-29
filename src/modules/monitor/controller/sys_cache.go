package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/monitor/model"
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
	c.JSON(200, resp.OkData(map[string]any{
		"info":          cache.Info(""),
		"db_size":       cache.KeySize(""),
		"command_stats": cache.CommandStats(""),
	}))
}

// Names 缓存名称列表
//
// GET /names
func (s SysCacheController) Names(c *gin.Context) {
	caches := []model.SysCache{
		model.NewNames("用户令牌", constants.CACHE_TOKEN_SESSION),
		model.NewNames("用户会话索引", constants.CACHE_TOKEN_USERS),
		model.NewNames("配置信息", constants.CACHE_SYS_CONFIG),
		model.NewNames("数据字典", constants.CACHE_SYS_DICT),
		model.NewNames("用户验证码", constants.CACHE_CAPTCHA_CODE),
		model.NewNames("防重提交", constants.CACHE_REPEAT_SUBMIT),
		model.NewNames("限流处理", constants.CACHE_RATE_LIMIT),
		model.NewNames("密码错误次数", constants.CACHE_PWD_ERR_COUNT),
	}
	c.JSON(200, resp.OkData(caches))
}

// Keys 缓存名称下键名列表
//
// GET /keys?cacheName=xxx
func (s SysCacheController) Keys(c *gin.Context) {
	var query struct {
		CacheName string `form:"cache_name" binding:"required"` // 键名列表中得到的缓存名称
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	caches := []model.SysCache{}
	cacheKeys, _ := cache.GetKeys("", query.CacheName+":*")
	for _, key := range cacheKeys {
		caches = append(caches, model.NewKeys(query.CacheName, key))
	}

	c.JSON(200, resp.OkData(caches))
}

// Value 缓存内容信息
//
// GET /value?cacheName=xxx&cacheKey=xxx
func (s SysCacheController) Value(c *gin.Context) {
	var query struct {
		CacheName string `form:"cache_name" binding:"required"` // 键名列表中得到的缓存名称
		CacheKey  string `form:"cache_key" binding:"required"`  // 键名列表中得到的缓存键名
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	cacheValue, err := cache.Get("", query.CacheName+":"+query.CacheKey)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	// token_sessions 缓存值脱敏：隐藏权限列表和密码等敏感字段
	if query.CacheName == constants.CACHE_TOKEN_SESSION && cacheValue != "" {
		var info map[string]any
		if json.Unmarshal([]byte(cacheValue), &info) == nil {
			delete(info, "permissions")
			if user, ok := info["user"].(map[string]any); ok {
				delete(user, "password")
			}
			if sanitized, err := json.Marshal(info); err == nil {
				cacheValue = string(sanitized)
			}
		}
	}
	sysCache := model.NewValue(query.CacheName, query.CacheKey, cacheValue)
	c.JSON(200, resp.OkData(sysCache))
}

// CleanNames 缓存名称列表安全删除
//
// DELETE /names
func (s SysCacheController) CleanNames(c *gin.Context) {
	caches := []model.SysCache{
		model.NewNames("配置信息", constants.CACHE_SYS_CONFIG),
		model.NewNames("数据字典", constants.CACHE_SYS_DICT),
		model.NewNames("用户验证码", constants.CACHE_CAPTCHA_CODE),
		model.NewNames("防重提交", constants.CACHE_REPEAT_SUBMIT),
		model.NewNames("限流处理", constants.CACHE_RATE_LIMIT),
		model.NewNames("密码错误次数", constants.CACHE_PWD_ERR_COUNT),
	}
	for _, v := range caches {
		cacheKeys, err := cache.GetKeys("", v.CacheName+":*")
		if err != nil {
			continue
		}
		_ = cache.DelKeys("", cacheKeys)
	}
	c.JSON(200, resp.Ok(nil))
}

// CleanKeys 缓存名称下键名删除
//
// DELETE /keys?cacheName=xxx
func (s SysCacheController) CleanKeys(c *gin.Context) {
	var query struct {
		CacheName string `form:"cache_name" binding:"required"` // 键名列表中得到的缓存名称
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	// 禁止从缓存管理界面直接删除
	if query.CacheName == constants.CACHE_TOKEN_SESSION || query.CacheName == constants.CACHE_TOKEN_USERS {
		c.JSON(200, resp.ErrMsg("不能删除令牌信息缓存"))
		return
	}

	cacheKeys, err := cache.GetKeys("", query.CacheName+":*")
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	if err = cache.DelKeys("", cacheKeys); err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, resp.Ok(nil))
}

// CleanValue 缓存内容删除
//
// DELETE /value?cacheName=xxx&cacheKey=xxx
func (s SysCacheController) CleanValue(c *gin.Context) {
	var query struct {
		CacheName string `form:"cache_name" binding:"required"` // 键名列表中得到的缓存名称
		CacheKey  string `form:"cache_key" binding:"required"`  // 键名列表中得到的缓存键名
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	// 禁止从缓存管理界面直接删除
	if query.CacheName == constants.CACHE_TOKEN_SESSION || query.CacheName == constants.CACHE_TOKEN_USERS {
		c.JSON(200, resp.ErrMsg("不能删除令牌信息缓存"))
		return
	}

	err := cache.Del("", query.CacheName+":"+query.CacheKey)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, resp.Ok(nil))
}
