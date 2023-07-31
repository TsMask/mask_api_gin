package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
)

// 实例化控制层 SysConfigController 结构体
var NewSysConfig = &SysConfigController{
	sysConfigService: service.NewSysConfigImpl,
}

// 参数配置信息
//
// PATH /system/config
type SysConfigController struct {
	// 参数配置服务
	sysConfigService service.ISysConfig
}

// 参数配置列表
//
// GET /list
func (s *SysConfigController) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysConfigService.SelectConfigPage(querys)
	c.JSON(200, result.Ok(data))
}

// 参数配置信息
//
// GET /:configId
func (s *SysConfigController) Info(c *gin.Context) {
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
func (s *SysConfigController) Add(c *gin.Context) {
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
func (s *SysConfigController) Edit(c *gin.Context) {
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
		c.JSON(200, result.ErrMsg("没有权限访问参数配置数据！"))
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
func (s *SysConfigController) RefreshCache(c *gin.Context) {
	s.sysConfigService.ResetConfigCache()
	c.JSON(200, result.Ok(nil))
}

// 参数配置根据参数键名
//
// GET /configKey/:configKey
func (s *SysConfigController) ConfigKey(c *gin.Context) {
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
func (s *SysConfigController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMapString(c)
	data := s.sysConfigService.SelectConfigPage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("config_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表
	sheet := "Sheet1"
	index, err := file.NewSheet(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	// 设置名为 Sheet1 工作表上 A 到 H 列的宽度为 20
	file.SetColWidth("Sheet1", "A", "H", 20)
	// 设置单元格的值
	file.SetCellValue(sheet, "A1", "参数编号")
	file.SetCellValue(sheet, "B1", "参数名称")
	file.SetCellValue(sheet, "C1", "参数键名")
	file.SetCellValue(sheet, "D1", "参数键值")
	file.SetCellValue(sheet, "E1", "系统内置")

	for i, row := range data["rows"].([]model.SysConfig) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.ConfigID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.ConfigName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.ConfigKey)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.ConfigValue)
		if row.ConfigType == "Y" {
			file.SetCellValue(sheet, "E"+strconv.Itoa(idx), "是")
		} else {
			file.SetCellValue(sheet, "E"+strconv.Itoa(idx), "否")
		}
	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
