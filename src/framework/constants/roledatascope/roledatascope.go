package roledatascope

// 系统角色数据范围常量

const (
	// 全部数据权限
	ALL = "1"

	// 自定数据权限
	CUSTOM = "2"

	// 部门数据权限
	DEPT = "3"

	// 部门及以下数据权限
	DEPT_AND_CHILD = "4"

	// 仅本人数据权限
	SELF = "5"
)

// 系统角色数据范围映射
var RoleDataScope = map[string]string{
	ALL:            "全部数据权限",
	CUSTOM:         "自定数据权限",
	DEPT:           "部门数据权限",
	DEPT_AND_CHILD: "部门及以下数据权限",
	SELF:           "仅本人数据权限",
}
