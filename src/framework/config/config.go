package config

import (
	"log"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// conf 配置上下文
var conf *viper.Viper

// InitConfig 初始化程序配置
func InitConfig() {
	conf = viper.New()
	initFlag()
	initViper()
}

// 指定参数绑定
func initFlag() {
	// --env prod
	pflag.String("env", "local", "指定运行环境配置，读取config配置文件 (local|prod)")
	pflag.Parse()
	_ = conf.BindPFlags(pflag.CommandLine)
}

// 配置文件读取
func initViper() {
	// 在当前工作目录中寻找配置
	conf.AddConfigPath("config")
	conf.AddConfigPath("src/config")
	// 如果配置文件名中没有扩展名，则需要设置Type
	conf.SetConfigType("yaml")
	// 配置文件的名称（无扩展名）
	conf.SetConfigName("config.default")
	// 读取默认配置文件
	if err := conf.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config default file: %s", err)
	}

	env := conf.GetString("env")
	if env != "local" && env != "prod" {
		log.Fatalf("fatal error config env for local or prod : %s", env)
	}
	log.Printf("当期服务环境运行配置 => %s \n", env)

	// 加载运行配置文件合并相同配置
	if env == "prod" {
		conf.SetConfigName("config.prod")
	} else {
		conf.SetConfigName("config.local")
	}
	if err := conf.MergeInConfig(); err != nil {
		log.Fatalf("fatal error config local file: %s", err)
	}

	// 记录程序开始运行的时间点
	conf.Set("runTime", time.Now())
}

// Env 获取运行服务环境
// local prod
func Env() string {
	return conf.GetString("env")
}

// RunTime 程序开始运行的时间
func RunTime() time.Time {
	return conf.GetTime("runTime")
}

// Get 获取配置信息
//
// Get("framework.name")
func Get(key string) any {
	return conf.Get(key)
}

// IsSysAdmin 用户是否为系统管理员
func IsSysAdmin(userID string) bool {
	if userID == "" {
		return false
	}
	// 从配置中获取系统管理员id列表
	admins := Get("user.sysAdminList").([]any)
	for _, s := range admins {
		if s.(string) == userID {
			return true
		}
	}
	return false
}
