package controller

import (
	"encoding/json"
	"fmt"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 SysOperLogController 结构体
var NewSysUserOnline = &SysUserOnlineController{
	sysUserOnlineService: service.NewSysUserOnlineImpl,
}

// 在线用户监控
//
// PATH /monitor/online
type SysUserOnlineController struct {
	// 在线用户服务
	sysUserOnlineService service.ISysUserOnline
}

// 在线用户列表
//
// GET /list
func (s *SysUserOnlineController) List(c *gin.Context) {
	ipaddr := c.Param("ipaddr")
	userName := c.Param("userName")

	// 获取所有在线用户key
	keys, _ := redis.GetKeys("", cachekey.LOGIN_TOKEN_KEY+"*")

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
	userOnlines := make([]model.SysUserOnline, 0)
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
			userOnlines = append(userOnlines, onlineUser)
		}
	}

	// 根据查询条件过滤
	filteredUserOnlines := make([]model.SysUserOnline, 0)
	if ipaddr != "" && userName != "" {
		for _, o := range userOnlines {
			if strings.Contains(o.IPAddr, ipaddr) && strings.Contains(o.UserName, userName) {
				filteredUserOnlines = append(filteredUserOnlines, o)
			}
		}
	} else if ipaddr != "" {
		for _, o := range userOnlines {
			if strings.Contains(o.IPAddr, ipaddr) {
				filteredUserOnlines = append(filteredUserOnlines, o)
			}
		}
	} else if userName != "" {
		for _, o := range userOnlines {
			if strings.Contains(o.UserName, userName) {
				filteredUserOnlines = append(filteredUserOnlines, o)
			}
		}
	} else {
		filteredUserOnlines = userOnlines
	}

	// 按登录时间排序
	sort.Slice(filteredUserOnlines, func(i, j int) bool {
		fmt.Println(i, j)
		fmt.Println(filteredUserOnlines[j].LoginTime, filteredUserOnlines[i].LoginTime)
		return filteredUserOnlines[j].LoginTime > filteredUserOnlines[i].LoginTime
	})

	c.JSON(200, result.Ok(map[string]any{
		"total": len(userOnlines),
		"rows":  filteredUserOnlines,
	}))
}

// 在线用户强制退出
//
// DELETE /:tokenId
func (s *SysUserOnlineController) ForceLogout(c *gin.Context) {
	tokenId := c.Param("tokenId")
	if tokenId == "" || tokenId == "*" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 删除token
	ok, _ := redis.Del("", cachekey.LOGIN_TOKEN_KEY+tokenId)
	if ok {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}
