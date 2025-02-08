package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/generate"
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
	tokenKey := constants.CACHE_LOGIN_TOKEN + ":" + uuid
	hasKey, err := redis.Has("", tokenKey)
	if hasKey > 0 && err == nil {
		_ = redis.Del("", tokenKey)
	}
	return claims[constants.JWT_USER_NAME].(string)
}

// Create 令牌生成
func Create(tokenInfo *TokenInfo, ilobArr [4]string) string {
	// 生成用户唯一token 32位
	tokenInfo.UUID = generate.Code(32)
	tokenInfo.LoginTime = time.Now().UnixMilli()

	// 设置请求用户登录客户端
	tokenInfo.LoginIp = ilobArr[0]
	tokenInfo.LoginLocation = ilobArr[1]
	tokenInfo.OS = ilobArr[2]
	tokenInfo.Browser = ilobArr[3]

	// 设置新登录IP和登录时间
	tokenInfo.User.LoginIp = tokenInfo.LoginIp
	tokenInfo.User.LoginTime = tokenInfo.LoginTime

	// 设置用户令牌有效期并存入缓存
	Cache(tokenInfo)

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
		constants.JWT_UUID:      tokenInfo.UUID,
		constants.JWT_USER_ID:   tokenInfo.UserId,
		constants.JWT_USER_NAME: tokenInfo.User.UserName,
		"ait":                   tokenInfo.LoginTime,
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
func Cache(tokenInfo *TokenInfo) {
	// 计算配置的有效期
	expTime := config.Get("jwt.expiresIn").(int)
	expTimestamp := time.Duration(expTime) * time.Minute
	iatTimestamp := time.Now().UnixMilli()
	tokenInfo.LoginTime = iatTimestamp
	tokenInfo.ExpireTime = iatTimestamp + expTimestamp.Milliseconds()
	tokenInfo.User.Password = ""
	// 登录信息标识缓存
	tokenKey := constants.CACHE_LOGIN_TOKEN + ":" + tokenInfo.UUID
	jsonBytes, err := json.Marshal(tokenInfo)
	if err != nil {
		return
	}
	_ = redis.SetByExpire("", tokenKey, string(jsonBytes), expTimestamp)
}

// RefreshIn 验证令牌有效期，相差不足xx分钟，自动刷新缓存
func RefreshIn(loginUser *TokenInfo) {
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

// Info 缓存的登录用户信息
func Info(claims jwt.MapClaims) TokenInfo {
	tokenInfo := TokenInfo{}
	uuid := claims[constants.JWT_UUID].(string)
	tokenKey := constants.CACHE_LOGIN_TOKEN + ":" + uuid
	hasKey, err := redis.Has("", tokenKey)
	if hasKey > 0 && err == nil {
		infoStr, err := redis.Get("", tokenKey)
		if infoStr == "" || err != nil {
			return tokenInfo
		}
		if err := json.Unmarshal([]byte(infoStr), &tokenInfo); err != nil {
			logger.Errorf("info json err : %v", err)
			return tokenInfo
		}
	}
	return tokenInfo
}
