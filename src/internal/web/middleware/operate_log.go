package middleware

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

// 业务操作类型，对应 sys_oper_log.business_type 列的取值契约。
const (
	OperateTypeOther  = "0" // 其它
	OperateTypeInsert = "1" // 新增
	OperateTypeUpdate = "2" // 修改
	OperateTypeDelete = "3" // 删除
	OperateTypeGrant  = "4" // 授权
	OperateTypeExport = "5" // 导出
	OperateTypeImport = "6" // 导入
	OperateTypeForce  = "7" // 强退
	OperateTypeClean  = "8" // 清空
)

// OperateLog 记录业务操作日志：请求参数脱敏后入库，响应记录摘要。
// title 为操作标题，operateType 为业务类型（空值按 OperateTypeOther）。
// 请在 AuthorizeUser 之后注册，以便获取登录用户信息。
func OperateLog(title, operateType string) gin.HandlerFunc {
	if operateType == "" {
		operateType = OperateTypeOther
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

		rec := OperateLogRecord{
			Title:           title,
			OperateType:     operateType,
			OperateFunction: funcName,
			OperateUrl:      c.Request.RequestURI,
			OperateMethod:   c.Request.Method,
			OperateIp:       ipaddr,
			OperateLocation: location,
			CreateBy:        userName,
		}

		// 请求参数（敏感字段脱敏后入库）
		params := reqctx.RequestParamsMap(c)
		maskSensitiveFields(params)
		if b, err := json.Marshal(params); err == nil {
			rec.OperateParam = truncateBytes(string(b), 2000)
		}

		c.Next()

		// 响应状态
		status := c.Writer.Status()
		if status == 200 {
			rec.OperateStatus = constants.STATUS_YES
		} else {
			rec.OperateStatus = constants.STATUS_NO
		}

		// 响应摘要
		contentType := c.Writer.Header().Get("Content-Type")
		contentDisposition := c.Writer.Header().Get("Content-Disposition")
		rec.OperateMsg = fmt.Sprintf(
			`{"status":"%d","size":%d,"content-type":"%s"}`,
			status, c.Writer.Size(), contentType+contentDisposition,
		)

		rec.OperateCostTime = time.Since(start).Milliseconds()

		if OnOperateLog != nil {
			OnOperateLog(rec)
		}
	}
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

// OperateLogRecord 操作日志记录，解耦中间件与业务模型。
// 由业务模块通过 OnOperateLog 钩子映射为具体的持久化对象。
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

// OnOperateLog 操作日志持久化钩子，由业务模块在启动时注入；为 nil 时丢弃记录。
var OnOperateLog func(OperateLogRecord)
