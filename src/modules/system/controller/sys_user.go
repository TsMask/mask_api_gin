package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/file"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	authService "mask_api_gin/src/modules/auth/service"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
)

// NewSysUser 实例化控制层
var NewSysUser = &SysUserController{
	sysUserService: service.NewSysUser,
	sysRoleService: service.NewSysRole,
	sysPostService: service.NewSysPost,
	accountService: authService.NewAccount,
}

// SysUserController 用户信息 控制层处理
//
// PATH /system/user
type SysUserController struct {
	sysUserService *service.SysUser     // 用户服务
	sysRoleService *service.SysRole     // 角色服务
	sysPostService *service.SysPost     // 岗位服务
	accountService *authService.Account // 账号身份操作服务
}

// List 用户信息列表
//
// GET /list
func (s SysUserController) List(c *gin.Context) {
	queryMap := reqctx.QueryMap(c)
	dataScopeSQL := service.DataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysUserService.FindByPage(queryMap, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 用户信息详情
//
// GET /:id
func (s SysUserController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id < 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 查询系统角色列表
	roles := s.sysRoleService.Find(model.SysRole{})
	// 查询系统岗位列表
	posts := s.sysPostService.Find(model.SysPost{})

	// 新增用户时，用户ID为0
	if id == 0 {
		c.JSON(200, resp.OkData(map[string]any{
			"user":     map[string]any{},
			"role_ids": []string{},
			"post_ids": []string{},
			"roles":    roles,
			"posts":    posts,
		}))
		return
	}

	// 检查用户是否存在
	userInfo := s.sysUserService.FindById(id)
	if userInfo.ID != id {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 角色ID组
	roleIds := make([]int64, 0)
	for _, r := range userInfo.Roles {
		roleIds = append(roleIds, r.ID)
	}

	// 岗位ID组
	postIds := make([]int64, 0)
	userPosts := s.sysPostService.FindByUserId(id)
	for _, p := range userPosts {
		postIds = append(postIds, p.ID)
	}

	c.JSON(200, resp.OkData(map[string]any{
		"user":     userInfo,
		"role_ids": roleIds,
		"post_ids": postIds,
		"roles":    roles,
		"posts":    posts,
	}))
}

// Add 用户信息新增
//
// POST /
func (s SysUserController) Add(c *gin.Context) {
	var body model.SysUser
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, resp.FormatBindError(err)))
		return
	}
	// 密码单独取，避免序列化输出
	var bodyPassword struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&bodyPassword); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	body.UserPassword = bodyPassword.Password

	if body.ID > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: userId not is empty"))
		return
	}
	if !regular.ValidUsername(body.UserName) {
		msg := fmt.Sprintf("新增用户【%s】失败，登录账号用户账号只能包含大写小写字母，数字，且不少于4位", body.UserName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}
	if !regular.ValidPassword(body.UserPassword) {
		msg := fmt.Sprintf("新增用户【%s】失败，登录密码至少包含大小写字母、数字、特殊符号，且不少于6位", body.UserName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(body.UserName, 0)
	if !uniqueUserName {
		msg := fmt.Sprintf("新增用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.UserPhone != "" {
		if regular.ValidMobile(body.UserPhone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.UserPhone, 0)
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
	if body.UserEmail != "" {
		if regular.ValidEmail(body.UserEmail) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.UserEmail, 0)
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
		UserName:     body.UserName,
		UserPassword: body.UserPassword,
		UserNick:     body.UserNick,
		UserEmail:    body.UserEmail,
		UserPhone:    body.UserPhone,
		UserGender:   body.UserGender,
		UserStatus:   body.UserStatus,
		Remark:       body.Remark,
		DeptId:       body.DeptId,  // 部门ID
		RoleIds:      body.RoleIds, // 角色ID组
		PostIds:      body.PostIds, // 岗位ID组
		UserAvatar:   body.UserAvatar,
		CreateBy:     reqctx.LoginUserToUserName(c),
	}
	insertId := s.sysUserService.Insert(userInfo)
	if insertId > 0 {
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
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: userId is empty"))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.ID) {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.ID)
	if userInfo.ID != body.ID {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.UserPhone != "" {
		if regular.ValidMobile(body.UserPhone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.UserPhone, body.ID)
			if !uniquePhone {
				msg := fmt.Sprintf("修改用户【%s】失败，手机号码已存在", userInfo.UserName)
				c.JSON(200, resp.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，手机号码格式错误", userInfo.UserName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.UserEmail != "" {
		if regular.ValidEmail(body.UserEmail) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.UserEmail, body.ID)
			if !uniqueEmail {
				msg := fmt.Sprintf("修改用户【%s】失败，邮箱已存在", userInfo.UserName)
				c.JSON(200, resp.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，邮箱格式错误", userInfo.UserName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	if body.UserAvatar != "" {
		userInfo.UserAvatar = body.UserAvatar
	}

	userInfo.UserPhone = body.UserPhone
	userInfo.UserEmail = body.UserEmail
	userInfo.UserGender = body.UserGender
	userInfo.UserStatus = body.UserStatus
	userInfo.Remark = body.Remark
	userInfo.DeptId = body.DeptId
	userInfo.RoleIds = body.RoleIds
	userInfo.PostIds = body.PostIds
	userInfo.UserPassword = "" // 忽略修改密码
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
// DELETE /:id
func (s SysUserController) Remove(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(idStr, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	loginUserId := reqctx.LoginUserToUserID(c)
	for _, id := range ids {
		// 不能删除自己
		if id == loginUserId {
			c.JSON(200, resp.ErrMsg("当前用户不能删除"))
			return
		}
		// 检查是否管理员用户
		if config.IsSystemUser(id) {
			c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
			return
		}
	}

	rows, err := s.sysUserService.DeleteByIds(ids)
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
		UserID   int64  `json:"user_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserID) {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	if !regular.ValidPassword(body.Password) {
		c.JSON(200, resp.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserID)
	if userInfo.ID != body.UserID {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	userInfo.UserPassword = body.Password
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
		UserID     int64  `json:"user_id" binding:"required"`
		UserStatus string `json:"user_status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserID) {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserID)
	if userInfo.ID != body.UserID {
		c.JSON(200, resp.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 与旧值相等不变更
	if userInfo.UserStatus == body.UserStatus {
		c.JSON(200, resp.ErrMsg("变更状态与旧值相等！"))
		return
	}

	userInfo.UserStatus = body.UserStatus
	userInfo.UserPassword = "" // 密码不更新
	userInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Unlock 用户账户解锁
//
// PUT /unlock/:user_name
func (s SysUserController) Unlock(c *gin.Context) {
	userName := c.Param("user_name")
	if userName == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: userName is empty"))
		return
	}

	ok := s.accountService.CleanLoginRecordCache(userName)
	if ok {
		c.JSON(200, resp.OkMsg("unlock success"))
		return
	}
	c.JSON(200, resp.ErrMsg("not found unlock user"))
}

// Export 用户信息列表导出
//
// GET /export
func (s SysUserController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	queryMap := reqctx.QueryMap(c)
	dataScopeSQL := service.DataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysUserService.FindByPage(queryMap, dataScopeSQL)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("user_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
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
	fileData, err := assetsDir.ReadFile("src/assets/template/csv/user_import_template.csv")
	if err != nil {
		c.String(400, "failed to read file")
		return
	}
	fileName := fmt.Sprintf("user_import_template_%d.csv", time.Now().UnixMilli())

	reqctx.HeaderContentOctetStream(c, fileName, fileData)
}

// Import 用户信息列表导入
//
// POST /import
func (s SysUserController) Import(c *gin.Context) {
	var body struct {
		FilePath string `json:"file_path" binding:"required"` // 上传的文件地址
		Update   bool   `json:"update"`                       // 允许进行更新
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	// 表格文件绝对地址
	filePath := file.UploadFileAbsPath(body.FilePath)
	// 读取表格数据
	rows, err := file.CSVReadAll(filePath)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	if len(rows) <= 0 {
		c.JSON(200, resp.ErrMsg("导入用户数据不能为空！"))
		return
	}

	// 获取操作人名称
	operaName := reqctx.LoginUserToUserName(c)
	// 导入数据
	message := s.sysUserService.ImportData(rows, operaName, body.Update)
	c.JSON(200, resp.OkMsg(message))
}
