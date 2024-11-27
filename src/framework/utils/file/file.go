package file

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/date"

	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TransferUploadFile 上传资源文件转存
//
// subPath 子路径，默认 UploadSubPath.DEFAULT
//
// allowExt 允许上传拓展类型（含“.”)，如 ['.png','.jpg']
func TransferUploadFile(file *multipart.FileHeader, subPath string, allowExt []string) (string, error) {
	// 上传文件检查
	err := isAllowWrite(file.Filename, allowExt, file.Size)
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
func ReadUploadFileStream(filePath, headerRange string) (map[string]any, error) {
	// 读取文件检查
	err := isAllowRead(filePath)
	if err != nil {
		return map[string]any{}, err
	}
	// 上传资源路径
	prefix, dir := resourceUpload()
	fileAsbPath := strings.Replace(filePath, prefix, dir, 1)

	// 响应结果
	result := map[string]any{
		"range":     "",
		"chunkSize": 0,
		"fileSize":  0,
		"data":      []byte{},
	}

	// 文件大小
	fileSize := getFileSize(fileAsbPath)
	if fileSize <= 0 {
		return result, fmt.Errorf("文件不存在")
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
			return map[string]any{}, fmt.Errorf("读取文件失败")
		}
		result["data"] = byteArr
		return result, nil
	}

	byteArr, err := getFileStream(fileAsbPath, 0, fileSize)
	if err != nil {
		return map[string]any{}, fmt.Errorf("读取文件失败")
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
	filePath := filepath.Join(constants.UPLOAD_CHUNK, date.ParseDatePath(time.Now()), identifier)
	writePathFile := filepath.Join(dir, filePath, index)
	// 存入新文件路径
	err = transferToNewFile(file, writePathFile)
	if err != nil {
		return "", err
	}
	urlPath := filepath.Join(prefix, filePath, index)
	return filepath.ToSlash(urlPath), nil
}

// ChunkCheckFile 上传资源切片文件检查
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
	dirPath := path.Join(constants.UPLOAD_CHUNK, date.ParseDatePath(time.Now()), identifier)
	readPath := path.Join(dir, dirPath)
	fileList, err := getDirFileNameList(readPath)
	if err != nil {
		return []string{}, fmt.Errorf("读取文件失败")
	}
	return fileList, nil
}

// ChunkMergeFile 上传资源切片文件检查
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
	dirPath := path.Join(constants.UPLOAD_CHUNK, date.ParseDatePath(time.Now()), identifier)
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

// ParseUploadFilePath  上传资源本地绝对资源路径
//
// filePath 上传文件路径
func ParseUploadFilePath(filePath string) string {
	prefix, dir := resourceUpload()
	absPath := strings.Replace(filePath, prefix, dir, 1)
	return filepath.ToSlash(absPath)
}
