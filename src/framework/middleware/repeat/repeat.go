package repeat

import (
	"encoding/json"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/ip2region"
	"mask_api_gin/src/framework/vo/result"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
// 注意之后JSON反序列使用：c.ShouldBindBodyWith(&params, binding.JSON)
func RepeatSubmit(interval int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if interval < 5 {
			interval = 5
		}

		// 提交参数
		params := ctx.RequestParamsMap(c)
		paramsJSONByte, err := json.Marshal(params)
		if err != nil {
			logger.Errorf("RepeatSubmit params json marshal err: %v", err)
		}
		paramsJSONStr := string(paramsJSONByte)

		// 唯一标识（指定key + 客户端IP + 请求地址）
		clientIP := ip2region.ClientIP(c.ClientIP())
		repeatKey := cachekey.REPEAT_SUBMIT_KEY + clientIP + ":" + c.Request.RequestURI

		// 在Redis查询并记录请求次数
		repeatStr, _ := redis.Get("", repeatKey)
		if repeatStr != "" {
			var rp repeatParam
			err := json.Unmarshal([]byte(repeatStr), &rp)
			if err != nil {
				logger.Errorf("RepeatSubmit repeatStr json unmarshal err: %v", err)
			}
			compareTime := time.Now().Unix() - rp.Time
			compareParams := rp.Params == paramsJSONStr

			// 设置重复提交声明响应头（毫秒）
			c.Header("X-RepeatSubmit-Rest", strconv.FormatInt(time.Now().Add(time.Duration(compareTime)*time.Second).UnixNano()/int64(time.Millisecond), 10))

			// 小于间隔时间且参数内容一致
			if compareTime < interval && compareParams {
				c.JSON(200, result.ErrMsg("不允许重复提交，请稍候再试"))
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
			logger.Errorf("RepeatSubmit rp json marshal err: %v", err)
		}
		// 保存请求时间和参数
		redis.SetByExpire("", repeatKey, string(rpJSON), time.Duration(interval)*time.Second)

		// 调用下一个处理程序
		c.Next()
	}
}
