package model

// SysDept 部门对象 sys_dept
type SysDept struct {
	DeptID     string `json:"deptId"`                      // 部门ID
	ParentID   string `json:"parentId" binding:"required"` // 父部门ID
	Ancestors  string `json:"ancestors"`                   // 祖级列表
	DeptName   string `json:"deptName" binding:"required"` // 部门名称
	OrderNum   int    `json:"orderNum"`                    // 显示顺序
	Leader     string `json:"leader"`                      // 负责人
	Phone      string `json:"phone"`                       // 联系电话
	Email      string `json:"email"`                       // 邮箱
	Status     string `json:"status"`                      // 部门状态（0正常 1停用）
	DelFlag    string `json:"delFlag"`                     // 删除标志（0存在 1删除）
	CreateBy   string `json:"createBy"`                    // 创建者
	CreateTime int64  `json:"createTime"`                  // 创建时间
	UpdateBy   string `json:"updateBy"`                    // 更新者
	UpdateTime int64  `json:"updateTime"`                  // 更新时间

	// ====== 非数据库字段属性 ======

	Children   []SysDept `json:"children,omitempty"`   // 子部门列表
	ParentName string    `json:"parentName,omitempty"` // 父部门名称
}
