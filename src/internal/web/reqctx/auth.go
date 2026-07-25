package reqctx

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/web/token"
)

// LoginUser 登录用户信息
func LoginUser(c *gin.Context) (token.UserInfo, error) {
	value, exists := c.Get(constants.CTX_LOGIN_USER)
	if exists && value != nil {
		return value.(token.UserInfo), nil
	}
	return token.UserInfo{}, fmt.Errorf("invalid login user information")
}

// LoginUserToUserID 登录用户信息-用户ID
func LoginUserToUserID(c *gin.Context) int64 {
	info, err := LoginUser(c)
	if err != nil {
		return 0
	}
	return info.UserId
}

// LoginUserToUserName 登录用户信息-用户名称
func LoginUserToUserName(c *gin.Context) string {
	info, err := LoginUser(c)
	if err != nil {
		return ""
	}
	return info.UserName
}

// LoginUserByContainRoles 登录用户信息-包含角色KEY
func LoginUserByContainRoles(c *gin.Context, target string) bool {
	info, err := LoginUser(c)
	if err != nil {
		return false
	}
	if config.IsSystemUser(info.UserId) {
		return true
	}
	return slices.Contains(info.RoleKeys, target)
}

// LoginUserByContainPerms 登录用户信息-包含权限标识
func LoginUserByContainPerms(c *gin.Context, target string) bool {
	info, err := LoginUser(c)
	if err != nil {
		return false
	}
	if config.IsSystemUser(info.UserId) {
		return true
	}
	perms := info.Permissions
	return slices.Contains(perms, target)
}
