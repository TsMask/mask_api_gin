package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/server/reqctx"
	"mask_api_gin/src/internal/server/resp"
	"mask_api_gin/src/internal/core/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
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
	query := reqctx.QueryMap(c)
	rows, total := s.sysDictTypeService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 字典类型信息
//
// GET /:id
func (s SysDictTypeController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.sysDictTypeService.FindById(id)
	if data.ID == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 字典类型新增
//
// POST /
func (s SysDictTypeController) Add(c *gin.Context) {
	var body model.SysDictType
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id not is empty"))
		return
	}

	// 检查字典名称唯一
	uniqueName := s.sysDictTypeService.CheckUniqueByName(body.DictName, 0)
	if !uniqueName {
		msg := fmt.Sprintf("字典新增【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueType := s.sysDictTypeService.CheckUniqueByType(body.DictType, 0)
	if !uniqueType {
		msg := fmt.Sprintf("字典新增【%s】失败，字典类型已存在", body.DictName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysDictTypeService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 字典类型修改
//
// PUT /
func (s SysDictTypeController) Edit(c *gin.Context) {
	var body model.SysDictType
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 检查数据是否存在
	dictType := s.sysDictTypeService.FindById(body.ID)
	if dictType.ID != body.ID {
		c.JSON(200, resp.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典名称唯一
	uniqueName := s.sysDictTypeService.CheckUniqueByName(body.DictName, body.ID)
	if !uniqueName {
		msg := fmt.Sprintf("字典修改【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueType := s.sysDictTypeService.CheckUniqueByType(body.DictType, body.ID)
	if !uniqueType {
		msg := fmt.Sprintf("字典修改【%s】失败，字典类型已存在", body.DictName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	dictType.DictName = body.DictName
	dictType.DictType = body.DictType
	dictType.DictStatus = body.DictStatus
	dictType.Remark = body.Remark
	dictType.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysDictTypeService.Update(dictType)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 字典类型删除
//
// DELETE /:id
func (s SysDictTypeController) Remove(c *gin.Context) {
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

	rows, err := s.sysDictTypeService.DeleteByIds(ids)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// Refresh 字典类型刷新缓存
//
// PUT /refresh
func (s SysDictTypeController) Refresh(c *gin.Context) {
	s.sysDictTypeService.CacheClean("*")
	s.sysDictTypeService.CacheLoad("*")
	c.JSON(200, resp.Ok(nil))
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
	c.JSON(200, resp.OkData(arr))
}

// Export 字典类型列表导出
//
// GET /export
func (s SysDictTypeController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	rows, total := s.sysDictTypeService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("dict_type_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysDictTypeService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
