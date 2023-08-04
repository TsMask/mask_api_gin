package repo

import (
	"fmt"
	"mask_api_gin/src/framework/utils/parse"
	"reflect"
	"strconv"
	"strings"
)

// PageNumSize 分页页码记录数
func PageNumSize(pageNum, pageSize any) (int64, int64) {
	// 记录起始索引
	num := parse.Number(pageNum)
	if num > 5000 {
		num = 5000
	}
	if num < 1 {
		num = 1
	}

	// 显示记录数
	size := parse.Number(pageSize)
	if size > 50000 {
		size = 50000
	}
	if size < 0 {
		size = 10
	}
	return num - 1, size
}

// SetFieldValue 判断结构体内是否存在指定字段并设置值
func SetFieldValue(obj any, fieldName string, value any) {
	// 获取结构体的反射值
	userValue := reflect.ValueOf(obj)

	// 获取字段的反射值
	fieldValue := userValue.Elem().FieldByName(fieldName)

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
		default:
			// 设置字段的值
			fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		}
	}
}

// ConvertIdsSlice 将 []string 转换为 []any
func ConvertIdsSlice(ids []string) []any {
	// 将 []string 转换为 []any
	arr := make([]any, len(ids))
	for i, v := range ids {
		arr[i] = v
	}
	return arr
}

// 查询-参数值的占位符
func KeyPlaceholderByQuery(sum int) string {
	placeholders := make([]string, sum)
	for i := 0; i < sum; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}

// 插入-参数映射键值占位符 keys, placeholder, values
func KeyPlaceholderValueByInsert(params map[string]any) ([]string, string, []any) {
	// 参数映射的键
	keys := make([]string, len(params))
	// 参数映射的值
	values := make([]any, len(params))
	sum := 0
	for k, v := range params {
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

// 更新-参数映射键值占位符 keys, values
func KeyValueByUpdate(params map[string]any) ([]string, []any) {
	// 参数映射的键
	keys := make([]string, len(params))
	// 参数映射的值
	values := make([]any, len(params))
	sum := 0
	for k, v := range params {
		keys[sum] = k + "=?"
		values[sum] = v
		sum++
	}
	return keys, values
}
