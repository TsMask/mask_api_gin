package config

import (
	"mask_api_gin/src/framework/logger"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 初始化程序配置
func InitConfig() {
	initFlag()
	initViper()
}

// 指定参数绑定
func initFlag() {
	// -env prod
	pflag.String("env", "local", "指定运行环境配置，读取config配置文件 local、prod")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

// 配置文件读取
func initViper() {
	env := viper.GetString("env")
	if env != "local" && env != "prod" {
		logger.Panicf("无效环境值 %s ，请指定local、prod", env)
	}
	logger.Warnf("当期服务环境运行配置 => %s", env)

	// 在当前工作目录中寻找配置
	viper.AddConfigPath("src/config")
	// 如果配置文件名中没有扩展名，则需要设置Type
	viper.SetConfigType("yaml")
	// 配置文件的名称（无扩展名）
	viper.SetConfigName("config.default")
	// 读取默认配置文件
	if err := viper.ReadInConfig(); err != nil {
		logger.Panicf("fatal error config default file: %s", err)
	}

	// 加载运行配置文件合并相同配置
	if env == "prod" {
		viper.SetConfigName("config.prod")
	} else {
		viper.SetConfigName("config.local")
	}
	if err := viper.MergeInConfig(); err != nil {
		logger.Panicf("fatal error config local file: %s", err)
	}

	// 记录程序开始运行的时间点
	viper.Set("runTime", time.Now())
}

// Env 获取运行服务环境
// local prod
func Env() string {
	return viper.GetString("env")
}

// RunTime 程序开始运行的时间
func RunTime() time.Time {
	return viper.GetTime("runTime")
}

// Get 获取配置信息
//
// Get("framework.name")
func Get(key string) any {
	return viper.Get(key)
}

// IsAdmin 用户是否为管理员
func IsAdmin(userID string) bool {
	if userID == "" {
		return false
	}
	// 从本地配置获取user信息
	admins := Get("user.adminList").([]any)
	for _, s := range admins {
		if s.(string) == userID {
			return true
		}
	}
	return false
}
