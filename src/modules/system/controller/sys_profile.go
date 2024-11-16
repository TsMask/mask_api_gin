package controller

import (
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/utils/token"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// NewSysProfile 实例化控制层
var NewSysProfile = &SysProfileController{
	sysUserService: service.NewSysUser,
	sysRoleService: service.NewSysRole,
	sysPostService: service.NewSysPost,
	sysMenuService: service.NewSysMenu,
}

// SysProfileController 个人信息
//
// PATH /system/user/profile
type SysProfileController struct {
	sysUserService *service.SysUser // 用户服务
	sysRoleService *service.SysRole // 角色服务
	sysPostService *service.SysPost // 岗位服务
	sysMenuService *service.SysMenu // 菜单服务
}

// Info 个人信息
//
// GET /
func (s SysProfileController) Info(c *gin.Context) {
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, response.CodeMsg(40003, err.Error()))
		return
	}

	// 查询用户所属角色组
	var roleGroup []string
	roles := s.sysRoleService.FindByUserId(loginUser.UserId)
	for _, role := range roles {
		roleGroup = append(roleGroup, role.RoleName)
	}

	// 查询用户所属岗位组
	var postGroup []string
	posts := s.sysPostService.FindByUserId(loginUser.UserId)
	for _, post := range posts {
		postGroup = append(postGroup, post.PostName)
	}

	c.JSON(200, response.OkData(map[string]any{
		"user":      loginUser.User,
		"roleGroup": parse.RemoveDuplicates(roleGroup),
		"postGroup": parse.RemoveDuplicates(postGroup),
	}))
}

// UpdateProfile 个人信息修改
//
// PUT /
func (s SysProfileController) UpdateProfile(c *gin.Context) {
	var body struct {
		NickName string `json:"nickName" binding:"required"`        // 昵称
		Sex      string `json:"sex" binding:"required,oneof=0 1 2"` // 性别
		Phone    string `json:"phone"`                              // 手机号
		Email    string `json:"email"`                              // 邮箱
		Avatar   string `json:"avatar"`                             // 头像地址
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 登录用户信息
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, response.CodeMsg(40003, err.Error()))
		return
	}
	userId := loginUser.UserId

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, userId)
			if !uniquePhone {
				c.JSON(200, response.CodeMsg(40018, "抱歉，手机号码已存在"))
				return
			}
		} else {
			c.JSON(200, response.CodeMsg(40018, "抱歉，手机号码格式错误"))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, userId)
			if !uniqueEmail {
				c.JSON(200, response.CodeMsg(40019, "抱歉，邮箱已存在"))
				return
			}
		} else {
			c.JSON(200, response.CodeMsg(40019, "抱歉，邮箱格式错误"))
			return
		}
	}

	// 查询当前登录用户信息
	userInfo := s.sysUserService.FindById(userId)
	if userInfo.UserId != userId {
		c.JSON(200, response.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 用户基本资料
	userInfo.NickName = body.NickName
	userInfo.Phone = body.Phone
	userInfo.Email = body.Email
	userInfo.Sex = body.Sex
	userInfo.Avatar = body.Avatar
	userInfo.Passwd = "" // 密码不更新
	userInfo.UpdateBy = userInfo.UserName
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		// 更新缓存用户信息
		loginUser.User = userInfo
		// 刷新令牌信息
		token.Cache(&loginUser)

		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.ErrMsg("上传图片异常"))
}

// UpdatePasswd 个人重置密码
//
// PUT /passwd
func (s SysProfileController) UpdatePasswd(c *gin.Context) {
	var body struct {
		OldPassword string `json:"oldPassword" binding:"required"` // 旧密码
		NewPassword string `json:"newPassword" binding:"required"` // 新密码
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 登录用户信息
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, response.CodeMsg(40003, err.Error()))
		return
	}
	userId := loginUser.UserId

	// 查询当前登录用户信息得到密码值
	userInfo := s.sysUserService.FindById(userId)
	if userInfo.UserId != userId {
		c.JSON(200, response.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 检查匹配用户密码
	oldCompare := crypto.BcryptCompare(body.OldPassword, userInfo.Passwd)
	if !oldCompare {
		c.JSON(200, response.ErrMsg("修改密码失败，旧密码错误"))
		return
	}
	newCompare := crypto.BcryptCompare(body.NewPassword, userInfo.Passwd)
	if newCompare {
		c.JSON(200, response.ErrMsg("新密码不能与旧密码相同"))
		return
	}

	// 修改新密码
	userInfo.Passwd = body.NewPassword
	userInfo.UpdateBy = userInfo.UserName
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}
