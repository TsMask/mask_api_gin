package controller

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	sysUserService     *service.SysUser     // 用户服务
	sysRoleService     *service.SysRole     // 角色服务
	sysPostService     *service.SysPost     // 岗位服务
	sysDictTypeService *service.SysDictType // 字典类型服务
	sysConfigService   *service.SysConfig   // 参数配置服务
}

// List 用户信息列表
//
// GET /list
func (s SysUserController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	rows, total := s.sysUserService.FindByPage(query, dataScopeSQL)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 用户信息详情
//
// GET /:userId
func (s SysUserController) Info(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: userId is empty"))
		return
	}

	// 查询系统角色列表
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	roles := s.sysRoleService.Find(model.SysRole{}, dataScopeSQL)
	// 查询系统岗位列表
	posts := s.sysPostService.Find(model.SysPost{})

	// 新增用户时，用户ID为0
	if userId == "0" {
		c.JSON(200, response.OkData(map[string]any{
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
		c.JSON(200, response.ErrMsg("没有权限访问用户数据！"))
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

	c.JSON(200, response.OkData(map[string]any{
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
		c.JSON(400, response.ErrMsg(response.FormatBindError(err)))
		return
	}
	if body.UserId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: userId not is empty"))
		return
	}

	// 密码单独取，避免序列化输出
	var bodyPassword struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&bodyPassword); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	body.Password = bodyPassword.Password

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(body.UserName, "")
	if !uniqueUserName {
		msg := fmt.Sprintf("新增用户【%s】失败，登录账号已存在", body.UserName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, "")
			if !uniquePhone {
				msg := fmt.Sprintf("新增用户【%s】失败，手机号码已存在", body.UserName)
				c.JSON(200, response.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("新增用户【%s】失败，手机号码格式错误", body.UserName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, "")
			if !uniqueEmail {
				msg := fmt.Sprintf("新增用户【%s】失败，邮箱已存在", body.UserName)
				c.JSON(200, response.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("新增用户【%s】失败，邮箱格式错误", body.UserName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysUserService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 用户信息修改
//
// POST /
func (s SysUserController) Edit(c *gin.Context) {
	var body model.SysUser
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.UserId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: userId is empty"))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserId) {
		c.JSON(200, response.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserId)
	if userInfo.UserId != body.UserId {
		c.JSON(200, response.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 检查手机号码格式并判断是否唯一
	if body.Phone != "" {
		if regular.ValidMobile(body.Phone) {
			uniquePhone := s.sysUserService.CheckUniqueByPhone(body.Phone, body.UserId)
			if !uniquePhone {
				msg := fmt.Sprintf("修改用户【%s】失败，手机号码已存在", body.UserName)
				c.JSON(200, response.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，手机号码格式错误", body.UserName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	// 检查邮箱格式并判断是否唯一
	if body.Email != "" {
		if regular.ValidEmail(body.Email) {
			uniqueEmail := s.sysUserService.CheckUniqueByEmail(body.Email, body.UserId)
			if !uniqueEmail {
				msg := fmt.Sprintf("修改用户【%s】失败，邮箱已存在", body.UserName)
				c.JSON(200, response.ErrMsg(msg))
				return
			}
		} else {
			msg := fmt.Sprintf("修改用户【%s】失败，邮箱格式错误", body.UserName)
			c.JSON(200, response.ErrMsg(msg))
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
	userInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.UpdateUserAndRolePost(userInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 用户信息删除
//
// DELETE /:userId
func (s SysUserController) Remove(c *gin.Context) {
	userIdsStr := c.Param("userId")
	userIds := parse.RemoveDuplicatesToArray(userIdsStr, ",")
	if userIdsStr == "" || len(userIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: userId is empty"))
		return
	}

	loginUserID := ctx.LoginUserToUserID(c)
	for _, id := range userIds {
		// 不能删除自己
		if id == loginUserID {
			c.JSON(200, response.ErrMsg("当前用户不能删除"))
			return
		}
		// 检查是否管理员用户
		if config.IsSystemUser(id) {
			c.JSON(200, response.ErrMsg("不允许操作系统管理员用户"))
			return
		}
	}

	rows, err := s.sysUserService.DeleteByIds(userIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
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
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserId) {
		c.JSON(200, response.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserId)
	if userInfo.UserId != body.UserId {
		c.JSON(200, response.ErrMsg("没有权限访问用户数据！"))
		return
	}

	if !regular.ValidPassword(body.Password) {
		c.JSON(200, response.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}

	userInfo.Password = body.Password
	userInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
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
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否系统管理员用户
	if config.IsSystemUser(body.UserId) {
		c.JSON(200, response.ErrMsg("不允许操作系统管理员用户"))
		return
	}

	// 检查是否存在
	userInfo := s.sysUserService.FindById(body.UserId)
	if userInfo.UserId != body.UserId {
		c.JSON(200, response.ErrMsg("没有权限访问用户数据！"))
		return
	}

	// 与旧值相等不变更
	if userInfo.StatusFlag == body.StatusFlag {
		c.JSON(200, response.ErrMsg("变更状态与旧值相等！"))
		return
	}

	userInfo.StatusFlag = body.StatusFlag
	userInfo.Password = "" // 密码不更新
	userInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysUserService.Update(userInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Export 用户信息列表导出
//
// GET /export
func (s SysUserController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	queryMap := ctx.QueryMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	rows, total := s.sysUserService.FindByPage(queryMap, dataScopeSQL)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

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
		idx := fmt.Sprint(i + 2)
		// 用户性别
		sysUserSex := "未知"
		for _, v := range dictSysUserSex {
			if row.Sex == v.DataValue {
				sysUserSex = v.DataLabel
				break
			}
		}
		// 账号状态
		statusValue := "停用"
		if row.StatusFlag == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.UserId,
			"B" + idx: row.UserName,
			"C" + idx: row.NickName,
			"D" + idx: row.Email,
			"E" + idx: row.Phone,
			"F" + idx: sysUserSex,
			"G" + idx: statusValue,
			"H" + idx: row.Dept.DeptId,
			"I" + idx: row.Dept.DeptName,
			"J" + idx: row.Dept.Leader,
			"K" + idx: row.LoginIp,
			"L" + idx: date.ParseDateToStr(row.LoginTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
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
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 表格文件绝对地址
	filePath := file.ParseUploadFilePath(body.FilePath)
	// 读取表格数据
	rows, err := file.ReadSheet(filePath, "")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	// 获取操作人名称
	operaName := ctx.LoginUserToUserName(c)
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
			if row["F"] == v.DataLabel {
				sysUserSex = v.DataValue
				break
			}
		}
		sysUserStatus := constants.STATUS_NO
		if row["G"] == "正常" {
			sysUserStatus = constants.STATUS_YES
		}

		// 验证是否存在这个用户
		newSysUser := s.sysUserService.FindByUserName(row["B"])
		newSysUser.Password = initPassword
		newSysUser.UserName = row["B"]
		newSysUser.NickName = row["C"]
		newSysUser.Phone = row["E"]
		newSysUser.Email = row["D"]
		newSysUser.StatusFlag = sysUserStatus
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

		if newSysUser.UserId == "" {
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
		if newSysUser.UserId != "" && body.Update {
			newSysUser.Password = "" // 密码不更新
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

	c.JSON(200, response.OkMsg(message))
}
