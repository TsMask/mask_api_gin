package token

import (
	"encoding/json"
	"errors"
	"mask_api_gin/src/framework/config"
	cachekeyConstants "mask_api_gin/src/framework/constants/cachekey"
	tokenConstants "mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/logger"
	redisCahe "mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/generate"
	"mask_api_gin/src/framework/vo"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Remove 清除登录用户信息UUID
func Remove(tokenStr string) string {
	claims, err := Verify(tokenStr)
	if err != nil {
		logger.Errorf("token verify err %v", err)
		return ""
	}
	// 清除缓存KEY
	uuid := claims[tokenConstants.JWT_UUID].(string)
	tokenKey := cachekeyConstants.LOGIN_TOKEN_KEY + uuid
	if redisCahe.Has(tokenKey) {
		redisCahe.Del(tokenKey)
	}
	return claims[tokenConstants.JWT_NAME].(string)
}

// Create 令牌生成
func Create(loginUser *vo.LoginUser, ilobArgs ...string) string {
	// 生成用户唯一tokne32位
	loginUser.UUID = generate.Code(32)

	// 设置请求用户登录客户端
	loginUser.IPAddr = ilobArgs[0]
	loginUser.LoginLocation = ilobArgs[1]
	loginUser.OS = ilobArgs[2]
	loginUser.Browser = ilobArgs[3]

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
		tokenConstants.JWT_UUID: loginUser.UUID,
		tokenConstants.JWT_KEY:  loginUser.UserID,
		tokenConstants.JWT_NAME: loginUser.User.UserName,
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
	iatTimestamp := date.NowTimestamp()
	loginUser.LoginTime = iatTimestamp
	loginUser.ExpireTime = iatTimestamp + expTimestamp.Milliseconds()
	// 根据登录标识将loginUser缓存
	tokenKey := cachekeyConstants.LOGIN_TOKEN_KEY + loginUser.UUID
	jsonBytes, err := json.Marshal(loginUser)
	if err != nil {
		return
	}
	redisCahe.SetByExpire(tokenKey, string(jsonBytes), expTimestamp)
}

// RefreshIn 验证令牌有效期，相差不足xx分钟，自动刷新缓存
func RefreshIn(loginUser *vo.LoginUser) {
	// 相差不足xx分钟，自动刷新缓存
	refreshTime := config.Get("jwt.refreshIn").(int)
	refreshTimestamp := time.Duration(refreshTime) * time.Minute
	// 过期时间
	expireTimestamp := loginUser.ExpireTime
	currentTimestamp := date.NowTimestamp()
	if expireTimestamp-currentTimestamp <= refreshTimestamp.Milliseconds() {
		Cache(loginUser)
	}
}

// Verify 校验令牌是否有效
func Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 判断加密算法是预期的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			secret := config.Get("jwt.secret").(string)
			return []byte(secret), nil
		}
		return nil, jwt.ErrSignatureInvalid
	})
	if err != nil {
		return nil, err
	}
	// 如果解析负荷成功并通过签名校验
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token Valid err")
}

// LoginUser 缓存的登录用户信息
func LoginUser(claims jwt.MapClaims) vo.LoginUser {
	uuid := claims[tokenConstants.JWT_UUID].(string)
	tokenKey := cachekeyConstants.LOGIN_TOKEN_KEY + uuid
	var loginUser vo.LoginUser
	if redisCahe.Has(tokenKey) {
		loginUserStr := redisCahe.Get(tokenKey)
		if loginUserStr == "" {
			return loginUser
		}
		err := json.Unmarshal([]byte(loginUserStr), &loginUser)
		if err != nil {
			logger.Errorf("loginuser info json err : %v", err)
			return loginUser
		}
		return loginUser
	}
	return loginUser
}
