package token

import (
	"encoding/json"
	"mask_api_gin/src/framework/cache/redis"
	"mask_api_gin/src/framework/config"
	cachekeyConstants "mask_api_gin/src/framework/constants/cachekey"
	tokenConstants "mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/model"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/generate"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Remove 清除登录用户信息UUID
func Remove(token string) bool {
	return true
}

// Create 令牌生成
func Create(loginUser *model.LoginUser, ilobArgs ...string) string {
	// 生成用户唯一tokne32位
	loginUser.UUID = generate.Code(32)
	// 设置请求用户登录客户端
	loginUser.IPAddr = ilobArgs[0]
	loginUser.LoginLocation = ilobArgs[1]
	loginUser.OS = ilobArgs[2]
	loginUser.Browser = ilobArgs[3]
	// 设置用户令牌有效期并存入缓存
	cacheTokenUUID(loginUser)
	// 生成令牌负荷绑定uuid标识
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		tokenConstants.JWT_UUID: loginUser.UUID,
		tokenConstants.JWT_KEY:  loginUser.UserID,
		"exp":                   loginUser.ExpireTime,
		"ait":                   loginUser.LoginTime,
	})
	// 生成令牌设置密钥
	key := config.Get("jwt.secret").(string)
	tokenStr, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		logger.Infof("jwt sign err : %v", err)
		return ""
	}
	return tokenStr
}

// cacheTokenUUID 缓存登录用户信息UUID
func cacheTokenUUID(loginUser *model.LoginUser) {
	// 计算配置的有效期
	expTimestamp := config.Get("jwt.expiresIn").(int)
	expTime := time.Duration(expTimestamp) * time.Minute
	iatTimestamp := date.NowTimestamp()
	loginUser.LoginTime = iatTimestamp
	loginUser.ExpireTime = iatTimestamp + expTime.Milliseconds()
	// 根据登录标识将loginUser缓存
	tokenKey := cachekeyConstants.LOGIN_TOKEN_KEY + loginUser.UUID
	jsonBytes, err := json.Marshal(loginUser)
	if err != nil {
		return
	}
	redis.SetByExpire(tokenKey, string(jsonBytes), expTime)
}
