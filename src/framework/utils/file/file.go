package file

import (
	"errors"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants/uploadsubpath"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/generate"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mime/multipart"
	"path"
	"path/filepath"
	"strconv"
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
// fileName 原始文件名称含后缀，如：midway1_logo_iipc68.png
//
// allowExts 允许上传拓展类型，['.png']
func isAllowWrite(fileName string, allowExts []string, fileSize int64) error {
	// 判断上传文件名称长度
	if len(fileName) > DEFAULT_FILE_NAME_LENGTH {
		return fmt.Errorf("上传文件名称长度限制最长为 %d", DEFAULT_FILE_NAME_LENGTH)
	}

	// 最大上传文件大小
	maxFileSize := uploadFileSize()
	if fileSize > maxFileSize {
		return fmt.Errorf("最大上传文件大小 %s", parse.Bit(float64(maxFileSize)))
	}

	// 判断文件拓展是否为允许的拓展类型
	fileExt := filepath.Ext(fileName)
	hasExt := false
	if len(allowExts) == 0 {
		allowExts = uploadWhiteList()
	}
	for _, ext := range allowExts {
		if ext == fileExt {
			hasExt = true
			break
		}
	}
	if !hasExt {
		return fmt.Errorf("上传文件类型不支持，仅支持以下类型：%s", strings.Join(allowExts, "、"))
	}

	return nil
}

// 检查文件允许本地读取
//
// filePath 文件存放资源路径，URL相对地址
func isAllowRead(filePath string) error {
	// 禁止目录上跳级别
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("禁止目录上跳级别")
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
		return fmt.Errorf("非法下载的文件规则：%s", fileExt)
	}

	return nil
}

// TransferUploadFile 上传资源文件转存
//
// subPath 子路径，默认 UploadSubPath.DEFAULT
//
// allowExts 允许上传拓展类型（含“.”)，如 ['.png','.jpg']
func TransferUploadFile(file *multipart.FileHeader, subPath string, allowExts []string) (string, error) {
	// 上传文件检查
	err := isAllowWrite(file.Filename, allowExts, file.Size)
	if err != nil {
		return "", err
	}
	// 上传资源路径
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

// ReadUploadFileStream 上传资源文件读取
//
// filePath 文件存放资源路径，URL相对地址 如：/upload/common/2023/06/xxx.png
//
// headerRange 断点续传范围区间，bytes=0-12131
func ReadUploadFileStream(filePath, headerRange string) (map[string]interface{}, error) {
	// 读取文件检查
	err := isAllowRead(filePath)
	if err != nil {
		return map[string]interface{}{}, err
	}
	// 上传资源路径
	prefix, dir := resourceUpload()
	fileAsbPath := strings.Replace(filePath, prefix, dir, 1)

	// 响应结果
	result := map[string]interface{}{
		"range":     "",
		"chunkSize": 0,
		"fileSize":  0,
		"data":      nil,
	}

	// 文件大小
	fileSize := getFileSize(fileAsbPath)
	if fileSize <= 0 {
		return result, nil
	}
	result["fileSize"] = fileSize

	if headerRange != "" {
		partsStr := strings.Replace(headerRange, "bytes=", "", 1)
		parts := strings.Split(partsStr, "-")
		start, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil || start > fileSize {
			start = 0
		}
		end, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end > fileSize {
			end = fileSize - 1
		}
		if start > end {
			start = end
		}

		// 分片结果
		result["range"] = fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize)
		result["chunkSize"] = end - start + 1
		byteArr, err := getFileStream(fileAsbPath, start, end)
		if err != nil {
			return map[string]interface{}{}, errors.New("读取文件失败")
		}
		result["data"] = byteArr
		return result, nil
	}

	byteArr, err := getFileStream(fileAsbPath, 0, fileSize)
	if err != nil {
		return map[string]interface{}{}, errors.New("读取文件失败")
	}
	result["data"] = byteArr
	return result, nil
}

// TransferChunkUploadFile 上传资源切片文件转存
//
// file 上传文件对象
//
// index 切片文件序号
//
// identifier 切片文件目录标识符
func TransferChunkUploadFile(file *multipart.FileHeader, index, identifier string) (string, error) {
	// 上传文件检查
	err := isAllowWrite(file.Filename, []string{}, file.Size)
	if err != nil {
		return "", err
	}
	// 上传资源路径
	prefix, dir := resourceUpload()
	// 新文件名称并组装文件地址
	filePath := filepath.Join(uploadsubpath.CHUNK, date.ParseDatePath(time.Now()), identifier)
	writePathFile := filepath.Join(dir, filePath, index)
	// 存入新文件路径
	err = transferToNewFile(file, writePathFile)
	if err != nil {
		return "", err
	}
	urlPath := filepath.Join(prefix, filePath, index)
	return filepath.ToSlash(urlPath), nil
}

// 上传资源切片文件检查
//
// identifier 切片文件目录标识符
//
// originalFileName 原始文件名称，如logo.png
func ChunkCheckFile(identifier, originalFileName string) ([]string, error) {
	// 读取文件检查
	err := isAllowWrite(originalFileName, []string{}, 0)
	if err != nil {
		return []string{}, err
	}
	// 上传资源路径
	_, dir := resourceUpload()
	dirPath := path.Join(uploadsubpath.CHUNK, date.ParseDatePath(time.Now()), identifier)
	readPath := path.Join(dir, dirPath)
	fileList, err := getDirFileNameList(readPath)
	if err != nil {
		return []string{}, errors.New("读取文件失败")
	}
	return fileList, nil
}

// 上传资源切片文件检查
//
// identifier 切片文件目录标识符
//
// originalFileName 原始文件名称，如logo.png
//
// subPath 子路径，默认 DEFAULT
func ChunkMergeFile(identifier, originalFileName, subPath string) (string, error) {
	// 读取文件检查
	err := isAllowWrite(originalFileName, []string{}, 0)
	if err != nil {
		return "", err
	}
	// 上传资源路径
	prefix, dir := resourceUpload()
	// 切片存放目录
	dirPath := path.Join(uploadsubpath.CHUNK, date.ParseDatePath(time.Now()), identifier)
	readPath := path.Join(dir, dirPath)
	// 组合存放文件路径
	fileName := generateFileName(originalFileName)
	filePath := path.Join(subPath, date.ParseDatePath(time.Now()))
	writePath := path.Join(dir, filePath)
	err = mergeToNewFile(readPath, writePath, fileName)
	if err != nil {
		return "", err
	}
	urlPath := filepath.Join(prefix, filePath, fileName)
	return filepath.ToSlash(urlPath), nil
}
