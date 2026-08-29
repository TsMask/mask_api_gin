package token

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/shared/parse"
)

// TokenPair 新登录会话生成的一组令牌。
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64
	RefreshExpiresIn int64
}

// newSessionId 生成新的会话ID
func newSessionId() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// signingMethod 签名方法
func signingMethod() jwt.SigningMethod {
	switch config.GetStringOr("jwt.algorithm", "HS256") {
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

// tokenSecret 令牌密钥
func tokenSecret(tokenType string) string {
	secret := config.GetStringOr("jwt.secret", "")
	if tokenType == "refresh" {
		return "User_Refresh:" + secret
	}
	return "User_Access:" + secret
}

// tokenTTL 令牌过期时间
func tokenTTL(tokenType string) time.Duration {
	minutes := config.GetIntOr("jwt.expiresIn", 0)
	if tokenType == "refresh" {
		minutes = config.GetIntOr("jwt.refreshIn", 0)
	}
	return time.Duration(minutes) * time.Minute
}

// CreateTokenPair 为同一会话生成访问与刷新令牌。
func CreateTokenPair(userId int64, deviceFingerprint string) (TokenPair, string, error) {
	sessionId, err := newSessionId()
	if err != nil {
		return TokenPair{}, "", fmt.Errorf("generate session id: %w", err)
	}

	accessToken, accessExp := UserTokenCreate(userId, deviceFingerprint, "access", sessionId)
	refreshToken, refreshExp := UserTokenCreate(userId, deviceFingerprint, "refresh", sessionId)
	if accessToken == "" || refreshToken == "" || accessExp <= 0 || refreshExp <= 0 {
		return TokenPair{}, "", errors.New("token generation failed")
	}
	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  accessExp,
		RefreshExpiresIn: refreshExp,
	}, sessionId, nil
}

// UserTokenCreate 生成指定会话的令牌。
func UserTokenCreate(userId int64, deviceFingerprint, tokenType, sessionId string) (string, int64) {
	if userId <= 0 || len(sessionId) != 32 {
		return "", 0
	}
	now := time.Now()
	exp := now.Add(tokenTTL(tokenType))
	token := jwt.NewWithClaims(signingMethod(), jwt.MapClaims{
		constants.JWT_DEVICE_ID:  deviceFingerprint,
		constants.JWT_SESSION_ID: sessionId,
		constants.JWT_USER_ID:    userId,
		"exp":                    exp.Unix(),
		"iat":                    now.Unix(),
		"nbf":                    now.Add(-10 * time.Second).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(tokenSecret(tokenType)))
	if err != nil {
		logger.Errorf("jwt sign err: %v", err)
		return "", 0
	}
	return tokenStr, int64(exp.Sub(now).Seconds())
}

// UserTokenVerify 校验令牌签名、算法与基础声明。
func UserTokenVerify(tokenStr, tokenType string) (jwt.MapClaims, error) {
	expected := signingMethod()
	jwtToken, err := jwt.Parse(tokenStr, func(jToken *jwt.Token) (any, error) {
		if jToken.Method != expected {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(tokenSecret(tokenType)), nil
	})
	if err != nil {
		return nil, errors.New("token invalid")
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("token valid error")
	}
	if parse.Number(fmt.Sprint(claims[constants.JWT_USER_ID])) <= 0 ||
		len(fmt.Sprint(claims[constants.JWT_SESSION_ID])) != 32 {
		return nil, errors.New("token claims invalid")
	}
	return claims, nil
}

// sessionKey 会话ID缓存键
func sessionKey(sessionId string) string {
	return constants.CACHE_TOKEN_SESSION + ":" + sessionId
}

// userIndexKey 用户会话索引缓存键
func userIndexKey(userId int64) string {
	return fmt.Sprintf("%s:%d", constants.CACHE_TOKEN_USERS, userId)
}

// UserInfoRemove 删除当前访问会话及其用户索引。
func UserInfoRemove(tokenStr string) (string, error) {
	claims, err := UserTokenVerify(tokenStr, "access")
	if err != nil {
		return "", err
	}
	info := UserInfoGet(claims)
	if info.UserId <= 0 {
		return "", errors.New("session not found")
	}
	sessionId := fmt.Sprint(claims[constants.JWT_SESSION_ID])
	if err := cache.Del("", sessionKey(sessionId)); err != nil {
		return info.UserName, err
	}
	_ = cache.SRem("", userIndexKey(info.UserId), sessionId)
	return info.UserName, nil
}

// UserInfoCreate 生成访问会话缓存，并维护用户会话索引。
func UserInfoCreate(info *UserInfo, ilobArr [4]string) error {
	if info == nil || info.UserId <= 0 ||
		len(info.SessionId) != 32 {
		return errors.New("invalid session")
	}
	info.LoginIp = ilobArr[0]
	info.LoginLocation = ilobArr[1]
	info.LoginOS = ilobArr[2]
	info.LoginBrowser = ilobArr[3]

	accessTTL := tokenTTL("access")
	refreshTTL := tokenTTL("refresh")
	now := time.Now()
	info.LoginTime = now.UnixMilli()
	info.ExpireTime = now.Add(accessTTL).UnixMilli()

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal user info: %w", err)
	}
	// 会话以 SessionId 为主键，存活到 refresh token 过期；access token 自身 exp 负责短有效期。
	if err := cache.Set("", sessionKey(info.SessionId), string(jsonBytes), refreshTTL); err != nil {
		return fmt.Errorf("set session cache: %w", err)
	}
	if err := cache.SAdd("", userIndexKey(info.UserId), refreshTTL, info.SessionId); err != nil {
		_ = cache.Del("", sessionKey(info.SessionId))
		return fmt.Errorf("set session index: %w", err)
	}
	return nil
}

// UserInfoUpdate 更新访问会话并保留原 TTL。
func UserInfoUpdate(info UserInfo) {
	if info.SessionId == "" {
		return
	}
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return
	}
	ttl, err := cache.SetWithExistingTTL("", sessionKey(info.SessionId), string(jsonBytes))
	if err != nil || ttl < 0 {
		logger.Errorf("UserInfoUpdate err: %v ttl: %d", err, ttl)
	}
}

// UserInfoGet 读取会话，并校验令牌中的用户与会话 ID。
func UserInfoGet(claims jwt.MapClaims) UserInfo {
	info := UserInfo{}
	sessionId := fmt.Sprint(claims[constants.JWT_SESSION_ID])
	deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
	if len(sessionId) != 32 {
		return info
	}
	infoStr, err := cache.Get("", sessionKey(sessionId))
	if err != nil || infoStr == "" {
		return info
	}
	if err := json.Unmarshal([]byte(infoStr), &info); err != nil {
		logger.Errorf("user info json err: %v", err)
		return UserInfo{}
	}
	if info.SessionId != sessionId ||
		info.DeviceId != deviceId ||
		info.UserId != parse.Number(fmt.Sprint(claims[constants.JWT_USER_ID])) {
		return UserInfo{}
	}
	return info
}

// SessionValid 判断刷新令牌对应的会话仍然有效。
func SessionValid(claims jwt.MapClaims) bool {
	return UserInfoGet(claims).UserId > 0
}

// RemoveByUserIds 按用户索引批量清理会话。
func RemoveByUserIds(userIds []int64) {
	for _, userId := range userIds {
		indexKey := userIndexKey(userId)
		sessions, err := cache.SMembers("", indexKey)
		if err != nil || len(sessions) == 0 {
			_ = cache.Del("", indexKey)
			continue
		}
		keys := make([]string, 0, len(sessions))
		for _, sessionId := range sessions {
			keys = append(keys, sessionKey(sessionId))
		}
		_ = cache.DelKeys("", keys)
		_ = cache.Del("", indexKey)
	}
}

// RemoveBySessionIds 按会话ID清理会话，同步移除用户索引。
func RemoveBySessionIds(sessionIds []string) {
	for _, sessionId := range sessionIds {
		if len(sessionId) != 32 {
			continue
		}
		raw, err := cache.Get("", sessionKey(sessionId))
		if err != nil || raw == "" {
			continue
		}
		var info UserInfo
		if json.Unmarshal([]byte(raw), &info) != nil || info.UserId <= 0 {
			continue
		}
		if cache.Del("", sessionKey(sessionId)) == nil {
			_ = cache.SRem("", userIndexKey(info.UserId), sessionId)
		}
	}
}
