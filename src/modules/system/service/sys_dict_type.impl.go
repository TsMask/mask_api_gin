package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// SysDictTypeImpl 字典类型 数据层处理
var SysDictTypeImpl = &sysDictTypeImpl{
	sysUserRepository: repository.SysDictTypeImpl,
}

type sysDictTypeImpl struct {
	// 字典类型服务
	sysUserRepository repository.ISysDictType
}

// SelectDictTypePage 根据条件分页查询字典类型
func (r *sysDictTypeImpl) SelectDictTypePage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectDictTypeList 根据条件查询字典类型
func (r *sysDictTypeImpl) SelectDictTypeList(sysDictType model.SysDictType) []model.SysDictType {
	return []model.SysDictType{}
}

// SelectDictTypeByID 根据字典类型ID查询信息
func (r *sysDictTypeImpl) SelectDictTypeByID(dictID string) model.SysDictType {
	return model.SysDictType{}
}

// SelectDictTypeByType 根据字典类型查询信息
func (r *sysDictTypeImpl) SelectDictTypeByType(dictType string) model.SysDictType {
	return model.SysDictType{}
}

// CheckUniqueDictName 校验字典名称是否唯一
func (r *sysDictTypeImpl) CheckUniqueDictName(dictName string) string {
	return ""
}

// CheckUniqueDictType 校验字典类型是否唯一
func (r *sysDictTypeImpl) CheckUniqueDictType(dictType string) string {
	return ""
}

// InsertDictType 新增字典类型信息
func (r *sysDictTypeImpl) InsertDictType(sysDictType model.SysDictType) string {
	return ""
}

// UpdateDictType 修改字典类型信息
func (r *sysDictTypeImpl) UpdateDictType(sysDictType model.SysDictType) int {
	return 0
}

// DeleteDictTypeByID 批量删除字典类型信息
func (r *sysDictTypeImpl) DeleteDictTypeByID(dictIDs []string) int {
	return 0
}
