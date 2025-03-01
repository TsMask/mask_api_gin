package service

import (
	"fmt"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysRole 实例化服务层
var NewSysRole = &SysRole{
	sysRoleRepository:     repository.NewSysRole,
	sysUserRoleRepository: repository.NewSysUserRole,
	sysRoleDeptRepository: repository.NewSysRoleDept,
	sysRoleMenuRepository: repository.NewSysRoleMenu,
}

// SysRole 角色 服务层处理
type SysRole struct {
	sysRoleRepository     *repository.SysRole     // 角色服务
	sysUserRoleRepository *repository.SysUserRole // 用户与角色关联服务
	sysRoleDeptRepository *repository.SysRoleDept // 角色与部门关联服务
	sysRoleMenuRepository *repository.SysRoleMenu // 角色与菜单关联服务
}

// FindByPage 分页查询列表数据
func (r SysRole) FindByPage(query map[string]string) ([]model.SysRole, int64) {
	return r.sysRoleRepository.SelectByPage(query)
}

// Find 查询列表数据
func (r SysRole) Find(sysRole model.SysRole) []model.SysRole {
	rows := r.sysRoleRepository.Select(sysRole)
	arr := make([]model.SysRole, 0)
	for _, v := range rows {
		if v.RoleId == constants.SYS_ROLE_SYSTEM_ID {
			continue
		}
		arr = append(arr, v)
	}
	return arr
}

// FindById 通过ID查询信息
func (r SysRole) FindById(roleId int64) model.SysRole {
	if roleId <= 0 {
		return model.SysRole{}
	}
	rows := r.sysRoleRepository.SelectByIds([]int64{roleId})
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysRole{}
}

// Insert 新增信息
func (r SysRole) Insert(sysRole model.SysRole) int64 {
	insertId := r.sysRoleRepository.Insert(sysRole)
	if insertId > 0 && len(sysRole.MenuIds) > 0 {
		r.insertRoleMenu(insertId, sysRole.MenuIds)
	}
	return insertId
}

// insertRoleMenu 新增角色菜单信息
func (r SysRole) insertRoleMenu(roleId int64, menuIds []int64) int64 {
	if roleId <= 0 || len(menuIds) <= 0 {
		return 0
	}
	sysRoleMenus := make([]model.SysRoleMenu, 0)
	for _, menuId := range menuIds {
		if menuId <= 0 {
			continue
		}
		sysRoleMenus = append(sysRoleMenus, model.SysRoleMenu{
			RoleId: roleId, MenuId: menuId,
		})
	}
	return r.sysRoleMenuRepository.BatchInsert(sysRoleMenus)
}

// Update 修改信息
func (r SysRole) Update(sysRole model.SysRole) int64 {
	rows := r.sysRoleRepository.Update(sysRole)
	if rows > 0 && len(sysRole.MenuIds) > 0 {
		// 删除角色与菜单关联
		r.sysRoleMenuRepository.DeleteByRoleIds([]int64{sysRole.RoleId})
		r.insertRoleMenu(sysRole.RoleId, sysRole.MenuIds)
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r SysRole) DeleteByIds(roleIds []int64) (int64, error) {
	// 检查是否存在
	roles := r.sysRoleRepository.SelectByIds(roleIds)
	if len(roles) <= 0 {
		return 0, fmt.Errorf("没有权限访问角色数据！")
	}
	for _, role := range roles {
		// 检查是否为已删除
		if role.DelFlag == "1" {
			return 0, fmt.Errorf("ID:%d 角色信息已经删除！", role.RoleId)
		}
		// 检查分配用户
		if useCount := r.sysUserRoleRepository.ExistUserByRoleId(role.RoleId); useCount > 0 {
			return 0, fmt.Errorf("【%s】已分配给用户,不能删除", role.RoleName)
		}
	}
	if len(roles) == len(roleIds) {
		r.sysRoleMenuRepository.DeleteByRoleIds(roleIds) // 删除角色与菜单关联
		r.sysRoleDeptRepository.DeleteByRoleIds(roleIds) // 删除角色与部门关联
		return r.sysRoleRepository.DeleteByIds(roleIds), nil
	}
	return 0, fmt.Errorf("删除角色信息失败！")
}

// FindByUserId 根据用户ID获取角色选择框列表
func (r SysRole) FindByUserId(userId int64) []model.SysRole {
	return r.sysRoleRepository.SelectByUserId(userId)
}

// CheckUniqueByName 检查角色名称是否唯一
func (r SysRole) CheckUniqueByName(roleName string, roleId int64) bool {
	uniqueId := r.sysRoleRepository.CheckUnique(model.SysRole{
		RoleName: roleName,
	})
	if uniqueId == roleId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueByKey 检查角色权限是否唯一
func (r SysRole) CheckUniqueByKey(roleKey string, roleId int64) bool {
	uniqueId := r.sysRoleRepository.CheckUnique(model.SysRole{
		RoleKey: roleKey,
	})
	if uniqueId == roleId {
		return true
	}
	return uniqueId == 0
}

// UpdateAndDataScope 修改信息同时更新数据权限信息
func (r SysRole) UpdateAndDataScope(sysRole model.SysRole) int64 {
	// 修改角色信息
	rows := r.sysRoleRepository.Update(sysRole)
	if rows > 0 {
		// 删除角色与部门关联
		r.sysRoleDeptRepository.DeleteByRoleIds([]int64{sysRole.RoleId})
		// 新增角色和部门信息
		if sysRole.DataScope == constants.ROLE_SCOPE_CUSTOM && len(sysRole.DeptIds) > 0 {
			arr := make([]model.SysRoleDept, 0)
			for _, deptId := range sysRole.DeptIds {
				if deptId <= 0 {
					continue
				}
				arr = append(arr, model.SysRoleDept{
					RoleId: sysRole.RoleId, DeptId: deptId,
				})
			}
			r.sysRoleDeptRepository.BatchInsert(arr)
		}
	}
	return rows
}

// InsertAuthUsers 批量新增授权用户角色
func (r SysRole) InsertAuthUsers(roleId int64, userIds []int64) int64 {
	if roleId <= 0 || len(userIds) <= 0 {
		return 0
	}
	sysUserRoles := make([]model.SysUserRole, 0)
	for _, userId := range userIds {
		if userId <= 0 {
			continue
		}
		sysUserRoles = append(sysUserRoles, model.SysUserRole{
			UserId: userId, RoleId: roleId,
		})
	}
	return r.sysUserRoleRepository.BatchInsert(sysUserRoles)
}

// DeleteAuthUsers 批量取消授权用户角色
func (r SysRole) DeleteAuthUsers(roleId int64, userIds []int64) int64 {
	return r.sysUserRoleRepository.DeleteByRoleId(roleId, userIds)
}

// ExportData 导出数据表格
func (s SysRole) ExportData(rows []model.SysRole, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "角色序号",
		"B1": "角色名称",
		"C1": "角色键值",
		"D1": "角色排序",
		"E1": "数据范围",
		"F1": "角色状态",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := fmt.Sprintf("%d", i+2)
		// 数据范围
		dataScope := "空"
		if v, ok := constants.ROLE_SCOPE_DATA[row.DataScope]; ok {
			dataScope = v
		}
		// 角色状态
		statusValue := "停用"
		if row.StatusFlag == constants.STATUS_YES {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.RoleId,
			"B" + idx: row.RoleName,
			"C" + idx: row.RoleKey,
			"D" + idx: row.RoleSort,
			"E" + idx: dataScope,
			"F" + idx: statusValue,
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
