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
func Number(str any) int64 {
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
func Boolean(str any) bool {
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

// FirstUpper 首字母转大写
//
// 字符串 abc_123!@# 结果 Abc_123
func FirstUpper(str string) string {
	if len(str) == 0 {
		return str
	}
	reg := regexp.MustCompile(`[^_\w]+`)
	str = reg.ReplaceAllString(str, "")
	return strings.ToUpper(str[:1]) + str[1:]
}

// Bit 比特位为单位
func Bit(bit float64) string {
	var GB, MB, KB string

	if bit > float64(1<<30) {
		GB = fmt.Sprintf("%0.2f", bit/(1<<30))
	}

	if bit > float64(1<<20) && bit < (1<<30) {
		MB = fmt.Sprintf("%.2f", bit/(1<<20))
	}

	if bit > float64(1<<10) && bit < (1<<20) {
		KB = fmt.Sprintf("%.2f", bit/(1<<10))
	}

	if GB != "" {
		return GB + "GB"
	} else if MB != "" {
		return MB + "MB"
	} else if KB != "" {
		return KB + "KB"
	} else {
		return fmt.Sprintf("%vB", bit)
	}
}

// CronExpression 解析 Cron 表达式，返回下一次执行的时间戳（毫秒）
//
// 【*/5 * * * * ?】 6个参数
func CronExpression(expression string) int64 {
	specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := specParser.Parse(expression)
	if err != nil {
		fmt.Println(err)
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
func RemoveDuplicates(ids []string) []string {
	uniqueIDs := make(map[string]bool)
	uniqueIDSlice := make([]string, 0)

	for _, id := range ids {
		_, ok := uniqueIDs[id]
		if !ok && id != "" {
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
