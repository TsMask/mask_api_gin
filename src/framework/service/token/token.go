package token

import (
	"mask_api_gin/src/framework/model"
	"mask_api_gin/src/framework/service/permission"
	systemSystem "mask_api_gin/src/modules/system/model"
)

// Remove 清除用户登录令牌
func Remove(token string) bool {
	_, err := getLoginUser(token)
	return err == nil
}

// CreateLoginUser 创建登录用户信息对象
func CreateLoginUser(user systemSystem.SysUser, isAdmin bool) model.LoginUser {
	loginUser := model.LoginUser{
		UserID: user.UserID,
		DeptID: user.DeptID,
		User:   user,
	}
	// 用户权限组标识
	loginUser.Permissions = permission.Menu(user.UserID, isAdmin)
	return model.LoginUser{}
}

// 获取用户身份信息
func getLoginUser(token string) (model.LoginUser, error) {
	return model.LoginUser{}, nil
}
