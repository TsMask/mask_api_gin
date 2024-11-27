package constants

// 文件上传-子路径类型常量
const (
	// UPLOAD_DEFAULT 默认
	UPLOAD_DEFAULT = "default"
	// UPLOAD_AVATAR 头像
	UPLOAD_AVATAR = "avatar"
	// UPLOAD_IMPORT 导入
	UPLOAD_IMPORT = "import"
	// UPLOAD_EXPORT 导出
	UPLOAD_EXPORT = "export"
	// UPLOAD_COMMON 通用上传
	UPLOAD_COMMON = "common"
	// UPLOAD_DOWNLOAD 下载
	UPLOAD_DOWNLOAD = "download"
	// UPLOAD_CHUNK 切片
	UPLOAD_CHUNK = "chunk"
)

// UPLOAD_SUB_PATH 子路径类型映射
var UPLOAD_SUB_PATH = map[string]string{
	UPLOAD_DEFAULT:  "默认",
	UPLOAD_AVATAR:   "头像",
	UPLOAD_IMPORT:   "导入",
	UPLOAD_EXPORT:   "导出",
	UPLOAD_COMMON:   "通用上传",
	UPLOAD_DOWNLOAD: "下载",
	UPLOAD_CHUNK:    "切片",
}
