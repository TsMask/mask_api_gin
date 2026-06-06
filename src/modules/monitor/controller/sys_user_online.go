package controller

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
)

// NewSysUserOnline 实例化控制层
var NewSysUserOnline = &SysUserOnlineController{
	sysUserOnlineService: service.NewSysUserOnline,
}

// SysUserOnlineController 在线用户信息 控制层处理
//
// PATH /monitor/user-online
type SysUserOnlineController struct {
	sysUserOnlineService *service.SysUserOnline
}

// List 在线用户列表
//
// GET /list
func (s SysUserOnlineController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)

	// SCAN 获取所有 key（仅名称），排序保证顺序确定性
	keys, _ := redis.GetKeys("", constants.CACHE_TOKEN_DEVICE+":*")
	sort.Strings(keys)

	total := int64(len(keys))
	if total <= 0 {
		c.JSON(200, resp.OkData(map[string]any{"total": 0, "rows": []model.SysUserOnline{}}))
		return
	}

	// 分页：只取当前页的 key，避免全量 mget
	pn := int(parse.Number(query["pageNum"]))
	if pn < 1 {
		pn = 1
	}
	ps := int(parse.Number(query["pageSize"]))
	if ps < 1 {
		ps = 10
	}
	start := (pn - 1) * ps
	end := pn * ps
	if start > len(keys) {
		start = len(keys)
	}
	if end > len(keys) {
		end = len(keys)
	}
	pageKeys := keys[start:end]

	values, _ := redis.GetBatch("", pageKeys)

	// 解析当前页数据
	list := make([]model.SysUserOnline, 0, len(values))
	for _, v := range values {
		str := fmt.Sprint(v)
		if str == "" {
			continue
		}
		var info token.UserInfo
		if err := json.Unmarshal([]byte(str), &info); err != nil {
			continue
		}
		onlineUser := s.sysUserOnlineService.UserInfoToUserOnline(info)
		if onlineUser.TokenID != "" {
			list = append(list, onlineUser)
		}
	}

	// 过滤
	if v, ok := query["tokenId"]; ok && v != "" {
		filtered := make([]model.SysUserOnline, 0, len(list))
		for _, o := range list {
			if strings.Contains(o.TokenID, v) {
				filtered = append(filtered, o)
			}
		}
		list = filtered
	}
	if v, ok := query["userName"]; ok && v != "" {
		filtered := make([]model.SysUserOnline, 0, len(list))
		for _, o := range list {
			if strings.Contains(o.UserName, v) {
				filtered = append(filtered, o)
			}
		}
		list = filtered
	}
	if v, ok := query["loginIp"]; ok && v != "" {
		filtered := make([]model.SysUserOnline, 0, len(list))
		for _, o := range list {
			if strings.Contains(o.LoginIp, v) {
				filtered = append(filtered, o)
			}
		}
		list = filtered
	}
	if v, ok := query["deptName"]; ok && v != "" {
		filtered := make([]model.SysUserOnline, 0, len(list))
		for _, o := range list {
			if strings.Contains(o.DeptName, v) {
				filtered = append(filtered, o)
			}
		}
		list = filtered
	}
	if v, ok := query["browser"]; ok && v != "" {
		filtered := make([]model.SysUserOnline, 0, len(list))
		for _, o := range list {
			if strings.Contains(o.Browser, v) {
				filtered = append(filtered, o)
			}
		}
		list = filtered
	}
	if v, ok := query["os"]; ok && v != "" {
		filtered := make([]model.SysUserOnline, 0, len(list))
		for _, o := range list {
			if strings.Contains(o.OS, v) {
				filtered = append(filtered, o)
			}
		}
		list = filtered
	}

	// 排序，格式: sort=loginTime,userName:asc（不加方向默认 DESC）
	if sortStr, ok := query["sort"]; ok && sortStr != "" {
		sort.SliceStable(list, func(i, j int) bool {
			return compareUserOnline(list[i], list[j], sortStr)
		})
	} else {
		// 默认登录时间降序
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].LoginTime > list[j].LoginTime
		})
	}

	c.JSON(200, resp.OkData(map[string]any{
		"total": total,
		"rows":  list,
	}))
}

// Logout 在线用户强制退出
//
// DELETE /logout/:tokenId
func (s SysUserOnlineController) Logout(c *gin.Context) {
	tokenIdStr := c.Param("tokenId")
	if tokenIdStr == "" || strings.Contains(tokenIdStr, "*") {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: tokenId is empty"))
		return
	}

	uniqueIDs := parse.RemoveDuplicatesToArray(tokenIdStr, ",")
	for _, v := range uniqueIDs {
		key := constants.CACHE_TOKEN_DEVICE + ":" + v
		if err := redis.Del("", key); err != nil {
			c.JSON(200, resp.ErrMsg(err.Error()))
			return
		}
	}

	c.JSON(200, resp.Ok(nil))
}

// compareUserOnline 多字段排序比较
func compareUserOnline(a, b model.SysUserOnline, sortStr string) bool {
	items := strings.Split(sortStr, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		field, dir, _ := strings.Cut(item, ":")
		field = strings.TrimSpace(field)
		asc := strings.EqualFold(strings.TrimSpace(dir), "asc")

		var cmp int
		switch field {
		case "tokenId":
			cmp = strings.Compare(a.TokenID, b.TokenID)
		case "userName":
			cmp = strings.Compare(a.UserName, b.UserName)
		case "loginIp":
			cmp = strings.Compare(a.LoginIp, b.LoginIp)
		case "deptName":
			cmp = strings.Compare(a.DeptName, b.DeptName)
		case "browser":
			cmp = strings.Compare(a.Browser, b.Browser)
		case "os":
			cmp = strings.Compare(a.OS, b.OS)
		case "loginTime":
			if a.LoginTime < b.LoginTime {
				cmp = -1
			} else if a.LoginTime > b.LoginTime {
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
