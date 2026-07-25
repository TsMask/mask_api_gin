package file

import (
	"fmt"
	"path/filepath"
	"strings"

	"mask_api_gin/src/internal/config"
)

// UploadResourceDir 读取上传目录配置 prefix, dir
func UploadResourceDir() (string, string) {
	upload := config.Get("staticfile.upload").(map[string]any)
	dir, _ := filepath.Abs(upload["dir"].(string))
	return upload["prefix"].(string), dir
}

// UploadWhiteList 扩展名白名单
func UploadWhiteList() []string {
	arr, _ := config.Get("upload.whitelist").([]any)
	list := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			list = append(list, s)
		}
	}
	return list
}

// UploadFileSize 最大上传文件大小
func UploadFileSize() int {
	size, ok := config.Get("upload.filesize").(int)
	if !ok || size < 1 {
		size = 1
	}
	return size
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
	for _, allow := range UploadWhiteList() {
		if ext == allow {
			return nil
		}
	}
	return fmt.Errorf("限定下载的文件规则：%s", ext)
}
