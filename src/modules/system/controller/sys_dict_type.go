package controller

import (
	"fmt"
	"mask_api_gin/src/framework/constants/common"
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

// 实例化控制层 SysDictTypeController 结构体
var NewSysDictType = &SysDictTypeController{
	sysDictTypeService: service.NewSysDictTypeImpl,
}

// 字典类型信息
//
// PATH /system/dict/type
type SysDictTypeController struct {
	// 字典类型服务
	sysDictTypeService service.ISysDictType
}

// 字典类型列表
//
// GET /list
func (s *SysDictTypeController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data := s.sysDictTypeService.SelectDictTypePage(querys)
	c.JSON(200, result.Ok(data))
}

// 字典类型信息
//
// GET /:dictId
func (s *SysDictTypeController) Info(c *gin.Context) {
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
func (s *SysDictTypeController) Add(c *gin.Context) {
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
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueDictType := s.sysDictTypeService.CheckUniqueDictType(body.DictType, "")
	if !uniqueDictType {
		msg := fmt.Sprintf("字典新增【%s】失败，字典类型已存在", body.DictType)
		c.JSON(200, result.ErrMsg(msg))
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
func (s *SysDictTypeController) Edit(c *gin.Context) {
	var body model.SysDictType
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DictID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查数据是否存在
	dictInfo := s.sysDictTypeService.SelectDictTypeByID(body.DictID)
	if dictInfo.DictID != body.DictID {
		c.JSON(200, result.ErrMsg("没有权限访问字典类型数据！"))
		return
	}

	// 检查字典名称唯一
	uniqueDictName := s.sysDictTypeService.CheckUniqueDictName(body.DictName, body.DictID)
	if !uniqueDictName {
		msg := fmt.Sprintf("字典修改【%s】失败，字典名称已存在", body.DictName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查字典类型唯一
	uniqueDictType := s.sysDictTypeService.CheckUniqueDictType(body.DictType, body.DictID)
	if !uniqueDictType {
		msg := fmt.Sprintf("字典修改【%s】失败，字典类型已存在", body.DictType)
		c.JSON(200, result.ErrMsg(msg))
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

// 字典类型删除
//
// DELETE /:dictIds
func (s *SysDictTypeController) Remove(c *gin.Context) {
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
func (s *SysDictTypeController) RefreshCache(c *gin.Context) {
	s.sysDictTypeService.ResetDictCache()
	c.JSON(200, result.Ok(nil))
}

// 字典类型选择框列表
//
// GET /getDictOptionselect
func (s *SysDictTypeController) DictOptionselect(c *gin.Context) {
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
func (s *SysDictTypeController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMap(c)
	data := s.sysDictTypeService.SelectDictTypePage(querys)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysDictType)

	// 导出文件名称
	fileName := fmt.Sprintf("dict_type_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "字典主键",
		"B1": "字典名称",
		"C1": "字典类型",
		"D1": "状态",
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
			"A" + idx: row.DictID,
			"B" + idx: row.DictName,
			"C" + idx: row.DictType,
			"D" + idx: statusValue,
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
