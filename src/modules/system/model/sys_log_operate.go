package model

// SysLogOperate 系统操作日志表（用 create_by 字符串软关联 sys_user.user_name）
type SysLogOperate struct {
	ID              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`    // 操作ID
	CreateBy        string `json:"createBy" gorm:"column:create_by"`                // 创建者
	CreateTime      int64  `json:"createTime" gorm:"column:create_time"`            // 创建时间
	Title           string `json:"title" gorm:"column:title"`                       // 模块标题
	OperateType     string `json:"operateType" gorm:"column:operate_type"`          // 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	OperateUrl      string `json:"operateUrl" gorm:"column:operate_url"`            // 请求URL
	OperateMethod   string `json:"operateMethod" gorm:"column:operate_method"`      // 请求方式
	OperateIp       string `json:"operateIp" gorm:"column:operate_ip"`              // 主机地址
	OperateLocation string `json:"operateLocation" gorm:"column:operate_location"`  // 操作地点
	OperateParam    string `json:"operateParam" gorm:"column:operate_param"`        // 请求参数
	OperateFunction string `json:"operateFunction" gorm:"column:operate_function"`  // 方法函数方法
	OperateCostTime int64  `json:"operateCostTime" gorm:"column:operate_cost_time"` // 消耗时间（毫秒）
	OperateStatus   string `json:"operateStatus" gorm:"column:operate_status"`      // 操作状态（0异常 1正常）
	OperateMsg      string `json:"operateMsg" gorm:"column:operate_msg"`            // 操作消息
}

// TableName 表名称
func (*SysLogOperate) TableName() string {
	return "sys_log_operate"
}
