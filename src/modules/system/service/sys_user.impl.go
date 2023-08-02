package service

import (
	"errors"
	"fmt"
	"mask_api_gin/src/framework/constants/admin"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
)

// 实例化服务层 SysUserImpl 结构体
var NewSysUserImpl = &SysUserImpl{
	sysUserRepository:     repository.NewSysUserImpl,
	sysUserRoleRepository: repository.NewSysUserRoleImpl,
	sysUserPostRepository: repository.NewSysUserPostImpl,
	sysDictDataService:    NewSysDictDataImpl,
	sysConfigService:      NewSysConfigImpl,
}

// SysUserImpl 用户 服务层处理
type SysUserImpl struct {
	// 用户服务
	sysUserRepository repository.ISysUser
	// 用户与角色服务
	sysUserRoleRepository repository.ISysUserRole
	// 用户与岗位服务
	sysUserPostRepository repository.ISysUserPost
	// 字典数据服务
	sysDictDataService ISysDictData
	// 参数配置服务
	sysConfigService ISysConfig
}

// SelectUserPage 根据条件分页查询用户列表
func (r *SysUserImpl) SelectUserPage(query map[string]any, dataScopeSQL string) map[string]any {
	return r.sysUserRepository.SelectUserPage(query, dataScopeSQL)
}

// SelectUserList 根据条件查询用户列表
func (r *SysUserImpl) SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	return []model.SysUser{}
}

// SelectAllocatedPage 根据条件分页查询分配用户角色列表
func (r *SysUserImpl) SelectAllocatedPage(query map[string]any, dataScopeSQL string) map[string]any {
	return r.sysUserRepository.SelectAllocatedPage(query, dataScopeSQL)
}

// SelectUserByUserName 通过用户名查询用户
func (r *SysUserImpl) SelectUserByUserName(userName string) model.SysUser {
	return r.sysUserRepository.SelectUserByUserName(userName)
}

// SelectUserById 通过用户ID查询用户
func (r *SysUserImpl) SelectUserById(userId string) model.SysUser {
	if userId == "" {
		return model.SysUser{}
	}
	users := r.sysUserRepository.SelectUserByIds([]string{userId})
	if len(users) > 0 {
		return users[0]
	}
	return model.SysUser{}
}

// InsertUser 新增用户信息
func (r *SysUserImpl) InsertUser(sysUser model.SysUser) string {
	// 新增用户信息
	insertId := r.sysUserRepository.InsertUser(sysUser)
	if insertId != "" {
		// 新增用户角色信息
		r.insertUserRole(insertId, sysUser.RoleIDs)
		// 新增用户岗位信息
		r.insertUserPost(insertId, sysUser.PostIDs)
	}
	return insertId
}

// insertUserRole 新增用户角色信息
func (r *SysUserImpl) insertUserRole(userId string, roleIds []string) int64 {
	if userId == "" || len(roleIds) <= 0 {
		return 0
	}

	sysUserRoles := []model.SysUserRole{}
	for _, roleId := range roleIds {
		// 管理员角色禁止操作，只能通过配置指定用户ID分配
		if roleId == "" || roleId == admin.ROLE_ID {
			continue
		}
		sysUserRoles = append(sysUserRoles, model.NewSysUserRole(userId, roleId))
	}

	return r.sysUserRoleRepository.BatchUserRole(sysUserRoles)
}

// insertUserPost 新增用户岗位信息
func (r *SysUserImpl) insertUserPost(userId string, postIds []string) int64 {
	if userId == "" || len(postIds) <= 0 {
		return 0
	}

	sysUserPosts := []model.SysUserPost{}
	for _, postId := range postIds {
		if postId == "" {
			continue
		}
		sysUserPosts = append(sysUserPosts, model.NewSysUserPost(userId, postId))
	}

	return r.sysUserPostRepository.BatchUserPost(sysUserPosts)
}

// UpdateUser 修改用户信息
func (r *SysUserImpl) UpdateUser(sysUser model.SysUser) int64 {
	return r.sysUserRepository.UpdateUser(sysUser)
}

// UpdateUserAndRolePost 修改用户信息同时更新角色和岗位
func (r *SysUserImpl) UpdateUserAndRolePost(sysUser model.SysUser) int64 {
	// 删除用户与角色关联
	r.sysUserRoleRepository.DeleteUserRole([]string{sysUser.UserID})
	// 新增用户角色信息
	r.insertUserRole(sysUser.UserID, sysUser.RoleIDs)
	// 删除用户与岗位关联
	r.sysUserPostRepository.DeleteUserPost([]string{sysUser.UserID})
	// 新增用户岗位信息
	r.insertUserPost(sysUser.UserID, sysUser.PostIDs)
	return r.sysUserRepository.UpdateUser(sysUser)
}

// DeleteUserByIds 批量删除用户信息
func (r *SysUserImpl) DeleteUserByIds(userIds []string) (int64, error) {
	// 检查是否存在
	users := r.sysUserRepository.SelectUserByIds(userIds)
	if len(users) <= 0 {
		return 0, errors.New("没有权限访问用户数据！")
	}
	if len(users) == len(userIds) {
		// 删除用户与角色关联
		r.sysUserRoleRepository.DeleteUserRole(userIds)
		// 删除用户与岗位关联
		r.sysUserPostRepository.DeleteUserPost(userIds)
		// ... 注意其他userId进行关联的表
		// 删除用户
		rows := r.sysUserRepository.DeleteUserByIds(userIds)
		return rows, nil
	}
	return 0, errors.New("删除用户信息失败！")
}

