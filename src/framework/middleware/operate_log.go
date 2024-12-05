package middleware

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// BUSINESS_TYPE_OTHER 业务操作类型-其它
	BUSINESS_TYPE_OTHER = "0"

	// BUSINESS_TYPE_INSERT 业务操作类型-新增
	BUSINESS_TYPE_INSERT = "1"

	// BUSINESS_TYPE_UPDATE 业务操作类型-修改
	BUSINESS_TYPE_UPDATE = "2"

	// BUSINESS_TYPE_DELETE 业务操作类型-删除
	BUSINESS_TYPE_DELETE = "3"

	// BUSINESS_TYPE_GRANT 业务操作类型-授权
	BUSINESS_TYPE_GRANT = "4"

	// BUSINESS_TYPE_EXPORT 业务操作类型-导出
	BUSINESS_TYPE_EXPORT = "5"

	// BUSINESS_TYPE_IMPORT 业务操作类型-导入
	BUSINESS_TYPE_IMPORT = "6"

	// BUSINESS_TYPE_FORCE 业务操作类型-强退
	BUSINESS_TYPE_FORCE = "7"

	// BUSINESS_TYPE_CLEAN 业务操作类型-清空数据
	BUSINESS_TYPE_CLEAN = "8"
)

// Options Option 操作日志参数
type Options struct {
	Title              string `json:"title"`              // 标题
	BusinessType       string `json:"businessType"`       // 类型，默认常量 BUSINESS_TYPE_OTHER
	IsSaveRequestData  bool   `json:"isSaveRequestData"`  // 是否保存请求的参数
	IsSaveResponseData bool   `json:"isSaveResponseData"` // 是否保存响应的参数
}

// OptionNew 操作日志参数默认值
//
// 标题 "title":"--"
//
// 类型 "businessType": BUSINESS_TYPE_OTHER
//
// 注意之后JSON反序列使用：c.ShouldBindBodyWithJSON(&params)
func OptionNew(title, businessType string) Options {
	return Options{
		Title:              title,
		BusinessType:       businessType,
		IsSaveRequestData:  true,
		IsSaveResponseData: true,
	}
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
			c.JSON(401, response.CodeMsg(401, "无效身份授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 操作日志记录
		operaLog := model.SysLogOperate{
			Title:          options.Title,
			BusinessType:   options.BusinessType,
			OperaMethod:    funcName,
			OperaUrl:       c.Request.RequestURI,
			OperaUrlMethod: c.Request.Method,
			OperaIp:        ipaddr,
			OperaLocation:  location,
			OperaBy:        loginUser.User.UserName,
		}

		// 是否需要保存request，参数和值
		if options.IsSaveRequestData {
			params := ctx.RequestParamsMap(c)
			// 敏感属性字段进行掩码
			processSensitiveFields(params)
			jsonStr, _ := json.Marshal(params)
			paramsStr := string(jsonStr)
			if len(paramsStr) > 2000 {
				paramsStr = paramsStr[:2000]
			}
			operaLog.OperaParam = paramsStr
		}

		// 调用下一个处理程序
		c.Next()

		// 响应状态
		status := c.Writer.Status()
		if status == 200 {
			operaLog.StatusFlag = constants.STATUS_YES
		} else {
			operaLog.StatusFlag = constants.STATUS_NO
		}

		// 是否需要保存response，参数和值
		if options.IsSaveResponseData {
			contentDisposition := c.Writer.Header().Get("Content-Disposition")
			contentType := c.Writer.Header().Get("Content-Type")
			content := contentType + contentDisposition
			msg := fmt.Sprintf(`{"status":"%d","size":"%d","content-type":"%s"}`, status, c.Writer.Size(), content)
			operaLog.OperaMsg = msg
		}

		// 日志记录时间
		duration := time.Since(c.GetTime("startTime"))
		operaLog.CostTime = duration.Milliseconds()
		operaLog.OperaTime = time.Now().UnixMilli()

		// 保存操作记录到数据库
		service.NewSysLogOperate.Insert(operaLog)
	}
}

// 敏感属性字段进行掩码
var maskProperties = []string{
	"password",
	"oldPassword",
	"newPassword",
	"confirmPassword",
}

// processSensitiveFields 处理敏感属性字段
func processSensitiveFields(obj interface{}) {
	val := reflect.ValueOf(obj)

	switch val.Kind() {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			keyStr := key.Interface().(string)

			// 遍历是否敏感属性
			hasMaskKey := false
			for _, v := range maskProperties {
				if v == keyStr {
					hasMaskKey = true
					break
				}
			}

			if hasMaskKey {
				valueStr := value.Interface().(string)
				if len(valueStr) > 100 {
					valueStr = valueStr[0:100]
				}
				val.SetMapIndex(key, reflect.ValueOf(parse.SafeContent(valueStr)))
			} else {
				processSensitiveFields(value.Interface())
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			processSensitiveFields(val.Index(i).Interface())
		}
		// default:
		// 	logger.Infof("processSensitiveFields unhandled case %v", val.Kind())
	}
}
