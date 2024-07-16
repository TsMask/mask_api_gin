package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// NewSysConfig 实例化控制层
var NewSysConfig = &SysConfigController{
	sysConfigService: service.NewSysConfig,
}

// SysConfigController 参数配置信息 控制层处理
//
// PATH /system/config
type SysConfigController struct {
	sysConfigService service.ISysConfigService // 参数配置服务
}

// List 参数配置列表
//
// GET /list
func (s *SysConfigController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	data := s.sysConfigService.FindByPage(query)
	c.JSON(200, result.Ok(data))
}

// Info 参数配置信息
//
// GET /:configId
func (s *SysConfigController) Info(c *gin.Context) {
	configId := c.Param("configId")
	if configId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysConfigService.FindById(configId)
	if data.ConfigID == configId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Add 参数配置新增
//
// POST /
func (s *SysConfigController) Add(c *gin.Context) {
	var body model.SysConfig
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.ConfigID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueByKey(body.ConfigKey, "")
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置新增【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysConfigService.Insert(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Edit 参数配置修改
//
// PUT /
func (s *SysConfigController) Edit(c *gin.Context) {
	var body model.SysConfig
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.ConfigID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	config := s.sysConfigService.FindById(body.ConfigID)
	if config.ConfigID != body.ConfigID {
		c.JSON(200, result.ErrMsg("没有权限访问参数配置数据！"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueByKey(body.ConfigKey, body.ConfigID)
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置修改【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysConfigService.Update(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 参数配置删除
//
// DELETE /:configIds
func (s *SysConfigController) Remove(c *gin.Context) {
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
	rows, err := s.sysConfigService.DeleteByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// RefreshCache 参数配置刷新缓存
//
// PUT /refreshCache
func (s *SysConfigController) RefreshCache(c *gin.Context) {
	s.sysConfigService.CacheClean("*")
	s.sysConfigService.CacheLoad("*")
	c.JSON(200, result.Ok(nil))
}

// ConfigKey 参数配置根据参数键名
//
// GET /configKey/:configKey
func (s *SysConfigController) ConfigKey(c *gin.Context) {
	configKey := c.Param("configKey")
	if configKey == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	key := s.sysConfigService.FindValueByKey(configKey)
	if key != "" {
		c.JSON(200, result.OkData(key))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Export 导出参数配置信息
//
// POST /export
func (s *SysConfigController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := ctx.BodyJSONMap(c)
	data := s.sysConfigService.FindByPage(query)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysConfig)

	// 导出文件名称
	fileName := fmt.Sprintf("config_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "参数编号",
		"B1": "参数名称",
		"C1": "参数键名",
		"D1": "参数键值",
		"E1": "系统内置",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		typeValue := "否"
		if row.ConfigType == "Y" {
			typeValue = "是"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.ConfigID,
			"B" + idx: row.ConfigName,
			"C" + idx: row.ConfigKey,
			"D" + idx: row.ConfigValue,
			"E" + idx: typeValue,
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
