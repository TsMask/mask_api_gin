package controller

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/internal/web/token"
)

// NewSysUserOnline 实例化控制层
var NewSysUserOnline = &SysUserOnlineController{}

// SysUserOnlineController 在线用户信息 控制层处理
//
// PATH /monitor/user-online
type SysUserOnlineController struct{}

// List 在线用户列表
//
// GET /list
func (s SysUserOnlineController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)

	// 取全部会话 key，排序保证顺序确定性
	keys, _ := cache.GetKeys("", constants.CACHE_TOKEN_DEVICE+":*")
	sort.Strings(keys)

	if len(keys) == 0 {
		c.JSON(200, resp.OkData(map[string]any{"total": 0, "rows": []any{}}))
		return
	}

	// 全量取出并解析：过滤/排序需基于完整数据，否则分页会漏命中
	values, _ := cache.GetBatch("", keys)
	list := make([]map[string]any, 0, len(values))
	for _, v := range values {
		str := fmt.Sprint(v)
		if str == "" {
			continue
		}
		var info token.UserInfo
		if err := json.Unmarshal([]byte(str), &info); err != nil {
			continue
		}
		if info.DeviceId == "" {
			continue
		}
		list = append(list, userOnlineRow(info))
	}

	// 过滤（基于全部数据）
	list = filterContains(list, "device_id", query["device_id"])
	list = filterContains(list, "user_name", query["user_name"])
	list = filterContains(list, "login_ip", query["login_ip"])
	list = filterContains(list, "login_location", query["login_location"])
	list = filterContains(list, "login_browser", query["browser"])
	list = filterContains(list, "login_os", query["os"])

	// 排序，格式: sort=loginTime,userName:asc（不加方向默认 DESC）
	if sortStr, ok := query["sort"]; ok && sortStr != "" {
		sort.SliceStable(list, func(i, j int) bool {
			return compareUserOnline(list[i], list[j], sortStr)
		})
	} else {
		// 默认登录时间降序
		sort.SliceStable(list, func(i, j int) bool {
			return parse.Number(list[i]["login_time"]) > parse.Number(list[j]["login_time"])
		})
	}

	// 分页（在过滤、排序之后，与 DB 列表一致；total 为过滤后的真实条数）
	total := int64(len(list))
	pn := int(parse.Number(query["page_num"]))
	if pn < 1 {
		pn = 1
	}
	ps := int(parse.Number(query["page_size"]))
	if ps < 1 {
		ps = 10
	}
	start := (pn - 1) * ps
	end := pn * ps
	if start > len(list) {
		start = len(list)
	}
	if end > len(list) {
		end = len(list)
	}

	c.JSON(200, resp.OkData(map[string]any{
		"total": total,
		"rows":  list[start:end],
	}))
}

// Logout 在线用户强制退出
//
// DELETE /logout/:device_id
func (s SysUserOnlineController) Logout(c *gin.Context) {
	deviceIdStr := c.Param("device_id")
	if deviceIdStr == "" || strings.Contains(deviceIdStr, "*") {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: deviceId is empty"))
		return
	}

	uniqueIDs := parse.RemoveDuplicatesToArray(deviceIdStr, ",")
	for _, v := range uniqueIDs {
		key := constants.CACHE_TOKEN_DEVICE + ":" + v
		if err := cache.Del("", key); err != nil {
			c.JSON(200, resp.ErrMsg(err.Error()))
			return
		}
	}

	c.JSON(200, resp.Ok(nil))
}

// userOnlineRow 将令牌信息转为列表行
func userOnlineRow(info token.UserInfo) map[string]any {
	return map[string]any{
		"device_id":      info.DeviceId,
		"user_id":        info.UserId,
		"user_name":      info.UserName,
		"login_ip":       info.LoginIp,
		"login_location": info.LoginLocation,
		"login_browser":  info.LoginBrowser,
		"login_os":       info.LoginOS,
		"login_time":     info.LoginTime,
		"role_keys":      info.RoleKeys,
	}
}

// filterContains 按值模糊匹配保留命中的行
func filterContains(list []map[string]any, key, val string) []map[string]any {
	if val == "" {
		return list
	}
	out := make([]map[string]any, 0, len(list))
	for _, o := range list {
		if strings.Contains(fmt.Sprintf("%v", o[key]), val) {
			out = append(out, o)
		}
	}
	return out
}

// compareUserOnline 多字段排序比较
func compareUserOnline(a, b map[string]any, sortStr string) bool {
	items := strings.SplitSeq(sortStr, ",")
	for item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		field, dir, _ := strings.Cut(item, ":")
		field = strings.TrimSpace(field)
		asc := strings.EqualFold(strings.TrimSpace(dir), "asc")

		var cmp int
		switch field {
		case "device_id":
			cmp = strings.Compare(fmt.Sprintf("%v", a["device_id"]), fmt.Sprintf("%v", b["device_id"]))
		case "user_name":
			cmp = strings.Compare(fmt.Sprintf("%v", a["user_name"]), fmt.Sprintf("%v", b["user_name"]))
		case "login_ip":
			cmp = strings.Compare(fmt.Sprintf("%v", a["login_ip"]), fmt.Sprintf("%v", b["login_ip"]))
		case "login_location":
			cmp = strings.Compare(fmt.Sprintf("%v", a["login_location"]), fmt.Sprintf("%v", b["login_location"]))
		case "login_browser":
			cmp = strings.Compare(fmt.Sprintf("%v", a["login_browser"]), fmt.Sprintf("%v", b["login_browser"]))
		case "login_os":
			cmp = strings.Compare(fmt.Sprintf("%v", a["login_os"]), fmt.Sprintf("%v", b["login_os"]))
		case "login_time":
			if parse.Number(a["login_time"]) < parse.Number(b["login_time"]) {
				cmp = -1
			} else if parse.Number(a["login_time"]) > parse.Number(b["login_time"]) {
				cmp = 1
			}
		}

		if cmp != 0 {
			if asc {
				return cmp < 0
			}
			return cmp > 0
		}
	}
	return false
}
