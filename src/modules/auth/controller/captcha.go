package controller

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/resp"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewCaptcha 实例化控制层
var NewCaptcha = &CaptchaController{
	sysConfigService: systemService.NewSysConfig,
}

// CaptchaController 验证码操作 控制层处理
//
// PATH /
type CaptchaController struct {
	sysConfigService *systemService.SysConfig // 参数配置服务
}

// Image 获取验证码-图片
//
// GET /captcha-image
func (s CaptchaController) Image(c *gin.Context) {
	// 从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captchaEnabled")
	captchaEnabled := parse.Boolean(captchaEnabledStr)
	if !captchaEnabled {
		c.JSON(200, resp.OkData(map[string]any{
			"enabled": captchaEnabled,
		}))
		return
	}

	// 生成唯一标识
	verifyKey := ""
	data := map[string]any{
		"enabled": captchaEnabled,
		"uuid":    "",
		"img":     "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7",
	}

	// 验证码有效期，单位秒
	captchaExpiration := 2 * 60 * time.Second
	// 从数据库配置获取验证码类型 math 数值计算 char 字符验证
	captchaType := s.sysConfigService.FindValueByKey("sys.account.captchaType")
	if captchaType == "math" {
		driverCaptcha := &base64Captcha.DriverMath{
			//Height png height in pixel.
			Height: config.GetIntOr("mathCaptcha.height", 0),
			// Width Captcha png width in pixel.
			Width: config.GetIntOr("mathCaptcha.width", 0),
			//NoiseCount text noise count.
			NoiseCount: config.GetIntOr("mathCaptcha.noise", 0),
			//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
			ShowLineOptions: base64Captcha.OptionShowHollowLine,
		}
		if config.GetBoolOr("mathCaptcha.color", false) {
			//BgColor captcha image background color (optional)
			driverCaptcha.BgColor = parse.Color(config.GetStringOr("mathCaptcha.background", ""))
		}
		// 验证码生成
		id, question, answer := driverCaptcha.GenerateIdQuestionAnswer()
		// 验证码表达式解析输出
		item, err := driverCaptcha.DrawCaptcha(question)
		if err != nil {
			logger.Infof("generate id question answer %s %s : %v", captchaType, question, err)
		} else {
			data["uuid"] = id
			data["img"] = item.EncodeB64string()
			verifyKey = constants.CACHE_CAPTCHA_CODE + ":" + id
			_ = cache.Set("", verifyKey, answer, captchaExpiration)
		}
	}
	if captchaType == "char" {
		driverCaptcha := &base64Captcha.DriverString{
			//Height png height in pixel.
			Height: config.GetIntOr("charCaptcha.height", 0),
			// Width Captcha png width in pixel.
			Width: config.GetIntOr("charCaptcha.width", 0),
			//NoiseCount text noise count.
			NoiseCount: config.GetIntOr("charCaptcha.noise", 0),
			//Length random string length.
			Length: config.GetIntOr("charCaptcha.size", 0),
			//Source is a unicode which is the rand string from.
			Source: config.GetStringOr("charCaptcha.chars", ""),
			//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
			ShowLineOptions: base64Captcha.OptionShowHollowLine,
		}
		if config.GetBoolOr("charCaptcha.color", false) {
			//BgColor captcha image background color (optional)
			driverCaptcha.BgColor = parse.Color(config.GetStringOr("charCaptcha.background", ""))
		}
		// 验证码生成
		id, question, answer := driverCaptcha.GenerateIdQuestionAnswer()
		// 验证码表达式解析输出
		item, err := driverCaptcha.DrawCaptcha(question)
		if err != nil {
			logger.Infof("generate id question answer %s %s : %v", captchaType, question, err)
		} else {
			data["uuid"] = id
			data["img"] = item.EncodeB64string()
			verifyKey = constants.CACHE_CAPTCHA_CODE + ":" + id
			_ = cache.Set("", verifyKey, strings.ToLower(answer), captchaExpiration)
		}
	}

	// 本地开发下返回验证码结果，方便接口调试
	if config.Env() == "local" {
		text, _ := cache.Get("", verifyKey)
		data["text"] = text
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.OkData(data))
}
