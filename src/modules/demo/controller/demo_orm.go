package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/demo/model"
	"mask_api_gin/src/modules/demo/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewDemoORM 实例化控制层 DemoORMController
var NewDemoORM = &DemoORMController{
	demoORMService: service.NewDemoORMService,
}

// DemoORMController 测试ORM
//
// PATH /demo
type DemoORMController struct {
	// 测试ORM信息服务
	demoORMService service.DemoORMService
}

// List 列表分页
//
// GET /list
func (s *DemoORMController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	data, err := s.demoORMService.SelectPage(query)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(data))
}

// All 列表无分页
//
// GET /all
func (s *DemoORMController) All(c *gin.Context) {
	demoORM := model.DemoORM{}

	title := c.Query("title")
	if title != "" {
		demoORM.Title = title
	}

	data, err := s.demoORMService.SelectList(demoORM)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.OkData(data))
}

// Info 信息
//
// GET /:id
func (s *DemoORMController) Info(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	data, err := s.demoORMService.SelectById(id)
	if err == nil {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Add 新增
//
// POST /
func (s *DemoORMController) Add(c *gin.Context) {
	var body model.DemoORM
	err := c.ShouldBindJSON(&body)
	if err != nil || body.ID != 0 {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	demoORM, err := s.demoORMService.Insert(body)
	if err == nil {
		c.JSON(200, result.OkData(demoORM.ID))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Edit 更新
//
// PUT /
func (s *DemoORMController) Edit(c *gin.Context) {
	var body model.DemoORM
	err := c.ShouldBindJSON(&body)
	if err != nil || body.ID == 0 {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	demoORM, err := s.demoORMService.Update(body)
	if err == nil {
		c.JSON(200, result.OkData(demoORM.ID))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 删除
//
// DELETE /:ids
func (s *DemoORMController) Remove(c *gin.Context) {
	ids := c.Param("ids")
	if ids == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组后去重
	idArr := strings.Split(ids, ",")
	uniqueIDs := parse.RemoveDuplicates(idArr)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows := s.demoORMService.DeleteByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Clean 清空
//
// DELETE /clean
func (s *DemoORMController) Clean(c *gin.Context) {
	_, err := s.demoORMService.Clean()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}
