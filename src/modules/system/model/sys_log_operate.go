package model

// SysLogOperate 系统操作日志表 sys_log_operate
type SysLogOperate struct {
	OperaID       string `json:"operaId"`       // 日志主键
	Title         string `json:"title"`         // 模块标题
	BusinessType  string `json:"businessType"`  // 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	Method        string `json:"method"`        // 方法名称
	RequestMethod string `json:"requestMethod"` // 请求方式
	OperatorType  string `json:"operatorType"`  // 操作人员类别（0其它 1后台用户 2手机端用户）
	OperaName     string `json:"operaName"`     // 操作人员
	DeptName      string `json:"deptName"`      // 部门名称
	OperaURL      string `json:"operaUrl"`      // 请求URL
	OperaIP       string `json:"operaIp"`       // 操作地址
	OperaLocation string `json:"operaLocation"` // 操作地点
	OperaParam    string `json:"operaParam"`    // 请求参数
	OperaMsg      string `json:"operaMsg"`      // 操作消息
	Status        string `json:"status"`        // 操作状态（0异常 1正常）
	OperaTime     int64  `json:"operaTime"`     // 操作时间
	CostTime      int64  `json:"costTime"`      // 消耗时间（毫秒）
}
