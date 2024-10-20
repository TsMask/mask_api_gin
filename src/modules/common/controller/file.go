package controller

import (
	"encoding/base64"
	"fmt"
	"mask_api_gin/src/framework/constants/uploadsubpath"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/vo/result"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 FileController 结构体
var NewFile = &FileController{}

// 文件操作处理
//
// PATH /file
type FileController struct{}

// 下载文件
//
// GET /download/:filePath
func (s *FileController) Download(c *gin.Context) {
	filePath := c.Param("filePath")
	if len(filePath) < 8 {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// base64解析出地址
	decodedBytes, err := base64.StdEncoding.DecodeString(filePath)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, err.Error()))
		return
	}
	routerPath := string(decodedBytes)
	// 地址文件名截取
	fileName := routerPath[strings.LastIndex(routerPath, "/")+1:]

	// 响应头
	c.Writer.Header().Set("Content-Disposition", `attachment; filename="`+url.QueryEscape(fileName)+`"`)
	c.Writer.Header().Set("Accept-Ranges", "bytes")
	c.Writer.Header().Set("Content-Type", "application/octet-stream")

	// 断点续传
	headerRange := c.GetHeader("Range")
	resultMap, err := file.ReadUploadFileStream(routerPath, headerRange)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	if headerRange != "" {
		c.Writer.Header().Set("Content-Range", fmt.Sprint(resultMap["range"]))
		c.Writer.Header().Set("Content-Length", fmt.Sprint(resultMap["chunkSize"]))
		c.Status(206)
	} else {
		c.Writer.Header().Set("Content-Length", fmt.Sprint(resultMap["fileSize"]))
		c.Status(200)

	}
	c.Writer.Write(resultMap["data"].([]byte))
}

// 上传文件
//
// POST /upload
func (s *FileController) Upload(c *gin.Context) {
	// 上传的文件
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 子路径需要在指定范围内
	subPath := c.PostForm("subPath")
	if _, ok := uploadsubpath.UploadSubpath[subPath]; subPath != "" && !ok {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	if subPath == "" {
		subPath = uploadsubpath.COMMON
	}

	// 上传文件转存
	upFilePath, err := file.TransferUploadFile(formFile, subPath, nil)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	newFileName := upFilePath[strings.LastIndex(upFilePath, "/")+1:]
	c.JSON(200, result.OkData(map[string]string{
		"url":              "http://" + c.Request.Host + upFilePath,
		"fileName":         upFilePath,
		"newFileName":      newFileName,
		"originalFileName": formFile.Filename,
	}))
}

// 切片文件检查
//
// POST /chunkCheck
func (s *FileController) ChunkCheck(c *gin.Context) {
	var body struct {
		// 唯一标识
		Identifier string `json:"identifier" binding:"required"`
		// 文件名
		FileName string `json:"fileName" binding:"required"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 读取标识目录
	chunks, err := file.ChunkCheckFile(body.Identifier, body.FileName)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.OkData(chunks))
}

// 切片文件合并
//
// POST /chunkMerge
func (s *FileController) ChunkMerge(c *gin.Context) {
	var body struct {
		// 唯一标识
		Identifier string `json:"identifier" binding:"required"`
		// 文件名
		FileName string `json:"fileName" binding:"required"`
		// 子路径类型
		SubPath string `json:"subPath"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 子路径需要在指定范围内
	if _, ok := uploadsubpath.UploadSubpath[body.SubPath]; body.SubPath != "" && !ok {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	if body.SubPath == "" {
		body.SubPath = uploadsubpath.COMMON
	}

	// 切片文件合并
	mergeFilePath, err := file.ChunkMergeFile(body.Identifier, body.FileName, body.SubPath)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	newFileName := mergeFilePath[strings.LastIndex(mergeFilePath, "/")+1:]
	c.JSON(200, result.OkData(map[string]string{
		"url":              "http://" + c.Request.Host + mergeFilePath,
		"fileName":         mergeFilePath,
		"newFileName":      newFileName,
		"originalFileName": body.FileName,
	}))
}

// 切片文件上传
//
// POST /chunkUpload
func (s *FileController) ChunkUpload(c *gin.Context) {
	// 切片编号
	index := c.PostForm("index")
	// 切片唯一标识
	identifier := c.PostForm("identifier")
	// 上传的文件
	formFile, err := c.FormFile("file")
	if index == "" || identifier == "" || err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 上传文件转存
	chunkFilePath, err := file.TransferChunkUploadFile(formFile, index, identifier)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(206, result.OkData(chunkFilePath))
}
