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

// NewSysPost 实例化控制层
var NewSysPost = &SysPostController{
	sysPostService: service.NewSysPost,
}

// SysPostController 岗位信息
//
// PATH /system/post
type SysPostController struct {
	sysPostService *service.SysPost // 岗位服务
}

// List 岗位列表
//
// GET /list
func (s SysPostController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	rows, total := s.sysPostService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 岗位信息
//
// GET /:postId
func (s SysPostController) Info(c *gin.Context) {
	postId := c.Param("postId")
	if postId == "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	data := s.sysPostService.FindById(postId)
	if data.PostId == postId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 岗位新增
//
// POST /
func (s SysPostController) Add(c *gin.Context) {
	var body model.SysPost
	if err := c.ShouldBindBodyWithJSON(&body); err != nil || body.PostId != "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 检查名称唯一
	uniquePostName := s.sysPostService.CheckUniqueByName(body.PostName, "")
	if !uniquePostName {
		msg := fmt.Sprintf("岗位新增【%s】失败，岗位名称已存在", body.PostName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查编码属性值唯一
	uniquePostCode := s.sysPostService.CheckUniqueByCode(body.PostCode, "")
	if !uniquePostCode {
		msg := fmt.Sprintf("岗位新增【%s】失败，岗位编码已存在", body.PostCode)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysPostService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 岗位修改
//
// PUT /
func (s SysPostController) Edit(c *gin.Context) {
	var body model.SysPost
	if err := c.ShouldBindBodyWithJSON(&body); err != nil || body.PostId <= "" {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 检查是否存在
	postInfo := s.sysPostService.FindById(body.PostId)
	if postInfo.PostId != body.PostId {
		c.JSON(200, response.ErrMsg("没有权限访问岗位数据！"))
		return
	}

	// 检查名称唯一
	uniquePostName := s.sysPostService.CheckUniqueByName(body.PostName, body.PostId)
	if !uniquePostName {
		msg := fmt.Sprintf("岗位修改【%s】失败，岗位名称已存在", body.PostName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查编码属性值唯一
	uniquePostCode := s.sysPostService.CheckUniqueByCode(body.PostCode, body.PostId)
	if !uniquePostCode {
		msg := fmt.Sprintf("岗位修改【%s】失败，岗位编码已存在", body.PostCode)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	postInfo.PostCode = body.PostCode
	postInfo.PostName = body.PostName
	postInfo.PostSort = body.PostSort
	postInfo.StatusFlag = body.StatusFlag
	postInfo.Remark = body.Remark
	postInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysPostService.Update(postInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 岗位删除
//
// DELETE /:postId
func (s SysPostController) Remove(c *gin.Context) {
	postIdsStr := c.Param("postId")
	postIds := parse.RemoveDuplicatesToArray(postIdsStr, ",")
	if postIdsStr == "" || len(postIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	rows, err := s.sysPostService.DeleteByIds(postIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}

// Export 导出岗位信息
//
// GET /export
func (s SysPostController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := ctx.QueryMap(c)
	rows, total := s.sysPostService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("post_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "岗位编号",
		"B1": "岗位编码",
		"C1": "岗位名称",
		"D1": "岗位排序",
		"E1": "岗位状态",
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
			"A" + idx: row.PostId,
			"B" + idx: row.PostCode,
			"C" + idx: row.PostName,
			"D" + idx: row.PostSort,
			"E" + idx: statusValue,
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
