package service

import (
	"testing"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/modules/system/model"
)

func TestBuildRouteMenusUsesFrontendRootLayout(t *testing.T) {
	menus := []model.SysMenu{
		{
			ID:           1,
			ParentId:     0,
			MenuName:     "系统管理",
			MenuPath:     "system",
			MenuType:     constants.MENU_TYPE_DIR,
			MenuVisible:  constants.STATUS_YES,
			MenuExternal: constants.STATUS_NO,
			Children: []model.SysMenu{
				{
					ID:            100,
					ParentId:      1,
					MenuName:      "用户管理",
					MenuPath:      "user",
					MenuComponent: "system/user/UserIndex",
					MenuType:      constants.MENU_TYPE_MENU,
					MenuCache:     constants.STATUS_YES,
					MenuVisible:   constants.STATUS_YES,
					MenuExternal:  constants.STATUS_NO,
				},
			},
		},
	}

	routers := NewSysMenu.BuildRouteMenus(menus, "")
	if len(routers) != 1 {
		t.Fatalf("BuildRouteMenus routers length = %d, want 1", len(routers))
	}

	root := routers[0]
	if root.Component != constants.MENU_COMPONENT_LAYOUT_BLANK {
		t.Errorf("root component = %q, want %q", root.Component, constants.MENU_COMPONENT_LAYOUT_BLANK)
	}
	if root.Redirect != "/system/user" {
		t.Errorf("root redirect = %q, want /system/user", root.Redirect)
	}
	if len(root.Children) != 1 {
		t.Fatalf("root children length = %d, want 1", len(root.Children))
	}

	page := root.Children[0]
	if page.Name != "UserIndex" {
		t.Errorf("page name = %q, want UserIndex", page.Name)
	}
	if page.Component != "system/user/UserIndex" {
		t.Errorf("page component = %q, want system/user/UserIndex", page.Component)
	}
	if !page.Meta.Cache {
		t.Error("page cache = false, want true")
	}
}
