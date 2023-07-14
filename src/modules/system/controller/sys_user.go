package controller

import (
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/framework/vo/result"
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
