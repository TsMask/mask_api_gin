package ctx

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constRoleDataScope "mask_api_gin/src/framework/constants/role_data_scope"
	constSystem "mask_api_gin/src/framework/constants/system"
	constToken "mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/utils/ip2region"
	"mask_api_gin/src/framework/utils/ua"
	"mask_api_gin/src/framework/vo"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// QueryMap 查询参数转换Map
func QueryMap(c *gin.Context) map[string]any {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]any, len(queryValues))
	for key, values := range queryValues {
		queryParams[key] = values[0]
	}
	return queryParams
}

// BodyJSONMap JSON参数转换Map
func BodyJSONMap(c *gin.Context) map[string]any {
	params := make(map[string]any)
	c.ShouldBindBodyWith(&params, binding.JSON)
	return params
}

// RequestParamsMap 请求参数转换Map
func RequestParamsMap(c *gin.Context) map[string]any {
	params := make(map[string]any)
	// json
	if strings.HasPrefix(c.ContentType(), "application/json") {
		c.ShouldBindBodyWith(&params, binding.JSON)
	}

	// 表单
	formParams := c.Request.PostForm
	for key, value := range formParams {
		if _, ok := params[key]; !ok {
			params[key] = value[0]
		}
	}

	// 查询
	queryParams := c.Request.URL.Query()
	for key, value := range queryParams {
		if _, ok := params[key]; !ok {
			params[key] = value[0]
		}
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
	authHeader := c.GetHeader(constToken.HEADER_KEY)
	if authHeader == "" {
		return ""
	}
	// 拆分 Authorization 请求头，提取 JWT 令牌部分
	tokenStr := strings.TrimPrefix(authHeader, constToken.HEADER_PREFIX)
	if len(tokenStr) > 99 {
		return strings.TrimSpace(tokenStr) // 去除可能存在的空格
	}
	return ""
}

// UaOsBrowser 解析请求用户代理信息
func UaOsBrowser(c *gin.Context) (string, string) {
	userAgent := c.GetHeader("user-agent")
	uaInfo := ua.Info(userAgent)

	browser := "未知"
	if bName, bVersion := uaInfo.Browser(); bName != "" {
		browser = bName
		if bVersion != "" {
			browser = bName + " " + bVersion
		}
	}

	os := "未知"
	if bos := uaInfo.OS(); bos != "" {
		os = bos
	}
	return os, browser
}

// LoginUser 登录用户信息
func LoginUser(c *gin.Context) (vo.LoginUser, error) {
	value, exists := c.Get(constSystem.CTX_LOGIN_USER)
	if exists {
		return value.(vo.LoginUser), nil
	}
	return vo.LoginUser{}, fmt.Errorf("无效登录用户信息")
}

// LoginUserToUserID 登录用户信息-用户ID
func LoginUserToUserID(c *gin.Context) string {
	value, exists := c.Get(constSystem.CTX_LOGIN_USER)
	if exists {
		loginUser := value.(vo.LoginUser)
		return loginUser.UserID
	}
	return ""
}

// LoginUserToUserName 登录用户信息-用户名称
func LoginUserToUserName(c *gin.Context) string {
	value, exists := c.Get(constSystem.CTX_LOGIN_USER)
	if exists {
		loginUser := value.(vo.LoginUser)
		return loginUser.User.UserName
	}
	return ""
}

// LoginUserByContainRoles 登录用户信息-包含角色KEY
func LoginUserByContainRoles(c *gin.Context, target string) bool {
	value, exists := c.Get(constSystem.CTX_LOGIN_USER)
	if !exists {
		return false
	}
	loginUser := value.(vo.LoginUser)
	if config.IsSysAdmin(loginUser.UserID) {
		return true
	}
	roles := loginUser.User.Roles
	for _, item := range roles {
		if item.RoleKey == target {
			return true
		}
	}
	return false
}

// LoginUserByContainPerms 登录用户信息-包含权限标识
func LoginUserByContainPerms(c *gin.Context, target string) bool {
	value, exists := c.Get(constSystem.CTX_LOGIN_USER)
	if !exists {
		return false
	}
	loginUser := value.(vo.LoginUser)
	if config.IsSysAdmin(loginUser.UserID) {
		return true
	}
	perms := loginUser.Permissions
	for _, str := range perms {
		if str == target {
			return true
		}
	}
	return false
}

// LoginUserToDataScopeSQL 登录用户信息-角色数据范围过滤SQL字符串
func LoginUserToDataScopeSQL(c *gin.Context, deptAlias string, userAlias string) string {
	dataScopeSQL := ""
	// 登录用户信息
	loginUser, err := LoginUser(c)
	if err != nil {
		return dataScopeSQL
	}
	userInfo := loginUser.User

	// 如果是系统管理员，则不过滤数据
	if config.IsSysAdmin(userInfo.UserId) {
		return dataScopeSQL
	}
	// 无用户角色
	if len(userInfo.Roles) <= 0 {
		return dataScopeSQL
	}

	// 记录角色权限范围定义添加过, 非自定数据权限不需要重复拼接SQL
	var scopeKeys []string
	var conditions []string
	for _, role := range userInfo.Roles {
		dataScope := role.DataScope

		if constRoleDataScope.ALL == dataScope {
			break
		}

		if constRoleDataScope.CUSTOM != dataScope {
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

		if constRoleDataScope.CUSTOM == dataScope {
			sql := fmt.Sprintf(`%s.dept_id IN ( SELECT dept_id FROM sys_role_dept WHERE role_id = '%s' )`, deptAlias, role.RoleId)
			conditions = append(conditions, sql)
		}

		if constRoleDataScope.DEPT == dataScope {
			sql := fmt.Sprintf("%s.dept_id = %s", deptAlias, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constRoleDataScope.DEPT_CHILD == dataScope {
			sql := fmt.Sprintf(`%s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = '%s' or find_in_set('%s' , ancestors ) )`, deptAlias, userInfo.DeptId, userInfo.DeptId)
			conditions = append(conditions, sql)
		}

		if constRoleDataScope.SELF == dataScope {
			// 数据权限为仅本人且没有userAlias别名不查询任何数据
			if userAlias == "" {
				sql := fmt.Sprintf(`%s.dept_id = '0'`, deptAlias)
				conditions = append(conditions, sql)
			} else {
				sql := fmt.Sprintf(`%s.user_id = '%s'`, userAlias, userInfo.UserId)
				conditions = append(conditions, sql)
			}
		}

		// 记录角色范围
		scopeKeys = append(scopeKeys, dataScope)
	}

	// 构建查询条件语句
	if len(conditions) > 0 {
		dataScopeSQL = fmt.Sprintf(" AND ( %s ) ", strings.Join(conditions, " OR "))
	}
	return dataScopeSQL
}
