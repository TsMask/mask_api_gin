package controller

import (
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NewSysConfig 实例化控制层
var NewSysConfig = &SysConfigController{
	sysConfigService: service.NewSysConfig,
}

// SysConfigController 参数配置信息 控制层处理
//
// PATH /system/config
type SysConfigController struct {
	sysConfigService *service.SysConfig // 参数配置服务
}

// List 参数配置列表
//
// GET /list
func (s SysConfigController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	rows, total := s.sysConfigService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 参数配置信息
//
// GET /:configId
func (s SysConfigController) Info(c *gin.Context) {
	configId := c.Param("configId")
	if configId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: configId is empty"))
		return
	}

	configInfo := s.sysConfigService.FindById(configId)
	if configInfo.ConfigId == configId {
		c.JSON(200, response.OkData(configInfo))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 参数配置新增
//
// POST /
func (s SysConfigController) Add(c *gin.Context) {
	var body model.SysConfig
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.ConfigId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: configId not is empty"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueByKey(body.ConfigKey, "")
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置新增【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysConfigService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 参数配置修改
//
// PUT /
func (s SysConfigController) Edit(c *gin.Context) {
	var body model.SysConfig
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.ConfigId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: configId is empty"))
		return
	}

	// 检查是否存在
	configInfo := s.sysConfigService.FindById(body.ConfigId)
	if configInfo.ConfigId != body.ConfigId {
		c.JSON(200, response.ErrMsg("没有权限访问参数配置数据！"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueByKey(body.ConfigKey, body.ConfigId)
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置修改【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	configInfo.ConfigType = body.ConfigType
	configInfo.ConfigName = body.ConfigName
	configInfo.ConfigKey = body.ConfigKey
	configInfo.ConfigValue = body.ConfigValue
	configInfo.Remark = body.Remark
	configInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysConfigService.Update(configInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 参数配置删除
//
// DELETE /:configId
func (s SysConfigController) Remove(c *gin.Context) {
	configIdStr := c.Param("configId")
	configIds := parse.RemoveDuplicatesToArray(configIdStr, ",")
	if configIdStr == "" || len(configIds) == 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: configId is empty"))
		return
	}

	rows, err := s.sysConfigService.DeleteByIds(configIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}

// Refresh 参数配置刷新缓存
//
// PUT /refresh
func (s SysConfigController) Refresh(c *gin.Context) {
	s.sysConfigService.CacheClean("*")
	s.sysConfigService.CacheLoad("*")
	c.JSON(200, response.Ok(nil))
}

// ConfigKey 参数配置根据参数键名
//
// GET /config-key/:configKey
func (s SysConfigController) ConfigKey(c *gin.Context) {
	configKey := c.Param("configKey")
	if configKey == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: configKey is empty"))
		return
	}
	key := s.sysConfigService.FindValueByKey(configKey)
	if key != "" {
		c.JSON(200, response.OkData(key))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Export 导出参数配置信息
//
// GET /export
func (s SysConfigController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := ctx.QueryMap(c)
	rows, total := s.sysConfigService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

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
			"A" + idx: row.ConfigId,
			"B" + idx: row.ConfigName,
			"C" + idx: row.ConfigKey,
			"D" + idx: row.ConfigValue,
			"E" + idx: typeValue,
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
