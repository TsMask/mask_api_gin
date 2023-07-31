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

// 实例化控制层 SysDictDataController 结构体
var NewSysDictData = &SysDictDataController{
	sysDictDataService: service.NewSysDictDataImpl,
	sysDictTypeService: service.NewSysDictTypeImpl,
}

// 字典类型对应的字典数据信息
//
// PATH /system/dict/data
type SysDictDataController struct {
	// 字典数据服务
	sysDictDataService service.ISysDictData
	// 字典类型服务
	sysDictTypeService service.ISysDictType
}

// 字典数据列表
//
// GET /list
func (s *SysDictDataController) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysDictDataService.SelectDictDataPage(querys)
	c.JSON(200, result.Ok(data))
}

// 字典数据详情
//
// GET /:dictCode
func (s *SysDictDataController) Info(c *gin.Context) {
	dictCode := c.Param("dictCode")
	if dictCode == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysDictDataService.SelectDictDataByCode(dictCode)
	if data.DictCode == dictCode {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 字典数据新增
//
// POST /
func (s *SysDictDataController) Add(c *gin.Context) {
	var body model.SysDictData
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DictCode != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查字典类型是否存在
	sysDictType := s.sysDictTypeService.SelectDictTypeByType(body.DictType)
	if sysDictType.DictType != body.DictType {
		c.JSON(200, result.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典标签唯一
	uniqueDictLabel := s.sysDictDataService.CheckUniqueDictLabel(body.DictType, body.DictLabel, "")
	if !uniqueDictLabel {
		msg := fmt.Sprintf("数据新增【%s】失败，该字典类型下标签名已存在", body.DictLabel)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查字典键值唯一
	uniqueDictValue := s.sysDictDataService.CheckUniqueDictValue(body.DictType, body.DictValue, "")
	if !uniqueDictValue {
		msg := fmt.Sprintf("数据新增【%s】失败，该字典类型下标签值已存在", body.DictValue)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysDictDataService.InsertDictData(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 字典类型修改
//
// PUT /
func (s *SysDictDataController) Edit(c *gin.Context) {
	var body model.SysDictData
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DictCode == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查字典类型是否存在
	sysDictType := s.sysDictTypeService.SelectDictTypeByType(body.DictType)
	if sysDictType.DictType != body.DictType {
		c.JSON(200, result.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典编码是否存在
	SysDictDataController := s.sysDictDataService.SelectDictDataByCode(body.DictCode)
	if SysDictDataController.DictCode != body.DictCode {
		c.JSON(200, result.ErrMsg("没有权限访问字典编码数据！"))
		return
	}

	// 检查字典标签唯一
	uniqueDictLabel := s.sysDictDataService.CheckUniqueDictLabel(body.DictType, body.DictLabel, body.DictCode)
	if !uniqueDictLabel {
		msg := fmt.Sprintf("数据修改【%s】失败，该字典类型下标签名已存在", body.DictLabel)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查字典键值唯一
	uniqueDictValue := s.sysDictDataService.CheckUniqueDictValue(body.DictType, body.DictValue, body.DictCode)
	if !uniqueDictValue {
		msg := fmt.Sprintf("数据修改【%s】失败，该字典类型下标签值已存在", body.DictValue)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysDictDataService.UpdateDictData(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 字典数据删除
//
// DELETE /:dictCodes
func (s *SysDictDataController) Remove(c *gin.Context) {
	dictCodes := c.Param("dictCodes")
	if dictCodes == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(dictCodes, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows, err := s.sysDictDataService.DeleteDictDataByCodes(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 字典数据列表（指定字典类型）
//
// GET /type/:dictType
func (s *SysDictDataController) DictType(c *gin.Context) {
	dictType := c.Param("dictType")
	if dictType == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	data := s.sysDictDataService.SelectDictDataByType(dictType)
	c.JSON(200, result.OkData(data))
}

// 字典数据列表导出
//
// POST /export
func (s *SysDictDataController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	data := s.sysDictDataService.SelectDictDataPage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("dict_data_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
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
	file.SetCellValue(sheet, "A1", "字典编码")
	file.SetCellValue(sheet, "B1", "字典排序")
	file.SetCellValue(sheet, "C1", "字典标签")
	file.SetCellValue(sheet, "D1", "字典键值")
	file.SetCellValue(sheet, "E1", "字典类型")
	file.SetCellValue(sheet, "F1", "状态")

	for i, row := range data["rows"].([]model.SysDictData) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.DictCode)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.DictSort)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.DictLabel)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.DictValue)
		file.SetCellValue(sheet, "E"+strconv.Itoa(idx), row.DictType)
		if row.Status == "0" {
			file.SetCellValue(sheet, "F"+strconv.Itoa(idx), "停用")
		} else {
			file.SetCellValue(sheet, "F"+strconv.Itoa(idx), "正常")
		}
	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
