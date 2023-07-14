package middleware

import (
	"fmt"
	AdminConstants "mask_api_gin/src/pkg/constants/admin"
	commonConstants "mask_api_gin/src/pkg/constants/common"
	"mask_api_gin/src/pkg/model/result"
	ctxUtils "mask_api_gin/src/pkg/utils/ctx"
	tokenUtils "mask_api_gin/src/pkg/utils/token"

	"github.com/gin-gonic/gin"
)

// PreAuthorize 用户身份授权认证校验
//
// 只需含有其中角色 hasRoles: []string{}
//
// 只需含有其中权限 hasPermissions: []string{}
//
// 同时匹配其中角色 matchRoles: []string{}
//
// 同时匹配其中权限 matchPermissions: []string{}
func PreAuthorize(options map[string][]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头标识信息
		tokenStr := ctxUtils.Authorization(c)
		if tokenStr == "" {
			c.JSON(401, result.ErrMsg("无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 验证令牌
		claims, err := tokenUtils.Verify(tokenStr)
		if err != nil {
			c.JSON(401, result.ErrMsg("无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 获取缓存的用户信息
		loginUser := tokenUtils.LoginUser(claims)
		if loginUser.UserID == "" {
			c.JSON(401, result.ErrMsg("无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 检查刷新有效期后存入上下文
		tokenUtils.Refresh(&loginUser)
		c.Set(commonConstants.CTX_LOGIN_USER, loginUser)

		// 登录用户角色权限校验
		if options != nil {
			var roles []string
			for _, item := range loginUser.User.Roles {
				roles = append(roles, item.RoleKey)
			}
			permissions := loginUser.Permissions
			verifyOk := verifyRolePermission(roles, permissions, options)
			if !verifyOk {
				msg := fmt.Sprintf("无权访问 %s %s", c.Request.Method, c.Request.RequestURI)
				c.JSON(403, result.ErrMsg(msg))
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
// permissions 权限字符数组
//
// metadata 装饰器参数身份
func verifyRolePermission(roles, permissions []string, options map[string][]string) bool {
	// 直接放行 管理员角色或任意权限
	if contains(roles, AdminConstants.ROLE_KEY) || contains(permissions, AdminConstants.PERMISSION) {
		return true
	}
	opts := make([]bool, 0, 4)

	// 只需含有其中角色
	hasRole := false
	if arr, ok := options["hasRoles"]; ok {
		hasRole = some(roles, arr)
		opts[0] = true
	}

	// 只需含有其中权限
	hasPermission := false
	if arr, ok := options["hasPermissions"]; ok {
		hasPermission = some(permissions, arr)
		opts[1] = true
	}

	// 同时匹配其中角色
	matchRoles := false
	if arr, ok := options["matchRoles"]; ok {
		matchRoles = every(roles, arr)
		opts[2] = true
	}

	// 同时匹配其中权限
	matchPermissions := false
	if arr, ok := options["matchPermissions"]; ok {
		matchPermissions = every(permissions, arr)
		opts[3] = true
	}

	// 同时判断 只需含有其中
	if opts[0] && opts[1] {
		return hasRole || hasPermission
	}
	// 同时判断 匹配其中
	if opts[2] && opts[3] {
		return matchRoles && matchPermissions
	}
	// 同时判断 含有其中且匹配其中
	if opts[0] && opts[3] {
		return hasRole && matchPermissions
	}
	if opts[1] && opts[2] {
		return hasPermission && matchRoles
	}

	return hasRole || hasPermission || matchRoles || matchPermissions
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