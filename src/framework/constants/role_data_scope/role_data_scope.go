package role_data_scope

// 系统角色数据范围常量

const (
	// All 全部数据权限
	All = "1"

	// Custom 自定数据权限
	Custom = "2"

	// Dept 部门数据权限
	Dept = "3"

	// DeptChild 部门及以下数据权限
	DeptChild = "4"

	// Self 仅本人数据权限
	Self = "5"
)

// RoleDataScope 系统角色数据范围映射
var RoleDataScope = map[string]string{
	All:       "全部数据权限",
	Custom:    "自定数据权限",
	Dept:      "部门数据权限",
	DeptChild: "部门及以下数据权限",
	Self:      "仅本人数据权限",
}
