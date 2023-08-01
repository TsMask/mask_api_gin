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

// 实例化控制层 SysPostController 结构体
var NewSysPost = &SysPostController{
	sysPostService: service.NewSysPostImpl,
}

// 岗位信息
//
// PATH /system/post
type SysPostController struct {
	// 岗位服务
	sysPostService service.ISysPost
}

// 岗位列表
//
// GET /list
func (s *SysPostController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data := s.sysPostService.SelectPostPage(querys)
	c.JSON(200, result.Ok(data))
}

// 岗位信息
//
// GET /:postId
func (s *SysPostController) Info(c *gin.Context) {
	postId := c.Param("postId")
	if postId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysPostService.SelectPostById(postId)
	if data.PostID == postId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 岗位新增
//
// POST /
func (s *SysPostController) Add(c *gin.Context) {
	var body model.SysPost
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.PostID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查名称唯一
	uniqueuPostName := s.sysPostService.CheckUniquePostName(body.PostName, "")
	if !uniqueuPostName {
		msg := fmt.Sprintf("岗位新增【%s】失败，岗位名称已存在", body.PostName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查编码属性值唯一
	uniquePostCode := s.sysPostService.CheckUniquePostCode(body.PostCode, "")
	if !uniquePostCode {
		msg := fmt.Sprintf("岗位新增【%s】失败，岗位编码已存在", body.PostCode)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysPostService.InsertPost(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 岗位修改
//
// PUT /
func (s *SysPostController) Edit(c *gin.Context) {
	var body model.SysPost
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.PostID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	post := s.sysPostService.SelectPostById(body.PostID)
	if post.PostID != body.PostID {
		c.JSON(200, result.ErrMsg("没有权限访问岗位数据！"))
		return
	}

	// 检查名称唯一
	uniqueuPostName := s.sysPostService.CheckUniquePostName(body.PostName, body.PostID)
	if !uniqueuPostName {
		msg := fmt.Sprintf("岗位修改【%s】失败，岗位名称已存在", body.PostName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查编码属性值唯一
	uniquePostCode := s.sysPostService.CheckUniquePostCode(body.PostCode, body.PostID)
	if !uniquePostCode {
		msg := fmt.Sprintf("岗位修改【%s】失败，岗位编码已存在", body.PostCode)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysPostService.UpdatePost(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 岗位删除
//
// DELETE /:postIds
func (s *SysPostController) Remove(c *gin.Context) {
	postIds := c.Param("postIds")
	if postIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(postIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows, err := s.sysPostService.DeletePostByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 导出岗位信息
//
// POST /export
func (s *SysPostController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMap(c)
	data := s.sysPostService.SelectPostPage(querys)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysPost)

	// 导出文件名称
	fileName := fmt.Sprintf("post_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "岗位编号",
		"B1": "岗位编码",
		"C1": "岗位名称",
		"D1": "岗位排序",
		"E1": "状态",
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
			"A" + idx: row.PostID,
			"B" + idx: row.PostCode,
			"C" + idx: row.PostName,
			"D" + idx: row.PostSort,
			"E" + idx: statusValue,
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
