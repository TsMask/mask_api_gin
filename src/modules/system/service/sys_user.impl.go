package service

import (
	"errors"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/date"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
)

// SysUserImpl 用户 数据层处理
var SysUserImpl = &sysUserImpl{
	sysUserRepository:     repository.SysUserImpl,
	sysUserRoleRepository: repository.SysUserRoleImpl,
	sysUserPostRepository: repository.SysUserPostImpl,
}

type sysUserImpl struct {
	// 用户服务
	sysUserRepository repository.ISysUser
	// 用户与角色服务
	sysUserRoleRepository repository.ISysUserRole
	// 用户与岗位服务
	sysUserPostRepository repository.ISysUserPost
}

// SelectUserPage 根据条件分页查询用户列表
func (r *sysUserImpl) SelectUserPage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return r.sysUserRepository.SelectUserPage(query, dataScopeSQL)
}

// SelectUserList 根据条件查询用户列表
func (r *sysUserImpl) SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	return []model.SysUser{}
}

// SelectAllocatedPage 根据条件分页查询分配用户角色列表
func (r *sysUserImpl) SelectAllocatedPage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectUserByUserName 通过用户名查询用户
func (r *sysUserImpl) SelectUserByUserName(userName string) model.SysUser {
	return r.sysUserRepository.SelectUserByUserName(userName)
}

// SelectUserById 通过用户ID查询用户
func (r *sysUserImpl) SelectUserById(userId string) model.SysUser {
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
func (r *sysUserImpl) InsertUser(sysUser model.SysUser) string {
	// 参数拼接
	paramMap := make(map[string]interface{})
	if sysUser.UserID != "" {
		paramMap["user_id"] = sysUser.UserID
	}
	if sysUser.DeptID != "" {
		paramMap["dept_id"] = sysUser.DeptID
	}
	if sysUser.UserName != "" {
		paramMap["user_name"] = sysUser.UserName
	}
	if sysUser.NickName != "" {
		paramMap["nick_name"] = sysUser.NickName
	}
	if sysUser.UserType != "" {
		paramMap["user_type"] = sysUser.UserType
	}
	if sysUser.Avatar != "" {
		paramMap["avatar"] = sysUser.Avatar
	}
	if sysUser.Email != "" {
		paramMap["email"] = sysUser.Email
	}
	if sysUser.PhoneNumber != "" {
		paramMap["phonenumber"] = sysUser.PhoneNumber
	}
	if sysUser.Sex != "" {
		paramMap["sex"] = sysUser.Sex
	}
	if sysUser.Password != "" {
		password := crypto.BcryptHash(sysUser.Password)
		paramMap["password"] = password
	}
	if sysUser.Status != "" {
		paramMap["status"] = sysUser.Status
	}
	if sysUser.Remark != "" {
		paramMap["remark"] = sysUser.Remark
	}
	if sysUser.CreateBy != "" {
		paramMap["create_by"] = sysUser.CreateBy
		paramMap["create_time"] = date.NowTimestamp()
	}

	// 构建执行语句
	keys, placeholder, values := repoUtils.KeyPlaceholderValueByInsert(paramMap)
	sql := "insert into sys_user (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

	db := datasource.DefaultDB()
	// 开启事务
	tx := db.Begin()
	// 执行插入
	err := tx.Exec(sql, values...).Error
	if err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return err.Error()
	}
	// 获取生成的自增 ID
	var insertedID string
	err = tx.Raw("select last_insert_id()").Row().Scan(&insertedID)
	if err != nil {
		logger.Errorf("insert last id : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 提交事务
	tx.Commit()
	return insertedID
}

// UpdateUser 修改用户信息
func (r *sysUserImpl) UpdateUser(sysUser model.SysUser) int64 {
	// 参数拼接
	paramMap := make(map[string]interface{})
	if sysUser.DeptID != "" {
		paramMap["dept_id"] = sysUser.DeptID
	}
	if sysUser.UserName != "" {
		paramMap["user_name"] = sysUser.UserName
	}
	if sysUser.NickName != "" {
		paramMap["nick_name"] = sysUser.NickName
	}
	if sysUser.UserType != "" {
		paramMap["user_type"] = sysUser.UserType
	}
	if sysUser.Avatar != "" {
		paramMap["avatar"] = sysUser.Avatar
	}
	if sysUser.Email != "" {
		paramMap["email"] = sysUser.Email
	}
	if sysUser.PhoneNumber != "" {
		paramMap["phonenumber"] = sysUser.PhoneNumber
	}
	if sysUser.Sex != "" {
		paramMap["sex"] = sysUser.Sex
	}
	if sysUser.Password != "" {
		password := crypto.BcryptHash(sysUser.Password)
		paramMap["password"] = password
	}
	if sysUser.Status != "" {
		paramMap["status"] = sysUser.Status
	}
	if sysUser.Remark != "" {
		paramMap["remark"] = sysUser.Remark
	}
	if sysUser.UpdateBy != "" {
		paramMap["update_by"] = sysUser.UpdateBy
		paramMap["update_time"] = date.NowTimestamp()
	}
	if sysUser.LoginIP != "" {
		paramMap["login_ip"] = sysUser.LoginIP
	}
	if sysUser.LoginDate > 0 {
		paramMap["login_date"] = sysUser.LoginDate
	}

	// 构建执行语句
	keys, values := repoUtils.KeyValueByUpdate(paramMap)
	sql := "update sys_user set " + strings.Join(keys, ",") + " where user_id = ?"

	// 执行更新
	values = append(values, sysUser.UserID)
	num, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return num
}

// DeleteUserByIds 批量删除用户信息
func (r *sysUserImpl) DeleteUserByIds(userIds []string) (int64, error) {
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
func (r *sysUserImpl) CheckUniqueUserName(userName, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		UserName: userName,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// CheckUniquePhone 校验手机号码是否唯一
func (r *sysUserImpl) CheckUniquePhone(phonenumber, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		PhoneNumber: phonenumber,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueEmail 校验email是否唯一
func (r *sysUserImpl) CheckUniqueEmail(email, userId string) bool {
	uniqueId := r.sysUserRepository.CheckUniqueUser(model.SysUser{
		Email: email,
	})
	if uniqueId == userId {
		return true
	}
	return uniqueId == ""
}
