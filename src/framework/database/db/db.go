package db

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"

	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
func DB(source string) *gorm.DB {
	// 不指定时获取默认实例
	if source == "" {
		source = config.Get("gorm.defaultDataSourceName").(string)
	}
	return dbMap[source]
}

// RawDB 原生查询语句
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

// ExecDB 原生执行语句
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

// Unmarshal 将数据从 []map[string]any 反序列化到指定的结构体中
func Unmarshal(data []map[string]any, v any) error {
	// 获取目标结构体的反射值
	val := reflect.ValueOf(v)
	// 检查 v 是否为指针类型
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("v must be a non-nil pointer to a slice of structs")
	}

	// 获取指向的结构体值
	elem := val.Elem()
	// 检查指向的类型是否为切片
	if elem.Kind() != reflect.Slice {
		return fmt.Errorf("v must be a pointer to a slice")
	}

	// 遍历 data 中的每个 map
	for _, item := range data {
		// 为每个 map 创建一个新的实例
		newElem := reflect.New(elem.Type().Elem()).Elem()

		// 遍历结构体的字段
		for i := 0; i < newElem.NumField(); i++ {
			field := newElem.Type().Field(i)
			fieldValue := newElem.Field(i)

			// 获取 GORM column 标签
			tagStr := field.Tag.Get("gorm")
			columnName := extractColumnName(tagStr)
			if columnName == "" {
				continue
			}

			// 检查是否有对应的键值
			if value, ok := item[columnName]; ok {
				setFieldValue(fieldValue, value)
			}
		}

		// 将新的实例添加到目标切片中
		elem.Set(reflect.Append(elem, newElem))
	}

	return nil
}

// extractColumnName 提取以 "column:" 开头的字段名
func extractColumnName(tag string) string {
	if !strings.Contains(tag, "column:") {
		return ""
	}
	// 使用分号分割标签
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		// 查找以 "column:" 开头的部分, 提取并返回字段名
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

// setFieldValue 判断结构体内是否存在指定字段并设置值
func setFieldValue(fieldValue reflect.Value, value any) {
	// 检查字段是否存在
	if fieldValue.IsValid() && fieldValue.CanSet() {
		// 获取字段的类型
		fieldType := fieldValue.Type()

		// 转换传入的值类型为字段类型
		switch fieldType.Kind() {
		case reflect.String:
			if value == nil {
				fieldValue.SetString("")
			} else {
				fieldValue.SetString(fmt.Sprintf("%v", value))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if err != nil {
				intValue = 0
			}
			fieldValue.SetInt(intValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 64)
			if err != nil {
				uintValue = 0
			}
			fieldValue.SetUint(uintValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
			if err != nil {
				floatValue = 0
			}
			fieldValue.SetFloat(floatValue)
		case reflect.Struct:
			logger.Infof("%v 时间解析 %s %v \n", fieldValue, fieldValue.Type(), value)
			if fieldValue.Type() == reflect.TypeOf(time.Time{}) && value != nil {
				// 解析 value 并转换为 time.Time 类型
				parsedTime, err := time.Parse("2006-01-02 15:04:05 +0800 CST", fmt.Sprintf("%v", value))
				if err != nil {
					logger.Errorf("时间解析出错: %v", err)
				} else {
					// 设置字段的值
					fieldValue.Set(reflect.ValueOf(parsedTime))
				}
			}
		default:
			// 设置字段的值
			fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		}
	}
}

// ConvertIdsSlice 将 []string 转换为 []any
func ConvertIdsSlice(ids []int64) []any {
	// 将 []string 转换为 []any
	arr := make([]any, len(ids))
	for i, v := range ids {
		arr[i] = v
	}
	return arr
}

// KeyPlaceholderByQuery 查询-参数值的占位符
func KeyPlaceholderByQuery(sum int) string {
	placeholders := make([]string, sum)
	for i := 0; i < sum; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}

// KeyValuePlaceholderByInsert 插入-键值数据与参数映射键值占位符 keys, values, placeholder
func KeyValuePlaceholderByInsert(params map[string]any) (string, []any, string) {
	// 参数映射的键
	var keys []string
	// 参数映射的值
	var values []any
	// 参数值的占位符
	var placeholders []string

	for k, v := range params {
		keys = append(keys, k)
		values = append(values, v)
		placeholders = append(placeholders, "?")
	}

	return strings.Join(keys, ","), values, strings.Join(placeholders, ",")
}

// KeyValueByUpdate 更新-键值数据 keys, values
func KeyValueByUpdate(params map[string]any) (string, []any) {
	// 参数映射的键
	var keys []string
	// 参数映射的值
	var values []any

	for k, v := range params {
		keys = append(keys, k+"=?")
		values = append(values, v)
	}

	return strings.Join(keys, ","), values
}
