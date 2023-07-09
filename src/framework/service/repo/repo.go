package repo

import (
	"mask_api_gin/src/framework/logger"
	"reflect"
	"strconv"
	"strings"
)

// 分页页码记录数
func PageNumSize(pageNum, pageSize string) (int, int) {
	// 记录起始索引
	num, err := strconv.Atoi(pageNum)
	if err != nil {
		logger.Errorf("PageNumSize strconv int num err %v", err)
		num = 0
	}
	if num > 5000 {
		num = 5000
	}
	if num < 1 {
		num = 1
	}

	// 显示记录数
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		logger.Errorf("PageNumSize strconv int size err %v", err)
		size = 10
	}
	if size > 50000 {
		size = 50000
	}
	if size < 0 {
		size = 10
	}
	return num - 1, size
}

// 转换记录结果
func ConvertResultRows(results interface{}) []interface{} {
	s := reflect.ValueOf(results)
	if s.Kind() != reflect.Slice {
		logger.Errorf("ConvertResultRows not a slice")
	}

	rows := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		rows[i] = s.Index(i).Interface()
	}
	return rows
}

// 参数映射键值占位符 keys, placeholder, values
func KeyValuePlaceholder(m map[string]interface{}) ([]string, string, []interface{}) {
	// 参数映射的键
	keys := make([]string, len(m))
	// 参数映射的值
	values := make([]interface{}, len(m))
	sum := 0
	for k, v := range m {
		keys[sum] = k
		values[sum] = v
		sum++
	}
	// 参数值的占位符
	placeholders := make([]string, sum)
	for i := 0; i < sum; i++ {
		placeholders[i] = "?"
	}
	return keys, strings.Join(placeholders, ","), values
}
