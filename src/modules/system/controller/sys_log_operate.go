package controller

import (
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/system/service"

	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NewSysLogOperate 实例化控制层
var NewSysLogOperate = &SysLogOperateController{
	SysLogOperateService: service.NewSysLogOperate,
}

// SysLogOperateController 操作日志记录信息
//
// PATH /system/log/operate
type SysLogOperateController struct {
	SysLogOperateService *service.SysLogOperate // 操作日志服务
}

// List 操作日志列表
//
// GET /list
func (s SysLogOperateController) List(c *gin.Context) {
	query := context.QueryMap(c)
	rows, total := s.SysLogOperateService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Clean 操作日志清空
//
// DELETE /clean
func (s SysLogOperateController) Clean(c *gin.Context) {
	rows := s.SysLogOperateService.Clean()
	c.JSON(200, response.OkData(rows))
}

// Export 导出操作日志
//
// GET /export
func (s SysLogOperateController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := context.QueryMap(c)
	rows, total := s.SysLogOperateService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("sys_log_operate_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "操作序号",
		"B1": "操作模块",
		"C1": "业务类型",
		"D1": "请求方法",
		"E1": "请求方式",
		"F1": "操作类别",
		"G1": "操作人员",
		"H1": "部门名称",
		"I1": "请求地址",
		"J1": "操作地址",
		"K1": "操作地点",
		"L1": "请求参数",
		"M1": "操作消息",
		"N1": "状态",
		"O1": "消耗时间（毫秒）",
		"P1": "操作时间",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 业务类型
		businessType := ""
		// 操作类别
		operatorType := ""
		// 状态
		statusValue := "失败"
		if row.StatusFlag == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.ID,
			"B" + idx: row.Title,
			"C" + idx: businessType,
			"D" + idx: row.OperaMethod,
			"E" + idx: row.OperaUrlMethod,
			"F" + idx: operatorType,
			"G" + idx: row.OperaBy,
			"I" + idx: row.OperaUrl,
			"J" + idx: row.OperaIp,
			"K" + idx: row.OperaLocation,
			"L" + idx: row.OperaParam,
			"M" + idx: row.OperaMsg,
			"N" + idx: statusValue,
			"O" + idx: row.CostTime,
			"P" + idx: date.ParseDateToStr(row.OperaTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
