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
	sysRoleRepository:     repository.NewSysRole,
	sysDeptRepository:     repository.NewSysDept,
	sysUserRoleRepository: repository.NewSysUserRole,
	sysUserPostRepository: repository.NewSysUserPost,
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
	sysDictDataService    *SysDictData            // 字典数据服务
	sysConfigService      *SysConfig              // 参数配置服务
}

// FindByPage 分页查询列表数据
func (s SysUser) FindByPage(query map[string]any, dataScopeSQL string) ([]model.SysUser, int64) {
	rows, total := s.sysUserRepository.SelectByPage(query, dataScopeSQL)
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

// Find 查询列表数据
func (s SysUser) Find(sysUser model.SysUser) []model.SysUser {
	return s.sysUserRepository.Select(sysUser)
}

// FindById 通过ID查询信息
func (s SysUser) FindById(userId int64) model.SysUser {
	userInfo := model.SysUser{}
	if userId == 0 {
		return userInfo
	}
	users := s.sysUserRepository.SelectByIds([]int64{userId})
	if len(users) > 0 {
		userInfo = users[0]
		// 部门
		deptInfo := s.sysDeptRepository.SelectById(userInfo.DeptId)
		userInfo.Dept = &deptInfo
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
		if roleId == 0 || roleId == constSystem.ROLE_SYSTEM_ID {
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
	if userId == 0 || len(postIds) <= 0 {
		return 0
	}

	var arr []model.SysUserPost
	for _, postId := range postIds {
		if postId == 0 {
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
	if userinfo.UserName == userName {
		// 部门
		deptInfo := s.sysDeptRepository.SelectById(userinfo.DeptId)
		userinfo.Dept = &deptInfo
		// 角色
		roleArr := s.sysRoleRepository.SelectByUserId(userinfo.UserId)
		roleIds := make([]int64, 0)
		roles := make([]*model.SysRole, 0)
		for _, role := range roleArr {
			roles = append(roles, &role)
			roleIds = append(roleIds, role.RoleId)
		}
		userinfo.Roles = roles
		userinfo.RoleIds = roleIds
	}
	return userinfo
}

// FindAuthUsersPage 根据条件分页查询分配用户角色列表
func (s SysUser) FindAuthUsersPage(query map[string]any, dataScopeSQL string) ([]model.SysUser, int64) {
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
