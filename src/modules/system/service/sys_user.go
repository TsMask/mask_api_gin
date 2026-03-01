package service

import (
	"fmt"
	"strings"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysUser 实例化服务层
var NewSysUser = &SysUser{
	sysUserRepository:     repository.NewSysUser,
	sysRoleRepository:     repository.NewSysRole,
	sysDeptRepository:     repository.NewSysDept,
	sysUserRoleRepository: repository.NewSysUserRole,
	sysUserPostRepository: repository.NewSysUserPost,
	sysDictTypeService:    NewSysDictType,
	sysDictDataService:    NewSysDictData,
	sysConfigService:      NewSysConfig,
}

// SysUser 用户 服务层处理
type SysUser struct {
	sysUserRepository     *repository.SysUser     // 用户服务
	sysRoleRepository     *repository.SysRole     // 角色服务
	sysDeptRepository     *repository.SysDept     // 部门服务
	sysUserRoleRepository *repository.SysUserRole // 用户与角色服务
	sysUserPostRepository *repository.SysUserPost // 用户与岗位服务
	sysDictTypeService    *SysDictType            // 字典类型服务
	sysDictDataService    *SysDictData            // 字典数据服务
	sysConfigService      *SysConfig              // 参数配置服务
}

// FindByPage 分页查询列表数据
func (s SysUser) FindByPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64) {
	rows, total := s.sysUserRepository.SelectByPage(query, dataScopeSQL)
	for i, v := range rows {
		// 部门
		if rows[i].DeptId > 0 && rows[i].Dept == nil {
			deptInfo := s.sysDeptRepository.SelectById(v.DeptId)
			rows[i].Dept = &deptInfo
		}
		// 角色
		roleArr := s.sysRoleRepository.SelectByUserId(v.UserId)
		roleIds := make([]int64, 0)
		roles := make([]*model.SysRole, 0)
		for _, role := range roleArr {
			roles = append(roles, &role)
			roleIds = append(roleIds, role.RoleId)
		}
		rows[i].Roles = roles
		rows[i].RoleIds = roleIds
	}
	return rows, total
}

// Find 查询数据
func (s SysUser) Find(sysUser model.SysUser) []model.SysUser {
	return s.sysUserRepository.Select(sysUser)
}

// FindById 通过ID查询信息
func (s SysUser) FindById(userId int64) model.SysUser {
	userInfo := model.SysUser{}
	if userId <= 0 {
		return userInfo
	}
	users := s.sysUserRepository.SelectByIds([]int64{userId})
	if len(users) > 0 {
		userInfo = users[0]
		// 部门
		if userInfo.DeptId > 0 && userInfo.Dept == nil {
			deptInfo := s.sysDeptRepository.SelectById(userInfo.DeptId)
			userInfo.Dept = &deptInfo
		}
		// 角色
		roleArr := s.sysRoleRepository.SelectByUserId(userInfo.UserId)
		roleIds := make([]int64, 0)
		roles := make([]*model.SysRole, 0)
		for _, role := range roleArr {
			roles = append(roles, &role)
			roleIds = append(roleIds, role.RoleId)
		}
		userInfo.Roles = roles
		userInfo.RoleIds = roleIds
	}
	return userInfo
}

// Insert 新增信息
func (s SysUser) Insert(sysUser model.SysUser) int64 {
	// 新增用户信息
	insertId := s.sysUserRepository.Insert(sysUser)
	if insertId > 0 {
		s.insertUserRole(insertId, sysUser.RoleIds) // 新增用户角色信息
		s.insertUserPost(insertId, sysUser.PostIds) // 新增用户岗位信息
	}
	return insertId
}

// insertUserRole 新增用户角色信息
func (s SysUser) insertUserRole(userId int64, roleIds []int64) int64 {
	if userId <= 0 || len(roleIds) <= 0 {
		return 0
	}

	var arr []model.SysUserRole
	for _, roleId := range roleIds {
		// 系统管理员角色禁止操作，只能通过配置指定用户ID分配
		if roleId <= 0 || roleId == constants.SYS_ROLE_SYSTEM_ID {
			continue
		}
		arr = append(arr, model.SysUserRole{
			UserId: userId, RoleId: roleId,
		})
	}

	return s.sysUserRoleRepository.BatchInsert(arr)
}

// insertUserPost 新增用户岗位信息
func (s SysUser) insertUserPost(userId int64, postIds []int64) int64 {
	if userId <= 0 || len(postIds) <= 0 {
		return 0
	}

	var arr []model.SysUserPost
	for _, postId := range postIds {
		if postId <= 0 {
			continue
		}
		arr = append(arr, model.SysUserPost{
			UserId: userId, PostId: postId,
		})
	}

	return s.sysUserPostRepository.BatchInsert(arr)
}

// Update 修改信息
func (s SysUser) Update(sysUser model.SysUser) int64 {
	return s.sysUserRepository.Update(sysUser)
}

// UpdateUserAndRolePost 修改用户信息同时更新角色和岗位
func (s SysUser) UpdateUserAndRolePost(sysUser model.SysUser) int64 {
	// 删除用户与角色关联
	s.sysUserRoleRepository.DeleteByUserIds([]int64{sysUser.UserId})
	// 新增用户角色信息
	s.insertUserRole(sysUser.UserId, sysUser.RoleIds)
	// 删除用户与岗位关联
	s.sysUserPostRepository.DeleteByUserIds([]int64{sysUser.UserId})
	// 新增用户岗位信息
	s.insertUserPost(sysUser.UserId, sysUser.PostIds)
	return s.sysUserRepository.Update(sysUser)
}

