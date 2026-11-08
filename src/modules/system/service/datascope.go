package service

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/modules/system/model"
)

// DataScopeSQL 登录用户信息-角色数据范围过滤SQL字符串
// deptAlias 数据范围过滤目标表的部门别名/表名；userAlias 仅本人数据权限时用户表别名/表名，空则回退部门维度
func DataScopeSQL(c *gin.Context, deptAlias string, userAlias string) string {
	dataScopeSQL := ""
	// 登录用户信息
	info, err := reqctx.LoginUser(c)
	if err != nil {
		return dataScopeSQL
	}
	var userInfo model.SysUser
	if len(info.User) > 0 {
		_ = json.Unmarshal(info.User, &userInfo)
	}

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
			if slices.Contains(scopeKeys, dataScope) {
				continue
			}
		}

		if constants.ROLE_SCOPE_CUSTOM == dataScope {
			sql := fmt.Sprintf(`%s.%s IN
			( SELECT dept_id FROM sys_role_dept WHERE role_id = %d )
			AND %s.%s NOT IN (
			SELECT c.ancestor_id FROM sys_dept_closure c
			INNER JOIN sys_role_dept rd ON rd.dept_id = c.descendant_id
			AND rd.role_id = %d WHERE c.depth = 1
			)`, deptAlias, deptRefCol, role.ID, deptAlias, deptRefCol, role.ID)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_DEPT == dataScope {
			sql := fmt.Sprintf("%s.%s = %d", deptAlias, deptRefCol, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constants.ROLE_SCOPE_DEPT_CHILD == dataScope {
			sql := fmt.Sprintf(`%s.%s IN (
			SELECT descendant_id FROM sys_dept_closure WHERE ancestor_id = %d
			)`, deptAlias, deptRefCol, userInfo.DeptId)
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

// UserVisibleSQL 登录用户信息-用户查询可见范围过滤SQL字符串
func UserVisibleSQL(c *gin.Context, userAlias string) string {
	dataScopeSQL := DataScopeSQL(c, userAlias, userAlias)
	if config.IsSystemUser(reqctx.LoginUserToUserID(c)) {
		return dataScopeSQL
	}

	hiddenUserIds := config.SystemUserIDs()
	if len(hiddenUserIds) <= 0 {
		return dataScopeSQL
	}

	values := make([]string, 0, len(hiddenUserIds))
	for _, id := range hiddenUserIds {
		values = append(values, fmt.Sprintf("%d", id))
	}
	hiddenSQL := fmt.Sprintf("%s.id not in (%s)", userAlias, strings.Join(values, ","))
	if dataScopeSQL == "" {
		return hiddenSQL
	}
	return fmt.Sprintf("(%s and %s)", dataScopeSQL, hiddenSQL)
}
