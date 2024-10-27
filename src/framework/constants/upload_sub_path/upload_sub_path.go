package upload_sub_path

// 文件上传-子路径类型常量

const (
	// DEFAULT 默认
	DEFAULT = "default"

	// AVATAR 头像
	AVATAR = "avatar"

	//IMPORT 导入
	IMPORT = "import"

	// EXPORT 导出
	EXPORT = "export"

	// Common 通用上传
	COMMON = "common"

	// DOWNLOAD 下载
	DOWNLOAD = "download"

	// CHUNK 切片
	CHUNK = "chunk"
)

// UploadSubPath 子路径类型映射
var UploadSubPath = map[string]string{
	DEFAULT:  "默认",
	AVATAR:   "头像",
	IMPORT:   "导入",
	EXPORT:   "导出",
	COMMON:   "通用上传",
	DOWNLOAD: "下载",
	CHUNK:    "切片",
}
