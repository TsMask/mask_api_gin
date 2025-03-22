package service

import (
	"fmt"

	"mask_api_gin/src/framework/utils/generate"
	"mask_api_gin/src/modules/oauth2/model"
	"mask_api_gin/src/modules/oauth2/repository"
)

// NewOauth2ClientService 实例化服务层
var NewOauth2ClientService = &Oauth2ClientService{
	oauth2ClientRepository: repository.NewOauth2Client,
}

// Oauth2ClientService 用户授权第三方应用信息 服务层处理
type Oauth2ClientService struct {
	oauth2ClientRepository *repository.Oauth2Client // 用户授权第三方应用表
}

// FindByPage 分页查询
func (s Oauth2ClientService) FindByPage(query map[string]string) ([]model.Oauth2Client, int64) {
	return s.oauth2ClientRepository.SelectByPage(query)
}

// Find 查询集合
func (s Oauth2ClientService) FindByClientId(clientId string) model.Oauth2Client {
	return s.oauth2ClientRepository.SelectByClientId(clientId)
}

// Insert 新增
func (s Oauth2ClientService) Insert(param model.Oauth2Client) int64 {
	param.ClientId = generate.Code(16)
	param.ClientSecret = generate.Code(32)
	return s.oauth2ClientRepository.Insert(param)
}

// Update 更新
func (s Oauth2ClientService) Update(param model.Oauth2Client) int64 {
	return s.oauth2ClientRepository.Update(param)
}

// DeleteByIds 批量删除
func (s Oauth2ClientService) DeleteByIds(ids []int64) (int64, error) {
	// 检查是否存在
	arr := s.oauth2ClientRepository.SelectByIds(ids)
	if len(arr) <= 0 {
		return 0, fmt.Errorf("没有权限访问用户授权第三方应用数据！")
	}
	if len(arr) == len(ids) {
		return s.oauth2ClientRepository.DeleteByIds(ids), nil
	}
	return 0, fmt.Errorf("删除用户授权第三方应用信息失败！")
}
