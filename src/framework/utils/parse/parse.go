package parse

import (
	"fmt"
	"image/color"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// Number 解析数值型
func Number(value any) int64 {
	switch v := value.(type) {
	case string:
		if v == "" {
			return 0
		}
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0
		}
		return num
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(v).Uint())
	case float32, float64:
		return int64(reflect.ValueOf(v).Float())
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// Boolean 解析布尔型
func Boolean(value any) bool {
	switch v := value.(type) {
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return b
	case int, int8, int16, int32, int64:
		num := reflect.ValueOf(v).Int()
		return num != 0
	case uint, uint8, uint16, uint32, uint64:
		num := int64(reflect.ValueOf(v).Uint())
		return num != 0
	case float32, float64:
		num := reflect.ValueOf(v).Float()
		return num != 0
	case bool:
		return v
	default:
		return false
	}
}

// ConvertToCamelCase 字符串转换驼峰形式
//
// 字符串 dict/inline/data/:dictId 结果 DictInlineDataDictId
func ConvertToCamelCase(str string) string {
	if len(str) == 0 {
		return str
	}
	reg := regexp.MustCompile(`[-_:/]\w`)
	result := reg.ReplaceAllStringFunc(str, func(match string) string {
		return strings.ToUpper(string(match[1]))
	})

	words := strings.Fields(result)
	for i, word := range words {
		str := word[1:]
		str = strings.ReplaceAll(str, "/", "")
		words[i] = strings.ToUpper(word[:1]) + str
	}

	return strings.Join(words, "")
}

// Bit 比特位为单位 1023.00 B --> 1.00 KB
func Bit(bit float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	for i := 0; i < len(units); i++ {
		if bit < 1024 || i == len(units)-1 {
			return fmt.Sprintf("%.2f %s", bit, units[i])
		}
		bit /= 1024
	}
	return ""
}

// CronExpression 解析 Cron 表达式，返回下一次执行的时间戳（毫秒）
//
// 【*/5 * * * * ?】 6个参数
func CronExpression(expression string) int64 {
	specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := specParser.Parse(expression)
	if err != nil {
		return 0
	}
	return schedule.Next(time.Now()).UnixMilli()
}

// SafeContent 内容值进行安全掩码
func SafeContent(value string) string {
	if len(value) < 3 {
		return strings.Repeat("*", len(value))
	} else if len(value) < 6 {
		return string(value[0]) + strings.Repeat("*", len(value)-1)
	} else if len(value) < 10 {
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	} else if len(value) < 15 {
		return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
	} else {
		return value[:3] + strings.Repeat("*", len(value)-6) + value[len(value)-3:]
	}
}

// RemoveDuplicates 数组内字符串去重
func RemoveDuplicates(arr []string) []string {
	uniqueIDs := make(map[string]bool)
	uniqueIDSlice := make([]string, 0)

	for _, id := range arr {
		_, ok := uniqueIDs[id]
		if !ok && id != "" {
			uniqueIDs[id] = true
			uniqueIDSlice = append(uniqueIDSlice, id)
		}
	}

	return uniqueIDSlice
}

// RemoveDuplicatesToNumber 数组内字符串分隔去重转为整型数组
func RemoveDuplicatesToNumber(keyStr, sep string) []int64 {
	arr := make([]int64, 0)
	if keyStr == "" {
		return arr
	}
	if strings.Contains(keyStr, sep) {
		// 处理字符转数组后去重
		strArr := strings.Split(keyStr, sep)
		uniqueKeys := make(map[string]bool)
		for _, str := range strArr {
			_, ok := uniqueKeys[str]
			if !ok && str != "" {
				uniqueKeys[str] = true
				val := Number(str)
				if val != 0 {
					arr = append(arr, val)
				}
			}
		}
	} else {
		val := Number(keyStr)
		if val != 0 {
			arr = append(arr, val)
		}
	}
	return arr
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
