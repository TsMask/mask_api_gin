package result

import "mask_api_gin/src/framework/constants"

// 响应成功结果 map[string]interface{}{}
func Ok(v map[string]interface{}) map[string]interface{} {
	args := make(map[string]interface{})
	args["code"] = constants.RESULT_CODE_SUCCESS
	args["msg"] = constants.RESULT_MSG_SUCCESS
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// 响应成功结果信息
func OkMsg(msg string) map[string]interface{} {
	args := make(map[string]interface{})
	args["code"] = constants.RESULT_CODE_SUCCESS
	args["msg"] = msg
	return args
}

// 响应成功结果数据
func OkData(data interface{}) map[string]interface{} {
	args := make(map[string]interface{})
	args["code"] = constants.RESULT_CODE_SUCCESS
	args["msg"] = constants.RESULT_MSG_SUCCESS
	args["data"] = data
	return args
}

// 响应失败结果 map[string]interface{}{}
func Err(v map[string]interface{}) map[string]interface{} {
	args := make(map[string]interface{})
	args["code"] = constants.RESULT_CODE_ERROR
	args["msg"] = constants.RESULT_MSG_ERROR
	// v合并到args
	for key, value := range v {
		args[key] = value
	}
	return args
}

// 响应失败结果信息
func ErrMsg(msg string) map[string]interface{} {
	args := make(map[string]interface{})
	args["code"] = constants.RESULT_CODE_ERROR
	args["msg"] = msg
	return args
}

// 响应失败结果数据
func ErrData(data interface{}) map[string]interface{} {
	args := make(map[string]interface{})
	args["code"] = constants.RESULT_CODE_ERROR
	args["msg"] = constants.RESULT_MSG_ERROR
	args["data"] = data
	return args
}
