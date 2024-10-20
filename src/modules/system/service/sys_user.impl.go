package service

import (
	"errors"
	"mask_api_gin/src/framework/constants/admin"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysUserImpl 实例化服务层 SysUserImpl 结构体
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
func (r *SysUserImpl) SelectUserPage(queryMap map[string]any, dataScopeSQL string) map[string]any {
	return r.sysUserRepository.SelectUserPage(queryMap, dataScopeSQL)
}

// SelectUserList 根据条件查询用户列表
func (r *SysUserImpl) SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	return r.sysUserRepository.SelectUserList(sysUser, dataScopeSQL)
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

	var sysUserRoles []model.SysUserRole
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

	var sysUserPosts []model.SysUserPost
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
func (r *SysUserImpl) CheckUniquePhone(phone, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		PhoneNumber: phone,
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