// CheckUniqueUserName 校验用户名称是否唯一
func (r *SysUserImpl) CheckUniqueUserName(userName, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		UserName: userName,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// CheckUniquePhone 校验手机号码是否唯一
func (r *SysUserImpl) CheckUniquePhone(phonenumber, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		PhoneNumber: phonenumber,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueEmail 校验email是否唯一
func (r *SysUserImpl) CheckUniqueEmail(email, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		Email: email,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// ImportUser 导入用户数据
func (r *SysUserImpl) ImportUser(rows []map[string]string, isUpdateSupport bool, operName string) string {
	// 读取默认初始密码
	initPassword := r.sysConfigService.SelectConfigValueByKey("sys.user.initPassword")
	// 读取用户性别字典数据
	dictSysUserSex := r.sysDictDataService.SelectDictDataByType("sys_user_sex")

	// 导入记录
	successNum := 0
	failureNum := 0
	successMsgArr := []string{}
	failureMsgArr := []string{}
	mustItemArr := []string{"C", "D"}
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
			failureMsgArr = append(failureMsgArr, fmt.Sprintf("表格中必填列表项，%s}", mustItemArrStr))
			continue
		}

		// 用户性别转值
		sysUserSex := "0"
		for _, v := range dictSysUserSex {
			if row["G"] == v.DictLabel {
				sysUserSex = v.DictValue
				break
			}
		}
		sysUserStatus := common.STATUS_NO
		if row["H"] == "正常" {
			sysUserStatus = common.STATUS_YES
		}

		// 构建用户实体信息
		newSysUser := model.SysUser{
			UserType:    "sys",
			Password:    initPassword,
			DeptID:      row["B"],
			UserName:    row["C"],
			NickName:    row["D"],
			PhoneNumber: row["F"],
			Email:       row["E"],
			Status:      sysUserStatus,
			Sex:         sysUserSex,
		}

		// 检查手机号码格式并判断是否唯一
		if newSysUser.PhoneNumber != "" {
			if regular.ValidMobile(newSysUser.PhoneNumber) {
				uniquePhone := r.CheckUniquePhone(newSysUser.PhoneNumber, "")
				if !uniquePhone {
					msg := fmt.Sprintf("序号：%s 手机号码 %s 已存在", row["A"], row["F"])
					failureNum++
					failureMsgArr = append(failureMsgArr, msg)
					continue
				}
			} else {
				msg := fmt.Sprintf("序号：%s 手机号码 %s 格式错误", row["A"], row["F"])
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
				continue
			}
		}

		// 检查邮箱格式并判断是否唯一
		if newSysUser.Email != "" {
			if regular.ValidEmail(newSysUser.Email) {
				uniqueEmail := r.CheckUniqueEmail(newSysUser.Email, "")
				if !uniqueEmail {
					msg := fmt.Sprintf("序号：%s 用户邮箱 %s 已存在", row["A"], row["E"])
					failureNum++
					failureMsgArr = append(failureMsgArr, msg)
					continue
				}
			} else {
				msg := fmt.Sprintf("序号：%s 用户邮箱 %s 格式错误", row["A"], row["E"])
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
				continue
			}
		}

		// 验证是否存在这个用户
		userInfo := r.sysUserRepository.SelectUserByUserName(newSysUser.UserName)
		if userInfo.UserName != newSysUser.UserName {
			newSysUser.CreateBy = operName
			insertId := r.InsertUser(newSysUser)
			if insertId != "" {
				msg := fmt.Sprintf("序号：%s 登录名称 %s 导入成功", row["A"], row["C"])
				successNum++
				successMsgArr = append(successMsgArr, msg)
			} else {
				msg := fmt.Sprintf("序号：%s 登录名称 %s 导入失败", row["A"], row["E"])
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
			}
			continue
		}

		// 如果用户已存在 同时 是否更新支持
		if userInfo.UserName == newSysUser.UserName && isUpdateSupport {
			newSysUser.UserID = userInfo.UserID
			newSysUser.UpdateBy = operName
			rows := r.UpdateUser(newSysUser)
			if rows > 0 {
				msg := fmt.Sprintf("序号：%s 登录名称 %s 更新成功", row["A"], row["C"])
				successNum++
				successMsgArr = append(successMsgArr, msg)
			} else {
				msg := fmt.Sprintf("序号：%s 登录名称 %s 更新失败", row["A"], row["E"])
				failureNum++
				failureMsgArr = append(failureMsgArr, msg)
			}
			continue
		}
	}

	if failureNum > 0 {
		failureMsgArr = append([]string{fmt.Sprintf("很抱歉，导入失败！共 %d 条数据格式不正确，错误如下：", failureNum)}, failureMsgArr...)
		return strings.Join(failureMsgArr, "<br/>")
	}

	successMsgArr = append([]string{fmt.Sprintf("恭喜您，数据已全部导入成功！共 %d 条，数据如下：", successNum)}, successMsgArr...)
	return strings.Join(successMsgArr, "<br/>")
}
