package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id 参数（取 OWASP / RFC 9106 推荐区间，后续需要可提升为配置）
const (
	argonVersion     = argon2.Version // 19
	argonMemory      = 64 * 1024      // KiB，约 64MB
	argonIterations  = 3
	argonParallelism = 2
	argonSaltLen     = 16
	argonKeyLen      = 32
)

// Argon2idHash 使用 Argon2id 生成自描述哈希串。
// 输出 PHC 格式：$argon2id$v=19$m=65536,t=3,p=2$<盐>$<哈希>
func Argon2idHash(originStr string) string {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return ""
	}
	key := argon2.IDKey([]byte(originStr), salt, argonIterations, argonMemory, argonParallelism, argonKeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argonVersion, argonMemory, argonIterations, argonParallelism, b64Salt, b64Key)
}

// Argon2idCompare 校验明文与 Argon2id 哈希是否匹配（常数时间比较）。
func Argon2idCompare(originStr, hashStr string) bool {
	p, salt, key, err := decodeHash(hashStr)
	if err != nil {
		return false
	}
	otherKey := argon2.IDKey([]byte(originStr), salt, p.iter, p.mem, p.par, uint32(len(key)))
	return subtle.ConstantTimeCompare(key, otherKey) == 1
}

type argonParams struct {
	mem  uint32
	iter uint32
	par  uint8
}

// decodeHash 解析 PHC 格式 argon2id 哈希串。
// 期望结构：$argon2id$v=19$m=65536,t=3,p=2$<b64盐>$<b64哈希>
func decodeHash(hashStr string) (argonParams, []byte, []byte, error) {
	parts := strings.Split(hashStr, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return argonParams{}, nil, nil, fmt.Errorf("invalid argon2id hash")
	}
	p := argonParams{}
	var version int
	for _, seg := range parts[2:4] {
		for kv := range strings.SplitSeq(seg, ",") {
			kvParts := strings.SplitN(kv, "=", 2)
			if len(kvParts) != 2 {
				return argonParams{}, nil, nil, fmt.Errorf("invalid argon2id params")
			}
			switch kvParts[0] {
			case "v":
				fmt.Sscanf(kvParts[1], "%d", &version)
			case "m":
				fmt.Sscanf(kvParts[1], "%d", &p.mem)
			case "t":
				fmt.Sscanf(kvParts[1], "%d", &p.iter)
			case "p":
				var pp int
				fmt.Sscanf(kvParts[1], "%d", &pp)
				p.par = uint8(pp)
			}
		}
	}
	if version != argonVersion {
		return argonParams{}, nil, nil, fmt.Errorf("unsupported argon2id version %d", version)
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return argonParams{}, nil, nil, err
	}
	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return argonParams{}, nil, nil, err
	}
	return p, salt, key, nil
}