// DeleteByIds 批量删除信息
func (s SysUser) DeleteByIds(userIds []int64) (int64, error) {
	// 检查是否存在
	users := s.sysUserRepository.SelectByIds(userIds)
	if len(users) <= 0 {
		return 0, fmt.Errorf("没有权限访问用户数据！")
	}
	if len(users) == len(userIds) {
		s.sysUserRoleRepository.DeleteByUserIds(userIds) // 删除用户与角色关联
		s.sysUserPostRepository.DeleteByUserIds(userIds) // 删除用户与岗位关联
		return s.sysUserRepository.DeleteByIds(userIds), nil
	}
	return 0, fmt.Errorf("删除用户信息失败！")
}

// CheckUniqueByUserName 检查用户名称是否唯一
func (s SysUser) CheckUniqueByUserName(userName string, userId int64) bool {
	uniqueId := s.sysUserRepository.CheckUnique(model.SysUser{
		UserName: userName,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueByPhone 检查手机号码是否唯一
func (s SysUser) CheckUniqueByPhone(phone string, userId int64) bool {
	uniqueId := s.sysUserRepository.CheckUnique(model.SysUser{
		Phone: phone,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueByEmail 检查Email是否唯一
func (s SysUser) CheckUniqueByEmail(email string, userId int64) bool {
	uniqueId := s.sysUserRepository.CheckUnique(model.SysUser{
		Email: email,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == 0
}

// FindByUserName 通过用户名查询用户信息
func (s SysUser) FindByUserName(userName string) model.SysUser {
	userinfo := s.sysUserRepository.SelectByUserName(userName)
	if userinfo.UserName != userName {
		return userinfo
	}
	// 部门
	deptInfo := s.sysDeptRepository.SelectById(userinfo.DeptId)
	userinfo.Dept = &deptInfo
	// 角色
	roleArr := s.sysRoleRepository.SelectByUserId(userinfo.UserId)
	roles := make([]*model.SysRole, 0)
	roleIds := make([]int64, 0)
	for _, role := range roleArr {
		roles = append(roles, &role)
		roleIds = append(roleIds, role.RoleId)
	}
	userinfo.Roles = roles
	userinfo.RoleIds = roleIds
	return userinfo
}

// FindAuthUsersPage 根据条件分页查询分配用户角色列表
func (s SysUser) FindAuthUsersPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64) {
	rows, total := s.sysUserRepository.SelectAuthUsersByPage(query, dataScopeSQL)
	for i, v := range rows {
		// 部门
		deptInfo := s.sysDeptRepository.SelectById(v.DeptId)
		rows[i].Dept = &deptInfo
		// 角色
		roleArr := s.sysRoleRepository.SelectByUserId(v.UserId)
		roleIds := make([]int64, 0)
		roles := make([]*model.SysRole, 0)
		for _, role := range roleArr {
			roles = append(roles, &role)
			roleIds = append(roleIds, role.RoleId)
		}
		rows[i].Roles = roles
		rows[i].RoleIds = roleIds
	}
	return rows, total
}

// ExportData 导出数据表格
func (s SysUser) ExportData(rows []model.SysUser, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "用户编号",
		"B1": "登录账号",
		"C1": "用户昵称",
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
		if row.StatusFlag == constants.STATUS_YES {
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
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}

// ImportData 导入数据表格
func (s SysUser) ImportData(rows []map[string]string, operaName string, updateSupport bool) string {
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
		var sysUserDeptId int64 = 100
		if row["H"] != "" {
			sysUserDeptId = parse.Number(row["H"])
		}

		// 验证是否存在这个用户
		newSysUser := s.FindByUserName(row["B"])
		newSysUser.Password = initPassword
		newSysUser.UserName = row["B"]
		newSysUser.NickName = row["C"]
		newSysUser.Phone = row["E"]
		newSysUser.Email = row["D"]
		newSysUser.StatusFlag = sysUserStatus
		newSysUser.Sex = sysUserSex
		newSysUser.DeptId = sysUserDeptId

		// 行用户编号
		rowNo := row["A"]

		// 检查手机号码格式并判断是否唯一
		if newSysUser.Phone != "" {
			if regular.ValidMobile(newSysUser.Phone) {
				uniquePhone := s.CheckUniqueByPhone(newSysUser.Phone, newSysUser.UserId)
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
				uniqueEmail := s.CheckUniqueByEmail(newSysUser.Email, newSysUser.UserId)
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

		if newSysUser.UserId <= 0 {
			newSysUser.CreateBy = operaName
			if s.Insert(newSysUser) > 0 {
				msg := fmt.Sprintf("用户编号：%s 登录账号：%s 导入成功", rowNo, newSysUser.UserName)
				successNum++
				successMsgArr = append(successMsgArr, msg)
			} else {
				msg := fmt.Sprintf("用户编号：%s 登录账号：%s 导入失败", rowNo, newSysUser.UserName)
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
			}
			continue
		}

		// 如果用户已存在 同时 是否更新支持
		if newSysUser.UserId > 0 && updateSupport {
			newSysUser.Password = "" // 密码不更新
			newSysUser.UpdateBy = operaName
			rows := s.Update(newSysUser)
			if rows > 0 {
				msg := fmt.Sprintf("用户编号：%s 登录账号：%s 更新成功", rowNo, newSysUser.UserName)
				successNum++
				successMsgArr = append(successMsgArr, msg)
			} else {
				msg := fmt.Sprintf("用户编号：%s 登录账号：%s 更新失败", rowNo, newSysUser.UserName)
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
	return message
}
