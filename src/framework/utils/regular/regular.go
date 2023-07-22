package regular

import (
	"regexp"

	"github.com/dlclark/regexp2"
)

// Replace 正则替换
func Replace(originStr, pattern, repStr string) string {
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(originStr, repStr)
}

// 判断是否为有效用户名格式
//
// 用户名不能以数字开头，可包含大写小写字母，数字，且不少于5位
func ValidUsername(username string) bool {
	if username == "" {
		return false
	}
	pattern := `^[a-zA-Z][a-z0-9A-Z]{5,}`
	match, err := regexp.MatchString(pattern, username)
	if err != nil {
		return false
	}
	return match
}

// 判断是否为有效密码格式
//
// 密码至少包含大小写字母、数字、特殊符号，且不少于6位
func ValidPassword(password string) bool {
	if password == "" {
		return false
	}
	pattern := `^(?![A-Za-z0-9]+$)(?![a-z0-9\W]+$)(?![A-Za-z\W]+$)(?![A-Z0-9\W]+$)[a-zA-Z0-9\W]{6,}$`
	re := regexp2.MustCompile(pattern, 0)
	match, err := re.MatchString(password)
	if err != nil {
		return false
	}
	return match
}

// 判断是否为有效手机号格式，1开头的11位手机号
func ValidMobile(mobile string) bool {
	if mobile == "" {
		return false
	}
	pattern := `^1[3|4|5|6|7|8|9][0-9]\d{8}$`
	match, err := regexp.MatchString(pattern, mobile)
	if err != nil {
		return false
	}
	return match
}

// 判断是否为有效邮箱格式
func ValidEmail(email string) bool {
	if email == "" {
		return false
	}
	pattern := `^(([^<>()\\.,;:\s@"]+(\.[^<>()\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+\.)+[a-zA-Z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]{2,}))$`
	re := regexp2.MustCompile(pattern, 0)
	match, err := re.MatchString(email)
	if err != nil {
		return false
	}
	return match
}

// 判断是否为http(s)://开头
//
// link 网络链接
func ValidHttp(link string) bool {
	if link == "" {
		return false
	}
	pattern := `^http(s)?:\/\/+`
	match, err := regexp.MatchString(pattern, link)
	if err != nil {
		return false
	}
	return match
}
