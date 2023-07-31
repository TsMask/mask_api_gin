package datasource

import (
	"fmt"
	"log"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"os"
	"regexp"
	"time"

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

// 载入数据库连接
func loadDialect() map[string]dialectInfo {
	dialects := make(map[string]dialectInfo, 0)

	// 读取数据源配置
	datasource := config.Get("gorm.datasource").(map[string]interface{})
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

	return dialects
}

// 载入连接日志配置
func loadLogger() gormLog.Interface {
	newLogger := gormLog.New(
		log.New(os.Stdout, "[GORM] ", log.LstdFlags), // 将日志输出到控制台
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
	// 遍历进行连接数据库实例
	for key, info := range loadDialect() {
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
		// 获取底层 SQL 数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			logger.Panicf("Failed to get underlying SQL database: %v", err)
		}
		// 测试数据库连接
		err = sqlDB.Ping()
		if err != nil {
			logger.Panicf("Failed to ping database: %v", err)
		}
		logger.Infof("Database %s connection is successful.", key)
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
func DefaultDB() *gorm.DB {
	source := config.Get("gorm.defaultDataSourceName").(string)
	return dbMap[source]
}

// 获取数据源
func DB(source string) *gorm.DB {
	return dbMap[source]
}

// RawDB 原生查询语句
func RawDB(source string, sql string, parameters []interface{}) ([]map[string]interface{}, error) {
	// 数据源
	db := DefaultDB()
	if source != "" {
		db = DB(source)
	}
	// 使用正则表达式替换连续的空白字符为单个空格
	fmtSql := regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")

	// logger.Infof("sql=> %v", fmtSql)
	// logger.Infof("parameters=> %v", parameters)

	// 查询结果
	var rows []map[string]interface{}
	res := db.Raw(fmtSql, parameters...).Scan(&rows)
	if res.Error != nil {
		return nil, res.Error
	}
	return rows, nil
}

// ExecDB 原生执行语句
func ExecDB(source string, sql string, parameters []interface{}) (int64, error) {
	// 数据源
	db := DefaultDB()
	if source != "" {
		db = DB(source)
	}
	// 使用正则表达式替换连续的空白字符为单个空格
	fmtSql := regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	// 执行结果
	res := db.Exec(fmtSql, parameters...)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
