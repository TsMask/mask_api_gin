package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/demo/model"
	"mask_api_gin/src/modules/demo/service"
)

// NewDemoORM 实例化控制层
var NewDemoORM = &DemoORMController{
	demoORMService: service.NewDemoORMService, // 测试ORM信息服务
}

// DemoORMController 测试ORM 控制层处理
//
// PATH /demo/orm
type DemoORMController struct {
	demoORMService *service.DemoORMService // 测试ORM信息服务
}

// List 列表分页
//
// GET /list
func (s DemoORMController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	rows, total := s.demoORMService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// All 列表无分页
//
// GET /all
func (s DemoORMController) All(c *gin.Context) {
	var query struct {
		Title      string `form:"title" binding:"omitempty"`
		StatusFlag string `form:"statusFlag" binding:"omitempty,oneof=0 1"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	demoORM := model.DemoORM{}
	if query.Title != "" {
		demoORM.Title = query.Title
	}
	if query.StatusFlag != "" {
		demoORM.StatusFlag = query.StatusFlag
	}

	data := s.demoORMService.Find(demoORM)
	c.JSON(200, resp.OkData(data))
}

// Info 信息
//
// GET /:id
func (s DemoORMController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.demoORMService.FindById(id)
	if data.Id == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 新增
//
// POST /
func (s DemoORMController) Add(c *gin.Context) {
	var body model.DemoORM
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.Id > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id not is empty"))
		return
	}

	insertId := s.demoORMService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 更新
//
// PUT /
func (s DemoORMController) Edit(c *gin.Context) {
	var body model.DemoORM
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.Id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	rowsAffected := s.demoORMService.Update(body)
	if rowsAffected > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 删除
//
// DELETE /:id
func (s DemoORMController) Remove(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(id, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	rowsAffected := s.demoORMService.DeleteByIds(ids)
	if rowsAffected > 0 {
		msg := fmt.Sprintf("删除成功：%d", rowsAffected)
		c.JSON(200, resp.OkMsg(msg))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Clean 清空
//
// DELETE /clean
func (s DemoORMController) Clean(c *gin.Context) {
	rows := s.demoORMService.Clean()
	c.JSON(200, resp.OkData(rows))
}
