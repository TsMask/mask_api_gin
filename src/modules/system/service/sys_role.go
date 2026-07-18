package service

import (
	"fmt"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/framework/utils/date"
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
		if v.ID == constants.SYS_ROLE_SYSTEM_ID {
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
	if sysRole.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysRole.UpdateBy = sysRole.CreateBy
		sysRole.UpdateTime = ms
		sysRole.CreateTime = ms
	}
	var insertId int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		insertId = r.sysRoleRepository.Insert(tx, sysRole)
		if insertId <= 0 {
			return fmt.Errorf("insert role failed")
		}
		if len(sysRole.MenuIds) > 0 {
			r.insertRoleMenu(tx, insertId, sysRole.MenuIds)
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return insertId
}

// insertRoleMenu 新增角色菜单信息
func (r SysRole) insertRoleMenu(d *gorm.DB, roleId int64, menuIds []int64) int64 {
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
	return r.sysRoleMenuRepository.BatchInsert(d, sysRoleMenus)
}

// Update 修改信息
func (r SysRole) Update(sysRole model.SysRole) int64 {
	if sysRole.UpdateBy != "" {
		sysRole.UpdateTime = time.Now().UnixMilli()
	}
	var rows int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		rows = r.sysRoleRepository.Update(tx, sysRole)
		if rows > 0 {
			r.sysRoleMenuRepository.DeleteByRoleIds(tx, []int64{sysRole.ID})
			r.insertRoleMenu(tx, sysRole.ID, sysRole.MenuIds)
		}
		return nil
	})
	if err != nil {
		return 0
	}
	// 角色菜单权限变更后清理持有该角色用户的 Token 缓存
	userIds := r.sysUserRoleRepository.SelectUserIdsByRoleIds([]int64{sysRole.ID})
	token.RemoveByUserIds(userIds)
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
		// 检查分配用户
		if useCount := r.sysUserRoleRepository.ExistUserByRoleId(role.ID); useCount > 0 {
			return 0, fmt.Errorf("【%s】已分配给用户,不能删除", role.RoleName)
		}
	}
	if len(roles) != len(roleIds) {
		return 0, fmt.Errorf("删除角色信息失败！")
	}
	var rows int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		r.sysRoleMenuRepository.DeleteByRoleIds(tx, roleIds)
		r.sysRoleDeptRepository.DeleteByRoleIds(tx, roleIds)
		rows = r.sysRoleRepository.DeleteByIds(tx, roleIds)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return rows, nil
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
	if sysRole.UpdateBy != "" {
		sysRole.UpdateTime = time.Now().UnixMilli()
	}
	var rows int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		rows = r.sysRoleRepository.Update(tx, sysRole)
		if rows > 0 {
			r.sysRoleDeptRepository.DeleteByRoleIds(tx, []int64{sysRole.ID})
			if sysRole.RoleDataScope == constants.ROLE_SCOPE_CUSTOM && len(sysRole.DeptIds) > 0 {
				arr := make([]model.SysRoleDept, 0)
				for _, deptId := range sysRole.DeptIds {
					if deptId <= 0 {
						continue
					}
					arr = append(arr, model.SysRoleDept{
						RoleId: sysRole.ID, DeptId: deptId,
					})
				}
				r.sysRoleDeptRepository.BatchInsert(tx, arr)
			}
		}
		return nil
	})
	if err != nil {
		return 0
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
	rows := r.sysUserRoleRepository.BatchInsert(nil, sysUserRoles)
	if rows > 0 {
		token.RemoveByUserIds(userIds)
	}
	return rows
}

// DeleteAuthUsers 批量取消授权用户角色
func (r SysRole) DeleteAuthUsers(roleId int64, userIds []int64) int64 {
	rows := r.sysUserRoleRepository.DeleteByRoleId(nil, roleId, userIds)
	if rows > 0 {
		token.RemoveByUserIds(userIds)
	}
	return rows
}

// ExportData 导出数据表格
func (s SysRole) ExportData(rows []model.SysRole, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"角色序号", "角色名称", "角色键值", "角色排序", "数据范围", "角色状态"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		// 数据范围
		dataScope := "空"
		if v, ok := constants.ROLE_SCOPE_DATA[row.RoleDataScope]; ok {
			dataScope = v
		}
		// 角色状态
		statusValue := "停用"
		if row.RoleStatus == constants.STATUS_YES {
			statusValue = "正常"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.RoleName,
			row.RoleKey,
			fmt.Sprintf("%d", row.RoleSort),
			dataScope,
			statusValue,
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
