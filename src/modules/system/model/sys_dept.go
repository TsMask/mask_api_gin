package model

// SysDept 部门对象 sys_dept
type SysDept struct {
	DeptId     string `json:"dept_id" gorm:"column:dept_id"`                        // 部门id
	ParentId   string `json:"parent_id" gorm:"column:parent_id" binding:"required"` // 父部门id 默认0
	Ancestors  string `json:"ancestors" gorm:"column:ancestors"`                    // 祖级列表
	DeptName   string `json:"dept_name" gorm:"column:dept_name" binding:"required"` // 部门名称
	OrderNum   int64  `json:"order_num" gorm:"column:order_num"`                    // 显示顺序
	Leader     string `json:"leader" gorm:"column:leader"`                          // 负责人
	Phone      string `json:"phone" gorm:"column:phone"`                            // 联系电话
	Email      string `json:"email" gorm:"column:email"`                            // 邮箱
	Status     string `json:"status" gorm:"column:status"`                          // 部门状态（0停用 1正常）
	DelFlag    string `json:"del_flag" gorm:"column:del_flag"`                      // 删除标志（0存在 1删除）
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                // 更新时间

	// ====== 非数据库字段属性 ======

	Children []SysDept `json:"children,omitempty" gorm:"-"` // 子部门列表
}

// TableName 表名称
func (*SysDept) TableName() string {
	return "sys_dept"
}
