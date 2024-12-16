package context

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/token"

	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoginUser 登录用户信息
func LoginUser(c *gin.Context) (token.LoginUser, error) {
	value, exists := c.Get(constants.CTX_LOGIN_USER)
	if exists && value != nil {
		return value.(token.LoginUser), nil
	}
	return token.LoginUser{}, fmt.Errorf("invalid login user information")
}

// LoginUserToUserID 登录用户信息-用户ID
func LoginUserToUserID(c *gin.Context) string {
	loginUser, err := LoginUser(c)
	if err != nil {
		return ""
	}
	return loginUser.UserId
}

// LoginUserToUserName 登录用户信息-用户名称
func LoginUserToUserName(c *gin.Context) string {
	loginUser, err := LoginUser(c)
	if err != nil {
		return ""
	}
	return loginUser.User.UserName
}

// LoginUserByContainRoles 登录用户信息-包含角色KEY
func LoginUserByContainRoles(c *gin.Context, target string) bool {
	loginUser, err := LoginUser(c)
	if err != nil {
		return false
	}
	if config.IsSystemUser(loginUser.UserId) {
		return true
	}
	roles := loginUser.User.Roles
	for _, item := range roles {
		if item.RoleKey == target {
			return true
		}
	}
	return false
}

// LoginUserByContainPerms 登录用户信息-包含权限标识
func LoginUserByContainPerms(c *gin.Context, target string) bool {
	loginUser, err := LoginUser(c)
	if err != nil {
		return false
	}
	if config.IsSystemUser(loginUser.UserId) {
		return true
	}
	perms := loginUser.Permissions
	for _, str := range perms {
		if str == target {
			return true
		}
	}
	return false
}

// LoginUserToDataScopeSQL 登录用户信息-角色数据范围过滤SQL字符串
func LoginUserToDataScopeSQL(c *gin.Context, deptAlias string, userAlias string) string {
	dataScopeSQL := ""
	// 登录用户信息
	loginUser, err := LoginUser(c)
	if err != nil {
		return dataScopeSQL
	}
	userInfo := loginUser.User

	// 如果是系统管理员，则不过滤数据
	if config.IsSystemUser(userInfo.UserId) {
		return dataScopeSQL
	}
	// 无用户角色
	if len(userInfo.Roles) <= 0 {
		return dataScopeSQL
	}

	// 记录角色权限范围定义添加过, 非自定数据权限不需要重复拼接SQL
	var scopeKeys []string
	var conditions []string
	for _, role := range userInfo.Roles {
		dataScope := role.DataScope

		if constants.ROLE_SCOPE_ALL == dataScope {
			break
		}

		if constants.ROLE_SCOPE_CUSTOM != dataScope {
			hasKey := false
			for _, key := range scopeKeys {
				if key == dataScope {
					hasKey = true
					break
				}
			}
			if hasKey {
				continue
			}
		}

		if constants.ROLE_SCOPE_CUSTOM == dataScope {
			sql := fmt.Sprintf(`%s.dept_id IN ( SELECT dept_id FROM sys_role_dept WHERE role_id = '%s' )`, deptAlias, role.RoleId)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_DEPT == dataScope {
			sql := fmt.Sprintf("%s.dept_id = '%s'", deptAlias, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_DEPT_CHILD == dataScope {
			sql := fmt.Sprintf(`%s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = '%s' or find_in_set('%s' , ancestors ) )`, deptAlias, userInfo.DeptId, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_SELF == dataScope {
			// 数据权限为仅本人且没有userAlias别名不查询任何数据
			if userAlias == "" {
				sql := fmt.Sprintf(`%s.dept_id = '0'`, deptAlias)
				conditions = append(conditions, sql)
			} else {
				sql := fmt.Sprintf(`%s.user_id = '%s'`, userAlias, userInfo.UserId)
				conditions = append(conditions, sql)
			}
		}

		// 记录角色范围
		scopeKeys = append(scopeKeys, dataScope)
	}

	// 构建查询条件语句
	if len(conditions) > 0 {
		dataScopeSQL = fmt.Sprintf(" AND ( %s ) ", strings.Join(conditions, " OR "))
	}
	return dataScopeSQL
}
