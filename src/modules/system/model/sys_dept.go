package model

// SysDept 部门对象 sys_dept
type SysDept struct {
	// 部门ID
	DeptID string `json:"deptId"`
	// 父部门ID
	ParentID string `json:"parentId"`
	// 祖级列表
	Ancestors string `json:"ancestors"`
	// 部门名称
	DeptName string `json:"deptName"`
	// 显示顺序
	OrderNum int `json:"orderNum"`
	// 负责人
	Leader string `json:"leader"`
	// 联系电话
	Phone string `json:"phone"`
	// 邮箱
	Email string `json:"email"`
	// 部门状态（0正常 1停用）
	Status string `json:"status"`
	// 删除标志（0代表存在 1代表删除）
	DelFlag string `json:"delFlag"`
	// 创建者
	CreateBy string `json:"createBy"`
	// 创建时间
	CreateTime int64 `json:"createTime"`
	// 更新者
	UpdateBy string `json:"updateBy"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`

	// ====== 非数据库字段属性 ======

	// 子部门列表
	Children []SysDept `json:"children,omitempty"`
}
