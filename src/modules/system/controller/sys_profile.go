package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constUploadSubPath "mask_api_gin/src/framework/constants/upload_sub_path"
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
	sysUserService service.ISysUserService // 用户服务
	sysRoleService service.ISysRoleService // 角色服务
	sysPostService service.ISysPostService // 岗位服务
	sysMenuService service.ISysMenuService // 菜单服务
}

// Info 个人信息
//
// GET /
func (s *SysProfileController) Info(c *gin.Context) {
	loginUser, err := ctx.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}

	// 查询用户所属角色组
	var roleGroup []string
	roles := s.sysRoleService.FindByUserId(loginUser.UserID)
	for _, role := range roles {
		roleGroup = append(roleGroup, role.RoleName)
	}
	isAdmin := config.IsAdmin(loginUser.UserID)
	if isAdmin {
		roleGroup = append(roleGroup, "管理员")
	}

	// 查询用户所属岗位组
	var postGroup []string
	posts := s.sysPostService.FindByUserId(loginUser.UserID)
	for _, post := range posts {
		postGroup = append(postGroup, post.PostName)
	}

	c.JSON(200, result.OkData(map[string]any{
		"user":      loginUser.User,
		"roleGroup": parse.RemoveDuplicates(roleGroup),
		"postGroup": parse.RemoveDuplicates(postGroup),
	}))
}

// UpdateProfile 个人信息修改
//
// PUT /
func (s *SysProfileController) UpdateProfile(c *gin.Context) {
	var body struct {
		NickName    string `json:"nickName" binding:"required"`        // 昵称
		Sex         string `json:"sex" binding:"required,oneof=0 1 2"` // 性别
		PhoneNumber string `json:"phonenumber"`                        // 手机号
		Email       string `json:"email"`                              // 邮箱
	}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
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
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.PhoneNumber, userId)
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
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, userId)
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
	}

	// 查询当前登录用户信息
	user := s.sysUserService.FindById(userId)
	if user.UserID != userId {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 用户基本资料
	user.UpdateBy = userName
	user.NickName = body.NickName
	user.Phone = body.PhoneNumber
	user.Email = body.Email
	user.Sex = body.Sex
	rows := s.sysUserService.Update(user)
	if rows > 0 {
		// 更新缓存用户信息
		loginUser.User = user
		// 刷新令牌信息
		token.Cache(&loginUser)

		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.ErrMsg("上传图片异常"))
}

// UpdatePwd 个人重置密码
//
// PUT /updatePwd
func (s *SysProfileController) UpdatePwd(c *gin.Context) {
	var body struct {
		OldPassword string `json:"oldPassword" binding:"required"` // 旧密码
		NewPassword string `json:"newPassword" binding:"required"` // 新密码
	}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
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
	user := s.sysUserService.FindById(userId)
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
	rows := s.sysUserService.Update(user)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Avatar 个人头像上传
//
// POST /avatar
func (s *SysProfileController) Avatar(c *gin.Context) {
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 上传文件转存
	filePath, err := file.TransferUploadFile(formFile, constUploadSubPath.Avatar, []string{".jpg", ".jpeg", ".png"})
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
	user := s.sysUserService.FindById(userId)
	if user.UserID != userId {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 更新头像地址
	user.Avatar = filePath
	user.UpdateBy = userName
	rows := s.sysUserService.Update(user)
	if rows > 0 {
		// 更新缓存用户信息
		loginUser.User = user
		// 刷新令牌信息
		token.Cache(&loginUser)

		c.JSON(200, result.OkData(filePath))
		return
	}
	c.JSON(200, result.Err(nil))
}
