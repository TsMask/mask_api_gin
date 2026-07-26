package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
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
	query := reqctx.QueryMap(c)
	rows, total := s.sysConfigService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 参数配置信息
//
// GET /:id
func (s SysConfigController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.sysConfigService.FindById(id)
	if data.ID == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 参数配置新增
//
// POST /
func (s SysConfigController) Add(c *gin.Context) {
	var body model.SysConfig
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id not is empty"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueByKey(body.ConfigKey, 0)
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置新增【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysConfigService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 参数配置修改
//
// PUT /
func (s SysConfigController) Edit(c *gin.Context) {
	var body model.SysConfig
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 检查是否存在
	configInfo := s.sysConfigService.FindById(body.ID)
	if configInfo.ID != body.ID {
		c.JSON(200, resp.ErrMsg("没有权限访问参数配置数据！"))
		return
	}

	// 检查属性值唯一
	uniqueConfigKey := s.sysConfigService.CheckUniqueByKey(body.ConfigKey, body.ID)
	if !uniqueConfigKey {
		msg := fmt.Sprintf("参数配置修改【%s】失败，参数键名已存在", body.ConfigKey)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	configInfo.ConfigValueType = body.ConfigValueType
	configInfo.ConfigGroup = body.ConfigGroup
	configInfo.ConfigSort = body.ConfigSort
	configInfo.ConfigStatus = body.ConfigStatus
	configInfo.ConfigName = body.ConfigName
	configInfo.ConfigKey = body.ConfigKey
	configInfo.ConfigValue = body.ConfigValue
	configInfo.Remark = body.Remark
	configInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysConfigService.Update(configInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 参数配置删除
//
// DELETE /:id
func (s SysConfigController) Remove(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(idStr, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	rows, err := s.sysConfigService.DeleteByIds(ids)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// Refresh 参数配置刷新缓存
//
// PUT /refresh
func (s SysConfigController) Refresh(c *gin.Context) {
	s.sysConfigService.CacheClean("*")
	s.sysConfigService.CacheLoad("*")
	c.JSON(200, resp.Ok(nil))
}

// ConfigKey 参数配置根据参数键名
//
// GET /config-key/:config_key
func (s SysConfigController) ConfigKey(c *gin.Context) {
	configKey := c.Param("config_key")
	if configKey == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: configKey is empty"))
		return
	}
	key := s.sysConfigService.FindValueByKey(configKey)
	if key != "" {
		c.JSON(200, resp.OkData(key))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Export 导出参数配置信息
//
// GET /export
func (s SysConfigController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	rows, total := s.sysConfigService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("config_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysConfigService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
