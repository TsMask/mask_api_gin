package controller

import (
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// 用户信息
//
// PATH /system/user
var SysUser = &sysUserController{
	sysUserService: service.SysUserImpl,
}

type sysUserController struct {
	// 用户服务
	sysUserService service.ISysUser
}

// 用户信息列表
//
// GET /list
func (s *sysUserController) List(c *gin.Context) {
	// 查询参数转换map
	querys := ctxUtils.QueryMapString(c)
	dataScopeSQL := repoUtils.DataScopeSQL("d", "u")
	list := s.sysUserService.SelectUserPage(querys, dataScopeSQL)
	c.JSON(200, result.Ok(list))
}

// 用户状态修改
//
// PUT /changeStatus
func (s *sysUserController) Status(c *gin.Context) {
	loginUser, err := ctxUtils.LoginUser(c)
	if err != nil {
		c.JSON(401, result.Err(map[string]interface{}{
			"code": 401,
			"msg":  err.Error(),
		}))
		c.Abort() // 停止执行后续的处理函数
		return
	}

	var body struct {
		UserID string `json:"userId"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 必要字段
	if body.UserID == "" || body.Status == "" {
		c.JSON(200, result.Err(nil))
		return
	}

	user := s.sysUserService.SelectUserById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}
	// 与旧值相等不变更
	if user.Status == body.Status {
		c.JSON(200, result.ErrMsg("与旧值相等！"))
		return
	}
	sysUser := model.SysUser{
		UserID:   body.UserID,
		Status:   body.Status,
		UpdateBy: loginUser.User.UserName,
	}
	num := s.sysUserService.UpdateUser(sysUser)
	if num > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}
