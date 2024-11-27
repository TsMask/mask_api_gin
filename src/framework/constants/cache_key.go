package constants

// 缓存的key常量
const (
	// CACHE_LOGIN_TOKEN 登录用户
	CACHE_LOGIN_TOKEN = "login_tokens:"
	// CACHE_CAPTCHA_CODE 验证码
	CACHE_CAPTCHA_CODE = "captcha_codes:"
	// CACHE_SYS_CONFIG 参数管理
	CACHE_SYS_CONFIG = "sys_config:"
	// CACHE_SYS_DICT 字典管理
	CACHE_SYS_DICT = "sys_dict:"
	// CACHE_REPEAT_SUBMIT  防重提交
	CACHE_REPEAT_SUBMIT = "repeat_submit:"
	// CACHE_RATE_LIMIT 限流
	CACHE_RATE_LIMIT = "rate_limit:"
	// CACHE_PWD_ERR_COUNT 登录账户密码错误次数
	CACHE_PWD_ERR_COUNT = "pwd_err_count:"
)
