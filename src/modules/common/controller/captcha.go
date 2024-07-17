package controller

import (
	"mask_api_gin/src/framework/config"
	constCachekey "mask_api_gin/src/framework/constants/cache_key"
	constCaptcha "mask_api_gin/src/framework/constants/captcha"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	systemService "mask_api_gin/src/modules/system/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// NewCaptcha 实例化控制层
var NewCaptcha = &CaptchaController{
	sysConfigService: systemService.NewSysConfig,
}

// CaptchaController 验证码操作 控制层处理
//
// PATH /
type CaptchaController struct {
	sysConfigService systemService.ISysConfigService // 参数配置服务
}

// Image 获取验证码
//
// GET /captchaImage
func (s *CaptchaController) Image(c *gin.Context) {
	// 从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captchaEnabled")
	captchaEnabled := parse.Boolean(captchaEnabledStr)
	if !captchaEnabled {
		c.JSON(200, result.Ok(map[string]any{
			"captchaEnabled": captchaEnabled,
		}))
		return
	}

	// 生成唯一标识
	verifyKey := ""
	data := map[string]any{
		"captchaEnabled": captchaEnabled,
		"uuid":           "",
		"img":            "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7",
	}

	// 从数据库配置获取验证码类型 math 数值计算 char 字符验证
	captchaType := s.sysConfigService.FindValueByKey("sys.account.captchaType")
	if captchaType == constCaptcha.TypeMath {
		math := config.Get("mathCaptcha").(map[string]any)
		driverCaptcha := &base64Captcha.DriverMath{
			//Height png height in pixel.
			Height: math["height"].(int),
			// Width Captcha png width in pixel.
			Width: math["width"].(int),
			//NoiseCount text noise count.
			NoiseCount: math["noise"].(int),
			//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
			ShowLineOptions: base64Captcha.OptionShowHollowLine,
		}
		if math["color"].(bool) {
			//BgColor captcha image background color (optional)
			driverCaptcha.BgColor = parse.Color(math["background"].(string))
		}
		// 验证码生成
		id, question, answer := driverCaptcha.GenerateIdQuestionAnswer()
		// 验证码表达式解析输出
		item, err := driverCaptcha.DrawCaptcha(question)
		if err != nil {
			logger.Infof("Generate Id Question Answer %s %s : %v", captchaType, question, err)
		} else {
			data["uuid"] = id
			data["img"] = item.EncodeB64string()
			expiration := constCaptcha.Expiration * time.Second
			verifyKey = constCachekey.CaptchaCodeKey + id
			_ = redis.SetByExpire("", verifyKey, answer, expiration)
		}
	}
	if captchaType == constCaptcha.TypeChar {
		char := config.Get("charCaptcha").(map[string]any)
		driverCaptcha := &base64Captcha.DriverString{
			//Height png height in pixel.
			Height: char["height"].(int),
			// Width Captcha png width in pixel.
			Width: char["width"].(int),
			//NoiseCount text noise count.
			NoiseCount: char["noise"].(int),
			//Length random string length.
			Length: char["size"].(int),
			//Source is a unicode which is the rand string from.
			Source: char["chars"].(string),
			//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
			ShowLineOptions: base64Captcha.OptionShowHollowLine,
		}
		if char["color"].(bool) {
			//BgColor captcha image background color (optional)
			driverCaptcha.BgColor = parse.Color(char["background"].(string))
		}
		// 验证码生成
		id, question, answer := driverCaptcha.GenerateIdQuestionAnswer()
		// 验证码表达式解析输出
		item, err := driverCaptcha.DrawCaptcha(question)
		if err != nil {
			logger.Infof("Generate Id Question Answer %s %s : %v", captchaType, question, err)
		} else {
			data["uuid"] = id
			data["img"] = item.EncodeB64string()
			expiration := constCaptcha.Expiration * time.Second
			verifyKey = constCachekey.CaptchaCodeKey + id
			_ = redis.SetByExpire("", verifyKey, answer, expiration)
		}
	}

	// 本地开发下返回验证码结果，方便接口调试
	if config.Env() == "local" {
		text, _ := redis.Get("", verifyKey)
		data["text"] = text
		c.JSON(200, result.Ok(data))
		return
	}
	c.JSON(200, result.Ok(data))
}
