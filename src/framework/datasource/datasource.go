package datasource

import (
	"fmt"
	"log"
	"mask_api_gin/src/framework/logger"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

// 数据库连接实例
var dbMap = make(map[string]*gorm.DB)
var Default = dbMap

type dialectInfo struct {
	dialector gorm.Dialector
	logging   bool
}

// 数据库连接
var dialects = make(map[string]dialectInfo)

// 载入数据库连接
func loadDialect() {
	// 读取数据源配置
	datasource := viper.GetStringMap("gorm.datasource")
	for key, value := range datasource {
		item := value.(map[string]interface{})
		// 数据库类型对应的数据库连接
		switch item["type"] {
		case "mysql":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				item["username"],
				item["password"],
				item["host"],
				item["port"],
				item["database"],
			)
			dialects[key] = dialectInfo{
				dialector: mysql.Open(dsn),
				logging:   item["logging"].(bool),
			}
		default:
			logger.Warnf("%s: %v\n Not Load DB Config Type", key, item)
		}
	}
}

// 载入连接日志配置
func loadLogger() gormLog.Interface {
	newLogger := gormLog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 将日志输出到控制台
		gormLog.Config{
			SlowThreshold:        time.Second,  // Slow SQL 阈值
			LogLevel:             gormLog.Info, // 日志级别 Silent不输出任何日志
			ParameterizedQueries: false,        // 参数化查询SQL 用实际值带入?的执行语句
			Colorful:             false,        // 彩色日志输出
		},
	)
	return newLogger
}

// 连接数据库实例
func Connect() {
	loadDialect()
	// 遍历进行连接数据库实例
	for key, info := range dialects {
		opts := &gorm.Config{}
		// 是否需要日志输出
		if info.logging {
			opts.Logger = loadLogger()
		}
		// 创建连接
		db, err := gorm.Open(info.dialector, opts)
		if err != nil {
			logger.Panicf("fatal error db connect: %s", err)
		}
		dbMap[key] = db
	}
}

// 关闭数据库实例
func Close() {
	for _, db := range dbMap {
		sqlDB, err := db.DB()
		if err != nil {
			continue
		}
		if err := sqlDB.Close(); err != nil {
			logger.Panicf("fatal error db close: %s", err)
		}
	}
}

// 获取默认数据源
func GetDefaultDB() *gorm.DB {
	source := viper.GetString("gorm.defaultDataSourceName")
	return dbMap[source]
}

// 获取数据源
func GetDB(source string) *gorm.DB {
	return dbMap[source]
}
