package controller

import (
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// 参数配置信息
//
// PATH /system/config
var SysConfig = &sysConfig{
	sysConfigService: service.SysConfigImpl,
}

type sysConfig struct {
	sysConfigService service.ISysConfig
}

// 导出参数配置信息
//
// POST /export
func (s *sysConfig) Export(c *gin.Context) {
	c.JSON(200, result.OkMsg("export"))
}

// 参数配置列表
//
// GET /list
func (s *sysConfig) List(c *gin.Context) {
	// 查询参数转换map
	querys := ctxUtils.QueryMapString(c)
	list := s.sysConfigService.SelectConfigPage(querys)
	c.JSON(200, result.Ok(list))
}

// 参数配置根据参数键名
//
// GET /configKey/:configKey
func (s *sysConfig) ConfigKey(c *gin.Context) {
	configKey := c.Param("configKey")
	if configKey == "" {
		c.JSON(200, result.Err(nil))
		return
	}
	key := s.sysConfigService.SelectConfigValueByKey(configKey)
	if key != "" {
		c.JSON(200, result.OkData(key))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 参数配置信息
//
// GET /:configId
func (s *sysConfig) Info(c *gin.Context) {
	configId := c.Param("configId")
	if configId == "" {
		c.JSON(200, result.Err(nil))
		return
	}
	data := s.sysConfigService.SelectConfigById(configId)
	if data.ConfigID == configId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 参数配置新增
//
// POST /
func (s *sysConfig) Add(c *gin.Context) {
	var config model.SysConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	//
	if config.ConfigID != "" || config.ConfigName == "" || config.ConfigKey == "" || config.ConfigValue == "" {
		c.JSON(200, result.Err(nil))
		return
	}
	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueConfigKey(config)
	if !uniqueConfigKey {
		c.JSON(200, result.ErrMsg("参数配置新增【"+config.ConfigKey+"】失败，参数键名已存在"))
		return
	}

	config.CreateBy = "oks"
	insertId := s.sysConfigService.InsertConfig(config)

	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 参数配置修改
//
// PUT /
func (s *sysConfig) Edit(c *gin.Context) {
	var config model.SysConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	//
	if config.ConfigID == "" || config.ConfigName == "" || config.ConfigKey == "" || config.ConfigValue == "" {
		c.JSON(200, result.Err(nil))
		return
	}
	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueConfigKey(config)
	if !uniqueConfigKey {
		c.JSON(200, result.ErrMsg("参数配置修改【"+config.ConfigKey+"】失败，参数键名已存在"))
		return
	}

	// 检查是否存在
	configInfo := s.sysConfigService.SelectConfigById(config.ConfigID)
	if configInfo.ConfigID != config.ConfigID {
		c.JSON(200, result.OkMsg("没有权限访问参数配置数据！"))
		return
	}

	config.UpdateBy = "oks"
	rows := s.sysConfigService.UpdateConfig(config)

	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 参数配置删除
//
// DELETE /:configIds
func (s *sysConfig) Remove(c *gin.Context) {
	configIds := c.Param("configIds")
	if configIds == "" {
		c.JSON(200, result.Err(nil))
		return
	}
	// 处理字符转id数组
	ids := strings.Split(configIds, ",")
	if len(ids) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	// 去重id
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows := s.sysConfigService.DeleteConfigByIds(uniqueIDs)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 参数配置刷新缓存
//
// PUT /refreshCache
func (s *sysConfig) RefreshCache(c *gin.Context) {
	s.sysConfigService.ResetConfigCache()
	c.JSON(200, result.Ok(nil))
}
