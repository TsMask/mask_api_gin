package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 文件操作处理
var File = new(file)

type file struct{}

// 下载文件
//
// GET /download/:filePath
func (s *file) Download(c *gin.Context) {
	filePath := c.Param("filePath")
	c.String(200, filePath)
}

// 上传文件
//
// POST /upload
func (s *file) Upload(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	dst := "./" + file.Filename
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(file, dst)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// 切片文件检查
//
// POST /chunkCheck
func (s *file) ChunkCheck(c *gin.Context) {
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
func (s *file) ChunkMerge(c *gin.Context) {
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
func (s *file) ChunkUpload(c *gin.Context) {
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
