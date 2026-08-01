package file

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"mask_api_gin/src/internal/config"
)

// UploadResourceDir 读取上传目录配置 prefix, dir
func UploadResourceDir() (string, string) {
	prefix := config.GetStringOr("staticFile.upload.prefix", "")
	dir, _ := filepath.Abs(config.GetStringOr("staticFile.upload.dir", ""))
	return prefix, dir
}

// UploadWhiteList 扩展名白名单
func UploadWhiteList() []string {
	return config.GetStringSliceOr("upload.whitelist", []string{})
}

// UploadFileSize 最大上传文件大小
func UploadFileSize() int {
	if size := config.GetIntOr("upload.filesize", 1); size >= 1 {
		return size
	}
	return 1
}

// UploadFileAbsPath  上传资源本地绝对资源路径
//
// filePath 上传文件路径
func UploadFileAbsPath(filePath string) string {
	prefix, dir := UploadResourceDir()
	absPath := strings.Replace(filePath, prefix, dir, 1)
	return filepath.ToSlash(absPath)
}

// UploadAllowReadFilePath 校验文件路径安全性
func UploadAllowReadFilePath(filePath string) error {
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("禁止目录上跳级别")
	}
	ext := filepath.Ext(filePath)
	if slices.Contains(UploadWhiteList(), ext) {
		return nil
	}
	return fmt.Errorf("限定下载的文件规则：%s", ext)
}
