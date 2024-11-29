package controller

import (
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/demo/model"
	"mask_api_gin/src/modules/demo/service"

	"fmt"

	"github.com/gin-gonic/gin"
)

// NewDemoORM 实例化控制层
var NewDemoORM = &DemoORMController{
	demoORMService: service.NewDemoORMService, // 测试ORM信息服务
}

// DemoORMController 测试ORM 控制层处理
//
// PATH /demo
type DemoORMController struct {
	demoORMService *service.DemoORMService // 测试ORM信息服务
}

// List 列表分页
//
// GET /list
func (s DemoORMController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	rows, total := s.demoORMService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// All 列表无分页
//
// GET /all
func (s DemoORMController) All(c *gin.Context) {
	demoORM := model.DemoORM{}

	title := c.Query("title")
	if title != "" {
		demoORM.Title = title
	}
	statusFlag := c.Query("statusFlag")
	if statusFlag != "" {
		demoORM.StatusFlag = statusFlag
	}

	data := s.demoORMService.Find(demoORM)
	c.JSON(200, response.OkData(data))
}

// Info 信息
//
// GET /:id
func (s DemoORMController) Info(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: id is empty"))
		return
	}

	data := s.demoORMService.FindById(id)
	if data.Id == id {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 新增
//
// POST /
func (s DemoORMController) Add(c *gin.Context) {
	var body model.DemoORM
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.Id != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: id not is empty"))
		return
	}

	InsertId := s.demoORMService.Insert(body)
	if InsertId != "" {
		c.JSON(200, response.OkData(InsertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 更新
//
// PUT /
func (s DemoORMController) Edit(c *gin.Context) {
	var body model.DemoORM
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.Id == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: id is empty"))
		return
	}

	rowsAffected := s.demoORMService.Update(body)
	if rowsAffected > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 删除
//
// DELETE /:id
func (s DemoORMController) Remove(c *gin.Context) {
	idStr := c.Param("id")
	ids := parse.RemoveDuplicatesToArray(idStr, ",")
	if idStr == "" || len(ids) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: id is empty"))
		return
	}

	rows := s.demoORMService.DeleteByIds(ids)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, response.OkMsg(msg))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Clean 清空
//
// DELETE /clean
func (s DemoORMController) Clean(c *gin.Context) {
	_, err := s.demoORMService.Clean()
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, response.Ok(nil))
}
