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
	querys := ctx.QueryMap(c)
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
	querys := ctx.BodyJSONMap(c)
	data := s.sysDictDataService.SelectDictDataPage(querys)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysDictData)

	// 导出文件名称
	fileName := fmt.Sprintf("dict_data_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "字典编码",
		"B1": "字典排序",
		"C1": "字典标签",
		"D1": "字典键值",
		"E1": "字典类型",
		"F1": "状态",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		statusValue := "停用"
		if row.Status == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.DictCode,
			"B" + idx: row.DictSort,
			"C" + idx: row.DictLabel,
			"D" + idx: row.DictValue,
			"E" + idx: row.DictType,
			"F" + idx: statusValue,
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
