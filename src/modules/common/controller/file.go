package controller

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
)

// NewFile 实例化控制层
var NewFile = &FileController{}

// FileController 文件操作 控制层处理
//
// PATH /file
type FileController struct{}

// Download 下载文件
//
// GET /download/:filePath
func (s FileController) Download(c *gin.Context) {
	filePath := c.Param("filePath")
	if len(filePath) < 8 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: filePath not is base64 string"))
		return
	}
	// base64解析出地址
	decodedBytes, err := base64.StdEncoding.DecodeString(filePath)
	if err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "decode err: filePath not is base64 string"))
		return
	}
	routerPath := string(decodedBytes)
	if len(routerPath) > 500 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "file path too long"))
		return
	}

	if err := file.UploadAllowReadFilePath(routerPath); err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	absPath := file.UploadFileAbsPath(routerPath)
	headerRange := c.GetHeader("Range") // 断点续传

	var result file.FileRange
	if headerRange != "" {
		partsStr := strings.Replace(headerRange, "bytes=", "", 1)
		parts := strings.Split(partsStr, "-")
		start, _ := strconv.ParseInt(parts[0], 10, 64)
		end, _ := strconv.ParseInt(parts[1], 10, 64)
		result, err = file.ReadStream(absPath, start, end)
	} else {
		result, err = file.ReadStream(absPath, 0, 0)
	}
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	// 响应头
	c.Writer.Header().Set("Accept-Ranges", "bytes")
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", `attachment; filename="`+url.QueryEscape(filepath.Base(routerPath))+`"`)

	if headerRange != "" {
		c.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", result.Start, result.End, result.Total))
		c.Writer.Header().Set("Content-Length", strconv.Itoa(len(result.Data)))
		c.Status(206)
	} else {
		c.Writer.Header().Set("Content-Length", strconv.FormatInt(result.Total, 10))
		c.Status(200)
	}
	_, _ = c.Writer.Write(result.Data)
}

// Upload 上传文件
//
// POST /upload
func (s FileController) Upload(c *gin.Context) {
	// 上传的文件
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: file is empty"))
		return
	}
	// 子路径需要在指定范围内
	subPath := c.PostForm("subPath")
	if _, ok := constants.UPLOAD_SUB_PATH[subPath]; subPath != "" && !ok {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: subPath not in range"))
		return
	}
	if subPath == "" {
		subPath = constants.UPLOAD_COMMON
	}

	src, err := formFile.Open()
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	defer src.Close()

	prefix, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()

	dst, err := file.FileUpload{
		Name:   formFile.Filename,
		Reader: src,
	}.Save(file.FileConfig{
		Dir:       filepath.Join(dir, subPath, datePath),
		MaxSize:   file.UploadFileSize(),
		WhiteExts: file.UploadWhiteList(),
	})
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	urlPath := strings.Replace(dst, dir, prefix, 1)
	urlPath = filepath.ToSlash(urlPath)

	c.JSON(200, resp.OkData(map[string]string{
		"url":              "//" + c.Request.Host + urlPath,
		"filePath":         urlPath,
		"newFileName":      filepath.Base(dst),
		"originalFileName": formFile.Filename,
	}))
}

// ChunkCheck 切片文件检查
//
// POST /chunk-check
func (s FileController) ChunkCheck(c *gin.Context) {
	var body struct {
		Identifier string `json:"identifier" binding:"required"` // 唯一标识
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, fmt.Sprintf("bind err: %s", resp.FormatBindError(err))))
		return
	}

	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	chunkDir := filepath.Join(dir, constants.UPLOAD_CHUNK, datePath)
	// 读取标识目录
	chunks, err := file.FileConfig{Dir: chunkDir}.ChunkList(body.Identifier)
	if err != nil {
		chunks = []string{}
	}
	c.JSON(200, resp.OkData(chunks))
}

// ChunkMerge 切片文件合并
//
// POST /chunk-merge
func (s FileController) ChunkMerge(c *gin.Context) {
	var body struct {
		Identifier string `json:"identifier" binding:"required"` // 唯一标识
		FileName   string `json:"fileName" binding:"required"`   // 文件名
		SubPath    string `json:"subPath"`                       // 子路径类型
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, fmt.Sprintf("bind err: %s", resp.FormatBindError(err))))
		return
	}
	// 子路径需要在指定范围内
	if _, ok := constants.UPLOAD_SUB_PATH[body.SubPath]; body.SubPath != "" && !ok {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: subPath not in range"))
		return
	}
	if body.SubPath == "" {
		body.SubPath = constants.UPLOAD_COMMON
	}

	prefix, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()

	// 合并到 chunk 日期目录下
	dst, err := file.FileConfig{
		Dir:       filepath.Join(dir, constants.UPLOAD_CHUNK, datePath),
		MaxSize:   file.UploadFileSize(),
		WhiteExts: file.UploadWhiteList(),
	}.ChunkMerge(
		body.Identifier,
		body.FileName,
	)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	// 移动到目标子路径目录
	finalDir := filepath.Join(dir, body.SubPath, datePath)
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	finalPath := filepath.Join(finalDir, filepath.Base(dst))
	if err := os.Rename(dst, finalPath); err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	urlPath := strings.Replace(finalPath, dir, prefix, 1)
	urlPath = filepath.ToSlash(urlPath)

	c.JSON(200, resp.OkData(map[string]string{
		"url":              "//" + c.Request.Host + urlPath,
		"filePath":         urlPath,
		"newFileName":      filepath.Base(finalPath),
		"originalFileName": body.FileName,
	}))
}

// ChunkUpload 切片文件上传
//
// POST /chunk-upload
func (s FileController) ChunkUpload(c *gin.Context) {
	index := c.PostForm("index")           // 切片编号
	identifier := c.PostForm("identifier") // 切片唯一标识
	if index == "" || identifier == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: index and identifier must be set"))
		return
	}
	// 上传的文件
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: file is empty"))
		return
	}

	src, err := formFile.Open()
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	defer src.Close()

	prefix, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	chunkDir := filepath.Join(dir, constants.UPLOAD_CHUNK, datePath)

	upload := file.FileUpload{
		Reader: src,
		Name:   formFile.Filename,
		Id:     identifier,
		Index:  index,
	}
	cfg := file.FileConfig{
		Dir:     chunkDir,
		MaxSize: file.UploadFileSize(),
	}
	if err := upload.ChunkSave(cfg); err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	urlPath := strings.Replace(chunkDir, dir, prefix, 1)
	urlPath = filepath.ToSlash(filepath.Join(urlPath, identifier, index))
	c.JSON(206, resp.OkData(urlPath))
}
