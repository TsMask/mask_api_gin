package controller

import (
	"encoding/json"
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewSysUserOnline 实例化控制层 SysUserOnlineController 结构体
var NewSysUserOnline = &SysUserOnlineController{
	sysUserOnlineService: service.NewSysUserOnline,
}

// SysUserOnlineController 在线用户信息 控制层处理
//
// PATH /monitor/online
type SysUserOnlineController struct {
	sysUserOnlineService *service.SysUserOnline // 在线用户服务
}

// List 在线用户列表
//
// GET /list
func (s SysUserOnlineController) List(c *gin.Context) {
	ipaddr := c.Query("ipaddr")
	userName := c.Query("userName")

	// 获取所有在线用户key
	keys, _ := redis.GetKeys("", constCacheKey.LOGIN_TOKEN_KEY+"*")

	// 分批获取
	arr := make([]string, 0)
	for i := 0; i < len(keys); i += 20 {
		end := i + 20
		if end > len(keys) {
			end = len(keys)
		}
		chunk := keys[i:end]
		values, _ := redis.GetBatch("", chunk)
		for _, v := range values {
			arr = append(arr, v.(string))
		}
	}

	// 遍历字符串信息解析组合可用对象
	var userOnline []model.SysUserOnline
	for _, str := range arr {
		if str == "" {
			continue
		}

		var loginUser vo.LoginUser
		err := json.Unmarshal([]byte(str), &loginUser)
		if err != nil {
			continue
		}

		onlineUser := s.sysUserOnlineService.LoginUserToUserOnline(loginUser)
		if onlineUser.TokenID != "" {
			userOnline = append(userOnline, onlineUser)
		}
	}

	// 根据查询条件过滤
	filteredUserOnline := make([]model.SysUserOnline, 0)
	if ipaddr != "" && userName != "" {
		for _, o := range userOnline {
			if strings.Contains(o.IPAddr, ipaddr) && strings.Contains(o.UserName, userName) {
				filteredUserOnline = append(filteredUserOnline, o)
			}
		}
	} else if ipaddr != "" {
		for _, o := range userOnline {
			if strings.Contains(o.IPAddr, ipaddr) {
				filteredUserOnline = append(filteredUserOnline, o)
			}
		}
	} else if userName != "" {
		for _, o := range userOnline {
			if strings.Contains(o.UserName, userName) {
				filteredUserOnline = append(filteredUserOnline, o)
			}
		}
	} else {
		filteredUserOnline = userOnline
	}

	// 按登录时间排序
	sort.Slice(filteredUserOnline, func(i, j int) bool {
		return filteredUserOnline[j].LoginTime > filteredUserOnline[i].LoginTime
	})

	c.JSON(200, result.Ok(map[string]any{
		"total": len(filteredUserOnline),
		"rows":  filteredUserOnline,
	}))
}

// Logout 在线用户强制退出
//
// DELETE /?tokenId=xxxxx
func (s SysUserOnlineController) Logout(c *gin.Context) {
	tokenId, ok := c.GetQuery("tokenId")
	if !ok || tokenId == "" || strings.Contains(tokenId, "*") {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组后去重
	ids := strings.Split(tokenId, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 批量删除token
	for _, v := range uniqueIDs {
		key := constCacheKey.LOGIN_TOKEN_KEY + v
		if err := redis.Del("", key); err != nil {
			c.JSON(200, result.ErrMsg(err.Error()))
			return
		}
	}

	c.JSON(200, result.Ok(nil))
}
