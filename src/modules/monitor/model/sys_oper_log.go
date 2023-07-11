package model

// SysOperLog 操作日志记录表 sys_oper_log
type SysOperLog struct {
	// 日志主键
	OperID string `json:"operId"`
	// 模块标题
	Title string `json:"title"`
	// 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	BusinessType string `json:"businessType"`
	// 方法名称
	Method string `json:"method"`
	// 请求方式
	RequestMethod string `json:"requestMethod"`
	// 操作类别（0其它 1后台用户 2手机端用户）
	OperatorType string `json:"operatorType"`
	// 操作人员
	OperName string `json:"operName"`
	// 部门名称
	DeptName string `json:"deptName"`
	// 请求URL
	OperURL string `json:"operUrl"`
	// 主机地址
	OperIP string `json:"operIp"`
	// 操作地点
	OperLocation string `json:"operLocation"`
	// 请求参数
	OperParam string `json:"operParam"`
	// 操作消息
	OperMsg string `json:"operMsg"`
	// 操作状态（0异常 1正常）
	Status string `json:"status"`
	// 操作时间
	OperTime int64 `json:"operTime"`
	// 消耗时间（毫秒）
	CostTime int64 `json:"costTime"`
}
