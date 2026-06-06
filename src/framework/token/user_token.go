package token

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
)

// validateDeviceId 校验设备标识格式，防止注入特殊字符
func validateDeviceId(deviceId string) bool {
	if len(deviceId) == 0 || len(deviceId) > 128 {
		return false
	}
	if strings.ContainsAny(deviceId, "*? \t\n\r") {
		return false
	}
	return true
}

// UserTokenCreate 生成令牌
// userId 用户ID
// deviceFingerprint 设备指纹 SHA256
// tokenType 令牌类型 access:访问令牌 refresh:刷新令牌
func UserTokenCreate(userId int64, deviceFingerprint, tokenType string) (string, int64) {
	if !validateDeviceId(deviceFingerprint) {
		return "", 0
	}

	// 令牌算法 HS256 HS384 HS512，默认 HS256
	algorithm := config.Get("jwt.algorithm").(string)
	method := jwt.SigningMethodHS256
	switch algorithm {
	case "HS512":
		method = jwt.SigningMethodHS512
	case "HS384":
		method = jwt.SigningMethodHS384
	}

	// 生成令牌设置密钥
	secret := fmt.Sprint(config.Get("jwt.secret"))
	// 设置令牌过期时间
	now := time.Now()
	exp := now
	if tokenType == "access" {
		expiresIn := time.Duration(parse.Number(config.Get("jwt.expiresIn")))
		exp = now.Add(expiresIn * time.Minute)
		secret = "User_Access:" + secret
	}
	if tokenType == "refresh" {
		refreshIn := time.Duration(parse.Number(config.Get("jwt.refreshIn")))
		exp = now.Add(refreshIn * time.Minute)
		secret = "User_Refresh:" + secret
	}

	// 生成令牌负载绑定设备标识
	jwtToken := jwt.NewWithClaims(method, jwt.MapClaims{
		constants.JWT_DEVICE_ID: deviceFingerprint,
		constants.JWT_USER_ID:   userId,
		"exp":                   exp.Unix(),                        // 过期时间
		"iat":                   now.Unix(),                        // 签发时间
		"nbf":                   now.Add(-10 * time.Second).Unix(), // 生效时间
	})

	tokenStr, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		logger.Infof("jwt sign err : %v", err)
		return "", 0
	}
	expSeconds := int64(exp.Sub(now).Seconds())
	return tokenStr, expSeconds
}

// UserTokenVerify 校验令牌是否有效
// tokenType 令牌类型 access:访问令牌 refresh:刷新令牌
func UserTokenVerify(tokenStr string, tokenType string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(tokenStr, func(jToken *jwt.Token) (any, error) {
		// 判断加密算法是预期的加密算法
		if _, ok := jToken.Method.(*jwt.SigningMethodHMAC); ok {
			secret := config.Get("jwt.secret").(string)
			if tokenType == "access" {
				secret = "User_Access:" + secret
			}
			if tokenType == "refresh" {
				secret = "User_Refresh:" + secret
			}
			return []byte(secret), nil
		}
		return nil, jwt.ErrSignatureInvalid
	})
	if err != nil {
		logger.Errorf("Token Verify Err: %v", err)
		return nil, fmt.Errorf("token invalid")
	}
	// 如果解析负载成功并通过签名校验
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token valid error")
}

// UserInfoRemove 清除访问用户信息缓存
func UserInfoRemove(tokenStr string) (string, error) {
	claims, err := UserTokenVerify(tokenStr, "access")
	if err != nil {
		logger.Errorf("token verify err %v", err)
		return "", err
	}
	info := UserInfoGet(claims)
	if info.User.UserName != "" {
		// 清除缓存KEY
		deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
		tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + deviceId
		return info.User.UserName, redis.Del("", tokenKey)
	}
	return "", fmt.Errorf("token invalid")
}

// UserInfoCreate 生成访问用户信息缓存
func UserInfoCreate(info *UserInfo, deviceFingerprint string, ilobArr [4]string) error {
	if !validateDeviceId(deviceFingerprint) {
		return fmt.Errorf("invalid device fingerprint")
	}
	info.DeviceId = deviceFingerprint

	// 设置请求用户登录客户端
	info.LoginIp = ilobArr[0]
	info.LoginLocation = ilobArr[1]
	info.OS = ilobArr[2]
	info.Browser = ilobArr[3]

	expiresIn := time.Duration(parse.Number(config.Get("jwt.expiresIn")))
	now := time.Now()
	exp := now.Add(expiresIn * time.Minute)
	info.LoginTime = now.UnixMilli()
	info.ExpireTime = exp.UnixMilli()
	// 设置新登录IP和登录时间
	info.User.LoginIp = info.LoginIp
	info.User.LoginTime = info.LoginTime
	info.User.Password = ""
	// 登录信息标识缓存
	tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + info.DeviceId
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("UserInfoCreate marshal err: %w", err)
	}
	if err := redis.Set("",
		tokenKey,
		string(jsonBytes),
		expiresIn*time.Minute,
	); err != nil {
		return fmt.Errorf("UserInfoCreate redis set err: %w", err)
	}
	return nil
}

// UserInfoUpdate 更新访问用户信息缓存
func UserInfoUpdate(info UserInfo) {
	info.User.Password = ""
	tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + info.DeviceId
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return
	}
	// 保留原 TTL 原子更新值，避免 TTL 竞争
	ttl, err := redis.SetWithExistingTTL("", tokenKey, string(jsonBytes))
	if err != nil || ttl < 0 {
		logger.Errorf("UserInfoUpdate SetWithExistingTTL err: %v ttl: %d", err, ttl)
	}
}

// UserInfoGet 缓存的访问用户消息
func UserInfoGet(claims jwt.MapClaims) UserInfo {
	info := UserInfo{}
	deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
	if !validateDeviceId(deviceId) {
		return info
	}
	tokenKey := constants.CACHE_TOKEN_DEVICE + ":" + deviceId
	infoStr, err := redis.Get("", tokenKey)
	if err != nil || infoStr == "" {
		return info
	}
	if err := json.Unmarshal([]byte(infoStr), &info); err != nil {
		logger.Errorf("info json err : %v", err)
		return info
	}
	return info
}

// RemoveByUserIds 根据 userId 列表清理对应的 Token 缓存
func RemoveByUserIds(userIds []int64) {
	if len(userIds) == 0 {
		return
	}
	userIdSet := make(map[int64]struct{}, len(userIds))
	for _, id := range userIds {
		userIdSet[id] = struct{}{}
	}
	keys, err := redis.GetKeys("", constants.CACHE_TOKEN_DEVICE+":*")
	if err != nil || len(keys) == 0 {
		return
	}
	// 批量 MGET 替代逐个 GET，减少网络往返
	var delKeys []string
	for i := 0; i < len(keys); i += 50 {
		end := i + 50
		if end > len(keys) {
			end = len(keys)
		}
		batch := keys[i:end]
		values, err := redis.GetBatch("", batch)
		if err != nil {
			continue
		}
		for j, v := range values {
			if v == nil {
				continue
			}
			str := fmt.Sprint(v)
			if str == "" {
				continue
			}
			var info UserInfo
			if json.Unmarshal([]byte(str), &info) != nil {
				continue
			}
			if _, ok := userIdSet[info.UserId]; ok {
				delKeys = append(delKeys, batch[j])
			}
		}
	}
	if len(delKeys) > 0 {
		_ = redis.DelKeys("", delKeys)
	}
}
