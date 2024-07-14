package service

import "mask_api_gin/src/modules/system/model"

// ISysDictDataService 字典类型数据 服务层接口
type ISysDictDataService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// Find 查询数据
	Find(sysDictData model.SysDictData) []model.SysDictData

	// FindByCode 通过Code查询信息
	FindByCode(dictCode string) model.SysDictData

	// FindByType 根据字典类型查询信息
	FindByType(dictType string) []model.SysDictData

	// Insert 新增信息
	Insert(sysDictData model.SysDictData) string

	// Update 修改信息
	Update(sysDictData model.SysDictData) int64

	// DeleteByCodes 批量删除信息
	DeleteByCodes(dictCodes []string) (int64, error)

	// CheckUniqueTypeByLabel 检查同字典类型下字典标签是否唯一
	CheckUniqueTypeByLabel(dictType, dictLabel, dictCode string) bool

	// CheckUniqueTypeByValue 检查同字典类型下字典键值是否唯一
	CheckUniqueTypeByValue(dictType, dictValue, dictCode string) bool
}
