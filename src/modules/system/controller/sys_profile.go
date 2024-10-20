package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants/admin"
	"mask_api_gin/src/framework/constants/uploadsubpath"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/utils/token"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 实例化控制层 SysProfileController 结构体
var NewSysProfile = &SysProfileController{
	sysUserService: service.NewSysUserImpl,
	sysRoleService: service.NewSysRoleImpl,
	sysPostService: service.NewSysPostImpl,
	sysMenuService: service.NewSysMenuImpl,
}

// 个人信息
//
// PATH /system/user/profile
type SysProfileController struct {
	// 用户服务
	sysUserService service.ISysUser
	// 角色服务
	sysRoleService service.ISysRole
	// 岗位服务
	sysPostService service.ISysPost
	// 菜单服务
	sysMenuService service.ISysMenu
}

// 个人信息
//
// GET /
func (s *SysProfileController) Info(c *gin.Context) {
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}

	// 查询用户所属角色组
	roleGroup := []string{}
	roles := s.sysRoleService.SelectRoleListByUserId(loginUser.UserID)
	for _, role := range roles {
		roleGroup = append(roleGroup, role.RoleName)
	}
	isAdmin := config.IsAdmin(loginUser.UserID)
	if isAdmin {
		roleGroup = append(roleGroup, "管理员")
	}

	// 查询用户所属岗位组
	postGroup := []string{}
	posts := s.sysPostService.SelectPostListByUserId(loginUser.UserID)
	for _, post := range posts {
		postGroup = append(postGroup, post.PostName)
	}

	c.JSON(200, result.OkData(map[string]any{
		"user":      loginUser.User,
		"roleGroup": parse.RemoveDuplicates(roleGroup),
		"postGroup": parse.RemoveDuplicates(postGroup),
	}))
}

// 个人信息修改
//
// PUT /
func (s *SysProfileController) UpdateProfile(c *gin.Context) {
	var body struct {
		// 昵称
		NickName string `json:"nickName" binding:"required"`
		// 性别
		Sex string `json:"sex" binding:"required"`
		// 手机号
		PhoneNumber string `json:"phonenumber"`
		// 邮箱
		Email string `json:"email"`
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.Sex == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 登录用户信息
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}
	userId := loginUser.UserID
	userName := loginUser.User.UserName

	// 检查手机号码格式并判断是否唯一
	if body.PhoneNumber != "" {
		if regular.ValidMobile(body.PhoneNumber) {
			uniquePhone := s.sysUserService.CheckUniquePhone(body.PhoneNumber, userId)
			if !uniquePhone {
				msg := fmt.Sprintf("修改用户【%s】失败，手机号码已存在", userName)
				c.JSON(200, result.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，手机号码格式错误", userName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	} else {
		body.PhoneNumber = "nil"
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueEmail(body.Email, userId)
			if !uniqueEmail {
				msg := fmt.Sprintf("修改用户【%s】失败，邮箱已存在", userName)
				c.JSON(200, result.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，邮箱格式错误", userName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	} else {
		body.Email = "nil"
	}

	// 查询当前登录用户信息
	user := s.sysUserService.SelectUserById(userId)
	if user.UserID != userId {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 用户基本资料
	user.UpdateBy = userName
	user.NickName = body.NickName
	user.PhoneNumber = body.PhoneNumber
	user.Email = body.Email
	user.Sex = body.Sex
	rows := s.sysUserService.UpdateUser(user)
	if rows > 0 {
		// 更新缓存用户信息
		loginUser.User = s.sysUserService.SelectUserByUserName(userName)
		// 用户权限组标识
		isAdmin := config.IsAdmin(userId)
		if isAdmin {
			loginUser.Permissions = []string{admin.PERMISSION}
		} else {
			perms := s.sysMenuService.SelectMenuPermsByUserId(userId)
			loginUser.Permissions = parse.RemoveDuplicates(perms)
		}
		// 刷新令牌信息
		token.Cache(&loginUser)

		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.ErrMsg("上传图片异常"))
}

// 个人重置密码
//
// PUT /updatePwd
func (s *SysProfileController) UpdatePwd(c *gin.Context) {
	var body struct {
		// 旧密码
		OldPassword string `json:"oldPassword" binding:"required"`
		// 新密码
		NewPassword string `json:"newPassword" binding:"required"`
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 登录用户信息
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}
	userId := loginUser.UserID
	userName := loginUser.User.UserName

	// 查询当前登录用户信息得到密码值
	user := s.sysUserService.SelectUserById(userId)
	if user.UserID != userId {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 检查匹配用户密码
	oldCompare := crypto.BcryptCompare(body.OldPassword, user.Password)
	if !oldCompare {
		c.JSON(200, result.ErrMsg("修改密码失败，旧密码错误"))
		return
	}
	newCompare := crypto.BcryptCompare(body.NewPassword, user.Password)
	if newCompare {
		c.JSON(200, result.ErrMsg("新密码不能与旧密码相同"))
		return
	}

	// 修改新密码
	user.UpdateBy = userName
	user.Password = body.NewPassword
	rows := s.sysUserService.UpdateUser(user)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 个人头像上传
//
// POST /avatar
func (s *SysProfileController) Avatar(c *gin.Context) {
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 上传文件转存
	filePath, err := file.TransferUploadFile(formFile, uploadsubpath.AVATART, []string{".jpg", ".jpeg", ".png"})
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 登录用户信息
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}
	userId := loginUser.UserID
	userName := loginUser.User.UserName

	// 查询当前登录用户信息
	user := s.sysUserService.SelectUserById(userId)
	if user.UserID != userId {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 更新头像地址
	user.Avatar = filePath
	user.UpdateBy = userName
	rows := s.sysUserService.UpdateUser(user)
	if rows > 0 {
		// 更新缓存用户信息
		loginUser.User = s.sysUserService.SelectUserByUserName(userName)
		// 用户权限组标识
		isAdmin := config.IsAdmin(userId)
		if isAdmin {
			loginUser.Permissions = []string{admin.PERMISSION}
		} else {
			perms := s.sysMenuService.SelectMenuPermsByUserId(userId)
			loginUser.Permissions = parse.RemoveDuplicates(perms)
		}
		// 刷新令牌信息
		token.Cache(&loginUser)

		c.JSON(200, result.OkData(filePath))
		return
	}
	c.JSON(200, result.Err(nil))
}
