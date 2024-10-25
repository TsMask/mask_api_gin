package service

import (
	"fmt"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysUser 实例化服务层
var NewSysUser = &SysUser{
	sysUserRepository:     repository.NewSysUser,
	sysUserRoleRepository: repository.NewSysUserRole,
	sysUserPostRepository: repository.NewSysUserPost,
	sysDictDataService:    NewSysDictData,
	sysConfigService:      NewSysConfig,
}

// SysUser 用户 服务层处理
type SysUser struct {
	sysUserRepository     *repository.SysUser               // 用户服务
	sysUserRoleRepository repository.ISysUserRoleRepository // 用户与角色服务
	sysUserPostRepository repository.ISysUserPostRepository // 用户与岗位服务
	sysDictDataService    ISysDictDataService               // 字典数据服务
	sysConfigService      *SysConfig                        // 参数配置服务
}

// FindByPage 分页查询列表数据
func (s SysUser) FindByPage(queryMap map[string]any, dataScopeSQL string) ([]model.SysUser, int64) {
	return s.sysUserRepository.SelectByPage(queryMap, dataScopeSQL)
}

// Find 查询列表数据
func (s SysUser) Find(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	return s.sysUserRepository.Select(sysUser, dataScopeSQL)
}

// FindById 通过ID查询信息
func (s SysUser) FindById(userId string) model.SysUser {
	if userId == "" {
		return model.SysUser{}
	}
	users := s.sysUserRepository.SelectByIds([]string{userId})
	if len(users) > 0 {
		return users[0]
	}
	return model.SysUser{}
}

// Insert 新增信息
func (s SysUser) Insert(sysUser model.SysUser) string {
	// 新增用户信息
	insertId := s.sysUserRepository.Insert(sysUser)
	if insertId != "" {
		s.insertUserRole(insertId, sysUser.RoleIDs) // 新增用户角色信息
		s.insertUserPost(insertId, sysUser.PostIDs) // 新增用户岗位信息
	}
	return insertId
}

// insertUserRole 新增用户角色信息
func (s SysUser) insertUserRole(userId string, roleIds []string) int64 {
	if userId == "" || len(roleIds) <= 0 {
		return 0
	}

	var arr []model.SysUserRole
	for _, roleId := range roleIds {
		// 系统管理员角色禁止操作，只能通过配置指定用户ID分配
		if roleId == "" || roleId == constSystem.RoleId {
			continue
		}
		arr = append(arr, model.SysUserRole{
			UserID: userId, RoleID: roleId,
		})
	}

	return s.sysUserRoleRepository.BatchInsert(arr)
}

// insertUserPost 新增用户岗位信息
func (s SysUser) insertUserPost(userId string, postIds []string) int64 {
	if userId == "" || len(postIds) <= 0 {
		return 0
	}

	var arr []model.SysUserPost
	for _, postId := range postIds {
		if postId == "" {
			continue
		}
		arr = append(arr, model.SysUserPost{
			UserID: userId, PostID: postId,
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
	s.sysUserRoleRepository.DeleteByUserIds([]string{sysUser.UserID})
	// 新增用户角色信息
	s.insertUserRole(sysUser.UserID, sysUser.RoleIDs)
	// 删除用户与岗位关联
	s.sysUserPostRepository.DeleteByUserIds([]string{sysUser.UserID})
	// 新增用户岗位信息
	s.insertUserPost(sysUser.UserID, sysUser.PostIDs)
	return s.sysUserRepository.Update(sysUser)
}

// DeleteByIds 批量删除信息
func (s SysUser) DeleteByIds(userIds []string) (int64, error) {
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
func (s SysUser) CheckUniqueByUserName(userName, userId string) bool {
	uniqueId := s.sysUserRepository.CheckUnique(model.SysUser{
		UserName: userName,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueByPhone 检查手机号码是否唯一
func (s SysUser) CheckUniqueByPhone(phone, userId string) bool {
	uniqueId := s.sysUserRepository.CheckUnique(model.SysUser{
		Phone: phone,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueByEmail 检查Email是否唯一
func (s SysUser) CheckUniqueByEmail(email, userId string) bool {
	uniqueId := s.sysUserRepository.CheckUnique(model.SysUser{
		Email: email,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// FindByUserName 通过用户名查询用户信息
func (s SysUser) FindByUserName(userName string) model.SysUser {
	return s.sysUserRepository.SelectByUserName(userName)
}

// FindAllocatedPage 根据条件分页查询分配用户角色列表
func (s SysUser) FindAllocatedPage(query map[string]any, dataScopeSQL string) ([]model.SysUser, int64) {
	return s.sysUserRepository.SelectAllocatedByPage(query, dataScopeSQL)
}
