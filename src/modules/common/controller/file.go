package controller

import (
	"encoding/base64"
	"fmt"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/file"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
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
func (s *FileController) Download(c *gin.Context) {
	filePath := c.Param("filePath")
	if len(filePath) < 8 {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}
	// base64解析出地址
	decodedBytes, err := base64.StdEncoding.DecodeString(filePath)
	if err != nil {
		c.JSON(400, response.CodeMsg(400, err.Error()))
		return
	}
	routerPath := string(decodedBytes)

	// 断点续传
	headerRange := c.GetHeader("Range")
	resultMap, err := file.ReadUploadFileStream(routerPath, headerRange)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	// 响应头
	c.Writer.Header().Set("Content-Disposition", `attachment; filename="`+url.QueryEscape(filepath.Base(routerPath))+`"`)
	c.Writer.Header().Set("Accept-Ranges", "bytes")
	c.Writer.Header().Set("Content-Type", "application/octet-stream")

	if headerRange != "" {
		c.Writer.Header().Set("Content-Range", fmt.Sprint(resultMap["range"]))
		c.Writer.Header().Set("Content-Length", fmt.Sprint(resultMap["chunkSize"]))
		c.Status(206)
	} else {
		c.Writer.Header().Set("Content-Length", fmt.Sprint(resultMap["fileSize"]))
		c.Status(200)

	}
	_, _ = c.Writer.Write(resultMap["data"].([]byte))
}

// Upload 上传文件
//
// POST /upload
func (s *FileController) Upload(c *gin.Context) {
	// 上传的文件
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}
	// 子路径需要在指定范围内
	subPath := c.PostForm("subPath")
	if _, ok := constants.UPLOAD_SUB_PATH[subPath]; subPath != "" && !ok {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}
	if subPath == "" {
		subPath = constants.UPLOAD_COMMON
	}

	// 上传文件转存
	uploadFilePath, err := file.TransferUploadFile(formFile, subPath, nil)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	c.JSON(200, response.OkData(map[string]string{
		"url":              "http://" + c.Request.Host + uploadFilePath,
		"filePath":         uploadFilePath,
		"newFileName":      filepath.Base(uploadFilePath),
		"originalFileName": formFile.Filename,
	}))
}

// ChunkCheck 切片文件检查
//
// POST /chunk-check
func (s *FileController) ChunkCheck(c *gin.Context) {
	var body struct {
		Identifier string `json:"identifier" binding:"required"` // 唯一标识
		FileName   string `json:"fileName" binding:"required"`   // 文件名
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 读取标识目录
	chunks, err := file.ChunkCheckFile(body.Identifier, body.FileName)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, response.OkData(chunks))
}

// ChunkMerge 切片文件合并
//
// POST /chunk-merge
func (s *FileController) ChunkMerge(c *gin.Context) {
	var body struct {
		Identifier string `json:"identifier" binding:"required"` // 唯一标识
		FileName   string `json:"fileName" binding:"required"`   // 文件名
		SubPath    string `json:"subPath"`                       // 子路径类型
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}
	// 子路径需要在指定范围内
	if _, ok := constants.UPLOAD_SUB_PATH[body.SubPath]; body.SubPath != "" && !ok {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}
	if body.SubPath == "" {
		body.SubPath = constants.UPLOAD_COMMON
	}

	// 切片文件合并
	mergeFilePath, err := file.ChunkMergeFile(body.Identifier, body.FileName, body.SubPath)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	c.JSON(200, response.OkData(map[string]string{
		"url":              "http://" + c.Request.Host + mergeFilePath,
		"filePath":         mergeFilePath,
		"newFileName":      filepath.Base(mergeFilePath),
		"originalFileName": body.FileName,
	}))
}

// ChunkUpload 切片文件上传
//
// POST /chunk-upload
func (s *FileController) ChunkUpload(c *gin.Context) {
	// 切片编号
	index := c.PostForm("index")
	// 切片唯一标识
	identifier := c.PostForm("identifier")
	// 上传的文件
	formFile, err := c.FormFile("file")
	if index == "" || identifier == "" || err != nil {
		c.JSON(400, response.CodeMsg(40010, "params error"))
		return
	}

	// 上传文件转存
	chunkFilePath, err := file.TransferChunkUploadFile(formFile, index, identifier)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	c.JSON(206, response.OkData(chunkFilePath))
}
