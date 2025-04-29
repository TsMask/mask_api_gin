package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/oauth2/model"
	"mask_api_gin/src/modules/oauth2/service"
)

// NewOauth2Client 实例化控制层
var NewOauth2Client = &Oauth2ClientController{
	oauth2ClientService: service.NewOauth2ClientService,
}

// Oauth2ClientController 客户端授权管理 控制层处理
//
// PATH /oauth2/client
type Oauth2ClientController struct {
	oauth2ClientService *service.Oauth2ClientService // 用户授权第三方应用信息服务
}

// List 列表
//
// GET /list
func (s Oauth2ClientController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	rows, total := s.oauth2ClientService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 信息
//
// GET /:clientId
func (s Oauth2ClientController) Info(c *gin.Context) {
	clientId := c.Param("clientId")
	if clientId == "" {
		c.JSON(422, resp.CodeMsg(422002, "bind err: clientId is empty"))
		return
	}

	info := s.oauth2ClientService.FindByClientId(clientId)
	if info.ClientId == clientId {
		c.JSON(200, resp.OkData(info))
		return
	}
	c.JSON(200, resp.ErrMsg("clientId does not exist"))
}

// Add 新增
//
// POST /
func (s Oauth2ClientController) Add(c *gin.Context) {
	var body model.Oauth2Client
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(422001, errMsgs))
		return
	}
	if body.Id > 0 {
		c.JSON(422, resp.CodeMsg(422002, "bind err: id not is empty"))
		return
	}

	// 本地IP地址不支持
	localHosts := []string{"127.0.0.1", "localhost", "::ffff:", "::1"}
	localHost := false
	for _, host := range localHosts {
		if strings.Contains(body.IPWhite, host) {
			localHost = true
			break
		}
	}
	if localHost {
		c.JSON(200, resp.ErrMsg("no support local host"))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.oauth2ClientService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 更新
//
// PUT /
func (s Oauth2ClientController) Edit(c *gin.Context) {
	var body model.Oauth2Client
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(422001, errMsgs))
		return
	}
	if body.Id <= 0 {
		c.JSON(422, resp.CodeMsg(422002, "bind err: id is empty"))
		return
	}

	// 本地IP地址不支持
	localHosts := []string{"127.0.0.1", "localhost", "::ffff:", "::1"}
	localHost := false
	for _, host := range localHosts {
		if strings.Contains(body.IPWhite, host) {
			localHost = true
			break
		}
	}
	if localHost {
		c.JSON(200, resp.ErrMsg("no support local host"))
		return
	}

	// 查询信息
	info := s.oauth2ClientService.FindById(body.Id)
	if info.ClientId == "" || info.Id != body.Id {
		c.JSON(200, resp.ErrMsg("modification failed, data not exist"))
		return
	}

	info.Title = body.Title
	info.IPWhite = body.IPWhite
	info.Remark = body.Remark
	info.UpdateBy = reqctx.LoginUserToUserName(c)
	rowsAffected := s.oauth2ClientService.Update(info)
	if rowsAffected > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 删除
//
// DELETE /:id
func (s Oauth2ClientController) Remove(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(422, resp.CodeMsg(422002, "bind err: id is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(id, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	rows, err := s.oauth2ClientService.DeleteByIds(ids)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}
