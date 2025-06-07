package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/ip2region"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/crypto"
)

// repeatParam 重复提交参数的类型定义
type repeatParam struct {
	Time   int64  `json:"time"`
	Params string `json:"params"`
}

// RepeatSubmit 防止表单重复提交，小于间隔时间视为重复提交
//
// 间隔时间(单位秒) 默认:5
//
// 注意之后JSON反序列使用：c.ShouldBindBodyWithJSON(&params)
func RepeatSubmit(interval int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if interval < 5 {
			interval = 5
		}

		// 提交参数
		params := reqctx.RequestParamsMap(c)
		paramsJSONByte, err := json.Marshal(params)
		if err != nil {
			logger.Errorf("repeat submit json marshal err: %v", err)
		}
		paramsJSONStr := string(paramsJSONByte)

		// 获取执行函数名称
		funcName := c.HandlerName()
		// 唯一标识（指定key + 客户端IP + 请求地址）
		clientIP := ip2region.ClientIP(c.ClientIP())
		funcMd5 := crypto.MD5(fmt.Sprintf("%s:%s", clientIP, funcName))
		repeatKey := constants.CACHE_REPEAT_SUBMIT + ":" + funcMd5

		// 在Redis查询并记录请求次数
		repeatStr, _ := redis.Get("", repeatKey)
		if repeatStr != "" {
			var rp repeatParam
			err := json.Unmarshal([]byte(repeatStr), &rp)
			if err != nil {
				logger.Errorf("repeat submit json unmarshal err: %v", err)
			}
			compareTime := time.Now().Unix() - rp.Time
			compareParams := rp.Params == paramsJSONStr

			// 设置重复提交声明响应头（毫秒）
			t := time.Now().Add(time.Duration(compareTime) * time.Second)
			c.Header("X-RepeatSubmit-Rest", fmt.Sprint(t.UnixMilli()))

			// 小于间隔时间且参数内容一致
			if compareTime < interval && compareParams {
				c.JSON(200, resp.ErrMsg("不允许重复提交，请稍候再试"))
				c.Abort()
				return
			}
		}

		// 当前请求参数
		rp := repeatParam{
			Time:   time.Now().Unix(),
			Params: paramsJSONStr,
		}
		rpJSON, err := json.Marshal(rp)
		if err != nil {
			logger.Errorf("repeat submit json marshal err: %v", err)
		}
		// 保存请求时间和参数
		_ = redis.Set("", repeatKey, string(rpJSON), time.Duration(interval)*time.Second)

		// 调用下一个处理程序
		c.Next()
	}
}
