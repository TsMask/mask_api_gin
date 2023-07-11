package controller

import (
	"fmt"
	"mask_api_gin/src/framework/model/result"
	"mask_api_gin/src/modules/common/model"
	"mask_api_gin/src/modules/common/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 账号身份操作处理
var Account = &accountController{
	accountService: service.AccountImpl,
}

type accountController struct {
	// 账号身份操作服务
	accountService service.IAccount
}

// 系统登录
//
// POST /login
func (s *accountController) Login(c *gin.Context) {
	var loginBody model.LoginBody
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.JSON(200, result.Ok(map[string]interface{}{
		"data": loginBody,
	}))
}

// 登录用户信息
//
// GET /getInfo
func (s *accountController) Info(c *gin.Context) {
	name := viper.GetString("pkg.name")
	version := viper.GetString("pkg.version")
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}

// 登录用户路由信息
//
// GET /getRouters
func (s *accountController) Router(c *gin.Context) {
	name := viper.GetString("pkg.name")
	version := viper.GetString("pkg.version")
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}

// 系统登出
//
// POST /logout
func (s *accountController) Logout(c *gin.Context) {
	c.JSON(200, result.OkMsg("退出成功"))
}
