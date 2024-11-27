package constants

// 系统角色数据范围常量
const (
	// ROLE_SCOPE_ALL 全部数据权限
	ROLE_SCOPE_ALL = "1"
	// ROLE_SCOPE_CUSTOM 自定数据权限
	ROLE_SCOPE_CUSTOM = "2"
	// ROLE_SCOPE_DEPT 部门数据权限
	ROLE_SCOPE_DEPT = "3"
	// ROLE_SCOPE_DEPT_CHILD 部门及以下数据权限
	ROLE_SCOPE_DEPT_CHILD = "4"
	// ROLE_SCOPE_SELF 仅本人数据权限
	ROLE_SCOPE_SELF = "5"
)

// ROLE_SCOPE_DATA 系统角色数据范围映射
var ROLE_SCOPE_DATA = map[string]string{
	ROLE_SCOPE_ALL:        "全部数据权限",
	ROLE_SCOPE_CUSTOM:     "自定数据权限",
	ROLE_SCOPE_DEPT:       "部门数据权限",
	ROLE_SCOPE_DEPT_CHILD: "部门及以下数据权限",
	ROLE_SCOPE_SELF:       "仅本人数据权限",
}
