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
	"mask_api_gin/src/framework/utils/parse"
)

// Create 生成令牌
// userId 用户ID
// deviceFingerprint 设备指纹 SHA256
// tokenType 令牌类型 access:访问令牌 refresh:刷新令牌
func Create(userId int64, deviceFingerprint, tokenType string) (string, int64) {
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

	// 设置令牌过期时间
	now := time.Now()
	exp := now
	if tokenType == "access" {
		expiresIn := time.Duration(parse.Number(config.Get("jwt.expiresIn")))
		exp = now.Add(expiresIn * time.Minute)
	}
	if tokenType == "refresh" {
		refreshIn := time.Duration(parse.Number(config.Get("jwt.refreshIn")))
		exp = now.Add(refreshIn * time.Minute)
	}

	// 生成令牌负荷绑定uuid标识
	jwtToken := jwt.NewWithClaims(method, jwt.MapClaims{
		constants.JWT_DEVICE_ID: deviceFingerprint,
		constants.JWT_USER_ID:   userId,
		"exp":                   exp.Unix(), // 过期时间
		"iat":                   now.Unix(), // 签发时间
		"nbf":                   now.Unix(), // 生效时间
	})

	// 生成令牌设置密钥
	secret := fmt.Sprint(config.Get("jwt.secret"))
	tokenStr, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		logger.Infof("jwt sign err : %v", err)
		return "", 0
	}
	expSeconds := int64(exp.Sub(now).Seconds())
	return tokenStr, expSeconds
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
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token valid error")
}

// RemoveInfo 清除登录用户信息
func RemoveInfo(token string) (string, error) {
	claims, err := Verify(token)
	if err != nil {
		logger.Errorf("token verify err %v", err)
		return "", err
	}
	info := Info(claims)
	if info.User.UserName != "" {
		// 清除缓存KEY
		deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
		tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + deviceId
		return info.User.UserName, redis.Del("", tokenKey)
	}
	return "", fmt.Errorf("token invalid")
}

// CreateInfo 生成访问令牌
func CreateInfo(tokenInfo *TokenInfo, deviceFingerprint string, ilobArr [4]string) {
	tokenInfo.DeviceId = deviceFingerprint
	tokenInfo.LoginTime = time.Now().UnixMilli()

	// 设置请求用户登录客户端
	tokenInfo.LoginIp = ilobArr[0]
	tokenInfo.LoginLocation = ilobArr[1]
	tokenInfo.OS = ilobArr[2]
	tokenInfo.Browser = ilobArr[3]

	// 设置新登录IP和登录时间
	tokenInfo.User.LoginIp = tokenInfo.LoginIp
	tokenInfo.User.LoginTime = tokenInfo.LoginTime

	expiresIn := time.Duration(parse.Number(config.Get("jwt.expiresIn")))
	now := time.Now()
	exp := now.Add(expiresIn * time.Minute)
	tokenInfo.LoginTime = now.UnixMilli()
	tokenInfo.ExpireTime = exp.UnixMilli()
	tokenInfo.User.Password = ""
	// 登录信息标识缓存
	tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + tokenInfo.DeviceId
	jsonBytes, err := json.Marshal(tokenInfo)
	if err != nil {
		return
	}
	_ = redis.Set("", tokenKey, string(jsonBytes), expiresIn*time.Minute)
}

// UpdateInfo 生成访问令牌
func UpdateInfo(tokenInfo TokenInfo) {
	tokenInfo.User.Password = ""
	// 登录信息标识缓存
	tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + tokenInfo.DeviceId
	jsonBytes, err := json.Marshal(tokenInfo)
	if err != nil {
		return
	}
	expiresIn, _ := redis.GetExpire("", tokenKey)
	expiration := time.Duration(expiresIn) * time.Second
	_ = redis.Set("", tokenKey, string(jsonBytes), expiration)
}

// Info 缓存的登录用户信息
func Info(claims jwt.MapClaims) TokenInfo {
	tokenInfo := TokenInfo{}
	deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
	tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + deviceId
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
