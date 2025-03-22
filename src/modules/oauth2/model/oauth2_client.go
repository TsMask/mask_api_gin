package model

// Oauth2Client 用户授权第三方应用表 oauth2_client
type Oauth2Client struct {
	Id           int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // 应用ID
	ClientId     string `gorm:"column:client_id" json:"clientId"`             // 应用的唯一标识
	ClientSecret string `gorm:"column:client_secret" json:"clientSecret"`     // 应用的凭证秘钥
	Title        string `gorm:"column:title" json:"title"`                    // 应用名称
	IPWhite      string `gorm:"column:ip_white" json:"ipWhite"`               // IP白名单
	DelFlag      string `gorm:"column:del_flag" json:"delFlag"`               // 删除标记（0存在 1删除）
	LoginIp      string `gorm:"column:login_ip" json:"loginIp"`               // 最后登录IP
	LoginTime    int64  `gorm:"column:login_time" json:"loginTime"`           // 最后登录时间
	CreateBy     string `gorm:"column:create_by" json:"createBy"`             // 创建者
	CreateTime   int64  `gorm:"column:create_time" json:"createTime"`         // 创建时间
	UpdateBy     string `gorm:"column:update_by" json:"updateBy"`             // 更新者
	UpdateTime   int64  `gorm:"column:update_time" json:"updateTime"`         // 更新时间
	Remark       string `gorm:"column:remark" json:"remark"`                  // 备注
}

func (*Oauth2Client) TableName() string {
	return "oauth2_client"
}
