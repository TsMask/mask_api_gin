package file

import (
	"errors"
	"fmt"
	"io"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/generate"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/**最大文件名长度 */
const DEFAULT_FILE_NAME_LENGTH = 100

// 文件上传路径 prefix, dir
func resourceUpload() (string, string) {
	upload := config.Get("staticFile.upload").(map[string]interface{})
	dir, err := filepath.Abs(upload["dir"].(string))
	if err != nil {
		logger.Errorf("file resourceUpload err %v", err)
	}
	return upload["prefix"].(string), dir
}

// 最大上传文件大小
func uploadFileSize() int64 {
	fileSize := 1 * 1024 * 1024
	size := config.Get("upload.fileSize").(int)
	if size > 1 {
		fileSize = size * 1024 * 1024
	}
	return int64(fileSize)
}

// 文件上传扩展名白名单
func uploadWhiteList() []string {
	arr := config.Get("upload.whitelist").([]interface{})
	strings := make([]string, len(arr))
	for i, val := range arr {
		if str, ok := val.(string); ok {
			strings[i] = str
		}
	}
	return strings
}

// 生成文件名称 fileName_随机值.extName
//
// fileName 原始文件名称含后缀，如：logo.png
func generateFileName(fileName string) string {
	fileExt := filepath.Ext(fileName)
	// 替换掉后缀和特殊字符保留文件名
	newFileName := regular.Replace(fileName, fileExt, "")
	newFileName = regular.Replace(newFileName, `[<>:"\\|?*]+`, "")
	newFileName = strings.TrimSpace(newFileName)
	return fmt.Sprintf("%s_%s%s", newFileName, generate.Code(6), fileExt)
}

// 检查文件允许写入本地
//
// allowExts 允许上传拓展类型，['.png']
//
// fileName 原始文件名称含后缀，如：midway1_logo_iipc68.png
func isAllowWrite(fileName string, fileSize int64, allowExts []string) error {
	// 判断上传文件名称长度
	if len(fileName) > DEFAULT_FILE_NAME_LENGTH {
		msg := fmt.Sprintf("上传文件名称长度限制最长为 %d", DEFAULT_FILE_NAME_LENGTH)
		return errors.New(msg)
	}

	// 最大上传文件大小
	maxFileSize := uploadFileSize()
	if fileSize > maxFileSize {
		msg := fmt.Sprintf("最大上传文件大小 %s", parse.Bit(float64(maxFileSize)))
		return errors.New(msg)
	}

	// 判断文件拓展是否为允许的拓展类型
	fileExt := filepath.Ext(fileName)
	hasExt := false
	for _, ext := range allowExts {
		if ext == fileExt {
			hasExt = true
			break
		}
	}
	if !hasExt {
		msg := fmt.Sprintf("上传文件类型不支持，仅支持以下类型：%s", strings.Join(allowExts, "、"))
		return errors.New(msg)
	}

	return nil
}

// 检查文件允许本地读取
//
// filePath 文件存放资源路径，URL相对地址
func isAllowRead(filePath string) error {
	// 禁止目录上跳级别
	if strings.Contains(filePath, "..") {
		return errors.New("禁止目录上跳级别")
	}

	// 检查允许下载的文件规则
	fileExt := filepath.Ext(filePath)
	hasExt := false
	for _, ext := range uploadWhiteList() {
		if ext == fileExt {
			hasExt = true
			break
		}
	}
	if !hasExt {
		msg := fmt.Sprintf("非法下载的文件规则：%s", fileExt)
		return errors.New(msg)
	}

	return nil
}

// transferToNewFile 读取目标文件转移到新路径下
func transferToNewFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// 上传资源文件转存
//
// subPath 子路径，默认 UploadSubPath.DEFAULT
//
// allowExts 允许上传拓展类型（含“.”)，如 ['.png','.jpg']
func TransferUploadFile(file *multipart.FileHeader, subPath string, allowExts []string) (string, error) {
	if len(allowExts) == 0 {
		allowExts = uploadWhiteList()
	}
	// 上传文件检查
	err := isAllowWrite(file.Filename, file.Size, allowExts)
	if err != nil {
		return "", err
	}
	prefix, dir := resourceUpload()
	// 新文件名称并组装文件地址
	fileName := generateFileName(file.Filename)
	filePath := filepath.Join(subPath, date.ParseDatePath(time.Now()))
	writePathFile := filepath.Join(dir, filePath, fileName)
	// 存入新文件路径
	err = transferToNewFile(file, writePathFile)
	if err != nil {
		return "", err
	}
	urlPath := filepath.Join(prefix, filePath, fileName)
	return filepath.ToSlash(urlPath), nil
}

// SaveUploadedFile uploads the form file to specific dst.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	isAllowRead(file.Filename)
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
