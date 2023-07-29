package model

// ZzOrm 测试ORM表 zz_orm
//
// https://gorm.io/zh_CN/docs/query.html
type ZzOrm struct {
	ID         uint   `gorm:"column:id;primaryKey" json:"id"`       // 测试ID
	Title      string `gorm:"column:title" json:"title"`            // 测试标题
	OrmType    string `gorm:"column:orm_type" json:"ormType"`       // orm类型
	Status     string `gorm:"column:status" json:"status"`          // 状态（0关闭 1正常）
	CreateBy   string `gorm:"column:create_by" json:"createBy"`     // 创建者
	CreateTime int64  `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateBy   string `gorm:"column:update_by" json:"updateBy"`     // 更新者
	UpdateTime int64  `gorm:"column:update_time" json:"updateTime"` // 更新时间
	Remark     string `gorm:"column:remark;size:500" json:"remark"` // 备注
}

func (ZzOrm) TableName() string {
	return "zz_orm"
}
