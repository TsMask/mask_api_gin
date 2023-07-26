package ctx

import (
	"errors"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/constants/roledatascope"
	"mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/utils/ip2region"
	"mask_api_gin/src/framework/utils/ua"
	"mask_api_gin/src/framework/vo"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// QueryMapString 查询参数转换MapString
func QueryMapString(c *gin.Context) map[string]string {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]string)
	for key, values := range queryValues {
		queryParams[key] = values[0]
	}
	return queryParams
}

// BodyJSONMapString JSON参数转换MapString
func BodyJSONMapString(c *gin.Context) map[string]string {
	params := make(map[string]string)
	c.ShouldBindBodyWith(&params, binding.JSON)
	return params
}

// RequestParamsMap 请求参数转换Map
func RequestParamsMap(c *gin.Context) map[string]any {
	params := make(map[string]interface{})
	// json
	if strings.HasPrefix(c.ContentType(), "application/json") {
		c.ShouldBindBodyWith(&params, binding.JSON)
	}

	// 表单
	bodyParams := c.Request.PostForm
	for key, value := range bodyParams {
		params[key] = value[0]
	}

	// 查询
	queryParams := c.Request.URL.Query()
	for key, value := range queryParams {
		params[key] = value[0]
	}
	return params
}

// IPAddrLocation 解析ip地址
func IPAddrLocation(c *gin.Context) (string, string) {
	ip := ip2region.ClientIP(c.ClientIP())
	location := ip2region.RealAddressByIp(ip)
	return ip, location
}

// Authorization 解析请求头
func Authorization(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	// 拆分 Authorization 请求头，提取 JWT 令牌部分
	arr := strings.Split(authHeader, token.HEADER_PREFIX)
	if len(arr) == 2 && arr[1] == "" {
		return ""
	}
	return arr[1]
}

// UaOsBrowser 解析请求用户代理信息
func UaOsBrowser(c *gin.Context) (string, string) {
	userAgent := c.GetHeader("user-agent")
	uaInfo := ua.Info(userAgent)

	browser := "未知 未知"
	bName, bVersion := uaInfo.Browser()
	if bName != "" && bVersion != "" {
		browser = bName + " " + bVersion
	}

	os := "未知 未知"
	bos := uaInfo.OS()
	if bos != "" {
		os = bos
	}
	return os, browser
}

// LoginUser 登录用户信息
func LoginUser(c *gin.Context) (vo.LoginUser, error) {
	value, exists := c.Get(common.CTX_LOGIN_USER)
	if exists {
		return value.(vo.LoginUser), nil
	}
	return vo.LoginUser{}, errors.New("无效登录用户信息")
}

// LoginUserToUserID 登录用户信息-用户ID
func LoginUserToUserID(c *gin.Context) string {
	value, exists := c.Get(common.CTX_LOGIN_USER)
	if exists {
		loginUser := value.(vo.LoginUser)
		return loginUser.UserID
	}
	return ""
}

// LoginUserToUserName 登录用户信息-用户名称
func LoginUserToUserName(c *gin.Context) string {
	value, exists := c.Get(common.CTX_LOGIN_USER)
	if exists {
		loginUser := value.(vo.LoginUser)
		return loginUser.User.UserName
	}
	return ""
}

// LoginUserToDataScopeSQL 登录用户信息-角色数据范围过滤SQL字符串
func LoginUserToDataScopeSQL(c *gin.Context, deptAlias string, userAlias string) string {
	loginUser, err := LoginUser(c)
	if err != nil {
		return ""
	}

	// 登录用户信息
	userInfo := loginUser.User

	// 如果是管理员，则不过滤数据
	if config.IsAdmin(userInfo.UserID) {
		return ""
	}
	// 无用户角色
	if len(userInfo.Roles) <= 0 {
		return ""
	}

	// 记录角色权限范围定义添加过, 非自定数据权限不需要重复拼接SQL
	var scopeKeys []string
	var conditions []string
	for _, role := range userInfo.Roles {
		dataScope := role.DataScope

		if roledatascope.ALL == dataScope {
			break
		}

		if roledatascope.CUSTOM != dataScope {
			hasKey := false
			for _, key := range scopeKeys {
				if key == dataScope {
					hasKey = true
					break
				}
			}
			if hasKey {
				continue
			}
		}

		if roledatascope.CUSTOM == dataScope {
			sql := fmt.Sprintf(`%s.dept_id IN ( SELECT dept_id FROM sys_role_dept WHERE role_id = %s )`, deptAlias, role.RoleID)
			conditions = append(conditions, sql)
		}

		if roledatascope.DEPT_AND_CHILD == dataScope {
			sql := fmt.Sprintf(`%s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = %s or find_in_set(%s , ancestors ) )`, deptAlias, userInfo.DeptID, userInfo.DeptID)
			conditions = append(conditions, sql)
		}

		if roledatascope.SELF == dataScope {
			// 数据权限为仅本人且没有userAlias别名不查询任何数据
			if userAlias == "" {
				sql := fmt.Sprintf(`%s.dept_id = 0`, deptAlias)
				conditions = append(conditions, sql)
			} else {
				sql := fmt.Sprintf(`%s.user_id = %s`, userAlias, userInfo.UserID)
				conditions = append(conditions, sql)
			}
		}

		// 记录角色范围
		scopeKeys = append(scopeKeys, dataScope)
	}

	// 构建查询条件语句
	dataScopeSQL := ""
	if len(conditions) > 0 {
		dataScopeSQL = fmt.Sprintf(" and ( %s ) ", strings.Join(conditions, " or "))
	}
	fmt.Println("dataScopeSQL======> ", dataScopeSQL)
	return dataScopeSQL
}
