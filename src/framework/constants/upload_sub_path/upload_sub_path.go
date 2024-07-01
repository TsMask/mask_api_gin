package upload_sub_path

// 文件上传-子路径类型常量

const (
	// Default 默认
	Default = "default"

	// Avatar 头像
	Avatar = "avatar"

	//Import 导入
	Import = "import"

	// Export 导出
	Export = "export"

	// Common 通用上传
	Common = "common"

	// Download 下载
	Download = "download"

	// Chunk 切片
	Chunk = "chunk"
)

// UploadSubPath 子路径类型映射
var UploadSubPath = map[string]string{
	Default:  "默认",
	Avatar:   "头像",
	Import:   "导入",
	Export:   "导出",
	Common:   "通用上传",
	Download: "下载",
	Chunk:    "切片",
}
