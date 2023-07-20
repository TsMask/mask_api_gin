package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants/admin"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 用户信息
//
// PATH /system/user
var SysUser = &sysUserController{
	sysUserService: service.SysUserImpl,
	sysRoleService: service.SysRoleImpl,
	sysPostService: service.SysPostImpl,
}

type sysUserController struct {
	// 用户服务
	sysUserService service.ISysUser
	// 角色服务
	sysRoleService service.ISysRole
	// 岗位服务
	sysPostService service.ISysPost
}

// 用户信息列表
//
// GET /list
func (s *sysUserController) List(c *gin.Context) {
	// 查询参数转换map
	querys := ctxUtils.QueryMapString(c)
	dataScopeSQL := repoUtils.DataScopeSQL("d", "u")
	list := s.sysUserService.SelectUserPage(querys, dataScopeSQL)
	c.JSON(200, result.Ok(list))
}

// 用户信息详情
//
// GET /:userId
func (s *sysUserController) Info(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 查询系统角色列表
	dataScopeSQL := ctxUtils.LoginUserToDataScopeSQL(c, "d", "")
	roles := s.sysRoleService.SelectRoleList(model.SysRole{}, dataScopeSQL)

	// 不是系统指定管理员需要排除其角色
	if !config.IsAdmin(userId) {
		rolesFilter := make([]model.SysRole, 0)
		for _, r := range roles {
			if r.RoleID != admin.ROLE_ID {
				rolesFilter = append(rolesFilter, r)
			}
		}
		roles = rolesFilter
	}

	// 查询系统岗位列表
	posts := s.sysPostService.SelectPostList(model.SysPost{})

	// 新增用户时，用户ID为0
	if userId == "0" {
		c.JSON(200, result.OkData(map[string]interface{}{
			"user":    map[string]interface{}{},
			"roleIds": []string{},
			"postIds": []string{},
			"roles":   roles,
			"posts":   posts,
		}))
		return
	}

	// 检查用户是否存在
	user := s.sysUserService.SelectUserById(userId)
	if user.UserID != userId {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 角色ID组
	roleIds := make([]string, 0)
	for _, r := range user.Roles {
		roleIds = append(roleIds, r.RoleID)
	}

	// 岗位ID组
	postIds := make([]string, 0)
	userPosts := s.sysPostService.SelectPostListByUserId(userId)
	for _, p := range userPosts {
		postIds = append(postIds, p.PostID)
	}

	c.JSON(200, result.OkData(map[string]interface{}{
		"user":    user,
		"roleIds": roleIds,
		"postIds": postIds,
		"roles":   roles,
		"posts":   posts,
	}))
}

// 用户信息新增
//
// POST /
func (s *sysUserController) Add(c *gin.Context) {
	userIds := c.Param("userIds")
	if userIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(userIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows, err := s.sysUserService.DeleteUserByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 用户信息修改
//
// POST /
func (s *sysUserController) Edit(c *gin.Context) {
	userIds := c.Param("userIds")
	if userIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(userIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows, err := s.sysUserService.DeleteUserByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 用户信息删除
//
// DELETE /:userIds
func (s *sysUserController) Remove(c *gin.Context) {
	userIds := c.Param("userIds")
	if userIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(userIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows, err := s.sysUserService.DeleteUserByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 用户重置密码
//
// PUT /resetPwd
func (s *sysUserController) ResetPwd(c *gin.Context) {
	var body struct {
		UserID   string `json:"userId"  binding:"required"`
		Password string `json:"password"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否管理员用户
	if config.IsAdmin(body.UserID) {
		c.JSON(200, result.ErrMsg("不允许操作管理员用户"))
		return
	}

	user := s.sysUserService.SelectUserById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}
	if !regular.ValidPassword(body.Password) {
		c.JSON(200, result.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}

	userName := ctxUtils.LoginUserToUserName(c)
	sysUser := model.SysUser{
		UserID:   body.UserID,
		Password: body.Password,
		UpdateBy: userName,
	}
	rows := s.sysUserService.UpdateUser(sysUser)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 用户状态修改
//
// PUT /changeStatus
func (s *sysUserController) Status(c *gin.Context) {
	var body struct {
		UserID string `json:"userId"  binding:"required"`
		Status string `json:"status"  binding:"required"`
	}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	user := s.sysUserService.SelectUserById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}
	// 与旧值相等不变更
	if user.Status == body.Status {
		c.JSON(200, result.ErrMsg("变更状态与旧值相等！"))
		return
	}

	userName := ctxUtils.LoginUserToUserName(c)
	sysUser := model.SysUser{
		UserID:   body.UserID,
		Status:   body.Status,
		UpdateBy: userName,
	}
	rows := s.sysUserService.UpdateUser(sysUser)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}
