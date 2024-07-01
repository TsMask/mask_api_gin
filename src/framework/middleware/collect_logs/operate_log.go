package collect_logs

import (
	"encoding/json"
	"fmt"
	constCommon "mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// BusinessTypeOther 业务操作类型-其它
	BusinessTypeOther = "0"

	// BusinessTypeInsert 业务操作类型-新增
	BusinessTypeInsert = "1"

	// BusinessTypeUpdate 业务操作类型-修改
	BusinessTypeUpdate = "2"

	// BusinessTypeDelete 业务操作类型-删除
	BusinessTypeDelete = "3"

	// BusinessTypeGrant 业务操作类型-授权
	BusinessTypeGrant = "4"

	// BusinessTypeExport 业务操作类型-导出
	BusinessTypeExport = "5"

	// BusinessTypeImport 业务操作类型-导入
	BusinessTypeImport = "6"

	// BusinessTypeForce 业务操作类型-强退
	BusinessTypeForce = "7"

	// BusinessTypeClean 业务操作类型-清空数据
	BusinessTypeClean = "8"
)

const (
	// OperatorTypeOther 操作人类别-其它
	OperatorTypeOther = "0"

	// OperatorTypeManage 操作人类别-后台用户
	OperatorTypeManage = "1"

	// OperatorTypeMobile 操作人类别-手机端用户
	OperatorTypeMobile = "2"
)

// Options Option 操作日志参数
type Options struct {
	Title              string `json:"title"`              // 标题
	BusinessType       string `json:"businessType"`       // 类型，默认常量 BUSINESS_TYPE_OTHER
	OperatorType       string `json:"operatorType"`       // 操作人类别，默认常量 OPERATOR_TYPE_OTHER
	IsSaveRequestData  bool   `json:"isSaveRequestData"`  // 是否保存请求的参数
	IsSaveResponseData bool   `json:"isSaveResponseData"` // 是否保存响应的参数
}

// OptionNew 操作日志参数默认值
//
// 标题 "title":"--"
//
// 类型 "businessType": BUSINESS_TYPE_OTHER
//
// 注意之后JSON反序列使用：c.ShouldBindBodyWith(&params, binding.JSON)
func OptionNew(title, businessType string) Options {
	return Options{
		Title:              title,
		BusinessType:       businessType,
		OperatorType:       OperatorTypeOther,
		IsSaveRequestData:  true,
		IsSaveResponseData: true,
	}
}

// 敏感属性字段进行掩码
var maskProperties []string = []string{
	"password",
	"oldPassword",
	"newPassword",
	"confirmPassword",
}

// OperateLog 访问操作日志记录
//
// 请在用户身份授权认证校验后使用以便获取登录用户信息
func OperateLog(options Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("startTime", time.Now())

		// 函数名
		funcName := c.HandlerName()
		lastDotIndex := strings.LastIndex(funcName, "/")
		funcName = funcName[lastDotIndex+1:]

		// 解析ip地址
		ipaddr, location := ctx.IPAddrLocation(c)

		// 获取登录用户信息
		loginUser, err := ctx.LoginUser(c)
		if err != nil {
			c.JSON(401, result.CodeMsg(401, "无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 操作日志记录
		operLog := model.SysLogOperate{
			Title:         options.Title,
			BusinessType:  options.BusinessType,
			OperatorType:  options.OperatorType,
			Method:        funcName,
			OperURL:       c.Request.RequestURI,
			RequestMethod: c.Request.Method,
			OperIP:        ipaddr,
			OperLocation:  location,
			OperName:      loginUser.User.UserName,
			DeptName:      loginUser.User.Dept.DeptName,
		}

		if loginUser.User.UserType == "sys" {
			operLog.OperatorType = OperatorTypeManage
		}

		// 是否需要保存request，参数和值
		if options.IsSaveRequestData {
			params := ctx.RequestParamsMap(c)
			for k, v := range params {
				// 敏感属性字段进行掩码
				for _, s := range maskProperties {
					if s == k {
						params[k] = parse.SafeContent(v.(string))
						break
					}
				}
			}
			jsonStr, _ := json.Marshal(params)
			paramsStr := string(jsonStr)
			if len(paramsStr) > 2000 {
				paramsStr = paramsStr[:2000]
			}
			operLog.OperParam = paramsStr
		}

		// 调用下一个处理程序
		c.Next()

		// 响应状态
		status := c.Writer.Status()
		if status == 200 {
			operLog.Status = constCommon.StatusYes
		} else {
			operLog.Status = constCommon.StatusNo
		}

		// 是否需要保存response，参数和值
		if options.IsSaveResponseData {
			contentDisposition := c.Writer.Header().Get("Content-Disposition")
			contentType := c.Writer.Header().Get("Content-Type")
			content := contentType + contentDisposition
			msg := fmt.Sprintf(`{"status":"%d","size":"%d","content-type":"%s"}`, status, c.Writer.Size(), content)
			operLog.OperMsg = msg
		}

		// 日志记录时间
		duration := time.Since(c.GetTime("startTime"))
		operLog.CostTime = duration.Milliseconds()
		operLog.OperTime = time.Now().UnixMilli()

		// 保存操作记录到数据库
		service.NewSysLogOperateImpl.InsertSysLogOperate(operLog)
	}
}
