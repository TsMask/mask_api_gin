package controller

import (
	"fmt"
	"mask_api_gin/src/framework/constants/common"
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

// 字典类型信息
//
// PATH /system/dict/type
var SysDictType = &sysDictType{
	sysDictTypeService: service.SysDictTypeImpl,
}

type sysDictType struct {
	// 字典类型服务
	sysDictTypeService service.ISysDictType
}

// 字典类型列表
//
// GET /list
func (s *sysDictType) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysDictTypeService.SelectDictTypePage(querys)
	c.JSON(200, result.Ok(data))
}

// 字典类型信息
//
// GET /:dictId
func (s *sysDictType) Info(c *gin.Context) {
	dictId := c.Param("dictId")
	if dictId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysDictTypeService.SelectDictTypeByID(dictId)
	if data.DictID == dictId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 字典类型新增
//
// POST /
func (s *sysDictType) Add(c *gin.Context) {
	var body model.SysDictType
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DictID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查字典名称唯一
	uniqueDictName := s.sysDictTypeService.CheckUniqueDictName(body.DictName, "")
	if !uniqueDictName {
		msg := fmt.Sprintf("字典新增【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueDictType := s.sysDictTypeService.CheckUniqueDictType(body.DictType, "")
	if !uniqueDictType {
		msg := fmt.Sprintf("字典新增【%s】失败，字典类型已存在", body.DictType)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysDictTypeService.InsertDictType(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 字典类型修改
//
// PUT /
func (s *sysDictType) Edit(c *gin.Context) {
	var body model.SysDictType
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DictID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查数据是否存在
	dictInfo := s.sysDictTypeService.SelectDictTypeByID(body.DictID)
	if dictInfo.DictID != body.DictID {
		c.JSON(200, result.OkMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典名称唯一
	uniqueDictName := s.sysDictTypeService.CheckUniqueDictName(body.DictName, body.DictID)
	if !uniqueDictName {
		msg := fmt.Sprintf("字典修改【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueDictType := s.sysDictTypeService.CheckUniqueDictType(body.DictType, body.DictID)
	if !uniqueDictType {
		msg := fmt.Sprintf("字典修改【%s】失败，字典类型已存在", body.DictType)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysDictTypeService.UpdateDictType(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 部门删除
//
// DELETE /:dictIds
func (s *sysDictType) Remove(c *gin.Context) {
	dictIds := c.Param("dictIds")
	if dictIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(dictIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows, err := s.sysDictTypeService.DeleteDictTypeByIDs(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 字典类型刷新缓存
//
// PUT /refreshCache
func (s *sysDictType) RefreshCache(c *gin.Context) {
	s.sysDictTypeService.ResetDictCache()
	c.JSON(200, result.Ok(nil))
}

// 字典类型选择框列表
//
// GET /getDictOptionselect
func (s *sysDictType) DictOptionselect(c *gin.Context) {
	data := s.sysDictTypeService.SelectDictTypeList(model.SysDictType{
		Status: common.STATUS_YES,
	})

	type labelValue struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	// 数据组
	arr := []labelValue{}
	for _, v := range data {
		arr = append(arr, labelValue{
			Label: v.DictName,
			Value: v.DictType,
		})
	}
	c.JSON(200, result.OkData(arr))
}

// 字典类型列表导出
//
// POST /export
func (s *sysDictType) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	data := s.sysDictTypeService.SelectDictTypePage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("dict_type_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
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
	file.SetCellValue(sheet, "A1", "字典主键")
	file.SetCellValue(sheet, "B1", "字典名称")
	file.SetCellValue(sheet, "C1", "字典类型")
	file.SetCellValue(sheet, "D1", "状态")

	for i, row := range data["rows"].([]model.SysDictType) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.DictID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.DictName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.DictType)
		if row.Status == "0" {
			file.SetCellValue(sheet, "D"+strconv.Itoa(idx), "停用")
		} else {
			file.SetCellValue(sheet, "D"+strconv.Itoa(idx), "正常")
		}
	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
