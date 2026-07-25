package middleware

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

const (
	// OPERATE_TYPE_OTHER 业务操作类型-其它
	OPERATE_TYPE_OTHER = "0"

	// OPERATE_TYPE_INSERT 业务操作类型-新增
	OPERATE_TYPE_INSERT = "1"

	// OPERATE_TYPE_UPDATE 业务操作类型-修改
	OPERATE_TYPE_UPDATE = "2"

	// OPERATE_TYPE_DELETE 业务操作类型-删除
	OPERATE_TYPE_DELETE = "3"

	// OPERATE_TYPE_GRANT 业务操作类型-授权
	OPERATE_TYPE_GRANT = "4"

	// OPERATE_TYPE_EXPORT 业务操作类型-导出
	OPERATE_TYPE_EXPORT = "5"

	// OPERATE_TYPE_IMPORT 业务操作类型-导入
	OPERATE_TYPE_IMPORT = "6"

	// OPERATE_TYPE_FORCE 业务操作类型-强退
	OPERATE_TYPE_FORCE = "7"

	// OPERATE_TYPE_CLEAN 业务操作类型-清空数据
	OPERATE_TYPE_CLEAN = "8"
)

// Options Option 操作日志参数
type Options struct {
	Title              string `json:"title"`              // 标题
	OperateType        string `json:"operateType"`        // 类型，默认常量 OPERATE_TYPE_OTHER
	IsSaveRequestData  bool   `json:"isSaveRequestData"`  // 是否保存请求的参数
	IsSaveResponseData bool   `json:"isSaveResponseData"` // 是否保存响应的参数
}

// OptionNew 操作日志参数默认值
//
// 标题 "title":"--"
//
// 类型 "operateType": OPERATE_TYPE_OTHER
//
// 注意之后JSON反序列使用：c.ShouldBindBodyWithJSON(&params)
func OptionNew(title, operateType string) Options {
	return Options{
		Title:              title,
		OperateType:        operateType,
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
		ipaddr, location := reqctx.IPAddrLocation(c)

		// 获取登录用户信息
		userName := reqctx.LoginUserToUserName(c)
		if userName == "" {
			c.JSON(401, resp.CodeMsg(resp.CODE_AUTH_INVALID, "invalid login user information"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 操作日志记录
		operaLog := OperateLogRecord{
			Title:           options.Title,
			OperateType:     options.OperateType,
			OperateFunction: funcName,
			OperateUrl:      c.Request.RequestURI,
			OperateMethod:   c.Request.Method,
			OperateIp:       ipaddr,
			OperateLocation: location,
			CreateBy:        userName,
		}

		// 是否需要保存request，参数和值
		if options.IsSaveRequestData {
			params := reqctx.RequestParamsMap(c)
			// 敏感属性字段进行掩码
			processSensitiveFields(params)
			jsonStr, _ := json.Marshal(params)
			paramsStr := string(jsonStr)
			if len(paramsStr) > 2000 {
				paramsStr = paramsStr[:2000]
			}
			operaLog.OperateParam = paramsStr
		}

		// 调用下一个处理程序
		c.Next()

		// 响应状态
		status := c.Writer.Status()
		if status == 200 {
			operaLog.OperateStatus = constants.STATUS_YES
		} else {
			operaLog.OperateStatus = constants.STATUS_NO
		}

		// 是否需要保存response，参数和值
		if options.IsSaveResponseData {
			contentDisposition := c.Writer.Header().Get("Content-Disposition")
			contentType := c.Writer.Header().Get("Content-Type")
			content := contentType + contentDisposition
			msg := fmt.Sprintf(`{"status":"%d","size":%d,"content-type":"%s"}`, status, c.Writer.Size(), content)
			operaLog.OperateMsg = msg
		}

		// 日志记录时间
		duration := time.Since(c.GetTime("startTime"))
		operaLog.OperateCostTime = duration.Milliseconds()

		// 保存操作记录到数据库
		if OnOperateLog != nil {
			OnOperateLog(operaLog)
		}
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
func processSensitiveFields(obj any) {
	val := reflect.ValueOf(obj)

	switch val.Kind() {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			keyStr := key.Interface().(string)

			// 遍历是否敏感属性
			if slices.Contains(maskProperties, keyStr) {
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

// OperateLogRecord decouples the operate-log middleware from the system business model.
type OperateLogRecord struct {
	Title           string
	OperateType     string
	OperateFunction string
	OperateUrl      string
	OperateMethod   string
	OperateIp       string
	OperateLocation string
	CreateBy        string
	OperateParam    string
	OperateStatus   string
	OperateMsg      string
	OperateCostTime int64
}

// OnOperateLog persists an operate-log record. Injected by the business module at startup; nil drops the record.
var OnOperateLog func(OperateLogRecord)
