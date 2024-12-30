package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
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
	query := reqctx.QueryMap(c)
	rows, total := s.sysDictDataService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 字典数据详情
//
// GET /:dataId
func (s SysDictDataController) Info(c *gin.Context) {
	dataId := parse.Number(c.Param("dataId"))
	if dataId <= 0 {
		c.JSON(400, resp.CodeMsg(40010, "bind err: deptId is empty"))
		return
	}

	data := s.sysDictDataService.FindById(dataId)
	if data.DataId == dataId {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 字典数据新增
//
// POST /
func (s SysDictDataController) Add(c *gin.Context) {
	var body model.SysDictData
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(40422, errMsgs))
		return
	}
	if body.DataId > 0 {
		c.JSON(400, resp.CodeMsg(40010, "bind err: dataId not is empty"))
		return
	}

	// 检查字典类型是否存在
	sysDictType := s.sysDictTypeService.FindByType(body.DictType)
	if sysDictType.DictType != body.DictType {
		c.JSON(200, resp.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典标签唯一
	uniqueDictLabel := s.sysDictDataService.CheckUniqueTypeByLabel(body.DictType, body.DataLabel, 0)
	if !uniqueDictLabel {
		msg := fmt.Sprintf("数据新增【%s】失败，该字典类型下标签名已存在", body.DataLabel)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查字典键值唯一
	uniqueDictValue := s.sysDictDataService.CheckUniqueTypeByValue(body.DictType, body.DataValue, 0)
	if !uniqueDictValue {
		msg := fmt.Sprintf("数据新增【%s】失败，该字典类型下标签值已存在", body.DataValue)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysDictDataService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 字典类型修改
//
// PUT /
func (s SysDictDataController) Edit(c *gin.Context) {
	var body model.SysDictData
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(40422, errMsgs))
		return
	}
	if body.DataId <= 0 {
		c.JSON(400, resp.CodeMsg(40010, "bind err: dataId is empty"))
		return
	}

	// 检查字典类型是否存在
	sysDictType := s.sysDictTypeService.FindByType(body.DictType)
	if sysDictType.DictType != body.DictType {
		c.JSON(200, resp.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典编码是否存在
	dataInfo := s.sysDictDataService.FindById(body.DataId)
	if dataInfo.DataId != body.DataId {
		c.JSON(200, resp.ErrMsg("没有权限访问字典编码数据！"))
		return
	}

	// 检查字典标签唯一
	uniqueDictLabel := s.sysDictDataService.CheckUniqueTypeByLabel(body.DictType, body.DataLabel, body.DataId)
	if !uniqueDictLabel {
		msg := fmt.Sprintf("数据修改【%s】失败，该字典类型下标签名已存在", body.DataLabel)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查字典键值唯一
	uniqueDictValue := s.sysDictDataService.CheckUniqueTypeByValue(body.DictType, body.DataValue, body.DataId)
	if !uniqueDictValue {
		msg := fmt.Sprintf("数据修改【%s】失败，该字典类型下标签值已存在", body.DataValue)
		c.JSON(200, resp.ErrMsg(msg))
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
	dataInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysDictDataService.Update(dataInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 字典数据删除
//
// DELETE /:dataId
func (s SysDictDataController) Remove(c *gin.Context) {
	dataIdsStr := c.Param("dataId")
	arr := parse.RemoveDuplicatesToArray(dataIdsStr, ",")
	if dataIdsStr == "" || len(arr) <= 0 {
		c.JSON(400, resp.CodeMsg(40010, "bind err: dataId is empty"))
		return
	}

	// 转换成int64数组类型
	dataIds := make([]int64, 0)
	for _, v := range arr {
		dataIds = append(dataIds, parse.Number(v))
	}

	rows, err := s.sysDictDataService.DeleteByIds(dataIds)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// DictType 字典数据列表（指定字典类型）
//
// GET /type/:dictType
func (s SysDictDataController) DictType(c *gin.Context) {
	dictType := c.Param("dictType")
	if dictType == "" {
		c.JSON(400, resp.CodeMsg(40010, "bind err: dictType is empty"))
		return
	}

	data := s.sysDictDataService.FindByType(dictType)
	c.JSON(200, resp.OkData(data))
}

// Export 字典数据列表导出
//
// GET /export
func (s SysDictDataController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	rows, total := s.sysDictDataService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("dict_data_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysDictDataService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
