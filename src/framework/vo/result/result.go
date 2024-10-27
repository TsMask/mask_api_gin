package result

const (
	// CODE_ERROOR 响应-code错误失败
	CODE_ERROOR = -1
	// MSG_ERROR 响应-msg错误失败
	MSG_ERROR = "失败"

	// CODE_SUCCESS 响应-msg正常成功
	CODE_SUCCESS = 0
	// MSG_SUCCCESS 响应-code正常成功
	MSG_SUCCCESS = "成功"
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
	args["code"] = CODE_SUCCESS
	args["msg"] = MSG_SUCCCESS
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// OkMsg 响应成功结果信息
func OkMsg(msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = CODE_SUCCESS
	args["msg"] = msg
	return args
}

// OkData 响应成功结果数据
func OkData(data any) map[string]any {
	args := make(map[string]any)
	args["code"] = CODE_SUCCESS
	args["msg"] = MSG_SUCCCESS
	args["data"] = data
	return args
}

// Err 响应失败结果 map[string]any{}
func Err(v map[string]any) map[string]any {
	args := make(map[string]any)
	args["code"] = CODE_ERROOR
	args["msg"] = MSG_ERROR
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// ErrMsg 响应失败结果信息
func ErrMsg(msg string) map[string]any {
	args := make(map[string]any)
	args["code"] = CODE_ERROOR
	args["msg"] = msg
	return args
}

// ErrData 响应失败结果数据
func ErrData(data any) map[string]any {
	args := make(map[string]any)
	args["code"] = CODE_ERROOR
	args["msg"] = MSG_ERROR
	args["data"] = data
	return args
}
