package model

// DemoORM 测试ORM表 demo_orm
//
// https://gorm.io/zh_CN/docs/query.html
type DemoORM struct {
	Id         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                     // 测试ID
	CreateBy   string `gorm:"column:create_by" json:"create_by"`                                // 创建者
	CreateTime int64  `gorm:"column:create_time" json:"create_time"`                            // 创建时间
	UpdateBy   string `gorm:"column:update_by" json:"update_by"`                                // 更新者
	UpdateTime int64  `gorm:"column:update_time" json:"update_time"`                            // 更新时间
	Remark     string `gorm:"column:remark" json:"remark"`                                      // 备注
	Title      string `gorm:"column:title" json:"title" binding:"required"`                     // 测试标题
	OrmType    string `gorm:"column:orm_type" json:"orm_type" binding:"required"`               // orm类型
	OrmStatus  string `gorm:"column:orm_status" json:"orm_status" binding:"required,oneof=0 1"` // 状态（0关闭 1正常）
}

func (*DemoORM) TableName() string {
	return "demo_orm"
}
