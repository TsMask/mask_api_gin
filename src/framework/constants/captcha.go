package constants

// 验证码常量信息
const (
	// CAPTCHA_EXPIRATION 验证码有效期，单位秒
	CAPTCHA_EXPIRATION = 2 * 60
	// CAPTCHA_TYPE_CHAR 验证码类型-数值计算
	CAPTCHA_TYPE_CHAR = "char"
	// CAPTCHA_TYPE_MATH 验证码类型-字符验证
	CAPTCHA_TYPE_MATH = "math"
)
