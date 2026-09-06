package constants

// 缓存的key常量
const (
	// CACHE_TOKEN_SESSION 登录用户会话标识
	CACHE_TOKEN_SESSION = "token_sessions"
	// CACHE_TOKEN_USERS 用户维度会话索引
	CACHE_TOKEN_USERS = "token_users"
	// CACHE_CAPTCHA_CODE 验证码
	CACHE_CAPTCHA_CODE = "captcha_codes"
	// CACHE_SMS_CODE 短信登录验证码
	CACHE_SMS_CODE = "sms_codes"
	// CACHE_SMS_CODE_COOLDOWN 短信登录验证码发送冷却
	CACHE_SMS_CODE_COOLDOWN = "sms_codes_cooldown"
	// CACHE_PASSWORD_RESET_CODE 找回密码验证码
	CACHE_PASSWORD_RESET_CODE = "password_reset_codes"
	// CACHE_PASSWORD_RESET_COOLDOWN 找回密码验证码发送冷却
	CACHE_PASSWORD_RESET_COOLDOWN = "password_reset_cooldown"
	// CACHE_SYS_CONFIG 参数管理
	CACHE_SYS_CONFIG = "sys_config"
	// CACHE_SYS_DICT 字典管理
	CACHE_SYS_DICT = "sys_dict"
	// CACHE_REPEAT_SUBMIT  防重提交
	CACHE_REPEAT_SUBMIT = "repeat_submit"
	// CACHE_RATE_LIMIT 限流
	CACHE_RATE_LIMIT = "rate_limit"
	// CACHE_PWD_ERR_COUNT 登录账户密码错误次数
	CACHE_PWD_ERR_COUNT = "pwd_err_count"
)
