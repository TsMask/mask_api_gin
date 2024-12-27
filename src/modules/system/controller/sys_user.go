package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
)

// NewSysUser 实例化控制层
var NewSysUser = &SysUserController{
	sysUserService: service.NewSysUser,
	sysRoleService: service.NewSysRole,
	sysPostService: service.NewSysPost,
}

// SysUserController 用户信息
//
// PATH /system/user
type SysUserController struct {
	sysUserService *service.SysUser // 用户服务
	sysRoleService *service.SysRole // 角色服务
	sysPostService *service.SysPost // 岗位服务
}

// List 用户信息列表
//
// GET /list
func (s SysUserController) List(c *gin.Context) {
	queryMap := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysUserService.FindByPage(queryMap, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 用户信息详情
//
// GET /:userId
func (s SysUserController) Info(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, resp.CodeMsg(40010, "bind err: userId is empty"))
		return
	}

	// 查询系统角色列表
	roles := s.sysRoleService.Find(model.SysRole{})
	// 查询系统岗位列表
	posts := s.sysPostService.Find(model.SysPost{})

	// 新增用户时，用户ID为0
	if userId == "0" {
		c.JSON(200, resp.OkData(map[string]any{
			"user":    map[string]any{},
			"roleIds": []string{},
			"postIds": []string{},
			"roles":   roles,
			"posts":   posts,
		}))
		return
	}

	// 检查用户是否存在
	userInfo := s.sysUserService.FindById(userId)
	if userInfo.UserId != userId {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 角色ID组
	roleIds := make([]string, 0)
	for _, r := range userInfo.Roles {
		roleIds = append(roleIds, r.RoleId)
	}

	// 岗位ID组
	postIds := make([]string, 0)
	userPosts := s.sysPostService.FindByUserId(userId)
	for _, p := range userPosts {
		postIds = append(postIds, p.PostId)
	}

	c.JSON(200, resp.OkData(map[string]any{
		"user":    userInfo,
		"roleIds": roleIds,
		"postIds": postIds,
		"roles":   roles,
		"posts":   posts,
	}))
}

// Add 用户信息新增
//
// POST /
func (s SysUserController) Add(c *gin.Context) {
	var body model.SysUser
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, resp.ErrMsg(resp.FormatBindError(err)))
		return
	}
	if body.UserId != "" {
		c.JSON(400, resp.CodeMsg(40010, "bind err: userId not is empty"))
		return
	}

	// 密码单独取，避免序列化输出
	var bodyPassword struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&bodyPassword); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(400, resp.CodeMsg(40010, errMsgs))
		return
	}
	body.Password = bodyPassword.Password

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(body.UserName, "")
	if !uniqueUserName {
		msg := fmt.Sprintf("新增用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, "")
			if !uniquePhone {
				msg := fmt.Sprintf("新增用户【%s】失败，手机号码已存在", body.UserName)
				c.JSON(200, resp.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("新增用户【%s】失败，手机号码格式错误", body.UserName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, "")
			if !uniqueEmail {
				msg := fmt.Sprintf("新增用户【%s】失败，邮箱已存在", body.UserName)
				c.JSON(200, resp.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("新增用户【%s】失败，邮箱格式错误", body.UserName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	userInfo := model.SysUser{
		UserName:   body.UserName,
		Password:   body.Password,
		NickName:   body.NickName,
		Email:      body.Email,
		Phone:      body.Phone,
		Sex:        body.Sex,
		StatusFlag: body.StatusFlag,
		Remark:     body.Remark,
		DeptId:     body.DeptId,  // 部门ID
		RoleIds:    body.RoleIds, // 角色ID组
		PostIds:    body.PostIds, // 岗位ID组
		Avatar:     body.Avatar,
		CreateBy:   reqctx.LoginUserToUserName(c),
	}
	insertId := s.sysUserService.Insert(userInfo)
	if insertId != "" {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 用户信息修改
//
// POST /
func (s SysUserController) Edit(c *gin.Context) {
	var body model.SysUser
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(400, resp.CodeMsg(40010, errMsgs))
		return
	}
	if body.UserId == "" {
		c.JSON(400, resp.CodeMsg(40010, "bind err: userId is empty"))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserId) {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserId)
	if userInfo.UserId != body.UserId {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, body.UserId)
			if !uniquePhone {
				msg := fmt.Sprintf("修改用户【%s】失败，手机号码已存在", body.UserName)
				c.JSON(200, resp.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，手机号码格式错误", body.UserName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, body.UserId)
			if !uniqueEmail {
				msg := fmt.Sprintf("修改用户【%s】失败，邮箱已存在", body.UserName)
				c.JSON(200, resp.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，邮箱格式错误", body.UserName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	if body.Avatar != "" {
		userInfo.Avatar = body.Avatar
	}

	userInfo.Phone = body.Phone
	userInfo.Email = body.Email
	userInfo.Sex = body.Sex
	userInfo.StatusFlag = body.StatusFlag
	userInfo.Remark = body.Remark
	userInfo.DeptId = body.DeptId
	userInfo.RoleIds = body.RoleIds
	userInfo.PostIds = body.PostIds
	userInfo.Password = "" // 忽略修改密码
	userInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysUserService.UpdateUserAndRolePost(userInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 用户信息删除
//
// DELETE /:userId
func (s SysUserController) Remove(c *gin.Context) {
	userIdsStr := c.Param("userId")
	userIds := parse.RemoveDuplicatesToArray(userIdsStr, ",")
	if userIdsStr == "" || len(userIds) <= 0 {
		c.JSON(400, resp.CodeMsg(40010, "bind err: userId is empty"))
		return
	}

	loginUserID := reqctx.LoginUserToUserID(c)
	for _, id := range userIds {
		// 不能删除自己
		if id == loginUserID {
			c.JSON(200, resp.ErrMsg("当前用户不能删除"))
			return
		}
		// 检查是否管理员用户
		if config.IsSystemUser(id) {
			c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
			return
		}
	}

	rows, err := s.sysUserService.DeleteByIds(userIds)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// Password 用户密码修改
//
// PUT /password
func (s SysUserController) Password(c *gin.Context) {
	var body struct {
		UserId   string `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(400, resp.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserId) {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserId)
	if userInfo.UserId != body.UserId {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	if !regular.ValidPassword(body.Password) {
		c.JSON(200, resp.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}

	userInfo.Password = body.Password
	userInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Status 用户状态修改
//
// PUT /status
func (s SysUserController) Status(c *gin.Context) {
	var body struct {
		UserId     string `json:"userId" binding:"required"`
		StatusFlag string `json:"statusFlag" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(400, resp.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserId) {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserId)
	if userInfo.UserId != body.UserId {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 与旧值相等不变更
	if userInfo.StatusFlag == body.StatusFlag {
		c.JSON(200, resp.ErrMsg("变更状态与旧值相等！"))
		return
	}

	userInfo.StatusFlag = body.StatusFlag
	userInfo.Password = "" // 密码不更新
	userInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Export 用户信息列表导出
//
// GET /export
func (s SysUserController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	queryMap := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysUserService.FindByPage(queryMap, dataScopeSQL)
	if total == 0 {
		c.JSON(200, resp.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("user_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysUserService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}

// Template 用户信息列表导入模板下载
//
// GET /import/template
func (s SysUserController) Template(c *gin.Context) {
	// 从 embed.FS 中读取内嵌文件
	assetsDir := config.GetAssetsDirFS()
	fileData, err := assetsDir.ReadFile("src/assets/template/excel/user_import_template.xlsx")
	if err != nil {
		c.String(400, "failed to read file")
		return
	}
	fileName := fmt.Sprintf("user_import_template_%d.xlsx", time.Now().UnixMilli())

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")

	// 返回响应体
	c.Data(200, "application/octet-stream", fileData)
}

// Import 用户信息列表导入
//
// POST /import
func (s SysUserController) Import(c *gin.Context) {
	var body struct {
		FilePath string `json:"filePath" binding:"required"` // 上传的文件地址
		Update   bool   `json:"update"`                      // 允许进行更新
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(400, resp.CodeMsg(40010, errMsgs))
		return
	}

	// 表格文件绝对地址
	filePath := file.ParseUploadFileAbsPath(body.FilePath)
	// 读取表格数据
	rows, err := file.ReadSheet(filePath, "")
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	// 获取操作人名称
	operaName := reqctx.LoginUserToUserName(c)
	// 导入数据
	message := s.sysUserService.ImportData(rows, operaName, body.Update)
	c.JSON(200, resp.OkMsg(message))
}
