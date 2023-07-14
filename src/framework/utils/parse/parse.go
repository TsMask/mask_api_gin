package parse

import (
	"image/color"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Number 解析数值型
func Number(str interface{}) int64 {
	switch str := str.(type) {
	case string:
		if str == "" {
			return 0
		}
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return num
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(str).Int()
	case float32, float64:
		return int64(reflect.ValueOf(str).Float())
	default:
		return 0
	}
}

// Boolean 解析布尔型
func Boolean(str interface{}) bool {
	switch str := str.(type) {
	case string:
		if str == "" || str == "false" || str == "0" {
			return false
		}
		// 尝试将字符串解析为数字
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return num != 0
		}
		return true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		num := reflect.ValueOf(str).Int()
		return num != 0
	case float32, float64:
		num := reflect.ValueOf(str).Float()
		return num != 0
	default:
		return false
	}
}

// 解析首字母转大写
//
// 字符串 abc_123!@# 结果 Abc_123
func ParseFirstUpper(str string) string {
	if len(str) == 0 {
		return str
	}
	reg := regexp.MustCompile(`[^_\w]+`)
	str = reg.ReplaceAllString(str, "")
	return strings.ToUpper(str[:1]) + str[1:]
}

// RemoveDuplicates 数组内字符串去重
func RemoveDuplicates(ids []string) []string {
	uniqueIDs := make(map[string]bool)
	uniqueIDSlice := make([]string, 0, len(ids))

	for _, id := range ids {
		if _, ok := uniqueIDs[id]; !ok {
			uniqueIDs[id] = true
			uniqueIDSlice = append(uniqueIDSlice, id)
		}
	}

	return uniqueIDSlice
}

// Color 解析颜色 #fafafa
func Color(colorStr string) *color.RGBA {
	// 去除 # 号
	colorStr = colorStr[1:]

	// 将颜色字符串拆分为 R、G、B 分量
	r, _ := strconv.ParseInt(colorStr[0:2], 16, 0)
	g, _ := strconv.ParseInt(colorStr[2:4], 16, 0)
	b, _ := strconv.ParseInt(colorStr[4:6], 16, 0)

	return &color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255, // 不透明
	}
}
