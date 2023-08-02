package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants/admin"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 实例化控制层 SysUserController 结构体
var NewSysUser = &SysUserController{
	sysUserService:     service.NewSysUserImpl,
	sysRoleService:     service.NewSysRoleImpl,
	sysPostService:     service.NewSysPostImpl,
	sysDictDataService: service.NewSysDictDataImpl,
}

// 用户信息
//
// PATH /system/user
type SysUserController struct {
	// 用户服务
	sysUserService service.ISysUser
	// 角色服务
	sysRoleService service.ISysRole
	// 岗位服务
	sysPostService service.ISysPost
	// 字典数据服务
	sysDictDataService service.ISysDictData
}

// 用户信息列表
//
// GET /list
func (s *SysUserController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	data := s.sysUserService.SelectUserPage(querys, dataScopeSQL)
	c.JSON(200, result.Ok(data))
}

// 用户信息详情
//
// GET /:userId
func (s *SysUserController) Info(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 查询系统角色列表
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
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
		c.JSON(200, result.OkData(map[string]any{
			"user":    map[string]any{},
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

	c.JSON(200, result.OkData(map[string]any{
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
func (s *SysUserController) Add(c *gin.Context) {
	var body model.SysUser
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.UserID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueUserName(body.UserName, "")
	if !uniqueUserName {
		msg := fmt.Sprintf("新增用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.PhoneNumber != "" {
		if regular.ValidMobile(body.PhoneNumber) {
			uniquePhone := s.sysUserService.CheckUniquePhone(body.PhoneNumber, "")
			if !uniquePhone {
				msg := fmt.Sprintf("新增用户【%s】失败，手机号码已存在", body.UserName)
				c.JSON(200, result.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("新增用户【%s】失败，手机号码格式错误", body.UserName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueEmail(body.Email, "")
			if !uniqueEmail {
				msg := fmt.Sprintf("新增用户【%s】失败，邮箱已存在", body.UserName)
				c.JSON(200, result.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("新增用户【%s】失败，邮箱格式错误", body.UserName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysUserService.InsertUser(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 用户信息修改
//
// POST /
func (s *SysUserController) Edit(c *gin.Context) {
	var body model.SysUser
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.UserID == "" {
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

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueUserName(body.UserName, body.UserID)
	if !uniqueUserName {
		msg := fmt.Sprintf("修改用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.PhoneNumber != "" {
		if regular.ValidMobile(body.PhoneNumber) {
			uniquePhone := s.sysUserService.CheckUniquePhone(body.PhoneNumber, body.UserID)
			if !uniquePhone {
				msg := fmt.Sprintf("修改用户【%s】失败，手机号码已存在", body.UserName)
				c.JSON(200, result.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，手机号码格式错误", body.UserName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueEmail(body.Email, body.UserID)
			if !uniqueEmail {
				msg := fmt.Sprintf("修改用户【%s】失败，邮箱已存在", body.UserName)
				c.JSON(200, result.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，邮箱格式错误", body.UserName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	body.UserName = "" // 忽略修改登录用户名称
	body.Password = "" // 忽略修改密码
	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.UpdateUserAndRolePost(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 用户信息删除
//
// DELETE /:userIds
func (s *SysUserController) Remove(c *gin.Context) {
	userIds := c.Param("userIds")
	if userIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(userIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
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
func (s *SysUserController) ResetPwd(c *gin.Context) {
	var body struct {
		UserID   string `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
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

	userName := ctx.LoginUserToUserName(c)
	SysUserController := model.SysUser{
		UserID:   body.UserID,
		Password: body.Password,
		UpdateBy: userName,
	}
	rows := s.sysUserService.UpdateUser(SysUserController)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 用户状态修改
//
// PUT /changeStatus
func (s *SysUserController) Status(c *gin.Context) {
	var body struct {
		UserID string `json:"userId" binding:"required"`
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
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

	userName := ctx.LoginUserToUserName(c)
	SysUserController := model.SysUser{
		UserID:   body.UserID,
		Status:   body.Status,
		UpdateBy: userName,
	}
	rows := s.sysUserService.UpdateUser(SysUserController)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 用户信息列表导出
//
// POST /export
func (s *SysUserController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	data := s.sysUserService.SelectUserPage(querys, dataScopeSQL)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysUser)

	// 导出文件名称
	fileName := fmt.Sprintf("user_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "用户序号",
		"B1": "登录名称",
		"C1": "用户名称",
		"D1": "用户邮箱",
		"E1": "手机号码",
		"F1": "用户性别",
		"G1": "帐号状态",
		"H1": "最后登录IP",
		"I1": "最后登录时间",
		"J1": "部门名称",
		"K1": "部门负责人",
	}
	// 读取用户性别字典数据
	dictSysUserSex := s.sysDictDataService.SelectDictDataByType("sys_user_sex")
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 用户性别
		sysUserSex := "未知"
		for _, v := range dictSysUserSex {
			if row.Sex == v.DictValue {
				sysUserSex = v.DictLabel
				break
			}
		}
		// 帐号状态
		statusValue := "停用"
		if row.Status == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.UserID,
			"B" + idx: row.UserName,
			"C" + idx: row.NickName,
			"D" + idx: row.Email,
			"E" + idx: row.PhoneNumber,
			"F" + idx: sysUserSex,
			"G" + idx: statusValue,
			"H" + idx: row.LoginIP,
			"I" + idx: date.ParseDateToStr(row.LoginDate, date.YYYY_MM_DD_HH_MM_SS),
			"J" + idx: row.Dept.DeptName,
			"K" + idx: row.Dept.Leader,
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}

// 用户信息列表导入模板下载
//
// GET /importTemplate
func (s *SysUserController) Template(c *gin.Context) {
	fileName := fmt.Sprintf("user_import_template_%d.xlsx", time.Now().UnixMilli())
	asserPath := "assets/template/excel/user_import_template.xlsx"
	c.FileAttachment(asserPath, fileName)
}

// 用户信息列表导入
//
// POST /importData
func (s *SysUserController) ImportData(c *gin.Context) {
	// 允许进行更新
	updateSupport := c.PostForm("updateSupport")
	// 上传的文件
	formFile, err := c.FormFile("file")
	if err != nil || updateSupport == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 保存表格文件
	filePath, err := file.TransferExeclUploadFile(formFile)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 读取表格数据
	rows, err := file.ReadSheet(filePath, "")
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 获取操作人名称
	operName := ctx.LoginUserToUserName(c)
	isUpdateSupport := parse.Boolean(updateSupport)
	message := s.sysUserService.ImportUser(rows, isUpdateSupport, operName)
	c.JSON(200, result.OkMsg(message))
}
