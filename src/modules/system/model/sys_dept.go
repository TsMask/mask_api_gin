package model

// SysDept 部门表
type SysDept struct {
	DeptId     string `json:"deptId" gorm:"column:dept_id;primaryKey;type:int;autoIncrement"` // 部门ID
	ParentId   string `json:"parentId" gorm:"column:parent_id" binding:"required"`            // 父部门ID 默认0
	Ancestors  string `json:"ancestors" gorm:"column:ancestors"`                              // 祖级列表
	DeptName   string `json:"deptName" gorm:"column:dept_name" binding:"required"`            // 部门名称
	DeptSort   int64  `json:"deptSort" gorm:"column:dept_sort"`                               // 显示顺序
	Leader     string `json:"leader" gorm:"column:leader"`                                    // 负责人
	Phone      string `json:"phone" gorm:"column:phone"`                                      // 联系电话
	Email      string `json:"email" gorm:"column:email"`                                      // 邮箱
	StatusFlag string `json:"statusFlag" gorm:"column:status_flag"`                           // 部门状态（0停用 1正常）
	DelFlag    string `json:"-" gorm:"column:del_flag"`                                       // 删除标记（0存在 1删除）
	CreateBy   string `json:"createBy" gorm:"column:create_by"`                               // 创建者
	CreateTime int64  `json:"createTime" gorm:"column:create_time"`                           // 创建时间
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`                               // 更新者
	UpdateTime int64  `json:"updateTime" gorm:"column:update_time"`                           // 更新时间

	// ====== 非数据库字段属性 ======

	Children []SysDept `json:"children,omitempty" gorm:"-"` // 子部门列表
}

// TableName 表名称
func (*SysDept) TableName() string {
	return "sys_dept"
}
