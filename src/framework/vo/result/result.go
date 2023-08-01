package result

import (
	"mask_api_gin/src/framework/constants/result"
)

// CodeMsg 响应结果
func CodeMsg(code int, msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = code
	args["msg"] = msg
	return args
}

// 响应成功结果 map[string]any{}
func Ok(v map[string]any) map[string]any {
	args := make(map[string]any)
	args["code"] = result.CODE_SUCCESS
	args["msg"] = result.MSG_SUCCESS
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// 响应成功结果信息
func OkMsg(msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = result.CODE_SUCCESS
	args["msg"] = msg
	return args
}

// 响应成功结果数据
func OkData(data any) map[string]any {
	args := make(map[string]any)
	args["code"] = result.CODE_SUCCESS
	args["msg"] = result.MSG_SUCCESS
	args["data"] = data
	return args
}

// 响应失败结果 map[string]any{}
func Err(v map[string]any) map[string]any {
	args := make(map[string]any)
	args["code"] = result.CODE_ERROR
	args["msg"] = result.MSG_ERROR
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// 响应失败结果信息
func ErrMsg(msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = result.CODE_ERROR
	args["msg"] = msg
	return args
}

// 响应失败结果数据
func ErrData(data any) map[string]any {
	args := make(map[string]any)
	args["code"] = result.CODE_ERROR
	args["msg"] = result.MSG_ERROR
	args["data"] = data
	return args
}
