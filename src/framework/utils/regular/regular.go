package regular

import (
	"regexp"

	"github.com/dlclark/regexp2"
)

// regexMatch 正则匹配字符串
func regexMatch(pattern, str string) bool {
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}

// ValidHttp 判断是否为http(s)://开头
func ValidHttp(s string) bool {
	pattern := `^http(s)?://`
	return regexMatch(pattern, s)
}

// ValidUsername 判断是否为有效用户名格式
//
// 用户账号只能包含大写小写字母，数字，且不少于4位
func ValidUsername(s string) bool {
	pattern := `[a-z0-9A-Z]{3,}$`
	return regexMatch(pattern, s)
}

// ValidMobile 判断是否为有效手机号格式，1开头的11位手机号
func ValidMobile(s string) bool {
	pattern := `^1[3-9]\d{9}$`
	return regexMatch(pattern, s)
}

// regex2Match 第三方正则匹配字符串
func regex2Match(pattern, str string) bool {
	re := regexp2.MustCompile(pattern, 0)
	match, err := re.MatchString(str)
	if err != nil {
		return false
	}
	return match
}

// ValidPassword 判断是否为有效密码格式
//
// 密码至少包含大小写字母、数字、特殊符号，且不少于6位
func ValidPassword(s string) bool {
	pattern := `^(?![A-Za-z0-9]+$)(?![a-z0-9\W]+$)(?![A-Za-z\W]+$)(?![A-Z0-9\W]+$)[a-zA-Z0-9\W]{6,}$`
	return regex2Match(pattern, s)
}

// ValidEmail 判断是否为有效邮箱格式
func ValidEmail(s string) bool {
	pattern := `^(([^<>()\\.,;:\s@"]+(\.[^<>()\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+\.)+[a-zA-Z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]{2,}))$`
	return regex2Match(pattern, s)
}

// Replace 正则替换字符串中匹配的字符串
func Replace(pattern, src, repl string) string {
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(src, repl)
}
