package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

// 参数配置列表
//
// GET /list
func (s *sysConfig) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysConfigService.SelectConfigPage(querys)
	c.JSON(200, result.Ok(data))
}

// 参数配置信息
//
// GET /:configId
func (s *sysConfig) Info(c *gin.Context) {
	configId := c.Param("configId")
	if configId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
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
	var body model.SysConfig
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.ConfigID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueConfigKey(body.ConfigKey, "")
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置新增【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysConfigService.InsertConfig(body)
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
	var body model.SysConfig
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.ConfigID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueConfigKey(body.ConfigKey, body.ConfigID)
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置修改【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查是否存在
	config := s.sysConfigService.SelectConfigById(body.ConfigID)
	if config.ConfigID != body.ConfigID {
		c.JSON(200, result.OkMsg("没有权限访问参数配置数据！"))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysConfigService.UpdateConfig(body)
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
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(configIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows, err := s.sysConfigService.DeleteConfigByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 参数配置刷新缓存
//
// PUT /refreshCache
func (s *sysConfig) RefreshCache(c *gin.Context) {
	s.sysConfigService.ResetConfigCache()
	c.JSON(200, result.Ok(nil))
}

// 参数配置根据参数键名
//
// GET /configKey/:configKey
func (s *sysConfig) ConfigKey(c *gin.Context) {
	configKey := c.Param("configKey")
	if configKey == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	key := s.sysConfigService.SelectConfigValueByKey(configKey)
	if key != "" {
		c.JSON(200, result.OkData(key))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 导出参数配置信息
//
// POST /export
func (s *sysConfig) Export(c *gin.Context) {
	c.JSON(200, result.OkMsg("export"))
}
