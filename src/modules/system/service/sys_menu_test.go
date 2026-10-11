package service

import (
	"testing"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/modules/system/model"
)

func TestBuildAuthMenusUsesFrontendMenuContract(t *testing.T) {
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

	authMenus := NewSysMenu.BuildAuthMenus(menus)
	if len(authMenus) != 1 {
		t.Fatalf("BuildAuthMenus menus length = %d, want 1", len(authMenus))
	}

	root := authMenus[0]
	if root.Kind != "directory" {
		t.Errorf("root kind = %q, want directory", root.Kind)
	}
	if root.Component != "" {
		t.Errorf("root component = %q, want empty", root.Component)
	}
	if root.Path != "system" {
		t.Errorf("root path = %q, want system", root.Path)
	}
	if len(root.Children) != 1 {
		t.Fatalf("root children length = %d, want 1", len(root.Children))
	}

	page := root.Children[0]
	if page.Kind != "page" {
		t.Errorf("page kind = %q, want page", page.Kind)
	}
	if page.Component != "system/user/UserIndex" {
		t.Errorf("page component = %q, want system/user/UserIndex", page.Component)
	}
	if !page.Cache {
		t.Error("page cache = false, want true")
	}
}
