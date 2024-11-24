package generate

import (
	"mask_api_gin/src/framework/logger"

	"crypto/rand"
	"math/big"
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
)

// Code 生成随机Code
// 包含数字、小写字母
// 不保证长度满足
func Code(size int) string {
	str, err := nanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", size)
	if err != nil {
		logger.Infof("Code %d : %v", size, err)
		return ""
	}
	return str
}

// String 生成随机字符串
// 包含数字、大小写字母、下划线、横杠
// 不保证长度满足
func String(size int) string {
	str, err := nanoid.New(size)
	if err != nil {
		logger.Errorf("String %d : %v", size, err)
		return ""
	}
	return str
}

// Number 随机数 纯数字0-9
func Number(size int) string {
	if size <= 0 {
		return ""
	}

	const digits = "0123456789"
	var sb strings.Builder
	sb.Grow(size) // 预分配字符串空间

	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			logger.Errorf("Number %d : %v", size, err)
			return ""
		}
		sb.WriteByte(digits[n.Int64()])
	}

	return sb.String()
}
