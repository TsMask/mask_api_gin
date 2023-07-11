package permission

import (
	"mask_api_gin/src/framework/constants/admin"
)

// Role 获取角色数据权限
func Role(userID string, isAdmin bool) []string {
	if isAdmin {
		return []string{admin.ROLE_KEY}
	}
	return []string{}
}

// Menu 获取菜单数据权限
func Menu(userID string, isAdmin bool) []string {
	if isAdmin {
		return []string{admin.PERMISSION}
	}
	return []string{}
}
