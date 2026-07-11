package model

// SysDept 部门表（树形: parent_id 自关联 + ancestors 物化祖级路径）
type SysDept struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`        // 部门ID
	CreateBy   string `json:"createBy" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"createTime" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"updateTime" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                         // 备注
	ParentId   int64  `json:"parentId" gorm:"column:parent_id"`                    // 父部门ID 默认0
	Ancestors  string `json:"ancestors" gorm:"column:ancestors"`                   // 祖级列表
	DeptName   string `json:"deptName" gorm:"column:dept_name" binding:"required"` // 部门名称
	DeptLeader string `json:"deptLeader" gorm:"column:dept_leader"`                // 负责人
	DeptPhone  string `json:"deptPhone" gorm:"column:dept_phone"`                  // 联系电话
	DeptEmail  string `json:"deptEmail" gorm:"column:dept_email"`                  // 邮箱
	DeptSort   int64  `json:"deptSort" gorm:"column:dept_sort"`                    // 显示顺序
	DeptStatus string `json:"deptStatus" gorm:"column:dept_status"`                // 部门状态（0停用 1正常）
	DeptDel    string `json:"-" gorm:"column:dept_del"`                            // 删除标记（0存在 1删除）

	// ====== 非数据库字段属性 ======

	Children []SysDept `json:"children,omitempty" gorm:"-"` // 子部门列表
}

// TableName 表名称
func (*SysDept) TableName() string {
	return "sys_dept"
}
