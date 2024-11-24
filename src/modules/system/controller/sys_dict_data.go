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

// NewSysDictData 实例化控制层
var NewSysDictData = &SysDictDataController{
	sysDictDataService: service.NewSysDictData,
	sysDictTypeService: service.NewSysDictType,
}

// SysDictDataController 字典类型对应的字典数据信息 控制层处理
//
// PATH /system/dict/data
type SysDictDataController struct {
	sysDictDataService *service.SysDictData // 字典数据服务
	sysDictTypeService *service.SysDictType // 字典类型服务
}

// List 字典数据列表
//
// GET /list
func (s SysDictDataController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	rows, total := s.sysDictDataService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 字典数据详情
//
// GET /:dataId
func (s SysDictDataController) Info(c *gin.Context) {
	dataId := c.Param("dataId")
	if dataId == "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	data := s.sysDictDataService.FindById(dataId)
	if data.DataId == dataId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 字典数据新增
//
// POST /
func (s SysDictDataController) Add(c *gin.Context) {
	var body model.SysDictData
	if err := c.ShouldBindBodyWithJSON(&body); err != nil || body.DataId != "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 检查字典类型是否存在
	sysDictType := s.sysDictTypeService.FindByType(body.DictType)
	if sysDictType.DictType != body.DictType {
		c.JSON(200, response.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典标签唯一
	uniqueDictLabel := s.sysDictDataService.CheckUniqueTypeByLabel(body.DictType, body.DataLabel, "")
	if !uniqueDictLabel {
		msg := fmt.Sprintf("数据新增【%s】失败，该字典类型下标签名已存在", body.DataLabel)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查字典键值唯一
	uniqueDictValue := s.sysDictDataService.CheckUniqueTypeByValue(body.DictType, body.DataValue, "")
	if !uniqueDictValue {
		msg := fmt.Sprintf("数据新增【%s】失败，该字典类型下标签值已存在", body.DataValue)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysDictDataService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 字典类型修改
//
// PUT /
func (s SysDictDataController) Edit(c *gin.Context) {
	var body model.SysDictData
	if err := c.ShouldBindBodyWithJSON(&body); err != nil || body.DataId == "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 检查字典类型是否存在
	sysDictType := s.sysDictTypeService.FindByType(body.DictType)
	if sysDictType.DictType != body.DictType {
		c.JSON(200, response.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典编码是否存在
	dataInfo := s.sysDictDataService.FindById(body.DataId)
	if dataInfo.DataId != body.DataId {
		c.JSON(200, response.ErrMsg("没有权限访问字典编码数据！"))
		return
	}

	// 检查字典标签唯一
	uniqueDictLabel := s.sysDictDataService.CheckUniqueTypeByLabel(body.DictType, body.DataLabel, body.DataId)
	if !uniqueDictLabel {
		msg := fmt.Sprintf("数据修改【%s】失败，该字典类型下标签名已存在", body.DataLabel)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查字典键值唯一
	uniqueDictValue := s.sysDictDataService.CheckUniqueTypeByValue(body.DictType, body.DataValue, body.DataId)
	if !uniqueDictValue {
		msg := fmt.Sprintf("数据修改【%s】失败，该字典类型下标签值已存在", body.DataValue)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	dataInfo.DictType = body.DictType
	dataInfo.DataLabel = body.DataLabel
	dataInfo.DataValue = body.DataValue
	dataInfo.DataSort = body.DataSort
	dataInfo.TagClass = body.TagClass
	dataInfo.TagType = body.TagType
	dataInfo.StatusFlag = body.StatusFlag
	dataInfo.Remark = body.Remark
	dataInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysDictDataService.Update(dataInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 字典数据删除
//
// DELETE /:dataId
func (s SysDictDataController) Remove(c *gin.Context) {
	dataIdsStr := c.Param("dataId")
	dataIds := parse.RemoveDuplicatesToArray(dataIdsStr, ",")
	if dataIdsStr == "" || len(dataIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	rows, err := s.sysDictDataService.DeleteByIds(dataIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}

// DictType 字典数据列表（指定字典类型）
//
// GET /type/:dictType
func (s SysDictDataController) DictType(c *gin.Context) {
	dictType := c.Param("dictType")
	if dictType == "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	data := s.sysDictDataService.FindByType(dictType)
	c.JSON(200, response.OkData(data))
}

// Export 字典数据列表导出
//
// GET /export
func (s SysDictDataController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := ctx.QueryMap(c)
	rows, total := s.sysDictDataService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("dict_data_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "字典类型",
		"B1": "数据排序",
		"C1": "数据编号",
		"D1": "数据标签",
		"E1": "数据键值",
		"F1": "数据状态",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		statusValue := "停用"
		if row.StatusFlag == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.DictType,
			"B" + idx: row.DataSort,
			"C" + idx: row.DataId,
			"D" + idx: row.DataLabel,
			"E" + idx: row.DataValue,
			"F" + idx: statusValue,
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
