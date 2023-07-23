package uploadsubpath

// 文件上传-子路径类型常量

const (
	// 默认
	DEFAULT = "default"

	// 头像
	AVATART = "avatar"

	// 导入
	IMPORT = "import"

	// 导出
	EXPORT = "export"

	// 通用上传
	COMMON = "common"

	// 下载
	DOWNLOAD = "download"

	// 切片
	CHUNK = "chunk"
)

// 子路径类型映射
var UploadSubpath = map[string]string{
	DEFAULT:  "默认",
	AVATART:  "头像",
	IMPORT:   "导入",
	EXPORT:   "导出",
	COMMON:   "通用上传",
	DOWNLOAD: "下载",
	CHUNK:    "切片",
}
