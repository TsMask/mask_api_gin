package crypto

import (
	"crypto/sha256"
	"encoding/base64"
)

// SHA256ToBase64 编码字符串
func SHA256ToBase64(str string) string {
	hash := sha256.Sum256([]byte(str))
	return base64.URLEncoding.EncodeToString(hash[:])
}
