package controller

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"

	"encoding/json"
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
// PATH /monitor/user-online
type SysUserOnlineController struct {
	sysUserOnlineService *service.SysUserOnline // 在线用户服务
}

// List 在线用户列表
//
// GET /list
func (s SysUserOnlineController) List(c *gin.Context) {
	loginIp := c.Query("loginIp")
	userName := c.Query("userName")

	// 获取所有在线用户key
	keys, _ := redis.GetKeys("", constants.CACHE_LOGIN_TOKEN+"*")

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
	if loginIp != "" && userName != "" {
		for _, o := range userOnline {
			if strings.Contains(o.LoginIp, loginIp) && strings.Contains(o.UserName, userName) {
				filteredUserOnline = append(filteredUserOnline, o)
			}
		}
	} else if loginIp != "" {
		for _, o := range userOnline {
			if strings.Contains(o.LoginIp, loginIp) {
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

	c.JSON(200, response.OkData(map[string]any{
		"total": len(filteredUserOnline),
		"rows":  filteredUserOnline,
	}))
}

// Logout 在线用户强制退出
//
// DELETE /logout/:tokenId
func (s SysUserOnlineController) Logout(c *gin.Context) {
	tokenIdStr := c.Param("tokenId")
	if tokenIdStr == "" || strings.Contains(tokenIdStr, "*") {
		c.JSON(400, response.CodeMsg(40010, "bind err: tokenId is empty"))
		return
	}

	// 处理字符转id数组后去重
	ids := strings.Split(tokenIdStr, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	for _, v := range uniqueIDs {
		key := constants.CACHE_LOGIN_TOKEN + v
		if err := redis.Del("", key); err != nil {
			c.JSON(200, response.ErrMsg(err.Error()))
			return
		}
	}

	c.JSON(200, response.Ok(nil))
}
