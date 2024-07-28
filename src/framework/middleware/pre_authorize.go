package middleware

import (
	"fmt"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/token"
	"mask_api_gin/src/framework/vo/result"

	"github.com/gin-gonic/gin"
)

// PreAuthorize 用户身份授权认证校验
//
// 只需含有其中角色 "hasRoles": {"xxx"},
//
// 只需含有其中权限 "hasPerms": {"xxx"},
//
// 同时匹配其中角色 "matchRoles": {"xxx"},
//
// 同时匹配其中权限 "matchPerms": {"xxx"},
func PreAuthorize(options map[string][]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头标识信息
		tokenStr := ctx.Authorization(c)
		if tokenStr == "" {
			c.JSON(401, result.CodeMsg(401, "无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 验证令牌
		claims, err := token.Verify(tokenStr)
		if err != nil {
			c.JSON(401, result.CodeMsg(401, err.Error()))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 获取缓存的用户信息
		loginUser := token.LoginUser(claims)
		if loginUser.UserID == "" {
			c.JSON(401, result.CodeMsg(401, "无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 检查刷新有效期后存入上下文
		token.RefreshIn(&loginUser)
		c.Set(constSystem.CtxLoginUser, loginUser)

		// 登录用户角色权限校验
		if options != nil {
			var roles []string
			for _, item := range loginUser.User.Roles {
				roles = append(roles, item.RoleKey)
			}
			perms := loginUser.Permissions
			verifyOk := verifyRolePermission(roles, perms, options)
			if !verifyOk {
				msg := fmt.Sprintf("无权访问 %s %s", c.Request.Method, c.Request.RequestURI)
				c.JSON(403, result.CodeMsg(403, msg))
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
	if contains(roles, constSystem.RoleKey) || contains(perms, constSystem.Permission) {
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

// contains 检查字符串数组中是否包含指定的字符串
func contains(arr []string, target string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}

// some 检查字符串数组中含有其中一项
func some(origin []string, target []string) bool {
	has := false
	for _, t := range target {
		for _, o := range origin {
			if t == o {
				has = true
				break
			}
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
		found := false
		for _, o := range origin {
			if t == o {
				found = true
				break
			}
		}
		if !found {
			match = false
			break
		}
	}
	return match
}
