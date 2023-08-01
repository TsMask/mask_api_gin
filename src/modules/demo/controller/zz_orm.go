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

// 演示-GORM基本使用
var ZzOrm = &zzorm{
	zzOrmService: *service.NewZzOrmService,
}

type zzorm struct {
	// 测试ORM信息服务
	zzOrmService service.ZzOrmService
}

// 列表分页
//
// GET /list
func (s *zzorm) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data, err := s.zzOrmService.SelectPage(querys)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(data))
}

// 列表无分页
//
// GET /all
func (s *zzorm) All(c *gin.Context) {
	zzOrm := model.ZzOrm{}

	title := c.Query("title")
	if title != "" {
		zzOrm.Title = title
	}

	data, err := s.zzOrmService.SelectList(zzOrm)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.OkData(data))
}

// 信息
//
// GET /:id
func (s *zzorm) Info(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	data, err := s.zzOrmService.SelectById(id)
	if err == nil {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 新增
//
// POST /
func (s *zzorm) Add(c *gin.Context) {
	var body model.ZzOrm
	err := c.ShouldBindJSON(&body)
	if err != nil || body.ID != 0 {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	zzOrm, err := s.zzOrmService.Insert(body)
	if err == nil {
		c.JSON(200, result.OkData(zzOrm.ID))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 更新
//
// PUT /
func (s *zzorm) Edit(c *gin.Context) {
	var body model.ZzOrm
	err := c.ShouldBindJSON(&body)
	if err != nil || body.ID == 0 {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	zzOrm, err := s.zzOrmService.Update(body)
	if err == nil {
		c.JSON(200, result.OkData(zzOrm.ID))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 删除
//
// DELETE /:ids
func (s *zzorm) Remove(c *gin.Context) {
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
	rows := s.zzOrmService.DeleteByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 清空
//
// DELETE /clean
func (s *zzorm) Clean(c *gin.Context) {
	_, err := s.zzOrmService.Clean()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}
