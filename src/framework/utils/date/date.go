package date

import (
	"mask_api_gin/src/framework/logger"
	"time"
)

const (
	// 年 列如：2022
	YYYY = "2006"
	// 年-月 列如：2022-12
	YYYY_MM = "2006-01"
	// 年-月-日 列如：2022-12-30
	YYYY_MM_DD = "2006-01-02"
	// 年月日时分秒 列如：20221230010159
	YYYYMMDDHHMMSS = "20060102150405"
	// 年-月-日 时:分:秒 列如：2022-12-30 01:01:59
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
)

// 格式时间字符串
//
// dateStr 时间字符串
//
// formatStr 时间格式 默认YYYY-MM-DD HH:mm:ss
func ParseStrToDate(dateStr, formatStr string) time.Time {
	t, err := time.Parse(formatStr, dateStr)
	if err != nil {
		logger.Infof("utils ParseStrToDate err %v", err)
		return time.Time{}
	}
	return t
}

// 格式时间
//
// date 可转的Date对象
//
// formatStr 时间格式 默认YYYY-MM-DD HH:mm:ss
func ParseDateToStr(date any, formatStr string) string {
	t, ok := date.(time.Time)
	if !ok {
		switch v := date.(type) {
		case int64:
			if v == 0 {
				return ""
			}
			t = time.UnixMilli(v)
		case string:
			parsedTime, err := time.Parse(formatStr, v)
			if err != nil {
				logger.Infof("utils ParseDateToStr err %v", err)
				return ""
			}
			t = parsedTime
		default:
			return ""
		}
	}
	return t.Format(formatStr)
}

// 格式时间成日期路径
//
// 年/月 列如：2022/12
func ParseDatePath(date time.Time) string {
	return date.Format("2006/01")
}
