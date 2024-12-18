package db

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"

	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

// 数据库连接实例
var dbMap = make(map[string]*gorm.DB)

type dialectInfo struct {
	dialectic gorm.Dialector
	logging   bool
}

// 载入数据库连接
func loadDialect() map[string]dialectInfo {
	dialects := make(map[string]dialectInfo)

	// 读取数据源配置
	datasource := config.Get("gorm.datasource").(map[string]any)
	for key, value := range datasource {
		item := value.(map[string]any)
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
				dialectic: mysql.Open(dsn),
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

// Connect 连接数据库实例
func Connect() {
	// 遍历进行连接数据库实例
	for key, info := range loadDialect() {
		opts := &gorm.Config{}
		// 是否需要日志输出
		if info.logging {
			opts.Logger = loadLogger()
		}
		// 创建连接
		db, err := gorm.Open(info.dialectic, opts)
		if err != nil {
			logger.Fatalf("failed error db connect: %s", err)
		}
		// 获取底层 SQL 数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			logger.Fatalf("failed error underlying SQL database: %v", err)
		}
		// 测试数据库连接
		err = sqlDB.Ping()
		if err != nil {
			logger.Fatalf("failed error ping database: %v", err)
		}
		// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)
		logger.Infof("database %s connection is successful.", key)
		dbMap[key] = db
	}
}

// Close 关闭数据库实例
func Close() {
	for _, db := range dbMap {
		sqlDB, err := db.DB()
		if err != nil {
			continue
		}
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("fatal error db close: %s", err)
		}
	}
}

// DB 获取数据源
//
// source-数据源
func DB(source string) *gorm.DB {
	// 不指定时获取默认实例
	if source == "" {
		source = config.Get("gorm.defaultDataSourceName").(string)
	}
	return dbMap[source]
}

// RawDB 原生语句查询
//
// source-数据源
// sql-预编译的SQL语句
// parameters-预编译的SQL语句参数
func RawDB(source string, sql string, parameters []any) ([]map[string]any, error) {
	var rows []map[string]any
	// 数据源
	db := DB(source)
	if db == nil {
		return rows, fmt.Errorf("not database source")
	}
	// 使用正则表达式替换连续的空白字符为单个空格
	fmtSql := regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	// 查询结果
	res := db.Raw(fmtSql, parameters...).Scan(&rows)
	if res.Error != nil {
		return nil, res.Error
	}
	return rows, nil
}

// ExecDB 原生语句执行
//
// source-数据源
// sql-预编译的SQL语句
// parameters-预编译的SQL语句参数
func ExecDB(source string, sql string, parameters []any) (int64, error) {
	// 数据源
	db := DB(source)
	if db == nil {
		return 0, fmt.Errorf("not database source")
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

// PageNumSize 分页页码记录数
//
// pageNum-页码
// pageSize-记录数
func PageNumSize(pageNum, pageSize any) (int, int) {
	// 记录起始索引
	num := parse.Number(pageNum)
	if num < 1 {
		num = 1
	}

	// 显示记录数
	size := parse.Number(pageSize)
	if size < 0 {
		size = 10
	}
	return int(num - 1), int(size)
}
