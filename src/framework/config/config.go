package config

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// 程序名
	Name string = "-"
	// 程序版本
	Version string = "-"
	// 编译时间
	BuildTime string = "-"
	// Go版本
	GoVer string = "-"
)

// conf 配置上下文
var conf *viper.Viper

// InitConfig 初始化程序配置
func InitConfig(configDir *embed.FS) {
	conf = viper.New()
	initFlag()
	initViper(configDir)
}

// 指定参数绑定
func initFlag() {
	// --env prod
	pflag.String("env", "prod", "指定运行环境配置,读取config配置文件 (local|prod)")
	// --config /etc/config.yaml
	// -c ./config.yaml
	pflag.StringP("config", "c", "", "指定配置文件覆盖默认配置")
	// --version
	// -V
	pVersion := pflag.BoolP("version", "V", false, "程序版本信息")
	// --help
	pHelp := pflag.Bool("help", false, "查看帮助信息")
	pflag.Parse()

	// 参数固定输出
	if *pVersion {
		fmt.Printf("Name:%s\nVersion:%s\nBuildTime:%s\nBuildGoVer:%s\n\n", Name, Version, BuildTime, GoVer)
		os.Exit(1)
	}
	if *pHelp {
		pflag.Usage()
		os.Exit(1)
	}

	_ = conf.BindPFlags(pflag.CommandLine)
}

// 配置文件读取
func initViper(configDir *embed.FS) {
	// 如果配置文件名中没有扩展名，则需要设置Type
	conf.SetConfigType("yaml")
	// 读取默认配置文件
	configDefaultByte, err := configDir.ReadFile("src/config/config.default.yaml")
	if err != nil {
		log.Fatalf("config default file read error: %s", err)
		return
	}
	if err = conf.ReadConfig(bytes.NewReader(configDefaultByte)); err != nil {
		log.Fatalf("config default file read error: %s", err)
		return
	}

	// 当期服务环境运行配置 => local
	env := conf.GetString("env")
	log.Printf("current service environment operation configuration => %s \n", env)

	// 加载运行配置文件合并相同配置
	envConfigPath := fmt.Sprintf("src/config/config.%s.yaml", env)
	configEnvByte, err := configDir.ReadFile(envConfigPath)
	if err != nil {
		log.Fatalf("config env %s file read error: %s", env, err)
		return
	}
	if err = conf.MergeConfig(bytes.NewReader(configEnvByte)); err != nil {
		log.Fatalf("config env %s file read error: %s", env, err)
		return
	}

	// 外部文件配置
	if externalConfig := conf.GetString("config"); externalConfig != "" {
		readExternalConfig(externalConfig)
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

// GetAssetsDirFS 访问程序内全局资源访问
func GetAssetsDirFS() *embed.FS {
	return conf.Get("AssetsDir").(*embed.FS)
}

// SetAssetsDirFS 设置程序内全局资源访问
func SetAssetsDirFS(assetsDir *embed.FS) {
	conf.Set("AssetsDir", assetsDir)
}

// readExternalConfig 读取外部文件配置
func readExternalConfig(configPaht string) {
	f, err := os.Open(configPaht)
	if err != nil {
		log.Fatalf("config external file read error: %s", err)
		return
	}
	defer f.Close()

	if err = conf.MergeConfig(f); err != nil {
		log.Fatalf("config external file read error: %s", err)
		return
	}
}

// IsSystemUser 用户是否为系统管理员
func IsSystemUser(userId string) bool {
	if userId == "" {
		return false
	}
	// 从配置中获取系统管理员ID列表
	arr := Get("user.system").([]any)
	for _, v := range arr {
		if fmt.Sprint(v) == userId {
			return true
		}
	}
	return false
}
