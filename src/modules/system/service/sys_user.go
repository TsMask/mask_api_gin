package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/token"
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

// fillUsersDeptAndRoles 批量填充用户的部门和角色信息（消除 N+1 查询）
func (s SysUser) fillUsersDeptAndRoles(users []model.SysUser) {
	if len(users) == 0 {
		return
	}

	// 收集所有唯一的 deptId 和 userId
	deptIdSet := make(map[int64]struct{})
	userIdSet := make(map[int64]struct{})
	for _, u := range users {
		if u.DeptId > 0 {
			deptIdSet[u.DeptId] = struct{}{}
		}
		userIdSet[u.UserId] = struct{}{}
	}

	// 批量查询部门
	deptIds := make([]int64, 0, len(deptIdSet))
	for id := range deptIdSet {
		deptIds = append(deptIds, id)
	}
	deptMap := s.sysDeptRepository.SelectMapByIds(deptIds)

	// 批量查询角色
	userIds := make([]int64, 0, len(userIdSet))
	for id := range userIdSet {
		userIds = append(userIds, id)
	}
	roleMap := s.sysRoleRepository.SelectByUserIds(userIds)

	// 填充
	for i := range users {
		if dept, ok := deptMap[users[i].DeptId]; ok {
			users[i].Dept = &dept
		}
		if roleArr, ok := roleMap[users[i].UserId]; ok {
			roles := make([]*model.SysRole, 0, len(roleArr))
			roleIds := make([]int64, 0, len(roleArr))
			for idx := range roleArr {
				roles = append(roles, &roleArr[idx])
				roleIds = append(roleIds, roleArr[idx].RoleId)
			}
			users[i].Roles = roles
			users[i].RoleIds = roleIds
		}
	}
}

// FindByPage 分页查询列表数据
func (s SysUser) FindByPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64) {
	rows, total := s.sysUserRepository.SelectByPage(query, dataScopeSQL)
	s.fillUsersDeptAndRoles(rows)
	return rows, total
}

// Find 查询数据
func (s SysUser) Find(sysUser model.SysUser) []model.SysUser {
	return s.sysUserRepository.Select(sysUser)
}

// FindById 通过ID查询信息
func (s SysUser) FindById(userId int64) model.SysUser {
	if userId <= 0 {
		return model.SysUser{}
	}
	users := s.sysUserRepository.SelectByIds([]int64{userId})
	if len(users) > 0 {
		s.fillUsersDeptAndRoles(users)
		return users[0]
	}
	return model.SysUser{}
}

// Insert 新增信息
func (s SysUser) Insert(sysUser model.SysUser) int64 {
	var insertId int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		insertId = s.sysUserRepository.Insert(tx, sysUser)
		if insertId <= 0 {
			return fmt.Errorf("insert user failed")
		}
		s.insertUserRole(tx, insertId, sysUser.RoleIds)
		s.insertUserPost(tx, insertId, sysUser.PostIds)
		return nil
	})
	if err != nil {
		return 0
	}
	return insertId
}

// insertUserRole 新增用户角色信息
func (s SysUser) insertUserRole(d *gorm.DB, userId int64, roleIds []int64) int64 {
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

	return s.sysUserRoleRepository.BatchInsert(d, arr)
}

// insertUserPost 新增用户岗位信息
func (s SysUser) insertUserPost(d *gorm.DB, userId int64, postIds []int64) int64 {
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

	return s.sysUserPostRepository.BatchInsert(d, arr)
}

// Update 修改信息
func (s SysUser) Update(sysUser model.SysUser) int64 {
	return s.sysUserRepository.Update(nil, sysUser)
}

// UpdateUserAndRolePost 修改用户信息同时更新角色和岗位
func (s SysUser) UpdateUserAndRolePost(sysUser model.SysUser) int64 {
	var rows int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		s.sysUserRoleRepository.DeleteByUserIds(tx, []int64{sysUser.UserId})
		s.insertUserRole(tx, sysUser.UserId, sysUser.RoleIds)
		s.sysUserPostRepository.DeleteByUserIds(tx, []int64{sysUser.UserId})
		s.insertUserPost(tx, sysUser.UserId, sysUser.PostIds)
		rows = s.sysUserRepository.Update(tx, sysUser)
		return nil
	})
	if err != nil {
		return 0
	}
	// 角色或状态变更后清理 Token 缓存，强制重新登录
	token.RemoveByUserIds([]int64{sysUser.UserId})
	return rows
}

// DeleteByIds 批量删除信息
func (s SysUser) DeleteByIds(userIds []int64) (int64, error) {
	// 检查是否存在
	users := s.sysUserRepository.SelectByIds(userIds)
	if len(users) <= 0 {
		return 0, fmt.Errorf("没有权限访问用户数据！")
	}
	if len(users) != len(userIds) {
		return 0, fmt.Errorf("删除用户信息失败！")
	}
	var rows int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		s.sysUserRoleRepository.DeleteByUserIds(tx, userIds)
		s.sysUserPostRepository.DeleteByUserIds(tx, userIds)
		rows = s.sysUserRepository.DeleteByIds(tx, userIds)
		return nil
	})
	if err != nil {
		return 0, err
	}
	// 删除用户后清理 Token 缓存
	token.RemoveByUserIds(userIds)
	return rows, nil
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
	s.fillUsersDeptAndRoles([]model.SysUser{userinfo})
	return userinfo
}

// FindAuthUsersPage 根据条件分页查询分配用户角色列表
func (s SysUser) FindAuthUsersPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64) {
	rows, total := s.sysUserRepository.SelectAuthUsersByPage(query, dataScopeSQL)
	s.fillUsersDeptAndRoles(rows)
	return rows, total
}

// ExportData 导出数据表格
func (s SysUser) ExportData(rows []model.SysUser, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"用户编号", "登录账号", "用户昵称", "用户邮箱", "手机号码", "用户性别", "帐号状态", "部门编号", "部门名称", "部门负责人", "最后登录IP", "最后登录时间"},
	}
	// 读取用户性别字典数据
	dictSysUserSex := s.sysDictTypeService.FindDataByType("sys_user_sex")
	// 从第二行开始的数据
	for _, row := range rows {
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
		data = append(data, []string{
			fmt.Sprintf("%d", row.UserId),
			row.UserName,
			row.NickName,
			row.Email,
			row.Phone,
			sysUserSex,
			statusValue,
			fmt.Sprintf("%d", row.Dept.DeptId),
			row.Dept.DeptName,
			row.Dept.Leader,
			row.LoginIp,
			date.ParseNumberToDate(row.LoginTime).Format(time.RFC3339),
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
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
	mustItemArr := []string{"username", "nickname"}
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
			if row["gender"] == v.DataLabel {
				sysUserSex = v.DataValue
				break
			}
		}
		sysUserStatus := constants.STATUS_NO
		if row["status"] == "正常" {
			sysUserStatus = constants.STATUS_YES
		}
		var sysUserDeptId int64 = 100
		if row["dept_id"] != "" {
			sysUserDeptId = parse.Number(row["dept_id"])
		}

		// 验证是否存在这个用户
		newSysUser := s.FindByUserName(row["username"])
		newSysUser.Password = initPassword
		newSysUser.UserName = row["username"]
		newSysUser.NickName = row["nickname"]
		newSysUser.Phone = row["phone"]
		newSysUser.Email = row["email"]
		newSysUser.StatusFlag = sysUserStatus
		newSysUser.Sex = sysUserSex
		newSysUser.DeptId = sysUserDeptId

		// 行用户编号
		rowNo := row["user_id"]

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
