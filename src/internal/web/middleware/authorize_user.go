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

// AuthOption 描述路由所需的角色/权限要求，零值表示仅校验登录态。
//
//	HasRoles / HasPerms   持有列表中【任意一项】即满足该类别
//	MatchRoles / MatchPerms 持有列表中【全部项】才满足该类别
type AuthOption struct {
	HasRoles   []string // 满足其一即可
	HasPerms   []string // 满足其一即可
	MatchRoles []string // 需全部满足
	MatchPerms []string // 需全部满足
}

// IsZero 报告是否未指定任何角色/权限要求。
func (o AuthOption) IsZero() bool {
	return len(o.HasRoles) == 0 && len(o.HasPerms) == 0 &&
		len(o.MatchRoles) == 0 && len(o.MatchPerms) == 0
}

// AuthorizeLogin 仅校验登录态：解析令牌、加载登录用户，不做角色/权限校验。
// 等价于 AuthorizeUser(AuthOption{})。
func AuthorizeLogin() gin.HandlerFunc {
	return AuthorizeUser(AuthOption{})
}

// AuthorizeUser 校验请求携带的访问令牌，解析登录用户并按 option 做角色/权限校验。
// 传零值 AuthOption 时仅校验登录态；请在需要登录用户信息的其它中间件之前注册。
func AuthorizeUser(option AuthOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := reqctx.Authorization(c)
		if tokenStr == "" {
			abortAuth(c, resp.CODE_AUTH_NOTOKEN, "authorization token is empty")
			return
		}

		claims, err := token.UserTokenVerify(tokenStr, "access")
		if err != nil {
			abortAuth(c, resp.CODE_AUTH, err.Error())
			return
		}

		info := token.UserInfoGet(claims)
		if info.UserId <= 0 {
			abortAuth(c, resp.CODE_AUTH_INVALID, "invalid login user information")
			return
		}
		c.Set(constants.CTX_LOGIN_USER, info)

		// 角色权限校验
		if !option.IsZero() && !checkAuth(info.RoleKeys, info.Permissions, option) {
			msg := fmt.Sprintf("unauthorized access %s %s", c.Request.Method, c.Request.RequestURI)
			c.JSON(403, resp.CodeMsg(resp.CODE_PERMISSION, msg))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		c.Next()
	}
}

// abortAuth 统一写出 401 鉴权失败响应并终止后续处理。
func abortAuth(c *gin.Context, code int, msg string) {
	c.JSON(401, resp.CodeMsg(code, msg))
	c.Abort()
}

// checkAuth 按 AuthOption 语义校验角色/权限。
// 系统管理员角色或全权权限直接放行；否则各指定类别需全部满足。
func checkAuth(roles, perms []string, opt AuthOption) bool {
	if slices.Contains(roles, constants.SYS_ROLE_SYSTEM_KEY) ||
		slices.Contains(perms, constants.SYS_PERMISSION_SYSTEM) {
		return true
	}
	if len(opt.HasRoles) > 0 && !containsAny(roles, opt.HasRoles) {
		return false
	}
	if len(opt.HasPerms) > 0 && !containsAny(perms, opt.HasPerms) {
		return false
	}
	if len(opt.MatchRoles) > 0 && !containsAll(roles, opt.MatchRoles) {
		return false
	}
	if len(opt.MatchPerms) > 0 && !containsAll(perms, opt.MatchPerms) {
		return false
	}
	return true
}

// containsAny 报告 origin 是否含有 required 中的任意一项。
func containsAny(origin, required []string) bool {
	for _, r := range required {
		if slices.Contains(origin, r) {
			return true
		}
	}
	return false
}

// containsAll 报告 origin 是否含有 required 的全部项。
func containsAll(origin, required []string) bool {
	for _, r := range required {
		if !slices.Contains(origin, r) {
			return false
		}
	}
	return true
}
