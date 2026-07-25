package middleware

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/internal/web/token"
)

// AuthorizeUser 用户身份授权认证校验
//
// 只需含有其中角色 "hasRoles": {"xxx"},
//
// 只需含有其中权限 "hasPerms": {"xxx"},
//
// 同时匹配其中角色 "matchRoles": {"xxx"},
//
// 同时匹配其中权限 "matchPerms": {"xxx"},
func AuthorizeUser(options map[string][]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头标识信息
		tokenStr := reqctx.Authorization(c)
		if tokenStr == "" {
			c.JSON(401, resp.CodeMsg(resp.CODE_AUTH_NOTOKEN, "authorization token is empty"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 验证令牌
		claims, err := token.UserTokenVerify(tokenStr, "access")
		if err != nil {
			c.JSON(401, resp.CodeMsg(resp.CODE_AUTH, err.Error()))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 获取缓存的用户信息
		info := token.UserInfoGet(claims)
		if info.UserId <= 0 {
			c.JSON(401, resp.CodeMsg(resp.CODE_AUTH_INVALID, "invalid login user information"))
			c.Abort() // 停止执行后续的处理函数
			return
		}
		c.Set(constants.CTX_LOGIN_USER, info)

		// 登录用户角色权限校验
		if options != nil {
			roles := info.RoleKeys
			perms := info.Permissions
			verifyOk := verifyRolePermission(roles, perms, options)
			if !verifyOk {
				msg := fmt.Sprintf("unauthorized access %s %s", c.Request.Method, c.Request.RequestURI)
				c.JSON(403, resp.CodeMsg(resp.CODE_PERMISSION, msg))
				c.Abort() // 停止执行后续的处理函数
				return
			}
		}

		// 调用下一个处理程序
		c.Next()
	}
}

// verifyRolePermission 校验角色权限是否满足
//
// roles 角色字符数组
//
// perms 权限字符数组
//
// options 参数
func verifyRolePermission(roles, perms []string, options map[string][]string) bool {
	// 直接放行 系统管理员角色或任意权限
	if slices.Contains(roles, constants.SYS_ROLE_SYSTEM_KEY) || slices.Contains(perms, constants.SYS_PERMISSION_SYSTEM) {
		return true
	}
	opts := make([]bool, 4)

	// 只需含有其中角色
	hasRole := false
	if arr, ok := options["hasRoles"]; ok && len(arr) > 0 {
		hasRole = some(roles, arr)
		opts[0] = true
	}

	// 只需含有其中权限
	hasPerms := false
	if arr, ok := options["hasPerms"]; ok && len(arr) > 0 {
		hasPerms = some(perms, arr)
		opts[1] = true
	}

	// 同时匹配其中角色
	matchRoles := false
	if arr, ok := options["matchRoles"]; ok && len(arr) > 0 {
		matchRoles = every(roles, arr)
		opts[2] = true
	}

	// 同时匹配其中权限
	matchPerms := false
	if arr, ok := options["matchPerms"]; ok && len(arr) > 0 {
		matchPerms = every(perms, arr)
		opts[3] = true
	}

	// 同时判断 含有其中
	if opts[0] && opts[1] {
		return hasRole || hasPerms
	}
	// 同时判断 匹配其中
	if opts[2] && opts[3] {
		return matchRoles && matchPerms
	}
	// 同时判断 含有其中且匹配其中
	if opts[0] && opts[3] {
		return hasRole && matchPerms
	}
	if opts[1] && opts[2] {
		return hasPerms && matchRoles
	}

	return hasRole || hasPerms || matchRoles || matchPerms
}

// some 检查字符串数组中含有其中一项
func some(origin []string, target []string) bool {
	has := false
	for _, t := range target {
		if slices.Contains(origin, t) {
			has = true
		}
		if has {
			break
		}
	}
	return has
}

// every 检查字符串数组中同时包含所有项
func every(origin []string, target []string) bool {
	match := true
	for _, t := range target {
		if !slices.Contains(origin, t) {
			match = false
			break
		}
	}
	return match
}
