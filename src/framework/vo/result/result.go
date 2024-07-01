package result

import (
	constResult "mask_api_gin/src/framework/constants/result"
)

// CodeMsg 响应结果
func CodeMsg(code int, msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = code
	args["msg"] = msg
	return args
}

// Ok 响应成功结果
func Ok(v map[string]any) map[string]any {
	args := make(map[string]any)
	args["code"] = constResult.CodeSuccess
	args["msg"] = constResult.MsgSuccess
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// OkMsg 响应成功结果信息
func OkMsg(msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = constResult.CodeSuccess
	args["msg"] = msg
	return args
}

// OkData 响应成功结果数据
func OkData(data any) map[string]any {
	args := make(map[string]any)
	args["code"] = constResult.CodeSuccess
	args["msg"] = constResult.MsgSuccess
	args["data"] = data
	return args
}

// Err 响应失败结果 map[string]any{}
func Err(v map[string]any) map[string]any {
	args := make(map[string]any)
	args["code"] = constResult.CodeError
	args["msg"] = constResult.MsgError
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// ErrMsg 响应失败结果信息
func ErrMsg(msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = constResult.CodeError
	args["msg"] = msg
	return args
}

// ErrData 响应失败结果数据
func ErrData(data any) map[string]any {
	args := make(map[string]any)
	args["code"] = constResult.CodeError
	args["msg"] = constResult.MsgError
	args["data"] = data
	return args
}
