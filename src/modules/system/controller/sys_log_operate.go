package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 SysLogOperateController 结构体
var NewSysLogOperate = &SysLogOperateController{
	SysLogOperateService: service.NewSysLogOperateImpl,
}

// 操作日志记录信息
//
// PATH /system/log/operate
type SysLogOperateController struct {
	// 操作日志服务
	SysLogOperateService service.ISysLogOperate
}

// 操作日志列表
//
// GET /list
func (s *SysLogOperateController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data := s.SysLogOperateService.SelectSysLogOperatePage(querys)
	c.JSON(200, result.Ok(data))
}

// 操作日志删除
//
// DELETE /:operIds
func (s *SysLogOperateController) Remove(c *gin.Context) {
	operIds := c.Param("operIds")
	if operIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组后去重
	ids := strings.Split(operIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows := s.SysLogOperateService.DeleteSysLogOperateByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 操作日志清空
//
// DELETE /clean
func (s *SysLogOperateController) Clean(c *gin.Context) {
	err := s.SysLogOperateService.CleanSysLogOperate()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// 导出操作日志
//
// POST /export
func (s *SysLogOperateController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMap(c)
	data := s.SysLogOperateService.SelectSysLogOperatePage(querys)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysLogOperate)

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
		if row.Status == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.OperID,
			"B" + idx: row.Title,
			"C" + idx: businessType,
			"D" + idx: row.Method,
			"E" + idx: row.RequestMethod,
			"F" + idx: operatorType,
			"G" + idx: row.OperName,
			"H" + idx: row.DeptName,
			"I" + idx: row.OperURL,
			"J" + idx: row.OperIP,
			"K" + idx: row.OperLocation,
			"L" + idx: row.OperParam,
			"M" + idx: row.OperMsg,
			"N" + idx: statusValue,
			"O" + idx: row.CostTime,
			"P" + idx: date.ParseDateToStr(row.OperTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
