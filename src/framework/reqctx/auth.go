package reqctx

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/token"
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
	return info.User.UserName
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
	roles := info.User.Roles
	for _, item := range roles {
		if item.RoleKey == target {
			return true
		}
	}
	return false
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
	for _, str := range perms {
		if str == target {
			return true
		}
	}
	return false
}

// LoginUserToDataScopeSQL 登录用户信息-角色数据范围过滤SQL字符串
// deptAlias 数据范围过滤目标表的部门别名/表名；userAlias 仅本人数据权限时用户表别名/表名，空则回退部门维度
func LoginUserToDataScopeSQL(c *gin.Context, deptAlias string, userAlias string) string {
	dataScopeSQL := ""
	// 登录用户信息
	info, err := LoginUser(c)
	if err != nil {
		return dataScopeSQL
	}
	userInfo := info.User

	// 数据范围过滤引用的部门列：目标表为 sys_dept 时为其主键 id，其余表为外键 dept_id
	deptRefCol := "dept_id"
	if deptAlias == "sys_dept" {
		deptRefCol = "id"
	}

	// 如果是系统管理员，则不过滤数据
	if config.IsSystemUser(userInfo.ID) {
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
		dataScope := role.RoleDataScope

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
			sql := fmt.Sprintf(`%s.%s IN
			( SELECT dept_id FROM sys_role_dept WHERE role_id = %d )
			AND %s.%s NOT IN
			(
			SELECT d.parent_id FROM sys_dept d
			INNER JOIN sys_role_dept rd ON rd.dept_id = d.id
			AND rd.role_id = %d
			)`, deptAlias, deptRefCol, role.ID, deptAlias, deptRefCol, role.ID)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_DEPT == dataScope {
			sql := fmt.Sprintf("%s.%s = %d", deptAlias, deptRefCol, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_DEPT_CHILD == dataScope {
			sql := fmt.Sprintf("%s.%s IN ( SELECT id FROM sys_dept WHERE id = %d OR find_in_set(%d, ancestors ) )", deptAlias, deptRefCol, userInfo.DeptId, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_SELF == dataScope {
			if userAlias == "" {
				sql := fmt.Sprintf("%s.%s = %d", deptAlias, deptRefCol, userInfo.DeptId)
				conditions = append(conditions, sql)
			} else {
				sql := fmt.Sprintf("%s.id = %d", userAlias, userInfo.ID)
				conditions = append(conditions, sql)
			}
		}

		// 记录角色范围
		scopeKeys = append(scopeKeys, dataScope)
	}

	// 构建查询条件语句
	if len(conditions) > 0 {
		dataScopeSQL = fmt.Sprintf(" ( %s ) ", strings.Join(conditions, " OR "))
	}
	return dataScopeSQL
}
