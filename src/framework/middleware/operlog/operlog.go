package operlog

import (
	"encoding/json"
	"fmt"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// 业务操作类型-其它
	BUSINESS_TYPE_OTHER = "0"

	// 业务操作类型-新增
	BUSINESS_TYPE_INSERT = "1"

	// 业务操作类型-修改
	BUSINESS_TYPE_UPDATE = "2"

	// 业务操作类型-删除
	BUSINESS_TYPE_DELETE = "3"

	// 业务操作类型-授权
	BUSINESS_TYPE_GRANT = "4"

	// 业务操作类型-导出
	BUSINESS_TYPE_EXPORT = "5"

	// 业务操作类型-导入
	BUSINESS_TYPE_IMPORT = "6"

	// 业务操作类型-强退
	BUSINESS_TYPE_FORCE = "7"

	// 业务操作类型-清空数据
	BUSINESS_TYPE_CLEAN = "8"
)

const (
	// 操作人类别-其它
	OPERATOR_TYPE_OTHER = "0"

	// 操作人类别-后台用户
	OPERATOR_TYPE_MANAGE = "1"

	// 操作人类别-手机端用户
	OPERATOR_TYPE_MOBILE = "2"
)

// Option 操作日志参数
type Option struct {
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
func OptionNew(title, businessType string) Option {
	return Option{
		Title:              title,
		BusinessType:       businessType,
		OperatorType:       OPERATOR_TYPE_OTHER,
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

// OperLog 访问操作日志记录
//
// 请在用户身份授权认证校验后使用以便获取登录用户信息
func OperLog(option Option) gin.HandlerFunc {
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
		operLog := model.SysOperLog{
			Title:         option.Title,
			BusinessType:  option.BusinessType,
			OperatorType:  option.OperatorType,
			Method:        funcName,
			OperURL:       c.Request.RequestURI,
			RequestMethod: c.Request.Method,
			OperIP:        ipaddr,
			OperLocation:  location,
			OperName:      loginUser.User.UserName,
			DeptName:      loginUser.User.Dept.DeptName,
		}

		if loginUser.User.UserType == "sys" {
			operLog.OperatorType = OPERATOR_TYPE_MANAGE
		}

		// 是否需要保存request，参数和值
		if option.IsSaveRequestData {
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

		// 是否需要保存response，参数和值
		if option.IsSaveResponseData {
			contentDisposition := c.Writer.Header().Get("Content-Disposition")
			contentType := c.Writer.Header().Get("Content-Type")
			content := contentType + contentDisposition
			msg := fmt.Sprintf(`{"status":"%d","size":"%d","content-type":"%s"}`, c.Writer.Status(), c.Writer.Size(), content)
			operLog.OperMsg = msg
		}

		// 日志记录时间
		duration := time.Since(c.GetTime("startTime"))
		operLog.CostTime = duration.Milliseconds()
		operLog.OperTime = time.Now().UnixNano() / 1e6
		operLog.Status = common.STATUS_YES

		// 保存操作记录到数据库
		service.SysOperLogImpl.InsertOperLog(operLog)
	}
}
