package token

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/generate"
	"mask_api_gin/src/framework/vo"

	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Remove 清除登录用户信息UUID
func Remove(token string) string {
	claims, err := Verify(token)
	if err != nil {
		logger.Errorf("token verify err %v", err)
		return ""
	}
	// 清除缓存KEY
	uuid := claims[constants.JWT_UUID].(string)
	tokenKey := constants.CACHE_LOGIN_TOKEN + uuid
	hasKey, err := redis.Has("", tokenKey)
	if hasKey > 0 && err == nil {
		_ = redis.Del("", tokenKey)
	}
	return claims[constants.JWT_USER_NAME].(string)
}

// Create 令牌生成
func Create(loginUser *vo.LoginUser, ilobArr [4]string) string {
	// 生成用户唯一token 32位
	loginUser.UUID = generate.Code(32)
	loginUser.LoginTime = time.Now().UnixMilli()

	// 设置请求用户登录客户端
	loginUser.LoginIp = ilobArr[0]
	loginUser.LoginLocation = ilobArr[1]
	loginUser.OS = ilobArr[2]
	loginUser.Browser = ilobArr[3]

	// 设置新登录IP和登录时间
	loginUser.User.LoginIp = loginUser.LoginIp
	loginUser.User.LoginTime = loginUser.LoginTime

	// 设置用户令牌有效期并存入缓存
	Cache(loginUser)

	// 令牌算法 HS256 HS384 HS512
	algorithm := config.Get("jwt.algorithm").(string)
	var method *jwt.SigningMethodHMAC
	switch algorithm {
	case "HS512":
		method = jwt.SigningMethodHS512
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS256":
	default:
		method = jwt.SigningMethodHS256
	}
	// 生成令牌负荷绑定uuid标识
	jwtToken := jwt.NewWithClaims(method, jwt.MapClaims{
		constants.JWT_UUID:      loginUser.UUID,
		constants.JWT_USER_ID:   loginUser.UserId,
		constants.JWT_USER_NAME: loginUser.User.UserName,
		"exp":                   loginUser.ExpireTime,
		"ait":                   loginUser.LoginTime,
	})

	// 生成令牌设置密钥
	secret := config.Get("jwt.secret").(string)
	tokenStr, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		logger.Infof("jwt sign err : %v", err)
		return ""
	}
	return tokenStr
}

// Cache 缓存登录用户信息
func Cache(loginUser *vo.LoginUser) {
	// 计算配置的有效期
	expTime := config.Get("jwt.expiresIn").(int)
	expTimestamp := time.Duration(expTime) * time.Minute
	iatTimestamp := time.Now().UnixMilli()
	loginUser.LoginTime = iatTimestamp
	loginUser.ExpireTime = iatTimestamp + expTimestamp.Milliseconds()
	loginUser.User.Password = ""
	// 根据登录标识将loginUser缓存
	tokenKey := constants.CACHE_LOGIN_TOKEN + loginUser.UUID
	jsonBytes, err := json.Marshal(loginUser)
	if err != nil {
		return
	}
	_ = redis.SetByExpire("", tokenKey, string(jsonBytes), expTimestamp)
}

// RefreshIn 验证令牌有效期，相差不足xx分钟，自动刷新缓存
func RefreshIn(loginUser *vo.LoginUser) {
	// 相差不足xx分钟，自动刷新缓存
	refreshTime := config.Get("jwt.refreshIn").(int)
	refreshTimestamp := time.Duration(refreshTime) * time.Minute
	// 过期时间
	expireTimestamp := loginUser.ExpireTime
	currentTimestamp := time.Now().UnixMilli()
	if expireTimestamp-currentTimestamp <= refreshTimestamp.Milliseconds() {
		Cache(loginUser)
	}
}

// Verify 校验令牌是否有效
func Verify(token string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(token, func(jToken *jwt.Token) (any, error) {
		// 判断加密算法是预期的加密算法
		if _, ok := jToken.Method.(*jwt.SigningMethodHMAC); ok {
			secret := config.Get("jwt.secret").(string)
			return []byte(secret), nil
		}
		return nil, jwt.ErrSignatureInvalid
	})
	if err != nil {
		logger.Errorf("Token Verify Err: %v", err)
		return nil, fmt.Errorf("token invalid")
	}
	// 如果解析负荷成功并通过签名校验
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token valid error")
}

// LoginUser 缓存的登录用户信息
func LoginUser(claims jwt.MapClaims) vo.LoginUser {
	loginUser := vo.LoginUser{}
	uuid := claims[constants.JWT_UUID].(string)
	tokenKey := constants.CACHE_LOGIN_TOKEN + uuid
	hasKey, err := redis.Has("", tokenKey)
	if hasKey > 0 && err == nil {
		loginUserStr, err := redis.Get("", tokenKey)
		if loginUserStr == "" || err != nil {
			return loginUser
		}
		if err := json.Unmarshal([]byte(loginUserStr), &loginUser); err != nil {
			logger.Errorf("loginuser info json err : %v", err)
			return loginUser
		}
	}
	return loginUser
}
