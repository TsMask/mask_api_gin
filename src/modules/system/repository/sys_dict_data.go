package repository

import "mask_api_gin/src/modules/system/model"

// ISysDictDataRepository 字典类型数据表 数据层接口
type ISysDictDataRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysDictData model.SysDictData) []model.SysDictData

	// SelectByCodes 通过Code查询信息
	SelectByCodes(dictCodes []string) []model.SysDictData

	// Insert 新增信息
	Insert(sysDictData model.SysDictData) string

	// Update 修改信息
	Update(sysDictData model.SysDictData) int64

	// DeleteByCodes 批量删除信息
	DeleteByCodes(dictCodes []string) int64

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysDictData model.SysDictData) string

	// ExistDataByDictType 存在数据数量
	ExistDataByDictType(dictType string) int64

	// UpdateDataByDictType 更新一组字典类型
	UpdateDataByDictType(oldDictType string, newDictType string) int64
}
