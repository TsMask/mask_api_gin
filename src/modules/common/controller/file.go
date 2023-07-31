package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 FileController 结构体
var NewFile = &FileController{}

// 文件操作处理
//
// PATH /
type FileController struct{}

// 下载文件
//
// GET /download/:filePath
func (s *FileController) Download(c *gin.Context) {
	filePath := c.Param("filePath")
	c.String(200, filePath)
}

// 上传文件
//
// POST /upload
func (s *FileController) Upload(c *gin.Context) {
	// 单文件
	FileController, _ := c.FormFile("FileController")
	log.Println(FileController.Filename)

	dst := "./" + FileController.Filename
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(FileController, dst)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", FileController.Filename))
}

// 切片文件检查
//
// POST /chunkCheck
func (s *FileController) ChunkCheck(c *gin.Context) {
	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	name, ok := jsonData["identifier"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid or missing 'name' field"})
		return
	}

	email, ok := jsonData["fileName"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid or missing 'email' field"})
		return
	}

	fmt.Println(name)
	fmt.Println(email)

	c.JSON(200, gin.H{
		"message": "User created",
		"name":    name,
		"email":   email,
	})
}

// 切片文件合并
//
// POST /chunkMerge
func (s *FileController) ChunkMerge(c *gin.Context) {
	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	name, ok := jsonData["identifier"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid or missing 'name' field"})
		return
	}

	email, ok := jsonData["fileName"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid or missing 'email' field"})
		return
	}

	fmt.Println(name)
	fmt.Println(email)

	c.JSON(200, gin.H{
		"message": "User created",
		"name":    name,
		"email":   email,
	})
}

// 切片文件上传
//
// POST /chunkUpload
func (s *FileController) ChunkUpload(c *gin.Context) {
	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	name, ok := jsonData["identifier"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid or missing 'name' field"})
		return
	}

	email, ok := jsonData["fileName"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid or missing 'email' field"})
		return
	}

	fmt.Println(name)
	fmt.Println(email)

	c.JSON(200, gin.H{
		"message": "User created",
		"name":    name,
		"email":   email,
	})
}
