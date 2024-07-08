package repository

import "mask_api_gin/src/modules/system/model"

// ISysDictTypeRepository 字典类型表 数据层接口
type ISysDictTypeRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysDictType model.SysDictType) []model.SysDictType

	// SelectByIds 通过ID查询信息
	SelectByIds(dictIds []string) []model.SysDictType

	// Insert 新增信息
	Insert(sysDictType model.SysDictType) string

	// Update 修改信息
	Update(sysDictType model.SysDictType) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(dictIds []string) int64

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysDictType model.SysDictType) string

	// SelectByDictType 通过字典类型查询信息
	SelectByDictType(dictType string) model.SysDictType
}
