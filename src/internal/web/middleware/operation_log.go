package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

// 业务操作类型，对应 sys_operation_log.operation_type 列的取值契约。
const (
	OperationTypeOther  = "0" // 其它
	OperationTypeInsert = "1" // 新增
	OperationTypeUpdate = "2" // 修改
	OperationTypeDelete = "3" // 删除
	OperationTypeGrant  = "4" // 授权
	OperationTypeExport = "5" // 导出
	OperationTypeImport = "6" // 导入
	OperationTypeForce  = "7" // 强退
	OperationTypeClean  = "8" // 清空
)

// OperationLog 记录业务操作日志：请求参数脱敏后入库，响应记录摘要。
// title 为操作标题，operationType 为业务类型（空值按 OperationTypeOther）。
// 请在 AuthorizeUser 之后注册，以便获取登录用户信息。
func OperationLog(title, operationType string) gin.HandlerFunc {
	if operationType == "" {
		operationType = OperationTypeOther
	}
	return func(c *gin.Context) {
		start := time.Now()

		// 处理函数名（取包路径最后一段之后的部分）
		funcName := c.HandlerName()
		if i := strings.LastIndex(funcName, "/"); i >= 0 {
			funcName = funcName[i+1:]
		}

		ipaddr, location := reqctx.IPAddrLocation(c)

		userName := reqctx.LoginUserToUserName(c)
		if userName == "" {
			c.JSON(401, resp.CodeMsg(resp.CODE_AUTH_INVALID, "invalid login user information"))
			c.Abort()
			return
		}

		rec := OperationLogRecord{
			Title:             title,
			OperationType:     operationType,
			OperationFunction: funcName,
			OperationUrl:      c.Request.RequestURI,
			OperationMethod:   c.Request.Method,
			OperationIp:       ipaddr,
			OperationLocation: location,
			CreateBy:          userName,
		}

		// 请求参数（敏感字段脱敏后入库）
		params := reqctx.RequestParamsMap(c)
		maskSensitiveFields(params)
		if b, err := json.Marshal(params); err == nil {
			rec.OperationParam = truncateBytes(string(b), 2000)
		}

		c.Next()

		// 响应状态
		status := c.Writer.Status()
		if isHTTPSuccess(status) {
			rec.OperationStatus = constants.STATUS_YES
		} else {
			rec.OperationStatus = constants.STATUS_NO
		}

		// 响应摘要
		contentType := c.Writer.Header().Get("Content-Type")
		contentDisposition := c.Writer.Header().Get("Content-Disposition")
		rec.OperationMsg = buildOperationMsg(status, c.Writer.Size(), contentType, contentDisposition)

		rec.OperationCostTime = time.Since(start).Milliseconds()

		if OnOperationLog != nil {
			OnOperationLog(rec)
		}
	}
}

// isHTTPSuccess 判断 HTTP 状态码是否属于成功区间。
func isHTTPSuccess(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

// buildOperationMsg 生成规范的响应摘要 JSON。
func buildOperationMsg(status, size int, contentType, contentDisposition string) string {
	msg, err := json.Marshal(struct {
		Status             int    `json:"status"`
		Size               int    `json:"size"`
		ContentType        string `json:"content_type"`
		ContentDisposition string `json:"content_disposition,omitempty"`
	}{
		Status:             status,
		Size:               size,
		ContentType:        contentType,
		ContentDisposition: contentDisposition,
	})
	if err != nil {
		return ""
	}
	return string(msg)
}

// sensitiveFields 需要脱敏的字段名。
// 同时覆盖 snake_case（当前 API 契约）与历史 camelCase，避免遗漏导致明文落库。
var sensitiveFields = []string{
	"password",
	"old_password", "new_password", "confirm_password",
	"oldPassword", "newPassword", "confirmPassword",
}

// maskSensitiveFields 就地脱敏 params 中的敏感字段，并递归处理嵌套的 map / slice。
func maskSensitiveFields(params map[string]any) {
	for key, val := range params {
		if slices.Contains(sensitiveFields, key) {
			params[key] = parse.SafeContent(truncateRunes(stringify(val), 100))
			continue
		}
		maskAny(val)
	}
}

// maskAny 递归脱敏任意结构中的敏感字段。
func maskAny(v any) {
	switch t := v.(type) {
	case map[string]any:
		maskSensitiveFields(t)
	case []any:
		for _, item := range t {
			maskAny(item)
		}
	}
}

// stringify 将值转为字符串，字符串原样返回，其余用 fmt.Sprint。
func stringify(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// truncateRunes 按 rune 截断，避免切断多字节字符。
func truncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

// truncateBytes 按字节截断（用于入库前的 JSON 长度限制）。
func truncateBytes(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// OperationLogRecord 操作日志记录，解耦中间件与业务模型。
// 由业务模块通过 OnOperationLog 钩子映射为具体的持久化对象。
type OperationLogRecord struct {
	Title             string
	OperationType     string
	OperationFunction string
	OperationUrl      string
	OperationMethod   string
	OperationIp       string
	OperationLocation string
	CreateBy          string
	OperationParam    string
	OperationStatus   string
	OperationMsg      string
	OperationCostTime int64
}

// OnOperationLog 操作日志持久化钩子，由业务模块在启动时注入；为 nil 时丢弃记录。
var OnOperationLog func(OperationLogRecord)
