package crypto

import "golang.org/x/crypto/bcrypt"

// BcryptHash Bcrypt密码加密
func BcryptHash(originStr string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(originStr), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// BcryptCompare Bcrypt密码匹配检查
func BcryptCompare(originStr, hashStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(originStr))
	return err == nil
}
