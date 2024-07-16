package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constAdmin "mask_api_gin/src/framework/constants/admin"
	constCommon "mask_api_gin/src/framework/constants/common"
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

// NewSysUser 实例化控制层
var NewSysUser = &SysUserController{
	sysUserService:     service.NewSysUser,
	sysRoleService:     service.NewSysRole,
	sysPostService:     service.NewSysPost,
	sysDictTypeService: service.NewSysDictType,
	sysConfigService:   service.NewSysConfig,
}

// SysUserController 用户信息
//
// PATH /system/user
type SysUserController struct {
	sysUserService     service.ISysUserService     // 用户服务
	sysRoleService     service.ISysRoleService     // 角色服务
	sysPostService     service.ISysPostService     // 岗位服务
	sysDictTypeService service.ISysDictTypeService // 字典类型服务
	sysConfigService   service.ISysConfigService   // 参数配置服务
}

// List 用户信息列表
//
// GET /list
func (s *SysUserController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	data := s.sysUserService.FindByPage(query, dataScopeSQL)
	c.JSON(200, result.Ok(data))
}

// Info 用户信息详情
//
// GET /:userId
func (s *SysUserController) Info(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 查询系统角色列表
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	roles := s.sysRoleService.Find(model.SysRole{}, dataScopeSQL)

	// 不是系统指定管理员需要排除其角色
	if !config.IsAdmin(userId) {
		rolesFilter := make([]model.SysRole, 0)
		for _, r := range roles {
			if r.RoleID != constAdmin.RoleId {
				rolesFilter = append(rolesFilter, r)
			}
		}
		roles = rolesFilter
	}

	// 查询系统岗位列表
	posts := s.sysPostService.Find(model.SysPost{})

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
	user := s.sysUserService.FindById(userId)
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
	userPosts := s.sysPostService.FindByUserId(userId)
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

// Add 用户信息新增
//
// POST /
func (s *SysUserController) Add(c *gin.Context) {
	var body model.SysUser
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.UserID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 密码单独取，避免序列化输出
	var bodyPassword struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWith(&bodyPassword, binding.JSON); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	body.Password = bodyPassword.Password

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(body.UserName, "")
	if !uniqueUserName {
		msg := fmt.Sprintf("新增用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, "")
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
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, "")
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
	insertId := s.sysUserService.Insert(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Edit 用户信息修改
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

	// 检查是否存在
	user := s.sysUserService.FindById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(body.UserName, body.UserID)
	if !uniqueUserName {
		msg := fmt.Sprintf("修改用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, body.UserID)
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
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, body.UserID)
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
	body.LoginIP = ""  // 忽略登录IP
	body.LoginDate = 0 // 忽略登录时间
	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.UpdateUserAndRolePost(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 用户信息删除
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

	// 检查是否管理员用户
	loginUserID := ctx.LoginUserToUserID(c)
	for _, id := range uniqueIDs {
		if id == loginUserID {
			c.JSON(200, result.ErrMsg("当前用户不能删除"))
			return
		}
		if config.IsAdmin(id) {
			c.JSON(200, result.ErrMsg("不允许操作管理员用户"))
			return
		}
	}

	rows, err := s.sysUserService.DeleteByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// ResetPwd 用户重置密码
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

	// 检查是否存在
	user := s.sysUserService.FindById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	if !regular.ValidPassword(body.Password) {
		c.JSON(200, result.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}

	user.Password = body.Password
	user.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(user)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Status 用户状态修改
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

	// 检查是否管理员用户
	if config.IsAdmin(body.UserID) {
		c.JSON(200, result.ErrMsg("不允许操作管理员用户"))
		return
	}

	// 检查是否存在
	user := s.sysUserService.FindById(body.UserID)
	if user.UserID != body.UserID {
		c.JSON(200, result.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 与旧值相等不变更
	if user.Status == body.Status {
		c.JSON(200, result.ErrMsg("变更状态与旧值相等！"))
		return
	}

	user.Status = body.Status
	user.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(user)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Export 用户信息列表导出
//
// POST /export
func (s *SysUserController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	queryMap := ctx.BodyJSONMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	data := s.sysUserService.FindByPage(queryMap, dataScopeSQL)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysUser)

	// 导出文件名称
	fileName := fmt.Sprintf("user_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "用户编号",
		"B1": "登录名称",
		"C1": "用户名称",
		"D1": "用户邮箱",
		"E1": "手机号码",
		"F1": "用户性别",
		"G1": "帐号状态",
		"H1": "部门编号",
		"I1": "部门名称",
		"J1": "部门负责人",
		"K1": "最后登录IP",
		"L1": "最后登录时间",
	}
	// 读取用户性别字典数据
	dictSysUserSex := s.sysDictTypeService.FindDataByType("sys_user_sex")
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
		// 账号状态
		statusValue := "停用"
		if row.Status == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.UserID,
			"B" + idx: row.UserName,
			"C" + idx: row.NickName,
			"D" + idx: row.Email,
			"E" + idx: row.Phone,
			"F" + idx: sysUserSex,
			"G" + idx: statusValue,
			"H" + idx: row.Dept.DeptID,
			"I" + idx: row.Dept.DeptName,
			"J" + idx: row.Dept.Leader,
			"K" + idx: row.LoginIP,
			"L" + idx: date.ParseDateToStr(row.LoginDate, date.YYYY_MM_DD_HH_MM_SS),
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

// Template 用户信息列表导入模板下载
//
// GET /importTemplate
func (s *SysUserController) Template(c *gin.Context) {
	fileName := fmt.Sprintf("user_import_template_%d.xlsx", time.Now().UnixMilli())
	asserPath := "assets/template/excel/user_import_template.xlsx"
	c.FileAttachment(asserPath, fileName)
}

// ImportData 用户信息列表导入
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
	filePath, err := file.TransferExcelUploadFile(formFile)
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
	operaName := ctx.LoginUserToUserName(c)
	isUpdateSupport := parse.Boolean(updateSupport)

	// 读取默认初始密码
	initPassword := s.sysConfigService.FindValueByKey("sys.user.initPassword")
	// 读取用户性别字典数据
	dictSysUserSex := s.sysDictTypeService.FindDataByType("sys_user_sex")

	// 导入记录
	successNum := 0
	failureNum := 0
	var successMsgArr []string
	var failureMsgArr []string
	mustItemArr := []string{"B", "C"}
	for _, row := range rows {
		// 检查必填列
		ownItem := true
		for _, item := range mustItemArr {
			if v, ok := row[item]; !ok || v == "" {
				ownItem = false
				break
			}
		}
		if !ownItem {
			mustItemArrStr := strings.Join(mustItemArr, "、")
			failureNum++
			msg := fmt.Sprintf("表格中必填列表项，%s", mustItemArrStr)
			failureMsgArr = append(failureMsgArr, msg)
			continue
		}

		// 用户性别转值
		sysUserSex := "0"
		for _, v := range dictSysUserSex {
			if row["F"] == v.DictLabel {
				sysUserSex = v.DictValue
				break
			}
		}
		sysUserStatus := constCommon.StatusNo
		if row["G"] == "正常" {
			sysUserStatus = constCommon.StatusYes
		}

		// 验证是否存在这个用户
		newSysUser := s.sysUserService.FindByUserName(row["B"])
		newSysUser.UserType = "sys"
		newSysUser.Password = initPassword
		newSysUser.DeptID = row["H"]
		newSysUser.UserName = row["B"]
		newSysUser.NickName = row["C"]
		newSysUser.Phone = row["E"]
		newSysUser.Email = row["D"]
		newSysUser.Status = sysUserStatus
		newSysUser.Sex = sysUserSex

		// 行用户编号
		rowNo := row["A"]

		// 检查手机号码格式并判断是否唯一
		if newSysUser.Phone != "" {
			if regular.ValidMobile(newSysUser.Phone) {
				uniquePhone := s.sysUserService.CheckUniqueByPhone(newSysUser.Phone, "")
				if !uniquePhone {
					msg := fmt.Sprintf("用户编号：%s 手机号码：%s 已存在", rowNo, newSysUser.Phone)
					failureNum++
					failureMsgArr = append(failureMsgArr, msg)
					continue
				}
			} else {
				msg := fmt.Sprintf("用户编号：%s 手机号码：%s 格式错误", rowNo, newSysUser.Phone)
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
				continue
			}
		}

		// 检查邮箱格式并判断是否唯一
		if newSysUser.Email != "" {
			if regular.ValidEmail(newSysUser.Email) {
				uniqueEmail := s.sysUserService.CheckUniqueByEmail(newSysUser.Email, "")
				if !uniqueEmail {
					msg := fmt.Sprintf("用户编号：%s 用户邮箱：%s 已存在", rowNo, newSysUser.Email)
					failureNum++
					failureMsgArr = append(failureMsgArr, msg)
					continue
				}
			} else {
				msg := fmt.Sprintf("用户编号：%s 用户邮箱：%s 格式错误", rowNo, newSysUser.Email)
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
				continue
			}
		}

		if newSysUser.UserID == "" {
			newSysUser.CreateBy = operaName
			insertId := s.sysUserService.Insert(newSysUser)
			if insertId != "" {
				msg := fmt.Sprintf("用户编号：%s 登录名称：%s 导入成功", rowNo, newSysUser.UserName)
				successNum++
				successMsgArr = append(successMsgArr, msg)
			} else {
				msg := fmt.Sprintf("用户编号：%s 登录名称：%s 导入失败", rowNo, newSysUser.UserName)
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
			}
			continue
		}

		// 如果用户已存在 同时 是否更新支持
		if newSysUser.UserID != "" && isUpdateSupport {
			newSysUser.UpdateBy = operaName
			rows := s.sysUserService.Update(newSysUser)
			if rows > 0 {
				msg := fmt.Sprintf("用户编号：%s 登录名称：%s 更新成功", rowNo, newSysUser.UserName)
				successNum++
				successMsgArr = append(successMsgArr, msg)
			} else {
				msg := fmt.Sprintf("用户编号：%s 登录名称：%s 更新失败", rowNo, newSysUser.UserName)
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
			}
			continue
		}
	}

	message := ""
	if failureNum > 0 {
		msg := fmt.Sprintf("很抱歉，导入失败！共 %d 条数据格式不正确，错误如下：", failureNum)
		message = strings.Join(append([]string{msg}, failureMsgArr...), "<br/>")
	} else {
		msg := fmt.Sprintf("恭喜您，数据已全部导入成功！共 %d 条，数据如下：", successNum)
		message = strings.Join(append([]string{msg}, successMsgArr...), "<br/>")
	}

	c.JSON(200, result.OkMsg(message))
}
