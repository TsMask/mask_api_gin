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

// NewSysPost 实例化控制层
var NewSysPost = &SysPostController{
	sysPostService: service.NewSysPost,
}

// SysPostController 岗位信息 控制层处理
//
// PATH /system/post
type SysPostController struct {
	sysPostService *service.SysPost // 岗位服务
}

// List 岗位列表
//
// GET /list
func (s SysPostController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	rows, total := s.sysPostService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 岗位信息
//
// GET /:id
func (s SysPostController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.sysPostService.FindById(id)
	if data.ID == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 岗位新增
//
// POST /
func (s SysPostController) Add(c *gin.Context) {
	var body model.SysPost
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: postId not is empty"))
		return
	}

	// 检查名称唯一
	uniqueuName := s.sysPostService.CheckUniqueByName(body.PostName, 0)
	if !uniqueuName {
		msg := fmt.Sprintf("岗位新增【%s】失败，岗位名称已存在", body.PostName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查编码属性值唯一
	uniqueCode := s.sysPostService.CheckUniqueByCode(body.PostCode, 0)
	if !uniqueCode {
		msg := fmt.Sprintf("岗位新增【%s】失败，岗位编码已存在", body.PostCode)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysPostService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 岗位修改
//
// PUT /
func (s SysPostController) Edit(c *gin.Context) {
	var body model.SysPost
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: postId is empty"))
		return
	}

	// 检查是否存在
	postInfo := s.sysPostService.FindById(body.ID)
	if postInfo.ID != body.ID {
		c.JSON(200, resp.ErrMsg("没有权限访问岗位数据！"))
		return
	}

	// 检查名称唯一
	uniquePostName := s.sysPostService.CheckUniqueByName(body.PostName, body.ID)
	if !uniquePostName {
		msg := fmt.Sprintf("岗位修改【%s】失败，岗位名称已存在", body.PostName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查编码属性值唯一
	uniquePostCode := s.sysPostService.CheckUniqueByCode(body.PostCode, body.ID)
	if !uniquePostCode {
		msg := fmt.Sprintf("岗位修改【%s】失败，岗位编码已存在", body.PostCode)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	postInfo.PostCode = body.PostCode
	postInfo.PostName = body.PostName
	postInfo.PostSort = body.PostSort
	postInfo.PostStatus = body.PostStatus
	postInfo.Remark = body.Remark
	postInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysPostService.Update(postInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 岗位删除
//
// DELETE /:id
func (s SysPostController) Remove(c *gin.Context) {
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

	rows, err := s.sysPostService.DeleteByIds(ids)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// Export 导出岗位信息
//
// GET /export
func (s SysPostController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	rows, total := s.sysPostService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("post_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysPostService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
