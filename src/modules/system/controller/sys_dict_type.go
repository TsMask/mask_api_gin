package controller

import (
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NewSysDictType 实例化控制层
var NewSysDictType = &SysDictTypeController{
	sysDictTypeService: service.NewSysDictType,
}

// SysDictTypeController 字典类型信息 控制层处理
//
// PATH /system/dict/type
type SysDictTypeController struct {
	sysDictTypeService *service.SysDictType // 字典类型服务
}

// List 字典类型列表
//
// GET /list
func (s SysDictTypeController) List(c *gin.Context) {
	query := context.QueryMap(c)
	rows, total := s.sysDictTypeService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 字典类型信息
//
// GET /:dictId
func (s SysDictTypeController) Info(c *gin.Context) {
	dictId := c.Param("dictId")
	if dictId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: dictId is empty"))
		return
	}

	data := s.sysDictTypeService.FindById(dictId)
	if data.DictId == dictId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 字典类型新增
//
// POST /
func (s SysDictTypeController) Add(c *gin.Context) {
	var body model.SysDictType
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.DictId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: dictId not is empty"))
		return
	}

	// 检查字典名称唯一
	uniqueDictName := s.sysDictTypeService.CheckUniqueByName(body.DictName, "")
	if !uniqueDictName {
		msg := fmt.Sprintf("字典新增【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueDictType := s.sysDictTypeService.CheckUniqueByType(body.DictType, "")
	if !uniqueDictType {
		msg := fmt.Sprintf("字典新增【%s】失败，字典类型已存在", body.DictType)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = context.LoginUserToUserName(c)
	insertId := s.sysDictTypeService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 字典类型修改
//
// PUT /
func (s SysDictTypeController) Edit(c *gin.Context) {
	var body model.SysDictType
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.DictId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: dictId is empty"))
		return
	}

	// 检查数据是否存在
	dictInfo := s.sysDictTypeService.FindById(body.DictId)
	if dictInfo.DictId != body.DictId {
		c.JSON(200, response.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典名称唯一
	uniqueDictName := s.sysDictTypeService.CheckUniqueByName(body.DictName, body.DictId)
	if !uniqueDictName {
		msg := fmt.Sprintf("字典修改【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueDictType := s.sysDictTypeService.CheckUniqueByType(body.DictType, body.DictId)
	if !uniqueDictType {
		msg := fmt.Sprintf("字典修改【%s】失败，字典类型已存在", body.DictType)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	dictInfo.DictName = body.DictName
	dictInfo.DictType = body.DictType
	dictInfo.StatusFlag = body.StatusFlag
	dictInfo.Remark = body.Remark
	dictInfo.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysDictTypeService.Update(dictInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 字典类型删除
//
// DELETE /:dictId
func (s SysDictTypeController) Remove(c *gin.Context) {
	dictIdStr := c.Param("dictId")
	dictIds := parse.RemoveDuplicatesToArray(dictIdStr, ",")
	if dictIdStr == "" || len(dictIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: dictId not is empty"))
		return
	}

	rows, err := s.sysDictTypeService.DeleteByIds(dictIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}

// Refresh 字典类型刷新缓存
//
// PUT /refresh
func (s SysDictTypeController) Refresh(c *gin.Context) {
	s.sysDictTypeService.CacheClean("*")
	s.sysDictTypeService.CacheLoad("*")
	c.JSON(200, response.Ok(nil))
}

// Options 字典类型选择框列表
//
// GET /options
func (s SysDictTypeController) Options(c *gin.Context) {
	data := s.sysDictTypeService.Find(model.SysDictType{})

	type labelValue struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	// 数据组
	arr := make([]labelValue, 0)
	for _, v := range data {
		arr = append(arr, labelValue{
			Label: v.DictName,
			Value: v.DictType,
		})
	}
	c.JSON(200, response.OkData(arr))
}

// Export 字典类型列表导出
//
// GET /export
func (s SysDictTypeController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := context.QueryMap(c)
	rows, total := s.sysDictTypeService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("dict_type_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "字典编号",
		"B1": "字典名称",
		"C1": "字典类型",
		"D1": "字典状态",
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
			"A" + idx: row.DictId,
			"B" + idx: row.DictName,
			"C" + idx: row.DictType,
			"D" + idx: statusValue,
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
