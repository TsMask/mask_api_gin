package controller

import (
	"mask_api_gin/src/framework/config"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

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

// 用户信息删除
//
// DELETE /:userIds
func (s *sysUserController) Remove(c *gin.Context) {
	userIds := c.Param("userIds")
	if userIds == "" {
		c.JSON(400, result.ErrMsg("参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(userIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows := s.sysUserService.DeleteUserByIds(uniqueIDs)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 用户重置密码
//
// PUT /resetPwd
func (s *sysUserController) ResetPwd(c *gin.Context) {
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
		UserID   string `json:"userId"  binding:"required"`
		Password string `json:"password"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, result.ErrMsg("参数错误"))
		return
	}

	// 检查是否管理员用户
	if config.IsAdmin(body.UserID) {
		c.JSON(200, result.ErrMsg("不允许操作管理员用户"))
		return
	}

	user := s.sysUserService.SelectUserById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}
	if !regular.ValidPassword(body.Password) {
		c.JSON(200, result.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}

	sysUser := model.SysUser{
		UserID:   body.UserID,
		Password: body.Password,
		UpdateBy: loginUser.User.UserName,
	}
	rows := s.sysUserService.UpdateUser(sysUser)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
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
		UserID string `json:"userId"  binding:"required"`
		Status string `json:"status"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	user := s.sysUserService.SelectUserById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}
	// 与旧值相等不变更
	if user.Status == body.Status {
		c.JSON(200, result.ErrMsg("变更状态与旧值相等！"))
		return
	}
	sysUser := model.SysUser{
		UserID:   body.UserID,
		Status:   body.Status,
		UpdateBy: loginUser.User.UserName,
	}
	rows := s.sysUserService.UpdateUser(sysUser)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}
